package entities

import (
	"orders-service/internal/domain/value_objects"
	"time"

	"github.com/google/uuid"
)

type Order struct {
	Id                 uuid.UUID                         `db:"id"`
	UserId             uuid.UUID                         `db:"user_id"`
	Amount             int64                             `db:"amount"`
	Status             value_objects.OrderStatus         `db:"status"`
	CancellationReason *value_objects.CancellationReason `db:"cancellation_reason"`
	Description        *string                           `db:"description"`
	CreatedAt          time.Time                         `db:"created_at"`
	UpdatedAt          time.Time                         `db:"updated_at"`
}
