package get_orders

import (
	"context"

	"github.com/google/uuid"

	"orders-service/internal/domain/entities"
)

type GetByUserUsecase struct {
	repo OrderRepoGetter
}

func NewGetByUserUsecase(repo OrderRepoGetter) *GetByUserUsecase {
	return &GetByUserUsecase{repo: repo}
}

func (uc *GetByUserUsecase) Execute(
	ctx context.Context,
	userId uuid.UUID,
	page, limit int,
) ([]entities.Order, int, error) {

	// TODO: maybe add redis for this later
	return uc.repo.GetByUserId(ctx, userId, page, limit)
}
