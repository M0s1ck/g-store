package get_orders

import (
	"context"
	"orders-service/internal/domain/entities"

	"github.com/google/uuid"
)

type OrderRepoGetter interface {
	GetById(ctx context.Context, id uuid.UUID) (*entities.Order, error)
}
