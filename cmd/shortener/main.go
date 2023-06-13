package main

import (
	"context"

	"github.com/MWT-proger/shortener/configs"
	"github.com/MWT-proger/shortener/internal/shortener/handlers"
	"github.com/MWT-proger/shortener/internal/shortener/router"
	"github.com/MWT-proger/shortener/internal/shortener/server"
	"github.com/MWT-proger/shortener/internal/shortener/storage"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())

	if err := run(ctx); err != nil {
		cancel()
		panic(err)
	}
}

// initProject() иницилизирует все необходимые переменный проекта
func initProject() error {
	configInit := configs.InitConfig()

	parseFlags(configInit)

	configs.SetConfigFromEnv()

	return nil
}

// run() выполняет все предворительные действия и вызывает функцию запуска сервера
func run(ctx context.Context) error {
	initProject()

	var (
		s    = &storage.Storage{}
		h, _ = handlers.NewAPIHandler(s)
		r    = router.Router(h)
	)

	err := s.InitJSONFile()

	if err != nil {
		return err
	}

	go s.BackupToJSONFile(ctx)

	err = server.Run(r)

	if err != nil {

		return err
	}

	return nil
}
