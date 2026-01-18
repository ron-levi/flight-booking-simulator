# Feature: Phase 2 - Temporal Workflows & Activities

The following plan should be complete, but it's important that you validate documentation and codebase patterns and task sanity before you start implementing.

Pay special attention to naming of existing utils, types, and models. Import from the right files etc.

## Feature Description

Implement the core Temporal workflow and activities for the flight booking system:
- **BookingWorkflow**: Long-running workflow managing seat reservation with 15-minute timer, signal-based seat updates, and payment processing
- **SeatReservationActivity**: Acquire/release/refresh Redis distributed locks for seats
- **PaymentValidationActivity**: Validate 5-digit payment codes with 15% simulated failure rate, 10-second timeout, and 3 retries
- **OrderManagementActivity**: Persist order state transitions to PostgreSQL
- **Workflow Queries**: Real-time status queries for frontend polling

This phase builds on the Phase 1 infrastructure to create the workflow orchestration layer.

## User Story

As a flight booking customer
I want my seat selection to be held for 15 minutes while I complete payment
So that I can secure my preferred seats without them being taken by someone else

## Problem Statement

The booking process requires:
- Distributed seat locking with automatic expiration
- Timer-based seat hold that resets on seat selection changes
- Payment validation with retry logic and failure simulation
- Real-time status updates for frontend display
- Compensation logic to release seats on failure/timeout

## Solution Statement

Use Temporal workflow orchestration to:
1. Manage long-running booking workflow with durable timers
2. Handle seat selection changes via signals that reset the hold timer
3. Process payments with configurable retry policies
4. Provide workflow queries for real-time status without side effects
5. Implement saga pattern for automatic cleanup on failure

## Feature Metadata

**Feature Type**: New Capability (Core Business Logic)
**Estimated Complexity**: High
**Primary Systems Affected**: Temporal workers, Redis locks, PostgreSQL orders
**Dependencies**: Phase 1 infrastructure (complete), Temporal SDK v1.26+

---

## CONTEXT REFERENCES

### Relevant Codebase Files - IMPORTANT: YOU MUST READ THESE FILES BEFORE IMPLEMENTING!

**Existing Infrastructure (from Phase 1):**
- `internal/config/config.go` (lines 44-49) - BookingConfig with timeout/retry settings
- `internal/domain/order.go` (lines 5-16) - OrderStatus enum for state transitions
- `internal/domain/order.go` (lines 44-71) - Order.IsTerminal() and CanTransitionTo() validation
- `internal/domain/errors.go` (lines 1-38) - Sentinel errors (ErrSeatUnavailable, ErrPaymentFailed, etc.)
- `internal/repository/seat_lock.go` (lines 28-170) - Redis seat locking operations
- `internal/repository/order_repo.go` (lines 26-193) - Order CRUD operations
- `internal/repository/flight_repo.go` (lines 114-131) - UpdateAvailableSeats method
- `cmd/worker/main.go` (lines 53-57) - Worker setup with TODO for registration

**PRD References:**
- `PRD.md` (lines 262-309) - BookingWorkflow specification with timer/signal logic
- `PRD.md` (lines 310-348) - Payment validation with retry policy
- `PRD.md` (lines 350-373) - Order status and query specification
- `PRD.md` (lines 669-686) - Phase 2 deliverables and validation criteria

### New Files to Create

**Temporal Package:**
- `internal/temporal/workflows/booking_workflow.go` - Main booking workflow
- `internal/temporal/activities/activities.go` - Activity struct with dependencies
- `internal/temporal/activities/seat_reservation.go` - Seat lock activities
- `internal/temporal/activities/payment_validation.go` - Payment processing activity
- `internal/temporal/activities/order_management.go` - Order state activities
- `internal/temporal/signals.go` - Signal and query type definitions
- `internal/temporal/errors.go` - Temporal-specific error types

**Worker Updates:**
- `cmd/worker/main.go` - Update to register workflows and activities

### Relevant Documentation - YOU SHOULD READ THESE BEFORE IMPLEMENTING!

