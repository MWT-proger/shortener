package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/MWT-proger/shortener/internal/shortener/logger"
	"github.com/MWT-proger/shortener/internal/shortener/models"
	"github.com/MWT-proger/shortener/internal/shortener/request"
	"github.com/MWT-proger/shortener/internal/shortener/storage"
	"github.com/MWT-proger/shortener/internal/shortener/utils"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
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

	for _, v := range data {
		h.DeletedChan <- models.DeletedShortURL{
			UserID:  userID,
			Payload: v,
		}
	}

	w.WriteHeader(http.StatusAccepted)

}

// // generator функция из предыдущего примера, делает то же, что и делала
// func generator(doneCh chan struct{}, input []string, userID ) chan int {
// 	inputCh := make(chan int)

// 	go func() {
// 		defer close(inputCh)

// 		for _, data := range input {
// 			select {
// 			case <-doneCh:
// 				return
// 			case inputCh <- data:
// 			}
// 		}
// 	}()

// 	return inputCh
// }

func (h *APIHandler) FlushDeleted() {
	// будем удалять, накопленные за последние 10 секунд
	ticker := time.NewTicker(10 * time.Second)

	var data []models.DeletedShortURL

	for {
		select {
		case d := <-h.DeletedChan:
			// добавим сообщение в слайс для последующего сохранения
			data = append(data, d)
		case <-ticker.C:
			// подождём, пока придёт хотя бы одно сообщение
			if len(data) == 0 {
				continue
			}
			fmt.Println("Удаление ____________")
			fmt.Println(data)
			// сохраним все пришедшие сообщения одновременно
			err := h.storage.DeleteList(data...)
			if err != nil {
				logger.Log.Debug("cannot deleted shortURL", zap.Error(err))
				// не будем стирать сообщения, попробуем отправить их чуть позже
				continue
			}
			// сотрём успешно отосланные сообщения
			data = nil
		}
	}
}
