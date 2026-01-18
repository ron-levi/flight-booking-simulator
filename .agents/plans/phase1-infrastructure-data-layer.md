# Feature: Phase 1 - Infrastructure & Data Layer

The following plan should be complete, but it's important that you validate documentation and codebase patterns and task sanity before you start implementing.

Pay special attention to naming of existing utils types and models. Import from the right files etc.

## Feature Description

Set up the foundational infrastructure for the flight booking system including:
- Go project structure with cmd/internal layout
- PostgreSQL schema for flights, orders, and seats
- Database migrations using golang-migrate
- Redis connection and seat lock helpers
- Docker Compose with PostgreSQL, Redis, and Temporal
- Basic domain models

This phase establishes all infrastructure required before implementing Temporal workflows in Phase 2.

## User Story

As a developer
I want to have a fully configured development environment with database schemas and infrastructure
So that I can build the flight booking system with proper data persistence and distributed locking

## Problem Statement

The flight booking system requires a robust infrastructure foundation including:
- Temporal server for workflow orchestration
- PostgreSQL for persistent storage of flights, orders, and seats
- Redis for distributed seat locks with TTL
- A well-organized Go project structure following best practices

## Solution Statement

Create a complete infrastructure layer with:
1. Docker Compose orchestrating all required services (Temporal, PostgreSQL, Redis)
2. Go project structure following cmd/internal pattern from CLAUDE.md
3. Database migrations for flights, orders, and seats tables
4. Domain models and repository layer with connection pooling
5. Configuration management via environment variables

## Feature Metadata

**Feature Type**: New Capability (Infrastructure Foundation)
**Estimated Complexity**: Medium
**Primary Systems Affected**: Database, Docker, Go project structure
**Dependencies**: Docker, Go 1.21+, PostgreSQL 15+, Redis 7+, Temporal 1.22+

---

## CONTEXT REFERENCES

### Relevant Codebase Files - IMPORTANT: YOU MUST READ THESE FILES BEFORE IMPLEMENTING!

- `CLAUDE.md` - Project coding guidelines, structure patterns, and conventions
- `PRD.md` (lines 174-227) - Directory structure specification
- `PRD.md` (lines 405-440) - Technology stack versions
- `PRD.md` (lines 443-466) - Configuration environment variables
- `PRD.md` (lines 480-617) - API specification (for understanding data models)

### New Files to Create

**Project Root:**
- `go.mod` - Go module definition
- `docker-compose.yml` - Container orchestration
- `Makefile` - Build and development commands
- `.env.example` - Environment variable template
- `.gitignore` - Git ignore patterns

**Entrypoints:**
- `cmd/server/main.go` - REST API server entrypoint
- `cmd/worker/main.go` - Temporal worker entrypoint

**Domain Models:**
- `internal/domain/flight.go` - Flight entity
- `internal/domain/order.go` - Order entity
- `internal/domain/seat.go` - Seat entity
- `internal/domain/errors.go` - Domain errors

**Database:**
- `internal/database/postgres.go` - PostgreSQL connection pool
- `internal/database/redis.go` - Redis client setup
- `internal/database/migrations/000001_create_flights_table.up.sql`
- `internal/database/migrations/000001_create_flights_table.down.sql`
- `internal/database/migrations/000002_create_orders_table.up.sql`
- `internal/database/migrations/000002_create_orders_table.down.sql`
- `internal/database/migrations/000003_create_seats_table.up.sql`
- `internal/database/migrations/000003_create_seats_table.down.sql`
- `internal/database/migrations/000004_seed_flights.up.sql`
- `internal/database/migrations/000004_seed_flights.down.sql`

**Repositories:**
- `internal/repository/flight_repo.go` - Flight data access
- `internal/repository/order_repo.go` - Order data access
- `internal/repository/seat_lock.go` - Redis seat lock operations

**Configuration:**
- `internal/config/config.go` - Configuration loading

### Relevant Documentation - YOU SHOULD READ THESE BEFORE IMPLEMENTING!

