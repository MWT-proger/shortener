package router

import (
	"github.com/MWT-proger/shortener/internal/shortener/handlers"
	"github.com/MWT-proger/shortener/internal/shortener/storage"
	"github.com/go-chi/chi"
)

func Router() *chi.Mux {

	r := chi.NewRouter()
	h, _ := handlers.NewAPIHandler(&storage.Storage{})

	r.Post("/", h.GenerateShortkeyHandler)
	r.Get("/{shortKey}", h.GetURLByKeyHandler)

	return r
}
