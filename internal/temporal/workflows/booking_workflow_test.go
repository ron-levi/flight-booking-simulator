package workflows_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/testsuite"

	"github.com/flight-booking-system/internal/domain"
	temporalpkg "github.com/flight-booking-system/internal/temporal"
	"github.com/flight-booking-system/internal/temporal/activities"
	"github.com/flight-booking-system/internal/temporal/workflows"
)

func TestBookingWorkflow_Success(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	// Register activities (nil struct is fine since we're mocking all calls)
	var a *activities.BookingActivities
	env.RegisterActivity(a)

	// Mock activities using activity function names
	env.OnActivity(a.ReserveSeats, mock.Anything, mock.Anything).Return(nil)
	env.OnActivity(a.CreateOrder, mock.Anything, mock.Anything).Return(nil)
	env.OnActivity(a.UpdateOrderStatus, mock.Anything, mock.Anything).Return(nil)
	env.OnActivity(a.ValidatePayment, mock.Anything, mock.Anything).Return(
		activities.ValidatePaymentOutput{Success: true, Message: "OK"}, nil,
	)
	env.OnActivity(a.ConfirmOrder, mock.Anything, mock.Anything).Return(nil)
	env.OnActivity(a.ReleaseSeats, mock.Anything, mock.Anything).Return(nil)

	// Send payment signal after workflow starts
	env.RegisterDelayedCallback(func() {
		env.SignalWorkflow(temporalpkg.SignalProceedToPay, temporalpkg.PaymentSignal{
			PaymentCode: "12345",
		})
	}, time.Second)

	// Execute workflow
	env.ExecuteWorkflow(workflows.BookingWorkflow, temporalpkg.BookingWorkflowInput{
		OrderID:  "test-order-1",
		FlightID: "test-flight-1",
		Seats:    []string{"1A", "1B"},
	})

	require.True(t, env.IsWorkflowCompleted())
	require.NoError(t, env.GetWorkflowError())

	var result temporalpkg.BookingWorkflowResult
	require.NoError(t, env.GetWorkflowResult(&result))
	require.Equal(t, domain.OrderStatusConfirmed, result.Status)
	require.Equal(t, "test-order-1", result.OrderID)
}

func TestBookingWorkflow_TimerExpired(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	// Register activities
	var a *activities.BookingActivities
	env.RegisterActivity(a)

	// Mock activities
	env.OnActivity(a.ReserveSeats, mock.Anything, mock.Anything).Return(nil)
	env.OnActivity(a.CreateOrder, mock.Anything, mock.Anything).Return(nil)
	env.OnActivity(a.ExpireOrder, mock.Anything, mock.Anything).Return(nil)
	env.OnActivity(a.ReleaseSeats, mock.Anything, mock.Anything).Return(nil)

	// Don't send payment signal - let timer expire

	// Execute workflow
	env.ExecuteWorkflow(workflows.BookingWorkflow, temporalpkg.BookingWorkflowInput{
		OrderID:  "test-order-2",
		FlightID: "test-flight-1",
		Seats:    []string{"2A"},
	})

	require.True(t, env.IsWorkflowCompleted())
	// Workflow returns an error when reservation expires
	workflowErr := env.GetWorkflowError()
	require.Error(t, workflowErr)
	require.Contains(t, workflowErr.Error(), "seat reservation expired")
}

