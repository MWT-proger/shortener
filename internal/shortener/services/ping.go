package services

// GenerateShortKeyHandler Принимает большой URL и возвращает маленький
func (s *ShortenerService) PingStorage() bool {

	if err := s.storage.Ping(); err != nil {
		return false
	}
	return true

}
