import { ReactNode } from 'react';

/**
 * Props for the LoadingState component.
 */
interface LoadingStateProps {
  /** Optional custom loading message */
  readonly message?: string;
  /** Optional custom className for the container */
  readonly className?: string;
}

/**
 * Reusable loading state component with animation.
 * Displays a pulsing loading indicator with customizable message.
 *
 * @example
 * ```tsx
 * if (isLoading) {
 *   return <LoadingState message="Loading agents..." />;
 * }
 * ```
 */
export function LoadingState({ message = 'Loading...', className = '' }: LoadingStateProps): ReactNode {
  return (
    <output className={`animate-pulse text-gray-500 dark:text-gray-400 ${className}`} aria-live="polite">
      <div className="flex items-center gap-2">
        <div className="h-4 w-4 rounded-full bg-gray-300 dark:bg-gray-600 animate-bounce" />
        <span>{message}</span>
      </div>
    </output>
  );
}

export default LoadingState;
