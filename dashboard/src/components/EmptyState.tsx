import { ReactNode, ComponentType } from 'react';

/**
 * Props for the EmptyState component.
 */
interface EmptyStateProps {
  /** Icon component to display (from @heroicons/react) */
  readonly icon: ComponentType<{ className?: string }>;
  /** Main title text */
  readonly title: string;
  /** Subtitle/description text */
  readonly description: string;
  /** Optional action button */
  readonly action?: {
    readonly label: string;
    readonly onClick: () => void;
  };
  /** Optional custom className for the container */
  readonly className?: string;
}

/**
 * Reusable empty state component for when lists are empty.
 * Displays an icon, title, description, and optional action button.
 *
 * @example
 * ```tsx
 * {items.length === 0 && (
 *   <EmptyState
 *     icon={ComputerDesktopIcon}
 *     title="No agents connected"
 *     description="Deploy an agent to get started"
 *     action={{ label: 'Add Agent', onClick: handleAddAgent }}
 *   />
 * )}
 * ```
 */
export function EmptyState({
  icon: Icon,
  title,
  description,
  action,
  className = '',
}: EmptyStateProps): ReactNode {
  return (
    <div className={`text-center py-12 ${className}`}>
      <Icon className="h-12 w-12 text-gray-400 mx-auto mb-4" aria-hidden="true" />
      <h3 className="text-lg font-medium text-gray-900">{title}</h3>
      <p className="text-gray-500 mt-1">{description}</p>
      {action && (
        <button
          className="btn-primary mt-4"
          onClick={action.onClick}
        >
          {action.label}
        </button>
      )}
    </div>
  );
}

export default EmptyState;
