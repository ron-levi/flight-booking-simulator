package domain

import "time"

// OrderStatus represents the current status of an order
type OrderStatus string

const (
	OrderStatusCreated           OrderStatus = "CREATED"
	OrderStatusSeatsReserved     OrderStatus = "SEATS_RESERVED"
	OrderStatusPaymentPending    OrderStatus = "PAYMENT_PENDING"
	OrderStatusPaymentProcessing OrderStatus = "PAYMENT_PROCESSING"
	OrderStatusConfirmed         OrderStatus = "CONFIRMED"
	OrderStatusFailed            OrderStatus = "FAILED"
	OrderStatusExpired           OrderStatus = "EXPIRED"
)

// Order represents a booking order
type Order struct {
	ID              string      `json:"id"`
	FlightID        string      `json:"flightId"`
	WorkflowID      string      `json:"workflowId"`
	Status          OrderStatus `json:"status"`
	Seats           []string    `json:"seats"`
	TotalPriceCents int64       `json:"totalPriceCents"`
	PaymentCode     *string     `json:"paymentCode,omitempty"`
	ExpiresAt       *time.Time  `json:"expiresAt,omitempty"`
	ConfirmedAt     *time.Time  `json:"confirmedAt,omitempty"`
	FailureReason   *string     `json:"failureReason,omitempty"`
	CreatedAt       time.Time   `json:"createdAt"`
	UpdatedAt       time.Time   `json:"updatedAt"`
}

// OrderStatusResponse represents the status response for polling
type OrderStatusResponse struct {
	OrderID         string      `json:"orderId"`
	Status          OrderStatus `json:"status"`
	Seats           []string    `json:"seats"`
	TimerRemaining  int         `json:"timerRemaining"` // seconds
	PaymentAttempts int         `json:"paymentAttempts"`
	LastError       string      `json:"lastError,omitempty"`
}

// IsTerminal returns true if the order is in a final state
func (o *Order) IsTerminal() bool {
	return o.Status == OrderStatusConfirmed ||
		o.Status == OrderStatusFailed ||
		o.Status == OrderStatusExpired
}

// CanTransitionTo checks if the order can transition to the given status
func (o *Order) CanTransitionTo(status OrderStatus) bool {
	validTransitions := map[OrderStatus][]OrderStatus{
		OrderStatusCreated:           {OrderStatusSeatsReserved, OrderStatusFailed},
		OrderStatusSeatsReserved:     {OrderStatusPaymentPending, OrderStatusExpired, OrderStatusFailed},
		OrderStatusPaymentPending:    {OrderStatusPaymentProcessing, OrderStatusExpired, OrderStatusFailed},
		OrderStatusPaymentProcessing: {OrderStatusConfirmed, OrderStatusFailed},
	}

	allowed, exists := validTransitions[o.Status]
	if !exists {
		return false
	}

	for _, s := range allowed {
		if s == status {
			return true
		}
	}
	return false
}
