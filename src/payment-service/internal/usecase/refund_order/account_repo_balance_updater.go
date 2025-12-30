package refund_order

import (
	"context"

	"github.com/google/uuid"

	"payment-service/internal/domain/entities"
)

type AccountRepoBalanceUpdater interface {
	GetByUserId(ctx context.Context, id uuid.UUID) (*entities.Account, error)
	GetById(ctx context.Context, id uuid.UUID) (*entities.Account, error)
	UpdateBalance(ctx context.Context, account *entities.Account) error
}
