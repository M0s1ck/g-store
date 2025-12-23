package mapper

import (
	"orders-service/internal/delivery/http/dto"
	"orders-service/internal/domain/entities"
)

func ToResponse(order *entities.Order) *dto.OrderResponse {
	return &dto.OrderResponse{
		Id:        order.Id,
		UserId:    order.UserId,
		Amount:    order.Amount,
		Status:    order.Status,
		CreatedAt: order.CreatedAt,
		UpdatedAt: order.UpdatedAt,
	}
}
