package outbox_publish

import (
	"context"

	"github.com/google/uuid"

	"orders-service/internal/domain/messages"
)

type OutboxRepoPublisher interface {
	GetUnsent(ctx context.Context, limit int) ([]messages.OutboxMessage, error)
	MarkAsSent(ctx context.Context, id uuid.UUID) error
	IncrementRetry(ctx context.Context, id uuid.UUID) error
}
