package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/MWT-proger/shortener/internal/shortener/logger"
	"github.com/MWT-proger/shortener/internal/shortener/models"
	"github.com/MWT-proger/shortener/internal/shortener/request"
	"github.com/MWT-proger/shortener/internal/shortener/storage"
	"github.com/MWT-proger/shortener/internal/shortener/utils"
)

type APIHandler struct {
	storage     storage.OperationStorager
	DeletedChan chan models.DeletedShortURL
	doneCh      chan struct{}
}

func NewAPIHandler(s storage.OperationStorager) (h *APIHandler, err error) {
	hh := &APIHandler{
		storage:     s,
		DeletedChan: make(chan models.DeletedShortURL, 1024),
		doneCh:      make(chan struct{}),
	}

	go hh.FlushDeleted()

	return hh, err
}

// GetURLByKeyHandler Возвращает по ключу длинный URL
func (h *APIHandler) GetURLByKeyHandler(w http.ResponseWriter, r *http.Request) {

	modelData, err := h.storage.Get(chi.URLParam(r, "shortKey"))

	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	if modelData.DeletedFlag {
		w.WriteHeader(http.StatusGone)
		return
	}
	if modelData.FullURL == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", modelData.FullURL)
	w.WriteHeader(http.StatusTemporaryRedirect)

}

// GetURLByKeyHandler Возвращает список URL-адресов пользователя
func (h *APIHandler) GetListUserURLsHandler(w http.ResponseWriter, r *http.Request) {

	userID, ok := request.UserIDFrom(r.Context())

	if !ok {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	listURLs, err := h.storage.GetList(userID)

	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	if len(listURLs) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	baseURL := utils.GetBaseShortURL(r.Host)
	for _, v := range listURLs {
		v.ShortURL = baseURL + v.ShortURL
	}

	resp, err := json.Marshal(listURLs)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)

}

// DeleteListUserURLsHandler  в теле запроса принимает
// список идентификаторов сокращённых URL для асинхронного удаления
// В случае успешного приёма запроса возвращает HTTP-статус 202 Accepted
func (h *APIHandler) DeleteListUserURLsHandler(w http.ResponseWriter, r *http.Request) {

	userID, ok := request.UserIDFrom(r.Context())

	if !ok {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	var data []string

	defer r.Body.Close()

	if err := h.unmarshalBody(r.Body, &data); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	inputCh := generator(h.doneCh, data, userID)
	channels := fanOut(h.doneCh, inputCh)
	fanIn(h.doneCh, h.DeletedChan, channels...)

	w.WriteHeader(http.StatusAccepted)

}

// generator функция из предыдущего примера, делает то же, что и делала
func generator(doneCh chan struct{}, input []string, userID uuid.UUID) chan models.DeletedShortURL {
	inputCh := make(chan models.DeletedShortURL)

	go func() {
		defer close(inputCh)

		for _, data := range input {
			select {
			case <-doneCh:
				return
			case inputCh <- models.DeletedShortURL{
				UserID:  userID,
				Payload: data,
			}:
			}
		}
	}()

	return inputCh
}

func (h *APIHandler) FlushDeleted() {
	// будем удалять, накопленные за последние 10 секунд
	ticker := time.NewTicker(10 * time.Second)

	var data []models.DeletedShortURL

	for {
		select {
		case d := <-h.DeletedChan:
			data = append(data, d)
		case <-ticker.C:
			if len(data) == 0 {
				continue
			}
			err := h.storage.DeleteList(data...)
			if err != nil {
				logger.Log.Debug("cannot deleted shortURL", zap.Error(err))
				continue
			}
			data = nil
		}
	}
}

// fanOut принимает канал данных, порождает 10 горутин
func fanOut(doneCh chan struct{}, inputCh chan models.DeletedShortURL) []chan models.DeletedShortURL {
	numWorkers := 10
	channels := make([]chan models.DeletedShortURL, numWorkers)

	for i := 0; i < numWorkers; i++ {
		channels[i] = inputCh
	}

	return channels
}

// fanIn объединяет несколько каналов resultChs в один.
func fanIn(doneCh chan struct{}, finalCh chan models.DeletedShortURL, resultChs ...chan models.DeletedShortURL) {

	for _, ch := range resultChs {
		chClosure := ch

		go func() {

			for data := range chClosure {
				select {
				case <-doneCh:
					return
				case finalCh <- data:
				}
			}
		}()
	}
}
