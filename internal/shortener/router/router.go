package router

import (
	"net/http"

	"github.com/MWT-proger/shortener/internal/shortener/handlers"
)

func Router() *http.ServeMux {
	r := http.NewServeMux()
	r.HandleFunc("/", handlers.BaseHandler)
	return r
}
