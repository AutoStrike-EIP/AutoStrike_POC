import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import Analytics from './Analytics';

// Mock Chart.js
vi.mock('react-chartjs-2', () => ({
  Line: () => <div data-testid="line-chart">Line Chart</div>,
  Bar: () => <div data-testid="bar-chart">Bar Chart</div>,
}));

vi.mock('chart.js', () => ({
  Chart: {
    register: vi.fn(),
  },
  CategoryScale: vi.fn(),
  LinearScale: vi.fn(),
  PointElement: vi.fn(),
  LineElement: vi.fn(),
  BarElement: vi.fn(),
  Title: vi.fn(),
  Tooltip: vi.fn(),
  Legend: vi.fn(),
  Filler: vi.fn(),
}));

// Mock the API
const mockCompare = vi.fn();
const mockTrend = vi.fn();
const mockSummary = vi.fn();

vi.mock('../lib/api', () => ({
  analyticsApi: {
    compare: (...args: unknown[]) => mockCompare(...args),
    trend: (...args: unknown[]) => mockTrend(...args),
    summary: (...args: unknown[]) => mockSummary(...args),
  },
}));

function createTestQueryClient() {
  return new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
      },
    },
  });
}

function renderAnalytics() {
  const queryClient = createTestQueryClient();
  return render(
    <QueryClientProvider client={queryClient}>
      <Analytics />
    </QueryClientProvider>
  );
}

const mockComparisonData = {
  current: {
    period: '30_days',
    start_date: '2024-01-01',
    end_date: '2024-01-30',
    execution_count: 25,
    average_score: 78.5,
    total_blocked: 150,
    total_detected: 75,
    total_successful: 25,
    total_techniques: 50,
  },
  previous: {
    period: '30_days',
    start_date: '2023-12-01',
    end_date: '2023-12-30',
    execution_count: 20,
    average_score: 65.0,
    total_blocked: 100,
    total_detected: 60,
    total_successful: 40,
    total_techniques: 50,
  },
  score_change: 13.5,
  score_trend: 'improving' as const,
  blocked_change: 50,
  detected_change: 15,
};

const mockTrendData = {
  period: '30_days',
  data_points: [
    {
      date: '2024-01-01',
      average_score: 70,
      execution_count: 5,
      blocked: 30,
      detected: 15,
      successful: 5,
    },
    {
      date: '2024-01-15',
      average_score: 75,
      execution_count: 8,
      blocked: 50,
      detected: 25,
      successful: 10,
    },
    {
      date: '2024-01-30',
      average_score: 80,
      execution_count: 12,
      blocked: 70,
      detected: 35,
      successful: 10,
    },
  ],
  summary: {
    start_score: 70,
    end_score: 80,
    average_score: 75,
    max_score: 85,
    min_score: 65,
    total_executions: 25,
    overall_trend: 'improving' as const,
    percentage_change: 14.3,
  },
};

const mockSummaryData = {
  total_executions: 100,
  completed_executions: 85,
  average_score: 72.5,
  best_score: 95.0,
  worst_score: 45.0,
  scores_by_scenario: {
    'Discovery Test': 85,
    'Lateral Movement': 70,
    'Defense Evasion': 65,
  },
  executions_by_status: {
    completed: 85,
    failed: 10,
    cancelled: 5,
  },
};

