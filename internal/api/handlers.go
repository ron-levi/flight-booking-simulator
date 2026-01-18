package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/flight-booking-system/internal/service"
)

// Handlers contains all HTTP handlers
type Handlers struct {
	flightService  *service.FlightService
	bookingService *service.BookingService
}

// NewHandlers creates a new Handlers instance
func NewHandlers(flightService *service.FlightService, bookingService *service.BookingService) *Handlers {
	return &Handlers{
		flightService:  flightService,
		bookingService: bookingService,
	}
}

// ListFlights handles GET /api/flights
func (h *Handlers) ListFlights(w http.ResponseWriter, r *http.Request) {
	flights, err := h.flightService.ListFlights(r.Context())
	if err != nil {
		HandleServiceError(w, err)
		return
	}

	response := FlightListResponse{
		Flights: make([]FlightResponse, len(flights)),
	}
	for i, f := range flights {
		response.Flights[i] = FlightResponse{
			ID:             f.ID,
			FlightNumber:   f.FlightNumber,
			Origin:         f.Origin,
			Destination:    f.Destination,
			DepartureTime:  f.DepartureTime,
			TotalSeats:     f.TotalSeats,
			AvailableSeats: f.AvailableSeats,
			PriceCents:     f.PriceCents,
		}
	}

	WriteJSON(w, http.StatusOK, response)
}

// GetFlight handles GET /api/flights/{flightId}
func (h *Handlers) GetFlight(w http.ResponseWriter, r *http.Request) {
	flightID := chi.URLParam(r, "flightId")
	if flightID == "" {
		WriteError(w, http.StatusBadRequest, ErrCodeInvalidRequest, "flight ID is required")
		return
	}

	flight, err := h.flightService.GetFlightWithSeats(r.Context(), flightID)
	if err != nil {
		HandleServiceError(w, err)
		return
	}

	// Build seat response
	seats := make([]SeatResponse, len(flight.SeatMap.Seats))
	for i, s := range flight.SeatMap.Seats {
		seats[i] = SeatResponse{
			ID:     s.ID,
			Row:    s.Row,
			Column: s.Column,
			Status: string(s.Status),
		}
	}

	response := FlightDetailResponse{
		FlightResponse: FlightResponse{
			ID:             flight.ID,
			FlightNumber:   flight.FlightNumber,
			Origin:         flight.Origin,
			Destination:    flight.Destination,
			DepartureTime:  flight.DepartureTime,
			TotalSeats:     flight.TotalSeats,
			AvailableSeats: flight.AvailableSeats,
			PriceCents:     flight.PriceCents,
		},
		SeatMap: SeatMapResponse{
			Rows:        flight.SeatMap.Rows,
			SeatsPerRow: flight.SeatMap.SeatsPerRow,
			Seats:       seats,
		},
	}

	WriteJSON(w, http.StatusOK, response)
}

// CreateOrder handles POST /api/orders
func (h *Handlers) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, ErrCodeInvalidRequest, "invalid request body")
		return
	}

	// Validate request
	if req.FlightID == "" {
		WriteError(w, http.StatusBadRequest, ErrCodeInvalidRequest, "flightId is required")
		return
	}
	if len(req.Seats) == 0 {
		WriteError(w, http.StatusBadRequest, ErrCodeInvalidSeats, "at least one seat must be selected")
		return
	}

	output, err := h.bookingService.CreateOrder(r.Context(), service.CreateOrderInput{
		FlightID: req.FlightID,
		Seats:    req.Seats,
	})
	if err != nil {
		HandleServiceError(w, err)
		return
	}

	response := CreateOrderResponse{
		OrderID:    output.OrderID,
		WorkflowID: output.WorkflowID,
		Status:     string(output.Status),
		ExpiresAt:  output.ExpiresAt,
	}

	WriteJSON(w, http.StatusCreated, response)
}

// UpdateSeats handles PUT /api/orders/{orderId}/seats
func (h *Handlers) UpdateSeats(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "orderId")
	if orderID == "" {
		WriteError(w, http.StatusBadRequest, ErrCodeInvalidRequest, "order ID is required")
		return
	}

	var req UpdateSeatsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, ErrCodeInvalidRequest, "invalid request body")
		return
	}

	if len(req.Seats) == 0 {
		WriteError(w, http.StatusBadRequest, ErrCodeInvalidSeats, "at least one seat must be selected")
		return
	}

	output, err := h.bookingService.UpdateSeats(r.Context(), orderID, req.Seats)
	if err != nil {
		HandleServiceError(w, err)
		return
	}

	response := UpdateSeatsResponse{
		OrderID:   output.OrderID,
		Status:    string(output.Status),
		Seats:     output.Seats,
		ExpiresAt: output.ExpiresAt,
	}

	WriteJSON(w, http.StatusOK, response)
}

// GetOrderStatus handles GET /api/orders/{orderId}/status
func (h *Handlers) GetOrderStatus(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "orderId")
	if orderID == "" {
		WriteError(w, http.StatusBadRequest, ErrCodeInvalidRequest, "order ID is required")
		return
	}

	status, err := h.bookingService.GetOrderStatus(r.Context(), orderID)
	if err != nil {
		HandleServiceError(w, err)
		return
	}

	response := OrderStatusResponse{
		OrderID:         status.OrderID,
		Status:          string(status.Status),
		Seats:           status.Seats,
		TimerRemaining:  status.TimerRemaining,
		PaymentAttempts: status.PaymentAttempts,
		LastError:       status.LastError,
	}

	WriteJSON(w, http.StatusOK, response)
}

// SubmitPayment handles POST /api/orders/{orderId}/pay
func (h *Handlers) SubmitPayment(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "orderId")
	if orderID == "" {
		WriteError(w, http.StatusBadRequest, ErrCodeInvalidRequest, "order ID is required")
		return
	}

	var req SubmitPaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, ErrCodeInvalidRequest, "invalid request body")
		return
	}

	if req.PaymentCode == "" {
		WriteError(w, http.StatusBadRequest, ErrCodePaymentFailed, "payment code is required")
		return
	}

	err := h.bookingService.SubmitPayment(r.Context(), orderID, req.PaymentCode)
	if err != nil {
		HandleServiceError(w, err)
		return
	}

	response := PaymentAcceptedResponse{
		OrderID: orderID,
		Status:  "PAYMENT_PROCESSING",
	}

	WriteJSON(w, http.StatusAccepted, response)
}

// CancelOrder handles DELETE /api/orders/{orderId}
func (h *Handlers) CancelOrder(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "orderId")
	if orderID == "" {
		WriteError(w, http.StatusBadRequest, ErrCodeInvalidRequest, "order ID is required")
		return
	}

	err := h.bookingService.CancelOrder(r.Context(), orderID)
	if err != nil {
		HandleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
