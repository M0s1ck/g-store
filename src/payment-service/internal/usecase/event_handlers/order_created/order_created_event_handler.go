package order_created

import (
	"context"
	"log"

	"payment-service/internal/domain/events/consumed"
	"payment-service/internal/domain/messages"
	"payment-service/internal/usecase/pay_for_order"
)

type OrderCreatedEventHandler struct {
	payUC          *pay_for_order.PayUsecase
	orderMapper    OrderCreatedEventMapper
	orderCrEvtType string
}

func NewEventHandler(
	payUC *pay_for_order.PayUsecase,
	mapper OrderCreatedEventMapper,
	orderCrEvtType string,
) *OrderCreatedEventHandler {

	return &OrderCreatedEventHandler{
		payUC:          payUC,
		orderMapper:    mapper,
		orderCrEvtType: orderCrEvtType,
	}
}

func (eh *OrderCreatedEventHandler) EventType() string {
	return eh.orderCrEvtType
}

func (eh *OrderCreatedEventHandler) Handle(ctx context.Context, msg messages.InboxMessage) error {
	log.Printf("Got msg: %v", msg.Topic)

	orderEvent, err := eh.orderMapper.ToOrderCreatedEvent(msg.Payload)
	if err != nil {
		return err
	}

	payCmd := payCmdFromEvent(orderEvent)
	return eh.payUC.Execute(ctx, payCmd)
}

func payCmdFromEvent(created *consumed_events.OrderCreatedEvent) *pay_for_order.Command {
	return &pay_for_order.Command{
		OrderId: created.OrderId,
		UserId:  created.UserId,
		Amount:  created.Amount,
	}
}
