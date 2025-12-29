package cancel_order

import (
	"context"
	"orders-service/internal/domain/entities"

	"github.com/google/uuid"
)

type OrderRepoCancel interface {
	Cancel(ctx context.Context, order *entities.Order) error
	GetById(ctx context.Context, id uuid.UUID) (*entities.Order, error)
}
