# Flight Booking System - Development Guidelines

## Project Context

Temporal-based flight booking system with Go backend and React frontend. See `.PRD.md` for full requirements.

## Core Principles

1. **Simplicity over cleverness** - Write obvious code that's easy to read and debug
2. **No premature abstraction** - Duplicate code 2-3 times before extracting
3. **Fail fast, fail loud** - Return errors early, log them clearly
4. **Thin layers** - Minimal indirection between request and business logic

---

## Go Guidelines

### Project Structure

```
cmd/           # Entrypoints only - minimal code, just wire things up
internal/      # All application code (not importable externally)
  api/         # HTTP handlers - thin, delegate to services
  domain/      # Plain structs, no methods beyond validation
  repository/  # Database/Redis access - no business logic
  service/     # Business logic orchestration
  temporal/    # Workflows and activities
```

### Error Handling

```go
// DO: Return errors with context
if err != nil {
    return fmt.Errorf("failed to reserve seat %s: %w", seatID, err)
}

// DON'T: Swallow errors or use panic for control flow
if err != nil {
    log.Println(err) // Lost context, caller doesn't know
    return nil
}
```

### Naming

```go
// DO: Short, clear names - especially for local scope
func (s *Service) GetFlight(id string) (*Flight, error)
func (r *Repo) FindByID(ctx context.Context, id string) (*Order, error)

// DON'T: Stutter or over-qualify
func (s *FlightService) GetFlightByFlightID(flightID string) (*FlightModel, error)
```

### Functions

- Max 40 lines per function - if longer, extract
- Max 3 parameters - use struct for more
- Single return path when possible (early returns for errors)

```go
// DO: Early returns, single happy path
func (s *Service) CreateOrder(ctx context.Context, req CreateOrderReq) (*Order, error) {
    if err := req.Validate(); err != nil {
        return nil, err
    }

    flight, err := s.flights.FindByID(ctx, req.FlightID)
    if err != nil {
        return nil, fmt.Errorf("flight lookup: %w", err)
    }

    if !flight.HasAvailableSeats(req.Seats) {
        return nil, ErrSeatsUnavailable
    }

    return s.orders.Create(ctx, flight.ID, req.Seats)
}
```

### HTTP Handlers

Keep handlers thin - parse request, call service, write response:

```go
func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
    var req CreateOrderRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.writeError(w, http.StatusBadRequest, "invalid request body")
        return
    }

    order, err := h.orderService.Create(r.Context(), req)
    if err != nil {
        h.handleServiceError(w, err)
        return
    }

    h.writeJSON(w, http.StatusCreated, order)
}
```

### Temporal Workflows

```go
// Workflows: Orchestration only, no I/O
func BookingWorkflow(ctx workflow.Context, orderID string) error {
    // Use activities for all external operations
    // Keep workflow logic deterministic
}

// Activities: Side effects happen here
func (a *Activities) ReserveSeats(ctx context.Context, orderID string, seats []string) error {
    // Database, Redis, external calls
}
```

### Testing

- Test behavior, not implementation
- One assertion per test when possible
- Use table-driven tests for multiple cases

```go
func TestPaymentValidation(t *testing.T) {
    tests := []struct {
        name    string
        code    string
        wantErr bool
    }{
        {"valid code", "12345", false},
        {"too short", "1234", true},
        {"too long", "123456", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidatePaymentCode(tt.code)
            if (err != nil) != tt.wantErr {
                t.Errorf("got err=%v, wantErr=%v", err, tt.wantErr)
            }
        })
    }
}
```

---

## React Guidelines

### Component Structure

```
web/src/
  components/    # Reusable UI components
  pages/         # Route-level components
  hooks/         # Custom hooks
  api/           # API client functions
  types/         # TypeScript types (if using TS)
```

### Component Design

- One component per file
- Max 100 lines per component - extract if larger
- Props interface at top of file

```jsx
// DO: Small, focused components
function Timer({ expiresAt }) {
  const remaining = useCountdown(expiresAt);
  return <span className="timer">{formatTime(remaining)}</span>;
}

// DON'T: Kitchen sink components
function BookingPage({ flight, order, user, settings, ... }) {
  // 500 lines of mixed concerns
}
```