- [Temporal Go SDK](https://docs.temporal.io/develop/go) - Worker setup patterns
- [pgx v5 Connection Pool](https://pkg.go.dev/github.com/jackc/pgx/v5/pgxpool) - Database pooling
- [go-redis v9](https://redis.uptrace.dev/guide/go-redis.html) - Redis client patterns
- [golang-migrate](https://github.com/golang-migrate/migrate) - Migration patterns
- [Chi Router v5](https://go-chi.io/) - HTTP routing
- [Temporal Docker Compose](https://github.com/temporalio/docker-compose) - Infrastructure setup

### Patterns to Follow

**Naming Conventions (from CLAUDE.md):**
```go
// Short, clear names
func (s *Service) GetFlight(id string) (*Flight, error)
func (r *Repo) FindByID(ctx context.Context, id string) (*Order, error)

// DON'T stutter
// BAD: func (s *FlightService) GetFlightByFlightID(flightID string)
```

**Error Handling (from CLAUDE.md):**
```go
// DO: Return errors with context
if err != nil {
    return fmt.Errorf("failed to reserve seat %s: %w", seatID, err)
}
```

**Function Guidelines:**
- Max 40 lines per function
- Max 3 parameters - use struct for more
- Single return path when possible (early returns for errors)

---

## IMPLEMENTATION PLAN

### Phase 1.1: Project Setup & Docker Infrastructure

Set up Go module, Docker Compose with all services, and basic configuration.

**Tasks:**
- Initialize Go module with dependencies
- Create Docker Compose with Temporal, PostgreSQL, Redis
- Set up environment configuration
- Create Makefile for common commands

### Phase 1.2: Database Schema & Migrations

Create PostgreSQL schema for flights, orders, and seats using golang-migrate.

**Tasks:**
- Create migrations directory structure
- Write migration files for all tables
- Add seed data for demo flights

### Phase 1.3: Domain Models

Define core domain models matching the database schema.

**Tasks:**
- Create Flight, Order, Seat domain structs
- Define order status constants
- Create domain-level errors

### Phase 1.4: Database Connections

Set up PostgreSQL connection pool and Redis client.

**Tasks:**
- Implement PostgreSQL pool with pgx
- Implement Redis client with go-redis
- Add health check methods

### Phase 1.5: Repository Layer

Implement data access layer for flights, orders, and seat locks.

**Tasks:**
- Flight repository with CRUD operations
- Order repository with CRUD operations
- Seat lock repository with Redis operations

### Phase 1.6: Entrypoints (Stubs)

Create minimal server and worker entrypoints for testing infrastructure.

**Tasks:**
- Server main.go with health endpoint
- Worker main.go stub for Phase 2

---

## STEP-BY-STEP TASKS

IMPORTANT: Execute every task in order, top to bottom. Each task is atomic and independently testable.

---

### Task 1: CREATE `.gitignore`

- **IMPLEMENT**: Standard Go gitignore with local env files
- **VALIDATE**: `cat .gitignore | head -20`

```gitignore
# Binaries
*.exe
*.exe~
*.dll
*.so
*.dylib
bin/

# Test binary
*.test

# Output
*.out

# Dependency directories
vendor/

# Go workspace
go.work
go.work.sum

# Environment
.env
.env.local

# IDE
.idea/
.vscode/
*.swp
*.swo

# OS
.DS_Store
Thumbs.db

# Build
dist/
tmp/
```

---

### Task 2: CREATE `go.mod`

- **IMPLEMENT**: Go module with all required dependencies
- **VALIDATE**: `go mod tidy && go mod verify`

```go
module github.com/flight-booking-system

go 1.21

require (
	github.com/go-chi/chi/v5 v5.0.12
	github.com/golang-migrate/migrate/v4 v4.17.0
	github.com/jackc/pgx/v5 v5.5.3
	github.com/redis/go-redis/v9 v9.4.0
	go.temporal.io/sdk v1.26.1
)
```

---

### Task 3: CREATE `.env.example`

- **IMPLEMENT**: Environment variable template matching PRD section 9
- **VALIDATE**: `cat .env.example`

```bash
# Server
SERVER_PORT=8080
SERVER_HOST=0.0.0.0

# Database
DATABASE_HOST=localhost
DATABASE_PORT=5433
DATABASE_USER=flightapp
DATABASE_PASSWORD=flightapp
DATABASE_NAME=flight_booking
DATABASE_SSLMODE=disable

# Redis
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

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

---

### Task 4: CREATE `docker-compose.yml`

- **IMPLEMENT**: Full infrastructure with Temporal, PostgreSQL, Redis
- **PATTERN**: Use health checks for service dependencies
- **VALIDATE**: `docker compose config`

```yaml
version: "3.8"

services:
  # Temporal's PostgreSQL (separate from app DB)
  temporal-postgresql:
    container_name: temporal-postgresql
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: temporal
      POSTGRES_PASSWORD: temporal
    volumes:
      - temporal-postgresql-data:/var/lib/postgresql/data
    networks:
      - flight-booking-network
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U temporal"]
      interval: 5s
      timeout: 5s
      retries: 5

  # Temporal Server
  temporal:
    container_name: temporal
    image: temporalio/auto-setup:1.22.4
    depends_on:
      temporal-postgresql:
        condition: service_healthy
    environment:
      - DB=postgres12
      - DB_PORT=5432
      - POSTGRES_USER=temporal
      - POSTGRES_PWD=temporal
      - POSTGRES_SEEDS=temporal-postgresql
    ports:
      - "7233:7233"
    networks:
      - flight-booking-network
    healthcheck:
      test: ["CMD", "tctl", "--address", "temporal:7233", "cluster", "health"]
      interval: 10s
      timeout: 5s
      retries: 10
      start_period: 30s

  # Temporal Web UI
  temporal-ui:
    container_name: temporal-ui
    image: temporalio/ui:2.22.3
    depends_on:
      temporal:
        condition: service_healthy
    environment:
      - TEMPORAL_ADDRESS=temporal:7233
      - TEMPORAL_CORS_ORIGINS=http://localhost:3000
    ports:
      - "8233:8080"
    networks:
      - flight-booking-network

  # Application PostgreSQL
  app-postgresql:
    container_name: flight-app-db
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: flightapp
      POSTGRES_PASSWORD: flightapp
      POSTGRES_DB: flight_booking
    ports:
      - "5433:5432"
    volumes:
      - app-postgresql-data:/var/lib/postgresql/data
    networks:
      - flight-booking-network
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U flightapp -d flight_booking"]
      interval: 5s
      timeout: 5s
      retries: 5

  # Redis for seat locks
  redis:
    container_name: flight-redis
    image: redis:7-alpine
    command: redis-server --appendonly yes
    ports:
      - "6379:6379"
    volumes:
      - redis-data:/data
    networks:
      - flight-booking-network
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 5s
      timeout: 3s
      retries: 5

networks:
  flight-booking-network:
    driver: bridge
    name: flight-booking-network

volumes:
  temporal-postgresql-data:
  app-postgresql-data:
  redis-data:
```

---

### Task 5: CREATE `Makefile`

- **IMPLEMENT**: Common development commands
- **VALIDATE**: `make help`

```makefile
.PHONY: help up down logs migrate-up migrate-down migrate-create build run test lint

# Default target
help:
	@echo "Flight Booking System - Development Commands"
	@echo ""
	@echo "Infrastructure:"
	@echo "  make up              - Start all Docker services"
	@echo "  make down            - Stop all Docker services"
	@echo "  make logs            - Tail Docker logs"
	@echo ""
	@echo "Database:"
	@echo "  make migrate-up      - Run all migrations"
	@echo "  make migrate-down    - Rollback last migration"
	@echo "  make migrate-create  - Create new migration (NAME=migration_name)"
	@echo ""
	@echo "Development:"
	@echo "  make build           - Build server and worker binaries"
	@echo "  make run-server      - Run API server"
	@echo "  make run-worker      - Run Temporal worker"
	@echo "  make test            - Run all tests"
	@echo "  make lint            - Run linter"

# Database URL for migrations
DATABASE_URL ?= postgres://flightapp:flightapp@localhost:5433/flight_booking?sslmode=disable

# Infrastructure
up:
	docker compose up -d
	@echo "Waiting for services to be healthy..."
	@sleep 5
	docker compose ps

down:
	docker compose down

logs:
	docker compose logs -f

# Migrations
migrate-up:
	migrate -database "$(DATABASE_URL)" -path internal/database/migrations up

migrate-down:
	migrate -database "$(DATABASE_URL)" -path internal/database/migrations down 1

migrate-create:
	@if [ -z "$(NAME)" ]; then echo "Usage: make migrate-create NAME=migration_name"; exit 1; fi
	migrate create -ext sql -dir internal/database/migrations -seq $(NAME)

# Build
build:
	go build -o bin/server ./cmd/server
	go build -o bin/worker ./cmd/worker

# Run
run-server:
	go run ./cmd/server

run-worker:
	go run ./cmd/worker

# Test
test:
	go test -v ./...

# Lint
lint:
	golangci-lint run ./...
```

---

### Task 6: CREATE `internal/config/config.go`

- **IMPLEMENT**: Configuration loading from environment variables
- **PATTERN**: Struct with defaults, environment override
- **VALIDATE**: `go build ./internal/config`

```go
package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds all application configuration
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Temporal TemporalConfig
	Booking  BookingConfig
}

type ServerConfig struct {
	Host string
	Port int
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type TemporalConfig struct {
	Host      string
	Namespace string
	TaskQueue string
}

type BookingConfig struct {
	SeatReservationTimeout   time.Duration
	PaymentValidationTimeout time.Duration
	PaymentMaxRetries        int
	PaymentFailureRate       float64
}

// Load reads configuration from environment variables with defaults
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Port: getEnvInt("SERVER_PORT", 8080),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DATABASE_HOST", "localhost"),
			Port:     getEnvInt("DATABASE_PORT", 5433),
			User:     getEnv("DATABASE_USER", "flightapp"),
			Password: getEnv("DATABASE_PASSWORD", "flightapp"),
			Name:     getEnv("DATABASE_NAME", "flight_booking"),
			SSLMode:  getEnv("DATABASE_SSLMODE", "disable"),
		},
		Redis: RedisConfig{
			Addr:     getEnv("REDIS_ADDR", "localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
		},
		Temporal: TemporalConfig{
			Host:      getEnv("TEMPORAL_HOST", "localhost:7233"),
			Namespace: getEnv("TEMPORAL_NAMESPACE", "default"),
			TaskQueue: getEnv("TEMPORAL_TASK_QUEUE", "booking-queue"),
		},
		Booking: BookingConfig{
			SeatReservationTimeout:   getEnvDuration("SEAT_RESERVATION_TIMEOUT", 15*time.Minute),
			PaymentValidationTimeout: getEnvDuration("PAYMENT_VALIDATION_TIMEOUT", 10*time.Second),
			PaymentMaxRetries:        getEnvInt("PAYMENT_MAX_RETRIES", 3),
			PaymentFailureRate:       getEnvFloat("PAYMENT_FAILURE_RATE", 0.15),
		},
	}
}

// DatabaseURL returns the PostgreSQL connection string
func (c *DatabaseConfig) DatabaseURL() string {
	return "postgres://" + c.User + ":" + c.Password + "@" + c.Host + ":" + strconv.Itoa(c.Port) + "/" + c.Name + "?sslmode=" + c.SSLMode
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvFloat(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
			return floatVal
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
```

---

### Task 7: CREATE `internal/domain/errors.go`

- **IMPLEMENT**: Domain-level error definitions
- **VALIDATE**: `go build ./internal/domain`

```go
package domain

import "errors"

var (
	// ErrNotFound indicates a resource was not found
	ErrNotFound = errors.New("resource not found")

	// ErrFlightNotFound indicates a flight was not found
	ErrFlightNotFound = errors.New("flight not found")

	// ErrOrderNotFound indicates an order was not found
	ErrOrderNotFound = errors.New("order not found")

	// ErrSeatNotFound indicates a seat was not found
	ErrSeatNotFound = errors.New("seat not found")

	// ErrSeatUnavailable indicates a seat is not available for booking
	ErrSeatUnavailable = errors.New("seat is not available")

	// ErrSeatsAlreadyLocked indicates seats are already locked by another order
	ErrSeatsAlreadyLocked = errors.New("seats are already locked")

	// ErrInsufficientSeats indicates not enough seats available
	ErrInsufficientSeats = errors.New("insufficient seats available")

	// ErrOrderExpired indicates the order reservation has expired
	ErrOrderExpired = errors.New("order reservation has expired")

	// ErrInvalidPaymentCode indicates the payment code format is invalid
	ErrInvalidPaymentCode = errors.New("invalid payment code format")

	// ErrPaymentFailed indicates payment validation failed
	ErrPaymentFailed = errors.New("payment validation failed")

	// ErrInvalidOrderStatus indicates an invalid order status transition
	ErrInvalidOrderStatus = errors.New("invalid order status transition")
)
```

---

### Task 8: CREATE `internal/domain/flight.go`

- **IMPLEMENT**: Flight domain model matching PRD API spec
- **VALIDATE**: `go build ./internal/domain`

```go
package domain

import "time"

// Flight represents a flight in the system
type Flight struct {
	ID             string    `json:"id"`
	FlightNumber   string    `json:"flightNumber"`
	Origin         string    `json:"origin"`
	Destination    string    `json:"destination"`
	DepartureTime  time.Time `json:"departureTime"`
	ArrivalTime    time.Time `json:"arrivalTime"`
	TotalSeats     int       `json:"totalSeats"`
	AvailableSeats int       `json:"availableSeats"`
	PriceCents     int64     `json:"priceCents"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

// FlightWithSeats represents a flight with its seat map
type FlightWithSeats struct {
	Flight
	SeatMap SeatMap `json:"seatMap"`
}

// SeatMap represents the seat configuration of a flight
type SeatMap struct {
	Rows        int    `json:"rows"`
	SeatsPerRow int    `json:"seatsPerRow"`
	Seats       []Seat `json:"seats"`
}
```

---

### Task 9: CREATE `internal/domain/seat.go`

- **IMPLEMENT**: Seat domain model with status enum
- **VALIDATE**: `go build ./internal/domain`

```go
package domain

import "time"

// SeatStatus represents the current status of a seat
type SeatStatus string

const (
	SeatStatusAvailable SeatStatus = "available"
	SeatStatusReserved  SeatStatus = "reserved"
	SeatStatusBooked    SeatStatus = "booked"
)

// Seat represents an individual seat on a flight
type Seat struct {
	ID        string     `json:"id"`
	FlightID  string     `json:"flightId"`
	Row       int        `json:"row"`
	Column    string     `json:"column"`
	Status    SeatStatus `json:"status"`
	OrderID   *string    `json:"orderId,omitempty"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

// SeatID returns the seat identifier (e.g., "12A")
func (s *Seat) SeatID() string {
	return s.ID
}

// IsAvailable checks if the seat can be selected
func (s *Seat) IsAvailable() bool {
	return s.Status == SeatStatusAvailable
}
```

---

### Task 10: CREATE `internal/domain/order.go`

- **IMPLEMENT**: Order domain model with status enum matching PRD
- **VALIDATE**: `go build ./internal/domain`

```go
package domain

import "time"

// OrderStatus represents the current status of an order
type OrderStatus string

const (
	OrderStatusCreated           OrderStatus = "CREATED"
	OrderStatusSeatsReserved     OrderStatus = "SEATS_RESERVED"
	OrderStatusPaymentPending    OrderStatus = "PAYMENT_PENDING"
	OrderStatusPaymentProcessing OrderStatus = "PAYMENT_PROCESSING"
	OrderStatusConfirmed         OrderStatus = "CONFIRMED"
	OrderStatusFailed            OrderStatus = "FAILED"
	OrderStatusExpired           OrderStatus = "EXPIRED"
)

// Order represents a booking order
type Order struct {
	ID              string      `json:"id"`
	FlightID        string      `json:"flightId"`
	WorkflowID      string      `json:"workflowId"`
	Status          OrderStatus `json:"status"`
	Seats           []string    `json:"seats"`
	TotalPriceCents int64       `json:"totalPriceCents"`
	PaymentCode     *string     `json:"paymentCode,omitempty"`
	ExpiresAt       *time.Time  `json:"expiresAt,omitempty"`
	ConfirmedAt     *time.Time  `json:"confirmedAt,omitempty"`
	FailureReason   *string     `json:"failureReason,omitempty"`
	CreatedAt       time.Time   `json:"createdAt"`
	UpdatedAt       time.Time   `json:"updatedAt"`
}

// OrderStatusResponse represents the status response for polling
type OrderStatusResponse struct {
	OrderID         string      `json:"orderId"`
	Status          OrderStatus `json:"status"`
	Seats           []string    `json:"seats"`
	TimerRemaining  int         `json:"timerRemaining"` // seconds
	PaymentAttempts int         `json:"paymentAttempts"`
	LastError       string      `json:"lastError,omitempty"`
}

// IsTerminal returns true if the order is in a final state
func (o *Order) IsTerminal() bool {
	return o.Status == OrderStatusConfirmed ||
		o.Status == OrderStatusFailed ||
		o.Status == OrderStatusExpired
}

// CanTransitionTo checks if the order can transition to the given status
func (o *Order) CanTransitionTo(status OrderStatus) bool {
	validTransitions := map[OrderStatus][]OrderStatus{
		OrderStatusCreated:           {OrderStatusSeatsReserved, OrderStatusFailed},
		OrderStatusSeatsReserved:     {OrderStatusPaymentPending, OrderStatusExpired, OrderStatusFailed},
		OrderStatusPaymentPending:    {OrderStatusPaymentProcessing, OrderStatusExpired, OrderStatusFailed},
		OrderStatusPaymentProcessing: {OrderStatusConfirmed, OrderStatusFailed},
	}

	allowed, exists := validTransitions[o.Status]
	if !exists {
		return false
	}

	for _, s := range allowed {
		if s == status {
			return true
		}
	}
	return false
}
```

---

### Task 11: CREATE `internal/database/postgres.go`

- **IMPLEMENT**: PostgreSQL connection pool using pgx
- **PATTERN**: Context-aware connection with health check
- **VALIDATE**: `go build ./internal/database`

```go
package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/flight-booking-system/internal/config"
)

// NewPostgresPool creates a new PostgreSQL connection pool
func NewPostgresPool(ctx context.Context, cfg config.DatabaseConfig) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.DatabaseURL())
	if err != nil {
		return nil, fmt.Errorf("parse database config: %w", err)
	}

	// Configure pool settings
	poolConfig.MaxConns = 25
	poolConfig.MinConns = 5
	poolConfig.MaxConnLifetime = time.Hour
	poolConfig.MaxConnIdleTime = 30 * time.Minute
	poolConfig.HealthCheckPeriod = time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("create database pool: %w", err)
	}

	// Verify connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return pool, nil
}

// HealthCheck verifies the database connection is healthy
func HealthCheck(ctx context.Context, pool *pgxpool.Pool) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return pool.Ping(ctx)
}
```

---

### Task 12: CREATE `internal/database/redis.go`

- **IMPLEMENT**: Redis client setup using go-redis
- **PATTERN**: Connection with health check
- **VALIDATE**: `go build ./internal/database`

```go
package database

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/flight-booking-system/internal/config"
)

// NewRedisClient creates a new Redis client
func NewRedisClient(ctx context.Context, cfg config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     10,
		MinIdleConns: 5,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		DialTimeout:  5 * time.Second,
		PoolTimeout:  4 * time.Second,
	})

	// Verify connection
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("ping redis: %w", err)
	}

	return client, nil
}

