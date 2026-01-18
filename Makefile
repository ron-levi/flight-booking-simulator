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

# Frontend
install-web:
	cd web && npm install

dev-web:
	cd web && npm run dev

build-web:
	cd web && npm run build

# Full stack development (run in separate terminals)
dev-all:
	@echo "Run these in separate terminals:"
	@echo "  Terminal 1: make up && make migrate-up"
	@echo "  Terminal 2: make run-worker"
	@echo "  Terminal 3: make run-server"
	@echo "  Terminal 4: make dev-web"
