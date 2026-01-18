BEGIN;

-- Insert demo flights
INSERT INTO flights (id, flight_number, origin, destination, departure_time, arrival_time, total_seats, available_seats, price_cents)
VALUES
    ('550e8400-e29b-41d4-a716-446655440001', 'FL101', 'NYC', 'LAX', NOW() + INTERVAL '2 days', NOW() + INTERVAL '2 days' + INTERVAL '6 hours', 120, 120, 35000),
    ('550e8400-e29b-41d4-a716-446655440002', 'FL102', 'LAX', 'NYC', NOW() + INTERVAL '3 days', NOW() + INTERVAL '3 days' + INTERVAL '5 hours', 120, 120, 32000),
    ('550e8400-e29b-41d4-a716-446655440003', 'FL201', 'SFO', 'CHI', NOW() + INTERVAL '1 day', NOW() + INTERVAL '1 day' + INTERVAL '4 hours', 90, 90, 28000),
    ('550e8400-e29b-41d4-a716-446655440004', 'FL202', 'CHI', 'SFO', NOW() + INTERVAL '4 days', NOW() + INTERVAL '4 days' + INTERVAL '4 hours', 90, 90, 27500);

-- Generate seats for FL101 and FL102 (20 rows x 6 seats = 120 seats)
INSERT INTO seats (id, flight_id, row_num, col, status)
SELECT
    row_num || col AS id,
    flight_id,
    row_num,
    col,
    'available'
FROM (
    SELECT
        f.id AS flight_id,
        r.row_num,
        c.col
    FROM flights f
    CROSS JOIN generate_series(1, 20) AS r(row_num)
    CROSS JOIN (VALUES ('A'), ('B'), ('C'), ('D'), ('E'), ('F')) AS c(col)
    WHERE f.flight_number IN ('FL101', 'FL102')
) AS seat_data;

-- Generate seats for FL201 and FL202 (15 rows x 6 seats = 90 seats)
INSERT INTO seats (id, flight_id, row_num, col, status)
SELECT
    row_num || col AS id,
    flight_id,
    row_num,
    col,
    'available'
FROM (
    SELECT
        f.id AS flight_id,
        r.row_num,
        c.col
    FROM flights f
    CROSS JOIN generate_series(1, 15) AS r(row_num)
    CROSS JOIN (VALUES ('A'), ('B'), ('C'), ('D'), ('E'), ('F')) AS c(col)
    WHERE f.flight_number IN ('FL201', 'FL202')
) AS seat_data;

COMMIT;
