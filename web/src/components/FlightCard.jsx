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
            {flight.origin} &rarr; {flight.destination}
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
