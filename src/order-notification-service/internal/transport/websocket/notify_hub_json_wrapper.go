package websocket

import (
	"context"
	"encoding/json"
	"log"
	"order-notification-service/internal/domain/events/consumed"

	"github.com/google/uuid"
)

type NotifyHubJSONWrapper struct {
	hub *Hub
}

func NewNotifyHubJSONWrapper(hub *Hub) *NotifyHubJSONWrapper {
	return &NotifyHubJSONWrapper{
		hub: hub,
	}
}

func (n *NotifyHubJSONWrapper) NotifyOrderStatusChanged(
	_ context.Context,
	event consumed_events.OrderStatusChangedEvent,
) {
	data, err := json.Marshal(struct {
		OrderID uuid.UUID `json:"orderId"`
		Status  string    `json:"status"`
	}{
		OrderID: event.OrderId,
		Status:  string(event.Status),
	})

	if err != nil {
		log.Printf("marshal error: %v", err)
		return
	}

	envelope := Envelope{
		Type: "order.status.changed",
		Data: data,
	}

	payload, err := json.Marshal(envelope)
	if err != nil {
		log.Printf("marshal envelope error: %v", err)
		return
	}

	n.hub.NotifyOrder(event.OrderId, payload)
}