// RedisHealthCheck verifies the Redis connection is healthy
func RedisHealthCheck(ctx context.Context, client *redis.Client) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return client.Ping(ctx).Err()
}
```

---

### Task 13: CREATE `internal/database/migrations/000001_create_flights_table.up.sql`

- **IMPLEMENT**: Flights table with indexes
- **VALIDATE**: Run migration after Docker is up

```sql
BEGIN;

CREATE TABLE IF NOT EXISTS flights (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    flight_number VARCHAR(10) NOT NULL,
    origin VARCHAR(3) NOT NULL,
    destination VARCHAR(3) NOT NULL,
    departure_time TIMESTAMPTZ NOT NULL,
    arrival_time TIMESTAMPTZ NOT NULL,
    total_seats INTEGER NOT NULL,
    available_seats INTEGER NOT NULL,
    price_cents BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT flights_flight_number_unique UNIQUE (flight_number),
    CONSTRAINT flights_seats_check CHECK (available_seats >= 0 AND available_seats <= total_seats)
);

CREATE INDEX idx_flights_departure ON flights(departure_time);
CREATE INDEX idx_flights_route ON flights(origin, destination);
CREATE INDEX idx_flights_available ON flights(available_seats) WHERE available_seats > 0;

COMMIT;
```

---

### Task 14: CREATE `internal/database/migrations/000001_create_flights_table.down.sql`

- **IMPLEMENT**: Rollback flights table
- **VALIDATE**: `cat internal/database/migrations/000001_create_flights_table.down.sql`

```sql
DROP TABLE IF EXISTS flights;
```

---

### Task 15: CREATE `internal/database/migrations/000002_create_orders_table.up.sql`

- **IMPLEMENT**: Orders table with status tracking
- **VALIDATE**: Run migration after Docker is up

```sql
BEGIN;

CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    flight_id UUID NOT NULL REFERENCES flights(id),
    workflow_id VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'CREATED',
    seats TEXT[] NOT NULL DEFAULT '{}',
    total_price_cents BIGINT NOT NULL DEFAULT 0,
    payment_code VARCHAR(5),
    expires_at TIMESTAMPTZ,
    confirmed_at TIMESTAMPTZ,
    failure_reason TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT orders_workflow_id_unique UNIQUE (workflow_id),
    CONSTRAINT orders_status_check CHECK (status IN (
        'CREATED', 'SEATS_RESERVED', 'PAYMENT_PENDING',
        'PAYMENT_PROCESSING', 'CONFIRMED', 'FAILED', 'EXPIRED'
    ))
);

CREATE INDEX idx_orders_flight ON orders(flight_id);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_expires ON orders(expires_at) WHERE status IN ('SEATS_RESERVED', 'PAYMENT_PENDING');

COMMIT;
```

---

### Task 16: CREATE `internal/database/migrations/000002_create_orders_table.down.sql`

- **IMPLEMENT**: Rollback orders table
- **VALIDATE**: `cat internal/database/migrations/000002_create_orders_table.down.sql`

```sql
DROP TABLE IF EXISTS orders;
```

---

### Task 17: CREATE `internal/database/migrations/000003_create_seats_table.up.sql`

- **IMPLEMENT**: Seats table with status tracking
- **VALIDATE**: Run migration after Docker is up

```sql
BEGIN;

