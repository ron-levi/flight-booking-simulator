package temporal

import (
	"time"

	"github.com/flight-booking-system/internal/domain"
)

// Signal names as constants
const (
	SignalUpdateSeats   = "update-seats"
	SignalProceedToPay  = "proceed-to-payment"
	SignalCancelBooking = "cancel-booking"
)

// Query names as constants
const (
	QueryBookingStatus = "booking-status"
)

// SeatUpdateSignal is sent when user changes seat selection
type SeatUpdateSignal struct {
	Seats []string `json:"seats"`
}

// PaymentSignal is sent when user submits payment
type PaymentSignal struct {
	PaymentCode string `json:"paymentCode"`
}

// BookingStatusResponse is returned by the status query
type BookingStatusResponse struct {
	OrderID         string             `json:"orderId"`
	FlightID        string             `json:"flightId"`
	Status          domain.OrderStatus `json:"status"`
	Seats           []string           `json:"seats"`
	ExpiresAt       time.Time          `json:"expiresAt"`
	TimerRemaining  int                `json:"timerRemaining"` // seconds
	PaymentAttempts int                `json:"paymentAttempts"`
	LastError       string             `json:"lastError,omitempty"`
}

// BookingWorkflowInput contains the initial workflow parameters
type BookingWorkflowInput struct {
	OrderID  string   `json:"orderId"`
	FlightID string   `json:"flightId"`
	Seats    []string `json:"seats"`
}

// BookingWorkflowResult contains the workflow completion result
type BookingWorkflowResult struct {
	OrderID string             `json:"orderId"`
	Status  domain.OrderStatus `json:"status"`
	Seats   []string           `json:"seats"`
	Error   string             `json:"error,omitempty"`
}
