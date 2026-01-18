BEGIN;

DELETE FROM seats WHERE flight_id IN (
    SELECT id FROM flights WHERE flight_number IN ('FL101', 'FL102', 'FL201', 'FL202')
);

DELETE FROM flights WHERE flight_number IN ('FL101', 'FL102', 'FL201', 'FL202');

COMMIT;
