package payment_processed

import (
	"time"

	"github.com/google/uuid"
)

type PaymentProcessedEvent struct {
	MessageId            uuid.UUID
	OrderId              uuid.UUID
	UserId               uuid.UUID
	Amount               int64
	Status               PaymentStatus
	PaymentFailureReason *PaymentFailureReason
	OccurredAt           time.Time
}

type PaymentStatus string

const (
	PaymentSuccess PaymentStatus = "SUCCESS"
	PaymentFailed  PaymentStatus = "FAILED"
)

type PaymentFailureReason string

const (
	FailureNoAccount         PaymentFailureReason = "NO_ACCOUNT"
	FailureInsufficientFunds PaymentFailureReason = "INSUFFICIENT_FUNDS"
	FailureInternal          PaymentFailureReason = "INTERNAL_ERROR"
)
