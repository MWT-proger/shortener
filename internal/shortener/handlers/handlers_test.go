package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/MWT-proger/shortener/configs"
)

type MockStorage struct {
	testData map[string]string
}

func (s *MockStorage) Set(fullURL string) (string, error) {
	return s.testData[fullURL], nil
}

func (s *MockStorage) Get(shortURL string) (string, error) {
	return s.testData[shortURL], nil
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string, bodyReader *strings.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, bodyReader)
	require.NoError(t, err)

	client := ts.Client()
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	resp, err := client.Do(req)

	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func TestAPIHandlerGetURLByKeyHandler(t *testing.T) {

	testCases := []struct {
		name             string
		method           string
		URL              string
		mapKeyValue      map[string]string
		expectedCode     int
		expectedLocation string
	}{
		{name: "Тест 1 - Не верный ключ", URL: "/", mapKeyValue: map[string]string{"testKeyNotValid": "http://example-full-url.com"}, method: http.MethodGet, expectedCode: http.StatusNotFound, expectedLocation: ""},
		{name: "Тест 2 - Не верный ключ", URL: "/testKey", mapKeyValue: map[string]string{"testKeyNotValid": "http://example-full-url.com"}, method: http.MethodGet, expectedCode: http.StatusBadRequest, expectedLocation: ""},
		{name: "Тест 3 - Не верный ключ", URL: "/testKey/testKey", mapKeyValue: map[string]string{"testKey": "http://example-full-url.com"}, method: http.MethodGet, expectedCode: http.StatusNotFound, expectedLocation: ""},
		{name: "Тест 4 - Успешный ответ", URL: "/testKey", mapKeyValue: map[string]string{"testKey": "https://practicum.yandex.ru/"}, method: http.MethodGet, expectedCode: http.StatusTemporaryRedirect, expectedLocation: "https://practicum.yandex.ru/"},
		{name: "Тест 5 - Не верный метод запроса", URL: "/testKey", mapKeyValue: map[string]string{"testKey": "http://example-full-url.com"}, method: http.MethodPost, expectedCode: http.StatusMethodNotAllowed, expectedLocation: ""},
		{name: "Тест 6 - Не верный метод запроса", URL: "/testKey", mapKeyValue: map[string]string{"testKey": "http://example-full-url.com"}, method: http.MethodPut, expectedCode: http.StatusMethodNotAllowed, expectedLocation: ""},
		{name: "Тест 7 - Не верный метод запроса", URL: "/testKey", mapKeyValue: map[string]string{"testKey": "http://example-full-url.com"}, method: http.MethodDelete, expectedCode: http.StatusMethodNotAllowed, expectedLocation: ""},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			m := &MockStorage{testData: tt.mapKeyValue}
			h := &APIHandler{m}

			router := chi.NewRouter()
			router.Get("/{shortKey}", h.GetURLByKeyHandler)

			ts := httptest.NewServer(router)

			result, _ := testRequest(t, ts, tt.method, tt.URL, strings.NewReader(""))
			defer result.Body.Close()

			assert.Equal(t, tt.expectedCode, result.StatusCode, "Код ответа не совпадает с ожидаемым")
			assert.Equal(t, tt.expectedLocation, result.Header.Get("Location"), "Location не совпадает с ожидаемым")
		})
	}
}

func TestAPIHandlerGenerateShortkeyHandler(t *testing.T) {

	testCases := []struct {
		name         string
		method       string
		URL          string
		key          string
		mapKeyValue  map[string]string
		expectedCode int
		expectedBody string
		envBaseURL   string
	}{
		{name: "Тест 1 - Не верный URL", URL: "/testKey", key: "http://example-full-url.com", mapKeyValue: map[string]string{"http://example-full-url.com": "testKey"}, method: http.MethodPost, expectedCode: http.StatusNotFound, expectedBody: "", envBaseURL: ""},
		{name: "Тест 2 - Успешный запрос", URL: "/", key: "http://example-full-url.com", mapKeyValue: map[string]string{"http://example-full-url.com": "testKey"}, method: http.MethodPost, expectedCode: http.StatusCreated, expectedBody: "%v%v/testKey", envBaseURL: ""},
		{name: "Тест 3 - Не верный метод запроса", URL: "/", key: "http://example-full-url.com", mapKeyValue: map[string]string{"http://example-full-url.com": "testKey"}, method: http.MethodGet, expectedCode: http.StatusMethodNotAllowed, expectedBody: "", envBaseURL: ""},
		{name: "Тест 4 - Не верный метод запроса", URL: "/", key: "http://example-full-url.com", mapKeyValue: map[string]string{"http://example-full-url.com": "testKey"}, method: http.MethodPut, expectedCode: http.StatusMethodNotAllowed, expectedBody: "", envBaseURL: ""},
		{name: "Тест 5 - Не верный метод запроса", URL: "/", key: "http://example-full-url.com", mapKeyValue: map[string]string{"http://example-full-url.com": "testKey"}, method: http.MethodDelete, expectedCode: http.StatusMethodNotAllowed, expectedBody: "", envBaseURL: ""},
		{name: "Тест 6 - Проверка BaseURL из ENV", URL: "/", key: "http://example-full-url.com", mapKeyValue: map[string]string{"http://example-full-url.com": "testKey"}, method: http.MethodPost, expectedCode: http.StatusCreated, expectedBody: "http://site.com/testKey", envBaseURL: "http://site.com"},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envBaseURL != "" {
				os.Setenv("BASE_URL", tt.envBaseURL)
			}
			configs.InitConfig()
			configs.SetConfigFromEnv()

			m := &MockStorage{testData: tt.mapKeyValue}
			h := &APIHandler{m}

			router := chi.NewRouter()
			router.Post("/", h.GenerateShortkeyHandler)

			ts := httptest.NewServer(router)
			bodyRequest := strings.NewReader(tt.key)

			result, bodyResponse := testRequest(t, ts, tt.method, tt.URL, bodyRequest)
			defer result.Body.Close()

			assert.Equal(t, tt.expectedCode, result.StatusCode, "Код ответа не совпадает с ожидаемым")
			if tt.expectedCode == http.StatusCreated {

				if tt.envBaseURL != "" {
					assert.Equal(t, tt.expectedBody, bodyResponse, "Тело ответа не совпадает с ожидаемым")
				} else {
					resp := fmt.Sprintf(tt.expectedBody, "http://", result.Request.URL.Host)
					assert.Equal(t, resp, bodyResponse, "Тело ответа не совпадает с ожидаемым")
				}

			}

		})
	}
	os.Setenv("BASE_URL", "")
}

