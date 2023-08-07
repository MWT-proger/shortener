package router

import (
	"github.com/go-chi/chi"

	"github.com/MWT-proger/shortener/internal/shortener/auth"
	"github.com/MWT-proger/shortener/internal/shortener/gzip"
	"github.com/MWT-proger/shortener/internal/shortener/handlers"
	"github.com/MWT-proger/shortener/internal/shortener/logger"
)

// Router() Перенаправляет запросы на необходимые хендлеры
func Router(h *handlers.APIHandler) *chi.Mux {

	r := chi.NewRouter()

	r.Use(logger.RequestLogger)
	r.Use(gzip.GzipMiddleware)

	r.Group(func(r chi.Router) {
		r.Use(auth.AuthCookieMiddleware)

		r.Post("/", h.GenerateShortkeyHandler)
		r.Get("/{shortKey}", h.GetURLByKeyHandler)

		r.Get("/ping", h.PingDB)

		r.Get("/api/user/urls", h.GetListUserURLsHandler)
		r.Delete("/api/user/urls", h.DeleteListUserURLsHandler)
		r.Post("/api/shorten", h.JSONGenerateShortkeyHandler)
		r.Post("/api/shorten/batch", h.JSONMultyGenerateShortkeyHandler)
	})

	r.Group(func(r chi.Router) {
		r.Use(auth.ValidAuthCookieMiddleware)
		r.Get("/api/user/urls", h.GetListUserURLsHandler)
	})

	return r
}
