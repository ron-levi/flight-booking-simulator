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
