package proto

import (
	"log"
	"orders-service/internal/domain/value_objects"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"orders-service/internal/domain/events"
)

type PayloadMapper struct {
}

func NewPayloadMapper() *PayloadMapper {
	return &PayloadMapper{}
}

func (p *PayloadMapper) OrderStatusChangedEventToPayload(event *events.OrderStatusChangedEvent) ([]byte, error) {
	evtProto := orderStatusChangedEventToProto(event)

	payload, err := proto.Marshal(evtProto)
	if err != nil {
		log.Println("Error while marshalling proto order-status-changed-event to payload")
		return nil, err
	}

	return payload, nil
}

func orderStatusChangedEventToProto(
	event *events.OrderStatusChangedEvent,
) *OrderStatusChangedEvent {

	// --- status ---
	var status OrderStatus
	switch event.Status {
	case value_objects.OrderPending:
		status = OrderStatus_ORDER_STATUS_PENDING
	case value_objects.OrderPaid:
		status = OrderStatus_ORDER_STATUS_PAID
	case value_objects.OrderCanceled:
		status = OrderStatus_ORDER_STATUS_CANCELED
	case value_objects.OrderAssembling:
		status = OrderStatus_ORDER_STATUS_ASSEMBLING
	case value_objects.OrderAssembled:
		status = OrderStatus_ORDER_STATUS_ASSEMBLED
	case value_objects.OrderDelivering:
		status = OrderStatus_ORDER_STATUS_DELIVERING
	case value_objects.OrderIssued:
		status = OrderStatus_ORDER_STATUS_ISSUED
	default:
		status = OrderStatus_ORDER_STATUS_UNSPECIFIED
	}

	// --- cancellation reason (optional string) ---
	var cancellationReason *string
	if event.CancellationReason != nil && *event.CancellationReason != "" {
		cancellationReason = event.CancellationReason
	}

	// --- created at ---
	var createdAt *timestamppb.Timestamp
	if !event.CreatedAt.IsZero() {
		createdAt = timestamppb.New(event.CreatedAt)
	}

	return &OrderStatusChangedEvent{
		MessageId:          event.MessageId.String(),
		OrderId:            event.OrderId.String(),
		UserId:             event.UserId.String(),
		Status:             status,
		CancellationReason: cancellationReason,
		CreatedAt:          createdAt,
	}
}
