package mapper

import (
	"orders-service/internal/delivery/http/dto"
	"orders-service/internal/domain/entities"
	"orders-service/internal/usecase/create_order"
)

func OrderToResponse(order *entities.Order) *dto.OrderResponse {
	return &dto.OrderResponse{
		Id:        order.Id,
		UserId:    order.UserId,
		Amount:    order.Amount,
		Status:    order.Status,
		CreatedAt: order.CreatedAt,
		UpdatedAt: order.UpdatedAt,
	}
}

func OrdersToResponse(orders []entities.Order) []dto.OrderResponse {
	res := make([]dto.OrderResponse, len(orders))
	for i := range orders {
		res[i] = *OrderToResponse(&orders[i])
	}

	return res
}

func OrderCreateRequestToApplication(req dto.CreateOrderRequest) *create_order.CreateOrderRequest {
	return &create_order.CreateOrderRequest{
		Amount: req.Amount,
	}
}

func OrderCreatedResponseToDto(resp *create_order.CreateOrderResponse) *dto.OrderCreatedResponse {
	return &dto.OrderCreatedResponse{
		Id: resp.Id,
	}
}