CREATE TABLE IF NOT EXISTS seats (
    id VARCHAR(10) NOT NULL,
    flight_id UUID NOT NULL REFERENCES flights(id) ON DELETE CASCADE,
    row_num INTEGER NOT NULL,
    col VARCHAR(1) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'available',
    order_id UUID REFERENCES orders(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (flight_id, id),
    CONSTRAINT seats_status_check CHECK (status IN ('available', 'reserved', 'booked'))
);

CREATE INDEX idx_seats_status ON seats(flight_id, status);
CREATE INDEX idx_seats_order ON seats(order_id) WHERE order_id IS NOT NULL;

COMMIT;
```

---

### Task 18: CREATE `internal/database/migrations/000003_create_seats_table.down.sql`

- **IMPLEMENT**: Rollback seats table
- **VALIDATE**: `cat internal/database/migrations/000003_create_seats_table.down.sql`

```sql
DROP TABLE IF EXISTS seats;
```

---

### Task 19: CREATE `internal/database/migrations/000004_seed_flights.up.sql`

- **IMPLEMENT**: Seed demo flights and seats as specified in PRD
- **VALIDATE**: Run migration and query flights table

```sql
BEGIN;

-- Insert demo flights
INSERT INTO flights (id, flight_number, origin, destination, departure_time, arrival_time, total_seats, available_seats, price_cents)
VALUES
    ('550e8400-e29b-41d4-a716-446655440001', 'FL101', 'NYC', 'LAX', NOW() + INTERVAL '2 days', NOW() + INTERVAL '2 days' + INTERVAL '6 hours', 120, 120, 35000),
    ('550e8400-e29b-41d4-a716-446655440002', 'FL102', 'LAX', 'NYC', NOW() + INTERVAL '3 days', NOW() + INTERVAL '3 days' + INTERVAL '5 hours', 120, 120, 32000),
    ('550e8400-e29b-41d4-a716-446655440003', 'FL201', 'SFO', 'CHI', NOW() + INTERVAL '1 day', NOW() + INTERVAL '1 day' + INTERVAL '4 hours', 90, 90, 28000),
    ('550e8400-e29b-41d4-a716-446655440004', 'FL202', 'CHI', 'SFO', NOW() + INTERVAL '4 days', NOW() + INTERVAL '4 days' + INTERVAL '4 hours', 90, 90, 27500);

-- Generate seats for FL101 and FL102 (20 rows x 6 seats = 120 seats)
INSERT INTO seats (id, flight_id, row_num, col, status)
SELECT
    row_num || col AS id,
    flight_id,
    row_num,
    col,
    'available'
FROM (
    SELECT
        f.id AS flight_id,
        r.row_num,
        c.col
    FROM flights f
    CROSS JOIN generate_series(1, 20) AS r(row_num)
    CROSS JOIN (VALUES ('A'), ('B'), ('C'), ('D'), ('E'), ('F')) AS c(col)
    WHERE f.flight_number IN ('FL101', 'FL102')
) AS seat_data;

-- Generate seats for FL201 and FL202 (15 rows x 6 seats = 90 seats)
INSERT INTO seats (id, flight_id, row_num, col, status)
SELECT
    row_num || col AS id,
    flight_id,
    row_num,
    col,
    'available'
FROM (
    SELECT
        f.id AS flight_id,
        r.row_num,
        c.col
    FROM flights f
    CROSS JOIN generate_series(1, 15) AS r(row_num)
    CROSS JOIN (VALUES ('A'), ('B'), ('C'), ('D'), ('E'), ('F')) AS c(col)
    WHERE f.flight_number IN ('FL201', 'FL202')
) AS seat_data;

COMMIT;
```

---

### Task 20: CREATE `internal/database/migrations/000004_seed_flights.down.sql`

- **IMPLEMENT**: Remove seed data
- **VALIDATE**: `cat internal/database/migrations/000004_seed_flights.down.sql`

```sql
BEGIN;

DELETE FROM seats WHERE flight_id IN (
    SELECT id FROM flights WHERE flight_number IN ('FL101', 'FL102', 'FL201', 'FL202')
);

DELETE FROM flights WHERE flight_number IN ('FL101', 'FL102', 'FL201', 'FL202');

COMMIT;
```

---

### Task 21: CREATE `internal/repository/flight_repo.go`

- **IMPLEMENT**: Flight repository with CRUD operations
- **PATTERN**: Context-aware queries with pgx
- **VALIDATE**: `go build ./internal/repository`

```go
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
```

---

### Task 22: CREATE `internal/repository/order_repo.go`

- **IMPLEMENT**: Order repository with CRUD operations
- **PATTERN**: Context-aware queries with pgx
- **VALIDATE**: `go build ./internal/repository`

```go
package repository

import (
	"context"
	"errors"
	"fmt"

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
```

Add the missing import at the top of order_repo.go:
```go
import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/flight-booking-system/internal/domain"
)
```

---

### Task 23: CREATE `internal/repository/seat_lock.go`

- **IMPLEMENT**: Redis-based seat locking with TTL
- **PATTERN**: Distributed locking as described in PRD section 6
- **VALIDATE**: `go build ./internal/repository`

```go
package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// SeatLockRepo handles distributed seat locking via Redis
type SeatLockRepo struct {
	client *redis.Client
}

// NewSeatLockRepo creates a new SeatLockRepo
func NewSeatLockRepo(client *redis.Client) *SeatLockRepo {
	return &SeatLockRepo{client: client}
}

// seatLockKey generates the Redis key for a seat lock
func seatLockKey(flightID, seatID string) string {
	return fmt.Sprintf("seat:lock:%s:%s", flightID, seatID)
}

// LockSeats attempts to lock multiple seats for an order
// Returns nil if all seats were locked, error otherwise
func (r *SeatLockRepo) LockSeats(ctx context.Context, flightID string, seatIDs []string, orderID string, ttl time.Duration) error {
	// Use a pipeline for atomic operations
	pipe := r.client.TxPipeline()

	// First, check if any seats are already locked
	for _, seatID := range seatIDs {
		key := seatLockKey(flightID, seatID)
		pipe.Get(ctx, key)
	}

	results, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return fmt.Errorf("check existing locks: %w", err)
	}

	// Check results - if any seat is already locked by a different order, fail
	for i, result := range results {
		if result.Err() == nil {
			existingOrderID, _ := result.(*redis.StringCmd).Result()
			if existingOrderID != orderID {
				return fmt.Errorf("seat %s already locked by order %s", seatIDs[i], existingOrderID)
			}
		}
	}

	// Now set all locks with NX (only if not exists) or update if same order
	pipe = r.client.TxPipeline()
	for _, seatID := range seatIDs {
		key := seatLockKey(flightID, seatID)
		pipe.Set(ctx, key, orderID, ttl)
	}

	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("set seat locks: %w", err)
	}

	return nil
}