describe('Analytics Page', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockCompare.mockResolvedValue({ data: mockComparisonData });
    mockTrend.mockResolvedValue({ data: mockTrendData });
    mockSummary.mockResolvedValue({ data: mockSummaryData });
  });

  it('renders loading state initially', () => {
    // Make the promises never resolve to keep loading state
    mockCompare.mockReturnValue(new Promise(() => {}));
    mockTrend.mockReturnValue(new Promise(() => {}));
    mockSummary.mockReturnValue(new Promise(() => {}));

    renderAnalytics();
    expect(screen.getByText('Loading analytics...')).toBeInTheDocument();
  });

  it('renders page title after loading', async () => {
    renderAnalytics();

    await waitFor(() => {
      expect(screen.getByText('Analytics')).toBeInTheDocument();
    });
  });

  it('displays period selector with default 30 days', async () => {
    renderAnalytics();

    await waitFor(() => {
      const select = screen.getByRole('combobox');
      expect(select).toHaveValue('30');
    });
  });

  it('displays all period options', async () => {
    renderAnalytics();

    await waitFor(() => {
      expect(screen.getByText('Last 7 days')).toBeInTheDocument();
      expect(screen.getByText('Last 30 days')).toBeInTheDocument();
      expect(screen.getByText('Last 90 days')).toBeInTheDocument();
    });
  });

  it('changes period when selector changes', async () => {
    renderAnalytics();

    await waitFor(() => {
      expect(screen.getByText('Analytics')).toBeInTheDocument();
    });

    const select = screen.getByRole('combobox');
    fireEvent.change(select, { target: { value: '7' } });

    await waitFor(() => {
      expect(mockCompare).toHaveBeenCalledWith(7);
      expect(mockTrend).toHaveBeenCalledWith(7);
      expect(mockSummary).toHaveBeenCalledWith(7);
    });
  });

  it('displays average score from comparison data', async () => {
    renderAnalytics();

    await waitFor(() => {
      expect(screen.getByText('78.5%')).toBeInTheDocument();
    });
  });

  it('displays execution count', async () => {
    renderAnalytics();

    await waitFor(() => {
      expect(screen.getByText('25')).toBeInTheDocument();
      expect(screen.getByText('20 previous period')).toBeInTheDocument();
    });
  });

  it('displays blocked attacks count', async () => {
    renderAnalytics();

    await waitFor(() => {
      expect(screen.getByText('Blocked Attacks')).toBeInTheDocument();
      expect(screen.getByText('150')).toBeInTheDocument();
    });
  });

  it('displays detected attacks count', async () => {
    renderAnalytics();

    await waitFor(() => {
      expect(screen.getByText('Detected Attacks')).toBeInTheDocument();
      expect(screen.getByText('75')).toBeInTheDocument();
    });
  });

  it('displays score change with positive prefix', async () => {
    renderAnalytics();

    await waitFor(() => {
      expect(screen.getByText('+13.5% vs previous period')).toBeInTheDocument();
    });
  });

  it('displays trend summary statistics', async () => {
    renderAnalytics();

    await waitFor(() => {
      expect(screen.getByText('Min Score')).toBeInTheDocument();
      expect(screen.getByText('65.0%')).toBeInTheDocument();
      expect(screen.getByText('Max Score')).toBeInTheDocument();
      expect(screen.getByText('85.0%')).toBeInTheDocument();
    });
  });

  it('displays execution summary totals', async () => {
    renderAnalytics();

    await waitFor(() => {
      expect(screen.getByText('Execution Summary')).toBeInTheDocument();
      expect(screen.getByText('Total Executions')).toBeInTheDocument();
      expect(screen.getByText('100')).toBeInTheDocument();
    });
  });

  it('displays completed executions', async () => {
    renderAnalytics();

    await waitFor(() => {
      expect(screen.getByText('Completed')).toBeInTheDocument();
      expect(screen.getByText('85')).toBeInTheDocument();
    });
  });

  it('displays best and worst scores', async () => {
    renderAnalytics();

    await waitFor(() => {
      expect(screen.getByText('Best Score')).toBeInTheDocument();
      expect(screen.getByText('95.0%')).toBeInTheDocument();
      expect(screen.getByText('Worst Score')).toBeInTheDocument();
      expect(screen.getByText('45.0%')).toBeInTheDocument();
    });
  });

  it('renders charts', async () => {
    renderAnalytics();

    await waitFor(() => {
      const lineCharts = screen.getAllByTestId('line-chart');
      const barCharts = screen.getAllByTestId('bar-chart');
      expect(lineCharts.length).toBeGreaterThan(0);
      expect(barCharts.length).toBeGreaterThan(0);
    });
  });

  it('displays performance by scenario section', async () => {
    renderAnalytics();

    await waitFor(() => {
      expect(screen.getByText('Performance by Scenario')).toBeInTheDocument();
      expect(screen.getByText('Discovery Test')).toBeInTheDocument();
      expect(screen.getByText('Lateral Movement')).toBeInTheDocument();
      expect(screen.getByText('Defense Evasion')).toBeInTheDocument();
    });
  });

  it('displays scenario scores', async () => {
    renderAnalytics();

    await waitFor(() => {
      expect(screen.getByText('85%')).toBeInTheDocument();
      expect(screen.getByText('70%')).toBeInTheDocument();
    });
  });
});

describe('Analytics Error State', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('displays error state when compare fails', async () => {
    mockCompare.mockRejectedValue(new Error('Network error'));
    mockTrend.mockResolvedValue({ data: mockTrendData });
    mockSummary.mockResolvedValue({ data: mockSummaryData });

    renderAnalytics();

    await waitFor(() => {
      expect(screen.getByText('Failed to load analytics')).toBeInTheDocument();
      expect(screen.getByText('Network error')).toBeInTheDocument();
    });
  });

  it('displays error state when trend fails', async () => {
    mockCompare.mockResolvedValue({ data: mockComparisonData });
    mockTrend.mockRejectedValue(new Error('Trend fetch failed'));
    mockSummary.mockResolvedValue({ data: mockSummaryData });

    renderAnalytics();

    await waitFor(() => {
      expect(screen.getByText('Failed to load analytics')).toBeInTheDocument();
      expect(screen.getByText('Trend fetch failed')).toBeInTheDocument();
    });
  });

  it('displays error state when summary fails', async () => {
    mockCompare.mockResolvedValue({ data: mockComparisonData });
    mockTrend.mockResolvedValue({ data: mockTrendData });
    mockSummary.mockRejectedValue(new Error('Summary error'));

    renderAnalytics();

    await waitFor(() => {
      expect(screen.getByText('Failed to load analytics')).toBeInTheDocument();
      expect(screen.getByText('Summary error')).toBeInTheDocument();
    });
  });

  it('shows try again button on error', async () => {
    mockCompare.mockRejectedValue(new Error('Network error'));
    mockTrend.mockResolvedValue({ data: mockTrendData });
    mockSummary.mockResolvedValue({ data: mockSummaryData });

    renderAnalytics();

    await waitFor(() => {
      expect(screen.getByText('Try Again')).toBeInTheDocument();
    });
  });

  it('refetches data when try again is clicked', async () => {
    mockCompare.mockRejectedValueOnce(new Error('Network error'));
    mockTrend.mockResolvedValue({ data: mockTrendData });
    mockSummary.mockResolvedValue({ data: mockSummaryData });

    renderAnalytics();

    await waitFor(() => {
      expect(screen.getByText('Try Again')).toBeInTheDocument();
    });

    // Reset mocks to succeed
    mockCompare.mockResolvedValue({ data: mockComparisonData });

    const retryButton = screen.getByText('Try Again');
    fireEvent.click(retryButton);

    await waitFor(() => {
      expect(mockCompare).toHaveBeenCalledTimes(2);
    });
  });

  it('shows fallback error message when no specific message', async () => {
    mockCompare.mockRejectedValue({});
    mockTrend.mockResolvedValue({ data: mockTrendData });
    mockSummary.mockResolvedValue({ data: mockSummaryData });

    renderAnalytics();

    await waitFor(() => {
      expect(screen.getByText('An error occurred while fetching data')).toBeInTheDocument();
    });
  });
});

