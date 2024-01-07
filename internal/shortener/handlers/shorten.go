package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"

	lErrors "github.com/MWT-proger/shortener/internal/shortener/errors"
	"github.com/MWT-proger/shortener/internal/shortener/models"
	"github.com/MWT-proger/shortener/internal/shortener/request"
	"github.com/MWT-proger/shortener/internal/shortener/utils"
)

// JSONShortenResponse - тело ответа для JSONGenerateShortkeyHandler
type JSONShortenResponse struct {
	Result string `json:"result"`
}

// JSONGenerateShortkeyHandler Принимает в теле запроса JSON-объект {"url":"<some_url>"}
// и возвращает в ответ объект {"result":"<short_url>"}
func (h *APIHandler) JSONGenerateShortkeyHandler(w http.ResponseWriter, r *http.Request) {

	var (
		data            models.JSONShortenRequest
		responseData    JSONShortenResponse
		finalStatusCode = http.StatusCreated
		serviceError    *lErrors.ServicesError
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

	userID, ok := request.UserIDFrom(r.Context())

	if !ok {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	shortURL, err := h.shortService.GenerateShortURL(userID, data.URL, r.Host)

	if err != nil {

		if errors.As(err, &serviceError) {

			if serviceError.ContentType != "" {
				finalStatusCode = http.StatusConflict
			} else {
				http.Error(w, serviceError.Error(), serviceError.HTTPCode)
				return
			}

		} else {
			http.Error(w, "Ошибка сервера, попробуйте позже.", http.StatusInternalServerError)
			return
		}

	}

	responseData.Result = shortURL

	resp, err := json.Marshal(responseData)

	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(finalStatusCode)
	w.Write(resp)

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
	userID, ok := request.UserIDFrom(r.Context())
	if !ok {
		userID = uuid.Nil
	}
	err := h.storage.SetMany(data, utils.GetBaseShortURL(r.Host), userID)

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
