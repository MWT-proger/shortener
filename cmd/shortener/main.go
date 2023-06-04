package main

import (
	"net/http"

	"github.com/MWT-proger/shortener/internal/shortener/router"
	"github.com/MWT-proger/shortener/internal/shortener/storage"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	storage.InitJSONFileStorage()
	return http.ListenAndServe(`:8080`, router.Router())
}
