package handlers

import (
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/MWT-proger/shortener/internal/shortener/storage"
)

func BaseHandler(res http.ResponseWriter, req *http.Request) {
	// Базовый хендлер для хардкорного первого инкремента

	ct := req.Header.Get("Content-Type")
	if ct != "text/plain" {
		http.Error(res, "Invalid Content-Type header type.", http.StatusBadRequest)
		return
	}

	if idx := strings.Index(req.URL.Path[1:], "/"); req.Method == http.MethodGet && idx < 0 && req.URL.Path != "/" {

		GetUrlByKeyHandler(res, req)

	} else if req.Method == http.MethodPost && req.URL.Path == "/" {

		GenerateShortkeyHandler(res, req)

	} else {
		http.Error(res, "Bad Request", http.StatusBadRequest)
	}

}

func GenerateShortkeyHandler(res http.ResponseWriter, req *http.Request) {
	var shortUrl string

	defer req.Body.Close()
	requestData, err := io.ReadAll((req.Body))
	if err != nil {
		log.Fatal(err)
	}

	stringRequestData := string(requestData)
	if stringRequestData == "" {
		http.Error(res, "Bad Request", http.StatusBadRequest)
		return
	}

	shortUrl = storage.SetInStorage(stringRequestData)

	res.Header().Set("content-type", "text/plain")
	res.WriteHeader(http.StatusCreated)

	res.Write([]byte(req.Host))
	res.Write([]byte("/"))
	res.Write([]byte(shortUrl))

}

func GetUrlByKeyHandler(res http.ResponseWriter, req *http.Request) {

	fullUrl := storage.GetFromStorage(req.URL.Path[1:])

	if fullUrl == "" {
		http.Error(res, "Bad Request", http.StatusBadRequest)
		return
	}

	res.Header().Set("Location", fullUrl)
	res.WriteHeader(http.StatusTemporaryRedirect)
}