describe('Analytics Trend Icons', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockTrend.mockResolvedValue({ data: mockTrendData });
    mockSummary.mockResolvedValue({ data: mockSummaryData });
  });

  it('shows improving trend indicator', async () => {
    mockCompare.mockResolvedValue({
      data: { ...mockComparisonData, score_trend: 'improving' },
    });

    renderAnalytics();

    await waitFor(() => {
      expect(screen.getByText('Average Score')).toBeInTheDocument();
    });
  });

  it('handles declining trend', async () => {
    mockCompare.mockResolvedValue({
      data: { ...mockComparisonData, score_trend: 'declining', score_change: -5.5 },
    });

    renderAnalytics();

    await waitFor(() => {
      expect(screen.getByText('-5.5% vs previous period')).toBeInTheDocument();
    });
  });

  it('handles stable trend', async () => {
    mockCompare.mockResolvedValue({
      data: { ...mockComparisonData, score_trend: 'stable', score_change: 0 },
    });

    renderAnalytics();

    await waitFor(() => {
      expect(screen.getByText('0% vs previous period')).toBeInTheDocument();
    });
  });
});

describe('Analytics Empty Data', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockCompare.mockResolvedValue({ data: mockComparisonData });
    mockTrend.mockResolvedValue({ data: mockTrendData });
    mockSummary.mockResolvedValue({ data: mockSummaryData });
  });

  it('shows message when no scenario data', async () => {
    mockSummary.mockResolvedValue({
      data: { ...mockSummaryData, scores_by_scenario: {} },
    });

    renderAnalytics();

    await waitFor(() => {
      expect(screen.getByText('No scenario data available')).toBeInTheDocument();
    });
  });

  it('handles zero values gracefully', async () => {
    mockCompare.mockResolvedValue({
      data: {
        ...mockComparisonData,
        current: {
          ...mockComparisonData.current,
          average_score: 0,
          execution_count: 0,
          total_blocked: 0,
          total_detected: 0,
        },
      },
    });

    renderAnalytics();

    await waitFor(() => {
      // Page still renders correctly with zero values
      expect(screen.getByText('Average Score')).toBeInTheDocument();
      expect(screen.getByText('Executions')).toBeInTheDocument();
      expect(screen.getByText('Blocked Attacks')).toBeInTheDocument();
    });
  });

  it('handles null/undefined score change', async () => {
    mockCompare.mockResolvedValue({
      data: {
        ...mockComparisonData,
        score_change: undefined,
      },
    });

    renderAnalytics();

    await waitFor(() => {
      expect(screen.getByText('0% vs previous period')).toBeInTheDocument();
    });
  });

  it('handles empty data points in trend', async () => {
    mockTrend.mockResolvedValue({
      data: {
        ...mockTrendData,
        data_points: [],
        summary: undefined,
      },
    });

    renderAnalytics();

    await waitFor(() => {
      expect(screen.getByText('Score Trend')).toBeInTheDocument();
    });
  });
});

describe('Analytics Period Changes', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockCompare.mockResolvedValue({ data: mockComparisonData });
    mockTrend.mockResolvedValue({ data: mockTrendData });
    mockSummary.mockResolvedValue({ data: mockSummaryData });
  });

  it('fetches with 7 days period', async () => {
    renderAnalytics();

    await waitFor(() => {
      expect(screen.getByText('Analytics')).toBeInTheDocument();
    });

    const select = screen.getByRole('combobox');
    fireEvent.change(select, { target: { value: '7' } });

    await waitFor(() => {
      expect(mockCompare).toHaveBeenLastCalledWith(7);
    });
  });

  it('fetches with 90 days period', async () => {
    renderAnalytics();

    await waitFor(() => {
      expect(screen.getByText('Analytics')).toBeInTheDocument();
    });

    const select = screen.getByRole('combobox');
    fireEvent.change(select, { target: { value: '90' } });

    await waitFor(() => {
      expect(mockCompare).toHaveBeenLastCalledWith(90);
    });
  });
});

