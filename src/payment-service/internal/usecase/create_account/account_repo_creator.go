package create_account

import (
	"context"

	"payment-service/internal/domain/entities"
)

type AccountRepoCreator interface {
	Create(ctx context.Context, account *entities.Account) error
}
