# Flight Booking System - Product Requirements Document

## 1. Executive Summary

This document outlines the requirements for a flight booking system built on Temporal workflow orchestration. The system enables users to search flights, reserve seats with time-limited holds, and complete payments through a validated code system.

The architecture leverages Temporal for managing complex booking workflows including seat reservation timeouts, payment retries, and order state management. The backend is implemented in Go with a RESTful API layer, while the frontend is a minimal React application providing real-time booking status and countdown timers.

**MVP Goal:** Deliver a functional flight booking system demonstrating Temporal workflow patterns for seat reservation (15-minute holds with refresh), payment validation (10-second timeout, 3 retries, 15% simulated failure), and real-time order tracking.

## 2. Mission

**Mission Statement:** Build a robust, fault-tolerant flight booking system that showcases Temporal's workflow orchestration capabilities for handling complex business processes with timeouts, retries, and state management.

**Core Principles:**
1. **Reliability First** - Workflows must handle failures gracefully with proper compensation logic
2. **Real-time Visibility** - Users always see current seat hold timers and order status
3. **Temporal-Native** - Leverage Temporal patterns (signals, queries, timers) rather than custom implementations
4. **Simplicity** - Minimal viable frontend; complexity lives in the workflow layer
5. **Demonstrable** - Clear simulation of failure scenarios (15% payment failures) for learning/demo purposes

## 3. Target Users

**Primary Persona: Flight Booking Customer**
- Books flights through web interface
- Selects specific seats from available inventory
- Completes payment using 5-digit validation codes
- Needs real-time feedback on reservation timer and payment status

**Secondary Persona: Developer/Learner**
- Studies Temporal workflow patterns
- Examines seat reservation timeout/refresh logic
- Observes payment retry and failure handling
- Reviews distributed locking patterns for inventory management

**Key User Needs:**
- Clear visual indication of seat hold timer (countdown)
- Immediate feedback when seat selection changes (timer refresh)
- Understanding of payment retry status during validation
- Clear error messages when booking fails

## 4. MVP Scope

### In Scope

**Core Functionality:**
- ✅ Flight listing with available seats
- ✅ Seat selection with visual seat map
- ✅ 15-minute seat reservation hold
- ✅ Timer refresh on seat selection changes
- ✅ Auto-release of seats when timer expires
- ✅ 5-digit payment code entry and validation
- ✅ 10-second payment validation timeout
- ✅ 3 retry attempts for failed payments
- ✅ 15% simulated payment failure rate
- ✅ Order confirmation on successful payment
- ✅ Real-time status updates via polling/SSE

**Technical:**
- ✅ Go RESTful API server
- ✅ Temporal workflows for booking orchestration
- ✅ Temporal workers for activity execution
- ✅ PostgreSQL for flight/order persistence
- ✅ Redis distributed locks for seat inventory
- ✅ Basic React frontend with timer display

**Deployment:**
- ✅ Local development with `temporal server start-dev`
- ✅ Docker Compose for full stack

### Out of Scope

**Deferred Features:**
- ❌ User authentication/accounts
- ❌ Flight search/filtering
- ❌ Multiple passengers per booking
- ❌ Actual payment gateway integration
- ❌ Email/SMS notifications
- ❌ Booking history
- ❌ Seat pricing tiers
- ❌ Flight cancellation/refunds
- ❌ Admin panel
- ❌ Production deployment (Kubernetes, Temporal Cloud)

## 5. User Stories

### Primary User Stories

**US-1: View Available Flights**
> As a customer, I want to see available flights with their seat availability, so that I can choose a flight to book.

*Example: User opens app and sees "Flight FL-101: NYC → LAX, 45/120 seats available"*

**US-2: Select Seats**
> As a customer, I want to select specific seats from a visual seat map, so that I can choose my preferred seating.

*Example: User clicks on seat 12A, seat turns from green (available) to blue (selected), 15-minute timer starts*

**US-3: See Reservation Timer**
> As a customer, I want to see a countdown timer for my seat hold, so that I know how much time I have to complete booking.

*Example: Timer displays "14:32 remaining" and counts down in real-time*

