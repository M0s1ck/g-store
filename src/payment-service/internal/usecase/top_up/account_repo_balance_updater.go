package top_up

import (
	"context"

	"github.com/google/uuid"

	"payment-service/internal/domain/entities"
)

type AccountRepoBalanceUpdater interface {
	GetById(ctx context.Context, id uuid.UUID) (*entities.Account, error)
	UpdateBalance(ctx context.Context, account *entities.Account) error
}
