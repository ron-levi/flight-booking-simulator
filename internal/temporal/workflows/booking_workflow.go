package workflows

import (
	"errors"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	"github.com/flight-booking-system/internal/domain"
	temporalpkg "github.com/flight-booking-system/internal/temporal"
	"github.com/flight-booking-system/internal/temporal/activities"
)

// BookingWorkflow manages the flight booking process
// - Reserves seats with 15-minute timer
// - Handles seat update signals (resets timer)
// - Processes payment on proceed signal
// - Releases seats on timeout/failure/cancellation
func BookingWorkflow(ctx workflow.Context, input temporalpkg.BookingWorkflowInput) (result temporalpkg.BookingWorkflowResult, err error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("BookingWorkflow started", "orderID", input.OrderID, "flightID", input.FlightID)

	// Initialize workflow state
	state := &bookingState{
		orderID:         input.OrderID,
		flightID:        input.FlightID,
		seats:           input.Seats,
		status:          domain.OrderStatusCreated,
		paymentAttempts: 0,
	}

	// Register query handler for status queries
	if err := workflow.SetQueryHandler(ctx, temporalpkg.QueryBookingStatus, func() (temporalpkg.BookingStatusResponse, error) {
		return state.toStatusResponse(), nil
	}); err != nil {
		return result, err
	}

	// Activity options for seat operations (short timeout, retries)
	seatActivityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    10 * time.Second,
			MaximumAttempts:    3,
		},
	}
	seatCtx := workflow.WithActivityOptions(ctx, seatActivityOptions)

	// Activity options for order operations (short timeout, retries)
	orderActivityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    10 * time.Second,
			MaximumAttempts:    3,
		},
	}
	orderCtx := workflow.WithActivityOptions(ctx, orderActivityOptions)

	// Activity options for payment (configured timeout and retries from PRD)
	paymentActivityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 1.5,
			MaximumInterval:    5 * time.Second,
			MaximumAttempts:    3,
			NonRetryableErrorTypes: []string{
				temporalpkg.ErrTypeInvalidPaymentCode,
				temporalpkg.ErrTypePaymentDeclined,
			},
		},
	}
	paymentCtx := workflow.WithActivityOptions(ctx, paymentActivityOptions)

	var a *activities.BookingActivities

	// Setup compensation for seat release on any failure
	defer func() {
		if err != nil || state.status == domain.OrderStatusExpired || state.status == domain.OrderStatusFailed {
			// Use disconnected context for cleanup (survives workflow cancellation)
			compensationCtx, _ := workflow.NewDisconnectedContext(ctx)
			compensationCtx = workflow.WithActivityOptions(compensationCtx, seatActivityOptions)

			releaseErr := workflow.ExecuteActivity(compensationCtx, a.ReleaseSeats, activities.ReleaseSeatsInput{
				OrderID:  state.orderID,
				FlightID: state.flightID,
				Seats:    state.seats,
			}).Get(compensationCtx, nil)

			if releaseErr != nil {
				logger.Error("Failed to release seats during compensation", "error", releaseErr)
			} else {
				logger.Info("Seats released during compensation", "seats", state.seats)
			}
		}
	}()

	// Phase 1: Reserve seats
	state.status = domain.OrderStatusSeatsReserved
	err = workflow.ExecuteActivity(seatCtx, a.ReserveSeats, activities.ReserveSeatInput{
		OrderID:  input.OrderID,
		FlightID: input.FlightID,
		Seats:    input.Seats,
	}).Get(seatCtx, nil)
	if err != nil {
		state.lastError = err.Error()
		state.status = domain.OrderStatusFailed
		return state.toResult(), err
	}
	logger.Info("Seats reserved", "seats", input.Seats)

	// Create order in database
	state.expiresAt = workflow.Now(ctx).Add(15 * time.Minute)
	err = workflow.ExecuteActivity(orderCtx, a.CreateOrder, activities.CreateOrderInput{
		OrderID:    input.OrderID,
		FlightID:   input.FlightID,
		WorkflowID: workflow.GetInfo(ctx).WorkflowExecution.ID,
		Seats:      input.Seats,
		ExpiresAt:  state.expiresAt,
	}).Get(orderCtx, nil)
	if err != nil {
		state.lastError = err.Error()
		state.status = domain.OrderStatusFailed
		return state.toResult(), err
	}

	// Phase 2: Wait for payment signal with 15-minute timeout
	// Handle seat update signals to reset timer
	seatUpdateChan := workflow.GetSignalChannel(ctx, temporalpkg.SignalUpdateSeats)
	paymentChan := workflow.GetSignalChannel(ctx, temporalpkg.SignalProceedToPay)
	cancelChan := workflow.GetSignalChannel(ctx, temporalpkg.SignalCancelBooking)

	var paymentSignal temporalpkg.PaymentSignal
	paymentReceived := false
	canceled := false

	for !paymentReceived && !canceled {
		// Create timer for remaining hold duration
		timerCtx, cancelTimer := workflow.WithCancel(ctx)
		timerDuration := state.expiresAt.Sub(workflow.Now(ctx))
		if timerDuration <= 0 {
			// Already expired
			state.status = domain.OrderStatusExpired
			state.lastError = "seat reservation expired"
			logger.Info("Seat hold expired")

			// Mark order as expired in database
			_ = workflow.ExecuteActivity(orderCtx, a.ExpireOrder, activities.ExpireOrderInput{
				OrderID: state.orderID,
			}).Get(orderCtx, nil)

			return state.toResult(), temporalpkg.ErrReservationExpired
		}

		holdTimer := workflow.NewTimer(timerCtx, timerDuration)

		selector := workflow.NewSelector(ctx)

		// Handle seat update signal
		selector.AddReceive(seatUpdateChan, func(c workflow.ReceiveChannel, more bool) {
			var signal temporalpkg.SeatUpdateSignal
			c.Receive(ctx, &signal)
			logger.Info("Received seat update signal", "newSeats", signal.Seats)

			// Update seat selection
			updateErr := workflow.ExecuteActivity(seatCtx, a.UpdateSeatSelection, activities.UpdateSeatSelectionInput{
				OrderID:  state.orderID,
				FlightID: state.flightID,
				OldSeats: state.seats,
				NewSeats: signal.Seats,
			}).Get(seatCtx, nil)

			if updateErr != nil {
				logger.Error("Failed to update seats", "error", updateErr)
				state.lastError = updateErr.Error()
			} else {
				state.seats = signal.Seats
				// Reset timer by updating expiration
				state.expiresAt = workflow.Now(ctx).Add(15 * time.Minute)

				// Update order in database
				_ = workflow.ExecuteActivity(orderCtx, a.UpdateOrderSeats, activities.UpdateOrderSeatsInput{
					OrderID:   state.orderID,
					Seats:     signal.Seats,
					ExpiresAt: state.expiresAt,
				}).Get(orderCtx, nil)

				logger.Info("Timer reset", "expiresAt", state.expiresAt)
			}

			cancelTimer() // Cancel current timer to restart with new duration
		})

		// Handle payment signal
		selector.AddReceive(paymentChan, func(c workflow.ReceiveChannel, more bool) {
			c.Receive(ctx, &paymentSignal)
			logger.Info("Received payment signal", "code", paymentSignal.PaymentCode[:2]+"***")
			paymentReceived = true
			cancelTimer()
		})

		// Handle cancel signal
		selector.AddReceive(cancelChan, func(c workflow.ReceiveChannel, more bool) {
			c.Receive(ctx, nil)
			logger.Info("Received cancel signal")
			canceled = true
			cancelTimer()
		})

		// Handle timer expiration
		selector.AddFuture(holdTimer, func(f workflow.Future) {
			timerErr := f.Get(timerCtx, nil)
			if timerErr == nil {
				// Timer actually expired (not canceled)
				state.status = domain.OrderStatusExpired
				state.lastError = "seat reservation expired"
				logger.Info("Seat hold timer expired")
			}
		})

		selector.Select(ctx)

		// Check if expired
		if state.status == domain.OrderStatusExpired {
			// Mark order as expired in database
			_ = workflow.ExecuteActivity(orderCtx, a.ExpireOrder, activities.ExpireOrderInput{
				OrderID: state.orderID,
			}).Get(orderCtx, nil)

			return state.toResult(), temporalpkg.ErrReservationExpired
		}
	}

	// Handle cancellation
	if canceled {
		state.status = domain.OrderStatusFailed
		state.lastError = "booking canceled by user"

		_ = workflow.ExecuteActivity(orderCtx, a.FailOrder, activities.FailOrderInput{
			OrderID: state.orderID,
			Reason:  state.lastError,
		}).Get(orderCtx, nil)

		return state.toResult(), temporalpkg.ErrWorkflowCanceled
	}

	// Phase 3: Process payment
	state.status = domain.OrderStatusPaymentProcessing
	_ = workflow.ExecuteActivity(orderCtx, a.UpdateOrderStatus, activities.UpdateOrderStatusInput{
		OrderID: state.orderID,
		Status:  domain.OrderStatusPaymentProcessing,
	}).Get(orderCtx, nil)

	var paymentResult activities.ValidatePaymentOutput
	err = workflow.ExecuteActivity(paymentCtx, a.ValidatePayment, activities.ValidatePaymentInput{
		OrderID:     state.orderID,
		PaymentCode: paymentSignal.PaymentCode,
	}).Get(paymentCtx, &paymentResult)

	// Track payment attempts (Temporal handles retries internally)
	state.paymentAttempts++

	if err != nil {
		state.status = domain.OrderStatusFailed
		state.lastError = "payment failed: " + err.Error()
		logger.Error("Payment validation failed", "error", err)

		// Check if it's a non-retryable error
		var appErr *temporal.ApplicationError
		if errors.As(err, &appErr) {
			state.lastError = "payment failed: " + appErr.Message()
		}

		_ = workflow.ExecuteActivity(orderCtx, a.FailOrder, activities.FailOrderInput{
			OrderID: state.orderID,
			Reason:  state.lastError,
		}).Get(orderCtx, nil)

		return state.toResult(), err
	}

	// Phase 4: Confirm booking
	state.status = domain.OrderStatusConfirmed
	err = workflow.ExecuteActivity(orderCtx, a.ConfirmOrder, activities.ConfirmOrderInput{
		OrderID:  state.orderID,
		FlightID: state.flightID,
		Seats:    state.seats,
	}).Get(orderCtx, nil)

	if err != nil {
		state.status = domain.OrderStatusFailed
		state.lastError = "confirmation failed: " + err.Error()
		logger.Error("Order confirmation failed", "error", err)

		_ = workflow.ExecuteActivity(orderCtx, a.FailOrder, activities.FailOrderInput{
			OrderID: state.orderID,
			Reason:  state.lastError,
		}).Get(orderCtx, nil)

		return state.toResult(), err
	}

	logger.Info("Booking confirmed", "orderID", state.orderID, "seats", state.seats)

	// Clear the error since compensation is not needed for successful bookings
	err = nil

	// Drain any remaining signals before completing
	drainSignals(ctx, seatUpdateChan, paymentChan, cancelChan)

	return state.toResult(), nil
}

