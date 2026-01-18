# Feature: Complete Phase 3 - REST API Layer with Services

The following plan should be complete, but its important that you validate documentation and codebase patterns and task sanity before you start implementing.

Pay special attention to naming of existing utils types and models. Import from the right files etc.

## Feature Description

Implement the complete REST API layer for the flight booking system, including a services layer that sits between HTTP handlers and repositories/Temporal. This implements Phase 3 of the PRD, exposing all booking functionality through HTTP endpoints that interact with Temporal workflows for booking orchestration.

## User Story

As a flight booking customer,
I want to browse flights, select seats, and complete bookings through a web API,
So that I can book flights with real-time feedback on reservation timers and payment status.

## Problem Statement

The backend has complete Temporal workflows and data repositories, but lacks the HTTP API layer needed for frontend integration. Currently, only a basic `/api/flights` endpoint exists inline in `main.go`. The system needs:
- Full RESTful API endpoints for flights, orders, seats, and payments
- A services layer to encapsulate business logic and Temporal client interactions
- Proper request validation and error handling
- JSON request/response handling following existing patterns

## Solution Statement

Create a complete API layer following the existing codebase patterns:
1. **Services layer** (`internal/service/`) - Business logic orchestration, Temporal client management
2. **Handlers layer** (`internal/api/handlers.go`) - HTTP request/response handling
3. **Routes configuration** (`internal/api/routes.go`) - Chi router setup with route groups
4. **Request/Response types** (`internal/api/types.go`) - API-specific DTOs
5. **Error handling** (`internal/api/errors.go`) - Standardized API error responses

## Feature Metadata

**Feature Type**: New Capability
**Estimated Complexity**: Medium-High
**Primary Systems Affected**: API layer, Service layer, Server entrypoint
**Dependencies**: Chi Router v5, Temporal Go SDK, existing repositories

---

## CONTEXT REFERENCES

### Relevant Codebase Files - IMPORTANT: YOU MUST READ THESE FILES BEFORE IMPLEMENTING!

**Domain Models** (understand existing structures):
- `internal/domain/flight.go` - Flight, FlightWithSeats, SeatMap structs
- `internal/domain/order.go` - Order, OrderStatus, OrderStatusResponse structs
- `internal/domain/seat.go` - Seat, SeatStatus structs
- `internal/domain/errors.go` - All domain error definitions

**Temporal Integration** (understand workflow interaction):
- `internal/temporal/signals.go` - Signal/query names, input/output types (BookingWorkflowInput, BookingStatusResponse, etc.)
- `internal/temporal/errors.go` - Error types for Temporal operations
- `internal/temporal/workflows/booking_workflow.go` - Workflow implementation details

**Repositories** (understand data access patterns):
- `internal/repository/flight_repo.go` - FlightRepo methods (FindAll, FindByID, FindSeats)
- `internal/repository/order_repo.go` - OrderRepo methods (Create, FindByID, UpdateStatus, etc.)
- `internal/repository/seat_lock.go` - SeatLockRepo for Redis locking

**Configuration** (understand config patterns):
- `internal/config/config.go` - Config struct with Server, Database, Redis, Temporal, Booking configs

**Current Server** (understand existing setup):
- `cmd/server/main.go` - Current server with Chi router, middleware, inline handlers

**Worker for Reference** (understand Temporal client creation):
- `cmd/worker/main.go` - Shows Temporal client.Dial pattern (lines 44-52)

### New Files to Create

```
internal/
├── api/
│   ├── handlers.go      # HTTP handlers for all endpoints
│   ├── routes.go        # Chi router configuration
│   ├── types.go         # Request/response DTOs
│   ├── errors.go        # API error handling utilities
│   └── middleware.go    # Custom middleware (CORS, etc.)
└── service/
    ├── flight_service.go   # Flight listing and details
    ├── booking_service.go  # Order management, Temporal interaction
    └── temporal_client.go  # Temporal client wrapper
```

### Relevant Documentation - YOU SHOULD READ THESE BEFORE IMPLEMENTING!