// TryLockSeat attempts to lock a single seat using SETNX
// Returns true if lock was acquired, false if already locked
func (r *SeatLockRepo) TryLockSeat(ctx context.Context, flightID, seatID, orderID string, ttl time.Duration) (bool, error) {
	key := seatLockKey(flightID, seatID)
	ok, err := r.client.SetNX(ctx, key, orderID, ttl).Result()
	if err != nil {
		return false, fmt.Errorf("setnx seat lock: %w", err)
	}
	return ok, nil
}

// ReleaseLocks releases all seat locks for an order
func (r *SeatLockRepo) ReleaseLocks(ctx context.Context, flightID string, seatIDs []string, orderID string) error {
	pipe := r.client.TxPipeline()

	for _, seatID := range seatIDs {
		key := seatLockKey(flightID, seatID)
		// Only delete if the lock belongs to this order (using Lua script)
		script := redis.NewScript(`
			if redis.call("get", KEYS[1]) == ARGV[1] then
				return redis.call("del", KEYS[1])
			else
				return 0
			end
		`)
		pipe.EvalSha(ctx, script.Hash(), []string{key}, orderID)
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("release seat locks: %w", err)
	}

	return nil
}

// ExtendLocks extends the TTL for all seat locks
func (r *SeatLockRepo) ExtendLocks(ctx context.Context, flightID string, seatIDs []string, orderID string, ttl time.Duration) error {
	pipe := r.client.TxPipeline()

	for _, seatID := range seatIDs {
		key := seatLockKey(flightID, seatID)
		// Only extend if the lock belongs to this order
		script := redis.NewScript(`
			if redis.call("get", KEYS[1]) == ARGV[1] then
				return redis.call("pexpire", KEYS[1], ARGV[2])
			else
				return 0
			end
		`)
		pipe.EvalSha(ctx, script.Hash(), []string{key}, orderID, ttl.Milliseconds())
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("extend seat locks: %w", err)
	}

	return nil
}

