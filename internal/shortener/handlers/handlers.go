package handlers

import (
	"io"
	"log"
	"net/http"

	"github.com/MWT-proger/shortener/internal/shortener/storage"
	"github.com/go-chi/chi"
)

type APIHandler struct {
	storage storage.OperationStorage
}

func NewAPIHandler(s storage.OperationStorage) (h *APIHandler, err error) {
	return &APIHandler{s}, err
}

func (h *APIHandler) GenerateShortkeyHandler(res http.ResponseWriter, req *http.Request) {
	var shortURL string

	defer req.Body.Close()
	requestData, err := io.ReadAll((req.Body))

	if err != nil {
		log.Fatal(err)
	}

	stringRequestData := string(requestData)
	if stringRequestData != "" {
		shortURL, err = h.storage.SetInStorage(stringRequestData)
		if err != nil {
			log.Fatal(err)
		}

		res.Header().Set("content-type", "text/plain")
		res.WriteHeader(http.StatusCreated)

		res.Write([]byte("http://" + req.Host + "/" + shortURL))
		return
	}

	http.Error(res, "Bad Request", http.StatusBadRequest)
}

func (h *APIHandler) GetURLByKeyHandler(w http.ResponseWriter, r *http.Request) {

	fullURL, err := h.storage.GetFromStorage(chi.URLParam(r, "shortKey"))
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}

	if fullURL != "" {
		w.Header().Set("Location", fullURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	http.Error(w, "Bad Request", http.StatusBadRequest)
}
