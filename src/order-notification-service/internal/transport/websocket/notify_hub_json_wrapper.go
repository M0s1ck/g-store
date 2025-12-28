package websocket

import (
	"context"
	"encoding/json"
	"log"
	"order-notification-service/internal/usecase/notify_order_status_changed"

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
	ctx context.Context,
	event notify_order_status_changed.Event,
) {
	data, err := json.Marshal(struct {
		OrderID            uuid.UUID `json:"orderId"`
		Status             string    `json:"status"`
		CancellationReason *string   `json:"cancellationReason,omitempty"`
	}{
		OrderID:            event.OrderID,
		Status:             string(event.Status),
		CancellationReason: event.CancellationReason,
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

	n.hub.NotifyOrder(event.OrderID, payload)
}
