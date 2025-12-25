package outbox

import "payment-service/internal/domain/events"

type PaymentProcessedEventPayloadMapper interface {
	PaymentProcessedEventToPayload(event *events.PaymentProcessedEvent) ([]byte, error)
}
