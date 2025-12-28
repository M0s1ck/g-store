package order_update_status

import (
	"context"
	"time"

	"github.com/google/uuid"

	"orders-service/internal/domain/entities"
	"orders-service/internal/domain/events"
	"orders-service/internal/usecase/common"
	"orders-service/internal/usecase/common/outbox"
)

type UpdateStatusUsecase struct {
	txManager        common.TxManager
	orderRepo        OrderRepoStatusUpdater
	outboxRepo       outbox.RepositoryCreator
	outboxMsgFactory *outbox.MessageFactory
}

func NewUpdateOrderStatusUsecase(
	txManager common.TxManager,
	orderRepo OrderRepoStatusUpdater,
	outboxRepo outbox.RepositoryCreator,
	outboxMsgFactory *outbox.MessageFactory,
) *UpdateStatusUsecase {

	return &UpdateStatusUsecase{
		orderRepo:        orderRepo,
		outboxRepo:       outboxRepo,
		txManager:        txManager,
		outboxMsgFactory: outboxMsgFactory,
	}
}

func (uc *UpdateStatusUsecase) Execute(ctx context.Context, request UpdateStatusRequest) error {
	return uc.txManager.WithinTx(ctx, func(ctx context.Context) error {
		order, err := uc.orderRepo.GetById(ctx, request.OrderID)
		if err != nil {
			return err
		}

		order.Status = request.Status
		order.CancellationReason = request.CancellationReason

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

func toOrderStatusChangedEvent(order *entities.Order) *events.OrderStatusChangedEvent {
	return &events.OrderStatusChangedEvent{
		MessageId:          uuid.New(),
		OrderId:            order.Id,
		UserId:             order.UserId,
		Status:             order.Status,
		CancellationReason: order.Description,
		CreatedAt:          time.Now(),
	}
}
