package get_account

import (
	"context"

	"github.com/google/uuid"

	"payment-service/internal/domain/entities"
)

type AccountRepoGetter interface {
	GetById(ctx context.Context, id uuid.UUID) (*entities.Account, error)
}
