import { useState, useCallback, useEffect } from 'react';
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
  // Note: Allows empty array to release all seats and reset timer
  const handleUpdateSeats = () => {
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
  useEffect(() => {
    if (orderStatus && isTerminalStatus(orderStatus.status) && bookingPhase !== 'complete') {
      setBookingPhase('complete');
    }
  }, [orderStatus, bookingPhase]);

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
          <div className="flex flex-col sm:flex-row justify-between items-start gap-4">
            <div>
              <h1 className="text-2xl font-bold">
                {flight?.flightNumber}: {flight?.origin} &rarr; {flight?.destination}
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
                  title={selectedSeats.length === 0 ? 'Release all seats and reset timer' : 'Update seat selection and reset timer'}
                >
                  {updateSeatsMutation.isPending ? 'Updating...' : 'Update Seats'}
                </button>
                <p className="text-xs text-gray-500 text-center -mt-2">
                  Updating seats resets the timer to 15 minutes
                </p>

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
                <span className="font-medium">Route:</span> {flight?.origin} &rarr; {flight?.destination}
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
