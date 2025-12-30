package order_status_changed

import "order-notification-service/internal/domain/events/consumed"

type PayloadMapper interface {
	OrderStatusChangedEventFromPayload(payload []byte) (*consumed_events.OrderStatusChangedEvent, error)
}
