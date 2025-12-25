CREATE TABLE inbox (
    id UUID PRIMARY KEY,
    key BYTEA,
    topic TEXT NOT NULL,
    event_type TEXT NOT NULL,
    payload BYTEA NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    processed_at TIMESTAMPTZ
);

-- fast search for unprocessed
CREATE INDEX idx_inbox_unprocessed
    ON inbox (created_at)
    WHERE processed_at IS NULL;