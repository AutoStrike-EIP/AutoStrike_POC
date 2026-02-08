import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { render, screen } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import Dashboard from './Dashboard';
import { api, techniqueApi } from '../lib/api';

// Mock the API
vi.mock('../lib/api', () => ({
  api: {
    get: vi.fn(),
  },
  techniqueApi: {
    getCoverage: vi.fn(),
  },
}));

// Mock chart.js to avoid canvas issues
vi.mock('react-chartjs-2', () => ({
  Doughnut: () => <div data-testid="doughnut-chart">Chart</div>,
}));

// Mock matchMedia for SecurityScore animation
const matchMediaMock = vi.fn().mockReturnValue({
  matches: true, // Prefer reduced motion to disable animation
  addEventListener: vi.fn(),
  removeEventListener: vi.fn(),
});

const createTestQueryClient = () =>
  new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
      },
    },
  });

function renderWithClient(ui: React.ReactElement) {
  const testQueryClient = createTestQueryClient();
  return render(
    <QueryClientProvider client={testQueryClient}>{ui}</QueryClientProvider>
  );
}

describe('Dashboard', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.stubGlobal('matchMedia', matchMediaMock);
  });

  afterEach(() => {
    vi.unstubAllGlobals();
  });

  it('renders dashboard title', async () => {
    vi.mocked(api.get).mockResolvedValue({ data: [] } as never);
    vi.mocked(techniqueApi.getCoverage).mockResolvedValue({ data: {} } as never);

    renderWithClient(<Dashboard />);
    expect(screen.getByText('Dashboard')).toBeInTheDocument();
  });

  it('renders stats cards', async () => {
    vi.mocked(api.get).mockResolvedValue({ data: [] } as never);
    vi.mocked(techniqueApi.getCoverage).mockResolvedValue({ data: {} } as never);

    renderWithClient(<Dashboard />);
    expect(screen.getByText('Agents Online')).toBeInTheDocument();
    // SecurityScore label appears in both stat card and main section
    expect(screen.getAllByText('Security Score').length).toBeGreaterThanOrEqual(1);
    expect(screen.getByText('Techniques Tested')).toBeInTheDocument();
    expect(screen.getByText('Executions Today')).toBeInTheDocument();
  });

  it('renders agent stats correctly', async () => {
    const mockAgents = [
      { paw: 'agent-1', hostname: 'PC1', status: 'online' },
      { paw: 'agent-2', hostname: 'PC2', status: 'offline' },
      { paw: 'agent-3', hostname: 'PC3', status: 'online' },
    ];
    vi.mocked(api.get).mockImplementation((url) => {
      if (url === '/agents') {
        return Promise.resolve({ data: mockAgents }) as never;
      }
      return Promise.resolve({ data: [] }) as never;
    });
    vi.mocked(techniqueApi.getCoverage).mockResolvedValue({ data: {} } as never);

    renderWithClient(<Dashboard />);

    expect(await screen.findByText('2')).toBeInTheDocument(); // online agents
    expect(screen.getByText('of 3 total')).toBeInTheDocument();
  });

  it('renders latest execution score', async () => {
    const mockExecutions = [
      {
        id: 'exec-1',
        scenario_id: 'scenario-1',
        status: 'completed',
        started_at: '2024-01-15T10:00:00Z',
        score: {
          overall: 85.5,
          blocked: 5,
          detected: 3,
          successful: 2,
          total: 10,
        },
      },
    ];
    vi.mocked(api.get).mockImplementation((url) => {
      if (url === '/executions') {
        return Promise.resolve({ data: mockExecutions }) as never;
      }
      return Promise.resolve({ data: [] }) as never;
    });
    vi.mocked(techniqueApi.getCoverage).mockResolvedValue({ data: {} } as never);

    renderWithClient(<Dashboard />);

    // SecurityScore displays score in multiple places (stat card + main section)
    const scores = await screen.findAllByText('85.5');
    expect(scores.length).toBe(2); // stat card + main SecurityScore
    expect(screen.getByText('10')).toBeInTheDocument(); // total techniques in breakdown
  });

  it('renders recent activity with executions', async () => {
    const mockExecutions = [
      {
        id: 'exec-1',
        scenario_id: 'Initial Access Test',
        status: 'completed',
        started_at: '2024-01-15T10:00:00Z',
        score: { overall: 80, blocked: 4, detected: 1, successful: 0, total: 5 },
      },
      {
        id: 'exec-2',
        scenario_id: 'Lateral Movement',
        status: 'running',
        started_at: '2024-01-15T11:00:00Z',
      },
    ];
    vi.mocked(api.get).mockImplementation((url) => {
      if (url === '/executions') {
        return Promise.resolve({ data: mockExecutions }) as never;
      }
      return Promise.resolve({ data: [] }) as never;
    });
    vi.mocked(techniqueApi.getCoverage).mockResolvedValue({ data: {} } as never);

    renderWithClient(<Dashboard />);

    expect(await screen.findByText('Initial Access Test')).toBeInTheDocument();
    expect(screen.getByText('Lateral Movement')).toBeInTheDocument();
    expect(screen.getByText('completed')).toBeInTheDocument();
    expect(screen.getByText('running')).toBeInTheDocument();
  });

  it('renders chart section', async () => {
    vi.mocked(api.get).mockResolvedValue({ data: [] } as never);
    vi.mocked(techniqueApi.getCoverage).mockResolvedValue({ data: {} } as never);

    renderWithClient(<Dashboard />);

    expect(screen.getByText('Detection Results')).toBeInTheDocument();
    expect(screen.getByText('Recent Activity')).toBeInTheDocument();
    expect(screen.getByTestId('doughnut-chart')).toBeInTheDocument();
  });

  it('handles empty state gracefully', async () => {
    vi.mocked(api.get).mockResolvedValue({ data: [] } as never);
    vi.mocked(techniqueApi.getCoverage).mockResolvedValue({ data: {} } as never);

    renderWithClient(<Dashboard />);

    // SecurityScore shows 0.0 in multiple places (stat card + main section)
    const zeroScores = await screen.findAllByText('0.0');
    expect(zeroScores.length).toBe(2);
    // Check for zeros in stat cards
    const zeros = screen.getAllByText('0');
    expect(zeros.length).toBeGreaterThan(0);
    // Check for "No recent executions" message
    expect(screen.getByText('No recent executions')).toBeInTheDocument();
  });

  it('renders up to 5 recent executions', async () => {
    const mockExecutions = [
      { id: 'exec-1', scenario_id: 'Scenario 1', status: 'completed', started_at: '2024-01-15T10:00:00Z' },
      { id: 'exec-2', scenario_id: 'Scenario 2', status: 'completed', started_at: '2024-01-15T11:00:00Z' },
      { id: 'exec-3', scenario_id: 'Scenario 3', status: 'completed', started_at: '2024-01-15T12:00:00Z' },
      { id: 'exec-4', scenario_id: 'Scenario 4', status: 'completed', started_at: '2024-01-15T13:00:00Z' },
      { id: 'exec-5', scenario_id: 'Scenario 5', status: 'completed', started_at: '2024-01-15T14:00:00Z' },
      { id: 'exec-6', scenario_id: 'Scenario 6', status: 'completed', started_at: '2024-01-15T15:00:00Z' },
    ];
    vi.mocked(api.get).mockImplementation((url) => {
      if (url === '/executions') {
        return Promise.resolve({ data: mockExecutions }) as never;
      }
      return Promise.resolve({ data: [] }) as never;
    });
    vi.mocked(techniqueApi.getCoverage).mockResolvedValue({ data: {} } as never);

    renderWithClient(<Dashboard />);

    // Only first 5 should be shown
    expect(await screen.findByText('Scenario 1')).toBeInTheDocument();
    expect(screen.getByText('Scenario 5')).toBeInTheDocument();
    expect(screen.queryByText('Scenario 6')).not.toBeInTheDocument();
  });

  it('renders MITRE Coverage section', async () => {
    vi.mocked(api.get).mockResolvedValue({ data: [] } as never);
    vi.mocked(techniqueApi.getCoverage).mockResolvedValue({ data: {} } as never);

    renderWithClient(<Dashboard />);

    expect(screen.getByText('MITRE Coverage')).toBeInTheDocument();
  });

  it('displays coverage data when available', async () => {
    const mockCoverage = {
      discovery: 9,
      execution: 3,
    };
    const mockTechniques = Array.from({ length: 12 }, (_, i) => ({ id: `T${1000 + i}`, name: `Tech ${i}` }));
    vi.mocked(api.get).mockImplementation((url) => {
      if (url === '/techniques') {
        return Promise.resolve({ data: mockTechniques }) as never;
      }
      return Promise.resolve({ data: [] }) as never;
    });
    vi.mocked(techniqueApi.getCoverage).mockResolvedValue({ data: mockCoverage } as never);

    renderWithClient(<Dashboard />);

    // Total techniques should be 12 (from techniques array length)
    expect(await screen.findByText('12')).toBeInTheDocument();
    expect(screen.getByText('discovery')).toBeInTheDocument();
    expect(screen.getByText('execution')).toBeInTheDocument();
  });
});
