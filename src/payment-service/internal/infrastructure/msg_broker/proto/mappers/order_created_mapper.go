package proto_mappers

import (
	"log"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"

	"payment-service/internal/domain/events/consumed"
	"payment-service/internal/infrastructure/msg_broker/proto/gen"
)

type OrderCreatedPayloadMapper struct {
}

func NewOrderCreatedPayloadMapper() *OrderCreatedPayloadMapper {
	return &OrderCreatedPayloadMapper{}
}

func (mapper *OrderCreatedPayloadMapper) ToOrderCreatedEvent(payload []byte) (*consumed_events.OrderCreatedEvent, error) {
	var protoEvent gen.OrderCreatedEvent

	if err := proto.Unmarshal(payload, &protoEvent); err != nil {
		return nil, err
	}

	return orderCreatedEventFromProto(&protoEvent)
}

func orderCreatedEventFromProto(prEvt *gen.OrderCreatedEvent) (*consumed_events.OrderCreatedEvent, error) {
	orderId, err := uuid.Parse(prEvt.OrderId)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	userId, err := uuid.Parse(prEvt.UserId)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &consumed_events.OrderCreatedEvent{
		OrderId:   orderId,
		UserId:    userId,
		Amount:    prEvt.Amount,
		CreatedAt: prEvt.CreatedAt.AsTime(),
	}, nil
}
