package order_status_changed

import (
	"time"

	"github.com/google/uuid"

	"order-notification-service/internal/usecase/domain/value_objects"
)

type Event struct {
	MessageId          uuid.UUID
	OrderId            uuid.UUID
	UserId             uuid.UUID
	Status             value_objects.OrderStatus
	CancellationReason *string
	CreatedAt          time.Time
}
