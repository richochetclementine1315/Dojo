import { forwardRef } from 'react';
import type { InputHTMLAttributes } from 'react';

interface InputProps extends InputHTMLAttributes<HTMLInputElement> {}

export const Input = forwardRef<HTMLInputElement, InputProps>(
  ({ className = '', ...props }, ref) => {
    return (
      <input
        ref={ref}
        className={`cursor-target w-full px-4 py-2 bg-dojo-black-800 border border-dojo-red-900/30 rounded-xl text-white placeholder:text-gray-500 focus:outline-none focus:ring-2 focus:ring-dojo-red-600 focus:border-transparent transition-all ${className}`}
        {...props}
      />
    );
  }
);

Input.displayName = 'Input';