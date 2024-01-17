package services

// PingStorage() Проверяет соединение с базой данных.
func (s *ShortenerService) PingStorage() bool {

	if err := s.storage.Ping(); err != nil {
		return false
	}
	return true

}
