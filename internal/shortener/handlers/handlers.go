package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"

	lErrors "github.com/MWT-proger/shortener/internal/shortener/errors"
	"github.com/MWT-proger/shortener/internal/shortener/models"
	"github.com/MWT-proger/shortener/internal/shortener/request"
	"github.com/MWT-proger/shortener/internal/shortener/storage"
)

// @Title Shortener API
// @Description Сервис сокращения ссылок.
// @Version 1.0

// @Contact.email support@localhost.ru

// @BasePath /
// @Host localhost:7000

// @SecurityDefinitions.apikey ApiKeyAuth
// @In cookie
// @Name token

// @Tag.name Short
// @Tag.description "API сокращения и получения ссылок"

// APIHandler Структура объеденяющая все эндпоинты
type APIHandler struct {
	storage storage.OperationStorager

	shortService ShortenerServicer
}

// NewAPIHandler
func NewAPIHandler(s storage.OperationStorager, ss ShortenerServicer) (h *APIHandler, err error) {
	hh := &APIHandler{
		storage:      s,
		shortService: ss,
	}

	return hh, err
}

type ShortenerServicer interface {
	GetFullURLByShortKey(shortKey string) (string, error)
	GetListUserURLs(userID uuid.UUID, requestHost string) ([]*models.JSONShortURL, error)
	DeleteListUserURLsHandler(userID uuid.UUID, data []string)
}

// GetURLByKeyHandler godoc
// @Tags Short
// @Summary Получить полный url по ключу
// @ID GetURLByKeyHandler
// @Success 307 {string} string
// @Failure 500 {string} string "Внутренняя ошибка"
// @Router /{shortKey} [get]
func (h *APIHandler) GetURLByKeyHandler(w http.ResponseWriter, r *http.Request) {

	fullURL, err := h.shortService.GetFullURLByShortKey(chi.URLParam(r, "shortKey"))

	if err != nil {
		h.setHTTPError(w, err)
		return
	}

	w.Header().Set("Location", fullURL)
	w.WriteHeader(http.StatusTemporaryRedirect)

}

// GetURLByKeyHandler Возвращает список URL-адресов пользователя
func (h *APIHandler) GetListUserURLsHandler(w http.ResponseWriter, r *http.Request) {

	userID, ok := request.UserIDFrom(r.Context())

	if !ok {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	listURLs, err := h.shortService.GetListUserURLs(userID, r.Host)

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
// список идентификаторов сокращённых URL для асинхронного удаления
// В случае успешного приёма запроса возвращает HTTP-статус 202 Accepted
func (h *APIHandler) DeleteListUserURLsHandler(w http.ResponseWriter, r *http.Request) {

	userID, ok := request.UserIDFrom(r.Context())

	if !ok {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	var data []string

	defer r.Body.Close()

	if err := h.unmarshalBody(r.Body, &data); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	h.shortService.DeleteListUserURLsHandler(userID, data)

	w.WriteHeader(http.StatusAccepted)

}

// setHTTPError(w http.ResponseWriter, err error) присваивает response статус ответа
// вынесен для исключения дублирования в коде
func (h *APIHandler) setHTTPError(w http.ResponseWriter, err error) {
	var serviceError *lErrors.ServicesError
	if errors.As(err, &serviceError) {
		http.Error(w, serviceError.Error(), serviceError.HTTPCode)
	} else {
		http.Error(w, "Ошибка сервера, попробуйте позже.", http.StatusInternalServerError)
	}
}
