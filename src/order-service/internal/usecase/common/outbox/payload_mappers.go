package outbox

import "orders-service/internal/domain/events"

type OrderCreatedEventPayloadMapper interface {
	OrderCreatedEventToPayload(event *events.OrderCreatedEvent) ([]byte, error)
}

type OrderStatusChangedPayloadMapper interface {
	OrderStatusChangedEventToPayload(event *events.OrderStatusChangedEvent) ([]byte, error)
}
