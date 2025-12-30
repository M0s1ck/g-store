package order_cancelled

import "payment-service/internal/domain/events/consumed"

type PayloadMapper interface {
	ToOrderCancelledEvent(payload []byte) (*consumed_events.OrderCancelledEvent, error)
}
