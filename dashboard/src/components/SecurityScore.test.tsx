import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { render, screen, act } from '@testing-library/react';
import { SecurityScore } from './SecurityScore';

describe('SecurityScore', () => {
  beforeEach(() => {
    // Mock matchMedia for reduced motion
    vi.stubGlobal('matchMedia', vi.fn().mockReturnValue({
      matches: false,
      addEventListener: vi.fn(),
      removeEventListener: vi.fn(),
    }));
  });

  afterEach(() => {
    vi.unstubAllGlobals();
  });

  it('renders score value correctly', () => {
    render(<SecurityScore score={75.5} animated={false} />);

    expect(screen.getByText('75.5')).toBeInTheDocument();
    expect(screen.getByText('%')).toBeInTheDocument();
  });

  it('applies danger color for score < 50', () => {
    render(<SecurityScore score={30} animated={false} />);

    const scoreElement = screen.getByText('30.0');
    expect(scoreElement).toHaveClass('text-danger-500');
  });

  it('applies warning color for score 50-79', () => {
    render(<SecurityScore score={65} animated={false} />);

    const scoreElement = screen.getByText('65.0');
    expect(scoreElement).toHaveClass('text-warning-500');
  });

  it('applies success color for score >= 80', () => {
    render(<SecurityScore score={90} animated={false} />);

    const scoreElement = screen.getByText('90.0');
    expect(scoreElement).toHaveClass('text-success-500');
  });

  it('renders breakdown when provided', () => {
    const breakdown = {
      blocked: 10,
      detected: 5,
      successful: 2,
      total: 17,
    };

    render(<SecurityScore score={75} breakdown={breakdown} animated={false} />);

    expect(screen.getByText('10')).toBeInTheDocument();
    expect(screen.getByText('5')).toBeInTheDocument();
    expect(screen.getByText('2')).toBeInTheDocument();
    expect(screen.getByText('Blocked')).toBeInTheDocument();
    expect(screen.getByText('Detected')).toBeInTheDocument();
    expect(screen.getByText('Success')).toBeInTheDocument();
  });

  it('does not render breakdown when not provided', () => {
    render(<SecurityScore score={75} animated={false} />);

    expect(screen.queryByText('Blocked')).not.toBeInTheDocument();
    expect(screen.queryByText('Detected')).not.toBeInTheDocument();
  });

  it('shows positive trend indicator when provided', () => {
    render(<SecurityScore score={75} trend={5.2} animated={false} />);

    expect(screen.getByText('+5.2%')).toBeInTheDocument();
  });

  it('shows negative trend indicator when provided', () => {
    render(<SecurityScore score={75} trend={-3.5} animated={false} />);

    expect(screen.getByText('-3.5%')).toBeInTheDocument();
  });

  it('does not show trend when trend is 0', () => {
    render(<SecurityScore score={75} trend={0} animated={false} />);

    expect(screen.queryByText('%', { selector: 'span.text-success-500' })).not.toBeInTheDocument();
    expect(screen.queryByText('%', { selector: 'span.text-danger-500' })).not.toBeInTheDocument();
  });

  it('respects size prop - sm', () => {
    const { container } = render(<SecurityScore score={75} size="sm" animated={false} />);

    expect(container.querySelector('.w-32')).toBeInTheDocument();
  });

  it('respects size prop - md (default)', () => {
    const { container } = render(<SecurityScore score={75} animated={false} />);

    expect(container.querySelector('.w-48')).toBeInTheDocument();
  });

  it('respects size prop - lg', () => {
    const { container } = render(<SecurityScore score={75} size="lg" animated={false} />);

    expect(container.querySelector('.w-64')).toBeInTheDocument();
  });

  it('has correct ARIA attributes', () => {
    const { container } = render(<SecurityScore score={75.5} animated={false} />);

    const meter = container.querySelector('meter');
    expect(meter).toBeInTheDocument();
    expect(meter).toHaveAttribute('value', '75.5');
    expect(meter).toHaveAttribute('min', '0');
    expect(meter).toHaveAttribute('max', '100');
    expect(meter).toHaveAttribute('aria-label', 'Security score: 75.5%');
  });

  it('handles edge case score of 0', () => {
    render(<SecurityScore score={0} animated={false} />);

    expect(screen.getByText('0.0')).toBeInTheDocument();
  });

  it('handles edge case score of 100', () => {
    render(<SecurityScore score={100} animated={false} />);

    expect(screen.getByText('100.0')).toBeInTheDocument();
    const scoreElement = screen.getByText('100.0');
    expect(scoreElement).toHaveClass('text-success-500');
  });

  it('clamps score to 0-100 range', () => {
    const { container } = render(<SecurityScore score={150} animated={false} />);

    const meter = container.querySelector('meter');
    expect(meter).toHaveAttribute('value', '100');
  });

  it('applies custom className', () => {
    const { container } = render(<SecurityScore score={75} className="custom-class" animated={false} />);

    expect(container.firstChild).toHaveClass('custom-class');
  });

  it('animates from current value on score prop change, not from 0', () => {
    // Mock requestAnimationFrame to run callback immediately with a time far enough for completion
    let rafCallback: ((time: number) => void) | null = null;
    vi.spyOn(globalThis, 'requestAnimationFrame').mockImplementation((cb) => {
      rafCallback = cb;
      return 1;
    });
    vi.spyOn(globalThis, 'cancelAnimationFrame').mockImplementation(() => {});
    vi.spyOn(performance, 'now').mockReturnValue(0);

    // Initial render with score=50, non-animated to set baseline
    const { rerender } = render(<SecurityScore score={50} animated={false} />);
    expect(screen.getByText('50.0')).toBeInTheDocument();

    // Re-render with score=80, animated
    act(() => {
      rerender(<SecurityScore score={80} animated={true} />);
    });

    // The animation should start - trigger the first frame at time=0 (start)
    if (rafCallback) {
      act(() => {
        (rafCallback as (time: number) => void)(0);
      });
    }

    // At time=0 (start of animation), the displayed score should be ~50 (the previous value),
    // NOT 0. This verifies the animation starts from the current displayed value.
    const displayedValue = screen.getByText(/\d+\.\d/);
    const numericValue = parseFloat(displayedValue.textContent || '0');
    expect(numericValue).toBeGreaterThanOrEqual(49);
  });
});
