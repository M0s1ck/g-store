package errors

import "fmt"

var ErrOrderNotFound = fmt.Errorf("order not found")

var ErrForbidden = fmt.Errorf("forbidden")
