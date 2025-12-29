package events

import (
	"orders-service/internal/domain/value_objects"
	"time"

	"github.com/google/uuid"
)

type OrderStatusChangedEvent struct {
	MessageId          uuid.UUID
	OrderId            uuid.UUID
	UserId             uuid.UUID
	Status             value_objects.OrderStatus
	CancellationReason *string
	CreatedAt          time.Time
}
