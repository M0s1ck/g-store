package create_order

import (
	"orders-service/internal/domain/events"
	"orders-service/internal/domain/messages"
)

type OutboxMessageFactory interface {
	CreateOutboxModelFromOrderCreatedEvent(event *events.OrderCreatedEvent) (*messages.OutboxMessage, error)
}
