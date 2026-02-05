import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { ThemeToggle } from './ThemeToggle';
import { ThemeProvider } from '../contexts/ThemeContext';

describe('ThemeToggle', () => {
  const localStorageMock = {
    getItem: vi.fn(),
    setItem: vi.fn(),
    removeItem: vi.fn(),
    clear: vi.fn(),
    length: 0,
    key: vi.fn(),
  };

  const matchMediaMock = vi.fn();

  beforeEach(() => {
    vi.stubGlobal('localStorage', localStorageMock);
    vi.stubGlobal('matchMedia', matchMediaMock);
    matchMediaMock.mockReturnValue({
      matches: false,
      addEventListener: vi.fn(),
      removeEventListener: vi.fn(),
    });
    localStorageMock.getItem.mockReturnValue('light');
    document.documentElement.classList.remove('dark');
  });

  afterEach(() => {
    vi.unstubAllGlobals();
    vi.clearAllMocks();
  });

  it('renders icon variant by default', () => {
    render(
      <ThemeProvider>
        <ThemeToggle />
      </ThemeProvider>
    );

    expect(screen.getByRole('button')).toBeInTheDocument();
    expect(screen.getByLabelText(/switch to dark mode/i)).toBeInTheDocument();
  });

  it('renders correct icon for light theme', () => {
    localStorageMock.getItem.mockReturnValue('light');

    render(
      <ThemeProvider>
        <ThemeToggle />
      </ThemeProvider>
    );

    // Moon icon should be visible in light mode (to switch to dark)
    expect(screen.getByLabelText(/switch to dark mode/i)).toBeInTheDocument();
  });

  it('renders correct icon for dark theme', () => {
    localStorageMock.getItem.mockReturnValue('dark');

    render(
      <ThemeProvider>
        <ThemeToggle />
      </ThemeProvider>
    );

    // Sun icon should be visible in dark mode (to switch to light)
    expect(screen.getByLabelText(/switch to light mode/i)).toBeInTheDocument();
  });

  it('toggles theme on click (icon variant)', () => {
    localStorageMock.getItem.mockReturnValue('light');

    render(
      <ThemeProvider>
        <ThemeToggle />
      </ThemeProvider>
    );

    const button = screen.getByRole('button');
    fireEvent.click(button);

    expect(localStorageMock.setItem).toHaveBeenCalledWith('autostrike_theme', 'dark');
  });

  it('renders full variant with all theme options', () => {
    render(
      <ThemeProvider>
        <ThemeToggle variant="full" />
      </ThemeProvider>
    );

    expect(screen.getByText('Light')).toBeInTheDocument();
    expect(screen.getByText('Dark')).toBeInTheDocument();
    expect(screen.getByText('System')).toBeInTheDocument();
  });

  it('shows correct option as selected in full variant', () => {
    localStorageMock.getItem.mockReturnValue('dark');

    render(
      <ThemeProvider>
        <ThemeToggle variant="full" />
      </ThemeProvider>
    );

    const darkRadio = screen.getByText('Dark').closest('label')?.querySelector('input[type="radio"]');
    expect(darkRadio).toBeChecked();
  });

  it('sets theme when clicking option in full variant', () => {
    localStorageMock.getItem.mockReturnValue('light');

    render(
      <ThemeProvider>
        <ThemeToggle variant="full" />
      </ThemeProvider>
    );

    fireEvent.click(screen.getByText('System'));

    expect(localStorageMock.setItem).toHaveBeenCalledWith('autostrike_theme', 'system');
  });

  it('has proper ARIA attributes', () => {
    render(
      <ThemeProvider>
        <ThemeToggle variant="full" />
      </ThemeProvider>
    );

    expect(screen.getByRole('radiogroup')).toHaveAttribute('aria-label', 'Theme selection');
  });

  it('applies custom className', () => {
    render(
      <ThemeProvider>
        <ThemeToggle className="custom-class" />
      </ThemeProvider>
    );

    expect(screen.getByRole('button')).toHaveClass('custom-class');
  });
});
