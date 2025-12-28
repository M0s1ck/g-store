package websocket

import (
	"log"

	"github.com/google/uuid"
)

type Hub struct {
	Clients    map[string]*Client
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan []byte
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[string]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan []byte),
	}
}

// TODO: understand this

func (h *Hub) Run() {
	for {
		select {

		case client := <-h.Register:
			h.Clients[client.ID] = client

		case client := <-h.Unregister:
			delete(h.Clients, client.ID)
			close(client.Send)

		case message := <-h.Broadcast:
			for _, client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					delete(h.Clients, client.ID)
				}
			}
		}
	}
}

func (h *Hub) NotifyOrder(userID uuid.UUID, payload []byte) {
	client, ok := h.Clients[userID.String()]
	if !ok {
		log.Printf("client not found in hub: %s", userID)
		return
	}

	client.Send <- payload
}
