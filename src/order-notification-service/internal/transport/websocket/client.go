package websocket

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID   string
	Conn *websocket.Conn
	Send chan []byte
}

func (c *Client) readPump(h *Hub) {
	defer func() {
		h.Unregister <- c
		err := c.Conn.Close()
		if err != nil {
			log.Printf("Error closing client id=%s: %v", c.ID, err)
		}
	}()

	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}
		h.Broadcast <- msg
	}
}

func (c *Client) writePump() {
	for msg := range c.Send {
		err := c.Conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Printf("write pump error with client=%s: %s", c.ID, err)
		}
	}
}
