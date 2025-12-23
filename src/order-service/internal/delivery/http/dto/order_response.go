package dto

import (
	"orders-service/internal/domain/entities"
	"time"

	"github.com/google/uuid"
)

type OrderResponse struct {
	Id        uuid.UUID
	UserId    uuid.UUID
	Amount    int64
	Status    entities.OrderStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}
