import { useEffect, useRef, useState } from 'react';
import { ArrowTrendingUpIcon, ArrowTrendingDownIcon } from '@heroicons/react/24/outline';

/**
 * Score breakdown interface
 */
interface ScoreBreakdown {
  readonly blocked: number;
  readonly detected: number;
  readonly successful: number;
  readonly total: number;
}

interface SecurityScoreProps {
  /** Overall security score (0-100) */
  readonly score: number;
  /** Breakdown of blocked/detected/successful attacks */
  readonly breakdown?: ScoreBreakdown;
  /** Trend indicator: positive = improving, negative = declining */
  readonly trend?: number;
  /** Size variant */
  readonly size?: 'sm' | 'md' | 'lg';
  /** Whether to animate the gauge on mount */
  readonly animated?: boolean;
  /** Additional CSS classes */
  readonly className?: string;
}

/**
 * Size configurations
 */
const sizeConfig = {
  sm: { container: 'w-32 h-32', strokeWidth: 8, fontSize: 'text-xl', labelSize: 'text-xs' },
  md: { container: 'w-48 h-48', strokeWidth: 10, fontSize: 'text-3xl', labelSize: 'text-sm' },
  lg: { container: 'w-64 h-64', strokeWidth: 12, fontSize: 'text-4xl', labelSize: 'text-base' },
};

/**
 * Get color based on score value
 */
function getScoreColor(score: number): string {
  if (score < 50) return '#ef4444'; // danger-500 (red)
  if (score < 80) return '#f59e0b'; // warning-500 (amber)
  return '#22c55e'; // success-500 (green)
}

/**
 * Get text color class based on score
 */
function getScoreTextClass(score: number): string {
  if (score < 50) return 'text-danger-500';
  if (score < 80) return 'text-warning-500';
  return 'text-success-500';
}

/**
 * Animated circular gauge component for displaying security scores.
 * Features SVG-based gauge with color-coded values and optional breakdown stats.
 */
export function SecurityScore({
  score,
  breakdown,
  trend,
  size = 'md',
  animated = true,
  className = '',
}: SecurityScoreProps) {
  const [displayScore, setDisplayScore] = useState(animated ? 0 : score);
  const displayScoreRef = useRef(displayScore);
  const config = sizeConfig[size];

  // SVG circle calculations
  const radius = 45;
  const circumference = 2 * Math.PI * radius;
  const normalizedScore = Math.min(Math.max(score, 0), 100);
  const offset = circumference - (displayScore / 100) * circumference;
  const color = getScoreColor(normalizedScore);

  // Animate score on mount
  useEffect(() => {
    if (!animated) {
      displayScoreRef.current = normalizedScore;
      setDisplayScore(normalizedScore);
      return;
    }

    // Check for reduced motion preference
    const prefersReducedMotion = globalThis.matchMedia('(prefers-reduced-motion: reduce)').matches;
    if (prefersReducedMotion) {
      displayScoreRef.current = normalizedScore;
      setDisplayScore(normalizedScore);
      return;
    }

    // Animate from current displayed value to target score
    const duration = 1000;
    const startTime = performance.now();
    const startValue = displayScoreRef.current;
    const endValue = normalizedScore;
    let animationId: number;

    function animate(currentTime: number) {
      const elapsed = currentTime - startTime;
      const progress = Math.min(elapsed / duration, 1);
      // Ease out cubic
      const eased = 1 - Math.pow(1 - progress, 3);
      const current = startValue + (endValue - startValue) * eased;
      displayScoreRef.current = current;
      setDisplayScore(current);

      if (progress < 1) {
        animationId = requestAnimationFrame(animate);
      }
    }

    animationId = requestAnimationFrame(animate);
    return () => cancelAnimationFrame(animationId);
  }, [animated, normalizedScore]);

  return (
    <div className={`flex flex-col items-center ${className}`}>
      {/* SVG Gauge */}
      <div className={`relative ${config.container}`}>
        <meter
          value={normalizedScore}
          min={0}
          max={100}
          aria-label={`Security score: ${normalizedScore.toFixed(1)}%`}
          className="sr-only"
        />
        <svg
          viewBox="0 0 100 100"
          className="transform -rotate-90 w-full h-full"
          aria-hidden="true"
        >
          {/* Background circle */}
          <circle
            cx="50"
            cy="50"
            r={radius}
            fill="none"
            stroke="currentColor"
            strokeWidth={config.strokeWidth}
            className="text-gray-200 dark:text-gray-700"
          />
          {/* Score arc */}
          <circle
            cx="50"
            cy="50"
            r={radius}
            fill="none"
            stroke={color}
            strokeWidth={config.strokeWidth}
            strokeLinecap="round"
            strokeDasharray={circumference}
            strokeDashoffset={offset}
            className="transition-all duration-1000 ease-out"
          />
        </svg>

        {/* Center content */}
        <div className="absolute inset-0 flex flex-col items-center justify-center">
          <div className="flex items-center gap-1">
            <span className={`font-bold ${config.fontSize} ${getScoreTextClass(normalizedScore)}`}>
              {displayScore.toFixed(1)}
            </span>
            <span className={`${config.labelSize} text-gray-500 dark:text-gray-400`}>%</span>
          </div>
          {trend !== undefined && trend !== 0 && (
            <div className={`flex items-center gap-0.5 ${config.labelSize}`}>
              {trend > 0 ? (
                <>
                  <ArrowTrendingUpIcon className="h-4 w-4 text-success-500" />
                  <span className="text-success-500">+{trend.toFixed(1)}%</span>
                </>
              ) : (
                <>
                  <ArrowTrendingDownIcon className="h-4 w-4 text-danger-500" />
                  <span className="text-danger-500">{trend.toFixed(1)}%</span>
                </>
              )}
            </div>
          )}
        </div>
      </div>

      {/* Breakdown stats */}
      {breakdown && (
        <div className="mt-4 grid grid-cols-3 gap-4 text-center">
          <div>
            <p className={`font-semibold text-success-600 dark:text-success-500 ${size === 'sm' ? 'text-lg' : 'text-xl'}`}>
              {breakdown.blocked}
            </p>
            <p className="text-xs text-gray-500 dark:text-gray-400">Blocked</p>
          </div>
          <div>
            <p className={`font-semibold text-warning-600 dark:text-warning-500 ${size === 'sm' ? 'text-lg' : 'text-xl'}`}>
              {breakdown.detected}
            </p>
            <p className="text-xs text-gray-500 dark:text-gray-400">Detected</p>
          </div>
          <div>
            <p className={`font-semibold text-danger-600 dark:text-danger-500 ${size === 'sm' ? 'text-lg' : 'text-xl'}`}>
              {breakdown.successful}
            </p>
            <p className="text-xs text-gray-500 dark:text-gray-400">Success</p>
          </div>
        </div>
      )}
    </div>
  );
}
