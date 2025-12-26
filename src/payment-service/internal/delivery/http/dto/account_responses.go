package dto

import (
	"time"

	"github.com/google/uuid"
)

type AccountResponse struct {
	ID        uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	UserID    uuid.UUID `json:"userId" example:"550e8400-e29b-41d4-a716-446655440000"`
	Balance   int64     `json:"balance" example:"19999"`
	CreatedAt time.Time `json:"createdAt" example:"2025-12-23T09:19:23.458426Z"`
	UpdatedAt time.Time `json:"updatedAt" example:"2025-12-23T09:19:23.458426Z"`
}

type AccountCreatedResponse struct {
	AccountID uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
}
