package payment_processed

import (
	"context"

	"github.com/google/uuid"

	"orders-service/internal/domain/entities"
)

type OrderRepoStatusUpdater interface {
	UpdateStatus(ctx context.Context, order *entities.Order) error
	GetById(ctx context.Context, id uuid.UUID) (*entities.Order, error)
}
