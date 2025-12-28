package proto

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"orders-service/internal/domain/entities"
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
	case entities.OrderPending:
		status = OrderStatus_ORDER_STATUS_PENDING
	case entities.OrderPaid:
		status = OrderStatus_ORDER_STATUS_PAID
	case entities.OrderCanceled:
		status = OrderStatus_ORDER_STATUS_CANCELED
	case entities.OrderAssembling:
		status = OrderStatus_ORDER_STATUS_ASSEMBLING
	case entities.OrderAssembled:
		status = OrderStatus_ORDER_STATUS_ASSEMBLED
	case entities.OrderDelivering:
		status = OrderStatus_ORDER_STATUS_DELIVERING
	case entities.OrderIssued:
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

// TODO: move from here

func OrderStatusChangedEventFromProto(
	prEvt *OrderStatusChangedEvent,
) (*events.OrderStatusChangedEvent, error) {

	// --- UUIDs ---
	messageId, err := uuid.Parse(prEvt.MessageId)
	if err != nil {
		return nil, fmt.Errorf("invalid message_id: %w", err)
	}

	orderId, err := uuid.Parse(prEvt.OrderId)
	if err != nil {
		return nil, fmt.Errorf("invalid order_id: %w", err)
	}

	userId, err := uuid.Parse(prEvt.UserId)
	if err != nil {
		return nil, fmt.Errorf("invalid user_id: %w", err)
	}

	// --- status ---
	var status entities.OrderStatus
	switch prEvt.Status {
	case OrderStatus_ORDER_STATUS_PENDING:
		status = entities.OrderPending
	case OrderStatus_ORDER_STATUS_PAID:
		status = entities.OrderPaid
	case OrderStatus_ORDER_STATUS_CANCELED:
		status = entities.OrderCanceled
	case OrderStatus_ORDER_STATUS_ASSEMBLING:
		status = entities.OrderAssembling
	case OrderStatus_ORDER_STATUS_ASSEMBLED:
		status = entities.OrderAssembled
	case OrderStatus_ORDER_STATUS_DELIVERING:
		status = entities.OrderDelivering
	case OrderStatus_ORDER_STATUS_ISSUED:
		status = entities.OrderIssued
	default:
		return nil, fmt.Errorf("unknown order status: %v", prEvt.Status)
	}

	// --- cancellation reason ---
	var cancellationReason *string
	if prEvt.CancellationReason != nil && *prEvt.CancellationReason != "" {
		cancellationReason = prEvt.CancellationReason
	}

	// --- created at ---
	var createdAt time.Time
	if prEvt.CreatedAt != nil {
		createdAt = prEvt.CreatedAt.AsTime()
	}

	return &events.OrderStatusChangedEvent{
		MessageId:          messageId,
		OrderId:            orderId,
		UserId:             userId,
		Status:             status,
		CancellationReason: cancellationReason,
		CreatedAt:          createdAt,
	}, nil
}
