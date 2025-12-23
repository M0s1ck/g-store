package proto

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"orders-service/internal/domain/events"
)

func OrderCreatedEventToProto(event *events.OrderCreatedEvent) *OrderCreatedEvent {
	return &OrderCreatedEvent{
		MessageId: event.MessageId.String(),
		OrderId:   event.OrderId.String(),
		UserId:    event.UserId.String(),
		Amount:    event.Amount,
		CreatedAt: timestamppb.New(event.CreatedAt),
	}
}
