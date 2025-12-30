package refund_order

import "github.com/google/uuid"

type Command struct {
	OrderId uuid.UUID
}
