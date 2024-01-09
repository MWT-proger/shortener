package main

import (
	"context"

	"github.com/MWT-proger/shortener/configs"
	"github.com/MWT-proger/shortener/internal/shortener/handlers"
	"github.com/MWT-proger/shortener/internal/shortener/logger"
	"github.com/MWT-proger/shortener/internal/shortener/router"
	"github.com/MWT-proger/shortener/internal/shortener/server"
	"github.com/MWT-proger/shortener/internal/shortener/services"
	"github.com/MWT-proger/shortener/internal/shortener/storage"
	"github.com/MWT-proger/shortener/internal/shortener/storage/filestorage"
	"github.com/MWT-proger/shortener/internal/shortener/storage/pgstorage"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := run(ctx); err != nil {
		cancel()
		panic(err)
	}
}

// initProject() иницилизирует все необходимые переменный проекта
func initProject(ctx context.Context) (storage.OperationStorager, error) {

	var (
		s          storage.OperationStorager
		configInit = configs.InitConfig()
	)

	parseFlags(configInit)

	conf := configs.SetConfigFromEnv()

	if conf.DatabaseDSN != "" {
		s = &pgstorage.PgStorage{}
	} else {
		s = &filestorage.FileStorage{}
	}

	if err := s.Init(ctx); err != nil {
		return nil, err
	}

	if err := logger.Initialize(conf.LogLevel); err != nil {
		return nil, err
	}

	return s, nil
}

// run() выполняет все предворительные действия и вызывает функцию запуска сервера
func run(ctx context.Context) error {

	storage, err := initProject(ctx)

	if err != nil {
		return err
	}

	defer storage.Close()

	service := services.NewShortenerService(storage)

	h, _ := handlers.NewAPIHandler(service)

	r := router.Router(h)

	err = server.Run(r)

	if err != nil {
		return err
	}

	return nil
}
