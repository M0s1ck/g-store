package outbox

import (
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"

	"orders-service/internal/domain/events"
	myproto "orders-service/internal/infrastructure/services/proto"
	"orders-service/internal/usecase/common/outbox"
)

type ModelProtoFactory struct {
}

func NewOutboxModelProtoFactory() *ModelProtoFactory {
	return &ModelProtoFactory{}
}

func (s *ModelProtoFactory) CreateOutboxModelFromOrderCreatedEvent(
	event *events.OrderCreatedEvent,
) (*outbox.Model, error) {

	eventProto := myproto.OrderCreatedEventToProto(event)

	payload, err := proto.Marshal(eventProto)
	if err != nil {
		return nil, err
	}

	return &outbox.Model{
		Id:          uuid.New(),
		Aggregate:   "order",
		AggregateID: event.OrderId,
		EventType:   "OrderCreated",
		Payload:     payload,
		CreatedAt:   time.Now(),
	}, nil
}
