package models

import (
	"github.com/google/uuid"
)

// ShortURL - служит связующим звеном между слоями сервисным и хранилища.
type ShortURL struct {
	ShortKey    string
	FullURL     string
	UserID      uuid.UUID
	DeletedFlag bool
}

// JSONShortURL - служит связующим звеном между слоями сервисным и хранилища.
type JSONShortURL struct {
	CorrelationID string    `json:"correlation_id,omitempty"`
	OriginalURL   string    `json:"original_url,omitempty" db:"full_url"`
	ShortURL      string    `json:"short_url,omitempty" db:"short_key"`
	UserID        uuid.UUID `json:"-" db:"user_id"`
	DeletedFlag   bool      `json:"-" db:"is_deleted"`
}

// IsValid проверяет на валидность JSONShortURL.
func (d *JSONShortURL) IsValid() bool {

	if d.OriginalURL == "" || d.CorrelationID == "" {
		return false
	}

	return true
}

// DeletedShortURL Используется при удаление строк из БД.
type DeletedShortURL struct {
	ID      int64
	UserID  uuid.UUID
	Payload string
}
