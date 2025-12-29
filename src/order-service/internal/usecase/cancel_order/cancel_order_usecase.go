package cancel_order

import (
	"context"
	"orders-service/internal/domain/value_objects"
)

// CancelOrderUsecase implies that order is already paid and not canceled yet
type CancelOrderUsecase struct {
	repo   OrderRepoCancel
	policy *CancelPolicy
}

func NewCancelOrderUsecase(repo OrderRepoCancel, policy *CancelPolicy) *CancelOrderUsecase {
	return &CancelOrderUsecase{
		repo:   repo,
		policy: policy,
	}
}

func (uc *CancelOrderUsecase) Execute(ctx context.Context, cmd *CancelOrderCommand) error {
	order, err := uc.repo.GetById(ctx, cmd.OrderID)
	if err != nil {
		return err
	}

	err = uc.policy.CanCancel(order, cmd.Actor, cmd.Reason)
	if err != nil {
		return err
	}

	order.Status = value_objects.OrderCanceled
	order.CancellationReason = &cmd.Reason

	// TODO: outbox to payment service XDDDDD

	return uc.repo.Cancel(ctx, order)
}
