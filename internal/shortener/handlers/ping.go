package handlers

import (
	"net/http"
)

// PingDB(w http.ResponseWriter, r *http.Request) Проверяет соединение с базой данных.
func (h *APIHandler) PingDB(w http.ResponseWriter, r *http.Request) {
	if err := h.storage.Ping(); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
