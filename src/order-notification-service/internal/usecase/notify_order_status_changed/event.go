package notify_order_status_changed

import (
	"github.com/google/uuid"

	"order-notification-service/internal/usecase/domain/value_objects"
)

// Event clean, for usecase
type Event struct {
	OrderID            uuid.UUID
	UserID             uuid.UUID
	Status             value_objects.OrderStatus
	CancellationReason *string
}
