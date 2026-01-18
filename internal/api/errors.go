package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/flight-booking-system/internal/domain"
)

// ErrorResponse represents an API error
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// Error codes
const (
	ErrCodeInvalidRequest   = "INVALID_REQUEST"
	ErrCodeInvalidSeats     = "INVALID_SEATS"
	ErrCodeFlightNotFound   = "FLIGHT_NOT_FOUND"
	ErrCodeOrderNotFound    = "ORDER_NOT_FOUND"
	ErrCodeOrderExpired     = "ORDER_EXPIRED"
	ErrCodeSeatsUnavailable = "SEATS_UNAVAILABLE"
	ErrCodePaymentFailed    = "PAYMENT_FAILED"
	ErrCodeInternalError    = "INTERNAL_ERROR"
	ErrCodeWorkflowError    = "WORKFLOW_ERROR"
)

// WriteError writes a JSON error response
func WriteError(w http.ResponseWriter, statusCode int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error:   code,
		Message: message,
	})
}

// WriteJSON writes a JSON response with the given status code
func WriteJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// MapDomainError maps domain errors to HTTP status codes and error codes
func MapDomainError(err error) (int, string, string) {
	switch {
	case errors.Is(err, domain.ErrFlightNotFound):
		return http.StatusNotFound, ErrCodeFlightNotFound, "Flight not found"
	case errors.Is(err, domain.ErrOrderNotFound):
		return http.StatusNotFound, ErrCodeOrderNotFound, "Order not found"
	case errors.Is(err, domain.ErrOrderExpired):
		return http.StatusConflict, ErrCodeOrderExpired, "Order reservation has expired"
	case errors.Is(err, domain.ErrSeatUnavailable), errors.Is(err, domain.ErrSeatsAlreadyLocked):
		return http.StatusConflict, ErrCodeSeatsUnavailable, "One or more seats are not available"
	case errors.Is(err, domain.ErrInvalidPaymentCode):
		return http.StatusBadRequest, ErrCodePaymentFailed, "Invalid payment code format"
	case errors.Is(err, domain.ErrPaymentFailed):
		return http.StatusBadRequest, ErrCodePaymentFailed, "Payment validation failed"
	default:
		return http.StatusInternalServerError, ErrCodeInternalError, "An internal error occurred"
	}
}

// HandleServiceError writes appropriate error response based on service error
func HandleServiceError(w http.ResponseWriter, err error) {
	statusCode, code, message := MapDomainError(err)
	WriteError(w, statusCode, code, message)
}