// IsLocked checks if a seat is currently locked
func (r *SeatLockRepo) IsLocked(ctx context.Context, flightID, seatID string) (bool, string, error) {
	key := seatLockKey(flightID, seatID)
	orderID, err := r.client.Get(ctx, key).Result()

	if err == redis.Nil {
		return false, "", nil
	}
	if err != nil {
		return false, "", fmt.Errorf("get seat lock: %w", err)
	}

	return true, orderID, nil
}

// GetLockedSeats returns all locked seat IDs for a flight
func (r *SeatLockRepo) GetLockedSeats(ctx context.Context, flightID string) (map[string]string, error) {
	pattern := fmt.Sprintf("seat:lock:%s:*", flightID)
	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, fmt.Errorf("get locked seat keys: %w", err)
	}

	if len(keys) == 0 {
		return make(map[string]string), nil
	}

	// Get all values
	pipe := r.client.Pipeline()
	cmds := make([]*redis.StringCmd, len(keys))
	for i, key := range keys {
		cmds[i] = pipe.Get(ctx, key)
	}

	_, err = pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("get locked seat values: %w", err)
	}

	result := make(map[string]string)
	for i, cmd := range cmds {
		if cmd.Err() == nil {
			// Extract seat ID from key (seat:lock:flightID:seatID)
			seatID := keys[i][len(fmt.Sprintf("seat:lock:%s:", flightID)):]
			result[seatID] = cmd.Val()
		}
	}

	return result, nil
}
```

---

### Task 24: CREATE `cmd/server/main.go`

- **IMPLEMENT**: Minimal API server with health endpoint
- **PATTERN**: Wire up config, database, and basic router
- **VALIDATE**: `go build ./cmd/server && ./bin/server` (with Docker running)

```go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/flight-booking-system/internal/config"
	"github.com/flight-booking-system/internal/database"
	"github.com/flight-booking-system/internal/repository"
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

	// Create repositories
	flightRepo := repository.NewFlightRepo(pool)
	_ = repository.NewOrderRepo(pool)
	_ = repository.NewSeatLockRepo(redisClient)

	// Create router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// Health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		// Check database
		if err := database.HealthCheck(r.Context(), pool); err != nil {
			http.Error(w, "database unhealthy", http.StatusServiceUnavailable)
			return
		}

		// Check Redis
		if err := database.RedisHealthCheck(r.Context(), redisClient); err != nil {
			http.Error(w, "redis unhealthy", http.StatusServiceUnavailable)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Temporary flights endpoint to verify database
	r.Get("/api/flights", func(w http.ResponseWriter, r *http.Request) {
		flights, err := flightRepo.FindAll(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"flights": flights,
		})
	})

	// Create server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
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

