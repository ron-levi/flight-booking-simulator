package temporal

import (
	"errors"

	"go.temporal.io/sdk/temporal"
)

// Workflow-level errors
var (
	// ErrReservationExpired indicates the 15-minute seat hold timer expired
	ErrReservationExpired = errors.New("seat reservation expired")

	// ErrPaymentTimeout indicates payment validation timed out
	ErrPaymentTimeout = errors.New("payment validation timed out")

	// ErrWorkflowCanceled indicates the workflow was canceled by user
	ErrWorkflowCanceled = errors.New("booking workflow canceled")
)

// Non-retryable error types for Temporal retry policy
const (
	ErrTypeSeatUnavailable    = "SEAT_UNAVAILABLE"
	ErrTypePaymentDeclined    = "PAYMENT_DECLINED"
	ErrTypeInvalidPaymentCode = "INVALID_PAYMENT_CODE"
	ErrTypeOrderExpired       = "ORDER_EXPIRED"
)

// NewSeatUnavailableError creates a non-retryable seat error
func NewSeatUnavailableError(seatID string) error {
	return temporal.NewApplicationErrorWithCause(
		"seat "+seatID+" is not available",
		ErrTypeSeatUnavailable,
		nil,
	)
}

// NewPaymentDeclinedError creates a non-retryable payment error
func NewPaymentDeclinedError(reason string) error {
	return temporal.NewApplicationErrorWithCause(
		reason,
		ErrTypePaymentDeclined,
		nil,
	)
}

// NewInvalidPaymentCodeError creates a non-retryable validation error
func NewInvalidPaymentCodeError() error {
	return temporal.NewApplicationErrorWithCause(
		"payment code must be 5 digits",
		ErrTypeInvalidPaymentCode,
		nil,
	)
}
