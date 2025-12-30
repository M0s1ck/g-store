package consumed_events

import (
	"time"

	"github.com/google/uuid"
)

type OrderCancelledEvent struct {
	OrderId      uuid.UUID
	UserId       uuid.UUID
	CancelReason OrderCancellationReason
	CancelSource CancelSource
	OccurredAt   time.Time
}

type CancelSource string

const (
	CancelSourceStore    CancelSource = "store"
	CancelSourceCustomer CancelSource = "customer"
)

type OrderCancellationReason string

const (
	CancellationNoPaymentAccount     OrderCancellationReason = "NO_PAYMENT_ACCOUNT"
	CancellationInsufficientFunds    OrderCancellationReason = "INSUFFICIENT_FUNDS"
	CancellationPaymentInternalError OrderCancellationReason = "PAYMENT_INTERNAL_ERROR"
	CancellationOutOfStock           OrderCancellationReason = "OUT_OF_STOCK"
	CancellationDeliveryUnavailable  OrderCancellationReason = "DELIVERY_UNAVAILABLE"
	CancellationChangedMind          OrderCancellationReason = "CHANGED_MIND"
)
