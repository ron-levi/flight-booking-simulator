# Feature: Phase 4 - React Frontend Implementation

The following plan should be complete, but it's important that you validate documentation and codebase patterns and task sanity before you start implementing.

Pay special attention to naming of existing utils types and models. Import from the right files etc.

## Feature Description

Build a minimal React frontend for the flight booking system that enables users to:
- View available flights with seat counts
- Select specific seats from a visual seat map
- See a real-time countdown timer for their seat reservation (15 minutes)
- Enter a 5-digit payment code to complete booking
- See payment retry status and booking confirmation/failure messages

This is Phase 4 of the PRD, completing the MVP with a functional user interface.

## User Story

As a flight booking customer,
I want to browse flights, select seats, and complete my booking through a web interface,
So that I can book flights with real-time feedback on my reservation timer and payment status.

## Problem Statement

The backend is fully implemented with REST API endpoints and Temporal workflows, but there is no frontend for users to interact with the system. The frontend needs to:
- Integrate with 6 API endpoints (flights, orders, seats, status, payment)
- Display real-time countdown timer synced with server expiration time
- Handle polling for order status updates during payment processing
- Provide clear visual feedback for booking states and errors

## Solution Statement

Create a React application using Vite + TanStack Query with the following:
1. **Pages**: FlightListPage, BookingPage (seat selection + payment)
2. **Components**: FlightCard, SeatMap, Timer, PaymentForm, OrderStatus
3. **API Layer**: Typed API client functions for all endpoints
4. **Hooks**: useOrderStatus (polling), useCountdown (timer)
5. **Styling**: Tailwind CSS for responsive, minimal UI

## Feature Metadata

**Feature Type**: New Capability (Frontend Implementation)
**Estimated Complexity**: Medium
**Primary Systems Affected**: Frontend (new), API integration
**Dependencies**: React 18, Vite 5, TanStack Query 5, React Router 6, Tailwind CSS

---

## CONTEXT REFERENCES

### Relevant Codebase Files - IMPORTANT: YOU MUST READ THESE FILES BEFORE IMPLEMENTING!

**Backend API Types** (understand request/response shapes):
- `internal/api/types.go` (lines 1-94) - All API request/response DTOs
- `internal/api/errors.go` (lines 1-71) - Error codes and response format
- `internal/api/routes.go` (lines 50-69) - Route definitions and URL patterns

**Domain Models** (understand data structures):
- `internal/domain/order.go` (lines 6-16) - OrderStatus enum values
- `internal/domain/seat.go` (lines 8-12) - SeatStatus enum values

**PRD Requirements**:
- `PRD.md` (lines 210-227) - Frontend directory structure
- `PRD.md` (lines 375-404) - Frontend polling implementation example
- `PRD.md` (lines 427-440) - Frontend technology versions
- `PRD.md` (lines 702-720) - Phase 4 deliverables

**Best Practices Reference**:
- `.claude/reference/react-frontend-best-practices.md` - Full React patterns guide

### New Files to Create

```
web/
├── index.html              # HTML entry point
├── package.json            # Dependencies and scripts
├── vite.config.js          # Vite configuration with proxy
├── tailwind.config.js      # Tailwind CSS configuration
├── postcss.config.js       # PostCSS for Tailwind
├── src/
│   ├── main.jsx            # React entry point
│   ├── App.jsx             # App component with routing
│   ├── index.css           # Global styles with Tailwind
│   ├── api/
│   │   └── client.js       # API client functions
│   ├── hooks/
│   │   ├── useOrderStatus.js   # Polling hook for order status
│   │   └── useCountdown.js     # Countdown timer hook
│   ├── components/
│   │   ├── Layout.jsx          # Page layout wrapper
│   │   ├── FlightCard.jsx      # Flight list item
│   │   ├── SeatMap.jsx         # Interactive seat grid
│   │   ├── Timer.jsx           # Countdown display
│   │   ├── PaymentForm.jsx     # Payment code input
│   │   ├── OrderStatus.jsx     # Status indicator
│   │   └── LoadingSpinner.jsx  # Loading indicator
│   └── pages/
│       ├── FlightListPage.jsx  # Flight listing page
│       └── BookingPage.jsx     # Seat selection + payment page
```

### Relevant Documentation - YOU SHOULD READ THESE BEFORE IMPLEMENTING!

