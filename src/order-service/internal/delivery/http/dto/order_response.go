package dto

import (
	"time"

	"github.com/google/uuid"

	"orders-service/internal/domain/entities"
)

type OrderResponse struct {
	Id        uuid.UUID            `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	UserId    uuid.UUID            `json:"userId" example:"550e8400-e29b-41d4-a716-446655440000"`
	Amount    int64                `json:"amount" example:"19999"`
	Status    entities.OrderStatus `json:"status" example:"PENDING"`
	CreatedAt time.Time            `json:"createdAt" example:"2025-12-23T09:19:23.458426Z"`
	UpdatedAt time.Time            `json:"updatedAt" example:"2025-12-23T09:19:23.458426Z"`
}

type OrdersResponse struct {
	Orders []OrderResponse `json:"orders"`
	Total  int             `json:"total"`
	Page   int             `json:"page"`
	Limit  int             `json:"limit"`
}

type OrderCreatedResponse struct {
	Id uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
}