### State Management

- Local state first (`useState`)
- Lift state only when needed
- Server state via TanStack Query (or similar)
- Avoid global state unless truly global

```jsx
// DO: Colocate state with usage
function SeatMap({ flightId, onSelect }) {
  const { data: seats } = useQuery(['seats', flightId], () => fetchSeats(flightId));
  const [selected, setSelected] = useState([]);
  // ...
}

// DON'T: Global state for local concerns
const globalSeatSelection = createGlobalState([]); // Unnecessary
```

### Hooks

- Custom hooks for reusable logic
- Prefix with `use`
- Single responsibility

```jsx
// DO: Focused, reusable hook
function useOrderStatus(orderId) {
  return useQuery(
    ['order', orderId, 'status'],
    () => fetchOrderStatus(orderId),
    { refetchInterval: 2000 }
  );
}

// Usage
function OrderStatus({ orderId }) {
  const { data, isLoading } = useOrderStatus(orderId);
  if (isLoading) return <Spinner />;
  return <StatusBadge status={data.status} />;
}
```

### API Calls

Centralize in `api/` directory:

```js
// api/orders.js
const API_BASE = '/api';

export async function createOrder(flightId, seats) {
  const res = await fetch(`${API_BASE}/orders`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ flightId, seats }),
  });
  if (!res.ok) throw new Error(await res.text());
  return res.json();
}
```

### Styling

- Keep it simple - plain CSS or CSS modules
- No complex styling libraries for MVP
- Consistent naming: `.component-name`, `.component-name__element`

---

## General Practices

### Git Commits

- Small, focused commits
- Present tense: "Add seat reservation workflow"
- Reference issue numbers when applicable

### Code Review Checklist

- [ ] Does it solve the problem simply?
- [ ] Can I understand it in 30 seconds?
- [ ] Are errors handled and logged?
- [ ] No dead code or commented-out blocks?
- [ ] No premature optimization?

### What NOT to Do

- Don't add features not in the PRD
- Don't create abstractions for single-use code
- Don't add comments that repeat the code
- Don't optimize before measuring
- Don't add dependencies without clear justification

### Dependencies

Before adding a dependency, ask:
1. Can I do this in <50 lines of code?
2. Is this a core, well-maintained library?
3. Does it solve a real problem I have now?

---

## File Templates

### Go Service

```go
package service

type FlightService struct {
    repo *repository.FlightRepo
}

func NewFlightService(repo *repository.FlightRepo) *FlightService {
    return &FlightService{repo: repo}
}

func (s *FlightService) GetByID(ctx context.Context, id string) (*domain.Flight, error) {
    return s.repo.FindByID(ctx, id)
}
```

### React Component

```jsx
import { useState } from 'react';
import './SeatMap.css';

export function SeatMap({ seats, onSelectionChange }) {
  const [selected, setSelected] = useState([]);

  const handleSeatClick = (seatId) => {
    const newSelection = selected.includes(seatId)
      ? selected.filter(s => s !== seatId)
      : [...selected, seatId];
    setSelected(newSelection);
    onSelectionChange(newSelection);
  };

  return (
    <div className="seat-map">
      {seats.map(seat => (
        <button
          key={seat.id}
          className={`seat ${selected.includes(seat.id) ? 'selected' : ''}`}
          onClick={() => handleSeatClick(seat.id)}
          disabled={seat.status !== 'available'}
        >
          {seat.id}
        </button>
      ))}
    </div>
  );
}
```

---

## Quick Reference

| Situation | Approach |
|-----------|----------|
| Need new endpoint | Handler (thin) → Service (logic) → Repo (data) |
| Need workflow state | Use Temporal queries, not external DB |
| Need real-time updates | Polling first, WebSocket if polling insufficient |
| Component too big | Extract smaller components, lift shared state |
| Repeated code (2x) | Leave it |
| Repeated code (3x+) | Consider extracting |
| Adding dependency | Justify it, prefer stdlib |
