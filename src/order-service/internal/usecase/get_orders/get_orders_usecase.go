package get_orders

import (
	"context"
	"orders-service/internal/domain/entities"

	"github.com/google/uuid"
)

type GetOrdersUsecase struct {
	repo OrderRepoGetter
}

func NewGetOrdersUsecase(repo OrderRepoGetter) *GetOrdersUsecase {
	return &GetOrdersUsecase{repo: repo}
}

func (uc *GetOrdersUsecase) GetById(ctx context.Context, id uuid.UUID) (*entities.Order, error) {
	return uc.repo.GetById(ctx, id)
}
