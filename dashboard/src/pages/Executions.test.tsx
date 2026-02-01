import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import Executions from './Executions';
import { api } from '../lib/api';

// Mock the API
vi.mock('../lib/api', () => ({
  api: {
    get: vi.fn(),
  },
}));

// Mock date-fns
vi.mock('date-fns', () => ({
  formatDistanceToNow: vi.fn(() => '2 hours ago'),
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

describe('Executions Page', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders loading state', () => {
    vi.mocked(api.get).mockReturnValue(new Promise(() => {}) as never);

    renderWithClient(<Executions />);
    expect(screen.getByText('Loading executions...')).toBeInTheDocument();
  });

  it('renders executions list with completed status', async () => {
    const mockExecutions = [
      {
        id: 'exec-12345678-abcd',
        scenario_id: 'scenario-1',
        status: 'completed',
        started_at: '2024-01-15T10:00:00Z',
        completed_at: '2024-01-15T10:30:00Z',
        safe_mode: true,
        score: {
          overall: 75.5,
          blocked: 5,
          detected: 3,
          successful: 2,
          total: 10,
        },
      },
    ];
    vi.mocked(api.get).mockResolvedValue({ data: mockExecutions } as never);

    renderWithClient(<Executions />);

    expect(await screen.findByText('scenario-1')).toBeInTheDocument();
    expect(screen.getByText('completed')).toBeInTheDocument();
    expect(screen.getByText('75.5%')).toBeInTheDocument();
    expect(screen.getByText('5 blocked')).toBeInTheDocument();
    expect(screen.getByText('3 detected')).toBeInTheDocument();
    expect(screen.getByText('2 success')).toBeInTheDocument();
    expect(screen.getByText('Safe')).toBeInTheDocument();
  });

  it('renders running execution with warning badge', async () => {
    const mockExecutions = [
      {
        id: 'exec-running',
        scenario_id: 'scenario-2',
        status: 'running',
        started_at: '2024-01-15T12:00:00Z',
        safe_mode: false,
      },
    ];
    vi.mocked(api.get).mockResolvedValue({ data: mockExecutions } as never);

    renderWithClient(<Executions />);

    expect(await screen.findByText('running')).toBeInTheDocument();
    expect(screen.getByText('Full')).toBeInTheDocument();
    expect(screen.getByText('-%')).toBeInTheDocument();
  });

  it('renders failed execution with danger badge', async () => {
    const mockExecutions = [
      {
        id: 'exec-failed',
        scenario_id: 'scenario-3',
        status: 'failed',
        started_at: '2024-01-15T11:00:00Z',
        safe_mode: true,
      },
    ];
    vi.mocked(api.get).mockResolvedValue({ data: mockExecutions } as never);

    renderWithClient(<Executions />);

    expect(await screen.findByText('failed')).toBeInTheDocument();
  });

  it('renders empty state when no executions', async () => {
    vi.mocked(api.get).mockResolvedValue({ data: [] } as never);

    renderWithClient(<Executions />);

    expect(await screen.findByText('No executions yet')).toBeInTheDocument();
    expect(screen.getByText('Run a scenario to see results here')).toBeInTheDocument();
  });

  it('renders page title and new execution button', async () => {
    vi.mocked(api.get).mockResolvedValue({ data: [] } as never);

    renderWithClient(<Executions />);

    expect(await screen.findByText('Executions')).toBeInTheDocument();
    expect(screen.getByText('New Execution')).toBeInTheDocument();
  });

  it('renders truncated execution ID', async () => {
    const mockExecutions = [
      {
        id: 'abcdefgh-1234-5678-9012-ijklmnopqrst',
        scenario_id: 'test-scenario',
        status: 'completed',
        started_at: '2024-01-15T10:00:00Z',
        safe_mode: true,
        score: { overall: 80, blocked: 4, detected: 1, successful: 0, total: 5 },
      },
    ];
    vi.mocked(api.get).mockResolvedValue({ data: mockExecutions } as never);

    renderWithClient(<Executions />);

    expect(await screen.findByText('abcdefgh...')).toBeInTheDocument();
  });

  it('handles execution without score', async () => {
    const mockExecutions = [
      {
        id: 'exec-no-score',
        scenario_id: 'pending-scenario',
        status: 'pending',
        started_at: '2024-01-15T14:00:00Z',
        safe_mode: true,
      },
    ];
    vi.mocked(api.get).mockResolvedValue({ data: mockExecutions } as never);

    renderWithClient(<Executions />);

    expect(await screen.findByText('pending-scenario')).toBeInTheDocument();
    expect(screen.getByText('-%')).toBeInTheDocument();
    expect(screen.getByText('0 blocked')).toBeInTheDocument();
    expect(screen.getByText('0 detected')).toBeInTheDocument();
    expect(screen.getByText('0 success')).toBeInTheDocument();
  });
});
