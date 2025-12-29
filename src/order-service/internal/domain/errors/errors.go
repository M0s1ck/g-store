package errors

import (
	"errors"
	"fmt"
)

var ErrOrderNotFound = fmt.Errorf("order not found")

var ErrForbidden = fmt.Errorf("forbidden")

var ErrAmountNotPositive = fmt.Errorf("order amount must be positive")

var ErrInvalidOrderStatus = errors.New("invalid order status")

var ErrInvalidOrderStatusChange = errors.New("invalid order status transition")

var ErrCancellationReasonRequired = errors.New("cancellation reason required")

var ErrOrderAlreadyCanceled = errors.New("order already canceled")

var ErrInvalidCancellationReason = errors.New("invalid cancellation reason")

var ErrOrderCantBeCancelled = errors.New("can't be cancelled")
