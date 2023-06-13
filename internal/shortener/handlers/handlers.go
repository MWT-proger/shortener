package handlers

import (
	"io"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/MWT-proger/shortener/configs"
	"github.com/MWT-proger/shortener/internal/shortener/storage"
)

type APIHandler struct {
	storage storage.OperationStorager
}

func NewAPIHandler(s storage.OperationStorager) (h *APIHandler, err error) {
	return &APIHandler{s}, err
}

// Принимает большой URL и возвращает маленький
func (h *APIHandler) GenerateShortkeyHandler(w http.ResponseWriter, r *http.Request) {
	var shortURL string
	conf := configs.GetConfig()

	defer r.Body.Close()
	requestData, err := io.ReadAll((r.Body))

	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	stringRequestData := string(requestData)

	if stringRequestData == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	shortURL, err = h.storage.Set(stringRequestData)

	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(http.StatusCreated)

	if conf.BaseURLShortener != "" {
		w.Write([]byte(conf.BaseURLShortener + "/" + shortURL))
	} else {
		w.Write([]byte("http://" + r.Host + "/" + shortURL))
	}

}

// Возвращает по ключу длинный URL
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
