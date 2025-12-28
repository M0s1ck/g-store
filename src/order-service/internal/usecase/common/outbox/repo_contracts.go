package outbox

import (
	"context"

	"orders-service/internal/domain/messages"
)

type RepositoryCreator interface {
	Create(ctx context.Context, model *messages.OutboxMessage) error
}