- [Chi Router Documentation](https://pkg.go.dev/github.com/go-chi/chi/v5)
  - Route grouping with `r.Route()`
  - URL parameters with `chi.URLParam(r, "id")`
- [Temporal Go SDK Client](https://pkg.go.dev/go.temporal.io/sdk/client)
  - `ExecuteWorkflow` - Start workflows
  - `SignalWorkflow` - Send signals to running workflows
  - `QueryWorkflow` - Query workflow state
- [Temporal Message Passing](https://docs.temporal.io/develop/go/message-passing)
  - Signal and query patterns from client side

### Patterns to Follow

**Naming Conventions:**
- Files: `snake_case.go`
- Structs: `PascalCase` (e.g., `FlightService`, `CreateOrderRequest`)
- JSON fields: `camelCase` (e.g., `flightId`, `totalSeats`)
- Error variables: `Err` prefix (e.g., `ErrFlightNotFound`)

**Repository Pattern:**
```go
type FlightRepo struct {
    pool *pgxpool.Pool
}

func NewFlightRepo(pool *pgxpool.Pool) *FlightRepo {
    return &FlightRepo{pool: pool}
}
```

**Error Handling Pattern:**
```go
if errors.Is(err, pgx.ErrNoRows) {
    return nil, domain.ErrFlightNotFound
}
if err != nil {
    return nil, fmt.Errorf("query flight: %w", err)
}
```

**JSON Response Pattern:**
```go
w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(response)
```

**Temporal Client Pattern (from worker):**
```go
temporalClient, err := client.Dial(client.Options{
    HostPort:  cfg.Temporal.Host,
    Namespace: cfg.Temporal.Namespace,
})
```

---

## IMPLEMENTATION PLAN

### Phase 1: Foundation - API Types and Error Handling

Create the foundational types for API requests/responses and standardized error handling.

**Tasks:**
- Create API request/response DTOs in `types.go`
- Create error response utilities in `errors.go`
- Create custom middleware for CORS in `middleware.go`

### Phase 2: Services Layer

Create business logic services that encapsulate repository access and Temporal interactions.

**Tasks:**
- Create Temporal client wrapper (`temporal_client.go`)
- Create flight service (`flight_service.go`)
- Create booking service (`booking_service.go`)

### Phase 3: HTTP Handlers

Create thin HTTP handlers that parse requests, call services, and write responses.

**Tasks:**
- Create handlers struct with dependencies
- Implement all endpoint handlers
- Create router configuration

### Phase 4: Integration

Wire everything together in the server entrypoint.

**Tasks:**
- Update `cmd/server/main.go` to use new API package
- Add Temporal client to server
- Remove inline handlers

---

## STEP-BY-STEP TASKS

IMPORTANT: Execute every task in order, top to bottom. Each task is atomic and independently testable.

### Task 1: CREATE `internal/api/types.go`

Create API-specific request and response types.

**IMPLEMENT:**
```go
package api

import "time"

// Request types

type CreateOrderRequest struct {
    FlightID string   `json:"flightId"`
    Seats    []string `json:"seats"`
}

type UpdateSeatsRequest struct {
    Seats []string `json:"seats"`
}

type SubmitPaymentRequest struct {
    PaymentCode string `json:"paymentCode"`
}

// Response types

type FlightListResponse struct {
    Flights []FlightResponse `json:"flights"`
}

type FlightResponse struct {
    ID             string    `json:"id"`
    FlightNumber   string    `json:"flightNumber"`
    Origin         string    `json:"origin"`
    Destination    string    `json:"destination"`
    DepartureTime  time.Time `json:"departureTime"`
    TotalSeats     int       `json:"totalSeats"`
    AvailableSeats int       `json:"availableSeats"`
    PriceCents     int64     `json:"priceCents"`
}

type FlightDetailResponse struct {
    FlightResponse
    SeatMap SeatMapResponse `json:"seatMap"`
}

type SeatMapResponse struct {
    Rows        int            `json:"rows"`
    SeatsPerRow int            `json:"seatsPerRow"`
    Seats       []SeatResponse `json:"seats"`
}

type SeatResponse struct {
    ID     string `json:"id"`
    Row    int    `json:"row"`
    Column string `json:"column"`
    Status string `json:"status"` // "available", "reserved", "booked"
}

type CreateOrderResponse struct {
    OrderID    string    `json:"orderId"`
    WorkflowID string    `json:"workflowId"`
    Status     string    `json:"status"`
    ExpiresAt  time.Time `json:"expiresAt"`
}

type OrderStatusResponse struct {
    OrderID         string   `json:"orderId"`
    Status          string   `json:"status"`
    Seats           []string `json:"seats"`
    TimerRemaining  int      `json:"timerRemaining"`
    PaymentAttempts int      `json:"paymentAttempts"`
    LastError       string   `json:"lastError,omitempty"`
}

type UpdateSeatsResponse struct {
    OrderID   string    `json:"orderId"`
    Status    string    `json:"status"`
    Seats     []string  `json:"seats"`
    ExpiresAt time.Time `json:"expiresAt"`
}

type PaymentAcceptedResponse struct {
    OrderID string `json:"orderId"`
    Status  string `json:"status"`
}
```

**IMPORTS:** `time`
**VALIDATE:** `go build ./internal/api/...`

---

### Task 2: CREATE `internal/api/errors.go`

Create standardized API error response utilities.

**IMPLEMENT:**
```go
package api

import (
    "encoding/json"
    "errors"
    "net/http"

    "github.com/flight-booking-system/internal/domain"
)

// ErrorResponse represents an API error
type ErrorResponse struct {
    Error   string `json:"error"`
    Message string `json:"message"`
}

// Error codes
const (
    ErrCodeInvalidRequest    = "INVALID_REQUEST"
    ErrCodeInvalidSeats      = "INVALID_SEATS"
    ErrCodeFlightNotFound    = "FLIGHT_NOT_FOUND"
    ErrCodeOrderNotFound     = "ORDER_NOT_FOUND"
    ErrCodeOrderExpired      = "ORDER_EXPIRED"
    ErrCodeSeatsUnavailable  = "SEATS_UNAVAILABLE"
    ErrCodePaymentFailed     = "PAYMENT_FAILED"
    ErrCodeInternalError     = "INTERNAL_ERROR"
    ErrCodeWorkflowError     = "WORKFLOW_ERROR"
)

// WriteError writes a JSON error response
func WriteError(w http.ResponseWriter, statusCode int, code, message string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    json.NewEncoder(w).Encode(ErrorResponse{
        Error:   code,
        Message: message,
    })
}

// WriteJSON writes a JSON response with the given status code
func WriteJSON(w http.ResponseWriter, statusCode int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    json.NewEncoder(w).Encode(data)
}

// MapDomainError maps domain errors to HTTP status codes and error codes
func MapDomainError(err error) (int, string, string) {
    switch {
    case errors.Is(err, domain.ErrFlightNotFound):
        return http.StatusNotFound, ErrCodeFlightNotFound, "Flight not found"
    case errors.Is(err, domain.ErrOrderNotFound):
        return http.StatusNotFound, ErrCodeOrderNotFound, "Order not found"
    case errors.Is(err, domain.ErrOrderExpired):
        return http.StatusConflict, ErrCodeOrderExpired, "Order reservation has expired"
    case errors.Is(err, domain.ErrSeatUnavailable), errors.Is(err, domain.ErrSeatsAlreadyLocked):
        return http.StatusConflict, ErrCodeSeatsUnavailable, "One or more seats are not available"
    case errors.Is(err, domain.ErrInvalidPaymentCode):
        return http.StatusBadRequest, ErrCodePaymentFailed, "Invalid payment code format"
    case errors.Is(err, domain.ErrPaymentFailed):
        return http.StatusBadRequest, ErrCodePaymentFailed, "Payment validation failed"
    default:
        return http.StatusInternalServerError, ErrCodeInternalError, "An internal error occurred"
    }
}

// HandleServiceError writes appropriate error response based on service error
func HandleServiceError(w http.ResponseWriter, err error) {
    statusCode, code, message := MapDomainError(err)
    WriteError(w, statusCode, code, message)
}
```

**IMPORTS:** `encoding/json`, `errors`, `net/http`, `github.com/flight-booking-system/internal/domain`
**VALIDATE:** `go build ./internal/api/...`

---

### Task 3: CREATE `internal/api/middleware.go`

Create CORS middleware for frontend origin support.

**IMPLEMENT:**
```go
package api

import "net/http"

// CORS middleware adds CORS headers for cross-origin requests
func CORS(allowedOrigins ...string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            origin := r.Header.Get("Origin")

            // Check if origin is allowed
            allowed := false
            for _, o := range allowedOrigins {
                if o == "*" || o == origin {
                    allowed = true
                    break
                }
            }

            if allowed {
                w.Header().Set("Access-Control-Allow-Origin", origin)
            } else if len(allowedOrigins) > 0 {
                w.Header().Set("Access-Control-Allow-Origin", allowedOrigins[0])
            }

            w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
            w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
            w.Header().Set("Access-Control-Max-Age", "86400")

            // Handle preflight
            if r.Method == http.MethodOptions {
                w.WriteHeader(http.StatusNoContent)
                return
            }

            next.ServeHTTP(w, r)
        })
    }
}
```

**IMPORTS:** `net/http`
**VALIDATE:** `go build ./internal/api/...`

---

### Task 4: CREATE `internal/service/temporal_client.go`

Create Temporal client wrapper for workflow interactions.

**IMPLEMENT:**
```go
package service

import (
    "context"
    "fmt"

    "go.temporal.io/sdk/client"

    "github.com/flight-booking-system/internal/config"
    temporalpkg "github.com/flight-booking-system/internal/temporal"
    "github.com/flight-booking-system/internal/temporal/workflows"
)

// TemporalClient wraps the Temporal SDK client for booking operations
type TemporalClient struct {
    client    client.Client
    taskQueue string
}

// NewTemporalClient creates a new Temporal client wrapper
func NewTemporalClient(cfg *config.TemporalConfig) (*TemporalClient, error) {
    c, err := client.Dial(client.Options{
        HostPort:  cfg.Host,
        Namespace: cfg.Namespace,
    })
    if err != nil {
        return nil, fmt.Errorf("dial temporal: %w", err)
    }

    return &TemporalClient{
        client:    c,
        taskQueue: cfg.TaskQueue,
    }, nil
}

// Close closes the Temporal client connection
func (tc *TemporalClient) Close() {
    tc.client.Close()
}

// StartBookingWorkflow starts a new booking workflow
func (tc *TemporalClient) StartBookingWorkflow(ctx context.Context, input temporalpkg.BookingWorkflowInput) (string, error) {
    workflowID := fmt.Sprintf("booking-%s", input.OrderID)

    opts := client.StartWorkflowOptions{
        ID:        workflowID,
        TaskQueue: tc.taskQueue,
    }

    run, err := tc.client.ExecuteWorkflow(ctx, opts, workflows.BookingWorkflow, input)
    if err != nil {
        return "", fmt.Errorf("start booking workflow: %w", err)
    }

    return run.GetID(), nil
}

// SignalUpdateSeats sends an update seats signal to a booking workflow
func (tc *TemporalClient) SignalUpdateSeats(ctx context.Context, orderID string, seats []string) error {
    workflowID := fmt.Sprintf("booking-%s", orderID)

    err := tc.client.SignalWorkflow(ctx, workflowID, "", temporalpkg.SignalUpdateSeats, temporalpkg.SeatUpdateSignal{
        Seats: seats,
    })
    if err != nil {
        return fmt.Errorf("signal update seats: %w", err)
    }

    return nil
}

// SignalProceedToPayment sends a proceed to payment signal with the payment code
func (tc *TemporalClient) SignalProceedToPayment(ctx context.Context, orderID string, paymentCode string) error {
    workflowID := fmt.Sprintf("booking-%s", orderID)

    err := tc.client.SignalWorkflow(ctx, workflowID, "", temporalpkg.SignalProceedToPay, temporalpkg.PaymentSignal{
        PaymentCode: paymentCode,
    })
    if err != nil {
        return fmt.Errorf("signal proceed to payment: %w", err)
    }

    return nil
}

// SignalCancelBooking sends a cancel signal to the booking workflow
func (tc *TemporalClient) SignalCancelBooking(ctx context.Context, orderID string) error {
    workflowID := fmt.Sprintf("booking-%s", orderID)

    err := tc.client.SignalWorkflow(ctx, workflowID, "", temporalpkg.SignalCancelBooking, nil)
    if err != nil {
        return fmt.Errorf("signal cancel booking: %w", err)
    }

    return nil
}

// QueryBookingStatus queries the current status of a booking workflow
func (tc *TemporalClient) QueryBookingStatus(ctx context.Context, orderID string) (*temporalpkg.BookingStatusResponse, error) {
    workflowID := fmt.Sprintf("booking-%s", orderID)

    result, err := tc.client.QueryWorkflow(ctx, workflowID, "", temporalpkg.QueryBookingStatus)
    if err != nil {
        return nil, fmt.Errorf("query booking status: %w", err)
    }

    var status temporalpkg.BookingStatusResponse
    if err := result.Get(&status); err != nil {
        return nil, fmt.Errorf("decode query result: %w", err)
    }

    return &status, nil
}
```

**IMPORTS:** `context`, `fmt`, `go.temporal.io/sdk/client`, `github.com/flight-booking-system/internal/config`, `github.com/flight-booking-system/internal/temporal`, `github.com/flight-booking-system/internal/temporal/workflows`
**GOTCHA:** Import temporal package as `temporalpkg` to avoid collision with sdk package names
**VALIDATE:** `go build ./internal/service/...`

---

### Task 5: CREATE `internal/service/flight_service.go`

Create flight service for listing and detail operations.

**IMPLEMENT:**
```go
package service

import (
    "context"

    "github.com/flight-booking-system/internal/domain"
    "github.com/flight-booking-system/internal/repository"
)

// FlightService handles flight-related business logic
type FlightService struct {
    flightRepo   *repository.FlightRepo
    seatLockRepo *repository.SeatLockRepo
}

// NewFlightService creates a new FlightService
func NewFlightService(flightRepo *repository.FlightRepo, seatLockRepo *repository.SeatLockRepo) *FlightService {
    return &FlightService{
        flightRepo:   flightRepo,
        seatLockRepo: seatLockRepo,
    }
}

// ListFlights returns all available flights
func (s *FlightService) ListFlights(ctx context.Context) ([]domain.Flight, error) {
    return s.flightRepo.FindAll(ctx)
}

// GetFlightWithSeats returns a flight with its seat map and real-time availability
func (s *FlightService) GetFlightWithSeats(ctx context.Context, flightID string) (*domain.FlightWithSeats, error) {
    // Get flight details
    flight, err := s.flightRepo.FindByID(ctx, flightID)
    if err != nil {
        return nil, err
    }

    // Get all seats for the flight
    seats, err := s.flightRepo.FindSeats(ctx, flightID)
    if err != nil {
        return nil, err
    }

    // Get currently locked seats from Redis
    lockedSeats, err := s.seatLockRepo.GetLockedSeats(ctx, flightID)
    if err != nil {
        return nil, err
    }

    // Update seat status based on locks
    for i := range seats {
        if _, isLocked := lockedSeats[seats[i].ID]; isLocked {
            if seats[i].Status == domain.SeatStatusAvailable {
                seats[i].Status = domain.SeatStatusReserved
            }
        }
    }

    // Calculate seat map dimensions
    rows := 0
    seatsPerRow := 0
    if len(seats) > 0 {
        rowMap := make(map[int]int)
        for _, seat := range seats {
            rowMap[seat.Row]++
            if seat.Row > rows {
                rows = seat.Row
            }
        }
        // Get seats per row from first row
        if count, ok := rowMap[1]; ok {
            seatsPerRow = count
        }
    }

    return &domain.FlightWithSeats{
        Flight: *flight,
        SeatMap: domain.SeatMap{
            Rows:        rows,
            SeatsPerRow: seatsPerRow,
            Seats:       seats,
        },
    }, nil
}
```

**IMPORTS:** `context`, `github.com/flight-booking-system/internal/domain`, `github.com/flight-booking-system/internal/repository`
**PATTERN:** Repository injection, context propagation
**VALIDATE:** `go build ./internal/service/...`

---

### Task 6: CREATE `internal/service/booking_service.go`

Create booking service for order management and Temporal workflow interaction.

**IMPLEMENT:**
```go
package service

import (
    "context"
    "fmt"
    "regexp"
    "time"

    "github.com/google/uuid"

    "github.com/flight-booking-system/internal/domain"
    "github.com/flight-booking-system/internal/repository"
)

// BookingService handles booking-related business logic
type BookingService struct {
    orderRepo      *repository.OrderRepo
    flightRepo     *repository.FlightRepo
    temporalClient *TemporalClient
}

// NewBookingService creates a new BookingService
func NewBookingService(
    orderRepo *repository.OrderRepo,
    flightRepo *repository.FlightRepo,
    temporalClient *TemporalClient,
) *BookingService {
    return &BookingService{
        orderRepo:      orderRepo,
        flightRepo:     flightRepo,
        temporalClient: temporalClient,
    }
}

// CreateOrderInput contains the parameters for creating an order
type CreateOrderInput struct {
    FlightID string
    Seats    []string
}

// CreateOrderOutput contains the result of order creation
type CreateOrderOutput struct {
    OrderID    string
    WorkflowID string
    Status     domain.OrderStatus
    ExpiresAt  time.Time
}

// CreateOrder creates a new booking order and starts the workflow
func (s *BookingService) CreateOrder(ctx context.Context, input CreateOrderInput) (*CreateOrderOutput, error) {
    // Validate flight exists
    flight, err := s.flightRepo.FindByID(ctx, input.FlightID)
    if err != nil {
        return nil, err
    }

    // Validate seats are not empty
    if len(input.Seats) == 0 {
        return nil, domain.ErrSeatUnavailable
    }

    // Generate order ID
    orderID := uuid.New().String()

    // Calculate expiration (15 minutes from now)
    expiresAt := time.Now().Add(15 * time.Minute)

    // Start the booking workflow
    workflowInput := struct {
        OrderID  string   `json:"orderId"`
        FlightID string   `json:"flightId"`
        Seats    []string `json:"seats"`
    }{
        OrderID:  orderID,
        FlightID: input.FlightID,
        Seats:    input.Seats,
    }

    // Import temporal package types
    temporalInput := toBookingWorkflowInput(orderID, input.FlightID, input.Seats)
    workflowID, err := s.temporalClient.StartBookingWorkflow(ctx, temporalInput)
    if err != nil {
        return nil, fmt.Errorf("start workflow: %w", err)
    }

    // Note: Order is created by the workflow's CreateOrder activity
    // We return optimistically assuming the workflow will create it
    _ = workflowInput // suppress unused warning
    _ = flight        // suppress unused warning

    return &CreateOrderOutput{
        OrderID:    orderID,
        WorkflowID: workflowID,
        Status:     domain.OrderStatusSeatsReserved,
        ExpiresAt:  expiresAt,
    }, nil
}

// GetOrderStatus queries the workflow for current order status
func (s *BookingService) GetOrderStatus(ctx context.Context, orderID string) (*domain.OrderStatusResponse, error) {
    // First try to query the workflow
    status, err := s.temporalClient.QueryBookingStatus(ctx, orderID)
    if err != nil {
        // If workflow query fails, try to get from database
        order, dbErr := s.orderRepo.FindByID(ctx, orderID)
        if dbErr != nil {
            return nil, domain.ErrOrderNotFound
        }

        // Return status from database (for completed/failed/expired orders)
        timerRemaining := 0
        if order.ExpiresAt != nil {
            remaining := time.Until(*order.ExpiresAt)
            if remaining > 0 {
                timerRemaining = int(remaining.Seconds())
            }
        }

        return &domain.OrderStatusResponse{
            OrderID:         order.ID,
            Status:          order.Status,
            Seats:           order.Seats,
            TimerRemaining:  timerRemaining,
            PaymentAttempts: 0,
            LastError:       stringValue(order.FailureReason),
        }, nil
    }

    return &domain.OrderStatusResponse{
        OrderID:         status.OrderID,
        Status:          status.Status,
        Seats:           status.Seats,
        TimerRemaining:  status.TimerRemaining,
        PaymentAttempts: status.PaymentAttempts,
        LastError:       status.LastError,
    }, nil
}

// UpdateSeats updates the seat selection for an order
func (s *BookingService) UpdateSeats(ctx context.Context, orderID string, seats []string) (*UpdateSeatsOutput, error) {
    // Validate seats are not empty
    if len(seats) == 0 {
        return nil, domain.ErrSeatUnavailable
    }

    // Send signal to workflow
    err := s.temporalClient.SignalUpdateSeats(ctx, orderID, seats)
    if err != nil {
        return nil, fmt.Errorf("signal update seats: %w", err)
    }

    // Query updated status
    status, err := s.temporalClient.QueryBookingStatus(ctx, orderID)
    if err != nil {
        return nil, fmt.Errorf("query status: %w", err)
    }

    return &UpdateSeatsOutput{
        OrderID:   status.OrderID,
        Status:    status.Status,
        Seats:     status.Seats,
        ExpiresAt: status.ExpiresAt,
    }, nil
}

// UpdateSeatsOutput contains the result of seat update
type UpdateSeatsOutput struct {
    OrderID   string
    Status    domain.OrderStatus
    Seats     []string
    ExpiresAt time.Time
}

// SubmitPayment submits a payment for an order
func (s *BookingService) SubmitPayment(ctx context.Context, orderID string, paymentCode string) error {
    // Validate payment code format (5 digits)
    if !isValidPaymentCode(paymentCode) {
        return domain.ErrInvalidPaymentCode
    }

    // Send payment signal to workflow
    err := s.temporalClient.SignalProceedToPayment(ctx, orderID, paymentCode)
    if err != nil {
        return fmt.Errorf("signal payment: %w", err)
    }

    return nil
}

// CancelOrder cancels an order
func (s *BookingService) CancelOrder(ctx context.Context, orderID string) error {
    err := s.temporalClient.SignalCancelBooking(ctx, orderID)
    if err != nil {
        return fmt.Errorf("signal cancel: %w", err)
    }

    return nil
}

// Helper functions

func isValidPaymentCode(code string) bool {
    matched, _ := regexp.MatchString(`^\d{5}$`, code)
    return matched
}

func stringValue(s *string) string {
    if s == nil {
        return ""
    }
    return *s
}

// toBookingWorkflowInput converts to the temporal package input type
func toBookingWorkflowInput(orderID, flightID string, seats []string) struct {
    OrderID  string   `json:"orderId"`
    FlightID string   `json:"flightId"`
    Seats    []string `json:"seats"`
} {
    return struct {
        OrderID  string   `json:"orderId"`
        FlightID string   `json:"flightId"`
        Seats    []string `json:"seats"`
    }{
        OrderID:  orderID,
        FlightID: flightID,
        Seats:    seats,
    }
}
```

**IMPORTS:** `context`, `fmt`, `regexp`, `time`, `github.com/google/uuid`, `github.com/flight-booking-system/internal/domain`, `github.com/flight-booking-system/internal/repository`
**GOTCHA:** Need to add `github.com/google/uuid` dependency: `go get github.com/google/uuid`
**VALIDATE:** `go get github.com/google/uuid && go build ./internal/service/...`

---

### Task 7: UPDATE `internal/service/booking_service.go` - Fix temporal import

Update the booking service to properly use the temporal package types.

**IMPLEMENT:** Replace the `toBookingWorkflowInput` function and update imports:

Find and replace the import section and helper function:

```go
// At top of file, add import:
import (
    // ... existing imports ...
    temporalpkg "github.com/flight-booking-system/internal/temporal"
)

// Replace toBookingWorkflowInput function with:
func toTemporalInput(orderID, flightID string, seats []string) temporalpkg.BookingWorkflowInput {
    return temporalpkg.BookingWorkflowInput{
        OrderID:  orderID,
        FlightID: flightID,
        Seats:    seats,
    }
}
```

Then update CreateOrder to use:
```go
temporalInput := toTemporalInput(orderID, input.FlightID, input.Seats)
```

**VALIDATE:** `go build ./internal/service/...`

---

### Task 8: CREATE `internal/api/handlers.go`

Create HTTP handlers for all API endpoints.

**IMPLEMENT:**
```go
package api

import (
    "encoding/json"
    "net/http"

    "github.com/go-chi/chi/v5"

    "github.com/flight-booking-system/internal/service"
)

// Handlers contains all HTTP handlers
type Handlers struct {
    flightService  *service.FlightService
    bookingService *service.BookingService
}

// NewHandlers creates a new Handlers instance
func NewHandlers(flightService *service.FlightService, bookingService *service.BookingService) *Handlers {
    return &Handlers{
        flightService:  flightService,
        bookingService: bookingService,
    }
}

// ListFlights handles GET /api/flights
func (h *Handlers) ListFlights(w http.ResponseWriter, r *http.Request) {
    flights, err := h.flightService.ListFlights(r.Context())
    if err != nil {
        HandleServiceError(w, err)
        return
    }

    response := FlightListResponse{
        Flights: make([]FlightResponse, len(flights)),
    }
    for i, f := range flights {
        response.Flights[i] = FlightResponse{
            ID:             f.ID,
            FlightNumber:   f.FlightNumber,
            Origin:         f.Origin,
            Destination:    f.Destination,
            DepartureTime:  f.DepartureTime,
            TotalSeats:     f.TotalSeats,
            AvailableSeats: f.AvailableSeats,
            PriceCents:     f.PriceCents,
        }
    }

    WriteJSON(w, http.StatusOK, response)
}

// GetFlight handles GET /api/flights/{flightId}
func (h *Handlers) GetFlight(w http.ResponseWriter, r *http.Request) {
    flightID := chi.URLParam(r, "flightId")
    if flightID == "" {
        WriteError(w, http.StatusBadRequest, ErrCodeInvalidRequest, "flight ID is required")
        return
    }

    flight, err := h.flightService.GetFlightWithSeats(r.Context(), flightID)
    if err != nil {
        HandleServiceError(w, err)
        return
    }

    // Build seat response
    seats := make([]SeatResponse, len(flight.SeatMap.Seats))
    for i, s := range flight.SeatMap.Seats {
        seats[i] = SeatResponse{
            ID:     s.ID,
            Row:    s.Row,
            Column: s.Column,
            Status: string(s.Status),
        }
    }

    response := FlightDetailResponse{
        FlightResponse: FlightResponse{
            ID:             flight.ID,
            FlightNumber:   flight.FlightNumber,
            Origin:         flight.Origin,
            Destination:    flight.Destination,
            DepartureTime:  flight.DepartureTime,
            TotalSeats:     flight.TotalSeats,
            AvailableSeats: flight.AvailableSeats,
            PriceCents:     flight.PriceCents,
        },
        SeatMap: SeatMapResponse{
            Rows:        flight.SeatMap.Rows,
            SeatsPerRow: flight.SeatMap.SeatsPerRow,
            Seats:       seats,
        },
    }

    WriteJSON(w, http.StatusOK, response)
}

// CreateOrder handles POST /api/orders
func (h *Handlers) CreateOrder(w http.ResponseWriter, r *http.Request) {
    var req CreateOrderRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        WriteError(w, http.StatusBadRequest, ErrCodeInvalidRequest, "invalid request body")
        return
    }

    // Validate request
    if req.FlightID == "" {
        WriteError(w, http.StatusBadRequest, ErrCodeInvalidRequest, "flightId is required")
        return
    }
    if len(req.Seats) == 0 {
        WriteError(w, http.StatusBadRequest, ErrCodeInvalidSeats, "at least one seat must be selected")
        return
    }

    output, err := h.bookingService.CreateOrder(r.Context(), service.CreateOrderInput{
        FlightID: req.FlightID,
        Seats:    req.Seats,
    })
    if err != nil {
        HandleServiceError(w, err)
        return
    }

    response := CreateOrderResponse{
        OrderID:    output.OrderID,
        WorkflowID: output.WorkflowID,
        Status:     string(output.Status),
        ExpiresAt:  output.ExpiresAt,
    }

    WriteJSON(w, http.StatusCreated, response)
}

// UpdateSeats handles PUT /api/orders/{orderId}/seats
func (h *Handlers) UpdateSeats(w http.ResponseWriter, r *http.Request) {
    orderID := chi.URLParam(r, "orderId")
    if orderID == "" {
        WriteError(w, http.StatusBadRequest, ErrCodeInvalidRequest, "order ID is required")
        return
    }

    var req UpdateSeatsRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        WriteError(w, http.StatusBadRequest, ErrCodeInvalidRequest, "invalid request body")
        return
    }

    if len(req.Seats) == 0 {
        WriteError(w, http.StatusBadRequest, ErrCodeInvalidSeats, "at least one seat must be selected")
        return
    }

    output, err := h.bookingService.UpdateSeats(r.Context(), orderID, req.Seats)
    if err != nil {
        HandleServiceError(w, err)
        return
    }

    response := UpdateSeatsResponse{
        OrderID:   output.OrderID,
        Status:    string(output.Status),
        Seats:     output.Seats,
        ExpiresAt: output.ExpiresAt,
    }

    WriteJSON(w, http.StatusOK, response)
}

// GetOrderStatus handles GET /api/orders/{orderId}/status
func (h *Handlers) GetOrderStatus(w http.ResponseWriter, r *http.Request) {
    orderID := chi.URLParam(r, "orderId")
    if orderID == "" {
        WriteError(w, http.StatusBadRequest, ErrCodeInvalidRequest, "order ID is required")
        return
    }

    status, err := h.bookingService.GetOrderStatus(r.Context(), orderID)
    if err != nil {
        HandleServiceError(w, err)
        return
    }

    response := OrderStatusResponse{
        OrderID:         status.OrderID,
        Status:          string(status.Status),
        Seats:           status.Seats,
        TimerRemaining:  status.TimerRemaining,
        PaymentAttempts: status.PaymentAttempts,
        LastError:       status.LastError,
    }

    WriteJSON(w, http.StatusOK, response)
}

// SubmitPayment handles POST /api/orders/{orderId}/pay
func (h *Handlers) SubmitPayment(w http.ResponseWriter, r *http.Request) {
    orderID := chi.URLParam(r, "orderId")
    if orderID == "" {
        WriteError(w, http.StatusBadRequest, ErrCodeInvalidRequest, "order ID is required")
        return
    }

    var req SubmitPaymentRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        WriteError(w, http.StatusBadRequest, ErrCodeInvalidRequest, "invalid request body")
        return
    }

    if req.PaymentCode == "" {
        WriteError(w, http.StatusBadRequest, ErrCodePaymentFailed, "payment code is required")
        return
    }

    err := h.bookingService.SubmitPayment(r.Context(), orderID, req.PaymentCode)
    if err != nil {
        HandleServiceError(w, err)
        return
    }

    response := PaymentAcceptedResponse{
        OrderID: orderID,
        Status:  "PAYMENT_PROCESSING",
    }

    WriteJSON(w, http.StatusAccepted, response)
}

// CancelOrder handles DELETE /api/orders/{orderId}
func (h *Handlers) CancelOrder(w http.ResponseWriter, r *http.Request) {
    orderID := chi.URLParam(r, "orderId")
    if orderID == "" {
        WriteError(w, http.StatusBadRequest, ErrCodeInvalidRequest, "order ID is required")
        return
    }

    err := h.bookingService.CancelOrder(r.Context(), orderID)
    if err != nil {
        HandleServiceError(w, err)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}
```

**IMPORTS:** `encoding/json`, `net/http`, `github.com/go-chi/chi/v5`, `github.com/flight-booking-system/internal/service`
**PATTERN:** Thin handlers - parse, validate, call service, respond
**VALIDATE:** `go build ./internal/api/...`

---

### Task 9: CREATE `internal/api/routes.go`

Create router configuration with all routes.

**IMPLEMENT:**
```go
package api

import (
    "net/http"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"

    "github.com/flight-booking-system/internal/database"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/redis/go-redis/v9"
)

// RouterConfig holds dependencies for router creation
type RouterConfig struct {
    Pool        *pgxpool.Pool
    RedisClient *redis.Client
    Handlers    *Handlers
}

// NewRouter creates a new Chi router with all routes configured
func NewRouter(cfg RouterConfig) *chi.Mux {
    r := chi.NewRouter()

    // Global middleware
    r.Use(middleware.RequestID)
    r.Use(middleware.RealIP)
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)
    r.Use(CORS("http://localhost:3000", "http://localhost:5173"))

    // Health check
    r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
        // Check database
        if err := database.HealthCheck(r.Context(), cfg.Pool); err != nil {
            http.Error(w, "database unhealthy", http.StatusServiceUnavailable)
            return
        }

        // Check Redis
        if err := database.RedisHealthCheck(r.Context(), cfg.RedisClient); err != nil {
            http.Error(w, "redis unhealthy", http.StatusServiceUnavailable)
            return
        }

        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    })

    // API routes
    r.Route("/api", func(r chi.Router) {
        // Flight routes
        r.Route("/flights", func(r chi.Router) {
            r.Get("/", cfg.Handlers.ListFlights)
            r.Get("/{flightId}", cfg.Handlers.GetFlight)
        })

        // Order routes
        r.Route("/orders", func(r chi.Router) {
            r.Post("/", cfg.Handlers.CreateOrder)

            r.Route("/{orderId}", func(r chi.Router) {
                r.Put("/seats", cfg.Handlers.UpdateSeats)
                r.Get("/status", cfg.Handlers.GetOrderStatus)
                r.Post("/pay", cfg.Handlers.SubmitPayment)
                r.Delete("/", cfg.Handlers.CancelOrder)
            })
        })
    })

    return r
}
```

**IMPORTS:** `net/http`, `github.com/go-chi/chi/v5`, `github.com/go-chi/chi/v5/middleware`, `github.com/flight-booking-system/internal/database`, `github.com/jackc/pgx/v5/pgxpool`, `github.com/redis/go-redis/v9`
**PATTERN:** Route grouping with `r.Route()`, URL params with `{paramName}`
**VALIDATE:** `go build ./internal/api/...`

---

### Task 10: UPDATE `cmd/server/main.go`

Replace inline handlers with new API package and add Temporal client.

**IMPLEMENT:** Full replacement of `cmd/server/main.go`:

```go
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/flight-booking-system/internal/api"
    "github.com/flight-booking-system/internal/config"
    "github.com/flight-booking-system/internal/database"
    "github.com/flight-booking-system/internal/repository"
    "github.com/flight-booking-system/internal/service"
)

func main() {
    // Load configuration
    cfg := config.Load()

    // Create context for graceful shutdown
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // Connect to PostgreSQL
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
    temporalClient, err := service.NewTemporalClient(&cfg.Temporal)
    if err != nil {
        log.Fatalf("Failed to connect to Temporal: %v", err)
    }
    defer temporalClient.Close()
    log.Println("Connected to Temporal")

    // Create repositories
    flightRepo := repository.NewFlightRepo(pool)
    orderRepo := repository.NewOrderRepo(pool)
    seatLockRepo := repository.NewSeatLockRepo(redisClient)

    // Create services
    flightService := service.NewFlightService(flightRepo, seatLockRepo)
    bookingService := service.NewBookingService(orderRepo, flightRepo, temporalClient)

    // Create handlers
    handlers := api.NewHandlers(flightService, bookingService)

    // Create router
    router := api.NewRouter(api.RouterConfig{
        Pool:        pool,
        RedisClient: redisClient,
        Handlers:    handlers,
    })

    // Create server
    addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
    srv := &http.Server{
        Addr:         addr,
        Handler:      router,
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 15 * time.Second,
        IdleTimeout:  60 * time.Second,
    }

    // Start server in goroutine
    go func() {
        log.Printf("Server starting on %s", addr)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Server failed: %v", err)
        }
    }()

    // Wait for interrupt signal
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    log.Println("Shutting down server...")

    // Graceful shutdown
    shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer shutdownCancel()

    if err := srv.Shutdown(shutdownCtx); err != nil {
        log.Fatalf("Server forced to shutdown: %v", err)
    }

    log.Println("Server stopped")
}
```

**IMPORTS:** Updated to use new packages
**VALIDATE:** `go build ./cmd/server/...`

---

### Task 11: ADD uuid dependency

Add the uuid package to go.mod.

**IMPLEMENT:**
```bash
go get github.com/google/uuid
```

**VALIDATE:** `go mod tidy && go build ./...`

---

### Task 12: BUILD and verify all packages

Full build verification.

**IMPLEMENT:**
```bash
go build ./...
```

**VALIDATE:** Build succeeds with no errors

---

## TESTING STRATEGY

### Unit Tests

Create tests in the following files:

**`internal/api/handlers_test.go`:**
- Test request parsing and validation
- Test error response mapping
- Mock service dependencies

**`internal/service/flight_service_test.go`:**
- Test flight listing
- Test seat availability calculation

**`internal/service/booking_service_test.go`:**
- Test order creation validation
- Test payment code validation

### Integration Tests

**`internal/api/integration_test.go`:**
- Test full request/response cycle
- Requires test database and Redis

### Edge Cases

- Empty seat selection
- Invalid flight ID
- Invalid order ID
- Invalid payment code (not 5 digits)
- Expired order queries
- Workflow query failures (fallback to database)

---

## VALIDATION COMMANDS

Execute every command to ensure zero regressions and 100% feature correctness.

### Level 1: Syntax & Build

```bash
# Build all packages
go build ./...

# Run go vet
go vet ./...

# Check formatting
gofmt -d .
```

### Level 2: Unit Tests

```bash
# Run all tests
go test -v ./...

# Run with race detector
go test -race ./...
```

### Level 3: Linting

```bash
# Run linter (if golangci-lint installed)
golangci-lint run ./...
```

### Level 4: Manual Validation

Start the infrastructure:
```bash
make up
make migrate-up
```

Start the worker (in terminal 1):
```bash
go run ./cmd/worker
```

Start the server (in terminal 2):
```bash
go run ./cmd/server
```

Test endpoints:
```bash
# List flights
curl http://localhost:8080/api/flights | jq

# Get flight details (use ID from previous response)
curl http://localhost:8080/api/flights/{flightId} | jq

# Create order
curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -d '{"flightId": "{flightId}", "seats": ["1A", "1B"]}' | jq

# Get order status
curl http://localhost:8080/api/orders/{orderId}/status | jq

# Update seats
curl -X PUT http://localhost:8080/api/orders/{orderId}/seats \
  -H "Content-Type: application/json" \
  -d '{"seats": ["2A", "2B"]}' | jq

# Submit payment
curl -X POST http://localhost:8080/api/orders/{orderId}/pay \
  -H "Content-Type: application/json" \
  -d '{"paymentCode": "12345"}' | jq

# Check status after payment
curl http://localhost:8080/api/orders/{orderId}/status | jq
```

---

## ACCEPTANCE CRITERIA

- [ ] All 6 API endpoints implemented and functional:
  - [ ] `GET /api/flights` - List all flights
  - [ ] `GET /api/flights/{flightId}` - Get flight with seat map
  - [ ] `POST /api/orders` - Create booking order
  - [ ] `PUT /api/orders/{orderId}/seats` - Update seat selection
  - [ ] `GET /api/orders/{orderId}/status` - Get order status
  - [ ] `POST /api/orders/{orderId}/pay` - Submit payment
- [ ] Services layer encapsulates business logic
- [ ] Temporal client properly starts workflows, sends signals, and queries status
- [ ] Error responses follow standardized format with error codes
- [ ] All validation commands pass with zero errors
- [ ] Code follows project conventions (CLAUDE.md guidelines)
- [ ] No regressions in existing functionality

---

## COMPLETION CHECKLIST

- [ ] All tasks completed in order
- [ ] Each task validation passed immediately
- [ ] All validation commands executed successfully
- [ ] Full build passes (`go build ./...`)
- [ ] No linting errors
- [ ] Manual testing confirms API works end-to-end
- [ ] Acceptance criteria all met

---

## NOTES

### Design Decisions

1. **Services Layer**: Chose to create a separate services layer rather than having handlers call repositories directly. This provides:
   - Clean separation of HTTP concerns from business logic
   - Single point for Temporal client interaction
   - Easier unit testing through service mocking

2. **Temporal Client Wrapper**: Created `TemporalClient` struct to encapsulate workflow interaction patterns:
   - Consistent workflow ID generation (`booking-{orderID}`)
   - Type-safe signal/query methods
   - Centralized connection management

3. **Error Mapping**: Implemented `MapDomainError` to translate domain errors to HTTP status codes, keeping error handling consistent across all endpoints.

4. **Optimistic Response on Order Creation**: The order creation returns immediately after starting the workflow, not waiting for the workflow to create the order in the database. This is acceptable because:
   - The workflow will create the order as its first activity
   - If the workflow fails, the order won't exist when queried
   - This keeps the API responsive

### Trade-offs

1. **Query Fallback**: When workflow query fails, we fall back to database lookup. This handles cases where workflows have completed but clients are still polling.

2. **No Request Timeout Middleware per Route**: Using global 60s timeout. For more fine-grained control, could add per-route timeout middleware.

### Future Improvements

- Add request validation middleware using a validation library
- Add rate limiting middleware
- Add structured logging with request correlation
- Add OpenAPI documentation generation
- Add metrics/tracing for observability
