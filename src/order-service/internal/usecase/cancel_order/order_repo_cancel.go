package cancel_order

import (
	"context"

	"github.com/google/uuid"

	"orders-service/internal/domain/entities"
)

type OrderRepoCancel interface {
	Cancel(ctx context.Context, order *entities.Order) error
	GetById(ctx context.Context, id uuid.UUID) (*entities.Order, error)
}
