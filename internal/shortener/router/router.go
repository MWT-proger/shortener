package router

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/MWT-proger/shortener/internal/shortener/auth"
	"github.com/MWT-proger/shortener/internal/shortener/gzip"
	"github.com/MWT-proger/shortener/internal/shortener/handlers"
	"github.com/MWT-proger/shortener/internal/shortener/logger"
)

// Router() Перенаправляет запросы на необходимые хендлеры
func Router(h *handlers.APIHandler) *chi.Mux {

	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(logger.RequestLogger)
		r.Use(gzip.GzipMiddleware)

		r.Use(auth.AuthCookieMiddleware)

		r.Post("/", h.GenerateShortkeyHandler)
		r.Get("/{shortKey}", h.GetURLByKeyHandler)

		r.Get("/api/user/urls", h.GetListUserURLsHandler)
		r.Delete("/api/user/urls", h.DeleteListUserURLsHandler)

		r.Post("/api/shorten", h.JSONGenerateShortkeyHandler)
		r.Post("/api/shorten/batch", h.JSONMultyGenerateShortkeyHandler)

		r.Get("/ping", h.PingDB)
	})

	r.Mount("/debug", middleware.Profiler())

	return r
}
