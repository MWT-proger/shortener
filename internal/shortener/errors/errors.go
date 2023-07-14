package errors

type ErrorDBNotConnection struct{}

func (m *ErrorDBNotConnection) Error() string {
	return "строка с адресом подключения к БД не указана"
}
