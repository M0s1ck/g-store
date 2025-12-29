package proto_mappers

import (
	"log"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"orders-service/internal/domain/events/produced"
	"orders-service/internal/domain/value_objects"
	"orders-service/internal/infrastructure/msg_broker/proto/gen"
)

type OrderStatusChangedPayloadMapper struct {
}

func NewOrderStatusChangedPayloadMapper() *OrderStatusChangedPayloadMapper {
	return &OrderStatusChangedPayloadMapper{}
}

func (p *OrderStatusChangedPayloadMapper) OrderStatusChangedEventToPayload(event *published_events.OrderStatusChangedEvent) ([]byte, error) {
	evtProto := orderStatusChangedEventToProto(event)

	payload, err := proto.Marshal(evtProto)
	if err != nil {
		log.Println("Error while marshalling proto order-status-changed-event to payload")
		return nil, err
	}

	return payload, nil
}

func orderStatusChangedEventToProto(
	event *published_events.OrderStatusChangedEvent,
) *gen.OrderStatusChangedEvent {

	// --- status ---
	var status gen.OrderStatus
	switch event.Status {
	case value_objects.OrderPending:
		status = gen.OrderStatus_ORDER_STATUS_PENDING
	case value_objects.OrderPaid:
		status = gen.OrderStatus_ORDER_STATUS_PAID
	case value_objects.OrderCanceled:
		status = gen.OrderStatus_ORDER_STATUS_CANCELLED
	case value_objects.OrderAssembling:
		status = gen.OrderStatus_ORDER_STATUS_ASSEMBLING
	case value_objects.OrderAssembled:
		status = gen.OrderStatus_ORDER_STATUS_ASSEMBLED
	case value_objects.OrderDelivering:
		status = gen.OrderStatus_ORDER_STATUS_DELIVERING
	case value_objects.OrderIssued:
		status = gen.OrderStatus_ORDER_STATUS_ISSUED
	default:
		status = gen.OrderStatus_ORDER_STATUS_UNSPECIFIED
	}

	// --- created at ---
	var createdAt *timestamppb.Timestamp
	if !event.OccurredAt.IsZero() {
		createdAt = timestamppb.New(event.OccurredAt)
	}

	return &gen.OrderStatusChangedEvent{
		OrderId:   event.OrderId.String(),
		UserId:    event.UserId.String(),
		Status:    status,
		CreatedAt: createdAt,
	}
}
