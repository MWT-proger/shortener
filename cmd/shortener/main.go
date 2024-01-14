package main

import (
	"context"

	"github.com/MWT-proger/shortener/configs"
	"github.com/MWT-proger/shortener/internal/shortener/handlers"
	"github.com/MWT-proger/shortener/internal/shortener/logger"
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

// run() выполняет все предворительные действия и вызывает функцию запуска сервера.
func run(ctx context.Context) error {

	var (
		conf    = configs.InitConfig()
		storage storage.OperationStorager
		err     error
	)

	if err = logger.Initialize(conf.LogLevel); err != nil {
		return err
	}

	if conf.DatabaseDSN != "" {
		storage, err = pgstorage.NewPgStorage(ctx, conf)
	} else {
		storage, err = filestorage.NewFileStorage(ctx, conf)
	}

	if err != nil {
		return err
	}

	defer storage.Close()

	service := services.NewShortenerService(ctx, conf, storage)

	apiHandler, err := handlers.NewAPIHandler(service)

	if err != nil {
		return err
	}

	if err := server.Run(apiHandler, conf); err != nil {
		return err
	}

	return nil
}
