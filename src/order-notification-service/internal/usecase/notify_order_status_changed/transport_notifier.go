package notify_order_status_changed

import "context"

type TransportNotifier interface {
	NotifyOrderStatusChanged(ctx context.Context, event Event)
}
