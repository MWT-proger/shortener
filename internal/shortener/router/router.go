package router

import (
	"net/http"

	"github.com/MWT-proger/shortener/internal/shortener/handlers"
)

func Router() *http.ServeMux {
	r := http.NewServeMux()
	h, _ := handlers.NewAPIHandler()
	r.HandleFunc("/", h.BaseHandler)
	return r
}
