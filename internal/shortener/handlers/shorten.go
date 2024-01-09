package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/MWT-proger/shortener/internal/shortener/auth"
	"github.com/MWT-proger/shortener/internal/shortener/models"
)

// JSONGenerateShortkeyHandler Принимает в теле запроса JSON-объект {"url":"<some_url>"}
// и возвращает в ответ объект {"result":"<short_url>"}.
func (h *APIHandler) JSONGenerateShortkeyHandler(w http.ResponseWriter, r *http.Request) {

	var (
		data            jsonShortenRequest
		responseData    jsonShortenResponse
		finalStatusCode = http.StatusCreated
		ctx             = r.Context()
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

	userID, ok := auth.UserIDFrom(ctx)

	if !ok {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	shortURL, err := h.shortService.GenerateShortURL(ctx, userID, data.URL, r.Host)

	if err != nil {
		finalStatusCode = h.setOrGetHTTPCode(w, err)
		if finalStatusCode == 0 {
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
	var (
		data []models.JSONShortURL
		ctx  = r.Context()
	)

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

	userID, ok := auth.UserIDFrom(ctx)

	if !ok {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if err := h.shortService.GenerateMultyShortURL(ctx, userID, data, r.Host); err != nil {
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

// jsonShortenResponse - тело ответа для JSONGenerateShortkeyHandler.
type jsonShortenResponse struct {
	Result string `json:"result"`
}

// jsonShortenRequest - тело запроса для JSONGenerateShortkeyHandler.
type jsonShortenRequest struct {
	URL string `json:"url"`
}

// IsValid проверяет на валидность JSONShortenRequest.
func (d *jsonShortenRequest) IsValid() bool {
	return d.URL != ""
}
