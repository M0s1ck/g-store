package websocket

import (
	"log"

	"github.com/google/uuid"
)

// Hub manages subscribers-clients notification, observer-subscribers pattern
type Hub struct {
	subscribers map[uuid.UUID]map[*Client]struct{} // orderID -> set of clients
	subscribe   chan subscription
	unsubscribe chan *Client
	notify      chan notification
}

func NewHub() *Hub {
	return &Hub{
		subscribers: make(map[uuid.UUID]map[*Client]struct{}),
		subscribe:   make(chan subscription),
		unsubscribe: make(chan *Client),
		notify:      make(chan notification),
	}
}

func (h *Hub) Run() {
	for {
		select {

		case sub := <-h.subscribe:
			clients := h.subscribers[sub.orderID]
			if clients == nil {
				clients = make(map[*Client]struct{})
				h.subscribers[sub.orderID] = clients
			}
			clients[sub.client] = struct{}{}
			log.Printf("client %v subscribed to order %v", sub.client.ID, sub.orderID)

		case client := <-h.unsubscribe:
			for orderID, clients := range h.subscribers {
				delete(clients, client)
				if len(clients) == 0 {
					delete(h.subscribers, orderID)
				}
			}

		// triggers sending notify.payload to every of subscribers of notify.orderID
		case n := <-h.notify:
			clients := h.subscribers[n.orderID]
			for c := range clients {
				select {
				case c.Send <- n.payload:
				default:
					// in case write pump ain't write pumping fast enough we drop it to prevent block
				}
			}
		}
	}
}

func (h *Hub) NotifyOrder(orderID uuid.UUID, payload []byte) {
	h.notify <- notification{orderID, payload}
}

type subscription struct {
	orderID uuid.UUID
	client  *Client
}

type notification struct {
	orderID uuid.UUID
	payload []byte
}
