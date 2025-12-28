package order_created

import (
	"log"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"orders-service/internal/domain/events"
)

type PayloadMapper struct {
}

func NewPayloadMapper() *PayloadMapper {
	return &PayloadMapper{}
}

func (p PayloadMapper) OrderCreatedEventToPayload(event *events.OrderCreatedEvent) ([]byte, error) {
	protoEvt := orderCreatedEventToProto(event)

	payload, err := proto.Marshal(protoEvt)
	if err != nil {
		log.Println("Error while marshalling proto order-created-event to payload")
		return nil, err
	}

	return payload, nil
}

func orderCreatedEventToProto(event *events.OrderCreatedEvent) *OrderCreatedEvent {
	return &OrderCreatedEvent{
		MessageId: event.MessageId.String(),
		OrderId:   event.OrderId.String(),
		UserId:    event.UserId.String(),
		Amount:    event.Amount,
		CreatedAt: timestamppb.New(event.CreatedAt),
	}
}
