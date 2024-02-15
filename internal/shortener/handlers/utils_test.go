package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	lErrors "github.com/MWT-proger/shortener/internal/shortener/errors"
)

// 1. Тест на успешное распаковывание тела запроса:
func TestUnmarshalBody_Success(t *testing.T) {
	handler := &APIHandler{}
	body := io.NopCloser(strings.NewReader(`{"name":"John","age":30}`))
	form := &struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{}

	err := handler.unmarshalBody(body, form)

	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if form.Name != "John" {
		t.Errorf("Expected Name to be 'John', but got: %s", form.Name)
	}

	if form.Age != 30 {
		t.Errorf("Expected Age to be 30, but got: %d", form.Age)
	}
}

// 2. Тест на ошибку при чтении тела запроса:
func TestUnmarshalBody_ReadError(t *testing.T) {
	handler := &APIHandler{}
	body := io.NopCloser(strings.NewReader(`{"name":"John","age":30}`))

	err := handler.unmarshalBody(body, nil)

	if err == nil {
		t.Error("Expected an error, but got nil")
	}
}

// 3. Тест на ошибку при распаковывании тела запроса:
func TestUnmarshalBody_UnmarshalError(t *testing.T) {
	handler := &APIHandler{}
	body := io.NopCloser(strings.NewReader(`{"name":"John","age":"thirty"}`))
	form := &struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{}

	err := handler.unmarshalBody(body, form)

	if err == nil {
		t.Error("Expected an error, but got nil")
	}
}

// 1. Тест на ошибку сервисного слоя:
func TestSetHTTPError_ServiceError(t *testing.T) {
	handler := &APIHandler{}
	w := httptest.NewRecorder()
	err := lErrors.NewServicesError("Service error", http.StatusBadRequest, false, 0)

	handler.setHTTPError(w, err)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, but got %d", http.StatusBadRequest, w.Code)
	}

	if w.Body.String() != "Service error\n" {
		t.Errorf("Expected response body 'Service error', but got '%s'", w.Body.String())
	}
}

// 2. Тест на ошибку сервера:
func TestSetHTTPError_InternalServerError(t *testing.T) {
	handler := &APIHandler{}
	w := httptest.NewRecorder()
	err := fmt.Errorf("Internal server error")

	handler.setHTTPError(w, err)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, but got %d", http.StatusInternalServerError, w.Code)
	}

	if w.Body.String() != "Ошибка сервера, попробуйте позже.\n" {
		t.Errorf("Expected response body 'Ошибка сервера, попробуйте позже.', but got '%s'", w.Body.String())
	}
}

// 1. Тест на ошибку сервисного слоя и досрочное выполнение кода:
func TestSetOrGetHTTPCode_ServiceError_ReturnTrue(t *testing.T) {
	handler := &APIHandler{}
	w := httptest.NewRecorder()
	err := lErrors.NewServicesError("Service error", http.StatusBadRequest, true, 0)
	expectedCode := 0

	code := handler.setOrGetHTTPCode(w, err)

	if code != expectedCode {
		t.Errorf("Expected HTTP code %d, but got %d", expectedCode, code)
	}

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, but got %d", expectedCode, w.Code)
	}

	if w.Body.String() != "Service error\n" {
		t.Errorf("Expected response body 'Service error', but got '%s'", w.Body.String())
	}
}
