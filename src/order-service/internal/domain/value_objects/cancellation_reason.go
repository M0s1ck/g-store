package value_objects

type CancellationReason string

const (
	CancellationNoPaymentAccount     CancellationReason = "NO_PAYMENT_ACCOUNT"
	CancellationInsufficientFunds    CancellationReason = "INSUFFICIENT_FUNDS"
	CancellationPaymentInternalError CancellationReason = "PAYMENT_INTERNAL_ERROR"
	CancellationOutOfStock           CancellationReason = "OUT_OF_STOCK"
	CancellationDeliveryUnavailable  CancellationReason = "DELIVERY_UNAVAILABLE"
	CancellationChangedMind          CancellationReason = "CHANGED_MIND"
)

// TODO: add as enums to db when settled

func (r CancellationReason) IsValid() bool {
	switch r {
	case
		CancellationNoPaymentAccount,
		CancellationInsufficientFunds,
		CancellationPaymentInternalError,
		CancellationOutOfStock,
		CancellationDeliveryUnavailable,
		CancellationChangedMind:
		return true
	default:
		return false
	}
}

func (r CancellationReason) String() string {
	return string(r)
}
