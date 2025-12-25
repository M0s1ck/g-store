package errors

import "fmt"

var ErrOrderNotFound = fmt.Errorf("order not found")

var ErrForbidden = fmt.Errorf("forbidden")

var ErrAmountNotPositive = fmt.Errorf("order amount must be positive")