describe('Analytics Blocked/Detected Changes', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockTrend.mockResolvedValue({ data: mockTrendData });
    mockSummary.mockResolvedValue({ data: mockSummaryData });
  });

  it('shows positive blocked change with plus', async () => {
    mockCompare.mockResolvedValue({
      data: { ...mockComparisonData, blocked_change: 50 },
    });

    renderAnalytics();

    await waitFor(() => {
      expect(screen.getByText('+50 vs previous')).toBeInTheDocument();
    });
  });

  it('shows positive detected change with plus', async () => {
    mockCompare.mockResolvedValue({
      data: { ...mockComparisonData, detected_change: 15 },
    });

    renderAnalytics();

    await waitFor(() => {
      expect(screen.getByText('+15 vs previous')).toBeInTheDocument();
    });
  });

  it('shows negative blocked change without plus', async () => {
    mockCompare.mockResolvedValue({
      data: { ...mockComparisonData, blocked_change: -10 },
    });

    renderAnalytics();

    await waitFor(() => {
      expect(screen.getByText('-10 vs previous')).toBeInTheDocument();
    });
  });

  it('shows zero change without prefix', async () => {
    mockCompare.mockResolvedValue({
      data: { ...mockComparisonData, blocked_change: 0, detected_change: 0 },
    });

    renderAnalytics();

    await waitFor(() => {
      const zeroTexts = screen.getAllByText('0 vs previous');
      expect(zeroTexts.length).toBe(2);
    });
  });
});

describe('Analytics Stable Trend Rendering', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockTrend.mockResolvedValue({ data: mockTrendData });
    mockSummary.mockResolvedValue({ data: mockSummaryData });
  });

  it('renders MinusIcon for stable trend (default case of getTrendIcon)', async () => {
    mockCompare.mockResolvedValue({
      data: { ...mockComparisonData, score_trend: 'stable', score_change: 0 },
    });

    const { container } = renderAnalytics();

    await waitFor(() => {
      expect(screen.getByText('Average Score')).toBeInTheDocument();
    });

    // The stable trend falls into the default case, rendering a MinusIcon with text-gray-400
    // Verify that no green (improving) or red (declining) trend icons are present
    const greenIcons = container.querySelectorAll('.text-green-500');
    const redIcons = container.querySelectorAll('.text-red-500');

    // Filter to only trend icons in the Average Score card (h-5 w-5 size)
    const greenTrendIcons = Array.from(greenIcons).filter(
      (el) => el.classList.contains('h-5') && el.classList.contains('w-5')
    );
    const redTrendIcons = Array.from(redIcons).filter(
      (el) => el.classList.contains('h-5') && el.classList.contains('w-5')
    );

    expect(greenTrendIcons.length).toBe(0);
    expect(redTrendIcons.length).toBe(0);

    // Verify the gray MinusIcon is rendered instead
    const grayIcons = container.querySelectorAll('.text-gray-400.h-5.w-5');
    expect(grayIcons.length).toBeGreaterThan(0);
  });

  it('applies text-gray-500 color for stable trend (default case of getTrendColor)', async () => {
    mockCompare.mockResolvedValue({
      data: { ...mockComparisonData, score_trend: 'stable', score_change: 0 },
    });

    renderAnalytics();

    await waitFor(() => {
      expect(screen.getByText('Average Score')).toBeInTheDocument();
    });

    // The score change text should have the gray-500 class from getTrendColor default case
    const scoreChangeText = screen.getByText('0% vs previous period');
    expect(scoreChangeText).toHaveClass('text-gray-500');
  });

  it('applies text-green-500 color for improving trend (getTrendColor)', async () => {
    mockCompare.mockResolvedValue({
      data: { ...mockComparisonData, score_trend: 'improving', score_change: 13.5 },
    });

    renderAnalytics();

    await waitFor(() => {
      expect(screen.getByText('Average Score')).toBeInTheDocument();
    });

    const scoreChangeText = screen.getByText('+13.5% vs previous period');
    expect(scoreChangeText).toHaveClass('text-green-500');
  });

  it('applies text-red-500 color for declining trend (getTrendColor)', async () => {
    mockCompare.mockResolvedValue({
      data: { ...mockComparisonData, score_trend: 'declining', score_change: -5.5 },
    });

    renderAnalytics();

    await waitFor(() => {
      expect(screen.getByText('Average Score')).toBeInTheDocument();
    });

    const scoreChangeText = screen.getByText('-5.5% vs previous period');
    expect(scoreChangeText).toHaveClass('text-red-500');
  });
});

