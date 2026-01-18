BEGIN;

CREATE TABLE IF NOT EXISTS seats (
    id VARCHAR(10) NOT NULL,
    flight_id UUID NOT NULL REFERENCES flights(id) ON DELETE CASCADE,
    row_num INTEGER NOT NULL,
    col VARCHAR(1) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'available',
    order_id UUID REFERENCES orders(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (flight_id, id),
    CONSTRAINT seats_status_check CHECK (status IN ('available', 'reserved', 'booked'))
);

CREATE INDEX idx_seats_status ON seats(flight_id, status);
CREATE INDEX idx_seats_order ON seats(order_id) WHERE order_id IS NOT NULL;

COMMIT;
