package handlers

import (
	"net/http"
)

// PingDB Проверяет соединение с базой данных.
func (h *APIHandler) PingDB(w http.ResponseWriter, r *http.Request) {

	if ok := h.shortService.PingStorage(); !ok {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
