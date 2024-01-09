package handlers

import (
	"context"
	"io"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"

	"github.com/MWT-proger/shortener/internal/shortener/models"
	"github.com/MWT-proger/shortener/internal/shortener/request"
)

// APIHandler Структура объеденяющая все эндпоинты.
type APIHandler struct {
	shortService ShortenerServicer
}

// NewAPIHandler создает новую структуру APIHandler.
func NewAPIHandler(service ShortenerServicer) (h *APIHandler, err error) {
	hh := &APIHandler{
		shortService: service,
	}

	return hh, err
}

// ShortenerServicer интерфейс описывающий необходимые методы для сервисного слоя.
type ShortenerServicer interface {
	GetFullURLByShortKey(ctx context.Context, shortKey string) (string, error)
	GetListUserURLs(ctx context.Context, userID uuid.UUID, requestHost string) ([]*models.JSONShortURL, error)

	GenerateShortURL(ctx context.Context, userID uuid.UUID, fullURL string, requestHost string) (string, error)
	GenerateMultyShortURL(ctx context.Context, userID uuid.UUID, data []models.JSONShortURL, requestHost string) error

	DeleteListUserURLsHandler(ctx context.Context, userID uuid.UUID, data []string)

	PingStorage() bool
}

// GenerateShortkeyHandler Принимает большой URL и возвращает маленький.
func (h *APIHandler) GenerateShortkeyHandler(w http.ResponseWriter, r *http.Request) {

	var (
		finalStatusCode = http.StatusCreated
		ctx             = r.Context()
	)

	defer r.Body.Close()

	requestBody, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	stringRequestBody := string(requestBody)

	if stringRequestBody == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	userID, ok := request.UserIDFrom(ctx)

	if !ok {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	shortURL, err := h.shortService.GenerateShortURL(ctx, userID, stringRequestBody, r.Host)

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

// GetURLByKeyHandler Принимает короткий ключ и делает rederict на полный URL
func (h *APIHandler) GetURLByKeyHandler(w http.ResponseWriter, r *http.Request) {
	var (
		ctx      = r.Context()
		shortKey = chi.URLParam(r, "shortKey")
	)

	fullURL, err := h.shortService.GetFullURLByShortKey(ctx, shortKey)

	if err != nil {
		h.setHTTPError(w, err)
		return
	}

	w.Header().Set("Location", fullURL)
	w.WriteHeader(http.StatusTemporaryRedirect)

}
