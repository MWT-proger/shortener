package main

import (
	"net/http"

	"github.com/MWT-proger/shortener/internal/shortener/router"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	return http.ListenAndServe(`:8080`, router.Router())
}
