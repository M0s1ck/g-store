package events

import (
	"time"

	"github.com/google/uuid"

	"orders-service/internal/domain/entities"
)

type OrderStatusChangedEvent struct {
	MessageId          uuid.UUID
	OrderId            uuid.UUID
	UserId             uuid.UUID
	Status             entities.OrderStatus
	CancellationReason *string
	CreatedAt          time.Time
}