func TestAPIHandlerJSONGenerateShortkeyHandler(t *testing.T) {

	testCases := []struct {
		name         string
		method       string
		URL          string
		requestData  string
		mapKeyValue  map[string]string
		expectedCode int
		expectedBody string
		envBaseURL   string
	}{
		{name: "Тест 1 - Не верный URL", URL: "/api/shorten/testKey", requestData: `{"url": "http://example-full-url.com"}`, mapKeyValue: map[string]string{"http://example-full-url.com": "testKey"}, method: http.MethodPost, expectedCode: http.StatusNotFound, expectedBody: "", envBaseURL: ""},
		{name: "Тест 2 - Успешный запрос", URL: "/api/shorten", requestData: `{"url": "http://example-full-url.com"}`, mapKeyValue: map[string]string{"http://example-full-url.com": "testKey"}, method: http.MethodPost, expectedCode: http.StatusCreated, expectedBody: `{"result": "%v%v/testKey"}`, envBaseURL: ""},
		{name: "Тест 3 - Не верный метод запроса", URL: "/api/shorten", requestData: `{"url": "http://example-full-url.com"}`, mapKeyValue: map[string]string{"http://example-full-url.com": "testKey"}, method: http.MethodGet, expectedCode: http.StatusMethodNotAllowed, expectedBody: "", envBaseURL: ""},
		{name: "Тест 4 - Не верный метод запроса", URL: "/api/shorten", requestData: `{"url": "http://example-full-url.com"}`, mapKeyValue: map[string]string{"http://example-full-url.com": "testKey"}, method: http.MethodPut, expectedCode: http.StatusMethodNotAllowed, expectedBody: "", envBaseURL: ""},
		{name: "Тест 5 - Не верный метод запроса", URL: "/api/shorten", requestData: `{"url": "http://example-full-url.com"}`, mapKeyValue: map[string]string{"http://example-full-url.com": "testKey"}, method: http.MethodDelete, expectedCode: http.StatusMethodNotAllowed, expectedBody: "", envBaseURL: ""},
		{name: "Тест 6 - Проверка BaseURL из ENV", URL: "/api/shorten", requestData: `{"url": "http://example-full-url.com"}`, mapKeyValue: map[string]string{"http://example-full-url.com": "testKey"}, method: http.MethodPost, expectedCode: http.StatusCreated, expectedBody: `{"result": "http://site.com/testKey"}`, envBaseURL: "http://site.com"},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envBaseURL != "" {
				os.Setenv("BASE_URL", tt.envBaseURL)
			}
			configs.InitConfig()
			configs.SetConfigFromEnv()

			m := &MockStorage{testData: tt.mapKeyValue}
			h := &APIHandler{m}

			router := chi.NewRouter()
			router.Post("/api/shorten", h.JSONGenerateShortkeyHandler)

			ts := httptest.NewServer(router)

			bodyRequest := strings.NewReader(tt.requestData)

			result, bodyResponse := testRequest(t, ts, tt.method, tt.URL, bodyRequest)

			defer result.Body.Close()

			assert.Equal(t, tt.expectedCode, result.StatusCode, "Код ответа не совпадает с ожидаемым")
			if tt.expectedCode == http.StatusCreated {
				if tt.envBaseURL != "" {
					assert.JSONEq(t, tt.expectedBody, bodyResponse, "Тело ответа не совпадает с ожидаемым")
				} else {
					resp := fmt.Sprintf(tt.expectedBody, "http://", result.Request.URL.Host)
					assert.JSONEq(t, resp, bodyResponse, "Тело ответа не совпадает с ожидаемым")
				}

			}

		})
	}

}
