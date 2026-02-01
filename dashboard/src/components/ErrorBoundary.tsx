import { Component, ErrorInfo, ReactNode } from 'react';
import { ExclamationTriangleIcon } from '@heroicons/react/24/outline';

/**
 * Props for the ErrorBoundary component.
 */
interface ErrorBoundaryProps {
  /** Child components to wrap */
  children: ReactNode;
  /** Optional fallback UI to render on error */
  fallback?: ReactNode;
  /** Optional callback when an error is caught */
  onError?: (error: Error, errorInfo: ErrorInfo) => void;
}

/**
 * State for the ErrorBoundary component.
 */
interface ErrorBoundaryState {
  /** Whether an error has been caught */
  hasError: boolean;
  /** The caught error, if any */
  error: Error | null;
}

/**
 * Error boundary component to catch and handle React errors gracefully.
 * Prevents the entire app from crashing when a child component throws.
 *
 * @example
 * ```tsx
 * <ErrorBoundary onError={(error) => logError(error)}>
 *   <ComponentThatMightFail />
 * </ErrorBoundary>
 * ```
 */
export class ErrorBoundary extends Component<ErrorBoundaryProps, ErrorBoundaryState> {
  constructor(props: ErrorBoundaryProps) {
    super(props);
    this.state = { hasError: false, error: null };
  }

  static getDerivedStateFromError(error: Error): ErrorBoundaryState {
    return { hasError: true, error };
  }

  componentDidCatch(error: Error, errorInfo: ErrorInfo): void {
    console.error('ErrorBoundary caught an error:', error, errorInfo);
    this.props.onError?.(error, errorInfo);
  }

  /**
   * Resets the error state to allow retry.
   */
  handleReset = (): void => {
    this.setState({ hasError: false, error: null });
  };

  render(): ReactNode {
    if (this.state.hasError) {
      if (this.props.fallback) {
        return this.props.fallback;
      }

      return (
        <div className="min-h-[200px] flex items-center justify-center">
          <div className="text-center p-8">
            <ExclamationTriangleIcon className="h-12 w-12 text-danger-500 mx-auto mb-4" />
            <h2 className="text-lg font-semibold text-gray-900 mb-2">Something went wrong</h2>
            <p className="text-gray-500 mb-4 max-w-md">
              {this.state.error?.message || 'An unexpected error occurred'}
            </p>
            <button
              className="btn-primary"
              onClick={this.handleReset}
            >
              Try Again
            </button>
          </div>
        </div>
      );
    }

    return this.props.children;
  }
}

export default ErrorBoundary;
