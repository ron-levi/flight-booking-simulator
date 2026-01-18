package activities

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/flight-booking-system/internal/config"
	"github.com/flight-booking-system/internal/repository"
)

// BookingActivities contains all activities for the booking workflow
type BookingActivities struct {
	orderRepo    *repository.OrderRepo
	flightRepo   *repository.FlightRepo
	seatLockRepo *repository.SeatLockRepo
	cfg          *config.BookingConfig
}

// NewBookingActivities creates a new BookingActivities instance
func NewBookingActivities(
	pool *pgxpool.Pool,
	redisClient *redis.Client,
	cfg *config.BookingConfig,
) *BookingActivities {
	return &BookingActivities{
		orderRepo:    repository.NewOrderRepo(pool),
		flightRepo:   repository.NewFlightRepo(pool),
		seatLockRepo: repository.NewSeatLockRepo(redisClient),
		cfg:          cfg,
	}
}
