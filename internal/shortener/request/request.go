package request

import (
	"context"

	"github.com/google/uuid"
)

type contextKey string

const userIDKey = contextKey("UserID")

// WithUserID кладет в context UserID
func WithUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// UserIDFrom достает UserID from context
func UserIDFrom(ctx context.Context) (uuid.UUID, bool) {
	v, ok := ctx.Value(userIDKey).(uuid.UUID)
	return v, ok
}
