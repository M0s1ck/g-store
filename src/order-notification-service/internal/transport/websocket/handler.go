package websocket

import (
	"net/http"

	"github.com/google/uuid"
)

func NewHandler(hub *Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}

		client := &Client{
			ID:   uuid.NewString(),
			Conn: conn,
			Send: make(chan []byte, 256),
		}

		go client.writePump()
		go client.readPump(hub)
	}
}
