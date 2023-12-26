package models

import (
	"github.com/google/uuid"
)

// Сущность
type ShortURL struct {
	ShortKey    string
	FullURL     string
	UserID      uuid.UUID
	DeletedFlag bool
}

// Сущность
type JSONShortURL struct {
	CorrelationID string    `json:"correlation_id,omitempty"`
	OriginalURL   string    `json:"original_url,omitempty" db:"full_url"`
	ShortURL      string    `json:"short_url,omitempty" db:"short_key"`
	UserID        uuid.UUID `json:"-" db:"user_id"`
	DeletedFlag   bool      `json:"-" db:"is_deleted"`
}

// IsValid проверяет на валидность
func (d *JSONShortURL) IsValid() bool {

	if d.OriginalURL == "" || d.CorrelationID == "" {
		return false
	}

	return true
}

// Сущность
type JSONShortenRequest struct {
	URL string `json:"url"`
}

// IsValid проверяет на валидность
func (d *JSONShortenRequest) IsValid() bool {
	return d.URL != ""
}

// Сущность
type DeletedShortURL struct {
	ID      int64
	UserID  uuid.UUID
	Payload string
}
