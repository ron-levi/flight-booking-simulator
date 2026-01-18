package service

import (
	"context"
	"fmt"

	"go.temporal.io/sdk/client"

	"github.com/flight-booking-system/internal/config"
	temporalpkg "github.com/flight-booking-system/internal/temporal"
	"github.com/flight-booking-system/internal/temporal/workflows"
)

// TemporalClient wraps the Temporal SDK client for booking operations
type TemporalClient struct {
	client    client.Client
	taskQueue string
}

// NewTemporalClient creates a new Temporal client wrapper
func NewTemporalClient(cfg *config.TemporalConfig) (*TemporalClient, error) {
	c, err := client.Dial(client.Options{
		HostPort:  cfg.Host,
		Namespace: cfg.Namespace,
	})
	if err != nil {
		return nil, fmt.Errorf("dial temporal: %w", err)
	}

	return &TemporalClient{
		client:    c,
		taskQueue: cfg.TaskQueue,
	}, nil
}

// Close closes the Temporal client connection
func (tc *TemporalClient) Close() {
	tc.client.Close()
}

// StartBookingWorkflow starts a new booking workflow
func (tc *TemporalClient) StartBookingWorkflow(ctx context.Context, input temporalpkg.BookingWorkflowInput) (string, error) {
	workflowID := fmt.Sprintf("booking-%s", input.OrderID)

	opts := client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: tc.taskQueue,
	}

	run, err := tc.client.ExecuteWorkflow(ctx, opts, workflows.BookingWorkflow, input)
	if err != nil {
		return "", fmt.Errorf("start booking workflow: %w", err)
	}

	return run.GetID(), nil
}

// SignalUpdateSeats sends an update seats signal to a booking workflow
func (tc *TemporalClient) SignalUpdateSeats(ctx context.Context, orderID string, seats []string) error {
	workflowID := fmt.Sprintf("booking-%s", orderID)

	err := tc.client.SignalWorkflow(ctx, workflowID, "", temporalpkg.SignalUpdateSeats, temporalpkg.SeatUpdateSignal{
		Seats: seats,
	})
	if err != nil {
		return fmt.Errorf("signal update seats: %w", err)
	}

	return nil
}

// SignalProceedToPayment sends a proceed to payment signal with the payment code
func (tc *TemporalClient) SignalProceedToPayment(ctx context.Context, orderID string, paymentCode string) error {
	workflowID := fmt.Sprintf("booking-%s", orderID)

	err := tc.client.SignalWorkflow(ctx, workflowID, "", temporalpkg.SignalProceedToPay, temporalpkg.PaymentSignal{
		PaymentCode: paymentCode,
	})
	if err != nil {
		return fmt.Errorf("signal proceed to payment: %w", err)
	}

	return nil
}

// SignalCancelBooking sends a cancel signal to the booking workflow
func (tc *TemporalClient) SignalCancelBooking(ctx context.Context, orderID string) error {
	workflowID := fmt.Sprintf("booking-%s", orderID)

	err := tc.client.SignalWorkflow(ctx, workflowID, "", temporalpkg.SignalCancelBooking, nil)
	if err != nil {
		return fmt.Errorf("signal cancel booking: %w", err)
	}

	return nil
}

// QueryBookingStatus queries the current status of a booking workflow
func (tc *TemporalClient) QueryBookingStatus(ctx context.Context, orderID string) (*temporalpkg.BookingStatusResponse, error) {
	workflowID := fmt.Sprintf("booking-%s", orderID)

	result, err := tc.client.QueryWorkflow(ctx, workflowID, "", temporalpkg.QueryBookingStatus)
	if err != nil {
		return nil, fmt.Errorf("query booking status: %w", err)
	}

	var status temporalpkg.BookingStatusResponse
	if err := result.Get(&status); err != nil {
		return nil, fmt.Errorf("decode query result: %w", err)
	}

	return &status, nil
}
