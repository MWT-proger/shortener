package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/MWT-proger/shortener/configs"
	"github.com/MWT-proger/shortener/internal/shortener/grpc"
	"github.com/MWT-proger/shortener/internal/shortener/handlers"
	"github.com/MWT-proger/shortener/internal/shortener/logger"
	"github.com/MWT-proger/shortener/internal/shortener/server"
	"github.com/MWT-proger/shortener/internal/shortener/services"
	"github.com/MWT-proger/shortener/internal/shortener/storage"
	"github.com/MWT-proger/shortener/internal/shortener/storage/filestorage"
	"github.com/MWT-proger/shortener/internal/shortener/storage/pgstorage"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	printBuild()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer stop()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if err := run(ctx); err != nil {
		cancel()
		time.Sleep(time.Second * 5)
		if !errors.Is(err, http.ErrServerClosed) {
			logger.Log.Error("Сервер остановлен из за критичиской ошибки", logger.ErrorField(err))
		}
		logger.Log.Info("Успешное завершение сервера")

	}
}

func printBuild() {
	if buildVersion == "" {
		buildVersion = "N/A"
	}
	fmt.Printf("Build version: %s\n", buildVersion)

	if buildDate == "" {
		buildDate = "N/A"
	}
	fmt.Printf("Build date: %s\n", buildDate)

	if buildCommit == "" {
		buildCommit = "N/A"
	}
	fmt.Printf("Build commit: %s\n", buildCommit)

}

// run() выполняет все предворительные действия и вызывает функцию запуска сервера.
func run(ctx context.Context) error {

	var (
		conf, err = configs.NewConfig()
		storage   storage.OperationStorager
	)

	if err != nil {
		return err
	}

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

	if conf.RunGRPC {

		gRPCServer, err := grpc.NewShortenerGRPCServer(service)

		if err != nil {
			return err
		}

		if err := gRPCServer.Run(ctx, conf); err != nil {
			return err
		}

	} else {
		apiHandler, err := handlers.NewAPIHandler(service)

		if err != nil {
			return err
		}

		if err := server.Run(ctx, apiHandler, conf); err != nil {
			return err
		}
	}

	return nil
}
