package domain

import "errors"

var (
	// ErrFlightNotFound indicates a flight was not found
	ErrFlightNotFound = errors.New("flight not found")

	// ErrOrderNotFound indicates an order was not found
	ErrOrderNotFound = errors.New("order not found")

	// ErrSeatUnavailable indicates a seat is not available for booking
	ErrSeatUnavailable = errors.New("seat is not available")

	// ErrSeatsAlreadyLocked indicates seats are already locked by another order
	ErrSeatsAlreadyLocked = errors.New("seats are already locked")

	// ErrInsufficientSeats indicates not enough seats available
	ErrInsufficientSeats = errors.New("insufficient seats available")

	// ErrOrderExpired indicates the order reservation has expired
	ErrOrderExpired = errors.New("order reservation has expired")

	// ErrInvalidPaymentCode indicates the payment code format is invalid
	ErrInvalidPaymentCode = errors.New("invalid payment code format")

	// ErrPaymentFailed indicates payment validation failed
	ErrPaymentFailed = errors.New("payment validation failed")
)
