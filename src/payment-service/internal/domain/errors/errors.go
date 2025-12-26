package errors

import (
	"errors"
	"fmt"
)

var ErrInsufficientFunds = errors.New("insufficient funds")

var ErrAccountNotFound = errors.New("account not found")

var ErrForbidden = fmt.Errorf("forbidden")

var ErrAmountNotPositive = fmt.Errorf("order amount must be positive")

var ErrAccountForUserAlreadyExists = errors.New("account for user already exists")