func TestBookingWorkflow_SeatUpdateResetsTimer(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	// Register activities
	var a *activities.BookingActivities
	env.RegisterActivity(a)

	// Mock activities
	env.OnActivity(a.ReserveSeats, mock.Anything, mock.Anything).Return(nil)
	env.OnActivity(a.CreateOrder, mock.Anything, mock.Anything).Return(nil)
	env.OnActivity(a.UpdateSeatSelection, mock.Anything, mock.Anything).Return(nil)
	env.OnActivity(a.UpdateOrderSeats, mock.Anything, mock.Anything).Return(nil)
	env.OnActivity(a.UpdateOrderStatus, mock.Anything, mock.Anything).Return(nil)
	env.OnActivity(a.ValidatePayment, mock.Anything, mock.Anything).Return(
		activities.ValidatePaymentOutput{Success: true, Message: "OK"}, nil,
	)
	env.OnActivity(a.ConfirmOrder, mock.Anything, mock.Anything).Return(nil)
	env.OnActivity(a.ReleaseSeats, mock.Anything, mock.Anything).Return(nil)

	// Send seat update signal at 14 minutes (would expire at 15 min)
	env.RegisterDelayedCallback(func() {
		env.SignalWorkflow(temporalpkg.SignalUpdateSeats, temporalpkg.SeatUpdateSignal{
			Seats: []string{"3A", "3B"},
		})
	}, 14*time.Minute)

	// Send payment signal at 16 minutes (after original timeout but before new timeout)
	env.RegisterDelayedCallback(func() {
		env.SignalWorkflow(temporalpkg.SignalProceedToPay, temporalpkg.PaymentSignal{
			PaymentCode: "12345",
		})
	}, 16*time.Minute)

	// Execute workflow
	env.ExecuteWorkflow(workflows.BookingWorkflow, temporalpkg.BookingWorkflowInput{
		OrderID:  "test-order-3",
		FlightID: "test-flight-1",
		Seats:    []string{"1A", "1B"},
	})

	require.True(t, env.IsWorkflowCompleted())
	require.NoError(t, env.GetWorkflowError())

	var result temporalpkg.BookingWorkflowResult
	require.NoError(t, env.GetWorkflowResult(&result))
	require.Equal(t, domain.OrderStatusConfirmed, result.Status)
	require.Equal(t, []string{"3A", "3B"}, result.Seats)
}

func TestBookingWorkflow_QueryStatus(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	// Register activities
	var a *activities.BookingActivities
	env.RegisterActivity(a)

	// Mock activities
	env.OnActivity(a.ReserveSeats, mock.Anything, mock.Anything).Return(nil)
	env.OnActivity(a.CreateOrder, mock.Anything, mock.Anything).Return(nil)
	env.OnActivity(a.UpdateOrderStatus, mock.Anything, mock.Anything).Return(nil)
	env.OnActivity(a.ValidatePayment, mock.Anything, mock.Anything).Return(
		activities.ValidatePaymentOutput{Success: true, Message: "OK"}, nil,
	)
	env.OnActivity(a.ConfirmOrder, mock.Anything, mock.Anything).Return(nil)
	env.OnActivity(a.ReleaseSeats, mock.Anything, mock.Anything).Return(nil)

	// Query status during workflow execution
	env.RegisterDelayedCallback(func() {
		result, err := env.QueryWorkflow(temporalpkg.QueryBookingStatus)
		require.NoError(t, err)

		var status temporalpkg.BookingStatusResponse
		require.NoError(t, result.Get(&status))
		require.Equal(t, "test-order-4", status.OrderID)
		require.Equal(t, domain.OrderStatusSeatsReserved, status.Status)
		require.True(t, status.TimerRemaining > 0)

		// Now send payment
		env.SignalWorkflow(temporalpkg.SignalProceedToPay, temporalpkg.PaymentSignal{
			PaymentCode: "12345",
		})
	}, time.Second)

	// Execute workflow
	env.ExecuteWorkflow(workflows.BookingWorkflow, temporalpkg.BookingWorkflowInput{
		OrderID:  "test-order-4",
		FlightID: "test-flight-1",
		Seats:    []string{"4A"},
	})

	require.True(t, env.IsWorkflowCompleted())
	require.NoError(t, env.GetWorkflowError())
}

func TestBookingWorkflow_Canceled(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	// Register activities
	var a *activities.BookingActivities
	env.RegisterActivity(a)

	// Mock activities
	env.OnActivity(a.ReserveSeats, mock.Anything, mock.Anything).Return(nil)
	env.OnActivity(a.CreateOrder, mock.Anything, mock.Anything).Return(nil)
	env.OnActivity(a.FailOrder, mock.Anything, mock.Anything).Return(nil)
	env.OnActivity(a.ReleaseSeats, mock.Anything, mock.Anything).Return(nil)

	// Send cancel signal
	env.RegisterDelayedCallback(func() {
		env.SignalWorkflow(temporalpkg.SignalCancelBooking, nil)
	}, time.Second)

	// Execute workflow
	env.ExecuteWorkflow(workflows.BookingWorkflow, temporalpkg.BookingWorkflowInput{
		OrderID:  "test-order-5",
		FlightID: "test-flight-1",
		Seats:    []string{"5A"},
	})

	require.True(t, env.IsWorkflowCompleted())
	// Workflow returns an error when canceled
	workflowErr := env.GetWorkflowError()
	require.Error(t, workflowErr)
	require.Contains(t, workflowErr.Error(), "booking workflow canceled")
}
