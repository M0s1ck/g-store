package get_account

import (
	"context"

	"github.com/google/uuid"

	"payment-service/internal/domain/entities"
	myerrors "payment-service/internal/domain/errors"
)

type GetByIdUsecase struct {
	repo AccountRepoGetter
}

func NewGetByIdUsecase(repo AccountRepoGetter) *GetByIdUsecase {
	return &GetByIdUsecase{repo: repo}
}

func (u *GetByIdUsecase) Execute(ctx context.Context, accID uuid.UUID, userID uuid.UUID) (*entities.Account, error) {
	account, err := u.repo.GetById(ctx, accID)
	if err != nil {
		return nil, err
	}

	if account.UserID != userID {
		return nil, myerrors.ErrForbidden
	}

	return account, nil
}
