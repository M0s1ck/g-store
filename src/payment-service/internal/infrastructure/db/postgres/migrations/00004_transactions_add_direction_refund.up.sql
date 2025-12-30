ALTER TYPE transaction_type ADD VALUE IF NOT EXISTS 'REFUND';

CREATE TYPE direction_type AS ENUM(
    'IN',
    'OUT'
);

ALTER TABLE balance_transactions
    ADD COLUMN direction direction_type NOT NULL DEFAULT 'IN';
