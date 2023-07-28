package storage

import (
	"context"

	"github.com/MWT-proger/shortener/internal/shortener/errors"
	"github.com/MWT-proger/shortener/internal/shortener/models"
	"github.com/google/uuid"
)

type OperationStorager interface {
	Set(newModel models.ShortURL) (string, error)
	SetMany(data []models.JSONShortURL, baseShortURL string, userID uuid.UUID) error
	Get(shortURL string) (string, error)
	Init(ctx context.Context) error
	Close() error
	Ping() error
}

type Storage struct{}

// Абстрактный метод
func (s *Storage) Init(ctx context.Context) error {
	return nil
}

// Абстрактный метод
func (s *Storage) Set(newModel models.ShortURL) (string, error) {
	return "", nil

}

// Абстрактный метод
func (s *Storage) SetMany(data []models.JSONShortURL, baseShortURL string, userID uuid.UUID) error {
	return nil

}

// Абстрактный метод
func (s *Storage) Get(shortURL string) (string, error) {
	return "", nil
}

// Абстрактный метод
func (s *Storage) Ping() error {
	return &errors.ErrorDBNotConnection{}
}

// Абстрактный метод
func (s *Storage) Close() error {
	return nil
}
