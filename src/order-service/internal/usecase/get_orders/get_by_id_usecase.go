package get_orders

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"orders-service/internal/domain/entities"
	derrors "orders-service/internal/domain/errors"
)

type GetByIdUsecase struct {
	repo OrderRepoGetter
}

func NewGetByIdUsecase(repo OrderRepoGetter) *GetByIdUsecase {
	return &GetByIdUsecase{repo: repo}
}

func (uc *GetByIdUsecase) Execute(
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
