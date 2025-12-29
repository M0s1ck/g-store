package common_outbox

import (
	"orders-service/internal/domain/events/produced"
)

type OrderCreatedEventPayloadMapper interface {
	OrderCreatedEventToPayload(event *published_events.OrderCreatedEvent) ([]byte, error)
}

type OrderStatusChangedPayloadMapper interface {
	OrderStatusChangedEventToPayload(event *published_events.OrderStatusChangedEvent) ([]byte, error)
}

type OrderCancelledPayloadMapper interface {
	OrderCancelledEventToPayload(event *published_events.OrderCancelledEvent) ([]byte, error)
}
