package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/MWT-proger/shortener/internal/shortener/logger"
	"github.com/MWT-proger/shortener/internal/shortener/request"
	"github.com/MWT-proger/shortener/internal/shortener/storage"
	"github.com/MWT-proger/shortener/internal/shortener/utils"
	"github.com/go-chi/chi"
)

type APIHandler struct {
	storage storage.OperationStorager
}

func NewAPIHandler(s storage.OperationStorager) (h *APIHandler, err error) {
	return &APIHandler{s}, err
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

	go func() {
		h.storage.DeleteList(data, userID)
		logger.Log.Debug("Удаление успешно завершено")
	}()

	w.WriteHeader(http.StatusAccepted)

}
