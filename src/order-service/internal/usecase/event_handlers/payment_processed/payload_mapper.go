package payment_processed

import "orders-service/internal/domain/events/consumed"

type PayloadMapper interface {
	ToPaymentProcessedEvent(payload []byte) (*consumed_events.PaymentProcessedEvent, error)
}
