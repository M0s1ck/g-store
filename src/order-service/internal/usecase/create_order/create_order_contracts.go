package create_order

import "github.com/google/uuid"

type CreateOrderRequest struct {
	Amount int64
}

type CreateOrderResponse struct {
	Id uuid.UUID
}
