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

// ReserveSeats acquires Redis locks and marks seats as reserved in DB atomically
// TTL is set to 16 minutes (1 min buffer over 15 min workflow timer)
// On failure, compensates by releasing any acquired locks
func (a *BookingActivities) ReserveSeats(ctx context.Context, input ReserveSeatInput) error {
	// Use configured timeout + 1 minute buffer for Redis TTL
	ttl := a.cfg.SeatReservationTimeout + time.Minute

	// Step 1: Acquire Redis locks
	err := a.seatLockRepo.LockSeats(ctx, input.FlightID, input.Seats, input.OrderID, ttl)
	if err != nil {
		return fmt.Errorf("lock seats for order %s: %w", input.OrderID, err)
	}

	// Step 2: Mark seats as reserved in DB
	err = a.flightRepo.MarkSeatsReserved(ctx, input.FlightID, input.Seats, input.OrderID)
	if err != nil {
		// Compensate: release Redis locks
		_ = a.seatLockRepo.ReleaseLocks(ctx, input.FlightID, input.Seats, input.OrderID)
		return fmt.Errorf("mark seats reserved in DB for order %s: %w", input.OrderID, err)
	}

	return nil
}

// ReleaseSeatsInput contains parameters for releasing seats
type ReleaseSeatsInput struct {
	OrderID  string
	FlightID string
	Seats    []string
}

// ReleaseSeats releases Redis locks and marks seats as available in DB
// Only releases if the lock is owned by this order (atomic via Lua script)
func (a *BookingActivities) ReleaseSeats(ctx context.Context, input ReleaseSeatsInput) error {
	// Step 1: Release Redis locks
	err := a.seatLockRepo.ReleaseLocks(ctx, input.FlightID, input.Seats, input.OrderID)
	if err != nil {
		return fmt.Errorf("release seats for order %s: %w", input.OrderID, err)
	}

	// Step 2: Mark seats as available in DB
	err = a.flightRepo.MarkSeatsAvailable(ctx, input.FlightID, input.Seats)
	if err != nil {
		return fmt.Errorf("mark seats available in DB for order %s: %w", input.OrderID, err)
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
// Updates both Redis locks and DB seat status
func (a *BookingActivities) UpdateSeatSelection(ctx context.Context, input UpdateSeatSelectionInput) error {
	ttl := a.cfg.SeatReservationTimeout + time.Minute

	// Release old seats first (Redis + DB)
	if len(input.OldSeats) > 0 {
		if err := a.seatLockRepo.ReleaseLocks(ctx, input.FlightID, input.OldSeats, input.OrderID); err != nil {
			return fmt.Errorf("release old seat locks: %w", err)
		}
		if err := a.flightRepo.MarkSeatsAvailable(ctx, input.FlightID, input.OldSeats); err != nil {
			return fmt.Errorf("mark old seats available: %w", err)
		}
	}

	// Acquire new seats (Redis + DB)
	if len(input.NewSeats) > 0 {
		if err := a.seatLockRepo.LockSeats(ctx, input.FlightID, input.NewSeats, input.OrderID, ttl); err != nil {
			// Try to re-acquire old seats on failure (best effort compensation)
			_ = a.seatLockRepo.LockSeats(ctx, input.FlightID, input.OldSeats, input.OrderID, ttl)
			_ = a.flightRepo.MarkSeatsReserved(ctx, input.FlightID, input.OldSeats, input.OrderID)
			return fmt.Errorf("lock new seats: %w", err)
		}
		if err := a.flightRepo.MarkSeatsReserved(ctx, input.FlightID, input.NewSeats, input.OrderID); err != nil {
			// Compensate: release Redis locks we just acquired
			_ = a.seatLockRepo.ReleaseLocks(ctx, input.FlightID, input.NewSeats, input.OrderID)
			// Re-acquire old seats (best effort)
			_ = a.seatLockRepo.LockSeats(ctx, input.FlightID, input.OldSeats, input.OrderID, ttl)
			_ = a.flightRepo.MarkSeatsReserved(ctx, input.FlightID, input.OldSeats, input.OrderID)
			return fmt.Errorf("mark new seats reserved: %w", err)
		}
	}

	return nil
}

// GetAllFlightIDs returns all flight IDs from the database
func (a *BookingActivities) GetAllFlightIDs(ctx context.Context) ([]string, error) {
	flightIDs, err := a.flightRepo.GetAllFlightIDs(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all flight IDs: %w", err)
	}
	return flightIDs, nil
}

// ReconcileSeatLocksInput contains parameters for reconciling seat locks
type ReconcileSeatLocksInput struct {
	FlightID string
}

// ReconcileSeatLocks reconciles Redis locks with DB seat status
// Releases orphaned Redis locks that don't match DB reserved/booked seats
// This runs periodically to clean up after failures or crashes
func (a *BookingActivities) ReconcileSeatLocks(ctx context.Context, input ReconcileSeatLocksInput) error {
	// Get all Redis locks for this flight
	redisLocks, err := a.seatLockRepo.GetLockedSeats(ctx, input.FlightID)
	if err != nil {
		return fmt.Errorf("get locked seats from Redis: %w", err)
	}

	// Get all DB seats for this flight
	dbSeats, err := a.flightRepo.FindSeats(ctx, input.FlightID)
	if err != nil {
		return fmt.Errorf("get seats from DB: %w", err)
	}

	// Build map of reserved/booked seats in DB with their order IDs
	dbReservedSeats := make(map[string]string)
	for _, seat := range dbSeats {
		if seat.Status == "reserved" || seat.Status == "booked" {
			if seat.OrderID != nil {
				dbReservedSeats[seat.ID] = *seat.OrderID
			}
		}
	}

	// Find orphaned locks (in Redis but not reserved/booked in DB)
	orphanedLocks := make([]string, 0)
	for seatID, redisOrderID := range redisLocks {
		dbOrderID, existsInDB := dbReservedSeats[seatID]
		if !existsInDB || dbOrderID != redisOrderID {
			// Orphaned lock: Redis lock exists but DB shows available or different order
			orphanedLocks = append(orphanedLocks, seatID)
		}
	}

	// Release orphaned locks
	if len(orphanedLocks) > 0 {
		for _, seatID := range orphanedLocks {
			orderID := redisLocks[seatID]
			err := a.seatLockRepo.ReleaseLocks(ctx, input.FlightID, []string{seatID}, orderID)
			if err != nil {
				// Log but continue - best effort cleanup
				continue
			}
		}
	}

	return nil
}
