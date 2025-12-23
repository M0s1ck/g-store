package get_orders

import (
	"context"
	"fmt"
	"orders-service/internal/domain/entities"
	derrors "orders-service/internal/domain/errors"

	"github.com/google/uuid"
)

type GetOrdersUsecase struct {
	repo OrderRepoGetter
}

func NewGetOrdersUsecase(repo OrderRepoGetter) *GetOrdersUsecase {
	return &GetOrdersUsecase{repo: repo}
}

func (uc *GetOrdersUsecase) GetById(
	ctx context.Context,
	orderId uuid.UUID,
	userId uuid.UUID,
) (*entities.Order, error) {

	order, err := uc.repo.GetById(ctx, orderId)
	if err != nil {
		return nil, err
	}

	if order.UserId != userId {
		return nil, fmt.Errorf("%w: user id=%s has no authority to get order id=%s",
			derrors.ErrForbidden, userId, orderId)
	}

	return order, nil
}
