package outbox

import (
	"payment-service/internal/domain/events/produced"
)

type PaymentProcessedEventPayloadMapper interface {
	PaymentProcessedEventToPayload(event *produced_events.PaymentProcessedEvent) ([]byte, error)
}
