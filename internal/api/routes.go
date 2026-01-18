package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/flight-booking-system/internal/database"
)

// RouterConfig holds dependencies for router creation
type RouterConfig struct {
	Pool        *pgxpool.Pool
	RedisClient *redis.Client
	Handlers    *Handlers
}

// NewRouter creates a new Chi router with all routes configured
func NewRouter(cfg RouterConfig) *chi.Mux {
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(CORS("http://localhost:3000", "http://localhost:5173"))

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		// Check database
		if err := database.HealthCheck(r.Context(), cfg.Pool); err != nil {
			http.Error(w, "database unhealthy", http.StatusServiceUnavailable)
			return
		}

		// Check Redis
		if err := database.RedisHealthCheck(r.Context(), cfg.RedisClient); err != nil {
			http.Error(w, "redis unhealthy", http.StatusServiceUnavailable)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// API routes
	r.Route("/api", func(r chi.Router) {
		// Flight routes
		r.Route("/flights", func(r chi.Router) {
			r.Get("/", cfg.Handlers.ListFlights)
			r.Get("/{flightId}", cfg.Handlers.GetFlight)
		})

		// Order routes
		r.Route("/orders", func(r chi.Router) {
			r.Post("/", cfg.Handlers.CreateOrder)

			r.Route("/{orderId}", func(r chi.Router) {
				r.Put("/seats", cfg.Handlers.UpdateSeats)
				r.Get("/status", cfg.Handlers.GetOrderStatus)
				r.Post("/pay", cfg.Handlers.SubmitPayment)
				r.Delete("/", cfg.Handlers.CancelOrder)
			})
		})
	})

	return r
}
