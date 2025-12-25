package payment_processed

import (
	"context"
	"orders-service/internal/domain/entities"

	"github.com/google/uuid"
)

type OrderRepoStatusUpdater interface {
	UpdateStatus(ctx context.Context, order *entities.Order) error
	GetById(ctx context.Context, id uuid.UUID) (*entities.Order, error)
}
