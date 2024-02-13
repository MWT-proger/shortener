package server

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/MWT-proger/shortener/configs"
	"github.com/MWT-proger/shortener/internal/shortener/auth"
	"github.com/MWT-proger/shortener/internal/shortener/gzip"
	"github.com/MWT-proger/shortener/internal/shortener/logger"
)

// Handler интерфейс определяет необходимые методы для инициализации маршрутизатора.
type Handler interface {
	DeleteListUserURLsHandler(w http.ResponseWriter, r *http.Request)
	GenerateShortkeyHandler(w http.ResponseWriter, r *http.Request)
	GetListUserURLsHandler(w http.ResponseWriter, r *http.Request)
	GetURLByKeyHandler(w http.ResponseWriter, r *http.Request)
	JSONGenerateShortkeyHandler(w http.ResponseWriter, r *http.Request)
	JSONMultyGenerateShortkeyHandler(w http.ResponseWriter, r *http.Request)
	PingDB(w http.ResponseWriter, r *http.Request)
	GetStats(w http.ResponseWriter, r *http.Request)
}

// initRouter() инициализирует и возвращает маршрутизатор.
func initRouter(conf configs.Config, h Handler) *chi.Mux {

	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(logger.RequestLoggerMiddleware)
		r.Use(gzip.GzipMiddleware)

		r.Group(func(r chi.Router) {
			r.Use(func(next http.Handler) http.Handler {
				return auth.CheckIPIncludedSubNetMiddleware(next, conf)
			})
			r.Get("/api/internal/statsn", h.GetStats)
		})

		r.Use(func(next http.Handler) http.Handler {
			return auth.AuthCookieMiddleware(next, conf)
		})

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