**US-4: Modify Seat Selection**
> As a customer, I want to change my seat selection and have the timer refresh, so that I'm not penalized for reconsidering.

*Example: User deselects 12A and selects 14B, timer resets to 15:00*

**US-5: Enter Payment Code**
> As a customer, I want to enter a 5-digit payment code to complete my booking, so that I can pay for my reservation.

*Example: User enters "12345", sees "Validating payment..." with spinner*

**US-6: See Payment Retry Status**
> As a customer, I want to see if payment validation is being retried, so that I understand why it's taking longer.

*Example: "Payment validation failed. Retrying (attempt 2 of 3)..."*

**US-7: Receive Booking Confirmation**
> As a customer, I want to receive confirmation when my booking succeeds, so that I know my seats are secured.

*Example: "Booking confirmed! Order #ORD-789. Seats 12A, 12B on Flight FL-101"*

**US-8: Understand Booking Failure**
> As a customer, I want to see a clear message when booking fails, so that I can take appropriate action.

*Example: "Payment validation failed after 3 attempts. Your seats have been released. Please try again."*

## 6. Core Architecture & Patterns

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         User Layer                               │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │                    React Web App                         │    │
│  │         (Real-time Status, Timer, Seat Map)             │    │
│  └─────────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Application Layer                           │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │                 RESTful Server (Go)                      │    │
│  │     /flights, /orders, /orders/{id}/seats, /pay         │    │
│  └─────────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Temporal Platform                           │
│  ┌──────────────────────┐    ┌──────────────────────────────┐   │
│  │   Temporal Server    │    │      Temporal Workers (Go)   │   │
│  │  (Workflow Engine)   │◄──►│  - SeatReservationActivity   │   │
│  │                      │    │  - PaymentValidationActivity │   │
│  └──────────────────────┘    │  - OrderManagementActivity   │   │
│                              └──────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                       Data Layer                                 │
│  ┌───────────────────┐         ┌────────────────────────────┐   │
│  │    PostgreSQL     │         │          Redis             │   │
│  │  (Flights, Orders)│         │  (Seat Locks with TTL)     │   │
│  └───────────────────┘         └────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
```

### Directory Structure

```
flight-booking-system/
├── cmd/
│   ├── server/          # REST API server entrypoint
│   │   └── main.go
│   └── worker/          # Temporal worker entrypoint
│       └── main.go
├── internal/
│   ├── api/             # HTTP handlers and routes
│   │   ├── handlers.go
│   │   ├── routes.go
│   │   └── middleware.go
│   ├── domain/          # Domain models
│   │   ├── flight.go
│   │   ├── order.go
│   │   └── seat.go
│   ├── repository/      # Data access layer
│   │   ├── flight_repo.go
│   │   ├── order_repo.go
│   │   └── seat_lock.go
│   ├── temporal/        # Temporal workflows and activities
│   │   ├── workflows/
│   │   │   └── booking_workflow.go
│   │   └── activities/
│   │       ├── seat_reservation.go
│   │       ├── payment_validation.go
│   │       └── order_management.go
│   └── service/         # Business logic services
│       ├── flight_service.go
│       ├── booking_service.go
│       └── payment_service.go
├── web/                 # React frontend
│   ├── src/
│   │   ├── components/
│   │   │   ├── FlightList.jsx
│   │   │   ├── SeatMap.jsx
│   │   │   ├── Timer.jsx
│   │   │   ├── PaymentForm.jsx
│   │   │   └── OrderStatus.jsx
│   │   ├── hooks/
│   │   │   └── useOrderStatus.js
│   │   ├── api/
│   │   │   └── client.js
│   │   ├── App.jsx
│   │   └── main.jsx
│   ├── package.json
│   └── vite.config.js
├── migrations/          # Database migrations
│   └── 001_initial.sql
├── docker-compose.yml
├── Makefile
├── go.mod
├── go.sum
└── README.md
```

### Key Design Patterns

**Temporal Workflow Patterns:**
- **Long-running workflow** - BookingWorkflow runs for duration of booking process (up to 15+ minutes)
- **Signals** - `UpdateSeats` signal to refresh seat selection and reset timer
- **Queries** - `GetStatus` query for real-time order state without side effects
- **Timers** - `workflow.NewTimer(15 * time.Minute)` for seat hold expiration
- **Activity retries** - Payment validation with `RetryPolicy{MaximumAttempts: 3}`
- **Saga pattern** - Compensating actions to release seats on payment failure

**Data Patterns:**
- **Distributed locking** - Redis `SET seat:flight:seatId NX EX 900` for seat holds
- **TTL-based expiration** - Redis keys auto-expire, Temporal timer as backup
- **Optimistic concurrency** - Version field on orders for conflict detection

## 7. Features

### Feature 1: Flight Listing

**Purpose:** Display available flights with seat availability counts

**Operations:**
- `GET /api/flights` - List all flights with availability
- `GET /api/flights/{id}` - Get flight details with seat map

**Key Features:**
- Shows flight number, route, departure time
- Real-time available seat count
- Seat map visualization (rows/columns with status)

### Feature 2: Seat Reservation Workflow

**Purpose:** Manage seat holds with 15-minute timeout and refresh capability

**Workflow: BookingWorkflow**
```go
// Simplified workflow logic
func BookingWorkflow(ctx workflow.Context, orderID string) error {
    // Initial seat reservation
    reserveSeats(ctx, orderID)

    // Timer with signal-based refresh
    for {
        selector := workflow.NewSelector(ctx)

        // 15-minute timeout
        timer := workflow.NewTimer(ctx, 15*time.Minute)
        selector.AddFuture(timer, func(f workflow.Future) {
            // Timer expired - release seats
            releaseSeats(ctx, orderID)
            return ErrReservationExpired
        })

        // Signal to update seats (resets timer)
        selector.AddReceive(updateSeatsChannel, func(c workflow.ReceiveChannel, more bool) {
            var newSeats []string
            c.Receive(ctx, &newSeats)
            updateSeatSelection(ctx, orderID, newSeats)
            // Loop continues, timer resets
        })

        // Signal to proceed to payment
        selector.AddReceive(proceedToPaymentChannel, func(...) {
            break // Exit loop, move to payment
        })

        selector.Select(ctx)
    }

    // Payment phase
    return processPayment(ctx, orderID)
}
```

**Key Features:**
- Redis distributed lock acquired on seat selection
- Lock TTL matches Temporal timer (15 minutes)
- Seat change signal refreshes both Temporal timer and Redis TTL
- Auto-release on timeout via both mechanisms (belt and suspenders)

### Feature 3: Payment Validation

**Purpose:** Validate 5-digit payment codes with timeout and retry logic

**Activity: ValidatePayment**
```go
func (a *Activities) ValidatePayment(ctx context.Context, code string) error {
    // Simulate 15% failure rate
    if rand.Float32() < 0.15 {
        return temporal.NewApplicationError("Payment validation failed", "PAYMENT_FAILED")
    }

    // Simulate processing time (up to 10 seconds)
    time.Sleep(time.Duration(rand.Intn(8)+1) * time.Second)

    return nil
}
```

**Retry Policy:**
```go
retryPolicy := &temporal.RetryPolicy{
    InitialInterval:    time.Second,
    BackoffCoefficient: 1.5,
    MaximumInterval:    5 * time.Second,
    MaximumAttempts:    3,
}

