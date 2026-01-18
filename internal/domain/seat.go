package domain

import "time"

// SeatStatus represents the current status of a seat
type SeatStatus string

const (
	SeatStatusAvailable SeatStatus = "available"
	SeatStatusReserved  SeatStatus = "reserved"
	SeatStatusBooked    SeatStatus = "booked"
)

// Seat represents an individual seat on a flight
type Seat struct {
	ID        string     `json:"id"`
	FlightID  string     `json:"flightId"`
	Row       int        `json:"row"`
	Column    string     `json:"column"`
	Status    SeatStatus `json:"status"`
	OrderID   *string    `json:"orderId,omitempty"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

// SeatID returns the seat identifier (e.g., "12A")
func (s *Seat) SeatID() string {
	return s.ID
}

// IsAvailable checks if the seat can be selected
func (s *Seat) IsAvailable() bool {
	return s.Status == SeatStatusAvailable
}
