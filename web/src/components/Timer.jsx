import useCountdown from '../hooks/useCountdown';

function Timer({ seconds, onExpire, className = '' }) {
  const { formatted, isExpired, minutes } = useCountdown(seconds, onExpire);

  // Determine urgency color
  const getTimerColor = () => {
    if (isExpired) return 'text-red-600 bg-red-50';
    if (minutes < 2) return 'text-red-600 bg-red-50 animate-pulse';
    if (minutes < 5) return 'text-yellow-600 bg-yellow-50';
    return 'text-green-600 bg-green-50';
  };

  return (
    <div className={`inline-flex items-center gap-2 px-4 py-2 rounded-lg ${getTimerColor()} ${className}`}>
      <svg
        className="w-5 h-5"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
      >
        <path
          strokeLinecap="round"
          strokeLinejoin="round"
          strokeWidth={2}
          d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
        />
      </svg>
      <span className="font-mono text-lg font-bold">
        {isExpired ? 'Expired' : formatted}
      </span>
      {!isExpired && (
        <span className="text-sm opacity-75">remaining</span>
      )}
    </div>
  );
}

export default Timer;
