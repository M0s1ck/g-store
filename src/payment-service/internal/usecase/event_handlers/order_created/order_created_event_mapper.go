package order_created

import "payment-service/internal/domain/events/consumed"

type OrderCreatedEventMapper interface {
	ToOrderCreatedEvent(payload []byte) (*consumed_events.OrderCreatedEvent, error)
}
