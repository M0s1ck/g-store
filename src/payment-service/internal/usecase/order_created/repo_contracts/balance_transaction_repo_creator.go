package repo_contracts

import (
	"context"
	"payment-service/internal/domain/entities"
)

type BalanceTransactionRepoCreator interface {
	Create(ctx context.Context, transaction *entities.BalanceTransaction) error
}
