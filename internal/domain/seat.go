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
