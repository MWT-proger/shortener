package storage

import (
	"context"

	"github.com/google/uuid"

	"github.com/MWT-proger/shortener/internal/shortener/models"
)

// OperationStorager Интерфейс хранилища.
type OperationStorager interface {
	Set(ctx context.Context, newModel models.ShortURL) (string, error)
	SetMany(ctx context.Context, data []models.JSONShortURL, baseShortURL string, userID uuid.UUID) error

	Get(ctx context.Context, shortURL string) (models.ShortURL, error)
	GetList(ctx context.Context, userID uuid.UUID) ([]*models.JSONShortURL, error)

	DeleteList(ctx context.Context, data ...models.DeletedShortURL) error

	Close() error
	Ping() error
}
