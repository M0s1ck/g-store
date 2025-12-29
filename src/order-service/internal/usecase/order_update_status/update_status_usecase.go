package order_update_status

import (
	"context"
	"time"

	"orders-service/internal/domain/entities"
	myerrors "orders-service/internal/domain/errors"
	"orders-service/internal/domain/events/produced"
	"orders-service/internal/usecase/common"
	"orders-service/internal/usecase/common/outbox"
)

type UpdateStatusUsecase struct {
	txManager        common.TxManager
	orderRepo        OrderRepoStatusUpdater
	outboxRepo       common_outbox.RepositoryCreator
	policy           *UpdateStatusPolicy
	outboxMsgFactory *common_outbox.MessageFactory
}

func NewUpdateOrderStatusUsecase(
	txManager common.TxManager,
	orderRepo OrderRepoStatusUpdater,
	outboxRepo common_outbox.RepositoryCreator,
	policy *UpdateStatusPolicy,
	outboxMsgFactory *common_outbox.MessageFactory,
) *UpdateStatusUsecase {

	return &UpdateStatusUsecase{
		orderRepo:        orderRepo,
		outboxRepo:       outboxRepo,
		txManager:        txManager,
		outboxMsgFactory: outboxMsgFactory,
		policy:           policy,
	}
}

func (uc *UpdateStatusUsecase) Execute(ctx context.Context, cmd *UpdateStatusCommand) error {
	return uc.txManager.WithinTx(ctx, func(ctx context.Context) error {
		order, err := uc.orderRepo.GetById(ctx, cmd.OrderID)
		if err != nil {
			return err
		}

		if !order.Status.CanTransitionTo(cmd.Status) {
			return myerrors.ErrInvalidOrderStatusChange
		}

		err = uc.policy.CanUpdateStatus(order, cmd.Actor, cmd.Status)
		if err != nil {
			return err
		}

		order.Status = cmd.Status

		err = uc.orderRepo.UpdateStatus(ctx, order)
		if err != nil {
			return err
		}

		statusEvent := toOrderStatusChangedEvent(order)
		outboxMsg, err := uc.outboxMsgFactory.CreateMessageOrderStatusChangedEvent(statusEvent)
		if err != nil {
			return err
		}

		return uc.outboxRepo.Create(ctx, outboxMsg)
	})

}

func toOrderStatusChangedEvent(order *entities.Order) *published_events.OrderStatusChangedEvent {
	return &published_events.OrderStatusChangedEvent{
		OrderId:    order.Id,
		UserId:     order.UserId,
		Status:     order.Status,
		OccurredAt: time.Now(),
	}
}
