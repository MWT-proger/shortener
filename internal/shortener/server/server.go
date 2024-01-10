package server

import (
	"net/http"

	"github.com/MWT-proger/shortener/configs"
	"github.com/MWT-proger/shortener/internal/shortener/logger"
)

// Run() запускает сервер и слушает его по указанному хосту.
func Run(h Handler, conf configs.Config) error {
	r := initRouter(conf, h)

	logger.Log.Info("Running server on", logger.StringField("host", conf.HostServer))
	return http.ListenAndServe(conf.HostServer, r)
}
