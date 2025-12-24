package create_order

import (
	"context"

	"orders-service/internal/domain/messages"
)

type Repository interface {
	Create(ctx context.Context, model *messages.OutboxMessage) error
}
