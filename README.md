# Flight Booking System

A full-stack flight booking application demonstrating **Temporal workflow orchestration** patterns for distributed systems.

## Features

- **Time-limited seat reservations** (15-minute hold with refresh)
- **Distributed seat locking** using Redis
- **Payment retry logic** with simulated failures (3 attempts, 15% failure rate)
- **Real-time status updates** via polling
- **Automatic compensation** on booking failure or timeout
- **Durable workflows** survive server restarts

## Tech Stack

**Backend:** Go, Temporal, PostgreSQL, Redis
**Frontend:** React, Vite, TanStack Query, Tailwind CSS

## Quick Start

```bash
# 1. Start infrastructure (Temporal, PostgreSQL, Redis)
make up && make migrate-up

# 2. Start Temporal worker (terminal 1)
make run-worker

# 3. Start API server (terminal 2)
make run-server

# 4. Start frontend (terminal 3)
cd web && npm install && npm run dev
```

Open http://localhost:3000 to book flights.

**Temporal UI:** http://localhost:8233 (view workflow executions)

## Architecture

```
React Frontend (port 3000)
      ↓
Go REST API (port 8080)
      ↓
Temporal Workflows
      ↓
PostgreSQL + Redis
```

**Workflow Pattern:** Long-running `BookingWorkflow` orchestrates seat reservation → payment validation → order confirmation, with automatic seat release on timeout or failure.

## Project Structure

```
cmd/
  server/          # HTTP API server
  worker/          # Temporal worker
internal/
  api/             # HTTP handlers and routes
  domain/          # Core entities (Flight, Order, Seat)
  repository/      # Data access layer
  service/         # Business logic
  temporal/
    workflows/     # BookingWorkflow
    activities/    # Seat reservation, payment, order management
web/
  src/
    components/    # UI components (SeatMap, Timer, etc.)
    pages/         # FlightListPage, BookingPage
    hooks/         # useOrderStatus (polling), useCountdown
```

## Key Workflows

**BookingWorkflow** (`internal/temporal/workflows/booking_workflow.go`):
- Reserves seats with Redis distributed lock
- 15-minute timer with signal-based refresh
- Payment validation with 3 retry attempts
- Automatic seat release on expiration/failure

**Signals:**
- `UpdateSeats` - Change seat selection (resets timer)
- `ProceedToPay` - Submit payment code
- `CancelBooking` - Cancel reservation

**Queries:**
- `GetStatus` - Real-time order status

## Development Commands

```bash
make up              # Start Docker services
make down            # Stop Docker services
make migrate-up      # Run database migrations
make run-server      # Run API server
make run-worker      # Run Temporal worker
make dev-web         # Run frontend dev server
make test            # Run tests
```

## Documentation

- **[PRD.md](PRD.md)** - Comprehensive product requirements
- **[CLAUDE.md](CLAUDE.md)** - Development guidelines and conventions

## Demo Flow

1. Browse available flights
2. Select seats from visual seat map
3. Watch 15-minute countdown timer
4. Change seats (timer resets)
5. Enter 5-digit payment code
6. View confirmation or retry on failure

