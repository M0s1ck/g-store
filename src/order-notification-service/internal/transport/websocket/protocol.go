package websocket

import "encoding/json"

type Envelope struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}
