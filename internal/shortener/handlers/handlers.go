package handlers

import (
	"net/http"

	"github.com/MWT-proger/shortener/internal/shortener/storage"
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

	fullURL, err := h.storage.Get(chi.URLParam(r, "shortKey"))

	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if fullURL == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", fullURL)
	w.WriteHeader(http.StatusTemporaryRedirect)

}
