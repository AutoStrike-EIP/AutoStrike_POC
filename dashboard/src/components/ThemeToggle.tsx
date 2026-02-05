import { SunIcon, MoonIcon, ComputerDesktopIcon } from '@heroicons/react/24/outline';
import { useTheme, Theme } from '../contexts/ThemeContext';

interface ThemeToggleProps {
  /** Display variant: icon for simple toggle, full for all options */
  readonly variant?: 'icon' | 'full';
  /** Additional CSS classes */
  readonly className?: string;
}

/**
 * Theme toggle component with icon and full variants.
 * Icon variant: Simple sun/moon toggle button
 * Full variant: Three buttons for Light/Dark/System selection
 */
export function ThemeToggle({ variant = 'icon', className = '' }: ThemeToggleProps) {
  const { theme, resolvedTheme, setTheme, toggleTheme } = useTheme();

  if (variant === 'icon') {
    return (
      <button
        onClick={toggleTheme}
        className={`p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors ${className}`}
        aria-label={`Switch to ${resolvedTheme === 'light' ? 'dark' : 'light'} mode`}
        title={`Switch to ${resolvedTheme === 'light' ? 'dark' : 'light'} mode`}
      >
        {resolvedTheme === 'light' ? (
          <MoonIcon className="h-5 w-5 text-gray-600 dark:text-gray-300" />
        ) : (
          <SunIcon className="h-5 w-5 text-gray-300" />
        )}
      </button>
    );
  }

  // Full variant with all theme options
  const options: Array<{ value: Theme; label: string; icon: typeof SunIcon }> = [
    { value: 'light', label: 'Light', icon: SunIcon },
    { value: 'dark', label: 'Dark', icon: MoonIcon },
    { value: 'system', label: 'System', icon: ComputerDesktopIcon },
  ];

  return (
    <div className={`flex gap-2 ${className}`} role="radiogroup" aria-label="Theme selection">
      {options.map(({ value, label, icon: Icon }) => (
        <button
          key={value}
          onClick={() => setTheme(value)}
          role="radio"
          aria-checked={theme === value}
          className={`flex items-center gap-2 px-4 py-2 rounded-lg transition-colors ${
            theme === value
              ? 'bg-primary-600 text-white'
              : 'bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-600'
          }`}
        >
          <Icon className="h-5 w-5" />
          <span>{label}</span>
        </button>
      ))}
    </div>
  );
}
