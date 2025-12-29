package cancel_order

import (
	"orders-service/internal/domain/entities"
	myerrors "orders-service/internal/domain/errors"
	"orders-service/internal/domain/value_objects"
)

// CancelPolicy domain rules to determine if actor can cancel having specific reason
type CancelPolicy struct {
}

func NewCancelPolicy() *CancelPolicy {
	return &CancelPolicy{}
}

// CanCancel domain rules to determine if actor can cancel the order having this specific reason
func (p *CancelPolicy) CanCancel(
	order *entities.Order,
	actor CancelActor,
	reason value_objects.CancellationReason,
) error {

	if !reason.IsValid() {
		return myerrors.ErrInvalidCancellationReason
	}

	if order.Status == value_objects.OrderCanceled {
		return myerrors.ErrOrderAlreadyCanceled
	}

	if order.Status == value_objects.OrderIssued {
		return myerrors.ErrOrderCantBeCancelled
	}

	switch actor.Type {

	case CancelActorCustomer:
		if actor.ID == nil || order.UserId != *actor.ID {
			return myerrors.ErrForbidden
		}

		if !isCustomerAllowedReason(reason) {
			return myerrors.ErrInvalidCancellationReason
		}

	case CancelActorStore:
		if !isStoreAllowedReason(reason) {
			return myerrors.ErrInvalidCancellationReason
		}

	case CancelActorPaymentService:
		if !isPaymentServiceAllowedReason(reason) {
			return myerrors.ErrInvalidCancellationReason
		}

	default:
		return myerrors.ErrForbidden
	}

	return nil
}

func isCustomerAllowedReason(r value_objects.CancellationReason) bool {
	switch r {
	case
		value_objects.CancellationChangedMind:
		return true
	default:
		return false
	}
}

func isStoreAllowedReason(r value_objects.CancellationReason) bool {
	switch r {
	case
		value_objects.CancellationOutOfStock,
		value_objects.CancellationDeliveryUnavailable:
		return true
	default:
		return false
	}
}

func isPaymentServiceAllowedReason(r value_objects.CancellationReason) bool {
	switch r {
	case
		value_objects.CancellationNoPaymentAccount,
		value_objects.CancellationInsufficientFunds,
		value_objects.CancellationPaymentInternalError:
		return true
	default:
		return false
	}
}
