package notify_order_status_changed

import "github.com/google/uuid"

type TransportNotifier interface {
	NotifyOrder(orderId uuid.UUID, payload []byte)
}
