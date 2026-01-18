package service

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/google/uuid"

	"github.com/flight-booking-system/internal/domain"
	"github.com/flight-booking-system/internal/repository"
	temporalpkg "github.com/flight-booking-system/internal/temporal"
)

// BookingService handles booking-related business logic
type BookingService struct {
	orderRepo      *repository.OrderRepo
	flightRepo     *repository.FlightRepo
	temporalClient *TemporalClient
}

// NewBookingService creates a new BookingService
func NewBookingService(
	orderRepo *repository.OrderRepo,
	flightRepo *repository.FlightRepo,
	temporalClient *TemporalClient,
) *BookingService {
	return &BookingService{
		orderRepo:      orderRepo,
		flightRepo:     flightRepo,
		temporalClient: temporalClient,
	}
}

// CreateOrderInput contains the parameters for creating an order
type CreateOrderInput struct {
	FlightID string
	Seats    []string
}

// CreateOrderOutput contains the result of order creation
type CreateOrderOutput struct {
	OrderID    string
	WorkflowID string
	Status     domain.OrderStatus
	ExpiresAt  time.Time
}

// CreateOrder creates a new booking order and starts the workflow
func (s *BookingService) CreateOrder(ctx context.Context, input CreateOrderInput) (*CreateOrderOutput, error) {
	// Validate flight exists
	_, err := s.flightRepo.FindByID(ctx, input.FlightID)
	if err != nil {
		return nil, err
	}

	// Validate seats are not empty
	if len(input.Seats) == 0 {
		return nil, domain.ErrSeatUnavailable
	}

	// Generate order ID
	orderID := uuid.New().String()

	// Calculate expiration (15 minutes from now)
	expiresAt := time.Now().Add(15 * time.Minute)

	// Start the booking workflow
	temporalInput := temporalpkg.BookingWorkflowInput{
		OrderID:  orderID,
		FlightID: input.FlightID,
		Seats:    input.Seats,
	}

	workflowID, err := s.temporalClient.StartBookingWorkflow(ctx, temporalInput)
	if err != nil {
		return nil, fmt.Errorf("start workflow: %w", err)
	}

	// Note: Order is created by the workflow's CreateOrder activity
	// We return optimistically assuming the workflow will create it

	return &CreateOrderOutput{
		OrderID:    orderID,
		WorkflowID: workflowID,
		Status:     domain.OrderStatusSeatsReserved,
		ExpiresAt:  expiresAt,
	}, nil
}

// GetOrderStatus queries the workflow for current order status
func (s *BookingService) GetOrderStatus(ctx context.Context, orderID string) (*domain.OrderStatusResponse, error) {
	// First try to query the workflow
	status, err := s.temporalClient.QueryBookingStatus(ctx, orderID)
	if err != nil {
		// If workflow query fails, try to get from database
		order, dbErr := s.orderRepo.FindByID(ctx, orderID)
		if dbErr != nil {
			return nil, domain.ErrOrderNotFound
		}

		// Return status from database (for completed/failed/expired orders)
		timerRemaining := 0
		if order.ExpiresAt != nil {
			remaining := time.Until(*order.ExpiresAt)
			if remaining > 0 {
				timerRemaining = int(remaining.Seconds())
			}
		}

		return &domain.OrderStatusResponse{
			OrderID:         order.ID,
			Status:          order.Status,
			Seats:           order.Seats,
			TimerRemaining:  timerRemaining,
			PaymentAttempts: 0,
			LastError:       stringValue(order.FailureReason),
		}, nil
	}

	return &domain.OrderStatusResponse{
		OrderID:         status.OrderID,
		Status:          status.Status,
		Seats:           status.Seats,
		TimerRemaining:  status.TimerRemaining,
		PaymentAttempts: status.PaymentAttempts,
		LastError:       status.LastError,
	}, nil
}

// UpdateSeatsOutput contains the result of seat update
type UpdateSeatsOutput struct {
	OrderID   string
	Status    domain.OrderStatus
	Seats     []string
	ExpiresAt time.Time
}

// UpdateSeats updates the seat selection for an order
func (s *BookingService) UpdateSeats(ctx context.Context, orderID string, seats []string) (*UpdateSeatsOutput, error) {
	// Validate seats are not empty
	if len(seats) == 0 {
		return nil, domain.ErrSeatUnavailable
	}

	// Send signal to workflow
	err := s.temporalClient.SignalUpdateSeats(ctx, orderID, seats)
	if err != nil {
		return nil, fmt.Errorf("signal update seats: %w", err)
	}

	// Query updated status
	status, err := s.temporalClient.QueryBookingStatus(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("query status: %w", err)
	}

	return &UpdateSeatsOutput{
		OrderID:   status.OrderID,
		Status:    status.Status,
		Seats:     status.Seats,
		ExpiresAt: status.ExpiresAt,
	}, nil
}

// SubmitPayment submits a payment for an order
func (s *BookingService) SubmitPayment(ctx context.Context, orderID string, paymentCode string) error {
	// Validate payment code format (5 digits)
	if !isValidPaymentCode(paymentCode) {
		return domain.ErrInvalidPaymentCode
	}

	// Send payment signal to workflow
	err := s.temporalClient.SignalProceedToPayment(ctx, orderID, paymentCode)
	if err != nil {
		return fmt.Errorf("signal payment: %w", err)
	}

	return nil
}

// CancelOrder cancels an order
func (s *BookingService) CancelOrder(ctx context.Context, orderID string) error {
	err := s.temporalClient.SignalCancelBooking(ctx, orderID)
	if err != nil {
		return fmt.Errorf("signal cancel: %w", err)
	}

	return nil
}

// Helper functions

func isValidPaymentCode(code string) bool {
	matched, _ := regexp.MatchString(`^\d{5}$`, code)
	return matched
}

func stringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
