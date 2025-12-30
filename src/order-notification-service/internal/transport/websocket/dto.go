package websocket

import "github.com/google/uuid"

// ClientConnectionMessage client sends this to modify connection: to subscribe / unsubscribe
type ClientConnectionMessage struct {
	Type    string    `json:"type"`
	OrderID uuid.UUID `json:"orderId"`
}

type StatusChanged struct {
	OrderID uuid.UUID `json:"orderId"`
	Status  string    `json:"status"`
}

type OrderCancelled struct {
	OrderID uuid.UUID `json:"orderId"`
	Reason  string    `json:"reason"`
}
