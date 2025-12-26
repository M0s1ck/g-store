package create_account

import (
	"context"

	"github.com/google/uuid"

	"payment-service/internal/domain/entities"
)

type CreateAccountUsecase struct {
	repo AccountRepoCreator
}

func NewCreateAccountUsecase(repo AccountRepoCreator) *CreateAccountUsecase {
	return &CreateAccountUsecase{repo: repo}
}

func (u *CreateAccountUsecase) Execute(ctx context.Context, userId uuid.UUID) (*Response, error) {
	accId := uuid.New()

	acc := &entities.Account{
		ID:     accId,
		UserID: userId,
	}

	err := u.repo.Create(ctx, acc)
	if err != nil {
		return nil, err
	}

	return &Response{ID: accId}, nil
}
