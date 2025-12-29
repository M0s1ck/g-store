package order_update_status

type UpdateStatusActorType string

const (
	UpdateStatusActorStaff          UpdateStatusActorType = "staff"
	UpdateStatusActorPaymentService UpdateStatusActorType = "payment_service"
)

type UpdateStatusActor struct {
	Type UpdateStatusActorType
}
