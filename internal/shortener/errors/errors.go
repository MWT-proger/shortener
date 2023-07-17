package errors

type ErrorDBNotConnection struct{}

func (m *ErrorDBNotConnection) Error() string {
	return "строка с адресом подключения к БД не указана"
}

type ErrorDuplicateShortKey struct{}

func (m *ErrorDuplicateShortKey) Error() string {
	return "повторяющееся значение short_key нарушает уникальное ограничение"
}
