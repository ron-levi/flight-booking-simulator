package domain

import "time"

// Flight represents a flight in the system
type Flight struct {
	ID             string    `json:"id"`
	FlightNumber   string    `json:"flightNumber"`
	Origin         string    `json:"origin"`
	Destination    string    `json:"destination"`
	DepartureTime  time.Time `json:"departureTime"`
	ArrivalTime    time.Time `json:"arrivalTime"`
	TotalSeats     int       `json:"totalSeats"`
	AvailableSeats int       `json:"availableSeats"`
	PriceCents     int64     `json:"priceCents"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

// FlightWithSeats represents a flight with its seat map
type FlightWithSeats struct {
	Flight
	SeatMap SeatMap `json:"seatMap"`
}

// SeatMap represents the seat configuration of a flight
type SeatMap struct {
	Rows        int    `json:"rows"`
	SeatsPerRow int    `json:"seatsPerRow"`
	Seats       []Seat `json:"seats"`
}
