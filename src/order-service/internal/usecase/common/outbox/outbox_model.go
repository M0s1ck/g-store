package outbox

import (
	"time"

	"github.com/google/uuid"
)

type Model struct {
	Id          uuid.UUID  `db:"id"`
	Aggregate   string     `db:"aggregate"` // type of aggregate, e. g. "order"
	AggregateID uuid.UUID  `db:"aggregate_id"`
	EventType   string     `db:"event_type"`
	Payload     []byte     `db:"payload"` // serialized event
	CreatedAt   time.Time  `db:"created_at"`
	SentAt      *time.Time `db:"sent_at"` // nil if not sent yet
}