activityOptions := workflow.ActivityOptions{
    StartToCloseTimeout: 10 * time.Second,
    RetryPolicy:         retryPolicy,
}
```

**Key Features:**
- 10-second timeout per validation attempt
- 3 maximum retry attempts
- 15% simulated failure rate for demo purposes
- Exponential backoff between retries

### Feature 4: Order Management

**Purpose:** Track order lifecycle and provide status updates

**Order States:**
```
CREATED → SEATS_RESERVED → PAYMENT_PENDING → PAYMENT_PROCESSING →
    → CONFIRMED (success)
    → FAILED (payment failed after retries)
    → EXPIRED (seat timer expired)
```

**Query: GetOrderStatus**
```go
func (w *BookingWorkflow) GetStatus() OrderStatus {
    return OrderStatus{
        State:           w.state,
        SelectedSeats:   w.seats,
        TimerRemaining:  w.timerEnd.Sub(time.Now()),
        PaymentAttempts: w.paymentAttempts,
        ErrorMessage:    w.lastError,
    }
}
```

### Feature 5: Real-time Frontend Updates

**Purpose:** Keep user informed of booking progress

**Implementation:** Polling with `useOrderStatus` hook
```javascript
// Poll every 2 seconds for status updates
const useOrderStatus = (orderId) => {
  const [status, setStatus] = useState(null);

  useEffect(() => {
    const interval = setInterval(async () => {
      const res = await fetch(`/api/orders/${orderId}/status`);
      setStatus(await res.json());
    }, 2000);

    return () => clearInterval(interval);
  }, [orderId]);

  return status;
};
```

**Displayed Information:**
- Countdown timer (MM:SS format)
- Selected seats list
- Current order state
- Payment attempt count (during processing)
- Error messages (on failure)

## 8. Technology Stack

### Backend

| Technology | Version | Purpose |
|------------|---------|---------|
| Go | 1.21+ | API server and Temporal workers |
| Temporal Go SDK | 1.25+ | Workflow and activity implementation |
| Chi Router | 5.x | HTTP routing |
| pgx | 5.x | PostgreSQL driver |
| go-redis | 9.x | Redis client for distributed locks |

### Infrastructure

| Technology | Version | Purpose |
|------------|---------|---------|
| Temporal Server | 1.22+ | Workflow orchestration |
| PostgreSQL | 15+ | Flight and order persistence |
| Redis | 7+ | Distributed seat locks with TTL |

### Frontend

| Technology | Version | Purpose |
|------------|---------|---------|
| React | 18.x | UI framework |
| Vite | 5.x | Build tool |
| TanStack Query | 5.x | Data fetching and caching |

### Development

| Technology | Purpose |
|------------|---------|
| Docker Compose | Local development environment |
| golang-migrate | Database migrations |
| Air | Go hot reload |

## 9. Security & Configuration

### Configuration (Environment Variables)

```bash
# Server
SERVER_PORT=8080
SERVER_HOST=0.0.0.0

