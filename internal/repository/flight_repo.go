package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/flight-booking-system/internal/domain"
)

// FlightRepo handles flight data access
type FlightRepo struct {
	pool *pgxpool.Pool
}

// NewFlightRepo creates a new FlightRepo
func NewFlightRepo(pool *pgxpool.Pool) *FlightRepo {
	return &FlightRepo{pool: pool}
}

// FindAll returns all flights
func (r *FlightRepo) FindAll(ctx context.Context) ([]domain.Flight, error) {
	query := `
		SELECT id, flight_number, origin, destination, departure_time, arrival_time,
		       total_seats, available_seats, price_cents, created_at, updated_at
		FROM flights
		ORDER BY departure_time ASC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query flights: %w", err)
	}
	defer rows.Close()

	var flights []domain.Flight
	for rows.Next() {
		var f domain.Flight
		err := rows.Scan(
			&f.ID, &f.FlightNumber, &f.Origin, &f.Destination,
			&f.DepartureTime, &f.ArrivalTime, &f.TotalSeats,
			&f.AvailableSeats, &f.PriceCents, &f.CreatedAt, &f.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan flight: %w", err)
		}
		flights = append(flights, f)
	}

	return flights, rows.Err()
}

// FindByID returns a flight by ID
func (r *FlightRepo) FindByID(ctx context.Context, id string) (*domain.Flight, error) {
	query := `
		SELECT id, flight_number, origin, destination, departure_time, arrival_time,
		       total_seats, available_seats, price_cents, created_at, updated_at
		FROM flights
		WHERE id = $1
	`

	var f domain.Flight
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&f.ID, &f.FlightNumber, &f.Origin, &f.Destination,
		&f.DepartureTime, &f.ArrivalTime, &f.TotalSeats,
		&f.AvailableSeats, &f.PriceCents, &f.CreatedAt, &f.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrFlightNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("query flight: %w", err)
	}

	return &f, nil
}

// FindSeats returns all seats for a flight
func (r *FlightRepo) FindSeats(ctx context.Context, flightID string) ([]domain.Seat, error) {
	query := `
		SELECT id, flight_id, row_num, col, status, order_id, created_at, updated_at
		FROM seats
		WHERE flight_id = $1
		ORDER BY row_num, col
	`

	rows, err := r.pool.Query(ctx, query, flightID)
	if err != nil {
		return nil, fmt.Errorf("query seats: %w", err)
	}
	defer rows.Close()

	var seats []domain.Seat
	for rows.Next() {
		var s domain.Seat
		err := rows.Scan(
			&s.ID, &s.FlightID, &s.Row, &s.Column,
			&s.Status, &s.OrderID, &s.CreatedAt, &s.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan seat: %w", err)
		}
		seats = append(seats, s)
	}

	return seats, rows.Err()
}

// UpdateAvailableSeats updates the available seat count
func (r *FlightRepo) UpdateAvailableSeats(ctx context.Context, flightID string, delta int) error {
	query := `
		UPDATE flights
		SET available_seats = available_seats + $1, updated_at = NOW()
		WHERE id = $2 AND available_seats + $1 >= 0
	`

	result, err := r.pool.Exec(ctx, query, delta, flightID)
	if err != nil {
		return fmt.Errorf("update available seats: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrInsufficientSeats
	}

	return nil
}

// BookSeats marks seats as booked and assigns them to an order
func (r *FlightRepo) BookSeats(ctx context.Context, flightID string, seatIDs []string, orderID string) error {
	query := `
		UPDATE seats
		SET status = 'booked', order_id = $1, updated_at = NOW()
		WHERE flight_id = $2 AND id = ANY($3)
	`

	result, err := r.pool.Exec(ctx, query, orderID, flightID, seatIDs)
	if err != nil {
		return fmt.Errorf("book seats: %w", err)
	}

	if result.RowsAffected() != int64(len(seatIDs)) {
		return fmt.Errorf("expected to book %d seats, but booked %d", len(seatIDs), result.RowsAffected())
	}

	return nil
}
