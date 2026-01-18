package api

import "time"

// Request types

// CreateOrderRequest is the request body for creating a new order
type CreateOrderRequest struct {
	FlightID string   `json:"flightId"`
	Seats    []string `json:"seats"`
}

// UpdateSeatsRequest is the request body for updating seat selection
type UpdateSeatsRequest struct {
	Seats []string `json:"seats"`
}

// SubmitPaymentRequest is the request body for submitting payment
type SubmitPaymentRequest struct {
	PaymentCode string `json:"paymentCode"`
}

// Response types

// FlightListResponse contains a list of flights
type FlightListResponse struct {
	Flights []FlightResponse `json:"flights"`
}

// FlightResponse represents a flight in API responses
type FlightResponse struct {
	ID             string    `json:"id"`
	FlightNumber   string    `json:"flightNumber"`
	Origin         string    `json:"origin"`
	Destination    string    `json:"destination"`
	DepartureTime  time.Time `json:"departureTime"`
	TotalSeats     int       `json:"totalSeats"`
	AvailableSeats int       `json:"availableSeats"`
	PriceCents     int64     `json:"priceCents"`
}

// FlightDetailResponse represents a flight with seat map
type FlightDetailResponse struct {
	FlightResponse
	SeatMap SeatMapResponse `json:"seatMap"`
}

// SeatMapResponse represents seat map configuration
type SeatMapResponse struct {
	Rows        int            `json:"rows"`
	SeatsPerRow int            `json:"seatsPerRow"`
	Seats       []SeatResponse `json:"seats"`
}

// SeatResponse represents a seat in API responses
type SeatResponse struct {
	ID     string `json:"id"`
	Row    int    `json:"row"`
	Column string `json:"column"`
	Status string `json:"status"` // "available", "reserved", "booked"
}

// CreateOrderResponse is the response for order creation
type CreateOrderResponse struct {
	OrderID    string    `json:"orderId"`
	WorkflowID string    `json:"workflowId"`
	Status     string    `json:"status"`
	ExpiresAt  time.Time `json:"expiresAt"`
}

// OrderStatusResponse is the response for order status queries
type OrderStatusResponse struct {
	OrderID         string   `json:"orderId"`
	Status          string   `json:"status"`
	Seats           []string `json:"seats"`
	TimerRemaining  int      `json:"timerRemaining"`
	PaymentAttempts int      `json:"paymentAttempts"`
	LastError       string   `json:"lastError,omitempty"`
}

// UpdateSeatsResponse is the response for seat update
type UpdateSeatsResponse struct {
	OrderID   string    `json:"orderId"`
	Status    string    `json:"status"`
	Seats     []string  `json:"seats"`
	ExpiresAt time.Time `json:"expiresAt"`
}

// PaymentAcceptedResponse is the response for payment submission
type PaymentAcceptedResponse struct {
	OrderID string `json:"orderId"`
	Status  string `json:"status"`
}
