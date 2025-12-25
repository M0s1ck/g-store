package create_order

import (
	"context"
	"orders-service/internal/domain/errors"
	"time"

	"github.com/google/uuid"

	"orders-service/internal/domain/entities"
	"orders-service/internal/domain/events"
	"orders-service/internal/usecase/common"
)

type CreateOrderUsecase struct {
	txManager          common.TxManager
	orderRepo          OrderRepoCreator
	outboxRepo         Repository
	outboxModelFactory OutboxMessageFactory
}

func NewCreateOrderUsecase(
	txManager common.TxManager,
	repo OrderRepoCreator,
	outboxRepo Repository,
	outboxModelFactory OutboxMessageFactory,
) *CreateOrderUsecase {

	return &CreateOrderUsecase{
		orderRepo:          repo,
		txManager:          txManager,
		outboxRepo:         outboxRepo,
		outboxModelFactory: outboxModelFactory,
	}
}

func (uc *CreateOrderUsecase) Execute(
	ctx context.Context,
	request *CreateOrderRequest,
	userId uuid.UUID,
) (*CreateOrderResponse, error) {

	order := uc.getOrderFromRequest(request, userId)

	if order.Amount <= 0 {
		return nil, errors.ErrAmountNotPositive
	}

	err := uc.txManager.WithinTx(ctx, func(ctx context.Context) error {

		err := uc.orderRepo.Create(ctx, order)
		if err != nil {
			return err
		}

		event := uc.getEventFromOrder(order)
		outboxModel, err := uc.outboxModelFactory.CreateOutboxModelFromOrderCreatedEvent(event)

		err = uc.outboxRepo.Create(ctx, outboxModel)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &CreateOrderResponse{Id: order.Id}, nil
}

func (uc *CreateOrderUsecase) getOrderFromRequest(req *CreateOrderRequest, userId uuid.UUID) *entities.Order {
	return &entities.Order{
		Id:     uuid.New(),
		UserId: userId,
		Amount: req.Amount,
		Status: entities.OrderPending,
	}
}

func (uc *CreateOrderUsecase) getEventFromOrder(order *entities.Order) *events.OrderCreatedEvent {
	return &events.OrderCreatedEvent{
		MessageId: uuid.New(),
		OrderId:   order.Id,
		UserId:    order.UserId,
		Amount:    order.Amount,
		CreatedAt: time.Now(),
	}
}
