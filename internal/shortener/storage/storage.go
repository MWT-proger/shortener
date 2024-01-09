package storage

import (
	"context"

	"github.com/MWT-proger/shortener/internal/shortener/errors"
	"github.com/MWT-proger/shortener/internal/shortener/models"
	"github.com/google/uuid"
)

// OperationStorager общее представление хранилища
type OperationStorager interface {
	Set(newModel models.ShortURL) (string, error)
	SetMany(data []models.JSONShortURL, baseShortURL string, userID uuid.UUID) error
	DeleteList(data ...models.DeletedShortURL) error
	Get(shortURL string) (models.ShortURL, error)
	GetList(userID uuid.UUID) ([]*models.JSONShortURL, error)
	Init(ctx context.Context) error
	Close() error
	Ping() error
}

// Storage базовое хранилище
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
func (s *Storage) Get(shortURL string) (models.ShortURL, error) {
	return models.ShortURL{}, nil
}

// Абстрактный метод
func (s *Storage) GetList(userID uuid.UUID) ([]*models.JSONShortURL, error) {
	return []*models.JSONShortURL{}, nil
}

// Абстрактный метод
func (s *Storage) Ping() error {
	return errors.ErrorDBNotConnection
}

// Абстрактный метод
func (s *Storage) Close() error {
	return nil
}

// Абстрактный метод
func (s *Storage) DeleteList(data ...models.DeletedShortURL) error {
	return nil
}