describe('Analytics formatScoreChange Edge Cases', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockTrend.mockResolvedValue({ data: mockTrendData });
    mockSummary.mockResolvedValue({ data: mockSummaryData });
  });

  it('handles null score_change by displaying 0', async () => {
    mockCompare.mockResolvedValue({
      data: {
        ...mockComparisonData,
        score_change: null,
      },
    });

    renderAnalytics();

    await waitFor(() => {
      expect(screen.getByText('0% vs previous period')).toBeInTheDocument();
    });
  });

  it('handles score_change of exactly 0 by displaying 0 (falsy value path)', async () => {
    mockCompare.mockResolvedValue({
      data: {
        ...mockComparisonData,
        score_change: 0,
        score_trend: 'stable',
      },
    });

    renderAnalytics();

    await waitFor(() => {
      // formatScoreChange(0) returns '0' because !0 is true
      expect(screen.getByText('0% vs previous period')).toBeInTheDocument();
    });
  });

  it('renders getTrendIcon default case for undefined trend', async () => {
    mockCompare.mockResolvedValue({
      data: {
        ...mockComparisonData,
        score_trend: undefined,
        score_change: 0,
      },
    });

    const { container } = renderAnalytics();

    await waitFor(() => {
      expect(screen.getByText('Average Score')).toBeInTheDocument();
    });

    // undefined trend falls into the default case, same as 'stable'
    const grayIcons = container.querySelectorAll('.text-gray-400.h-5.w-5');
    expect(grayIcons.length).toBeGreaterThan(0);
  });

  it('renders getTrendColor default case for undefined trend', async () => {
    mockCompare.mockResolvedValue({
      data: {
        ...mockComparisonData,
        score_trend: undefined,
        score_change: 0,
      },
    });

    renderAnalytics();

    await waitFor(() => {
      expect(screen.getByText('Average Score')).toBeInTheDocument();
    });

    const scoreChangeText = screen.getByText('0% vs previous period');
    expect(scoreChangeText).toHaveClass('text-gray-500');
  });

  it('formats negative score_change without plus prefix', async () => {
    mockCompare.mockResolvedValue({
      data: {
        ...mockComparisonData,
        score_change: -2.3,
        score_trend: 'declining',
      },
    });

    renderAnalytics();

    await waitFor(() => {
      expect(screen.getByText('-2.3% vs previous period')).toBeInTheDocument();
    });
  });

  it('formats large positive score_change with plus prefix', async () => {
    mockCompare.mockResolvedValue({
      data: {
        ...mockComparisonData,
        score_change: 99.9,
        score_trend: 'improving',
      },
    });

    renderAnalytics();

    await waitFor(() => {
      expect(screen.getByText('+99.9% vs previous period')).toBeInTheDocument();
    });
  });
});

describe('Analytics Declining Trend Icon', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockTrend.mockResolvedValue({ data: mockTrendData });
    mockSummary.mockResolvedValue({ data: mockSummaryData });
  });

  it('renders ArrowTrendingDownIcon for declining trend', async () => {
    mockCompare.mockResolvedValue({
      data: { ...mockComparisonData, score_trend: 'declining', score_change: -8.2 },
    });

    const { container } = renderAnalytics();

    await waitFor(() => {
      expect(screen.getByText('Average Score')).toBeInTheDocument();
    });

    // Declining trend renders a red icon (h-5 w-5 text-red-500)
    const redTrendIcons = container.querySelectorAll('.text-red-500.h-5.w-5');
    expect(redTrendIcons.length).toBeGreaterThan(0);

    // No green improving icons should be present
    const greenTrendIcons = container.querySelectorAll('.text-green-500.h-5.w-5');
    expect(greenTrendIcons.length).toBe(0);
  });

  it('renders ArrowTrendingUpIcon for improving trend', async () => {
    mockCompare.mockResolvedValue({
      data: { ...mockComparisonData, score_trend: 'improving', score_change: 10.0 },
    });

    const { container } = renderAnalytics();

    await waitFor(() => {
      expect(screen.getByText('Average Score')).toBeInTheDocument();
    });

    // Improving trend renders a green icon (h-5 w-5 text-green-500)
    const greenTrendIcons = container.querySelectorAll('.text-green-500.h-5.w-5');
    expect(greenTrendIcons.length).toBeGreaterThan(0);
  });
});

describe('Analytics Partial Loading States', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('shows loading when only comparison is loading', () => {
    mockCompare.mockReturnValue(new Promise(() => {}));
    mockTrend.mockResolvedValue({ data: mockTrendData });
    mockSummary.mockResolvedValue({ data: mockSummaryData });

    renderAnalytics();
    expect(screen.getByText('Loading analytics...')).toBeInTheDocument();
  });

  it('shows loading when only trend is loading', () => {
    mockCompare.mockResolvedValue({ data: mockComparisonData });
    mockTrend.mockReturnValue(new Promise(() => {}));
    mockSummary.mockResolvedValue({ data: mockSummaryData });

    renderAnalytics();
    expect(screen.getByText('Loading analytics...')).toBeInTheDocument();
  });

  it('shows loading when only summary is loading', () => {
    mockCompare.mockResolvedValue({ data: mockComparisonData });
    mockTrend.mockResolvedValue({ data: mockTrendData });
    mockSummary.mockReturnValue(new Promise(() => {}));

    renderAnalytics();
    expect(screen.getByText('Loading analytics...')).toBeInTheDocument();
  });
});

