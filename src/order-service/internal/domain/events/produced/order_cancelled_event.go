package published_events

import (
	"time"

	"github.com/google/uuid"

	"orders-service/internal/domain/value_objects"
)

type OrderCancelledEvent struct {
	OrderId      uuid.UUID
	UserId       uuid.UUID
	CancelReason value_objects.CancellationReason
	CancelSource CancelSource
	OccurredAt   time.Time
}

type CancelSource string

const (
	CancelSourceStore    CancelSource = "store"
	CancelSourceCustomer CancelSource = "customer"
)
