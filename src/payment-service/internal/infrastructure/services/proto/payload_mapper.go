package proto

import (
	"log"
	"payment-service/internal/usecase/order_created"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

type PayloadMapper struct {
}

func NewPayloadMapper() *PayloadMapper {
	return &PayloadMapper{}
}

func (mapper *PayloadMapper) ToOrderCreatedEvent(payload []byte) (*order_created.OrderCreatedEvent, error) {
	var protoEvent OrderCreatedEvent

	if err := proto.Unmarshal(payload, &protoEvent); err != nil {
		return nil, err
	}

	return orderCreatedEventFromProto(&protoEvent)
}

func orderCreatedEventFromProto(prEvt *OrderCreatedEvent) (*order_created.OrderCreatedEvent, error) {
	orderId, err := uuid.Parse(prEvt.OrderId)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	msgId, err := uuid.Parse(prEvt.MessageId)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	userId, err := uuid.Parse(prEvt.UserId)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &order_created.OrderCreatedEvent{
		OrderId:   orderId,
		MessageId: msgId,
		UserId:    userId,
		Amount:    prEvt.Amount,
		CreatedAt: prEvt.CreatedAt.AsTime(),
	}, nil
}
