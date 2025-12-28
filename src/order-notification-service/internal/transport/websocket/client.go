package websocket

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID   string
	Conn *websocket.Conn
	Send chan []byte
}

// handles received msgs from client (already in websocket), should work in a separate goroutine
// subscribe -> init client notification system
func (c *Client) readPump(h *Hub) {
	// unsubscribe by default: in case of stopped connection from client or an error
	defer func() {
		h.unsubscribe <- c
		_ = c.Conn.Close()
	}()

	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			// error or client stopped connection
			log.Printf("read pump for client %s stopped : %s", c.ID, err)
			break
		}

		var m ClientConnectionMessage
		if err := json.Unmarshal(msg, &m); err != nil {
			log.Printf("invalid connect message from client %s: %v", c.ID, err)
			continue
		}

		switch m.Type {

		case "subscribe":
			h.subscribe <- subscription{
				orderID: m.OrderID,
				client:  c,
			}

		default:
			log.Printf("unknown message type from client %s: %s", c.ID, m.Type)
		}
	}
}

func (c *Client) writePump() {
	defer func() {
		_ = c.Conn.Close()
	}()

	for msg := range c.Send {
		err := c.Conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Printf("write error for client %s: %v", c.ID, err)
			return
		}
	}
}
