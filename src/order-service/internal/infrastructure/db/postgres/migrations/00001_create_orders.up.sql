CREATE TYPE order_status AS ENUM (
    'PENDING',
    'PAID',
    'CANCELED'
);

CREATE TABLE orders (
    id          UUID PRIMARY KEY,
    user_id     UUID NOT NULL,
    amount      BIGINT NOT NULL CHECK (amount > 0),
    status      order_status NOT NULL,
    description TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_orders_user_id ON orders(user_id);

-- Trigger for auto updated_at update
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_orders_updated_at
    BEFORE UPDATE ON orders
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();
