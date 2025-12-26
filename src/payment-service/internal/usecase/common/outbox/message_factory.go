package outbox

import (
	"time"

	"payment-service/internal/domain/events"
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

func (f *MessageFactory) PaymentProcessedEventToOutboxMessage(event *events.PaymentProcessedEvent,
) *messages.OutboxMessage {

	payload, err := f.paymentProcessedEventPayloadMapper.PaymentProcessedEventToPayload(event)
	if err != nil {
		return nil
	}

	return &messages.OutboxMessage{
		Id:          event.MessageId,
		Aggregate:   "payment",
		AggregateID: event.OrderId,
		EventType:   f.paymentProcessedEventType,
		Payload:     payload,
		CreatedAt:   time.Now(),
		SentAt:      nil,
		RetryCount:  0,
	}
}
