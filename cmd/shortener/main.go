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

// initProject() иницилизирует все необходимые переменный проекта
func initProject() {
	configInit := configs.InitConfig()

	parseFlags(configInit)

	configs.SetConfigFromEnv()

	storage.InitJSONFile()
}

// run() запускает сервер
func run() error {
	initProject()

	conf := configs.GetConfig()
	fmt.Println("Running server on", conf.HostServer)
	return http.ListenAndServe(conf.HostServer, router.Router())
}
