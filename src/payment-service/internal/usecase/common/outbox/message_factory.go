package outbox

import (
	"time"

	"github.com/google/uuid"

	"payment-service/internal/domain/events/produced"
	"payment-service/internal/domain/messages"
)

type MessageFactory struct {
	paymentProcessedEventPayloadMapper PaymentProcessedEventPayloadMapper
	paymentProcessedEventType          string
}

func NewOutboxMessageFactory(
	paymentProcessedEventPayloadMapper PaymentProcessedEventPayloadMapper,
	paymentProcessedEventType string,
) *MessageFactory {

	return &MessageFactory{
		paymentProcessedEventPayloadMapper: paymentProcessedEventPayloadMapper,
		paymentProcessedEventType:          paymentProcessedEventType,
	}
}

func (f *MessageFactory) PaymentProcessedEventToOutboxMessage(event *produced_events.PaymentProcessedEvent,
) *messages.OutboxMessage {

	payload, err := f.paymentProcessedEventPayloadMapper.PaymentProcessedEventToPayload(event)
	if err != nil {
		return nil
	}

	return &messages.OutboxMessage{
		Id:          uuid.New(),
		Aggregate:   "payment",
		AggregateID: event.OrderId,
		EventType:   f.paymentProcessedEventType,
		Payload:     payload,
		CreatedAt:   time.Now(),
		SentAt:      nil,
		RetryCount:  0,
	}
}
