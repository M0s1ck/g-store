package cancel_order

import (
	"context"
	"log"
	"time"

	"orders-service/internal/domain/entities"
	"orders-service/internal/domain/events/produced"
	"orders-service/internal/domain/value_objects"
	"orders-service/internal/usecase/common"
	"orders-service/internal/usecase/common/outbox"
)

// CancelOrderUsecase implies that order is already paid and not canceled yet
type CancelOrderUsecase struct {
	orderRepo        OrderRepoCancel
	outboxRepo       common_outbox.RepositoryCreator
	txManager        common.TxManager
	outboxMsgFactory *common_outbox.MessageFactory
	policy           *CancelPolicy
}

func NewCancelOrderUsecase(
	orderRepo OrderRepoCancel,
	outboxRepo common_outbox.RepositoryCreator,
	txManager common.TxManager,
	outboxMsgFactory *common_outbox.MessageFactory,
	policy *CancelPolicy,
) *CancelOrderUsecase {

	return &CancelOrderUsecase{
		orderRepo:        orderRepo,
		outboxRepo:       outboxRepo,
		txManager:        txManager,
		outboxMsgFactory: outboxMsgFactory,
		policy:           policy,
	}
}

func (uc *CancelOrderUsecase) Execute(ctx context.Context, cmd *CancelOrderCommand) error {
	return uc.txManager.WithinTx(ctx, func(ctx context.Context) error {
		order, err := uc.orderRepo.GetById(ctx, cmd.OrderID)
		if err != nil {
			return err
		}

		err = uc.policy.CanCancel(order, cmd.Actor, cmd.Reason)
		if err != nil {
			return err
		}

		order.Status = value_objects.OrderCanceled
		order.CancellationReason = &cmd.Reason

		err = uc.orderRepo.Cancel(ctx, order)
		if err != nil {
			log.Printf("cancel order %v failed: %v", order.Id, err)
			return err
		}

		event := uc.createOrderCancelledEvent(order, cmd)
		outboxMsg, err := uc.outboxMsgFactory.CreateMessageOrderCancelledEvent(event)
		if err != nil {
			return err
		}

		err = uc.outboxRepo.Create(ctx, outboxMsg)
		if err != nil {
			log.Printf("outbox msg for event %v failed: %v", outboxMsg.EventType, err)
			return err
		}

		return nil
	})
}

func (uc *CancelOrderUsecase) createOrderCancelledEvent(
	order *entities.Order,
	cmd *CancelOrderCommand,
) *published_events.OrderCancelledEvent {

	var cancelSrc published_events.CancelSource

	switch cmd.Actor.Type {
	case CancelActorStore,
		CancelActorPaymentService:
		cancelSrc = published_events.CancelSourceStore
	case CancelActorCustomer:
		cancelSrc = published_events.CancelSourceCustomer
	}

	return &published_events.OrderCancelledEvent{
		OrderId:      order.Id,
		UserId:       order.UserId,
		CancelReason: cmd.Reason,
		CancelSource: cancelSrc,
		OccurredAt:   time.Now(),
	}
}
