package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// SeatLockRepo handles distributed seat locking via Redis
type SeatLockRepo struct {
	client *redis.Client
}

// NewSeatLockRepo creates a new SeatLockRepo
func NewSeatLockRepo(client *redis.Client) *SeatLockRepo {
	return &SeatLockRepo{client: client}
}

// seatLockKey generates the Redis key for a seat lock
func seatLockKey(flightID, seatID string) string {
	return fmt.Sprintf("seat:lock:%s:%s", flightID, seatID)
}

// LockSeats attempts to lock multiple seats for an order
// Returns nil if all seats were locked, error otherwise
func (r *SeatLockRepo) LockSeats(ctx context.Context, flightID string, seatIDs []string, orderID string, ttl time.Duration) error {
	// Use a pipeline for atomic operations
	pipe := r.client.TxPipeline()

	// First, check if any seats are already locked
	for _, seatID := range seatIDs {
		key := seatLockKey(flightID, seatID)
		pipe.Get(ctx, key)
	}

	results, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return fmt.Errorf("check existing locks: %w", err)
	}

	// Check results - if any seat is already locked by a different order, fail
	for i, result := range results {
		if result.Err() == nil {
			existingOrderID, _ := result.(*redis.StringCmd).Result()
			if existingOrderID != orderID {
				return fmt.Errorf("seat %s already locked by order %s", seatIDs[i], existingOrderID)
			}
		}
	}

	// Now set all locks with NX (only if not exists) or update if same order
	pipe = r.client.TxPipeline()
	for _, seatID := range seatIDs {
		key := seatLockKey(flightID, seatID)
		pipe.Set(ctx, key, orderID, ttl)
	}

	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("set seat locks: %w", err)
	}

	return nil
}

// ReleaseLocks releases all seat locks for an order
func (r *SeatLockRepo) ReleaseLocks(ctx context.Context, flightID string, seatIDs []string, orderID string) error {
	for _, seatID := range seatIDs {
		key := seatLockKey(flightID, seatID)
		// Only delete if the lock belongs to this order (using Lua script)
		script := redis.NewScript(`
			if redis.call("get", KEYS[1]) == ARGV[1] then
				return redis.call("del", KEYS[1])
			else
				return 0
			end
		`)
		_, err := script.Run(ctx, r.client, []string{key}, orderID).Result()
		if err != nil && err != redis.Nil {
			return fmt.Errorf("release seat lock %s: %w", seatID, err)
		}
	}

	return nil
}

// ExtendLocks extends the TTL for all seat locks
func (r *SeatLockRepo) ExtendLocks(ctx context.Context, flightID string, seatIDs []string, orderID string, ttl time.Duration) error {
	for _, seatID := range seatIDs {
		key := seatLockKey(flightID, seatID)
		// Only extend if the lock belongs to this order
		script := redis.NewScript(`
			if redis.call("get", KEYS[1]) == ARGV[1] then
				return redis.call("pexpire", KEYS[1], ARGV[2])
			else
				return 0
			end
		`)
		_, err := script.Run(ctx, r.client, []string{key}, orderID, ttl.Milliseconds()).Result()
		if err != nil && err != redis.Nil {
			return fmt.Errorf("extend seat lock %s: %w", seatID, err)
		}
	}

	return nil
}

// GetLockedSeats returns all locked seat IDs for a flight
func (r *SeatLockRepo) GetLockedSeats(ctx context.Context, flightID string) (map[string]string, error) {
	pattern := fmt.Sprintf("seat:lock:%s:*", flightID)
	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, fmt.Errorf("get locked seat keys: %w", err)
	}

	if len(keys) == 0 {
		return make(map[string]string), nil
	}

	// Get all values
	pipe := r.client.Pipeline()
	cmds := make([]*redis.StringCmd, len(keys))
	for i, key := range keys {
		cmds[i] = pipe.Get(ctx, key)
	}

	_, err = pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("get locked seat values: %w", err)
	}

	result := make(map[string]string)
	for i, cmd := range cmds {
		if cmd.Err() == nil {
			// Extract seat ID from key (seat:lock:flightID:seatID)
			seatID := keys[i][len(fmt.Sprintf("seat:lock:%s:", flightID)):]
			result[seatID] = cmd.Val()
		}
	}

	return result, nil
}