# Database
DATABASE_URL=postgres://user:pass@localhost:5432/flights?sslmode=disable

# Redis
REDIS_URL=redis://localhost:6379/0

# Temporal
TEMPORAL_HOST=localhost:7233
TEMPORAL_NAMESPACE=default
TEMPORAL_TASK_QUEUE=booking-queue

# Timeouts (configurable for testing)
SEAT_RESERVATION_TIMEOUT=15m
PAYMENT_VALIDATION_TIMEOUT=10s
PAYMENT_MAX_RETRIES=3
PAYMENT_FAILURE_RATE=0.15
```

### Security Scope

**In Scope (MVP):**
- ✅ Input validation on all API endpoints
- ✅ SQL injection prevention via parameterized queries
- ✅ CORS configuration for frontend origin

**Out of Scope (MVP):**
- ❌ User authentication
- ❌ Rate limiting
- ❌ HTTPS (local dev only)
- ❌ API key authentication

## 10. API Specification

### Endpoints

#### List Flights
```
GET /api/flights

Response 200:
{
  "flights": [
    {
      "id": "fl-101",
      "flightNumber": "FL101",
      "origin": "NYC",
      "destination": "LAX",
      "departureTime": "2024-03-15T10:00:00Z",
      "totalSeats": 120,
      "availableSeats": 45
    }
  ]
}
```

#### Get Flight Details with Seat Map
```
GET /api/flights/{flightId}

Response 200:
{
  "id": "fl-101",
  "flightNumber": "FL101",
  "origin": "NYC",
  "destination": "LAX",
  "departureTime": "2024-03-15T10:00:00Z",
  "seatMap": {
    "rows": 20,
    "seatsPerRow": 6,
    "seats": [
      {"id": "1A", "row": 1, "column": "A", "status": "available"},
      {"id": "1B", "row": 1, "column": "B", "status": "reserved"},
      ...
    ]
  }
}
```

#### Create Order (Start Booking)
```
POST /api/orders

