package events

import (
	"time"

	"github.com/google/uuid"
)

type OrderCreatedEvent struct {
	MessageId uuid.UUID
	OrderId   uuid.UUID
	UserId    uuid.UUID
	Amount    int64
	CreatedAt time.Time
}
