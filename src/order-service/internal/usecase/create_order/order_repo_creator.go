package create_order

import (
	"context"

	"orders-service/internal/domain/entities"
)

type OrderRepoCreator interface {
	Create(ctx context.Context, order *entities.Order) error
}
