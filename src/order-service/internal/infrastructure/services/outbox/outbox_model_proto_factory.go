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
	orderCreatedEventType string
}

func NewOutboxModelProtoFactory(orderCreatedEventName string) *ModelProtoFactory {
	return &ModelProtoFactory{
		orderCreatedEventType: orderCreatedEventName,
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
		EventType:   f.orderCreatedEventType,
		Payload:     payload,
		CreatedAt:   time.Now(),
	}, nil
}
