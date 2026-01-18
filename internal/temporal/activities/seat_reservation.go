package activities

import (
	"context"
	"fmt"
	"time"
)

// ReserveSeatInput contains parameters for seat reservation
type ReserveSeatInput struct {
	OrderID  string
	FlightID string
	Seats    []string
}

// ReserveSeats acquires Redis locks for the specified seats
// TTL is set to 16 minutes (1 min buffer over 15 min workflow timer)
func (a *BookingActivities) ReserveSeats(ctx context.Context, input ReserveSeatInput) error {
	// Use configured timeout + 1 minute buffer for Redis TTL
	ttl := a.cfg.SeatReservationTimeout + time.Minute

	err := a.seatLockRepo.LockSeats(ctx, input.FlightID, input.Seats, input.OrderID, ttl)
	if err != nil {
		return fmt.Errorf("lock seats for order %s: %w", input.OrderID, err)
	}

	return nil
}

// ReleaseSeatsInput contains parameters for releasing seats
type ReleaseSeatsInput struct {
	OrderID  string
	FlightID string
	Seats    []string
}

// ReleaseSeats releases Redis locks for the specified seats
// Only releases if the lock is owned by this order (atomic via Lua script)
func (a *BookingActivities) ReleaseSeats(ctx context.Context, input ReleaseSeatsInput) error {
	err := a.seatLockRepo.ReleaseLocks(ctx, input.FlightID, input.Seats, input.OrderID)
	if err != nil {
		return fmt.Errorf("release seats for order %s: %w", input.OrderID, err)
	}

	return nil
}

// RefreshSeatLocksInput contains parameters for refreshing seat locks
type RefreshSeatLocksInput struct {
	OrderID  string
	FlightID string
	Seats    []string
}

// RefreshSeatLocks extends the TTL for all seat locks
// Called when user updates seat selection to reset the hold timer
func (a *BookingActivities) RefreshSeatLocks(ctx context.Context, input RefreshSeatLocksInput) error {
	// Use configured timeout + 1 minute buffer
	ttl := a.cfg.SeatReservationTimeout + time.Minute

	err := a.seatLockRepo.ExtendLocks(ctx, input.FlightID, input.Seats, input.OrderID, ttl)
	if err != nil {
		return fmt.Errorf("refresh seat locks for order %s: %w", input.OrderID, err)
	}

	return nil
}

// UpdateSeatSelectionInput contains parameters for changing seat selection
type UpdateSeatSelectionInput struct {
	OrderID  string
	FlightID string
	OldSeats []string
	NewSeats []string
}

// UpdateSeatSelection releases old seats and acquires new ones atomically
func (a *BookingActivities) UpdateSeatSelection(ctx context.Context, input UpdateSeatSelectionInput) error {
	ttl := a.cfg.SeatReservationTimeout + time.Minute

	// Release old seats first
	if len(input.OldSeats) > 0 {
		if err := a.seatLockRepo.ReleaseLocks(ctx, input.FlightID, input.OldSeats, input.OrderID); err != nil {
			return fmt.Errorf("release old seats: %w", err)
		}
	}

	// Acquire new seats
	if len(input.NewSeats) > 0 {
		if err := a.seatLockRepo.LockSeats(ctx, input.FlightID, input.NewSeats, input.OrderID, ttl); err != nil {
			// Try to re-acquire old seats on failure (best effort)
			_ = a.seatLockRepo.LockSeats(ctx, input.FlightID, input.OldSeats, input.OrderID, ttl)
			return fmt.Errorf("lock new seats: %w", err)
		}
	}

	return nil
}
