package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	lErrors "github.com/MWT-proger/shortener/internal/shortener/errors"
	"github.com/MWT-proger/shortener/internal/shortener/models"
	"github.com/MWT-proger/shortener/internal/shortener/request"
	"github.com/MWT-proger/shortener/internal/shortener/utils"
)

type JSONShortenResponse struct {
	Result string `json:"result"`
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
	newModel := models.ShortURL{FullURL: stringRequestData}

	userID, _ := request.UserIDFrom(r.Context())
	newModel.UserID = userID
	shortURL, err = h.storage.Set(newModel)

	if err != nil {

		if !errors.Is(err, &lErrors.ErrorDuplicateFullURL{}) {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusConflict)
		isConflict = true

	}

	if !isConflict {
		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusCreated)
	}
	w.Write([]byte(utils.GetBaseShortURL(r.Host) + shortURL))

}

// JSONGenerateShortkeyHandler Принимает в теле запроса JSON-объект {"url":"<some_url>"}
// и возвращает в ответ объект {"result":"<short_url>"}
func (h *APIHandler) JSONGenerateShortkeyHandler(w http.ResponseWriter, r *http.Request) {

	var (
		data         models.JSONShortenRequest
		responseData JSONShortenResponse
		isConflict   bool
	)

	defer r.Body.Close()

	if err := h.unmarshalBody(r.Body, &data); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if ok := data.IsValid(); !ok {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	newModel := models.ShortURL{FullURL: data.URL}

	userID, ok := request.UserIDFrom(r.Context())
	if ok {
		newModel.UserID = userID
	}

	shortURL, err := h.storage.Set(newModel)

	if err != nil {

		if !errors.Is(err, &lErrors.ErrorDuplicateFullURL{}) {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		isConflict = true

	}

	responseData.Result = utils.GetBaseShortURL(r.Host) + shortURL

	resp, err := json.Marshal(responseData)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if !isConflict {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
	}
	w.Write(resp)

}
