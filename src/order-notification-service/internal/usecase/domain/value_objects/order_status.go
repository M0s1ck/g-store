package value_objects

type OrderStatus string

const (
	OrderPending    OrderStatus = "PENDING"
	OrderPaid       OrderStatus = "PAID"
	OrderCanceled   OrderStatus = "CANCELED"
	OrderAssembling OrderStatus = "ASSEMBLING"
	OrderAssembled  OrderStatus = "ASSEMBLED"
	OrderDelivering OrderStatus = "DELIVERING"
	OrderIssued     OrderStatus = "ISSUED"
)