// bookingState tracks the internal workflow state
type bookingState struct {
	orderID         string
	flightID        string
	seats           []string
	status          domain.OrderStatus
	expiresAt       time.Time
	paymentAttempts int
	lastError       string
}

// toStatusResponse converts state to query response
func (s *bookingState) toStatusResponse() temporalpkg.BookingStatusResponse {
	timerRemaining := 0
	if !s.expiresAt.IsZero() {
		remaining := time.Until(s.expiresAt)
		if remaining > 0 {
			timerRemaining = int(remaining.Seconds())
		}
	}

	return temporalpkg.BookingStatusResponse{
		OrderID:         s.orderID,
		FlightID:        s.flightID,
		Status:          s.status,
		Seats:           s.seats,
		ExpiresAt:       s.expiresAt,
		TimerRemaining:  timerRemaining,
		PaymentAttempts: s.paymentAttempts,
		LastError:       s.lastError,
	}
}

// toResult converts state to workflow result
func (s *bookingState) toResult() temporalpkg.BookingWorkflowResult {
	return temporalpkg.BookingWorkflowResult{
		OrderID: s.orderID,
		Status:  s.status,
		Seats:   s.seats,
		Error:   s.lastError,
	}
}

// drainSignals empties signal channels to prevent "unhandled signal" warnings
func drainSignals(_ workflow.Context, channels ...workflow.ReceiveChannel) {
	for _, ch := range channels {
		for {
			var discard any
			ok := ch.ReceiveAsync(&discard)
			if !ok {
				break
			}
		}
	}
}