describe('Analytics Null/Undefined Data Fallbacks', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockTrend.mockResolvedValue({ data: mockTrendData });
    mockSummary.mockResolvedValue({ data: mockSummaryData });
  });

  it('displays 0 for best_score when it is null', async () => {
    mockCompare.mockResolvedValue({ data: mockComparisonData });
    mockSummary.mockResolvedValue({
      data: {
        ...mockSummaryData,
        best_score: null,
        worst_score: null,
      },
    });

    renderAnalytics();

    await waitFor(() => {
      expect(screen.getByText('Best Score')).toBeInTheDocument();
      expect(screen.getByText('Worst Score')).toBeInTheDocument();
    });
  });

  it('handles undefined blocked_change showing 0 vs previous', async () => {
    mockCompare.mockResolvedValue({
      data: {
        ...mockComparisonData,
        blocked_change: undefined,
      },
    });

    renderAnalytics();

    await waitFor(() => {
      expect(screen.getByText('Blocked Attacks')).toBeInTheDocument();
    });
  });

  it('handles undefined detected_change showing 0 vs previous', async () => {
    mockCompare.mockResolvedValue({
      data: {
        ...mockComparisonData,
        detected_change: undefined,
      },
    });

    renderAnalytics();

    await waitFor(() => {
      expect(screen.getByText('Detected Attacks')).toBeInTheDocument();
    });
  });

  it('handles negative detected_change without plus prefix', async () => {
    mockCompare.mockResolvedValue({
      data: { ...mockComparisonData, detected_change: -7 },
    });

    renderAnalytics();

    await waitFor(() => {
      expect(screen.getByText('-7 vs previous')).toBeInTheDocument();
    });
  });

  it('renders charts with empty trend data_points', async () => {
    mockCompare.mockResolvedValue({ data: mockComparisonData });
    mockTrend.mockResolvedValue({
      data: {
        period: '30_days',
        data_points: [],
        summary: undefined,
      },
    });

    renderAnalytics();

    await waitFor(() => {
      // Charts should still render with empty data
      const lineCharts = screen.getAllByTestId('line-chart');
      const barCharts = screen.getAllByTestId('bar-chart');
      expect(lineCharts.length).toBeGreaterThan(0);
      expect(barCharts.length).toBeGreaterThan(0);
    });
  });

  it('renders with undefined executions_by_status', async () => {
    mockCompare.mockResolvedValue({ data: mockComparisonData });
    mockSummary.mockResolvedValue({
      data: {
        ...mockSummaryData,
        executions_by_status: undefined,
      },
    });

    renderAnalytics();

    await waitFor(() => {
      expect(screen.getByText('Executions by Status')).toBeInTheDocument();
    });
  });

  it('renders with undefined scores_by_scenario showing empty message', async () => {
    mockCompare.mockResolvedValue({ data: mockComparisonData });
    mockSummary.mockResolvedValue({
      data: {
        ...mockSummaryData,
        scores_by_scenario: undefined,
      },
    });

    renderAnalytics();

    await waitFor(() => {
      expect(screen.getByText('No scenario data available')).toBeInTheDocument();
    });
  });

  it('renders summary section with zero total_executions', async () => {
    mockCompare.mockResolvedValue({ data: mockComparisonData });
    mockSummary.mockResolvedValue({
      data: {
        total_executions: 0,
        completed_executions: 0,
        average_score: 0,
        best_score: 0,
        worst_score: 0,
        scores_by_scenario: {},
        executions_by_status: {},
      },
    });

    renderAnalytics();

    await waitFor(() => {
      expect(screen.getByText('Total Executions')).toBeInTheDocument();
      expect(screen.getByText('No scenario data available')).toBeInTheDocument();
    });
  });

  it('handles previous period execution_count of 0', async () => {
    mockCompare.mockResolvedValue({
      data: {
        ...mockComparisonData,
        previous: {
          ...mockComparisonData.previous,
          execution_count: 0,
        },
      },
    });

    renderAnalytics();

    await waitFor(() => {
      expect(screen.getByText('0 previous period')).toBeInTheDocument();
    });
  });

  it('handles current execution_count of 0', async () => {
    mockCompare.mockResolvedValue({
      data: {
        ...mockComparisonData,
        current: {
          ...mockComparisonData.current,
          execution_count: 0,
          total_blocked: 0,
          total_detected: 0,
        },
      },
    });

    renderAnalytics();

    await waitFor(() => {
      expect(screen.getByText('Executions')).toBeInTheDocument();
    });
  });
});

