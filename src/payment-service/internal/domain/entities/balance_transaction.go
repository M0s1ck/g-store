package entities

import (
	"time"

	"github.com/google/uuid"
)

type BalanceTransaction struct {
	ID        uuid.UUID       `db:"id"`
	AccountID uuid.UUID       `db:"account_id"`
	OrderID   *uuid.UUID      `db:"order_id"` // null for top-ups
	Amount    int64           `db:"amount"`   // + / -
	Type      TransactionType `db:"type"`
	CreatedAt time.Time       `db:"created_at"`
}

type TransactionType string

const (
	TransactionTopUp   TransactionType = "TOP_UP"
	TransactionPayment TransactionType = "PAYMENT"
)
