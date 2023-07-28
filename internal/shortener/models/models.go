package models

import "github.com/google/uuid"

type ShortURL struct {
	ShortKey string
	FullURL  string
	UserID   uuid.UUID
}

type JSONShortURL struct {
	CorrelationID string `json:"correlation_id,omitempty"`
	OriginalURL   string `json:"original_url,omitempty"`
	ShortURL      string `json:"short_url,omitempty"`
}

func (d *JSONShortURL) IsValid() bool {

	if d.OriginalURL == "" || d.CorrelationID == "" {
		return false
	}

	return true
}

type JSONShortenRequest struct {
	URL string `json:"url"`
}

func (d *JSONShortenRequest) IsValid() bool {
	return d.URL != ""
}
