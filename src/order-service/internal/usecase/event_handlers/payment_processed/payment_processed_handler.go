package payment_processed

import (
	"context"

	"orders-service/internal/domain/entities"
	"orders-service/internal/usecase/order_update_status"
)

const (
	successDesc = "the order is paid"
)

type PaymentProcessedEventHandler struct {
	updateStatusUC            *order_update_status.UpdateStatusUsecase
	payloadMapper             PayloadMapper
	paymentProcessedEventType string
}

func NewPaymentProcessedEventHandler(
	updateStatus *order_update_status.UpdateStatusUsecase,
	mapper PayloadMapper,
	paymentProcessedEventType string,
) *PaymentProcessedEventHandler {

	return &PaymentProcessedEventHandler{
		paymentProcessedEventType: paymentProcessedEventType,
		updateStatusUC:            updateStatus,
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

	updReq := order_update_status.UpdateStatusRequest{
		OrderID: event.OrderId,
	}

	if event.Status == PaymentSuccess {
		updReq.Status = entities.OrderPaid
	} else {
		updReq.Status = entities.OrderCanceled
		updReq.CancellationReason = (*string)(event.PaymentFailureReason)
	}

	return p.updateStatusUC.Execute(ctx, updReq)
}
