package activities

import (
	"context"
	"fmt"
	"time"

	"github.com/flight-booking-system/internal/domain"
)

// CreateOrderInput contains parameters for creating an order
type CreateOrderInput struct {
	OrderID    string
	FlightID   string
	WorkflowID string
	Seats      []string
	ExpiresAt  time.Time
}

// CreateOrder creates a new order in SEATS_RESERVED status
func (a *BookingActivities) CreateOrder(ctx context.Context, input CreateOrderInput) error {
	// Get flight to calculate price
	flight, err := a.flightRepo.FindByID(ctx, input.FlightID)
	if err != nil {
		return fmt.Errorf("get flight: %w", err)
	}

	// Calculate total price
	totalPrice := flight.PriceCents * int64(len(input.Seats))
	expiresAt := input.ExpiresAt

	order := &domain.Order{
		ID:              input.OrderID,
		FlightID:        input.FlightID,
		WorkflowID:      input.WorkflowID,
		Status:          domain.OrderStatusSeatsReserved,
		Seats:           input.Seats,
		TotalPriceCents: totalPrice,
		ExpiresAt:       &expiresAt,
	}

	if err := a.orderRepo.Create(ctx, order); err != nil {
		return fmt.Errorf("create order: %w", err)
	}

	return nil
}

// UpdateOrderStatusInput contains parameters for status update
type UpdateOrderStatusInput struct {
	OrderID string
	Status  domain.OrderStatus
}

// UpdateOrderStatus updates the order status
func (a *BookingActivities) UpdateOrderStatus(ctx context.Context, input UpdateOrderStatusInput) error {
	if err := a.orderRepo.UpdateStatus(ctx, input.OrderID, input.Status); err != nil {
		return fmt.Errorf("update order status: %w", err)
	}

	return nil
}

// UpdateOrderSeatsInput contains parameters for seat update
type UpdateOrderSeatsInput struct {
	OrderID   string
	Seats     []string
	ExpiresAt time.Time
}

// UpdateOrderSeats updates the order seats and expiration time
func (a *BookingActivities) UpdateOrderSeats(ctx context.Context, input UpdateOrderSeatsInput) error {
	expiresAt := input.ExpiresAt
	if err := a.orderRepo.UpdateSeats(ctx, input.OrderID, input.Seats, &expiresAt); err != nil {
		return fmt.Errorf("update order seats: %w", err)
	}

	return nil
}

// ConfirmOrderInput contains parameters for order confirmation
type ConfirmOrderInput struct {
	OrderID  string
	FlightID string
	Seats    []string
}

// ConfirmOrder marks the order as confirmed and updates flight availability
func (a *BookingActivities) ConfirmOrder(ctx context.Context, input ConfirmOrderInput) error {
	// Confirm the order
	if err := a.orderRepo.Confirm(ctx, input.OrderID); err != nil {
		return fmt.Errorf("confirm order: %w", err)
	}

	// Mark seats as booked in the database
	if err := a.flightRepo.BookSeats(ctx, input.FlightID, input.Seats, input.OrderID); err != nil {
		return fmt.Errorf("book seats: %w", err)
	}

	// Decrease available seats count
	seatCount := len(input.Seats)
	if err := a.flightRepo.UpdateAvailableSeats(ctx, input.FlightID, -seatCount); err != nil {
		return fmt.Errorf("update available seats: %w", err)
	}

	// Release Redis locks since seats are now permanently booked
	_ = a.seatLockRepo.ReleaseLocks(ctx, input.FlightID, input.Seats, input.OrderID)

	return nil
}

// FailOrderInput contains parameters for order failure
type FailOrderInput struct {
	OrderID string
	Reason  string
}

// FailOrder marks the order as failed with a reason
func (a *BookingActivities) FailOrder(ctx context.Context, input FailOrderInput) error {
	if err := a.orderRepo.Fail(ctx, input.OrderID, input.Reason); err != nil {
		return fmt.Errorf("fail order: %w", err)
	}

	return nil
}

// ExpireOrderInput contains parameters for order expiration
type ExpireOrderInput struct {
	OrderID string
}

// ExpireOrder marks the order as expired
func (a *BookingActivities) ExpireOrder(ctx context.Context, input ExpireOrderInput) error {
	if err := a.orderRepo.Expire(ctx, input.OrderID); err != nil {
		return fmt.Errorf("expire order: %w", err)
	}

	return nil
}