- [React 18 Documentation](https://react.dev/learn)
  - State management with hooks
  - Why: Core React patterns
- [TanStack Query v5](https://tanstack.com/query/v5/docs/react/overview)
  - useQuery for data fetching
  - useMutation for POST/PUT requests
  - refetchInterval for polling
  - Why: Server state management
- [React Router v6](https://reactrouter.com/en/main)
  - BrowserRouter, Routes, Route
  - useParams, useNavigate
  - Why: Client-side routing
- [Vite Configuration](https://vitejs.dev/config/)
  - Server proxy for API calls
  - Why: Development server setup
- [Tailwind CSS](https://tailwindcss.com/docs)
  - Utility classes
  - Why: Styling approach

### Patterns to Follow

**File Naming:**
- Components: `PascalCase.jsx` (e.g., `FlightCard.jsx`)
- Hooks: `camelCase.js` with `use` prefix (e.g., `useCountdown.js`)
- API: `camelCase.js` (e.g., `client.js`)

**Component Pattern:**
```jsx
// Functional component with destructured props
function FlightCard({ flight, onSelect }) {
  return (
    <div className="p-4 border rounded">
      <h3>{flight.flightNumber}</h3>
      <button onClick={() => onSelect(flight.id)}>Select</button>
    </div>
  );
}

export default FlightCard;
```

**TanStack Query Pattern:**
```jsx
// useQuery for GET requests
const { data, isLoading, error } = useQuery({
  queryKey: ['flights'],
  queryFn: fetchFlights,
});

// useMutation for POST/PUT
const { mutate, isPending } = useMutation({
  mutationFn: createOrder,
  onSuccess: (data) => { /* handle success */ },
});
```

**API Client Pattern:**
```javascript
const API_BASE = '/api';

export async function fetchFlights() {
  const res = await fetch(`${API_BASE}/flights`);
  if (!res.ok) throw new Error('Failed to fetch flights');
  return res.json();
}
```

**Error Response Format (from backend):**
```json
{
  "error": "ERROR_CODE",
  "message": "Human readable message"
}
```

---

## IMPLEMENTATION PLAN

### Phase 1: Project Setup

Initialize Vite + React project with all dependencies and configuration.

**Tasks:**
- Create package.json with all dependencies
- Configure Vite with API proxy
- Set up Tailwind CSS
- Create entry points (index.html, main.jsx)

### Phase 2: API Layer

Create typed API client functions matching backend endpoints.

**Tasks:**
- Create API client with fetch wrapper
- Implement all endpoint functions
- Handle error responses

### Phase 3: Custom Hooks

Create reusable hooks for order status polling and countdown timer.

**Tasks:**
- useOrderStatus hook with TanStack Query polling
- useCountdown hook for timer display

### Phase 4: UI Components

Build all presentational and interactive components.

**Tasks:**
- Layout component for consistent page structure
- FlightCard for flight list items
- SeatMap for interactive seat selection
- Timer for countdown display
- PaymentForm for payment code entry
- OrderStatus for booking state display
- LoadingSpinner for loading states

### Phase 5: Pages

Create page-level components with routing.

**Tasks:**
- FlightListPage with flight query
- BookingPage with full booking flow
- App.jsx with routing setup

### Phase 6: Integration & Polish

Connect everything and add finishing touches.

**Tasks:**
- Wire up all components
- Add error handling
- Add responsive design
- Test full flow

---

## STEP-BY-STEP TASKS

IMPORTANT: Execute every task in order, top to bottom. Each task is atomic and independently testable.

---

### Task 1: CREATE `web/package.json`

Initialize project with all dependencies.

**IMPLEMENT:**
```json
{
  "name": "flight-booking-web",
  "private": true,
  "version": "0.0.1",
  "type": "module",
  "scripts": {
    "dev": "vite",
    "build": "vite build",
    "preview": "vite preview",
    "lint": "eslint . --ext js,jsx --report-unused-disable-directives --max-warnings 0"
  },
  "dependencies": {
    "@tanstack/react-query": "^5.17.0",
    "react": "^18.2.0",
    "react-dom": "^18.2.0",
    "react-router-dom": "^6.21.0"
  },
  "devDependencies": {
    "@types/react": "^18.2.0",
    "@types/react-dom": "^18.2.0",
    "@vitejs/plugin-react": "^4.2.1",
    "autoprefixer": "^10.4.16",
    "eslint": "^8.55.0",
    "eslint-plugin-react": "^7.33.2",
    "eslint-plugin-react-hooks": "^4.6.0",
    "eslint-plugin-react-refresh": "^0.4.5",
    "postcss": "^8.4.32",
    "tailwindcss": "^3.4.0",
    "vite": "^5.0.10"
  }
}
```

**VALIDATE:** `cd web && cat package.json`

---

### Task 2: CREATE `web/vite.config.js`

Configure Vite with API proxy for development.

**IMPLEMENT:**
```javascript
import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';

export default defineConfig({
  plugins: [react()],
  server: {
    port: 3000,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      '/health': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
});
```

**VALIDATE:** `cat web/vite.config.js`

---

### Task 3: CREATE `web/tailwind.config.js`

Configure Tailwind CSS.

**IMPLEMENT:**
```javascript
/** @type {import('tailwindcss').Config} */
export default {
  content: [
    './index.html',
    './src/**/*.{js,jsx}',
  ],
  theme: {
    extend: {
      colors: {
        primary: '#3B82F6',
        success: '#10B981',
        warning: '#F59E0B',
        danger: '#EF4444',
      },
    },
  },
  plugins: [],
};
```

**VALIDATE:** `cat web/tailwind.config.js`

---

### Task 4: CREATE `web/postcss.config.js`

Configure PostCSS for Tailwind.

**IMPLEMENT:**
```javascript
export default {
  plugins: {
    tailwindcss: {},
    autoprefixer: {},
  },
};
```

**VALIDATE:** `cat web/postcss.config.js`

---

### Task 5: CREATE `web/index.html`

HTML entry point.

**IMPLEMENT:**
```html
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>Flight Booking System</title>
  </head>
  <body>
    <div id="root"></div>
    <script type="module" src="/src/main.jsx"></script>
  </body>
</html>
```

**VALIDATE:** `cat web/index.html`

---

### Task 6: CREATE `web/src/index.css`

Global styles with Tailwind directives.

**IMPLEMENT:**
```css
@tailwind base;
@tailwind components;
@tailwind utilities;

body {
  @apply bg-gray-50 text-gray-900;
}

/* Custom utility classes */
@layer components {
  .btn {
    @apply px-4 py-2 rounded font-medium transition-colors focus:outline-none focus:ring-2 focus:ring-offset-2;
  }

  .btn-primary {
    @apply btn bg-primary text-white hover:bg-blue-600 focus:ring-primary;
  }

  .btn-secondary {
    @apply btn bg-gray-200 text-gray-800 hover:bg-gray-300 focus:ring-gray-400;
  }

  .btn-success {
    @apply btn bg-success text-white hover:bg-green-600 focus:ring-success;
  }

  .btn-danger {
    @apply btn bg-danger text-white hover:bg-red-600 focus:ring-danger;
  }

  .card {
    @apply bg-white rounded-lg shadow-md p-6;
  }

  .input {
    @apply w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent;
  }
}
```

**VALIDATE:** `cat web/src/index.css`

---

### Task 7: CREATE `web/src/api/client.js`

API client with all endpoint functions.

**IMPLEMENT:**
```javascript
const API_BASE = '/api';

/**
 * Generic request helper
 */
async function request(endpoint, options = {}) {
  const response = await fetch(`${API_BASE}${endpoint}`, {
    headers: {
      'Content-Type': 'application/json',
      ...options.headers,
    },
    ...options,
  });

  // Handle non-JSON responses
  if (response.status === 204) {
    return null;
  }

  const data = await response.json();

  if (!response.ok) {
    const error = new Error(data.message || 'An error occurred');
    error.code = data.error;
    error.status = response.status;
    throw error;
  }

  return data;
}

// ============ Flight Endpoints ============

/**
 * List all available flights
 * GET /api/flights
 */
export async function fetchFlights() {
  const data = await request('/flights');
  return data.flights;
}

/**
 * Get flight details with seat map
 * GET /api/flights/{flightId}
 */
export async function fetchFlightDetails(flightId) {
  return request(`/flights/${flightId}`);
}

// ============ Order Endpoints ============

/**
 * Create a new booking order
 * POST /api/orders
 * @param {Object} params - { flightId: string, seats: string[] }
 */
export async function createOrder({ flightId, seats }) {
  return request('/orders', {
    method: 'POST',
    body: JSON.stringify({ flightId, seats }),
  });
}

/**
 * Update seat selection for an order
 * PUT /api/orders/{orderId}/seats
 * @param {Object} params - { orderId: string, seats: string[] }
 */
export async function updateSeats({ orderId, seats }) {
  return request(`/orders/${orderId}/seats`, {
    method: 'PUT',
    body: JSON.stringify({ seats }),
  });
}

/**
 * Get order status (for polling)
 * GET /api/orders/{orderId}/status
 */
export async function fetchOrderStatus(orderId) {
  return request(`/orders/${orderId}/status`);
}

/**
 * Submit payment for an order
 * POST /api/orders/{orderId}/pay
 * @param {Object} params - { orderId: string, paymentCode: string }
 */
export async function submitPayment({ orderId, paymentCode }) {
  return request(`/orders/${orderId}/pay`, {
    method: 'POST',
    body: JSON.stringify({ paymentCode }),
  });
}

/**
 * Cancel an order
 * DELETE /api/orders/{orderId}
 */
export async function cancelOrder(orderId) {
  return request(`/orders/${orderId}`, {
    method: 'DELETE',
  });
}

// ============ Type Definitions (for reference) ============

/**
 * @typedef {Object} Flight
 * @property {string} id
 * @property {string} flightNumber
 * @property {string} origin
 * @property {string} destination
 * @property {string} departureTime
 * @property {number} totalSeats
 * @property {number} availableSeats
 * @property {number} priceCents
 */

/**
 * @typedef {Object} FlightDetail
 * @property {string} id
 * @property {string} flightNumber
 * @property {string} origin
 * @property {string} destination
 * @property {string} departureTime
 * @property {number} totalSeats
 * @property {number} availableSeats
 * @property {number} priceCents
 * @property {SeatMap} seatMap
 */

/**
 * @typedef {Object} SeatMap
 * @property {number} rows
 * @property {number} seatsPerRow
 * @property {Seat[]} seats
 */

/**
 * @typedef {Object} Seat
 * @property {string} id
 * @property {number} row
 * @property {string} column
 * @property {string} status - 'available' | 'reserved' | 'booked'
 */

/**
 * @typedef {Object} OrderStatus
 * @property {string} orderId
 * @property {string} status - Order status enum
 * @property {string[]} seats
 * @property {number} timerRemaining - Seconds remaining
 * @property {number} paymentAttempts
 * @property {string} lastError
 */
```

**VALIDATE:** `cat web/src/api/client.js`

---

### Task 8: CREATE `web/src/hooks/useCountdown.js`

Countdown timer hook.

**IMPLEMENT:**
```javascript
import { useState, useEffect, useCallback } from 'react';

/**
 * Custom hook for countdown timer
 * @param {number} initialSeconds - Starting seconds
 * @param {Function} onExpire - Callback when timer reaches 0
 * @returns {Object} - { seconds, minutes, formatted, isExpired, reset }
 */
export function useCountdown(initialSeconds, onExpire) {
  const [seconds, setSeconds] = useState(initialSeconds);

  // Reset function to update timer from server
  const reset = useCallback((newSeconds) => {
    setSeconds(newSeconds);
  }, []);

  useEffect(() => {
    // Update when initial value changes
    if (initialSeconds > 0) {
      setSeconds(initialSeconds);
    }
  }, [initialSeconds]);

  useEffect(() => {
    if (seconds <= 0) {
      onExpire?.();
      return;
    }

    const interval = setInterval(() => {
      setSeconds((prev) => {
        if (prev <= 1) {
          onExpire?.();
          return 0;
        }
        return prev - 1;
      });
    }, 1000);

    return () => clearInterval(interval);
  }, [seconds, onExpire]);

  const minutes = Math.floor(seconds / 60);
  const remainingSeconds = seconds % 60;
  const formatted = `${minutes}:${remainingSeconds.toString().padStart(2, '0')}`;
  const isExpired = seconds <= 0;

  return {
    seconds,
    minutes,
    remainingSeconds,
    formatted,
    isExpired,
    reset,
  };
}

export default useCountdown;
```

**VALIDATE:** `cat web/src/hooks/useCountdown.js`

---

### Task 9: CREATE `web/src/hooks/useOrderStatus.js`

Order status polling hook using TanStack Query.

**IMPLEMENT:**
```javascript
import { useQuery } from '@tanstack/react-query';
import { fetchOrderStatus } from '../api/client';

/**
 * Custom hook for polling order status
 * @param {string} orderId - The order ID to poll
 * @param {Object} options - Hook options
 * @param {boolean} options.enabled - Whether polling is enabled
 * @param {number} options.refetchInterval - Polling interval in ms (default: 2000)
 * @returns {Object} - TanStack Query result with status data
 */
export function useOrderStatus(orderId, options = {}) {
  const {
    enabled = true,
    refetchInterval = 2000,
  } = options;

  return useQuery({
    queryKey: ['orderStatus', orderId],
    queryFn: () => fetchOrderStatus(orderId),
    enabled: enabled && !!orderId,
    refetchInterval: (query) => {
      // Stop polling if order is in terminal state
      const status = query.state.data?.status;
      if (status === 'CONFIRMED' || status === 'FAILED' || status === 'EXPIRED') {
        return false;
      }
      return refetchInterval;
    },
    staleTime: 0, // Always fetch fresh data
  });
}

/**
 * Check if order status is terminal (no more changes expected)
 */
export function isTerminalStatus(status) {
  return status === 'CONFIRMED' || status === 'FAILED' || status === 'EXPIRED';
}

/**
 * Get human-readable status message
 */
export function getStatusMessage(status) {
  const messages = {
    'CREATED': 'Creating order...',
    'SEATS_RESERVED': 'Seats reserved',
    'PAYMENT_PENDING': 'Awaiting payment',
    'PAYMENT_PROCESSING': 'Processing payment...',
    'CONFIRMED': 'Booking confirmed!',
    'FAILED': 'Booking failed',
    'EXPIRED': 'Reservation expired',
  };
  return messages[status] || status;
}

/**
 * Get status color class
 */
export function getStatusColor(status) {
  const colors = {
    'CREATED': 'text-gray-500',
    'SEATS_RESERVED': 'text-blue-500',
    'PAYMENT_PENDING': 'text-yellow-500',
    'PAYMENT_PROCESSING': 'text-yellow-500',
    'CONFIRMED': 'text-green-500',
    'FAILED': 'text-red-500',
    'EXPIRED': 'text-red-500',
  };
  return colors[status] || 'text-gray-500';
}

export default useOrderStatus;
```

**VALIDATE:** `cat web/src/hooks/useOrderStatus.js`

---

### Task 10: CREATE `web/src/components/LoadingSpinner.jsx`

Loading indicator component.

**IMPLEMENT:**
```jsx
function LoadingSpinner({ size = 'md', className = '' }) {
  const sizeClasses = {
    sm: 'h-4 w-4',
    md: 'h-8 w-8',
    lg: 'h-12 w-12',
  };

  return (
    <div className={`flex justify-center items-center ${className}`}>
      <div
        className={`${sizeClasses[size]} animate-spin rounded-full border-2 border-gray-300 border-t-primary`}
      />
    </div>
  );
}

export default LoadingSpinner;
```

**VALIDATE:** `cat web/src/components/LoadingSpinner.jsx`

---

### Task 11: CREATE `web/src/components/Layout.jsx`

Page layout wrapper component.

**IMPLEMENT:**
```jsx
import { Link } from 'react-router-dom';

function Layout({ children }) {
  return (
    <div className="min-h-screen flex flex-col">
      {/* Header */}
      <header className="bg-white shadow-sm">
        <div className="max-w-6xl mx-auto px-4 py-4">
          <Link to="/" className="text-xl font-bold text-primary">
            Flight Booking
          </Link>
        </div>
      </header>

      {/* Main Content */}
      <main className="flex-1 max-w-6xl mx-auto w-full px-4 py-8">
        {children}
      </main>

      {/* Footer */}
      <footer className="bg-gray-100 border-t">
        <div className="max-w-6xl mx-auto px-4 py-4 text-center text-sm text-gray-500">
          Demo application showcasing Temporal workflow patterns
        </div>
      </footer>
    </div>
  );
}

export default Layout;
```

**VALIDATE:** `cat web/src/components/Layout.jsx`

---

### Task 12: CREATE `web/src/components/FlightCard.jsx`

Flight list item component.

**IMPLEMENT:**
```jsx
function FlightCard({ flight, onSelect }) {
  const formatTime = (isoString) => {
    const date = new Date(isoString);
    return date.toLocaleString('en-US', {
      weekday: 'short',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  const formatPrice = (cents) => {
    return `$${(cents / 100).toFixed(2)}`;
  };

  return (
    <div className="card hover:shadow-lg transition-shadow">
      <div className="flex justify-between items-start">
        {/* Flight Info */}
        <div className="space-y-2">
          <div className="text-lg font-semibold">
            {flight.flightNumber}
          </div>
          <div className="text-2xl font-bold">
            {flight.origin} → {flight.destination}
          </div>
          <div className="text-gray-500">
            {formatTime(flight.departureTime)}
          </div>
        </div>

        {/* Price and Availability */}
        <div className="text-right space-y-2">
          <div className="text-2xl font-bold text-primary">
            {formatPrice(flight.priceCents)}
          </div>
          <div className="text-sm text-gray-500">
            {flight.availableSeats} of {flight.totalSeats} seats available
          </div>
        </div>
      </div>

      {/* Action Button */}
      <div className="mt-4 pt-4 border-t">
        <button
          onClick={() => onSelect(flight.id)}
          disabled={flight.availableSeats === 0}
          className={`w-full py-2 rounded font-medium transition-colors ${
            flight.availableSeats > 0
              ? 'btn-primary'
              : 'bg-gray-200 text-gray-500 cursor-not-allowed'
          }`}
        >
          {flight.availableSeats > 0 ? 'Select Seats' : 'Sold Out'}
        </button>
      </div>
    </div>
  );
}

export default FlightCard;
```

**VALIDATE:** `cat web/src/components/FlightCard.jsx`

---

### Task 13: CREATE `web/src/components/SeatMap.jsx`

Interactive seat selection grid.

**IMPLEMENT:**
```jsx
import { useMemo } from 'react';

function SeatMap({ seatMap, selectedSeats, onSeatClick, disabled = false }) {
  // Organize seats by row
  const seatsByRow = useMemo(() => {
    const rows = {};
    seatMap.seats.forEach((seat) => {
      if (!rows[seat.row]) {
        rows[seat.row] = [];
      }
      rows[seat.row].push(seat);
    });
    // Sort each row by column
    Object.values(rows).forEach((row) => {
      row.sort((a, b) => a.column.localeCompare(b.column));
    });
    return rows;
  }, [seatMap.seats]);

  const rowNumbers = Object.keys(seatsByRow)
    .map(Number)
    .sort((a, b) => a - b);

  const getSeatClass = (seat) => {
    const isSelected = selectedSeats.includes(seat.id);

    if (seat.status === 'booked') {
      return 'bg-gray-400 text-gray-600 cursor-not-allowed';
    }
    if (seat.status === 'reserved' && !isSelected) {
      return 'bg-yellow-200 text-yellow-800 cursor-not-allowed';
    }
    if (isSelected) {
      return 'bg-primary text-white ring-2 ring-primary ring-offset-2';
    }
    return 'bg-green-100 text-green-800 hover:bg-green-200 cursor-pointer';
  };

  const handleSeatClick = (seat) => {
    if (disabled) return;
    if (seat.status === 'booked') return;
    if (seat.status === 'reserved' && !selectedSeats.includes(seat.id)) return;
    onSeatClick(seat.id);
  };

  return (
    <div className="space-y-4">
      {/* Legend */}
      <div className="flex gap-4 text-sm justify-center">
        <div className="flex items-center gap-2">
          <div className="w-6 h-6 bg-green-100 rounded" />
          <span>Available</span>
        </div>
        <div className="flex items-center gap-2">
          <div className="w-6 h-6 bg-primary rounded" />
          <span>Selected</span>
        </div>
        <div className="flex items-center gap-2">
          <div className="w-6 h-6 bg-yellow-200 rounded" />
          <span>Reserved</span>
        </div>
        <div className="flex items-center gap-2">
          <div className="w-6 h-6 bg-gray-400 rounded" />
          <span>Booked</span>
        </div>
      </div>

      {/* Seat Grid */}
      <div className="bg-gray-100 p-6 rounded-lg">
        {/* Front of plane indicator */}
        <div className="text-center text-sm text-gray-500 mb-4 pb-4 border-b border-gray-300">
          ✈ Front of Plane
        </div>

        <div className="space-y-2">
          {rowNumbers.map((rowNum) => (
            <div key={rowNum} className="flex items-center justify-center gap-1">
              {/* Row number */}
              <div className="w-8 text-sm text-gray-500 text-right pr-2">
                {rowNum}
              </div>

              {/* Seats - split into left and right sections (3 + 3) */}
              <div className="flex gap-1">
                {seatsByRow[rowNum].slice(0, 3).map((seat) => (
                  <button
                    key={seat.id}
                    onClick={() => handleSeatClick(seat)}
                    disabled={disabled || seat.status === 'booked' || (seat.status === 'reserved' && !selectedSeats.includes(seat.id))}
                    className={`w-10 h-10 rounded text-xs font-medium transition-all ${getSeatClass(seat)}`}
                    title={`Seat ${seat.id} - ${seat.status}`}
                  >
                    {seat.column}
                  </button>
                ))}
              </div>

              {/* Aisle */}
              <div className="w-8" />

              {/* Right seats */}
              <div className="flex gap-1">
                {seatsByRow[rowNum].slice(3).map((seat) => (
                  <button
                    key={seat.id}
                    onClick={() => handleSeatClick(seat)}
                    disabled={disabled || seat.status === 'booked' || (seat.status === 'reserved' && !selectedSeats.includes(seat.id))}
                    className={`w-10 h-10 rounded text-xs font-medium transition-all ${getSeatClass(seat)}`}
                    title={`Seat ${seat.id} - ${seat.status}`}
                  >
                    {seat.column}
                  </button>
                ))}
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* Selected seats summary */}
      {selectedSeats.length > 0 && (
        <div className="text-center text-sm">
          Selected: <span className="font-semibold">{selectedSeats.join(', ')}</span>
        </div>
      )}
    </div>
  );
}

export default SeatMap;
```

**VALIDATE:** `cat web/src/components/SeatMap.jsx`

---

### Task 14: CREATE `web/src/components/Timer.jsx`

Countdown timer display component.

**IMPLEMENT:**
```jsx
import useCountdown from '../hooks/useCountdown';

function Timer({ seconds, onExpire, className = '' }) {
  const { formatted, isExpired, minutes } = useCountdown(seconds, onExpire);

  // Determine urgency color
  const getTimerColor = () => {
    if (isExpired) return 'text-red-600 bg-red-50';
    if (minutes < 2) return 'text-red-600 bg-red-50 animate-pulse';
    if (minutes < 5) return 'text-yellow-600 bg-yellow-50';
    return 'text-green-600 bg-green-50';
  };

  return (
    <div className={`inline-flex items-center gap-2 px-4 py-2 rounded-lg ${getTimerColor()} ${className}`}>
      <svg
        className="w-5 h-5"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
      >
        <path
          strokeLinecap="round"
          strokeLinejoin="round"
          strokeWidth={2}
          d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
        />
      </svg>
      <span className="font-mono text-lg font-bold">
        {isExpired ? 'Expired' : formatted}
      </span>
      {!isExpired && (
        <span className="text-sm opacity-75">remaining</span>
      )}
    </div>
  );
}

export default Timer;
```

**VALIDATE:** `cat web/src/components/Timer.jsx`

---

### Task 15: CREATE `web/src/components/PaymentForm.jsx`

Payment code input form.

**IMPLEMENT:**
```jsx
import { useState } from 'react';

function PaymentForm({ onSubmit, isLoading, disabled = false }) {
  const [paymentCode, setPaymentCode] = useState('');
  const [error, setError] = useState('');

  const handleChange = (e) => {
    const value = e.target.value.replace(/\D/g, '').slice(0, 5);
    setPaymentCode(value);
    setError('');
  };

  const handleSubmit = (e) => {
    e.preventDefault();

    if (paymentCode.length !== 5) {
      setError('Payment code must be exactly 5 digits');
      return;
    }

    onSubmit(paymentCode);
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div>
        <label htmlFor="paymentCode" className="block text-sm font-medium text-gray-700 mb-1">
          Payment Code
        </label>
        <input
          id="paymentCode"
          type="text"
          inputMode="numeric"
          pattern="\d*"
          value={paymentCode}
          onChange={handleChange}
          placeholder="Enter 5-digit code"
          disabled={disabled || isLoading}
          className={`input text-center text-2xl tracking-widest font-mono ${
            error ? 'border-red-500 focus:ring-red-500' : ''
          }`}
          autoComplete="off"
        />
        {error && (
          <p className="mt-1 text-sm text-red-600">{error}</p>
        )}
        <p className="mt-2 text-xs text-gray-500">
          Demo mode: Any 5-digit code works (15% simulated failure rate)
        </p>
      </div>

      <button
        type="submit"
        disabled={disabled || isLoading || paymentCode.length !== 5}
        className={`w-full py-3 rounded font-medium transition-colors ${
          disabled || isLoading || paymentCode.length !== 5
            ? 'bg-gray-200 text-gray-500 cursor-not-allowed'
            : 'btn-success'
        }`}
      >
        {isLoading ? (
          <span className="flex items-center justify-center gap-2">
            <span className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin" />
            Processing...
          </span>
        ) : (
          'Pay Now'
        )}
      </button>
    </form>
  );
}

export default PaymentForm;
```

**VALIDATE:** `cat web/src/components/PaymentForm.jsx`

---

### Task 16: CREATE `web/src/components/OrderStatus.jsx`

Order status indicator component.

**IMPLEMENT:**
```jsx
import { getStatusMessage, getStatusColor, isTerminalStatus } from '../hooks/useOrderStatus';
import LoadingSpinner from './LoadingSpinner';

function OrderStatus({ status, paymentAttempts = 0, lastError = '' }) {
  const isTerminal = isTerminalStatus(status);
  const statusMessage = getStatusMessage(status);
  const statusColor = getStatusColor(status);

  const getIcon = () => {
    switch (status) {
      case 'CONFIRMED':
        return (
          <svg className="w-8 h-8 text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
          </svg>
        );
      case 'FAILED':
      case 'EXPIRED':
        return (
          <svg className="w-8 h-8 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
          </svg>
        );
      case 'PAYMENT_PROCESSING':
        return <LoadingSpinner size="md" />;
      default:
        return null;
    }
  };

  return (
    <div className={`p-4 rounded-lg border ${
      status === 'CONFIRMED' ? 'bg-green-50 border-green-200' :
      status === 'FAILED' || status === 'EXPIRED' ? 'bg-red-50 border-red-200' :
      status === 'PAYMENT_PROCESSING' ? 'bg-yellow-50 border-yellow-200' :
      'bg-gray-50 border-gray-200'
    }`}>
      <div className="flex items-center gap-3">
        {getIcon()}
        <div>
          <div className={`font-semibold ${statusColor}`}>
            {statusMessage}
          </div>

          {/* Payment attempts indicator */}
          {status === 'PAYMENT_PROCESSING' && paymentAttempts > 0 && (
            <div className="text-sm text-gray-600">
              Attempt {paymentAttempts} of 3
              {paymentAttempts > 1 && ' (retrying...)'}
            </div>
          )}

          {/* Error message */}
          {lastError && (status === 'FAILED' || status === 'EXPIRED') && (
            <div className="text-sm text-red-600 mt-1">
              {lastError}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

export default OrderStatus;
```

**VALIDATE:** `cat web/src/components/OrderStatus.jsx`

---

### Task 17: CREATE `web/src/pages/FlightListPage.jsx`

Flight listing page with query.

**IMPLEMENT:**
```jsx
import { useQuery } from '@tanstack/react-query';
import { useNavigate } from 'react-router-dom';
import { fetchFlights } from '../api/client';
import FlightCard from '../components/FlightCard';
import LoadingSpinner from '../components/LoadingSpinner';
import Layout from '../components/Layout';

function FlightListPage() {
  const navigate = useNavigate();

  const { data: flights, isLoading, error } = useQuery({
    queryKey: ['flights'],
    queryFn: fetchFlights,
  });

  const handleSelectFlight = (flightId) => {
    navigate(`/book/${flightId}`);
  };

  if (isLoading) {
    return (
      <Layout>
        <div className="flex justify-center items-center h-64">
          <LoadingSpinner size="lg" />
        </div>
      </Layout>
    );
  }

  if (error) {
    return (
      <Layout>
        <div className="card text-center">
          <div className="text-red-500 mb-4">
            Failed to load flights: {error.message}
          </div>
          <button
            onClick={() => window.location.reload()}
            className="btn-primary"
          >
            Retry
          </button>
        </div>
      </Layout>
    );
  }

  return (
    <Layout>
      <div className="space-y-6">
        <div>
          <h1 className="text-3xl font-bold">Available Flights</h1>
          <p className="text-gray-500 mt-1">
            Select a flight to view seats and make a booking
          </p>
        </div>

        {flights?.length === 0 ? (
          <div className="card text-center text-gray-500">
            No flights available at this time.
          </div>
        ) : (
          <div className="grid gap-4 md:grid-cols-2">
            {flights?.map((flight) => (
              <FlightCard
                key={flight.id}
                flight={flight}
                onSelect={handleSelectFlight}
              />
            ))}
          </div>
        )}
      </div>
    </Layout>
  );
}

export default FlightListPage;
```

**VALIDATE:** `cat web/src/pages/FlightListPage.jsx`

---

### Task 18: CREATE `web/src/pages/BookingPage.jsx`

Complete booking page with seat selection and payment.

**IMPLEMENT:**
```jsx
import { useState, useCallback } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import {
  fetchFlightDetails,
  createOrder,
  updateSeats,
  submitPayment,
  cancelOrder,
} from '../api/client';
import { useOrderStatus, isTerminalStatus } from '../hooks/useOrderStatus';
import Layout from '../components/Layout';
import LoadingSpinner from '../components/LoadingSpinner';
import SeatMap from '../components/SeatMap';
import Timer from '../components/Timer';
import PaymentForm from '../components/PaymentForm';
import OrderStatus from '../components/OrderStatus';

function BookingPage() {
  const { flightId } = useParams();
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  // Local state
  const [selectedSeats, setSelectedSeats] = useState([]);
  const [orderId, setOrderId] = useState(null);
  const [bookingPhase, setBookingPhase] = useState('selecting'); // selecting | reserved | paying | complete
  const [error, setError] = useState(null);

  // Fetch flight details
  const { data: flight, isLoading: flightLoading, error: flightError } = useQuery({
    queryKey: ['flight', flightId],
    queryFn: () => fetchFlightDetails(flightId),
    enabled: !!flightId,
  });

  // Poll order status when we have an order
  const { data: orderStatus } = useOrderStatus(orderId, {
    enabled: !!orderId && bookingPhase !== 'selecting',
    refetchInterval: 2000,
  });

  // Create order mutation
  const createOrderMutation = useMutation({
    mutationFn: createOrder,
    onSuccess: (data) => {
      setOrderId(data.orderId);
      setBookingPhase('reserved');
      setError(null);
    },
    onError: (err) => {
      setError(err.message);
    },
  });

  // Update seats mutation
  const updateSeatsMutation = useMutation({
    mutationFn: updateSeats,
    onSuccess: () => {
      // Invalidate flight query to refresh seat map
      queryClient.invalidateQueries({ queryKey: ['flight', flightId] });
    },
    onError: (err) => {
      setError(err.message);
    },
  });

  // Submit payment mutation
  const paymentMutation = useMutation({
    mutationFn: submitPayment,
    onSuccess: () => {
      setBookingPhase('paying');
    },
    onError: (err) => {
      setError(err.message);
    },
  });

  // Cancel order mutation
  const cancelMutation = useMutation({
    mutationFn: () => cancelOrder(orderId),
    onSuccess: () => {
      setOrderId(null);
      setSelectedSeats([]);
      setBookingPhase('selecting');
      queryClient.invalidateQueries({ queryKey: ['flight', flightId] });
    },
  });

  // Handle seat click
  const handleSeatClick = useCallback((seatId) => {
    setSelectedSeats((prev) => {
      if (prev.includes(seatId)) {
        return prev.filter((id) => id !== seatId);
      }
      return [...prev, seatId];
    });
  }, []);

  // Handle reserve seats
  const handleReserveSeats = () => {
    if (selectedSeats.length === 0) {
      setError('Please select at least one seat');
      return;
    }
    createOrderMutation.mutate({ flightId, seats: selectedSeats });
  };

  // Handle update seats (after reservation)
  const handleUpdateSeats = () => {
    if (selectedSeats.length === 0) {
      setError('Please select at least one seat');
      return;
    }
    updateSeatsMutation.mutate({ orderId, seats: selectedSeats });
  };

  // Handle payment submission
  const handlePayment = (paymentCode) => {
    paymentMutation.mutate({ orderId, paymentCode });
  };

  // Handle timer expiry
  const handleTimerExpire = useCallback(() => {
    if (bookingPhase === 'reserved') {
      setBookingPhase('complete');
    }
  }, [bookingPhase]);

  // Handle back to flights
  const handleBackToFlights = () => {
    navigate('/');
  };

  // Check if order status changed to terminal
  if (orderStatus && isTerminalStatus(orderStatus.status) && bookingPhase !== 'complete') {
    setBookingPhase('complete');
  }

  // Loading state
  if (flightLoading) {
    return (
      <Layout>
        <div className="flex justify-center items-center h-64">
          <LoadingSpinner size="lg" />
        </div>
      </Layout>
    );
  }

  // Error state
  if (flightError) {
    return (
      <Layout>
        <div className="card text-center">
          <div className="text-red-500 mb-4">
            {flightError.message}
          </div>
          <button onClick={handleBackToFlights} className="btn-primary">
            Back to Flights
          </button>
        </div>
      </Layout>
    );
  }

  // Calculate total price
  const totalPrice = selectedSeats.length * (flight?.priceCents || 0);
  const formatPrice = (cents) => `$${(cents / 100).toFixed(2)}`;

  return (
    <Layout>
      <div className="space-y-6">
        {/* Flight Header */}
        <div className="card">
          <div className="flex justify-between items-start">
            <div>
              <h1 className="text-2xl font-bold">
                {flight?.flightNumber}: {flight?.origin} → {flight?.destination}
              </h1>
              <p className="text-gray-500">
                {new Date(flight?.departureTime).toLocaleString()}
              </p>
            </div>
            <div className="text-right">
              <div className="text-xl font-bold text-primary">
                {formatPrice(flight?.priceCents || 0)} / seat
              </div>
            </div>
          </div>
        </div>

        {/* Timer (when reserved) */}
        {bookingPhase === 'reserved' && orderStatus?.timerRemaining > 0 && (
          <div className="flex justify-center">
            <Timer
              seconds={orderStatus.timerRemaining}
              onExpire={handleTimerExpire}
            />
          </div>
        )}

        {/* Order Status (when paying or complete) */}
        {(bookingPhase === 'paying' || bookingPhase === 'complete') && orderStatus && (
          <OrderStatus
            status={orderStatus.status}
            paymentAttempts={orderStatus.paymentAttempts}
            lastError={orderStatus.lastError}
          />
        )}

        {/* Error Message */}
        {error && (
          <div className="p-4 bg-red-50 border border-red-200 rounded-lg text-red-600">
            {error}
          </div>
        )}

        <div className="grid gap-6 lg:grid-cols-3">
          {/* Seat Map - Takes 2 columns */}
          <div className="lg:col-span-2 card">
            <h2 className="text-xl font-semibold mb-4">Select Your Seats</h2>
            {flight?.seatMap && (
              <SeatMap
                seatMap={flight.seatMap}
                selectedSeats={selectedSeats}
                onSeatClick={handleSeatClick}
                disabled={bookingPhase === 'paying' || bookingPhase === 'complete'}
              />
            )}
          </div>

          {/* Booking Panel */}
          <div className="card space-y-4">
            <h2 className="text-xl font-semibold">Booking Summary</h2>

            {/* Selected Seats */}
            <div>
              <div className="text-sm text-gray-500">Selected Seats</div>
              <div className="font-semibold">
                {selectedSeats.length > 0
                  ? selectedSeats.join(', ')
                  : 'None selected'}
              </div>
            </div>

            {/* Total Price */}
            <div>
              <div className="text-sm text-gray-500">Total Price</div>
              <div className="text-2xl font-bold text-primary">
                {formatPrice(totalPrice)}
              </div>
            </div>

            {/* Action Buttons */}
            {bookingPhase === 'selecting' && (
              <button
                onClick={handleReserveSeats}
                disabled={selectedSeats.length === 0 || createOrderMutation.isPending}
                className={`w-full py-3 rounded font-medium ${
                  selectedSeats.length === 0 || createOrderMutation.isPending
                    ? 'bg-gray-200 text-gray-500 cursor-not-allowed'
                    : 'btn-primary'
                }`}
              >
                {createOrderMutation.isPending ? 'Reserving...' : 'Reserve Seats'}
              </button>
            )}

            {bookingPhase === 'reserved' && (
              <>
                <button
                  onClick={handleUpdateSeats}
                  disabled={updateSeatsMutation.isPending}
                  className="w-full btn-secondary"
                >
                  {updateSeatsMutation.isPending ? 'Updating...' : 'Update Seats'}
                </button>

                <div className="border-t pt-4">
                  <h3 className="font-semibold mb-3">Complete Payment</h3>
                  <PaymentForm
                    onSubmit={handlePayment}
                    isLoading={paymentMutation.isPending}
                    disabled={orderStatus?.timerRemaining <= 0}
                  />
                </div>

                <button
                  onClick={() => cancelMutation.mutate()}
                  disabled={cancelMutation.isPending}
                  className="w-full text-red-500 hover:text-red-700 text-sm"
                >
                  Cancel Booking
                </button>
              </>
            )}

            {bookingPhase === 'paying' && orderStatus?.status === 'PAYMENT_PROCESSING' && (
              <div className="text-center text-gray-500">
                Processing payment...
              </div>
            )}

            {bookingPhase === 'complete' && (
              <button
                onClick={handleBackToFlights}
                className="w-full btn-primary"
              >
                {orderStatus?.status === 'CONFIRMED'
                  ? 'Book Another Flight'
                  : 'Back to Flights'}
              </button>
            )}
          </div>
        </div>

        {/* Confirmation Details */}
        {bookingPhase === 'complete' && orderStatus?.status === 'CONFIRMED' && (
          <div className="card bg-green-50 border-green-200">
            <h2 className="text-xl font-semibold text-green-800 mb-4">
              Booking Confirmed!
            </h2>
            <div className="grid gap-2 text-green-700">
              <div>
                <span className="font-medium">Order ID:</span> {orderId}
              </div>
              <div>
                <span className="font-medium">Flight:</span> {flight?.flightNumber}
              </div>
              <div>
                <span className="font-medium">Route:</span> {flight?.origin} → {flight?.destination}
              </div>
              <div>
                <span className="font-medium">Seats:</span> {orderStatus?.seats?.join(', ')}
              </div>
            </div>
          </div>
        )}
      </div>
    </Layout>
  );
}

export default BookingPage;
```

**VALIDATE:** `cat web/src/pages/BookingPage.jsx`

---

### Task 19: CREATE `web/src/App.jsx`

Main app component with routing.

**IMPLEMENT:**
```jsx
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import FlightListPage from './pages/FlightListPage';
import BookingPage from './pages/BookingPage';

// Create query client
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 1000 * 60, // 1 minute
      retry: 1,
      refetchOnWindowFocus: false,
    },
  },
});

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<FlightListPage />} />
          <Route path="/book/:flightId" element={<BookingPage />} />
          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </BrowserRouter>
    </QueryClientProvider>
  );
}

export default App;
```

**VALIDATE:** `cat web/src/App.jsx`

---

### Task 20: CREATE `web/src/main.jsx`

React entry point.

**IMPLEMENT:**
```jsx
import React from 'react';
import ReactDOM from 'react-dom/client';
import App from './App';
import './index.css';

ReactDOM.createRoot(document.getElementById('root')).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
);
```

**VALIDATE:** `cat web/src/main.jsx`

---

### Task 21: UPDATE `Makefile` - Add frontend commands

Add npm commands to the Makefile.

**IMPLEMENT:** Add these targets to the end of the existing Makefile:

```makefile
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
```

**VALIDATE:** `grep -A 10 "Frontend" Makefile`

---

### Task 22: CREATE `web/.gitignore`

Git ignore for frontend.

**IMPLEMENT:**
```
node_modules/
dist/
.env.local
.DS_Store
```

**VALIDATE:** `cat web/.gitignore`

---

### Task 23: INSTALL dependencies and verify build

Install npm packages and verify everything builds.

**IMPLEMENT:**
```bash
cd web && npm install
npm run build
```

**VALIDATE:** `cd web && npm run build` succeeds without errors

---

## TESTING STRATEGY

### Unit Tests

For MVP, manual testing is sufficient. Future tests should cover:

**`web/src/hooks/__tests__/useCountdown.test.js`:**
- Timer decrements correctly
- onExpire callback fires at 0
- reset function works

**`web/src/components/__tests__/SeatMap.test.jsx`:**
- Renders correct number of seats
- Click handlers work
- Disabled states work

### Integration Tests

**Full flow tests:**
1. Can view flight list
2. Can select a flight and see seat map
3. Can select seats and reserve them
4. Can see countdown timer
5. Can enter payment code and submit
6. Can see confirmation or error

### Edge Cases

- Empty flight list
- No available seats on flight
- Timer expiration during booking
- Payment validation failure (retry behavior)
- Network errors during API calls
- Selecting already-reserved seats

---

## VALIDATION COMMANDS

Execute every command to ensure zero regressions and 100% feature correctness.

### Level 1: Build

```bash
# Install dependencies
cd web && npm install

# Build production bundle
npm run build

# Lint code
npm run lint
```

### Level 2: Development Server

```bash
# Start development server
npm run dev

# Verify it starts on port 3000
curl -s http://localhost:3000 | head -5
```

### Level 3: API Proxy

With backend running:
```bash
# Verify API proxy works
curl -s http://localhost:3000/api/flights | jq .
```

### Level 4: Manual Validation

Start the full stack:
```bash
# Terminal 1: Infrastructure
make up && make migrate-up

# Terminal 2: Worker
make run-worker

# Terminal 3: Server
make run-server

# Terminal 4: Frontend
cd web && npm run dev
```

Open http://localhost:3000 and test:
1. ✓ Flight list loads with all flights
2. ✓ Click "Select Seats" navigates to booking page
3. ✓ Seat map displays with correct colors
4. ✓ Can click seats to select/deselect
5. ✓ "Reserve Seats" creates order and starts timer
6. ✓ Timer counts down in real-time
7. ✓ Can update seat selection (timer resets)
8. ✓ Can enter 5-digit payment code
9. ✓ Payment processing shows loading state
10. ✓ Success shows confirmation with details
11. ✓ Failure shows error message
12. ✓ Timer expiration shows expired status

---

## ACCEPTANCE CRITERIA

- [ ] React frontend builds without errors
- [ ] Flight list page displays all available flights
- [ ] Seat map shows available, reserved, and booked seats
- [ ] Seat selection updates local state correctly
- [ ] Timer displays countdown in MM:SS format
- [ ] Timer shows urgency colors (green → yellow → red)
- [ ] Payment form validates 5-digit codes
- [ ] Order status polling works during payment
- [ ] Success confirmation shows order details
- [ ] Failure/expiry states show appropriate messages
- [ ] All pages are responsive (mobile + desktop)
- [ ] CORS works with backend on port 8080
- [ ] Navigation between pages works correctly

---

## COMPLETION CHECKLIST

- [ ] All tasks completed in order
- [ ] `npm install` succeeds
- [ ] `npm run build` succeeds without errors
- [ ] `npm run dev` starts development server
- [ ] API proxy routes to backend correctly
- [ ] Full booking flow works end-to-end
- [ ] Acceptance criteria all met
- [ ] Code follows project conventions

---

## NOTES

### Design Decisions

1. **Vite over CRA**: Vite provides faster development experience with native ES modules and better build performance.

2. **TanStack Query for server state**: Using React Query provides:
   - Automatic caching and deduplication
   - Built-in loading/error states
   - Configurable polling with `refetchInterval`
   - Automatic background refetching

3. **Polling over WebSocket**: For MVP, polling every 2 seconds is simpler to implement and sufficient for the demo. WebSocket can be added in future.

4. **Tailwind CSS**: Utility-first CSS provides:
   - Rapid development
   - Consistent design tokens
   - Small production bundle (purged CSS)
   - No custom CSS to maintain

5. **Component structure**: Kept components focused and single-purpose:
   - Layout: Page wrapper
   - FlightCard: Flight display
   - SeatMap: Seat grid (most complex)
   - Timer: Countdown display
   - PaymentForm: Payment input
   - OrderStatus: Status indicator

### Trade-offs

1. **No TypeScript**: For faster MVP development, using plain JavaScript with JSDoc comments for type hints. TypeScript can be added later.

2. **No form library**: Using native form handling for PaymentForm since it's simple (single input). Would use react-hook-form for more complex forms.

3. **Basic error handling**: Errors are displayed inline. Production would add toast notifications and error boundaries.

4. **No loading skeletons**: Using simple spinners. Skeletons would provide better UX.

### Future Improvements

- Add TypeScript for type safety
- Add React Testing Library tests
- Add toast notifications for better feedback
- Add skeleton loading states
- Add WebSocket for real-time updates
- Add seat pricing tiers (economy, business, first)
- Add booking history page
- Add responsive seat map for mobile
- Add animation transitions
