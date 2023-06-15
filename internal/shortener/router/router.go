package router

import (
	"github.com/go-chi/chi"

	"github.com/MWT-proger/shortener/internal/shortener/handlers"
)

// Router() Перенаправляет запросы на необходимые хендлеры
func Router(h *handlers.APIHandler) *chi.Mux {

	r := chi.NewRouter()

	r.Post("/", h.GenerateShortkeyHandler)
	r.Get("/{shortKey}", h.GetURLByKeyHandler)

	return r
}
