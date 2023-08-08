package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"github.com/MWT-proger/shortener/internal/shortener/models"
	"github.com/MWT-proger/shortener/internal/shortener/request"
	"github.com/MWT-proger/shortener/internal/shortener/utils"
)

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
