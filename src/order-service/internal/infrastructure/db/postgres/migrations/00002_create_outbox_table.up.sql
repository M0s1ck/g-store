CREATE TABLE outbox (
    id            UUID PRIMARY KEY,
    aggregate     TEXT NOT NULL,
    aggregate_id  UUID NOT NULL,
    event_type    TEXT NOT NULL,
    payload       BYTEA NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    sent_at       TIMESTAMPTZ
);

-- fast search for unsent
CREATE INDEX idx_outbox_unsent
    ON outbox (created_at)
    WHERE sent_at IS NULL;