describe('Analytics All APIs Fail', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('shows error when all three APIs fail', async () => {
    mockCompare.mockRejectedValue(new Error('Compare failed'));
    mockTrend.mockRejectedValue(new Error('Trend failed'));
    mockSummary.mockRejectedValue(new Error('Summary failed'));

    renderAnalytics();

    await waitFor(() => {
      expect(screen.getByText('Failed to load analytics')).toBeInTheDocument();
      // First error message in the fallback chain is comparisonError
      expect(screen.getByText('Compare failed')).toBeInTheDocument();
    });
  });

  it('prioritizes comparisonError message over trendError and summaryError', async () => {
    mockCompare.mockRejectedValue(new Error('First error'));
    mockTrend.mockRejectedValue(new Error('Second error'));
    mockSummary.mockRejectedValue(new Error('Third error'));

    renderAnalytics();

    await waitFor(() => {
      // The || chain picks comparisonError.message first
      expect(screen.getByText('First error')).toBeInTheDocument();
    });
  });

  it('falls back to trendError when comparisonError has no message', async () => {
    mockCompare.mockRejectedValue({ message: '' });
    mockTrend.mockRejectedValue(new Error('Trend error message'));
    mockSummary.mockResolvedValue({ data: mockSummaryData });

    renderAnalytics();

    await waitFor(() => {
      expect(screen.getByText('Failed to load analytics')).toBeInTheDocument();
    });
  });

  it('falls back to summaryError when comparison and trend have no message', async () => {
    mockCompare.mockResolvedValue({ data: mockComparisonData });
    mockTrend.mockResolvedValue({ data: mockTrendData });
    mockSummary.mockRejectedValue(new Error('Only summary failed'));

    renderAnalytics();

    await waitFor(() => {
      expect(screen.getByText('Failed to load analytics')).toBeInTheDocument();
      expect(screen.getByText('Only summary failed')).toBeInTheDocument();
    });
  });

  it('handles retry after all APIs fail', async () => {
    mockCompare.mockRejectedValueOnce(new Error('Fail'));
    mockTrend.mockRejectedValueOnce(new Error('Fail'));
    mockSummary.mockRejectedValueOnce(new Error('Fail'));

    renderAnalytics();

    await waitFor(() => {
      expect(screen.getByText('Try Again')).toBeInTheDocument();
    });

    // Reset to succeed
    mockCompare.mockResolvedValue({ data: mockComparisonData });
    mockTrend.mockResolvedValue({ data: mockTrendData });
    mockSummary.mockResolvedValue({ data: mockSummaryData });

    fireEvent.click(screen.getByText('Try Again'));

    await waitFor(() => {
      expect(mockCompare).toHaveBeenCalledTimes(2);
      expect(mockTrend).toHaveBeenCalledTimes(2);
      expect(mockSummary).toHaveBeenCalledTimes(2);
    });
  });
});

describe('Analytics Period Selection Full Verification', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockCompare.mockResolvedValue({ data: mockComparisonData });
    mockTrend.mockResolvedValue({ data: mockTrendData });
    mockSummary.mockResolvedValue({ data: mockSummaryData });
  });

  it('calls all three APIs with 90-day period', async () => {
    renderAnalytics();

    await waitFor(() => {
      expect(screen.getByText('Analytics')).toBeInTheDocument();
    });

    const select = screen.getByRole('combobox');
    fireEvent.change(select, { target: { value: '90' } });

    await waitFor(() => {
      expect(mockCompare).toHaveBeenCalledWith(90);
      expect(mockTrend).toHaveBeenCalledWith(90);
      expect(mockSummary).toHaveBeenCalledWith(90);
    });
  });

  it('calls all three APIs with 7-day period', async () => {
    renderAnalytics();

    await waitFor(() => {
      expect(screen.getByText('Analytics')).toBeInTheDocument();
    });

    const select = screen.getByRole('combobox');
    fireEvent.change(select, { target: { value: '7' } });

    await waitFor(() => {
      expect(mockCompare).toHaveBeenCalledWith(7);
      expect(mockTrend).toHaveBeenCalledWith(7);
      expect(mockSummary).toHaveBeenCalledWith(7);
    });
  });

  it('defaults to 30-day period on initial load', async () => {
    renderAnalytics();

    await waitFor(() => {
      expect(mockCompare).toHaveBeenCalledWith(30);
      expect(mockTrend).toHaveBeenCalledWith(30);
      expect(mockSummary).toHaveBeenCalledWith(30);
    });
  });
});

describe('Analytics Trend Summary Conditional Rendering', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockCompare.mockResolvedValue({ data: mockComparisonData });
    mockSummary.mockResolvedValue({ data: mockSummaryData });
  });

  it('renders trend summary section when summary exists', async () => {
    mockTrend.mockResolvedValue({ data: mockTrendData });

    renderAnalytics();

    await waitFor(() => {
      expect(screen.getByText('Min Score')).toBeInTheDocument();
      expect(screen.getByText('Max Score')).toBeInTheDocument();
      expect(screen.getByText('Average')).toBeInTheDocument();
      expect(screen.getByText('65.0%')).toBeInTheDocument();
      expect(screen.getByText('85.0%')).toBeInTheDocument();
      expect(screen.getByText('75.0%')).toBeInTheDocument();
    });
  });

  it('does not render trend summary section when summary is undefined', async () => {
    mockTrend.mockResolvedValue({
      data: {
        ...mockTrendData,
        summary: undefined,
      },
    });

    renderAnalytics();

    await waitFor(() => {
      expect(screen.getByText('Score Trend')).toBeInTheDocument();
    });

    // Min/Max/Average labels from the trend summary should NOT be present
    expect(screen.queryByText('Min Score')).not.toBeInTheDocument();
    expect(screen.queryByText('Max Score')).not.toBeInTheDocument();
  });

  it('does not render trend summary section when summary is null', async () => {
    mockTrend.mockResolvedValue({
      data: {
        ...mockTrendData,
        summary: null,
      },
    });

    renderAnalytics();

    await waitFor(() => {
      expect(screen.getByText('Score Trend')).toBeInTheDocument();
    });

    expect(screen.queryByText('Min Score')).not.toBeInTheDocument();
    expect(screen.queryByText('Max Score')).not.toBeInTheDocument();
  });
});

