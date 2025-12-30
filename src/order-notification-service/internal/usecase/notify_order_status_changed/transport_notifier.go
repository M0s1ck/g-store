package notify_order_status_changed

import (
	"context"
	"order-notification-service/internal/domain/events/consumed"
)

type TransportNotifier interface {
	NotifyOrderStatusChanged(ctx context.Context, event consumed_events.OrderStatusChangedEvent)
}
