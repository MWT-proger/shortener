package handlers

// import (
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"net/http/httptest"
// 	"strings"

// 	"github.com/MWT-proger/shortener/configs"
// 	"github.com/MWT-proger/shortener/internal/shortener/models"
// 	"github.com/go-chi/chi"
// )

// func ExampleAPIHandler_GetURLByKeyHandler() {

// 	m := &MockStorage{testData: map[string]string{"testKey": "https://practicum.yandex.ru/"}}

// 	h := &APIHandler{
// 		storage:     m,
// 		DeletedChan: make(chan models.DeletedShortURL, 1024),
// 	}

// 	router := chi.NewRouter()
// 	router.Get("/{shortKey}", h.GetURLByKeyHandler)

// 	ts := httptest.NewServer(router)

// 	result, _ := exampleRequest(ts, http.MethodGet, "/testKey", strings.NewReader(""))
// 	defer result.Body.Close()
// 	location, _ := result.Location()
// 	fmt.Println(location)
// 	// Output:
// 	// https://practicum.yandex.ru/
// }

// func ExampleAPIHandler_GenerateShortkeyHandler() {
// 	configs.InitConfig()

// 	m := &MockStorage{testData: map[string]string{"http://example-full-url.com": "testKey"}}
// 	h := &APIHandler{
// 		storage:     m,
// 		DeletedChan: make(chan models.DeletedShortURL, 1024),
// 	}

// 	router := chi.NewRouter()
// 	router.Post("/", h.GenerateShortkeyHandler)

// 	ts := httptest.NewServer(router)
// 	bodyRequest := strings.NewReader("http://example-full-url.com")

// 	result, _ := exampleRequest(ts, http.MethodPost, "/", bodyRequest)
// 	defer result.Body.Close()
// 	fmt.Println(result.Status)
// 	// Output:
// 	// 201 Created

// }

// func ExampleAPIHandler_JSONGenerateShortkeyHandler() {
// 	configs.InitConfig()

// 	m := &MockStorage{testData: map[string]string{"http://example-full-url.com": "testKey"}}
// 	h := &APIHandler{
// 		storage:     m,
// 		DeletedChan: make(chan models.DeletedShortURL, 1024),
// 	}

// 	router := chi.NewRouter()
// 	router.Post("/api/shorten", h.JSONGenerateShortkeyHandler)

// 	ts := httptest.NewServer(router)

// 	bodyRequest := strings.NewReader(`{"url": "http://example-full-url.com"}`)

// 	result, _ := exampleRequest(ts, http.MethodPost, "/api/shorten", bodyRequest)
// 	defer result.Body.Close()
// 	fmt.Println(result.Status)
// 	// Output:
// 	// 201 Created

// }

// func exampleRequest(ts *httptest.Server, method, path string, bodyReader *strings.Reader) (*http.Response, string) {
// 	req, err := http.NewRequest(method, ts.URL+path, bodyReader)

// 	if err != nil {
// 		return nil, ""
// 	}

// 	client := ts.Client()
// 	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
// 		return http.ErrUseLastResponse
// 	}

// 	resp, err := client.Do(req)

// 	if err != nil {
// 		return nil, ""
// 	}

// 	defer resp.Body.Close()

// 	respBody, err := io.ReadAll(resp.Body)

// 	if err != nil {
// 		return nil, ""
// 	}

// 	return resp, string(respBody)
// }
