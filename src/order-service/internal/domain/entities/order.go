package entities

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	Id                 uuid.UUID   `db:"id"`
	UserId             uuid.UUID   `db:"user_id"`
	Amount             int64       `db:"amount"`
	Status             OrderStatus `db:"status"`
	CancellationReason *string     `db:"cancellation_reason"`
	Description        *string     `db:"description"`
	CreatedAt          time.Time   `db:"created_at"`
	UpdatedAt          time.Time   `db:"updated_at"`
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
