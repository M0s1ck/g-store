package order_update_status

import (
	"github.com/google/uuid"

	"orders-service/internal/domain/entities"
)

type UpdateStatusRequest struct {
	OrderID            uuid.UUID
	Status             entities.OrderStatus
	CancellationReason *string
}
