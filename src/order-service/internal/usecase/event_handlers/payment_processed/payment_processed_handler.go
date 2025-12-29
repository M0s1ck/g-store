package payment_processed

import (
	"context"

	"orders-service/internal/domain/events/consumed"
	"orders-service/internal/domain/value_objects"
	"orders-service/internal/usecase/cancel_order"
	"orders-service/internal/usecase/order_update_status"
)

type PaymentProcessedEventHandler struct {
	cancelUC                  *cancel_order.CancelOrderUsecase
	updateStatusUC            *order_update_status.UpdateStatusUsecase
	payloadMapper             PayloadMapper
	paymentProcessedEventType string
}

func NewPaymentProcessedEventHandler(
	cancelUC *cancel_order.CancelOrderUsecase,
	updateStatusUC *order_update_status.UpdateStatusUsecase,
	mapper PayloadMapper,
	paymentProcessedEventType string,
) *PaymentProcessedEventHandler {

	return &PaymentProcessedEventHandler{
		paymentProcessedEventType: paymentProcessedEventType,
		cancelUC:                  cancelUC,
		updateStatusUC:            updateStatusUC,
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

	if event.Status == consumed_events.PaymentSuccess {
		updCmd := order_update_status.UpdateStatusCommand{
			OrderID: event.OrderId,
			Status:  value_objects.OrderPaid,
			Actor: order_update_status.UpdateStatusActor{
				Type: order_update_status.UpdateStatusActorPaymentService,
			},
		}

		return p.updateStatusUC.Execute(ctx, &updCmd)
	}

	reason := paymentFailureReasonToCancellationReason(*event.PaymentFailureReason)

	cmd := cancel_order.CancelOrderCommand{
		OrderID: event.OrderId,
		Reason:  reason,
		Actor: cancel_order.CancelActor{
			Type: cancel_order.CancelActorPaymentService,
		},
	}

	return p.cancelUC.Execute(ctx, &cmd)
}

func paymentFailureReasonToCancellationReason(pfr consumed_events.PaymentFailureReason) value_objects.CancellationReason {
	switch pfr {
	case consumed_events.FailureNoAccount:
		return value_objects.CancellationNoPaymentAccount
	case consumed_events.FailureInsufficientFunds:
		return value_objects.CancellationInsufficientFunds
	default:
		return value_objects.CancellationPaymentInternalError
	}
}
