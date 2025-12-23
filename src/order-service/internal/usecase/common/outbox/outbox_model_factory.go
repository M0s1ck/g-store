package outbox

import (
	"orders-service/internal/domain/events"
)

type ModelFactory interface {
	CreateOutboxModelFromOrderCreatedEvent(event *events.OrderCreatedEvent) (*Model, error)
}
