package order_update_status

import (
	"github.com/google/uuid"

	"orders-service/internal/domain/value_objects"
)

type UpdateStatusCommand struct {
	OrderID uuid.UUID
	Status  value_objects.OrderStatus
	Actor   UpdateStatusActor
}
