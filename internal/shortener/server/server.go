package server

import (
	"net/http"

	"golang.org/x/crypto/acme/autocert"

	"github.com/MWT-proger/shortener/configs"
	"github.com/MWT-proger/shortener/internal/shortener/logger"
)

// Run() запускает сервер и слушает его по указанному хосту.
func Run(h Handler, conf configs.Config) error {
	r := initRouter(conf, h)

	logger.Log.Info("Running server on", logger.StringField("host", conf.HostServer))

	if conf.EnableHTTPS {
		// конструируем менеджер TLS-сертификатов
		manager := &autocert.Manager{
			// директория для хранения сертификатов
			Cache: autocert.DirCache("cache-dir"),
			// функция, принимающая Terms of Service издателя сертификатов
			Prompt: autocert.AcceptTOS,
			// перечень доменов, для которых будут поддерживаться сертификаты
			HostPolicy: autocert.HostWhitelist(conf.HostServer),
		}
		// конструируем сервер с поддержкой TLS
		server := &http.Server{
			Addr:    ":443",
			Handler: r,
			// для TLS-конфигурации используем менеджер сертификатов
			TLSConfig: manager.TLSConfig(),
		}
		return server.ListenAndServeTLS("", "")
	}

	return http.ListenAndServe(conf.HostServer, r)
}
