package proto_mappers

import (
	"log"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"orders-service/internal/domain/events/produced"
	"orders-service/internal/infrastructure/msg_broker/proto/gen"
)

type OrderCreatedPayloadMapper struct {
}

func NewOrderCreatedPayloadMapper() *OrderCreatedPayloadMapper {
	return &OrderCreatedPayloadMapper{}
}

func (p OrderCreatedPayloadMapper) OrderCreatedEventToPayload(event *published_events.OrderCreatedEvent) ([]byte, error) {
	protoEvt := orderCreatedEventToProto(event)

	payload, err := proto.Marshal(protoEvt)
	if err != nil {
		log.Println("Error while marshalling proto order-created-event to payload")
		return nil, err
	}

	return payload, nil
}

func orderCreatedEventToProto(event *published_events.OrderCreatedEvent) *gen.OrderCreatedEvent {
	return &gen.OrderCreatedEvent{
		OrderId:   event.OrderId.String(),
		UserId:    event.UserId.String(),
		Amount:    event.Amount,
		CreatedAt: timestamppb.New(event.OccurredAt),
	}
}
