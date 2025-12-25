package messages

import (
	"time"

	"github.com/google/uuid"
)

type InboxMessage struct {
	Id          uuid.UUID  `db:"id"`
	Topic       string     `db:"topic"`
	EventType   string     `db:"event_type"`
	Key         []byte     `db:"key"`
	Payload     []byte     `db:"payload"`
	CreatedAt   time.Time  `db:"created_at"`
	ProcessedAt *time.Time `db:"processed_at"` // nil if not processed yet
}
