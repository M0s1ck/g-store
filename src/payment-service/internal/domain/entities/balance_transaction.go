package entities

import (
	"time"

	"github.com/google/uuid"
)

type BalanceTransaction struct {
	ID        uuid.UUID       `db:"id"`
	AccountID uuid.UUID       `db:"account_id"`
	OrderID   *uuid.UUID      `db:"order_id"` // null for top-ups
	Amount    int64           `db:"amount"`
	Direction Direction       `db:"direction"`
	Type      TransactionType `db:"type"`
	CreatedAt time.Time       `db:"created_at"`
}

type TransactionType string

const (
	TransactionTopUp   TransactionType = "TOP_UP"
	TransactionPayment TransactionType = "PAYMENT"
	TransactionRefund  TransactionType = "REFUND"
)

type Direction string

const (
	DirectionIn  Direction = "IN"
	DirectionOut Direction = "OUT"
)
