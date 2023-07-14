package storage

import (
	"context"

	"github.com/MWT-proger/shortener/internal/shortener/errors"
)

type OperationStorager interface {
	Set(fullURL string) (string, error)
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
func (s *Storage) Set(fullURL string) (string, error) {
	return "", nil

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
