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
