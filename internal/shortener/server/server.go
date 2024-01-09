package server

import (
	"net/http"

	"github.com/go-chi/chi"

	"github.com/MWT-proger/shortener/configs"
	"github.com/MWT-proger/shortener/internal/shortener/logger"
)

// Run() запускает сервер и слушает его по указанному хосту
func Run(r *chi.Mux) error {
	conf := configs.GetConfig()

	logger.Log.Info("Running server on", logger.StringField("host", conf.HostServer))
	return http.ListenAndServe(conf.HostServer, r)
}
