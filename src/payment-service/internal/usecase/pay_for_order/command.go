package pay_for_order

import "github.com/google/uuid"

type Command struct {
	OrderId uuid.UUID
	UserId  uuid.UUID
	Amount  int64
}
