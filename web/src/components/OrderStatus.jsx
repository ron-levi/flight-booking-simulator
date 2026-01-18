import { getStatusMessage, getStatusColor } from '../hooks/useOrderStatus';
import LoadingSpinner from './LoadingSpinner';

function OrderStatus({ status, paymentAttempts = 0, lastError = '' }) {
  const statusMessage = getStatusMessage(status);
  const statusColor = getStatusColor(status);

  const getIcon = () => {
    switch (status) {
      case 'CONFIRMED':
        return (
          <svg className="w-8 h-8 text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
          </svg>
        );
      case 'FAILED':
      case 'EXPIRED':
        return (
          <svg className="w-8 h-8 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
          </svg>
        );
      case 'PAYMENT_PROCESSING':
        return <LoadingSpinner size="md" />;
      default:
        return null;
    }
  };

  return (
    <div className={`p-4 rounded-lg border ${
      status === 'CONFIRMED' ? 'bg-green-50 border-green-200' :
      status === 'FAILED' || status === 'EXPIRED' ? 'bg-red-50 border-red-200' :
      status === 'PAYMENT_PROCESSING' ? 'bg-yellow-50 border-yellow-200' :
      'bg-gray-50 border-gray-200'
    }`}>
      <div className="flex items-center gap-3">
        {getIcon()}
        <div>
          <div className={`font-semibold ${statusColor}`}>
            {statusMessage}
          </div>

          {/* Payment attempts indicator */}
          {status === 'PAYMENT_PROCESSING' && paymentAttempts > 0 && (
            <div className="text-sm text-gray-600">
              Attempt {paymentAttempts} of 3
              {paymentAttempts > 1 && ' (retrying...)'}
            </div>
          )}

          {/* Error message */}
          {lastError && (status === 'FAILED' || status === 'EXPIRED') && (
            <div className="text-sm text-red-600 mt-1">
              {lastError}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

export default OrderStatus;
