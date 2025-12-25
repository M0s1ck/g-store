CREATE TABLE accounts (
    id         UUID PRIMARY KEY,
    user_id    UUID UNIQUE NOT NULL,
    balance    BIGINT      NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TYPE transaction_type AS ENUM(
    'TOP_UP',
    'PAYMENT'
);

CREATE TABLE balance_transactions (
    id UUID PRIMARY KEY,
    account_id UUID NOT NULL,
    order_id UUID,
    amount BIGINT NOT NULL CHECK (amount <> 0),
    type transaction_type NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    FOREIGN KEY (account_id) REFERENCES accounts(id)
        ON DELETE CASCADE
);

CREATE INDEX ix_balance_transactions_account
    ON balance_transactions(account_id);

CREATE INDEX ix_balance_transactions_account_created
    ON balance_transactions(account_id, created_at DESC);

-- Trigger for auto updated_at update
CREATE OR REPLACE FUNCTION set_updated_at()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_accounts_updated_at
    BEFORE UPDATE ON accounts
    FOR EACH ROW
EXECUTE FUNCTION set_updated_at();