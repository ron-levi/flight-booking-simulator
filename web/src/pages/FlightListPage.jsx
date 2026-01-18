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
