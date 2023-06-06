package handlers

import (
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/MWT-proger/shortener/internal/shortener/storage"
)

type APIHandler struct {
	storage storage.Storage
}

func NewAPIHandler() (h *APIHandler, err error) {
	return &APIHandler{}, err
}

func (h *APIHandler) BaseHandler(res http.ResponseWriter, req *http.Request) {
	switch req.Method {

	case http.MethodGet:
		h.GetURLByKeyHandler(res, req, &h.storage)
		return

	case http.MethodPost:
		h.GenerateShortkeyHandler(res, req, &h.storage)
		return

	default:
		http.Error(res, "Bad Request", http.StatusBadRequest)
	}
}

func (h *APIHandler) GenerateShortkeyHandler(res http.ResponseWriter, req *http.Request, storg storage.OperationStorage) {
	var shortURL string

	if req.Method == http.MethodPost {
		if req.URL.Path == "/" {

			defer req.Body.Close()
			requestData, err := io.ReadAll((req.Body))

			if err != nil {
				log.Fatal(err)
			}

			stringRequestData := string(requestData)
			if stringRequestData != "" {
				shortURL, err = storg.SetInStorage(stringRequestData)
				if err != nil {
					log.Fatal(err)
				}

				res.Header().Set("content-type", "text/plain")
				res.WriteHeader(http.StatusCreated)

				res.Write([]byte("http://" + req.Host + "/" + shortURL))
				return
			}

		}
	}

	http.Error(res, "Bad Request", http.StatusBadRequest)
}

func (h *APIHandler) GetURLByKeyHandler(w http.ResponseWriter, r *http.Request, s storage.OperationStorage) {

	if r.Method == http.MethodGet {
		if idx := strings.Index(r.URL.Path[1:], "/"); idx < 0 && r.URL.Path != "/" {

			fullURL, err := s.GetFromStorage(r.URL.Path[1:])
			if err != nil {
				http.Error(w, "Bad Request", http.StatusBadRequest)
			}

			if fullURL != "" {
				w.Header().Set("Location", fullURL)
				w.WriteHeader(http.StatusTemporaryRedirect)
				return
			}

		}
	}

	http.Error(w, "Bad Request", http.StatusBadRequest)
}
