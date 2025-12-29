package order_update_status

import (
	"orders-service/internal/domain/entities"
	myerrors "orders-service/internal/domain/errors"
	"orders-service/internal/domain/value_objects"
)

type UpdateStatusPolicy struct{}

func NewUpdateStatusPolicy() *UpdateStatusPolicy {
	return &UpdateStatusPolicy{}
}

func (p *UpdateStatusPolicy) CanUpdateStatus(
	_ *entities.Order,
	actor UpdateStatusActor,
	newStatus value_objects.OrderStatus,
) error {

	if !newStatus.IsValid() {
		return myerrors.ErrInvalidOrderStatus
	}

	switch newStatus {

	case value_objects.OrderCanceled:
		return myerrors.ErrInvalidOrderStatusChange

	case value_objects.OrderPending:
		return myerrors.ErrInvalidOrderStatusChange

	case value_objects.OrderPaid:
		if actor.Type != UpdateStatusActorPaymentService {
			return myerrors.ErrForbidden
		}

	default:
		if actor.Type != UpdateStatusActorStaff {
			return myerrors.ErrForbidden
		}
	}

	return nil
}
