package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/flight-booking-system/internal/domain"
)

// OrderRepo handles order data access
type OrderRepo struct {
	pool *pgxpool.Pool
}

// NewOrderRepo creates a new OrderRepo
func NewOrderRepo(pool *pgxpool.Pool) *OrderRepo {
	return &OrderRepo{pool: pool}
}

// Create creates a new order
func (r *OrderRepo) Create(ctx context.Context, order *domain.Order) error {
	query := `
		INSERT INTO orders (id, flight_id, workflow_id, status, seats, total_price_cents, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.pool.Exec(ctx, query,
		order.ID, order.FlightID, order.WorkflowID, order.Status,
		order.Seats, order.TotalPriceCents, order.ExpiresAt,
	)
	if err != nil {
		return fmt.Errorf("insert order: %w", err)
	}

	return nil
}

// FindByID returns an order by ID
func (r *OrderRepo) FindByID(ctx context.Context, id string) (*domain.Order, error) {
	query := `
		SELECT id, flight_id, workflow_id, status, seats, total_price_cents,
		       payment_code, expires_at, confirmed_at, failure_reason, created_at, updated_at
		FROM orders
		WHERE id = $1
	`

	var o domain.Order
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&o.ID, &o.FlightID, &o.WorkflowID, &o.Status, &o.Seats,
		&o.TotalPriceCents, &o.PaymentCode, &o.ExpiresAt,
		&o.ConfirmedAt, &o.FailureReason, &o.CreatedAt, &o.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrOrderNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("query order: %w", err)
	}

	return &o, nil
}

// FindByWorkflowID returns an order by workflow ID
func (r *OrderRepo) FindByWorkflowID(ctx context.Context, workflowID string) (*domain.Order, error) {
	query := `
		SELECT id, flight_id, workflow_id, status, seats, total_price_cents,
		       payment_code, expires_at, confirmed_at, failure_reason, created_at, updated_at
		FROM orders
		WHERE workflow_id = $1
	`

	var o domain.Order
	err := r.pool.QueryRow(ctx, query, workflowID).Scan(
		&o.ID, &o.FlightID, &o.WorkflowID, &o.Status, &o.Seats,
		&o.TotalPriceCents, &o.PaymentCode, &o.ExpiresAt,
		&o.ConfirmedAt, &o.FailureReason, &o.CreatedAt, &o.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrOrderNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("query order: %w", err)
	}

	return &o, nil
}

// UpdateStatus updates the order status
func (r *OrderRepo) UpdateStatus(ctx context.Context, id string, status domain.OrderStatus) error {
	query := `
		UPDATE orders
		SET status = $1, updated_at = NOW()
		WHERE id = $2
	`

	result, err := r.pool.Exec(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("update order status: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrOrderNotFound
	}

	return nil
}

// UpdateSeats updates the order seats and expiration
func (r *OrderRepo) UpdateSeats(ctx context.Context, id string, seats []string, expiresAt *time.Time) error {
	query := `
		UPDATE orders
		SET seats = $1, expires_at = $2, updated_at = NOW()
		WHERE id = $3
	`

	result, err := r.pool.Exec(ctx, query, seats, expiresAt, id)
	if err != nil {
		return fmt.Errorf("update order seats: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrOrderNotFound
	}

	return nil
}

// Confirm marks the order as confirmed
func (r *OrderRepo) Confirm(ctx context.Context, id string) error {
	query := `
		UPDATE orders
		SET status = 'CONFIRMED', confirmed_at = NOW(), updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("confirm order: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrOrderNotFound
	}

	return nil
}

// Fail marks the order as failed
func (r *OrderRepo) Fail(ctx context.Context, id string, reason string) error {
	query := `
		UPDATE orders
		SET status = 'FAILED', failure_reason = $1, updated_at = NOW()
		WHERE id = $2
	`

	result, err := r.pool.Exec(ctx, query, reason, id)
	if err != nil {
		return fmt.Errorf("fail order: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrOrderNotFound
	}

	return nil
}

// Expire marks the order as expired
func (r *OrderRepo) Expire(ctx context.Context, id string) error {
	query := `
		UPDATE orders
		SET status = 'EXPIRED', updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("expire order: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrOrderNotFound
	}

	return nil
}
