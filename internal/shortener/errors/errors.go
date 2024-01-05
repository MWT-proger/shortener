package errors

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

// Error
type ErrorDuplicateFullURL struct{}

// Error
func (m *ErrorDuplicateFullURL) Error() string {
	return "повторяющееся значение full_url нарушает уникальное ограничение"
}
