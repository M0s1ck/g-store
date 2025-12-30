package consumed_events

import (
	"time"

	"github.com/google/uuid"
)

type OrderStatusChangedEvent struct {
	OrderId    uuid.UUID
	UserId     uuid.UUID
	Status     OrderStatus
	OccurredAt time.Time
}

type OrderStatus string

const (
	OrderPending    OrderStatus = "PENDING"
	OrderPaid       OrderStatus = "PAID"
	OrderCanceled   OrderStatus = "CANCELED"
	OrderAssembling OrderStatus = "ASSEMBLING"
	OrderAssembled  OrderStatus = "ASSEMBLED"
	OrderDelivering OrderStatus = "DELIVERING"
	OrderIssued     OrderStatus = "ISSUED"
)
