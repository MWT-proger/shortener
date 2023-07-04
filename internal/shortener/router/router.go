package router

import (
	"github.com/go-chi/chi"

	"github.com/MWT-proger/shortener/internal/shortener/gzip"
	"github.com/MWT-proger/shortener/internal/shortener/handlers"
	"github.com/MWT-proger/shortener/internal/shortener/logger"
)

// Router() Перенаправляет запросы на необходимые хендлеры
func Router(h *handlers.APIHandler) *chi.Mux {

	r := chi.NewRouter()
	r.Use(logger.RequestLogger)
	r.Use(gzip.GzipMiddleware)
	r.Post("/", h.GenerateShortkeyHandler)
	r.Get("/{shortKey}", h.GetURLByKeyHandler)
	r.Post("/api/shorten", h.JSONGenerateShortkeyHandler)

	return r
}
