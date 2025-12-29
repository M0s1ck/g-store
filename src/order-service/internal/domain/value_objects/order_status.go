package value_objects

type OrderStatus string

const (
	OrderPending    OrderStatus = "PENDING"
	OrderPaid       OrderStatus = "PAID"
	OrderCanceled   OrderStatus = "CANCELLED"
	OrderAssembling OrderStatus = "ASSEMBLING"
	OrderAssembled  OrderStatus = "ASSEMBLED"
	OrderDelivering OrderStatus = "DELIVERING"
	OrderIssued     OrderStatus = "ISSUED"
)

func (s OrderStatus) IsValid() bool {
	switch s {
	case OrderPending,
		OrderPaid,
		OrderCanceled,
		OrderAssembling,
		OrderAssembled,
		OrderDelivering,
		OrderIssued:
		return true
	default:
		return false
	}
}

var allowedTransitions = map[OrderStatus][]OrderStatus{
	OrderPending: {
		OrderPaid,
		OrderCanceled,
	},
	OrderPaid: {
		OrderAssembling,
		OrderCanceled,
	},
	OrderAssembling: {
		OrderAssembled,
		OrderCanceled,
	},
	OrderAssembled: {
		OrderDelivering,
		OrderCanceled,
	},
	OrderDelivering: {
		OrderIssued,
		OrderCanceled,
	},
}

func (s OrderStatus) CanTransitionTo(to OrderStatus) bool {
	allowed, ok := allowedTransitions[s]
	if !ok {
		return false
	}

	for _, s := range allowed {
		if s == to {
			return true
		}
	}

	return false
}
