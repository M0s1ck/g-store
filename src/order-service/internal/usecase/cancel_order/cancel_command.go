package cancel_order

import (
	"orders-service/internal/domain/value_objects"

	"github.com/google/uuid"
)

type CancelOrderCommand struct {
	OrderID uuid.UUID
	Actor   CancelActor
	Reason  value_objects.CancellationReason
}
