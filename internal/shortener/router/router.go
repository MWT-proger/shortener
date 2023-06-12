package router

import (
	"github.com/go-chi/chi"

	"github.com/MWT-proger/shortener/internal/shortener/handlers"
	"github.com/MWT-proger/shortener/internal/shortener/storage"
)

// Router() Перенаправляет запросы на необходимые хендлеры
func Router() *chi.Mux {

	r := chi.NewRouter()
	h, _ := handlers.NewAPIHandler(&storage.Storage{})

	r.Post("/", h.GenerateShortkeyHandler)
	r.Get("/{shortKey}", h.GetURLByKeyHandler)

	return r
}