Request:
{
  "flightId": "fl-101",
  "seats": ["12A", "12B"]
}

Response 201:
{
  "orderId": "ord-abc123",
  "workflowId": "booking-ord-abc123",
  "status": "SEATS_RESERVED",
  "expiresAt": "2024-03-15T09:15:00Z"
}
```

#### Update Seat Selection (Signal Workflow)
```
PUT /api/orders/{orderId}/seats

Request:
{
  "seats": ["14A", "14B"]
}

Response 200:
{
  "orderId": "ord-abc123",
  "status": "SEATS_RESERVED",
  "seats": ["14A", "14B"],
  "expiresAt": "2024-03-15T09:30:00Z"  // Timer refreshed
}
```

#### Get Order Status (Query Workflow)
```
GET /api/orders/{orderId}/status

Response 200:
{
  "orderId": "ord-abc123",
  "status": "PAYMENT_PROCESSING",
  "seats": ["14A", "14B"],
  "timerRemaining": 845,  // seconds
  "paymentAttempts": 2,
  "lastError": "Payment validation failed"
}
```

#### Submit Payment
```
POST /api/orders/{orderId}/pay

Request:
{
  "paymentCode": "12345"
}

Response 202:
{
  "orderId": "ord-abc123",
  "status": "PAYMENT_PROCESSING"
}

// Poll GET /api/orders/{orderId}/status for result
```

#### Error Responses
```
400 Bad Request:
{
  "error": "INVALID_SEATS",
  "message": "Seats 12A, 12B are no longer available"
}

404 Not Found:
{
  "error": "ORDER_NOT_FOUND",
  "message": "Order ord-xyz not found"
}

