package common_outbox

import (
	"time"

	"github.com/google/uuid"

	"orders-service/internal/domain/events/produced"
	"orders-service/internal/domain/messages"
)

type MessageFactory struct {
	ordCrMapper         OrderCreatedEventPayloadMapper
	ordCancelMapper     OrderCancelledPayloadMapper
	ordStMapper         OrderStatusChangedPayloadMapper
	ordCrEvtType        string
	ordCancelEvtType    string
	ordStChangedEvtType string
}

func NewMessageFactory(
	ordCrMapper OrderCreatedEventPayloadMapper,
	ordCancelMapper OrderCancelledPayloadMapper,
	ordStMapper OrderStatusChangedPayloadMapper,
	ordCrEvtType string,
	ordCancelEvtType string,
	ordStChangedEvtType string,
) *MessageFactory {

	return &MessageFactory{
		ordCrMapper:         ordCrMapper,
		ordStMapper:         ordStMapper,
		ordCancelMapper:     ordCancelMapper,
		ordCancelEvtType:    ordCancelEvtType,
		ordStChangedEvtType: ordStChangedEvtType,
		ordCrEvtType:        ordCrEvtType,
	}
}

func (f *MessageFactory) CreateFromOrderCreatedEvent(
	event *published_events.OrderCreatedEvent,
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

func (f *MessageFactory) CreateMessageOrderCancelledEvent(
	event *published_events.OrderCancelledEvent,
) (*messages.OutboxMessage, error) {

	payload, err := f.ordCancelMapper.OrderCancelledEventToPayload(event)
	if err != nil {
		return nil, err
	}

	return &messages.OutboxMessage{
		Id:          uuid.New(),
		Aggregate:   "order",
		AggregateID: event.OrderId,
		EventType:   f.ordCancelEvtType,
		Key:         event.OrderId.String(),
		Payload:     payload,
		CreatedAt:   time.Now(),
	}, nil
}

func (f *MessageFactory) CreateMessageOrderStatusChangedEvent(
	event *published_events.OrderStatusChangedEvent,
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
		Key:         event.UserId.String(), // there will be partition by user
		Payload:     payload,
		CreatedAt:   time.Now(),
	}, nil
}
