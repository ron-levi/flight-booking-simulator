import { useState, useEffect, useCallback } from 'react';

/**
 * Custom hook for countdown timer
 * @param {number} initialSeconds - Starting seconds
 * @param {Function} onExpire - Callback when timer reaches 0
 * @returns {Object} - { seconds, minutes, formatted, isExpired, reset }
 */
export function useCountdown(initialSeconds, onExpire) {
  const [seconds, setSeconds] = useState(initialSeconds);

  // Reset function to update timer from server
  const reset = useCallback((newSeconds) => {
    setSeconds(newSeconds);
  }, []);

  useEffect(() => {
    // Update when initial value changes
    if (initialSeconds > 0) {
      setSeconds(initialSeconds);
    }
  }, [initialSeconds]);

  useEffect(() => {
    if (seconds <= 0) {
      onExpire?.();
      return;
    }

    const interval = setInterval(() => {
      setSeconds((prev) => {
        if (prev <= 1) {
          onExpire?.();
          return 0;
        }
        return prev - 1;
      });
    }, 1000);

    return () => clearInterval(interval);
  }, [seconds, onExpire]);

  const minutes = Math.floor(seconds / 60);
  const remainingSeconds = seconds % 60;
  const formatted = `${minutes}:${remainingSeconds.toString().padStart(2, '0')}`;
  const isExpired = seconds <= 0;

  return {
    seconds,
    minutes,
    remainingSeconds,
    formatted,
    isExpired,
    reset,
  };
}

export default useCountdown;
