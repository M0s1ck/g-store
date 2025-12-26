package outbox

import (
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"

	myproto "orders-service/internal/delivery/proto/order_created"
	"orders-service/internal/domain/events"
	"orders-service/internal/domain/messages"
)

type MessageProtoFactory struct {
	orderCreatedEventType string
}

func NewOutboxModelProtoFactory(orderCreatedEventName string) *MessageProtoFactory {
	return &MessageProtoFactory{
		orderCreatedEventType: orderCreatedEventName,
	}
}

func (f *MessageProtoFactory) CreateOutboxModelFromOrderCreatedEvent(
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