409 Conflict:
{
  "error": "ORDER_EXPIRED",
  "message": "Seat reservation has expired"
}
```

## 11. Success Criteria

### MVP Success Definition

The MVP is successful when a user can complete the full booking flow:
1. View flight with seat map
2. Select seats and see timer start
3. Change seats and see timer refresh
4. Enter payment code and see retry feedback
5. Receive confirmation or failure message

### Functional Requirements

- ✅ Seat selection acquires distributed lock
- ✅ 15-minute timer counts down accurately (within 1 second)
- ✅ Seat changes refresh timer to full 15 minutes
- ✅ Expired reservations release seats automatically
- ✅ Payment validation times out at 10 seconds
- ✅ Failed payments retry up to 3 times
- ✅ ~15% of payments fail (simulated)
- ✅ Successful payment confirms order and persists seats
- ✅ Failed payment (after retries) releases seats and shows error

### Quality Indicators

- API response time < 200ms for status queries
- Timer accuracy within 1 second of actual
- Zero seat double-bookings (distributed lock integrity)
- Graceful handling of Temporal server restarts (workflow durability)

## 12. Implementation Phases

### Phase 1: Infrastructure & Data Layer

**Goal:** Set up project structure, database, and basic data models

**Deliverables:**
- ✅ Go project structure with cmd/internal layout
- ✅ PostgreSQL schema for flights, orders, seats
- ✅ Database migrations
- ✅ Redis connection and seat lock helpers
- ✅ Docker Compose with Postgres, Redis, Temporal
- ✅ Basic domain models

**Validation:**
- `docker-compose up` starts all services
- Can connect to Temporal UI at localhost:8233
- Migration runs successfully

### Phase 2: Temporal Workflows & Activities

**Goal:** Implement core booking workflow with seat reservation and payment

**Deliverables:**
- ✅ BookingWorkflow with timer and signal handling
- ✅ SeatReservationActivity (acquire/release/refresh locks)
- ✅ PaymentValidationActivity (with retry policy)
- ✅ OrderManagementActivity (state transitions)
- ✅ Worker registration and startup
- ✅ Workflow queries for status

**Validation:**
- Worker connects to Temporal and registers workflows
- Can start workflow via `tctl` or Temporal UI
- Timer fires after 15 minutes (test with shorter duration)
- Signals update workflow state

### Phase 3: REST API Layer

**Goal:** Expose booking functionality via HTTP endpoints

**Deliverables:**
- ✅ Flight listing and detail endpoints
- ✅ Order creation (starts workflow)
- ✅ Seat update endpoint (sends signal)
- ✅ Status endpoint (queries workflow)
- ✅ Payment submission endpoint
- ✅ Error handling and validation

**Validation:**
- Can complete full booking flow via curl/Postman
- Status polling shows accurate timer countdown
- Seat updates refresh timer
- Payment retries visible in status

### Phase 4: React Frontend

**Goal:** Build minimal UI for booking flow

**Deliverables:**
- ✅ Flight list page
- ✅ Seat map component with selection
- ✅ Countdown timer component
- ✅ Payment form with code input
- ✅ Order status display
- ✅ Error and success messaging

**Validation:**
- Full user flow completable in browser
- Timer updates every second
- Seat selection shows visual feedback
- Payment retry status visible

## 13. Future Considerations

### Post-MVP Enhancements

**User Experience:**
- WebSocket/SSE for real-time updates (replace polling)
- Animated seat map with smooth transitions
- Mobile-responsive design

**Functionality:**
- User accounts and booking history
- Multiple flight search with filters
- Multi-passenger bookings
- Seat pricing tiers (economy, business, first class)
- Booking modification and cancellation

**Integration:**
- Real payment gateway (Stripe, etc.)
- Email confirmation via SendGrid
- SMS notifications via Twilio

**Operations:**
- Admin dashboard for flight management
- Temporal Cloud deployment
- Observability (metrics, tracing, logging)
- Load testing and performance optimization

## 14. Risks & Mitigations

### Risk 1: Distributed Lock Race Conditions

**Risk:** Redis lock and Temporal timer could get out of sync, leading to seat conflicts.

**Mitigation:**
- Use both mechanisms as defense in depth
- Redis TTL slightly longer than Temporal timer (16 min vs 15 min)
- Temporal workflow is source of truth; Redis is optimization
- Periodic reconciliation activity to clean up orphaned locks

### Risk 2: Temporal Server Unavailability

**Risk:** If Temporal server goes down, bookings cannot proceed.

**Mitigation:**
- Temporal workflows are durable; they resume when server recovers
- Frontend shows appropriate "service unavailable" message
- For production: multi-node Temporal cluster

### Risk 3: Payment Simulation Confusion

**Risk:** 15% failure rate may confuse users expecting real payment processing.

**Mitigation:**
- Clear UI messaging: "Demo Mode - Payments have 15% simulated failure rate"
- Configuration flag to disable simulation for demos
- Consistent behavior for specific test codes (e.g., "99999" always fails with retries, "00000" always succeeds)

### Risk 4: Timer Drift

**Risk:** Frontend timer could drift from server timer, showing incorrect remaining time.

**Mitigation:**
- Server returns absolute expiration time, not duration
- Frontend calculates remaining from server time
- Periodic sync every 30 seconds via status poll

### Risk 5: Seat Map Performance

**Risk:** Large flights (300+ seats) could cause slow rendering.

**Mitigation:**
- Virtual scrolling for large seat maps
- Paginate seat map by section (front/middle/rear)
- Cache seat status with short TTL

## 15. Appendix

### Related Documents

- System Flow Diagram: `system_flow.png`

### Key Dependencies

| Dependency | Documentation |
|------------|---------------|
| Temporal Go SDK | https://docs.temporal.io/dev-guide/go |
| Chi Router | https://go-chi.io/ |
| go-redis | https://redis.uptrace.dev/ |
| pgx | https://github.com/jackc/pgx |
| React | https://react.dev/ |
| TanStack Query | https://tanstack.com/query |

### Temporal Workflow Patterns Reference

- [Timers](https://docs.temporal.io/dev-guide/go/features#timers)
- [Signals](https://docs.temporal.io/dev-guide/go/features#signals)
- [Queries](https://docs.temporal.io/dev-guide/go/features#queries)
- [Activity Retries](https://docs.temporal.io/dev-guide/go/features#activity-retries)
