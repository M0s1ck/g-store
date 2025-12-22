package middleware

import (
	"context"

	"github.com/google/uuid"
)

func UUIDFromContext(ctx context.Context) uuid.UUID {
	id, ok := ctx.Value(ctxKeyUUID).(uuid.UUID)
	if !ok {
		panic("uuid missing from context")
	}

	return id
}
