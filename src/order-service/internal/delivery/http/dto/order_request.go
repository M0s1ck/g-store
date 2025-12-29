package dto

type CreateOrderRequest struct {
	Amount int64 `json:"amount"`
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status" example:"ASSEMBLING"`
}

type ConsumerCancelOrderRequest struct {
	CancellationReason string `json:"cancellation_reason,omitempty" example:"CHANGED_MIND"`
}

type StaffCancelOrderRequest struct {
	CancellationReason string `json:"cancellation_reason" example:"OUT_OF_STOCK"`
}
