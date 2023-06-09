package handlers

import (
	"bytes"
	"encoding/json"
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

// GenerateShortkeyHandler Принимает большой URL и возвращает маленький
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
		conf         = configs.GetConfig()
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
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if conf.BaseURLShortener != "" {
		responseData.Result = conf.BaseURLShortener + "/" + shortURL
	} else {
		responseData.Result = "http://" + r.Host + "/" + shortURL
	}

	resp, err := json.Marshal(responseData)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)

}
