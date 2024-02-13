package services

//go:generate mockgen -destination mock_test.go -package services github.com/MWT-proger/shortener/internal/shortener/services Storager

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/MWT-proger/shortener/configs"
	lErrors "github.com/MWT-proger/shortener/internal/shortener/errors"
	"github.com/MWT-proger/shortener/internal/shortener/logger"
	"github.com/MWT-proger/shortener/internal/shortener/models"
	"github.com/MWT-proger/shortener/internal/shortener/utils"
)

// Storager интерфейс хранилища.
type Storager interface {
	Set(ctx context.Context, newModel models.ShortURL) (string, error)
	SetMany(ctx context.Context, data []models.JSONShortURL, baseShortURL string, userID uuid.UUID) error
	Get(ctx context.Context, shortURL string) (models.ShortURL, error)
	GetList(ctx context.Context, userID uuid.UUID) ([]*models.JSONShortURL, error)
	DeleteList(ctx context.Context, data ...models.DeletedShortURL) error
	CountUsersAndUrls(ctx context.Context) (urls int, users int, err error)

	Ping() error
}

// ShortenerService сервис обработки Full and Short URLs.
type ShortenerService struct {
	conf        configs.Config
	storage     Storager
	deletedChan chan models.DeletedShortURL
	doneCh      chan struct{}
}

// NewShortenerService - создаёт новый экземпляр сервиса обработки Full and Short URLs.
func NewShortenerService(ctx context.Context, conf configs.Config, s Storager) *ShortenerService {

	ss := &ShortenerService{
		storage:     s,
		conf:        conf,
		deletedChan: make(chan models.DeletedShortURL, 1024),
		doneCh:      make(chan struct{}),
	}

	go ss.flushDeleted(ctx)

	return ss
}

// GenerateShortURL Принимает большой URL и возвращает маленький.
func (s *ShortenerService) GenerateShortURL(ctx context.Context, userID uuid.UUID, fullURL string, requestHost string) (string, error) {
	var responseErr error
	data := models.ShortURL{FullURL: fullURL}
	data.UserID = userID
	shortKey, err := s.storage.Set(ctx, data)

	if err != nil {

		if !errors.Is(err, lErrors.ErrorDuplicateFullURLServicesError) {
			return "", lErrors.InternalServicesError
		}
		responseErr = lErrors.ErrorDuplicateFullURLServicesError

	}
	shortURL := utils.GetBaseShortURL(s.conf, requestHost) + shortKey

	return shortURL, responseErr

}

// GenerateMultyShortURL Принимает  []models.JSONShortURL и  добавлет в каждый объект сокращенный URL.
func (s *ShortenerService) GenerateMultyShortURL(ctx context.Context, userID uuid.UUID, data []models.JSONShortURL, requestHost string) error {

	baseShortURL := utils.GetBaseShortURL(s.conf, requestHost)

	err := s.storage.SetMany(ctx, data, baseShortURL, userID)

	if err != nil {
		return lErrors.InternalServicesError
	}

	return nil
}

// GetFullURLByShortKey Возвращает полный URL по переданному ключу.
func (s *ShortenerService) GetFullURLByShortKey(ctx context.Context, shortKey string) (string, error) {

	data, err := s.storage.Get(ctx, shortKey)

	if err != nil {
		return "", lErrors.GetFullURLServicesError
	}

	if data.DeletedFlag {
		return "", lErrors.GoneServicesError
	}

	if data.FullURL == "" {
		return "", lErrors.NotFoundServicesError
	}

	return data.FullURL, nil

}

// GetListUserURLs Возвращает список URL-адресов пользователя.
func (s *ShortenerService) GetListUserURLs(ctx context.Context, userID uuid.UUID, requestHost string) ([]*models.JSONShortURL, error) {

	listURLs, err := s.storage.GetList(ctx, userID)

	if err != nil {
		return nil, lErrors.GetFullURLServicesError
	}

	if len(listURLs) == 0 {
		return nil, lErrors.NoContentUserServicesError
	}

	baseURL := utils.GetBaseShortURL(s.conf, requestHost)

	for _, v := range listURLs {
		v.ShortURL = baseURL + v.ShortURL
	}

	return listURLs, nil

}

// DeleteListUserURLs принимает список идентификаторов сокращённых URL для асинхронного удаления.
func (s *ShortenerService) DeleteListUserURLs(ctx context.Context, userID uuid.UUID, data []string) {

	go func() {
		for _, d := range data {
			select {
			case <-s.doneCh:
				return
			case s.deletedChan <- models.DeletedShortURL{
				UserID:  userID,
				Payload: d,
			}:
			}
		}
	}()

}

// GetStats возвращает количество пользователей и сокращенных URL в сервисе.
func (s *ShortenerService) GetStats(ctx context.Context) (urls int, users int, err error) {

	urls, users, err = s.storage.CountUsersAndUrls(ctx)

	if err != nil {
		return 0, 0, lErrors.InternalServicesError
	}

	return urls, users, nil

}

// flushDeleted запускается в горутине и удаляет ссылки.
func (s *ShortenerService) flushDeleted(ctx context.Context) {
	// будем удалять, накопленные за последние 10 секунд
	ticker := time.NewTicker(10 * time.Second)

	var data []models.DeletedShortURL

	for {
		select {
		case <-ctx.Done():

			for len(s.deletedChan) > 0 {
				d := <-s.deletedChan
				data = append(data, d)
			}

			if len(data) != 0 {
				if err := s.storage.DeleteList(ctx, data...); err != nil {
					logger.Log.Debug("cannot deleted shortURL", logger.ErrorField(err))
				}
			}
			logger.Log.Info("Завершение удаления ссылок - ОК")
			return

		case d := <-s.deletedChan:
			data = append(data, d)
		case <-ticker.C:
			if len(data) == 0 {
				continue
			}
			err := s.storage.DeleteList(ctx, data...)
			if err != nil {
				logger.Log.Debug("cannot deleted shortURL", logger.ErrorField(err))
				continue
			}
			data = nil
		}
	}
}
