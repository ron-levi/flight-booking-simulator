BEGIN;

-- Delete orders first (foreign key constraint)
DELETE FROM orders WHERE flight_id IN (
    SELECT id FROM flights WHERE flight_number IN ('FL101', 'FL102', 'FL201', 'FL202')
);

-- Then delete seats
DELETE FROM seats WHERE flight_id IN (
    SELECT id FROM flights WHERE flight_number IN ('FL101', 'FL102', 'FL201', 'FL202')
);

-- Finally delete flights
DELETE FROM flights WHERE flight_number IN ('FL101', 'FL102', 'FL201', 'FL202');

COMMIT;
