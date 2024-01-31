package server

import (
	"context"
	"net/http"

	"golang.org/x/crypto/acme/autocert"

	"github.com/MWT-proger/shortener/configs"
	"github.com/MWT-proger/shortener/internal/shortener/logger"
)

// Run() запускает сервер и слушает его по указанному хосту.
func Run(ctx context.Context, h Handler, conf configs.Config) error {
	r := initRouter(conf, h)
	server := http.Server{
		Handler: r,
		Addr:    conf.HostServer,
	}
	logger.Log.Info("Running server on", logger.StringField("host", conf.HostServer))

	go func() {
		<-ctx.Done()
		logger.Log.Info("Stopping server on", logger.StringField("host", conf.HostServer))
		if err := server.Shutdown(ctx); err != nil {
			logger.Log.Error("HTTP server Shutdown", logger.ErrorField(err))
		}
	}()

	if conf.EnableHTTPS {
		manager := &autocert.Manager{
			Cache:      autocert.DirCache("cache-dir"),
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(conf.HostServer),
		}
		server.Addr = ":443"
		server.TLSConfig = manager.TLSConfig()
		return server.ListenAndServeTLS("", "")
	}

	return server.ListenAndServe()
}
