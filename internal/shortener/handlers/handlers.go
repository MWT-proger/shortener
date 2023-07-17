package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/chi"

	lErrors "github.com/MWT-proger/shortener/internal/shortener/errors"
	"github.com/MWT-proger/shortener/internal/shortener/models"
	"github.com/MWT-proger/shortener/internal/shortener/storage"
	"github.com/MWT-proger/shortener/internal/shortener/utils"
)

type APIHandler struct {
	storage storage.OperationStorager
}

func NewAPIHandler(s storage.OperationStorager) (h *APIHandler, err error) {
	return &APIHandler{s}, err
}

// GenerateShortkeyHandler Принимает большой URL и возвращает маленький
func (h *APIHandler) GenerateShortkeyHandler(w http.ResponseWriter, r *http.Request) {
	var shortURL string
	var isConflict bool

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

		if !errors.Is(err, &lErrors.ErrorDuplicateFullURL{}) {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusConflict)
		isConflict = true

	}

	w.Header().Set("content-type", "text/plain")
	if !isConflict {
		w.WriteHeader(http.StatusCreated)
	}
	w.Write([]byte(utils.GetBaseShortURL(r.Host) + shortURL))

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

type JSONShortenRequest struct {
	URL string `json:"url"`
}
type JSONShortenResponse struct {
	Result string `json:"result"`
}

// JSONGenerateShortkeyHandler Принимает в теле запроса JSON-объект {"url":"<some_url>"}
// и возвращает в ответ объект {"result":"<short_url>"}
func (h *APIHandler) JSONGenerateShortkeyHandler(w http.ResponseWriter, r *http.Request) {
	var (
		shortURL     string
		buf          bytes.Buffer
		requestData  JSONShortenRequest
		responseData JSONShortenResponse
		isConflict   bool
	)
	defer r.Body.Close()

	_, err := buf.ReadFrom(r.Body)

	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &requestData); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if requestData.URL == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	shortURL, err = h.storage.Set(requestData.URL)

	if err != nil {
		if !errors.Is(err, &lErrors.ErrorDuplicateFullURL{}) {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusConflict)
		isConflict = true

	}

	responseData.Result = utils.GetBaseShortURL(r.Host) + shortURL

	resp, err := json.Marshal(responseData)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if !isConflict {
		w.WriteHeader(http.StatusCreated)
	}
	w.Write(resp)

}

func (h *APIHandler) unmarshalBody(body io.ReadCloser, form interface{}) error {

	defer body.Close()

	var buf bytes.Buffer
	_, err := buf.ReadFrom(body)

	if err != nil {
		return err
	}

	if err = json.Unmarshal(buf.Bytes(), form); err != nil {
		return err
	}

	return nil
}

// JSONMultyGenerateShortkeyHandler Принимает в теле запроса JSON-объект в виде списка
// и возвращает в ответ объект в виде списка
func (h *APIHandler) JSONMultyGenerateShortkeyHandler(w http.ResponseWriter, r *http.Request) {
	var data []models.JSONShortURL

	defer r.Body.Close()

	if err := h.unmarshalBody(r.Body, &data); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	for _, v := range data {

		if ok := v.IsValid(); !ok {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
	}
	err := h.storage.SetMany(data, utils.GetBaseShortURL(r.Host))

	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)

}
