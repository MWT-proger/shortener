package handlers

import (
	"io"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"

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
	GenerateShortURL(userID uuid.UUID, fullURL string, requestHost string) (string, error)
}

// GenerateShortkeyHandler Принимает большой URL и возвращает маленький
func (h *APIHandler) GenerateShortkeyHandler(w http.ResponseWriter, r *http.Request) {

	var finalStatusCode = http.StatusCreated

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

	userID, ok := request.UserIDFrom(r.Context())

	if !ok {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	shortURL, err := h.shortService.GenerateShortURL(userID, stringRequestData, r.Host)

	if err != nil {
		finalStatusCode = h.setOrGetHTTPCode(w, err)

		if finalStatusCode == 0 {
			return
		}
	}

	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(finalStatusCode)
	w.Write([]byte(shortURL))

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
