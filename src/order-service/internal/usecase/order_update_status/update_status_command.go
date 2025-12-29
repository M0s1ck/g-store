package order_update_status

import (
	"orders-service/internal/domain/value_objects"

	"github.com/google/uuid"
)

type UpdateStatusCommand struct {
	OrderID uuid.UUID
	Status  value_objects.OrderStatus
	Actor   UpdateStatusActor
}
