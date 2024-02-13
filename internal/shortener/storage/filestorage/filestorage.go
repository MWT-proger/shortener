package filestorage

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/google/uuid"

	"github.com/MWT-proger/shortener/configs"
	"github.com/MWT-proger/shortener/internal/shortener/errors"
	"github.com/MWT-proger/shortener/internal/shortener/logger"
	"github.com/MWT-proger/shortener/internal/shortener/models"
	"github.com/MWT-proger/shortener/internal/shortener/utils"
)

// fileStorage хранит данные в json файле.
type fileStorage struct {
	tempStorage map[string]string
}

// NewFileStorage создаёт и возвращает новый экземпляр fileStorage.
func NewFileStorage(ctx context.Context, conf configs.Config) (*fileStorage, error) {

	s := &fileStorage{}

	s.tempStorage = make(map[string]string, 0)
	if conf.JSONFileDB != "" {

		content, err := os.ReadFile(conf.JSONFileDB)

		if err != nil {
			str := "{}"
			if err = os.WriteFile(conf.JSONFileDB, []byte(str), 0644); err != nil {
				return nil, err
			}

		} else {

			if err = json.Unmarshal(content, &s.tempStorage); err != nil {
				return nil, err

			}
		}

		go s.workerBackupToJSONFile(ctx, conf)

	}
	return s, nil
}

// workerBackupToJSONFile() Запускается в gorutine и делаем бэкап по истечению таймера.
func (s *fileStorage) workerBackupToJSONFile(ctx context.Context, conf configs.Config) error {
	ticker := time.NewTicker(conf.TimebackupToJSONFile)

	for {
		select {
		case <-ctx.Done():
			if err := s.backupToJSONFile(ctx, conf); err != nil {
				return err
			}
			return nil

		case <-ticker.C:
			if err := s.backupToJSONFile(ctx, conf); err != nil {
				return err
			}
		}
	}

}

// backupToJSONFile() Делает резервное копирование переменной в файл.
func (s *fileStorage) backupToJSONFile(ctx context.Context, conf configs.Config) error {
	logger.Log.Info("старт - резервное копирование данных -> file.json")

	b, err := json.Marshal(s.tempStorage)

	if err != nil {
		return err
	}

	os.WriteFile(conf.JSONFileDB, b, 0644)
	logger.Log.Info("финиш - резервное копирование данных -> file.json")

	return nil
}

// Set Добавляет в хранилище полную ссылку и присваевает ей ключ.
func (s *fileStorage) Set(ctx context.Context, newModel models.ShortURL) (string, error) {

	newModel.ShortKey = utils.StringWithCharset(5)

	for {
		_, ok := s.tempStorage[newModel.ShortKey]
		if !ok {
			s.tempStorage[newModel.ShortKey] = newModel.FullURL
			break
		}
		newModel.ShortKey = utils.StringWithCharset(5)
	}

	return newModel.ShortKey, nil

}

// Get Достаёт из хранилища и возвращает полную ссылку по ключу.
func (s *fileStorage) Get(ctx context.Context, shortURL string) (models.ShortURL, error) {
	var model models.ShortURL

	fullURL, ok := s.tempStorage[shortURL]
	if !ok {
		return model, nil
	}

	model.FullURL = fullURL

	return model, nil
}

// SetMany Добавляет в хранилище полную ссылку и присваевает ей ключ.
func (s *fileStorage) SetMany(ctx context.Context, data []models.JSONShortURL, baseShortURL string, userID uuid.UUID) error {

	for i, v := range data {

		shortKey := utils.StringWithCharset(5)

		for {
			_, ok := s.tempStorage[shortKey]
			if !ok {
				s.tempStorage[shortKey] = v.OriginalURL
				break
			}
			shortKey = utils.StringWithCharset(5)
		}

		data[i].ShortURL = baseShortURL + shortKey
	}

	return nil

}

// DeleteList Удаляет значение по ключам.
func (s *fileStorage) DeleteList(ctx context.Context, data ...models.DeletedShortURL) error {

	for _, v := range data {

		_, ok := s.tempStorage[v.Payload]

		if ok {
			delete(s.tempStorage, v.Payload)
		}

	}

	return nil
}

// Абстрактный метод.
func (s fileStorage) GetList(ctx context.Context, userID uuid.UUID) ([]*models.JSONShortURL, error) {
	return []*models.JSONShortURL{}, nil
}

// Абстрактный метод.
func (s fileStorage) Ping() error {
	return errors.ErrorDBNotConnection
}

// Абстрактный метод.
func (s fileStorage) Close() error {
	return nil
}
