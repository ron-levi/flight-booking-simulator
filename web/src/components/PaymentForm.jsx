import { useState } from 'react';

function PaymentForm({ onSubmit, isLoading, disabled = false }) {
  const [paymentCode, setPaymentCode] = useState('');
  const [error, setError] = useState('');

  const handleChange = (e) => {
    const value = e.target.value.replace(/\D/g, '').slice(0, 5);
    setPaymentCode(value);
    setError('');
  };

  const handleSubmit = (e) => {
    e.preventDefault();

    if (paymentCode.length !== 5) {
      setError('Payment code must be exactly 5 digits');
      return;
    }

    onSubmit(paymentCode);
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div>
        <label htmlFor="paymentCode" className="block text-sm font-medium text-gray-700 mb-1">
          Payment Code
        </label>
        <input
          id="paymentCode"
          type="text"
          inputMode="numeric"
          pattern="\d*"
          value={paymentCode}
          onChange={handleChange}
          placeholder="Enter 5-digit code"
          disabled={disabled || isLoading}
          className={`input text-center text-2xl tracking-widest font-mono ${
            error ? 'border-red-500 focus:ring-red-500' : ''
          }`}
          autoComplete="off"
        />
        {error && (
          <p className="mt-1 text-sm text-red-600">{error}</p>
        )}
        <p className="mt-2 text-xs text-gray-500">
          Demo mode: Any 5-digit code works (15% simulated failure rate)
        </p>
      </div>

      <button
        type="submit"
        disabled={disabled || isLoading || paymentCode.length !== 5}
        className={`w-full py-3 rounded font-medium transition-colors ${
          disabled || isLoading || paymentCode.length !== 5
            ? 'bg-gray-200 text-gray-500 cursor-not-allowed'
            : 'btn-success'
        }`}
      >
        {isLoading ? (
          <span className="flex items-center justify-center gap-2">
            <span className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin" />
            Processing...
          </span>
        ) : (
          'Pay Now'
        )}
      </button>
    </form>
  );
}

export default PaymentForm;
