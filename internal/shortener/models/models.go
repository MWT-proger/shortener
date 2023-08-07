package models

type ShortURL struct {
	ShortKey string
	FullURL  string
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
