package order_created

import (
	"context"
	"log"

	"payment-service/internal/domain/messages"
)

type OrderCreatedEventHandler struct {
	orderCreatedEventType string
}

func NewOrderCreatedEventHandler(orderCreatedEventType string) *OrderCreatedEventHandler {
	return &OrderCreatedEventHandler{
		orderCreatedEventType: orderCreatedEventType,
	}
}

func (o OrderCreatedEventHandler) EventType() string {
	return o.orderCreatedEventType
}

func (o OrderCreatedEventHandler) Handle(ctx context.Context, msg messages.InboxMessage) error {
	//TODO implement me
	log.Printf("Got msg!!! : %v %v", msg.Topic, msg.Payload)
	return nil
}
