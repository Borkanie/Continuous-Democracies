import { useState, useEffect } from 'react';

/**
 * Custom hook to debounce a value, e.g., for search input.
 * @param value The value to debounce.
 * @param delay The debounce delay in milliseconds.
 * @return The debounced string value.
 */
export const useDebounce = (value: string, delay = 500) => {
  const [debouncedValue, setDebouncedValue] = useState(value);

  useEffect(() => {
    const handler = setTimeout(() => setDebouncedValue(value), delay);
    return () => clearTimeout(handler);
  }, [value, delay]);

  return debouncedValue;
};
