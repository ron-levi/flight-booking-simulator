BEGIN;

CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    flight_id UUID NOT NULL REFERENCES flights(id),
    workflow_id VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'CREATED',
    seats TEXT[] NOT NULL DEFAULT '{}',
    total_price_cents BIGINT NOT NULL DEFAULT 0,
    payment_code VARCHAR(5),
    expires_at TIMESTAMPTZ,
    confirmed_at TIMESTAMPTZ,
    failure_reason TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT orders_workflow_id_unique UNIQUE (workflow_id),
    CONSTRAINT orders_status_check CHECK (status IN (
        'CREATED', 'SEATS_RESERVED', 'PAYMENT_PENDING',
        'PAYMENT_PROCESSING', 'CONFIRMED', 'FAILED', 'EXPIRED'
    ))
);

CREATE INDEX idx_orders_flight ON orders(flight_id);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_expires ON orders(expires_at) WHERE status IN ('SEATS_RESERVED', 'PAYMENT_PENDING');

COMMIT;
