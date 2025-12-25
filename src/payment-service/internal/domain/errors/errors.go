package errors

import "errors"

var ErrInsufficientFunds = errors.New("insufficient funds")

var ErrAccountNotFound = errors.New("account not found")
