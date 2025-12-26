package outbox

import (
	"context"

	"github.com/google/uuid"

	"payment-service/internal/domain/messages"
)

type RepositoryCreator interface {
	Create(ctx context.Context, model *messages.OutboxMessage) error
}

type RepositoryPublisher interface {
	GetUnsent(ctx context.Context, limit int) ([]messages.OutboxMessage, error)
	MarkAsSent(ctx context.Context, id uuid.UUID) error
	IncrementRetry(ctx context.Context, id uuid.UUID) error
}
