package published_events

import (
	"time"

	"github.com/google/uuid"
)

type OrderCreatedEvent struct {
	OrderId    uuid.UUID
	UserId     uuid.UUID
	Amount     int64
	OccurredAt time.Time
}
