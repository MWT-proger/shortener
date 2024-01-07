package errors

import "net/http"

// Error
type ErrorDBNotConnection struct{}

// Error
func (m *ErrorDBNotConnection) Error() string {
	return "строка с адресом подключения к БД не указана"
}

// Error
type ErrorDuplicateShortKey struct{}

// Error
func (m *ErrorDuplicateShortKey) Error() string {
	return "повторяющееся значение short_key нарушает уникальное ограничение"
}

// ServicesError тип ошибки сервисного слоя.
// Содержит текст ошибки и http код ошибки.
type ServicesError struct {
	s           string
	HTTPCode    int
	ContentType string
}

// Error базовый метод ошибки.
func (e *ServicesError) Error() string {
	return e.s
}

// NewServicesError создает новую сервисную ошибку.
func NewServicesError(text string, httpCode int, contentType string) *ServicesError {
	return &ServicesError{text, httpCode, contentType}
}

// GetFullURLServicesError ошибка получения FullURL.
var GetFullURLServicesError = NewServicesError(
	"не получилось обработать запрос получения URL",
	http.StatusInternalServerError,
	"",
)

// GoneServicesError ошибка получения объекта. Объект Удален.
var GoneServicesError = NewServicesError(
	"запрашиваемый объект удален",
	http.StatusGone,
	"",
)

// NotFoundServicesErro ошибка получения объекта. Объект не найден.
var NotFoundServicesError = NewServicesError(
	"запрашиваемый объект не найден",
	http.StatusBadRequest,
	"",
)

// NoContentUserServicesError список URL-адресов пользователя пуст
var NoContentUserServicesError = NewServicesError(
	"список URL-адресов пользователя пуст",
	http.StatusNoContent,
	"",
)

// NoContentUserServicesError список URL-адресов пользователя пуст
var ErrorDuplicateFullURLServicesError = NewServicesError(
	"повторяющееся значение full_url нарушает уникальное ограничение",
	http.StatusInternalServerError,
	"",
)

// NoContentUserServicesError список URL-адресов пользователя пуст
var ErrorSetConflictServicesError = NewServicesError(
	"создаваемый объект уже существует",
	http.StatusConflict,
	"text/plain",
)
