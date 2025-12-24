package outbox

import (
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"

	"orders-service/internal/domain/events"
	"orders-service/internal/domain/messages"
	myproto "orders-service/internal/infrastructure/services/proto"
)

type ModelProtoFactory struct {
	orderCreatedEventTopic string
}

func NewOutboxModelProtoFactory(eventTopic string) *ModelProtoFactory {
	return &ModelProtoFactory{
		orderCreatedEventTopic: eventTopic,
	}
}

func (f *ModelProtoFactory) CreateOutboxModelFromOrderCreatedEvent(
	event *events.OrderCreatedEvent,
) (*messages.OutboxMessage, error) {

	eventProto := myproto.OrderCreatedEventToProto(event)

	payload, err := proto.Marshal(eventProto)
	if err != nil {
		return nil, err
	}

	return &messages.OutboxMessage{
		Id:          uuid.New(),
		Aggregate:   "order",
		AggregateID: event.OrderId,
		EventType:   f.orderCreatedEventTopic,
		Payload:     payload,
		CreatedAt:   time.Now(),
	}, nil
}
