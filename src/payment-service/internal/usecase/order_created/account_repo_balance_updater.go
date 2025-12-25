package order_created

import (
	"context"
	"payment-service/internal/domain/entities"

	"github.com/google/uuid"
)

type AccountRepoBalanceUpdater interface {
	GetByUserId(ctx context.Context, id uuid.UUID) (*entities.Account, error)
	UpdateBalance(ctx context.Context, account *entities.Account) error
}
