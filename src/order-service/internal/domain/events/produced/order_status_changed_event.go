package published_events

import (
	"time"

	"github.com/google/uuid"

	"orders-service/internal/domain/value_objects"
)

type OrderStatusChangedEvent struct {
	OrderId    uuid.UUID
	UserId     uuid.UUID
	Status     value_objects.OrderStatus
	OccurredAt time.Time
}
