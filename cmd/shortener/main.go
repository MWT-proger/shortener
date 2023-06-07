package main

import (
	"fmt"
	"net/http"

	"github.com/MWT-proger/shortener/configs"
	"github.com/MWT-proger/shortener/internal/shortener/router"
	"github.com/MWT-proger/shortener/internal/shortener/storage"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	configInit := configs.InitConfig()
	parseFlags(configInit)
	storage.InitJSONFileStorage()

	conf := configs.GetConfig()
	fmt.Println("Running server on", conf.HostServer)
	return http.ListenAndServe(conf.HostServer, router.Router())
}
