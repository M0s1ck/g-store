package order_cancelled

import (
	"context"
	"log"

	"payment-service/internal/domain/events/consumed"
	"payment-service/internal/domain/messages"
	"payment-service/internal/usecase/refund_order"
)

type CancelEventHandler struct {
	refundUC                *refund_order.RefundUsecase
	orderMapper             PayloadMapper
	orderCancelledEventType string
}

func NewEventHandler(
	refundUC *refund_order.RefundUsecase,
	orderMapper PayloadMapper,
	orderCancelledEventType string,
) *CancelEventHandler {

	return &CancelEventHandler{
		refundUC:                refundUC,
		orderMapper:             orderMapper,
		orderCancelledEventType: orderCancelledEventType,
	}
}

func (eh *CancelEventHandler) Handle(ctx context.Context, msg messages.InboxMessage) error {
	log.Printf("Got msg: %v", msg.Topic)

	event, err := eh.orderMapper.ToOrderCancelledEvent(msg.Payload)
	if err != nil {
		return err
	}

	// Here we can add logic to determine which usecase to use
	// For now we will just refund regardless

	refundCmd := eh.refundCmdFromEvent(event)
	err = eh.refundUC.Execute(ctx, refundCmd)

	if err != nil {
		// Rn do nothing
		log.Printf("Error executing refund usecase: %v", err)
	}
	return nil
}

func (eh *CancelEventHandler) EventType() string {
	return eh.orderCancelledEventType
}

func (eh *CancelEventHandler) refundCmdFromEvent(event *consumed_events.OrderCancelledEvent) *refund_order.Command {
	return &refund_order.Command{
		OrderId: event.OrderId,
	}
}
