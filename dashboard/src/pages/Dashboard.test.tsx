import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import Dashboard from './Dashboard';
import { api } from '../lib/api';

// Mock the API
vi.mock('../lib/api', () => ({
  api: {
    get: vi.fn(),
  },
}));

// Mock chart.js to avoid canvas issues
vi.mock('react-chartjs-2', () => ({
  Doughnut: () => <div data-testid="doughnut-chart">Chart</div>,
}));

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
  });

  it('renders dashboard title', async () => {
    vi.mocked(api.get).mockResolvedValue({ data: [] } as never);

    renderWithClient(<Dashboard />);
    expect(screen.getByText('Dashboard')).toBeInTheDocument();
  });

  it('renders stats cards', async () => {
    vi.mocked(api.get).mockResolvedValue({ data: [] } as never);

    renderWithClient(<Dashboard />);
    expect(screen.getByText('Agents Online')).toBeInTheDocument();
    expect(screen.getByText('Security Score')).toBeInTheDocument();
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

    renderWithClient(<Dashboard />);

    expect(await screen.findByText('85.5%')).toBeInTheDocument();
    expect(screen.getByText('10')).toBeInTheDocument(); // total techniques
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

    renderWithClient(<Dashboard />);

    expect(await screen.findByText('Initial Access Test')).toBeInTheDocument();
    expect(screen.getByText('Lateral Movement')).toBeInTheDocument();
    expect(screen.getByText('completed')).toBeInTheDocument();
    expect(screen.getByText('running')).toBeInTheDocument();
  });

  it('renders chart section', async () => {
    vi.mocked(api.get).mockResolvedValue({ data: [] } as never);

    renderWithClient(<Dashboard />);

    expect(screen.getByText('Detection Results')).toBeInTheDocument();
    expect(screen.getByText('Recent Activity')).toBeInTheDocument();
    expect(screen.getByTestId('doughnut-chart')).toBeInTheDocument();
  });

  it('handles empty state gracefully', async () => {
    vi.mocked(api.get).mockResolvedValue({ data: [] } as never);

    renderWithClient(<Dashboard />);

    // Default values when no data
    expect(await screen.findByText('0%')).toBeInTheDocument();
    const zeros = screen.getAllByText('0');
    expect(zeros.length).toBeGreaterThan(0);
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

    renderWithClient(<Dashboard />);

    // Only first 5 should be shown
    expect(await screen.findByText('Scenario 1')).toBeInTheDocument();
    expect(screen.getByText('Scenario 5')).toBeInTheDocument();
    expect(screen.queryByText('Scenario 6')).not.toBeInTheDocument();
  });
});
