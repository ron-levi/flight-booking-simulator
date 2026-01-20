package workflows

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	"github.com/flight-booking-system/internal/temporal/activities"
)

// SeatReconciliationWorkflow reconciles Redis locks with DB seat status
// This workflow runs on a cron schedule to clean up orphaned locks
func SeatReconciliationWorkflow(ctx workflow.Context) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("Starting seat reconciliation workflow")

	// Activity options for reconciliation
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 3,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	// Get list of all flight IDs from database
	var flightIDs []string
	err := workflow.ExecuteActivity(ctx, "GetAllFlightIDs").Get(ctx, &flightIDs)
	if err != nil {
		logger.Error("Failed to get flight IDs", "error", err)
		return err
	}

	if len(flightIDs) == 0 {
		logger.Info("No flights found to reconcile")
		return nil
	}

	logger.Info("Reconciling locks for flights", "count", len(flightIDs))

	// Reconcile each flight
	for _, flightID := range flightIDs {
		input := activities.ReconcileSeatLocksInput{
			FlightID: flightID,
		}

		err := workflow.ExecuteActivity(ctx, "ReconcileSeatLocks", input).Get(ctx, nil)
		if err != nil {
			logger.Error("Failed to reconcile locks for flight", "flightID", flightID, "error", err)
			// Continue with other flights even if one fails
			continue
		}

		logger.Info("Successfully reconciled locks for flight", "flightID", flightID)
	}

	logger.Info("Completed seat reconciliation workflow")
	return nil
}
