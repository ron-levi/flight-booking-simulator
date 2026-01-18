package service

import (
	"context"

	"github.com/flight-booking-system/internal/domain"
	"github.com/flight-booking-system/internal/repository"
)

// FlightService handles flight-related business logic
type FlightService struct {
	flightRepo   *repository.FlightRepo
	seatLockRepo *repository.SeatLockRepo
}

// NewFlightService creates a new FlightService
func NewFlightService(flightRepo *repository.FlightRepo, seatLockRepo *repository.SeatLockRepo) *FlightService {
	return &FlightService{
		flightRepo:   flightRepo,
		seatLockRepo: seatLockRepo,
	}
}

// ListFlights returns all available flights
func (s *FlightService) ListFlights(ctx context.Context) ([]domain.Flight, error) {
	return s.flightRepo.FindAll(ctx)
}

// GetFlightWithSeats returns a flight with its seat map and real-time availability
func (s *FlightService) GetFlightWithSeats(ctx context.Context, flightID string) (*domain.FlightWithSeats, error) {
	// Get flight details
	flight, err := s.flightRepo.FindByID(ctx, flightID)
	if err != nil {
		return nil, err
	}

	// Get all seats for the flight
	seats, err := s.flightRepo.FindSeats(ctx, flightID)
	if err != nil {
		return nil, err
	}

	// Get currently locked seats from Redis
	lockedSeats, err := s.seatLockRepo.GetLockedSeats(ctx, flightID)
	if err != nil {
		return nil, err
	}

	// Update seat status based on locks
	for i := range seats {
		if _, isLocked := lockedSeats[seats[i].ID]; isLocked {
			if seats[i].Status == domain.SeatStatusAvailable {
				seats[i].Status = domain.SeatStatusReserved
			}
		}
	}

	// Calculate seat map dimensions
	rows := 0
	seatsPerRow := 0
	if len(seats) > 0 {
		rowMap := make(map[int]int)
		for _, seat := range seats {
			rowMap[seat.Row]++
			if seat.Row > rows {
				rows = seat.Row
			}
		}
		// Get seats per row from first row
		if count, ok := rowMap[1]; ok {
			seatsPerRow = count
		}
	}

	return &domain.FlightWithSeats{
		Flight: *flight,
		SeatMap: domain.SeatMap{
			Rows:        rows,
			SeatsPerRow: seatsPerRow,
			Seats:       seats,
		},
	}, nil
}
