package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockStorage struct {
	testData map[string]string
}

func (s *MockStorage) SetInStorage(fullURL string) (string, error) {
	return s.testData[fullURL], nil
}

func (s *MockStorage) GetFromStorage(shortURL string) (string, error) {
	return s.testData[shortURL], nil
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
		{name: "Тест 1 - Не верный ключ", URL: "/", mapKeyValue: map[string]string{"testKeyNotValid": "http://example-full-url.com"}, method: http.MethodGet, expectedCode: http.StatusBadRequest, expectedLocation: ""},
		{name: "Тест 2 - Не верный ключ", URL: "/testKey", mapKeyValue: map[string]string{"testKeyNotValid": "http://example-full-url.com"}, method: http.MethodGet, expectedCode: http.StatusBadRequest, expectedLocation: ""},
		{name: "Тест 3 - Не верный ключ", URL: "/testKey/testKey", mapKeyValue: map[string]string{"testKey": "http://example-full-url.com"}, method: http.MethodGet, expectedCode: http.StatusBadRequest, expectedLocation: ""},
		{name: "Тест 4 - Успешный ответ", URL: "/testKey", mapKeyValue: map[string]string{"testKey": "http://example-full-url.com"}, method: http.MethodGet, expectedCode: http.StatusTemporaryRedirect, expectedLocation: "http://example-full-url.com"},
		{name: "Тест 5 - Не верный метод запроса", URL: "/testKey", mapKeyValue: map[string]string{"testKey": "http://example-full-url.com"}, method: http.MethodPost, expectedCode: http.StatusBadRequest, expectedLocation: ""},
		{name: "Тест 6 - Не верный метод запроса", URL: "/testKey", mapKeyValue: map[string]string{"testKey": "http://example-full-url.com"}, method: http.MethodPut, expectedCode: http.StatusBadRequest, expectedLocation: ""},
		{name: "Тест 7 - Не верный метод запроса", URL: "/testKey", mapKeyValue: map[string]string{"testKey": "http://example-full-url.com"}, method: http.MethodDelete, expectedCode: http.StatusBadRequest, expectedLocation: ""},
	}

	h := &APIHandler{}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(tt.method, tt.URL, nil)
			w := httptest.NewRecorder()
			m := &MockStorage{testData: tt.mapKeyValue}

			h.GetURLByKeyHandler(w, r, m)
			result := w.Result()
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
	}{
		{name: "Тест 1 - Не верный URL", URL: "/testKey", key: "http://example-full-url.com", mapKeyValue: map[string]string{"http://example-full-url.com": "testKey"}, method: http.MethodPost, expectedCode: http.StatusBadRequest, expectedBody: ""},
		{name: "Тест 2 - Успешный запрос", URL: "/", key: "http://example-full-url.com", mapKeyValue: map[string]string{"http://example-full-url.com": "testKey"}, method: http.MethodPost, expectedCode: http.StatusCreated, expectedBody: "http://example.com/testKey"},
		{name: "Тест 3 - Не верный метод запроса", URL: "/", key: "http://example-full-url.com", mapKeyValue: map[string]string{"http://example-full-url.com": "testKey"}, method: http.MethodGet, expectedCode: http.StatusBadRequest, expectedBody: ""},
		{name: "Тест 4 - Не верный метод запроса", URL: "/", key: "http://example-full-url.com", mapKeyValue: map[string]string{"http://example-full-url.com": "testKey"}, method: http.MethodPut, expectedCode: http.StatusBadRequest, expectedBody: ""},
		{name: "Тест 5 - Не верный метод запроса", URL: "/", key: "http://example-full-url.com", mapKeyValue: map[string]string{"http://example-full-url.com": "testKey"}, method: http.MethodDelete, expectedCode: http.StatusBadRequest, expectedBody: ""},
	}

	h := &APIHandler{}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {

			bodyReader := strings.NewReader(tt.key)

			r := httptest.NewRequest(tt.method, tt.URL, bodyReader)

			w := httptest.NewRecorder()
			m := &MockStorage{testData: tt.mapKeyValue}

			h.GenerateShortkeyHandler(w, r, m)
			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.expectedCode, result.StatusCode, "Код ответа не совпадает с ожидаемым")
			if tt.expectedCode == http.StatusOK {
				assert.Equal(t, tt.expectedBody, w.Body.String(), "Тело ответа не совпадает с ожидаемым")
			}

		})
	}
}
