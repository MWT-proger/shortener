package handlers

import (
	"encoding/json"
	"net/http"
)

// GetStats возвращает количество пользователей и сокращенных URL в сервисе.
func (h *APIHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	var (
		responseData    jsonStatsResponse
		finalStatusCode = http.StatusOK
		ctx             = r.Context()
	)

	urls, users, err := h.shortService.GetStats(ctx)

	if err != nil {
		finalStatusCode = h.setOrGetHTTPCode(w, err)
		if finalStatusCode == 0 {
			return
		}
	}

	responseData.URLs = urls
	responseData.Users = users

	resp, err := json.Marshal(responseData)

	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(finalStatusCode)
	w.Write(resp)

}

// jsonShortenResponse - тело ответа для GetStats.
type jsonStatsResponse struct {
	URLs  int `json:"urls"`
	Users int `json:"users"`
}
