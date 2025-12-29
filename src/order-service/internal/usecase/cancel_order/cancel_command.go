package cancel_order

import (
	"github.com/google/uuid"

	"orders-service/internal/domain/value_objects"
)

type CancelOrderCommand struct {
	OrderID uuid.UUID
	Actor   CancelActor
	Reason  value_objects.CancellationReason
}
