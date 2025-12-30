package dto

type CreateOrderRequest struct {
	Amount int64 `json:"amount"`
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status" example:"ASSEMBLING"`
}

type ConsumerCancelOrderRequest struct {
	CancellationReason string `json:"cancellationReason,omitempty" example:"CHANGED_MIND"`
}

type StaffCancelOrderRequest struct {
	CancellationReason string `json:"cancellationReason" example:"OUT_OF_STOCK"`
}
