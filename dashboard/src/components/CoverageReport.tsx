import { getTacticBarColor, formatTacticName } from '../lib/tacticColors';

interface CoverageReportProps {
  /** Coverage data by tactic - { tactic: count } */
  readonly coverage: Record<string, number>;
  /** Total techniques available in the system */
  readonly totalTechniques: number;
  /** Display variant */
  readonly variant?: 'compact' | 'detailed';
  /** Additional CSS classes */
  readonly className?: string;
}

/**
 * MITRE ATT&CK coverage visualization component.
 * Shows technique coverage by tactic with colored progress bars.
 */
export function CoverageReport({
  coverage,
  totalTechniques,
  variant = 'compact',
  className = '',
}: CoverageReportProps) {
  const tactics = Object.entries(coverage).sort((a, b) => b[1] - a[1]);
  const maxCount = Math.max(...Object.values(coverage), 1);
  const totalCovered = Object.values(coverage).reduce((sum, count) => sum + count, 0);

  if (tactics.length === 0) {
    return (
      <div className={`text-center py-8 text-gray-500 dark:text-gray-400 ${className}`}>
        No coverage data available
      </div>
    );
  }

  if (variant === 'compact') {
    return (
      <div className={className}>
        {/* Header stats */}
        <div className="flex items-baseline gap-2 mb-4">
          <span className="text-3xl font-bold text-gray-900 dark:text-gray-100">
            {totalTechniques}
          </span>
          <span className="text-sm text-gray-500 dark:text-gray-400">
            techniques across {tactics.length} tactics
          </span>
        </div>

        {/* Tactic bars */}
        <div className="space-y-3">
          {tactics.map(([tactic, count]) => (
            <TacticBar
              key={tactic}
              tactic={tactic}
              count={count}
              maxCount={maxCount}
            />
          ))}
        </div>
      </div>
    );
  }

  // Detailed variant
  return (
    <div className={className}>
      {/* Header with total coverage */}
      <div className="mb-6">
        <div className="flex items-baseline gap-2">
          <span className="text-4xl font-bold text-gray-900 dark:text-gray-100">
            {totalTechniques}
          </span>
          <span className="text-lg text-gray-500 dark:text-gray-400">techniques</span>
        </div>
        <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">
          Covering {tactics.length} MITRE ATT&CK tactics
        </p>
      </div>

      {/* Detailed tactic grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        {tactics.map(([tactic, count]) => (
          <TacticCard
            key={tactic}
            tactic={tactic}
            count={count}
            total={totalCovered}
          />
        ))}
      </div>
    </div>
  );
}

interface TacticBarProps {
  readonly tactic: string;
  readonly count: number;
  readonly maxCount: number;
}

/**
 * Compact tactic progress bar
 */
function TacticBar({ tactic, count, maxCount }: TacticBarProps) {
  const percentage = (count / maxCount) * 100;
  const barColor = getTacticBarColor(tactic);

  return (
    <div className="flex items-center gap-3">
      <span className="text-sm text-gray-600 dark:text-gray-400 w-28 truncate capitalize">
        {formatTacticName(tactic)}
      </span>
      <div className="flex-1 h-2 bg-gray-200 dark:bg-gray-700 rounded-full overflow-hidden">
        <div
          className={`h-full ${barColor} transition-all duration-500 ease-out rounded-full`}
          style={{ width: `${percentage}%` }}
        />
      </div>
      <span className="text-sm font-medium text-gray-900 dark:text-gray-100 w-8 text-right">
        {count}
      </span>
    </div>
  );
}

interface TacticCardProps {
  readonly tactic: string;
  readonly count: number;
  readonly total: number;
}

/**
 * Detailed tactic card with percentage
 */
function TacticCard({ tactic, count, total }: TacticCardProps) {
  const percentage = total > 0 ? (count / total) * 100 : 0;
  const barColor = getTacticBarColor(tactic);

  return (
    <div className="bg-gray-50 dark:bg-gray-800/50 rounded-lg p-4">
      <div className="flex items-center justify-between mb-2">
        <span className="text-sm font-medium text-gray-900 dark:text-gray-100 capitalize">
          {formatTacticName(tactic)}
        </span>
        <span className="text-sm text-gray-500 dark:text-gray-400">
          {percentage.toFixed(1)}%
        </span>
      </div>
      <div className="h-2 bg-gray-200 dark:bg-gray-700 rounded-full overflow-hidden">
        <div
          className={`h-full ${barColor} transition-all duration-500 ease-out rounded-full`}
          style={{ width: `${percentage}%` }}
        />
      </div>
      <p className="text-xs text-gray-500 dark:text-gray-400 mt-2">
        {count} technique{count === 1 ? '' : 's'}
      </p>
    </div>
  );
}
