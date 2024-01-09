package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/MWT-proger/shortener/internal/shortener/auth"
)

// GetListUserURLsHandler Возвращает список URL-адресов пользователя.
func (h *APIHandler) GetListUserURLsHandler(w http.ResponseWriter, r *http.Request) {
	var (
		ctx        = r.Context()
		userID, ok = auth.UserIDFrom(ctx)
	)

	if !ok {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	listURLs, err := h.shortService.GetListUserURLs(ctx, userID, r.Host)

	if err != nil {
		h.setHTTPError(w, err)
		return
	}

	resp, err := json.Marshal(listURLs)

	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)

}

// DeleteListUserURLsHandler  в теле запроса принимает
// список идентификаторов сокращённых URL для асинхронного удаления.
// В случае успешного приёма запроса возвращает HTTP-статус 202 Accepted.
func (h *APIHandler) DeleteListUserURLsHandler(w http.ResponseWriter, r *http.Request) {

	var (
		ctx        = r.Context()
		userID, ok = auth.UserIDFrom(ctx)
		data       []string
	)

	if !ok {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()

	if err := h.unmarshalBody(r.Body, &data); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	h.shortService.DeleteListUserURLs(ctx, userID, data)

	w.WriteHeader(http.StatusAccepted)

}
