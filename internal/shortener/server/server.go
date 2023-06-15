package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/MWT-proger/shortener/configs"
)

// Run() запускает сервер и слушает его по указанному хосту
func Run(r *chi.Mux) error {
	conf := configs.GetConfig()

	fmt.Println("Running server on", conf.HostServer)
	return http.ListenAndServe(conf.HostServer, r)
}
