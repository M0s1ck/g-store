package payment_processed

import (
	"context"
	"log"

	"orders-service/internal/domain/entities"
)

const (
	successDesc = "the order is paid"
)

type PaymentProcessedEventHandler struct {
	orderRepo                 OrderRepoStatusUpdater
	paymentProcessedEventType string
	payloadMapper             PayloadMapper
}

func NewPaymentProcessedEventHandler(
	paymentProcessedEventType string,
	orderRepo OrderRepoStatusUpdater,
	mapper PayloadMapper,
) *PaymentProcessedEventHandler {

	return &PaymentProcessedEventHandler{
		paymentProcessedEventType: paymentProcessedEventType,
		orderRepo:                 orderRepo,
		payloadMapper:             mapper,
	}
}

func (p *PaymentProcessedEventHandler) EventType() string {
	return p.paymentProcessedEventType
}

func (p *PaymentProcessedEventHandler) Handle(ctx context.Context, payload []byte) error {
	event, err := p.payloadMapper.ToPaymentProcessedEvent(payload)
	if err != nil {
		return err
	}

	order, err := p.orderRepo.GetById(ctx, event.OrderId)
	if err != nil {
		log.Printf("Could not get order by id: %s: %s", event.OrderId, err)
		return err
	}

	if event.Status == PaymentSuccess {
		order.Status = entities.OrderPaid
		return p.orderRepo.UpdateStatus(ctx, order)
	}

	order.Description = (*string)(event.PaymentFailureReason)
	order.Status = entities.OrderCanceled
	return p.orderRepo.UpdateStatus(ctx, order)
}
