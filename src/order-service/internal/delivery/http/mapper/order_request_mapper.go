package mapper

import (
	"orders-service/internal/delivery/http/dto"
	myerrors "orders-service/internal/domain/errors"
	"orders-service/internal/domain/value_objects"
	"orders-service/internal/usecase/cancel_order"
	"orders-service/internal/usecase/create_order"
	"orders-service/internal/usecase/order_update_status"

	"github.com/google/uuid"
)

func OrderCreateRequestToApplication(req dto.CreateOrderRequest) *create_order.CreateOrderRequest {
	return &create_order.CreateOrderRequest{
		Amount: req.Amount,
	}
}

func UpdateOrderStatusRequestToApplication(
	orderID uuid.UUID,
	dto dto.UpdateOrderStatusRequest,
) (*order_update_status.UpdateStatusCommand, error) {

	status := value_objects.OrderStatus(dto.Status)

	if !status.IsValid() {
		return nil, myerrors.ErrInvalidOrderStatus
	}

	return &order_update_status.UpdateStatusCommand{
		OrderID: orderID,
		Status:  status,
		Actor: order_update_status.UpdateStatusActor{
			Type: order_update_status.UpdateStatusActorStaff,
		},
	}, nil
}

func ConsumerCancelOrderRequestToApplication(
	orderID uuid.UUID, userID uuid.UUID,
	dto *dto.ConsumerCancelOrderRequest) (*cancel_order.CancelOrderCommand, error) {

	var reason value_objects.CancellationReason

	if dto != nil {
		reason = value_objects.CancellationReason(dto.CancellationReason)

		if !reason.IsValid() {
			return nil, myerrors.ErrInvalidCancellationReason
		}

	} else {
		reason = value_objects.CancellationChangedMind
	}

	return &cancel_order.CancelOrderCommand{
		OrderID: orderID,
		Reason:  reason,
		Actor: cancel_order.CancelActor{
			Type: cancel_order.CancelActorCustomer,
			ID:   &userID,
		},
	}, nil
}

func StaffCancelOrderRequestToApplication(
	orderID uuid.UUID,
	dto dto.StaffCancelOrderRequest) (*cancel_order.CancelOrderCommand, error) {

	reason := value_objects.CancellationReason(dto.CancellationReason)

	if !reason.IsValid() {
		return nil, myerrors.ErrInvalidCancellationReason
	}

	return &cancel_order.CancelOrderCommand{
		OrderID: orderID,
		Reason:  reason,
		Actor: cancel_order.CancelActor{
			Type: cancel_order.CancelActorStore,
		},
	}, nil
}
