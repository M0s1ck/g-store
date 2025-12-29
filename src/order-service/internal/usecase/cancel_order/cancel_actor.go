package cancel_order

import "github.com/google/uuid"

type CancelActorType string

const (
	CancelActorCustomer       CancelActorType = "customer"
	CancelActorStore          CancelActorType = "store"
	CancelActorPaymentService CancelActorType = "payment_service"
)

type CancelActor struct {
	Type CancelActorType
	ID   *uuid.UUID // nil for store/system
}
