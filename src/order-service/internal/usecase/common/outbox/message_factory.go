package outbox

import (
	"time"

	"github.com/google/uuid"

	"orders-service/internal/domain/events"
	"orders-service/internal/domain/messages"
)

type MessageFactory struct {
	ordCrMapper         OrderCreatedEventPayloadMapper
	ordStMapper         OrderStatusChangedPayloadMapper
	ordStChangedEvtType string
	ordCrEvtType        string
}

func NewMessageFactory(
	ordCrMapper OrderCreatedEventPayloadMapper,
	ordStMapper OrderStatusChangedPayloadMapper,
	ordStCrEvtType string,
	ordStChangedEvtType string,
) *MessageFactory {

	return &MessageFactory{
		ordCrMapper:         ordCrMapper,
		ordStMapper:         ordStMapper,
		ordStChangedEvtType: ordStChangedEvtType,
		ordCrEvtType:        ordStCrEvtType,
	}
}

func (f *MessageFactory) CreateFromOrderCreatedEvent(
	event *events.OrderCreatedEvent,
) (*messages.OutboxMessage, error) {

	payload, err := f.ordCrMapper.OrderCreatedEventToPayload(event)
	if err != nil {
		return nil, err
	}

	return &messages.OutboxMessage{
		Id:          uuid.New(),
		Aggregate:   "order",
		AggregateID: event.OrderId,
		EventType:   f.ordCrEvtType,
		Key:         event.OrderId.String(),
		Payload:     payload,
		CreatedAt:   time.Now(),
	}, nil
}

func (f *MessageFactory) CreateMessageOrderStatusChangedEvent(
	event *events.OrderStatusChangedEvent,
) (*messages.OutboxMessage, error) {

	payload, err := f.ordStMapper.OrderStatusChangedEventToPayload(event)
	if err != nil {
		return nil, err
	}

	return &messages.OutboxMessage{
		Id:          uuid.New(),
		Aggregate:   "order",
		AggregateID: event.OrderId,
		EventType:   f.ordStChangedEvtType,
		Key:         event.UserId.String(),
		Payload:     payload,
		CreatedAt:   time.Now(),
	}, nil
}
