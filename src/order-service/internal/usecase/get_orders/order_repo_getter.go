package get_orders

import (
	"context"

	"github.com/google/uuid"

	"orders-service/internal/domain/entities"
)

type OrderRepoGetter interface {
	GetById(ctx context.Context, id uuid.UUID) (*entities.Order, error)
	GetByUserId(ctx context.Context, userId uuid.UUID, page, limit int) ([]entities.Order, int, error)
}