---

### Task 25: CREATE `cmd/worker/main.go`

- **IMPLEMENT**: Stub Temporal worker for Phase 2
- **PATTERN**: Basic worker setup without workflows yet
- **VALIDATE**: `go build ./cmd/worker`

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

	// TODO: Register workflows and activities in Phase 2
	// w.RegisterWorkflow(temporal.BookingWorkflow)
	// w.RegisterActivity(&temporal.BookingActivities{})

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

## TESTING STRATEGY

### Unit Tests

For Phase 1, focus on repository tests with a test database.

```go
// Example test pattern for repository tests
func TestFlightRepo_FindByID(t *testing.T) {
    // Setup test database connection
    // Insert test data
    // Call FindByID
    // Assert results
    // Cleanup
}
```

### Integration Tests

Test Docker Compose services are running and connected:
- PostgreSQL accepts connections
- Redis accepts connections
- Temporal server is healthy
- Migrations run successfully

### Edge Cases

- Database connection failure handling
- Redis connection failure handling
- Invalid flight/order IDs
- Concurrent seat locking attempts

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

### Level 2: Infrastructure

```bash
# Start Docker services
make up

# Verify all services are healthy
docker compose ps

# Check Temporal UI is accessible
curl -s http://localhost:8233 | head -5

# Check PostgreSQL connection
docker exec flight-app-db pg_isready -U flightapp -d flight_booking
```

### Level 3: Migrations

```bash
# Run all migrations
make migrate-up

# Verify tables exist
docker exec flight-app-db psql -U flightapp -d flight_booking -c "\dt"

# Verify seed data
docker exec flight-app-db psql -U flightapp -d flight_booking -c "SELECT flight_number, origin, destination, available_seats FROM flights;"
```

### Level 4: Server Health

```bash
# Build and run server
make build
./bin/server &

# Test health endpoint
curl -s http://localhost:8080/health

# Test flights endpoint
curl -s http://localhost:8080/api/flights | jq .

# Stop server
pkill -f "bin/server"
```

### Level 5: Worker Connection

```bash
# Build and run worker (briefly)
timeout 10 ./bin/worker || true

# Verify worker connected to Temporal (check logs)
```

---

## ACCEPTANCE CRITERIA

- [x] Go project structure follows cmd/internal pattern
- [ ] Docker Compose starts all services (Temporal, PostgreSQL, Redis)
- [ ] Temporal UI accessible at localhost:8233
- [ ] Database migrations run successfully
- [ ] Seed data creates 4 flights with seats
- [ ] Server health endpoint returns OK
- [ ] Server can query flights from database
- [ ] Worker can connect to Temporal server
- [ ] Redis seat lock operations work correctly
- [ ] All code builds without errors
- [ ] Configuration loads from environment variables

---

## COMPLETION CHECKLIST

- [ ] All 25 tasks completed in order
- [ ] Each task validation passed
- [ ] `docker compose up` starts all services
- [ ] `make migrate-up` runs all migrations
- [ ] `curl localhost:8080/health` returns OK
- [ ] `curl localhost:8080/api/flights` returns flight data
- [ ] `go build ./...` succeeds
- [ ] Code follows CLAUDE.md conventions

---

## NOTES

### Design Decisions

1. **Separate PostgreSQL instances**: Temporal uses its own PostgreSQL for durability, app uses separate instance for clean separation
2. **Redis for seat locks only**: Redis TTL provides automatic lock expiration as backup to Temporal timers
3. **UUIDs for IDs**: All entities use UUIDs for globally unique identifiers
4. **TEXT[] for seats**: PostgreSQL array type for flexible seat list storage in orders

### Gotchas

1. **Port 8233 for Temporal UI**: Not 8080 (which is reserved for our API server)
2. **Port 5433 for app database**: Avoids conflict with Temporal's PostgreSQL on 5432
3. **Redis Lua scripts**: Required for atomic check-and-delete operations in seat locks
4. **go.mod module path**: Use `github.com/flight-booking-system` consistently in all imports

### Future Considerations (Phase 2)

- Add Temporal workflow and activity implementations
- Add proper HTTP handlers with error responses
- Add request/response validation
- Add API documentation (OpenAPI)
