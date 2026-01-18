package activities

import (
	"context"
	"fmt"
	"math/rand"
	"regexp"
	"time"

	"go.temporal.io/sdk/temporal"

	temporalpkg "github.com/flight-booking-system/internal/temporal"
)

// ValidatePaymentInput contains payment validation parameters
type ValidatePaymentInput struct {
	OrderID     string
	PaymentCode string
}

// ValidatePaymentOutput contains the validation result
type ValidatePaymentOutput struct {
	Success bool
	Message string
}

// 5-digit code pattern
var paymentCodePattern = regexp.MustCompile(`^\d{5}$`)

// ValidatePayment simulates payment code validation
// - 15% failure rate (configurable via cfg.PaymentFailureRate)
// - Random processing time 1-8 seconds
// - Returns non-retryable error for invalid code format
func (a *BookingActivities) ValidatePayment(ctx context.Context, input ValidatePaymentInput) (ValidatePaymentOutput, error) {
	// Validate payment code format (5 digits)
	if !paymentCodePattern.MatchString(input.PaymentCode) {
		return ValidatePaymentOutput{}, temporalpkg.NewInvalidPaymentCodeError()
	}

	// Special codes for testing
	switch input.PaymentCode {
	case "00000":
		// Always fails - useful for testing
		return ValidatePaymentOutput{}, temporal.NewApplicationError(
			"payment declined: insufficient funds",
			temporalpkg.ErrTypePaymentDeclined,
		)
	case "99999":
		// Always succeeds instantly - useful for testing
		return ValidatePaymentOutput{Success: true, Message: "Payment validated (test mode)"}, nil
	}

	// Simulate processing time (1-8 seconds)
	processingTime := time.Duration(rand.Intn(7)+1) * time.Second
	select {
	case <-time.After(processingTime):
		// Processing complete
	case <-ctx.Done():
		return ValidatePaymentOutput{}, ctx.Err()
	}

	// Simulate failure rate
	if rand.Float64() < a.cfg.PaymentFailureRate {
		// This error IS retryable (will be retried by Temporal)
		return ValidatePaymentOutput{}, fmt.Errorf("payment validation failed: temporary gateway error")
	}

	return ValidatePaymentOutput{
		Success: true,
		Message: "Payment validated successfully",
	}, nil
}
