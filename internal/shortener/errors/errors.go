package errors

import (
	"errors"
	"net/http"
)

// ErrorDBNotConnection ошибка подключения хранилища.
var ErrorDBNotConnection = errors.New("строка с адресом подключения к БД не указана")

// ErrorDuplicateShortKey данный short_key уже указан.
var ErrorDuplicateShortKey = errors.New("повторяющееся значение short_key нарушает уникальное ограничение")

// ServicesError тип ошибки сервисного слоя.
// Содержит текст ошибки и http код ошибки.
type ServicesError struct {
	s        string
	HTTPCode int
	IsReturn bool
}

// Error базовый метод ошибки.
func (e *ServicesError) Error() string {
	return e.s
}

// NewServicesError создает новую сервисную ошибку.
func NewServicesError(text string, httpCode int, isReturn bool) *ServicesError {
	return &ServicesError{text, httpCode, isReturn}
}

// GetFullURLServicesError ошибка получения FullURL.
var GetFullURLServicesError = NewServicesError(
	"не получилось обработать запрос получения URL",
	http.StatusInternalServerError,
	true,
)

// GoneServicesError ошибка получения объекта. Объект Удален.
var GoneServicesError = NewServicesError(
	"запрашиваемый объект удален",
	http.StatusGone,
	true,
)

// NotFoundServicesError ошибка получения объекта. Объект не найден.
var NotFoundServicesError = NewServicesError(
	"запрашиваемый объект не найден",
	http.StatusBadRequest,
	true,
)

// NoContentUserServicesError список URL-адресов пользователя пуст.
var NoContentUserServicesError = NewServicesError(
	"список URL-адресов пользователя пуст",
	http.StatusNoContent,
	true,
)

// ErrorDuplicateFullURLServicesError у данного пользователя full_url уже зарегистрирован.
var ErrorDuplicateFullURLServicesError = NewServicesError(
	"повторяющееся значение full_url нарушает уникальное ограничение",
	http.StatusConflict,
	false,
)

// InternalServicesError внутренняя ошибка сервиса.
var InternalServicesError = NewServicesError(
	"внутренняя ошибка сервиса ",
	http.StatusInternalServerError,
	true,
)