- [Temporal Go SDK - Workflows](https://docs.temporal.io/develop/go/core-application#develop-workflows)
  - Workflow determinism requirements
  - Why: Workflows must be deterministic (no rand, time.Now, etc. directly)

- [Temporal Go SDK - Timers](https://docs.temporal.io/develop/go/timers)
  - workflow.NewTimer and timer cancellation
  - Why: 15-minute seat hold timer implementation

- [Temporal Go SDK - Signals](https://docs.temporal.io/develop/go/message-passing#signals)
  - workflow.GetSignalChannel and signal handling
  - Why: Seat selection updates reset timer

- [Temporal Go SDK - Queries](https://docs.temporal.io/develop/go/message-passing#queries)
  - workflow.SetQueryHandler
  - Why: Real-time status for frontend polling

- [Temporal Go SDK - Activity Retries](https://docs.temporal.io/develop/go/failure-detection#activity-retries)
  - temporal.RetryPolicy configuration
  - Why: Payment validation retry logic

- [Temporal Go SDK - Cancellation](https://docs.temporal.io/develop/go/cancellation)
  - workflow.NewDisconnectedContext for cleanup
  - Why: Saga pattern compensation

- [Temporal Samples - Saga](https://github.com/temporalio/samples-go/tree/main/saga)
  - Compensation pattern implementation
  - Why: Automatic seat release on failure

### Patterns to Follow

**Naming Conventions (from existing codebase):**
```go
// Struct pattern from repository
type BookingActivities struct {
    orderRepo    *repository.OrderRepo
    flightRepo   *repository.FlightRepo
    seatLockRepo *repository.SeatLockRepo
    cfg          *config.BookingConfig
}

// Constructor pattern
func NewBookingActivities(...) *BookingActivities

// Method receiver pattern (short single-letter)
func (a *BookingActivities) ReserveSeats(ctx context.Context, ...) error
```

**Error Handling (from internal/repository):**
```go
// Wrap errors with context
if err != nil {
    return fmt.Errorf("reserve seats for order %s: %w", orderID, err)
}

// Use domain errors for business logic
if !seat.IsAvailable() {
    return domain.ErrSeatUnavailable
}
```

**Config Usage (from config.go:44-49):**
```go
// Access booking config
cfg.Booking.SeatReservationTimeout   // 15 minutes
cfg.Booking.PaymentValidationTimeout // 10 seconds
cfg.Booking.PaymentMaxRetries        // 3
cfg.Booking.PaymentFailureRate       // 0.15
```

**Order Status Transitions (from domain/order.go:52-71):**
```go
// Valid transitions:
// CREATED → SEATS_RESERVED → PAYMENT_PENDING → PAYMENT_PROCESSING → CONFIRMED
//                                                                 → FAILED
//                         → EXPIRED
```

---

## IMPLEMENTATION PLAN

### Phase 2.1: Type Definitions & Errors

Define signal/query types and Temporal-specific errors.

**Tasks:**
- Create signal payload structs (SeatUpdateSignal, PaymentSignal)
- Create query response struct (BookingStatusQuery)
- Define Temporal application errors (ErrReservationExpired, etc.)

### Phase 2.2: Activity Implementation

Implement all activities that interact with repositories.

**Tasks:**
- Create BookingActivities struct with dependencies
- Implement ReserveSeats activity (acquire Redis locks)
- Implement ReleaseSeats activity (release Redis locks)
- Implement RefreshSeatLocks activity (extend TTL)
- Implement ValidatePayment activity (simulated with failure rate)
- Implement order state activities (UpdateStatus, Confirm, Fail, Expire)

### Phase 2.3: Workflow Implementation

Implement the BookingWorkflow with timer, signals, and queries.

**Tasks:**
- Create workflow state struct to track booking progress
- Register query handler for status queries
- Implement seat reservation phase with timer
- Implement signal handling for seat updates (timer reset)
- Implement payment phase with proceed signal
- Implement saga compensation (defer-based cleanup)

### Phase 2.4: Worker Registration

Update worker to register workflows and activities.

**Tasks:**
- Import temporal package
- Create activity instances with dependencies
- Register workflow and activities
- Verify worker starts and connects

---

## STEP-BY-STEP TASKS

IMPORTANT: Execute every task in order, top to bottom. Each task is atomic and independently testable.

---

### Task 1: CREATE `internal/temporal/errors.go`

- **IMPLEMENT**: Temporal-specific error types using temporal.NewApplicationError
- **PATTERN**: Similar to domain/errors.go but for workflow-level errors
- **IMPORTS**: `go.temporal.io/sdk/temporal`
- **VALIDATE**: `go build ./internal/temporal`

```go
package temporal

import (
	"errors"

	"go.temporal.io/sdk/temporal"
)

// Workflow-level errors
var (
	// ErrReservationExpired indicates the 15-minute seat hold timer expired
	ErrReservationExpired = errors.New("seat reservation expired")

	// ErrPaymentTimeout indicates payment validation timed out
	ErrPaymentTimeout = errors.New("payment validation timed out")

	// ErrWorkflowCanceled indicates the workflow was canceled by user
	ErrWorkflowCanceled = errors.New("booking workflow canceled")
)

// Non-retryable error types for Temporal retry policy
const (
	ErrTypeSeatUnavailable    = "SEAT_UNAVAILABLE"
	ErrTypePaymentDeclined    = "PAYMENT_DECLINED"
	ErrTypeInvalidPaymentCode = "INVALID_PAYMENT_CODE"
	ErrTypeOrderExpired       = "ORDER_EXPIRED"
)

// NewSeatUnavailableError creates a non-retryable seat error
func NewSeatUnavailableError(seatID string) error {
	return temporal.NewApplicationErrorWithCause(
		"seat "+seatID+" is not available",
		ErrTypeSeatUnavailable,
		nil,
	)
}

// NewPaymentDeclinedError creates a non-retryable payment error
func NewPaymentDeclinedError(reason string) error {
	return temporal.NewApplicationErrorWithCause(
		reason,
		ErrTypePaymentDeclined,
		nil,
	)
}

// NewInvalidPaymentCodeError creates a non-retryable validation error
func NewInvalidPaymentCodeError() error {
	return temporal.NewApplicationErrorWithCause(
		"payment code must be 5 digits",
		ErrTypeInvalidPaymentCode,
		nil,
	)
}
```

---

### Task 2: CREATE `internal/temporal/signals.go`

- **IMPLEMENT**: Signal payload types and query response types
- **PATTERN**: Plain structs with JSON tags for serialization
- **IMPORTS**: `time`
- **VALIDATE**: `go build ./internal/temporal`

```go
package temporal

import (
	"time"

	"github.com/flight-booking-system/internal/domain"
)

// Signal names as constants
const (
	SignalUpdateSeats = "update-seats"
	SignalProceedToPay = "proceed-to-payment"
	SignalCancelBooking = "cancel-booking"
)

// Query names as constants
const (
	QueryBookingStatus = "booking-status"
)

// SeatUpdateSignal is sent when user changes seat selection
type SeatUpdateSignal struct {
	Seats []string `json:"seats"`
}

// PaymentSignal is sent when user submits payment
type PaymentSignal struct {
	PaymentCode string `json:"paymentCode"`
}

// BookingStatusResponse is returned by the status query
type BookingStatusResponse struct {
	OrderID         string             `json:"orderId"`
	FlightID        string             `json:"flightId"`
	Status          domain.OrderStatus `json:"status"`
	Seats           []string           `json:"seats"`
	ExpiresAt       time.Time          `json:"expiresAt"`
	TimerRemaining  int                `json:"timerRemaining"` // seconds
	PaymentAttempts int                `json:"paymentAttempts"`
	LastError       string             `json:"lastError,omitempty"`
}

// BookingWorkflowInput contains the initial workflow parameters
type BookingWorkflowInput struct {
	OrderID  string   `json:"orderId"`
	FlightID string   `json:"flightId"`
	Seats    []string `json:"seats"`
}

// BookingWorkflowResult contains the workflow completion result
type BookingWorkflowResult struct {
	OrderID   string             `json:"orderId"`
	Status    domain.OrderStatus `json:"status"`
	Seats     []string           `json:"seats"`
	Error     string             `json:"error,omitempty"`
}
```

---

### Task 3: CREATE `internal/temporal/activities/activities.go`

- **IMPLEMENT**: Activities struct with all repository dependencies
- **PATTERN**: Mirror repository constructor pattern (NewXxxRepo)
- **IMPORTS**: Repository types, config, redis client
- **GOTCHA**: Activities need access to all repos for different operations
- **VALIDATE**: `go build ./internal/temporal/activities`

```go
package activities

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/flight-booking-system/internal/config"
	"github.com/flight-booking-system/internal/repository"
)

// BookingActivities contains all activities for the booking workflow
type BookingActivities struct {
	orderRepo    *repository.OrderRepo
	flightRepo   *repository.FlightRepo
	seatLockRepo *repository.SeatLockRepo
	cfg          *config.BookingConfig
}

// NewBookingActivities creates a new BookingActivities instance
func NewBookingActivities(
	pool *pgxpool.Pool,
	redisClient *redis.Client,
	cfg *config.BookingConfig,
) *BookingActivities {
	return &BookingActivities{
		orderRepo:    repository.NewOrderRepo(pool),
		flightRepo:   repository.NewFlightRepo(pool),
		seatLockRepo: repository.NewSeatLockRepo(redisClient),
		cfg:          cfg,
	}
}
```

---

### Task 4: CREATE `internal/temporal/activities/seat_reservation.go`

- **IMPLEMENT**: Seat locking activities using SeatLockRepo
- **PATTERN**: Use existing seat_lock.go methods (lines 28-170)
- **IMPORTS**: context, domain, repository
- **GOTCHA**: TTL should be slightly longer than workflow timer (16 min vs 15 min)
- **VALIDATE**: `go build ./internal/temporal/activities`

```go
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
	OrderID     string
	FlightID    string
	OldSeats    []string
	NewSeats    []string
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
```

---

### Task 5: CREATE `internal/temporal/activities/payment_validation.go`

- **IMPLEMENT**: Payment validation with simulated failure rate
- **PATTERN**: Match PRD lines 314-327 for simulation logic
- **IMPORTS**: context, math/rand, time, temporal
- **GOTCHA**: Use rand from math/rand/v2 for better randomness, seed is auto-initialized in Go 1.20+
- **VALIDATE**: `go build ./internal/temporal/activities`

```go
package activities

import (
	"context"
	"fmt"
	"math/rand"
	"regexp"
	"time"

	"go.temporal.io/sdk/temporal"

	temporalpkg "github.com/flight-booking-system/internal/temporal"
)

// ValidatePaymentInput contains payment validation parameters
type ValidatePaymentInput struct {
	OrderID     string
	PaymentCode string
}

// ValidatePaymentOutput contains the validation result
type ValidatePaymentOutput struct {
	Success bool
	Message string
}

// 5-digit code pattern
var paymentCodePattern = regexp.MustCompile(`^\d{5}$`)

// ValidatePayment simulates payment code validation
// - 15% failure rate (configurable via cfg.PaymentFailureRate)
// - Random processing time 1-8 seconds
// - Returns non-retryable error for invalid code format
func (a *BookingActivities) ValidatePayment(ctx context.Context, input ValidatePaymentInput) (ValidatePaymentOutput, error) {
	// Validate payment code format (5 digits)
	if !paymentCodePattern.MatchString(input.PaymentCode) {
		return ValidatePaymentOutput{}, temporalpkg.NewInvalidPaymentCodeError()
	}

	// Special codes for testing
	switch input.PaymentCode {
	case "00000":
		// Always fails - useful for testing
		return ValidatePaymentOutput{}, temporal.NewApplicationError(
			"payment declined: insufficient funds",
			temporalpkg.ErrTypePaymentDeclined,
		)
	case "99999":
		// Always succeeds instantly - useful for testing
		return ValidatePaymentOutput{Success: true, Message: "Payment validated (test mode)"}, nil
	}

	// Simulate processing time (1-8 seconds)
	processingTime := time.Duration(rand.Intn(7)+1) * time.Second
	select {
	case <-time.After(processingTime):
		// Processing complete
	case <-ctx.Done():
		return ValidatePaymentOutput{}, ctx.Err()
	}

	// Simulate failure rate
	if rand.Float64() < a.cfg.PaymentFailureRate {
		// This error IS retryable (will be retried by Temporal)
		return ValidatePaymentOutput{}, fmt.Errorf("payment validation failed: temporary gateway error")
	}

	return ValidatePaymentOutput{
		Success: true,
		Message: "Payment validated successfully",
	}, nil
}
```

---

### Task 6: CREATE `internal/temporal/activities/order_management.go`

- **IMPLEMENT**: Order state management activities using OrderRepo
- **PATTERN**: Use existing order_repo.go methods (lines 96-193)
- **IMPORTS**: context, domain, time
- **VALIDATE**: `go build ./internal/temporal/activities`

```go
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

	// Decrease available seats count
	seatCount := len(input.Seats)
	if err := a.flightRepo.UpdateAvailableSeats(ctx, input.FlightID, -seatCount); err != nil {
		return fmt.Errorf("update available seats: %w", err)
	}

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
```

---

### Task 7: CREATE `internal/temporal/workflows/booking_workflow.go`

- **IMPLEMENT**: Main booking workflow with timer, signals, queries, and saga compensation
- **PATTERN**: Match PRD lines 262-309 workflow structure
- **IMPORTS**: workflow, temporal, time, errors
- **GOTCHA**: Use workflow.Now(ctx) instead of time.Now() for determinism
- **GOTCHA**: Must drain signal channels before workflow completion
- **GOTCHA**: Use workflow.NewDisconnectedContext for compensation after cancellation
- **VALIDATE**: `go build ./internal/temporal/workflows`

```go
package workflows

import (
	"errors"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	"github.com/flight-booking-system/internal/domain"
	temporalpkg "github.com/flight-booking-system/internal/temporal"
	"github.com/flight-booking-system/internal/temporal/activities"
)

// BookingWorkflow manages the flight booking process
// - Reserves seats with 15-minute timer
// - Handles seat update signals (resets timer)
// - Processes payment on proceed signal
// - Releases seats on timeout/failure/cancellation
func BookingWorkflow(ctx workflow.Context, input temporalpkg.BookingWorkflowInput) (result temporalpkg.BookingWorkflowResult, err error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("BookingWorkflow started", "orderID", input.OrderID, "flightID", input.FlightID)

	// Initialize workflow state
	state := &bookingState{
		orderID:         input.OrderID,
		flightID:        input.FlightID,
		seats:           input.Seats,
		status:          domain.OrderStatusCreated,
		paymentAttempts: 0,
	}

	// Register query handler for status queries
	if err := workflow.SetQueryHandler(ctx, temporalpkg.QueryBookingStatus, func() (temporalpkg.BookingStatusResponse, error) {
		return state.toStatusResponse(), nil
	}); err != nil {
		return result, err
	}

	// Activity options for seat operations (short timeout, retries)
	seatActivityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    10 * time.Second,
			MaximumAttempts:    3,
		},
	}
	seatCtx := workflow.WithActivityOptions(ctx, seatActivityOptions)

	// Activity options for order operations (short timeout, retries)
	orderActivityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    10 * time.Second,
			MaximumAttempts:    3,
		},
	}
	orderCtx := workflow.WithActivityOptions(ctx, orderActivityOptions)

	// Activity options for payment (configured timeout and retries from PRD)
	paymentActivityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:        time.Second,
			BackoffCoefficient:     1.5,
			MaximumInterval:        5 * time.Second,
			MaximumAttempts:        3,
			NonRetryableErrorTypes: []string{
				temporalpkg.ErrTypeInvalidPaymentCode,
				temporalpkg.ErrTypePaymentDeclined,
			},
		},
	}
	paymentCtx := workflow.WithActivityOptions(ctx, paymentActivityOptions)

	var a *activities.BookingActivities

	// Setup compensation for seat release on any failure
	defer func() {
		if err != nil || state.status == domain.OrderStatusExpired || state.status == domain.OrderStatusFailed {
			// Use disconnected context for cleanup (survives workflow cancellation)
			compensationCtx, _ := workflow.NewDisconnectedContext(ctx)
			compensationCtx = workflow.WithActivityOptions(compensationCtx, seatActivityOptions)

			releaseErr := workflow.ExecuteActivity(compensationCtx, a.ReleaseSeats, activities.ReleaseSeatsInput{
				OrderID:  state.orderID,
				FlightID: state.flightID,
				Seats:    state.seats,
			}).Get(compensationCtx, nil)

			if releaseErr != nil {
				logger.Error("Failed to release seats during compensation", "error", releaseErr)
			} else {
				logger.Info("Seats released during compensation", "seats", state.seats)
			}
		}
	}()

	// Phase 1: Reserve seats
	state.status = domain.OrderStatusSeatsReserved
	err = workflow.ExecuteActivity(seatCtx, a.ReserveSeats, activities.ReserveSeatInput{
		OrderID:  input.OrderID,
		FlightID: input.FlightID,
		Seats:    input.Seats,
	}).Get(seatCtx, nil)
	if err != nil {
		state.lastError = err.Error()
		state.status = domain.OrderStatusFailed
		return state.toResult(), err
	}
	logger.Info("Seats reserved", "seats", input.Seats)

	// Create order in database
	state.expiresAt = workflow.Now(ctx).Add(15 * time.Minute)
	err = workflow.ExecuteActivity(orderCtx, a.CreateOrder, activities.CreateOrderInput{
		OrderID:    input.OrderID,
		FlightID:   input.FlightID,
		WorkflowID: workflow.GetInfo(ctx).WorkflowExecution.ID,
		Seats:      input.Seats,
		ExpiresAt:  state.expiresAt,
	}).Get(orderCtx, nil)
	if err != nil {
		state.lastError = err.Error()
		state.status = domain.OrderStatusFailed
		return state.toResult(), err
	}

	// Phase 2: Wait for payment signal with 15-minute timeout
	// Handle seat update signals to reset timer
	seatUpdateChan := workflow.GetSignalChannel(ctx, temporalpkg.SignalUpdateSeats)
	paymentChan := workflow.GetSignalChannel(ctx, temporalpkg.SignalProceedToPay)
	cancelChan := workflow.GetSignalChannel(ctx, temporalpkg.SignalCancelBooking)

	var paymentSignal temporalpkg.PaymentSignal
	paymentReceived := false
	canceled := false

	for !paymentReceived && !canceled {
		// Create timer for remaining hold duration
		timerCtx, cancelTimer := workflow.WithCancel(ctx)
		timerDuration := state.expiresAt.Sub(workflow.Now(ctx))
		if timerDuration <= 0 {
			// Already expired
			state.status = domain.OrderStatusExpired
			state.lastError = "seat reservation expired"
			logger.Info("Seat hold expired")

			// Mark order as expired in database
			_ = workflow.ExecuteActivity(orderCtx, a.ExpireOrder, activities.ExpireOrderInput{
				OrderID: state.orderID,
			}).Get(orderCtx, nil)

			return state.toResult(), temporalpkg.ErrReservationExpired
		}

		holdTimer := workflow.NewTimer(timerCtx, timerDuration)

		selector := workflow.NewSelector(ctx)

		// Handle seat update signal
		selector.AddReceive(seatUpdateChan, func(c workflow.ReceiveChannel, more bool) {
			var signal temporalpkg.SeatUpdateSignal
			c.Receive(ctx, &signal)
			logger.Info("Received seat update signal", "newSeats", signal.Seats)

			// Update seat selection
			updateErr := workflow.ExecuteActivity(seatCtx, a.UpdateSeatSelection, activities.UpdateSeatSelectionInput{
				OrderID:  state.orderID,
				FlightID: state.flightID,
				OldSeats: state.seats,
				NewSeats: signal.Seats,
			}).Get(seatCtx, nil)

			if updateErr != nil {
				logger.Error("Failed to update seats", "error", updateErr)
				state.lastError = updateErr.Error()
			} else {
				state.seats = signal.Seats
				// Reset timer by updating expiration
				state.expiresAt = workflow.Now(ctx).Add(15 * time.Minute)

				// Update order in database
				_ = workflow.ExecuteActivity(orderCtx, a.UpdateOrderSeats, activities.UpdateOrderSeatsInput{
					OrderID:   state.orderID,
					Seats:     signal.Seats,
					ExpiresAt: state.expiresAt,
				}).Get(orderCtx, nil)

				logger.Info("Timer reset", "expiresAt", state.expiresAt)
			}

			cancelTimer() // Cancel current timer to restart with new duration
		})

		// Handle payment signal
		selector.AddReceive(paymentChan, func(c workflow.ReceiveChannel, more bool) {
			c.Receive(ctx, &paymentSignal)
			logger.Info("Received payment signal", "code", paymentSignal.PaymentCode[:2]+"***")
			paymentReceived = true
			cancelTimer()
		})

		// Handle cancel signal
		selector.AddReceive(cancelChan, func(c workflow.ReceiveChannel, more bool) {
			c.Receive(ctx, nil)
			logger.Info("Received cancel signal")
			canceled = true
			cancelTimer()
		})

		// Handle timer expiration
		selector.AddFuture(holdTimer, func(f workflow.Future) {
			timerErr := f.Get(timerCtx, nil)
			if timerErr == nil {
				// Timer actually expired (not canceled)
				state.status = domain.OrderStatusExpired
				state.lastError = "seat reservation expired"
				logger.Info("Seat hold timer expired")
			}
		})

		selector.Select(ctx)

		// Check if expired
		if state.status == domain.OrderStatusExpired {
			// Mark order as expired in database
			_ = workflow.ExecuteActivity(orderCtx, a.ExpireOrder, activities.ExpireOrderInput{
				OrderID: state.orderID,
			}).Get(orderCtx, nil)

			return state.toResult(), temporalpkg.ErrReservationExpired
		}
	}

	// Handle cancellation
	if canceled {
		state.status = domain.OrderStatusFailed
		state.lastError = "booking canceled by user"

		_ = workflow.ExecuteActivity(orderCtx, a.FailOrder, activities.FailOrderInput{
			OrderID: state.orderID,
			Reason:  state.lastError,
		}).Get(orderCtx, nil)

		return state.toResult(), temporalpkg.ErrWorkflowCanceled
	}

	// Phase 3: Process payment
	state.status = domain.OrderStatusPaymentProcessing
	_ = workflow.ExecuteActivity(orderCtx, a.UpdateOrderStatus, activities.UpdateOrderStatusInput{
		OrderID: state.orderID,
		Status:  domain.OrderStatusPaymentProcessing,
	}).Get(orderCtx, nil)

	var paymentResult activities.ValidatePaymentOutput
	err = workflow.ExecuteActivity(paymentCtx, a.ValidatePayment, activities.ValidatePaymentInput{
		OrderID:     state.orderID,
		PaymentCode: paymentSignal.PaymentCode,
	}).Get(paymentCtx, &paymentResult)

	// Track payment attempts (Temporal handles retries internally)
	state.paymentAttempts++

	if err != nil {
		state.status = domain.OrderStatusFailed
		state.lastError = "payment failed: " + err.Error()
		logger.Error("Payment validation failed", "error", err)

		// Check if it's a non-retryable error
		var appErr *temporal.ApplicationError
		if errors.As(err, &appErr) {
			state.lastError = "payment failed: " + appErr.Message()
		}

		_ = workflow.ExecuteActivity(orderCtx, a.FailOrder, activities.FailOrderInput{
			OrderID: state.orderID,
			Reason:  state.lastError,
		}).Get(orderCtx, nil)

		return state.toResult(), err
	}

	// Phase 4: Confirm booking
	state.status = domain.OrderStatusConfirmed
	err = workflow.ExecuteActivity(orderCtx, a.ConfirmOrder, activities.ConfirmOrderInput{
		OrderID:  state.orderID,
		FlightID: state.flightID,
		Seats:    state.seats,
	}).Get(orderCtx, nil)

	if err != nil {
		state.status = domain.OrderStatusFailed
		state.lastError = "confirmation failed: " + err.Error()
		logger.Error("Order confirmation failed", "error", err)

		_ = workflow.ExecuteActivity(orderCtx, a.FailOrder, activities.FailOrderInput{
			OrderID: state.orderID,
			Reason:  state.lastError,
		}).Get(orderCtx, nil)

		return state.toResult(), err
	}

	logger.Info("Booking confirmed", "orderID", state.orderID, "seats", state.seats)

	// Clear the error since compensation is not needed for successful bookings
	err = nil

	// Drain any remaining signals before completing
	drainSignals(ctx, seatUpdateChan, paymentChan, cancelChan)

	return state.toResult(), nil
}

// bookingState tracks the internal workflow state
type bookingState struct {
	orderID         string
	flightID        string
	seats           []string
	status          domain.OrderStatus
	expiresAt       time.Time
	paymentAttempts int
	lastError       string
}

// toStatusResponse converts state to query response
func (s *bookingState) toStatusResponse() temporalpkg.BookingStatusResponse {
	timerRemaining := 0
	if !s.expiresAt.IsZero() {
		remaining := time.Until(s.expiresAt)
		if remaining > 0 {
			timerRemaining = int(remaining.Seconds())
		}
	}

	return temporalpkg.BookingStatusResponse{
		OrderID:         s.orderID,
		FlightID:        s.flightID,
		Status:          s.status,
		Seats:           s.seats,
		ExpiresAt:       s.expiresAt,
		TimerRemaining:  timerRemaining,
		PaymentAttempts: s.paymentAttempts,
		LastError:       s.lastError,
	}
}

// toResult converts state to workflow result
func (s *bookingState) toResult() temporalpkg.BookingWorkflowResult {
	return temporalpkg.BookingWorkflowResult{
		OrderID: s.orderID,
		Status:  s.status,
		Seats:   s.seats,
		Error:   s.lastError,
	}
}

// drainSignals empties signal channels to prevent "unhandled signal" warnings
func drainSignals(ctx workflow.Context, channels ...workflow.ReceiveChannel) {
	for _, ch := range channels {
		for {
			var discard interface{}
			ok := ch.ReceiveAsync(&discard)
			if !ok {
				break
			}
		}
	}
}
```

---

### Task 8: UPDATE `cmd/worker/main.go`

- **IMPLEMENT**: Register workflows and activities with the worker
- **PATTERN**: Follow existing worker setup pattern
- **IMPORTS**: Add temporal packages
- **GOTCHA**: Activities must be registered as struct instance, not just methods
- **VALIDATE**: `go build ./cmd/worker && timeout 10 ./bin/worker || true`

```go
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"github.com/flight-booking-system/internal/config"
	"github.com/flight-booking-system/internal/database"
	"github.com/flight-booking-system/internal/temporal/activities"
	"github.com/flight-booking-system/internal/temporal/workflows"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Connect to PostgreSQL (workers need database access for activities)
	pool, err := database.NewPostgresPool(ctx, cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer pool.Close()
	log.Println("Connected to PostgreSQL")

	// Connect to Redis
	redisClient, err := database.NewRedisClient(ctx, cfg.Redis)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()
	log.Println("Connected to Redis")

	// Connect to Temporal
	temporalClient, err := client.Dial(client.Options{
		HostPort:  cfg.Temporal.Host,
		Namespace: cfg.Temporal.Namespace,
	})
	if err != nil {
		log.Fatalf("Failed to connect to Temporal: %v", err)
	}
	defer temporalClient.Close()
	log.Println("Connected to Temporal")

	// Create worker
	w := worker.New(temporalClient, cfg.Temporal.TaskQueue, worker.Options{})

	// Register workflows
	w.RegisterWorkflow(workflows.BookingWorkflow)

	// Create and register activities
	bookingActivities := activities.NewBookingActivities(pool, redisClient, &cfg.Booking)
	w.RegisterActivity(bookingActivities)

	log.Println("Registered workflows and activities")

	// Start worker in goroutine
	go func() {
		log.Printf("Worker starting on task queue: %s", cfg.Temporal.TaskQueue)
		if err := w.Run(worker.InterruptCh()); err != nil {
			log.Fatalf("Worker failed: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down worker...")
	w.Stop()
	log.Println("Worker stopped")
}
```

---

### Task 9: CREATE Unit Test `internal/temporal/workflows/booking_workflow_test.go`

- **IMPLEMENT**: Basic workflow test using Temporal's test framework
- **PATTERN**: Use testsuite.WorkflowTestSuite and mock activities
- **IMPORTS**: testing, testsuite, require
- **VALIDATE**: `go test -v ./internal/temporal/workflows/...`

```go
package workflows_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/testsuite"

	"github.com/flight-booking-system/internal/domain"
	temporalpkg "github.com/flight-booking-system/internal/temporal"
	"github.com/flight-booking-system/internal/temporal/activities"
	"github.com/flight-booking-system/internal/temporal/workflows"
)

func TestBookingWorkflow_Success(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	// Mock activities
	env.OnActivity((*activities.BookingActivities).ReserveSeats, mock.Anything, mock.Anything).Return(nil)
	env.OnActivity((*activities.BookingActivities).CreateOrder, mock.Anything, mock.Anything).Return(nil)
	env.OnActivity((*activities.BookingActivities).UpdateOrderStatus, mock.Anything, mock.Anything).Return(nil)
	env.OnActivity((*activities.BookingActivities).ValidatePayment, mock.Anything, mock.Anything).Return(
		activities.ValidatePaymentOutput{Success: true, Message: "OK"}, nil,
	)
	env.OnActivity((*activities.BookingActivities).ConfirmOrder, mock.Anything, mock.Anything).Return(nil)
	env.OnActivity((*activities.BookingActivities).ReleaseSeats, mock.Anything, mock.Anything).Return(nil)

	// Send payment signal after workflow starts
	env.RegisterDelayedCallback(func() {
		env.SignalWorkflow(temporalpkg.SignalProceedToPay, temporalpkg.PaymentSignal{
			PaymentCode: "12345",
		})
	}, time.Second)

	// Execute workflow
	env.ExecuteWorkflow(workflows.BookingWorkflow, temporalpkg.BookingWorkflowInput{
		OrderID:  "test-order-1",
		FlightID: "test-flight-1",
		Seats:    []string{"1A", "1B"},
	})

	require.True(t, env.IsWorkflowCompleted())
	require.NoError(t, env.GetWorkflowError())

	var result temporalpkg.BookingWorkflowResult
	require.NoError(t, env.GetWorkflowResult(&result))
	require.Equal(t, domain.OrderStatusConfirmed, result.Status)
	require.Equal(t, "test-order-1", result.OrderID)
}

func TestBookingWorkflow_TimerExpired(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	// Mock activities
	env.OnActivity((*activities.BookingActivities).ReserveSeats, mock.Anything, mock.Anything).Return(nil)
	env.OnActivity((*activities.BookingActivities).CreateOrder, mock.Anything, mock.Anything).Return(nil)
	env.OnActivity((*activities.BookingActivities).ExpireOrder, mock.Anything, mock.Anything).Return(nil)
	env.OnActivity((*activities.BookingActivities).ReleaseSeats, mock.Anything, mock.Anything).Return(nil)

	// Don't send payment signal - let timer expire

	// Execute workflow
	env.ExecuteWorkflow(workflows.BookingWorkflow, temporalpkg.BookingWorkflowInput{
		OrderID:  "test-order-2",
		FlightID: "test-flight-1",
		Seats:    []string{"2A"},
	})

	require.True(t, env.IsWorkflowCompleted())
	require.Error(t, env.GetWorkflowError())

	var result temporalpkg.BookingWorkflowResult
	require.NoError(t, env.GetWorkflowResult(&result))
	require.Equal(t, domain.OrderStatusExpired, result.Status)
}

func TestBookingWorkflow_SeatUpdateResetsTimer(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	// Mock activities
	env.OnActivity((*activities.BookingActivities).ReserveSeats, mock.Anything, mock.Anything).Return(nil)
	env.OnActivity((*activities.BookingActivities).CreateOrder, mock.Anything, mock.Anything).Return(nil)
	env.OnActivity((*activities.BookingActivities).UpdateSeatSelection, mock.Anything, mock.Anything).Return(nil)
	env.OnActivity((*activities.BookingActivities).UpdateOrderSeats, mock.Anything, mock.Anything).Return(nil)
	env.OnActivity((*activities.BookingActivities).UpdateOrderStatus, mock.Anything, mock.Anything).Return(nil)
	env.OnActivity((*activities.BookingActivities).ValidatePayment, mock.Anything, mock.Anything).Return(
		activities.ValidatePaymentOutput{Success: true, Message: "OK"}, nil,
	)
	env.OnActivity((*activities.BookingActivities).ConfirmOrder, mock.Anything, mock.Anything).Return(nil)
	env.OnActivity((*activities.BookingActivities).ReleaseSeats, mock.Anything, mock.Anything).Return(nil)

	// Send seat update signal at 14 minutes (would expire at 15 min)
	env.RegisterDelayedCallback(func() {
		env.SignalWorkflow(temporalpkg.SignalUpdateSeats, temporalpkg.SeatUpdateSignal{
			Seats: []string{"3A", "3B"},
		})
	}, 14*time.Minute)

	// Send payment signal at 16 minutes (after original timeout but before new timeout)
	env.RegisterDelayedCallback(func() {
		env.SignalWorkflow(temporalpkg.SignalProceedToPay, temporalpkg.PaymentSignal{
			PaymentCode: "12345",
		})
	}, 16*time.Minute)

	// Execute workflow
	env.ExecuteWorkflow(workflows.BookingWorkflow, temporalpkg.BookingWorkflowInput{
		OrderID:  "test-order-3",
		FlightID: "test-flight-1",
		Seats:    []string{"1A", "1B"},
	})

	require.True(t, env.IsWorkflowCompleted())
	require.NoError(t, env.GetWorkflowError())

	var result temporalpkg.BookingWorkflowResult
	require.NoError(t, env.GetWorkflowResult(&result))
	require.Equal(t, domain.OrderStatusConfirmed, result.Status)
	require.Equal(t, []string{"3A", "3B"}, result.Seats)
}

func TestBookingWorkflow_QueryStatus(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	// Mock activities
	env.OnActivity((*activities.BookingActivities).ReserveSeats, mock.Anything, mock.Anything).Return(nil)
	env.OnActivity((*activities.BookingActivities).CreateOrder, mock.Anything, mock.Anything).Return(nil)
	env.OnActivity((*activities.BookingActivities).UpdateOrderStatus, mock.Anything, mock.Anything).Return(nil)
	env.OnActivity((*activities.BookingActivities).ValidatePayment, mock.Anything, mock.Anything).Return(
		activities.ValidatePaymentOutput{Success: true, Message: "OK"}, nil,
	)
	env.OnActivity((*activities.BookingActivities).ConfirmOrder, mock.Anything, mock.Anything).Return(nil)
	env.OnActivity((*activities.BookingActivities).ReleaseSeats, mock.Anything, mock.Anything).Return(nil)

	// Query status during workflow execution
	env.RegisterDelayedCallback(func() {
		result, err := env.QueryWorkflow(temporalpkg.QueryBookingStatus)
		require.NoError(t, err)

		var status temporalpkg.BookingStatusResponse
		require.NoError(t, result.Get(&status))
		require.Equal(t, "test-order-4", status.OrderID)
		require.Equal(t, domain.OrderStatusSeatsReserved, status.Status)
		require.True(t, status.TimerRemaining > 0)

		// Now send payment
		env.SignalWorkflow(temporalpkg.SignalProceedToPay, temporalpkg.PaymentSignal{
			PaymentCode: "12345",
		})
	}, time.Second)

	// Execute workflow
	env.ExecuteWorkflow(workflows.BookingWorkflow, temporalpkg.BookingWorkflowInput{
		OrderID:  "test-order-4",
		FlightID: "test-flight-1",
		Seats:    []string{"4A"},
	})

	require.True(t, env.IsWorkflowCompleted())
	require.NoError(t, env.GetWorkflowError())
}
```

---

### Task 10: ADD Test Dependencies to `go.mod`

- **IMPLEMENT**: Add testify for test assertions
- **VALIDATE**: `go mod tidy && go mod verify`

Run the following command to add test dependency:

```bash
go get github.com/stretchr/testify@v1.9.0
go mod tidy
```

---

## TESTING STRATEGY

### Unit Tests

Test workflow logic using Temporal's test framework with mocked activities.

**Workflow Test Cases:**
1. `TestBookingWorkflow_Success` - Happy path with payment signal
2. `TestBookingWorkflow_TimerExpired` - Timer expires without payment
3. `TestBookingWorkflow_SeatUpdateResetsTimer` - Signal resets timer
4. `TestBookingWorkflow_PaymentFailed` - Payment validation fails after retries
5. `TestBookingWorkflow_Canceled` - User cancels booking
6. `TestBookingWorkflow_QueryStatus` - Query returns correct state

**Activity Test Cases:**
1. `TestReserveSeats_Success` - Locks acquired successfully
2. `TestReserveSeats_AlreadyLocked` - Returns error when seat locked
3. `TestValidatePayment_InvalidCode` - Non-retryable error for bad format
4. `TestValidatePayment_DeclinedCode` - Non-retryable error for "00000"

### Integration Tests

Test worker connects to Temporal and can execute workflows.

**Integration Test Cases:**
1. Worker starts and registers workflows
2. Workflow can be started via Temporal client
3. Signal can be sent to running workflow
4. Query returns workflow state

### Edge Cases

- Concurrent seat selection by multiple orders
- Network partition during payment processing
- Workflow continues after worker restart
- Timer drift over long durations

---

## VALIDATION COMMANDS

Execute every command to ensure zero regressions and 100% feature correctness.

### Level 1: Syntax & Build

```bash
# Format code
go fmt ./...

# Verify module
go mod tidy
go mod verify

# Build all packages
go build ./...
```

### Level 2: Unit Tests

```bash
# Run all tests
go test -v ./...

# Run workflow tests specifically
go test -v ./internal/temporal/workflows/...

# Run with race detector
go test -race ./...
```

### Level 3: Infrastructure Check

```bash
# Ensure Docker services are running
docker compose ps

# Verify Temporal is healthy
curl -s http://localhost:8233/api/v1/namespaces | head -5

# Verify Redis is ready
docker exec flight-redis redis-cli ping
```

### Level 4: Worker Integration

```bash
# Build and start worker (brief test)
make build
timeout 15 ./bin/worker 2>&1 | grep -E "(Connected|Registered|Worker starting)"

# Expected output should include:
# - Connected to PostgreSQL
# - Connected to Redis
# - Connected to Temporal
# - Registered workflows and activities
# - Worker starting on task queue: booking-queue
```

### Level 5: Manual Workflow Test

```bash
# Start worker in background
./bin/worker &
WORKER_PID=$!

# Use Temporal CLI to start a workflow
temporal workflow start \
  --task-queue booking-queue \
  --type BookingWorkflow \
  --input '{"orderId":"test-123","flightId":"550e8400-e29b-41d4-a716-446655440001","seats":["1A","1B"]}'

# Query workflow status
temporal workflow query \
  --workflow-id <workflow-id-from-above> \
  --type booking-status

# Send payment signal
temporal workflow signal \
  --workflow-id <workflow-id-from-above> \
  --name proceed-to-payment \
  --input '{"paymentCode":"12345"}'

# Check workflow completion
temporal workflow describe \
  --workflow-id <workflow-id-from-above>

# Stop worker
kill $WORKER_PID
```

---

## ACCEPTANCE CRITERIA

- [ ] BookingWorkflow implements timer-based 15-minute seat hold
- [ ] Seat update signal resets the hold timer
- [ ] Payment signal proceeds to payment validation
- [ ] Cancel signal releases seats and fails order
- [ ] Timer expiration releases seats and expires order
- [ ] Payment validation has 15% simulated failure rate
- [ ] Payment retries 3 times with exponential backoff
- [ ] Payment code "00000" always fails (for testing)
- [ ] Payment code "99999" always succeeds instantly (for testing)
- [ ] Query handler returns current workflow state
- [ ] Activities correctly interact with PostgreSQL and Redis
- [ ] Worker registers and runs workflows without errors
- [ ] Unit tests pass with mocked activities
- [ ] Code follows CLAUDE.md conventions

---

## COMPLETION CHECKLIST

- [ ] All 10 tasks completed in order
- [ ] Each task validation passed
- [ ] `go build ./...` succeeds with no errors
- [ ] `go test -v ./...` passes all tests
- [ ] Worker connects to Temporal successfully
- [ ] Workflow can be started via CLI
- [ ] Status query returns valid response
- [ ] Signals update workflow state correctly
- [ ] Timer expiration triggers compensation
- [ ] Code follows existing codebase patterns

---

## NOTES

### Design Decisions

1. **Activity struct pattern**: All activities belong to `BookingActivities` struct with repo dependencies, matching existing repository patterns

2. **Timer reset implementation**: Instead of trying to "reset" a timer, we cancel the existing timer context and create a new timer with the updated duration

3. **Saga compensation via defer**: Using Go's defer with workflow.NewDisconnectedContext ensures cleanup runs even on workflow cancellation

4. **Non-retryable errors**: Payment code format and "declined" errors are marked non-retryable to fail immediately

5. **Special test codes**: "00000" (always fail) and "99999" (always succeed) enable deterministic testing

### Gotchas

1. **Workflow determinism**: Never use `time.Now()`, `rand.Float64()` directly in workflow code - use `workflow.Now(ctx)` and activities for random behavior

2. **Signal draining**: Always drain signal channels before workflow completion to avoid "unhandled signal" warnings

3. **Timer state**: The `bookingState.expiresAt` uses workflow time for status queries, but actual timer uses `time.Until()` which Temporal handles correctly

4. **Activity retries**: Payment activity retry count in `state.paymentAttempts` only increments once since Temporal handles retries internally

5. **Import alias**: Use `temporalpkg` alias for `internal/temporal` to avoid conflict with `go.temporal.io/sdk/temporal`

### Future Considerations (Phase 3)

- Add HTTP handlers to start workflows and send signals via API
- Add real-time updates via SSE or polling endpoint
- Add workflow search/listing capabilities
- Add admin endpoints for manual workflow management
