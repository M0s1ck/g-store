package payment_processed

type PayloadMapper interface {
	ToPaymentProcessedEvent(payload []byte) (*PaymentProcessedEvent, error)
}