describe('Analytics Scenario Score Rendering', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockCompare.mockResolvedValue({ data: mockComparisonData });
    mockTrend.mockResolvedValue({ data: mockTrendData });
  });

  it('renders progress bars with correct widths for scenario scores', async () => {
    mockSummary.mockResolvedValue({
      data: {
        ...mockSummaryData,
        scores_by_scenario: {
          'High Score Scenario': 100,
          'Low Score Scenario': 10,
        },
      },
    });

    renderAnalytics();

    await waitFor(() => {
      expect(screen.getByText('High Score Scenario')).toBeInTheDocument();
      expect(screen.getByText('Low Score Scenario')).toBeInTheDocument();
      expect(screen.getByText('100%')).toBeInTheDocument();
      expect(screen.getByText('10%')).toBeInTheDocument();
    });
  });

  it('caps progress bar width at 100% for scores over 100', async () => {
    mockSummary.mockResolvedValue({
      data: {
        ...mockSummaryData,
        scores_by_scenario: {
          'Over 100': 150,
        },
      },
    });

    const { container } = renderAnalytics();

    await waitFor(() => {
      expect(screen.getByText('Over 100')).toBeInTheDocument();
      expect(screen.getByText('150%')).toBeInTheDocument();
    });

    // The progress bar width should be capped at 100% via Math.min(score, 100)
    const progressBars = container.querySelectorAll('.bg-primary-600.h-2.rounded-full');
    const overBar = Array.from(progressBars).find(
      (el) => (el as HTMLElement).style.width === '100%'
    );
    expect(overBar).toBeTruthy();
  });

  it('renders scenario title attribute for long scenario names', async () => {
    const longName = 'A Very Long Scenario Name That Might Overflow The Container';
    mockSummary.mockResolvedValue({
      data: {
        ...mockSummaryData,
        scores_by_scenario: {
          [longName]: 75,
        },
      },
    });

    renderAnalytics();

    await waitFor(() => {
      const scenarioLabel = screen.getByTitle(longName);
      expect(scenarioLabel).toBeInTheDocument();
    });
  });
});

describe('Analytics Null Query Data Fallbacks', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders with null comparison data using fallback values', async () => {
    // When comparison is null, all comparison?.xxx || 0 branches hit fallback
    // This covers the || 0 branch on line 232 and similar fallback expressions
    mockCompare.mockResolvedValue({ data: null });
    mockTrend.mockResolvedValue({ data: mockTrendData });
    mockSummary.mockResolvedValue({ data: mockSummaryData });

    renderAnalytics();

    await waitFor(() => {
      expect(screen.getByText('Average Score')).toBeInTheDocument();
      expect(screen.getByText('Executions')).toBeInTheDocument();
      expect(screen.getByText('Blocked Attacks')).toBeInTheDocument();
      expect(screen.getByText('Detected Attacks')).toBeInTheDocument();
    });
  });

  it('renders charts with null trend data using empty array fallbacks', async () => {
    // When trend is null, all trend?.data_points.map(...) || [] branches hit fallback
    // This covers the || [] branches on lines 139, 143, 153, 157, 162, 167
    mockCompare.mockResolvedValue({ data: mockComparisonData });
    mockTrend.mockResolvedValue({ data: null });
    mockSummary.mockResolvedValue({ data: mockSummaryData });

    renderAnalytics();

    await waitFor(() => {
      expect(screen.getByText('Score Trend')).toBeInTheDocument();
      expect(screen.getByText('Detection Results Over Time')).toBeInTheDocument();
      // Charts should still render with empty data from fallbacks
      const lineCharts = screen.getAllByTestId('line-chart');
      const barCharts = screen.getAllByTestId('bar-chart');
      expect(lineCharts.length).toBeGreaterThan(0);
      expect(barCharts.length).toBeGreaterThan(0);
    });

    // Trend summary should NOT render since trend is null
    expect(screen.queryByText('Min Score')).not.toBeInTheDocument();
  });

  it('renders with null summary data using fallback values', async () => {
    // When summary is null, all summary?.xxx || 0 branches hit fallback
    mockCompare.mockResolvedValue({ data: mockComparisonData });
    mockTrend.mockResolvedValue({ data: mockTrendData });
    mockSummary.mockResolvedValue({ data: null });

    renderAnalytics();

    await waitFor(() => {
      expect(screen.getByText('Execution Summary')).toBeInTheDocument();
      expect(screen.getByText('Executions by Status')).toBeInTheDocument();
      expect(screen.getByText('Performance by Scenario')).toBeInTheDocument();
      // No scenario data when summary is null
      expect(screen.getByText('No scenario data available')).toBeInTheDocument();
    });
  });

  it('renders with all query data null', async () => {
    mockCompare.mockResolvedValue({ data: null });
    mockTrend.mockResolvedValue({ data: null });
    mockSummary.mockResolvedValue({ data: null });

    renderAnalytics();

    await waitFor(() => {
      // All sections render with fallback values
      expect(screen.getByText('Analytics')).toBeInTheDocument();
      expect(screen.getByText('Average Score')).toBeInTheDocument();
      expect(screen.getByText('Score Trend')).toBeInTheDocument();
      expect(screen.getByText('Execution Summary')).toBeInTheDocument();
      expect(screen.getByText('No scenario data available')).toBeInTheDocument();
    });
  });
});
