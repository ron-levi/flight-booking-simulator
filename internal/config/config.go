package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds all application configuration
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Temporal TemporalConfig
	Booking  BookingConfig
}

type ServerConfig struct {
	Host string
	Port int
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type TemporalConfig struct {
	Host      string
	Namespace string
	TaskQueue string
}

type BookingConfig struct {
	SeatReservationTimeout   time.Duration
	PaymentValidationTimeout time.Duration
	PaymentMaxRetries        int
	PaymentFailureRate       float64
}

// Load reads configuration from environment variables with defaults
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Port: getEnvInt("SERVER_PORT", 8080),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DATABASE_HOST", "localhost"),
			Port:     getEnvInt("DATABASE_PORT", 5433),
			User:     getEnv("DATABASE_USER", "flightapp"),
			Password: getEnv("DATABASE_PASSWORD", "flightapp"),
			Name:     getEnv("DATABASE_NAME", "flight_booking"),
			SSLMode:  getEnv("DATABASE_SSLMODE", "disable"),
		},
		Redis: RedisConfig{
			Addr:     getEnv("REDIS_ADDR", "localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
		},
		Temporal: TemporalConfig{
			Host:      getEnv("TEMPORAL_HOST", "localhost:7233"),
			Namespace: getEnv("TEMPORAL_NAMESPACE", "default"),
			TaskQueue: getEnv("TEMPORAL_TASK_QUEUE", "booking-queue"),
		},
		Booking: BookingConfig{
			SeatReservationTimeout:   getEnvDuration("SEAT_RESERVATION_TIMEOUT", 15*time.Minute),
			PaymentValidationTimeout: getEnvDuration("PAYMENT_VALIDATION_TIMEOUT", 10*time.Second),
			PaymentMaxRetries:        getEnvInt("PAYMENT_MAX_RETRIES", 3),
			PaymentFailureRate:       getEnvFloat("PAYMENT_FAILURE_RATE", 0.15),
		},
	}
}

// DatabaseURL returns the PostgreSQL connection string
func (c *DatabaseConfig) DatabaseURL() string {
	return "postgres://" + c.User + ":" + c.Password + "@" + c.Host + ":" + strconv.Itoa(c.Port) + "/" + c.Name + "?sslmode=" + c.SSLMode
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvFloat(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
			return floatVal
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
