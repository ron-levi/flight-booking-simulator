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
      <div className="flex gap-4 text-sm justify-center flex-wrap">
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
      <div className="bg-gray-100 p-6 rounded-lg overflow-x-auto">
        {/* Front of plane indicator */}
        <div className="text-center text-sm text-gray-500 mb-4 pb-4 border-b border-gray-300">
          &#9992; Front of Plane
        </div>

        <div className="space-y-2 min-w-fit">
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
