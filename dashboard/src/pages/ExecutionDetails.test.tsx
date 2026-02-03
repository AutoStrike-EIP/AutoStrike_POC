import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import ExecutionDetails from './ExecutionDetails';
import { executionApi } from '../lib/api';

// Mock the API
vi.mock('../lib/api', () => ({
  api: {
    get: vi.fn(),
  },
  executionApi: {
    get: vi.fn(),
    getResults: vi.fn(),
  },
}));

// Mock date-fns
vi.mock('date-fns', () => ({
  formatDistanceToNow: vi.fn(() => '2 hours ago'),
  format: vi.fn(() => 'Jan 15, 2024, 10:00 AM'),
}));

const createTestQueryClient = () =>
  new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
      },
    },
  });

function renderWithRouter(executionId: string) {
  const testQueryClient = createTestQueryClient();
  return render(
    <QueryClientProvider client={testQueryClient}>
      <MemoryRouter initialEntries={[`/executions/${executionId}`]}>
        <Routes>
          <Route path="/executions/:id" element={<ExecutionDetails />} />
          <Route path="/executions" element={<div>Executions List</div>} />
        </Routes>
      </MemoryRouter>
    </QueryClientProvider>
  );
}

describe('ExecutionDetails Page', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders loading state', () => {
    vi.mocked(executionApi.get).mockReturnValue(new Promise(() => {}) as never);
    vi.mocked(executionApi.getResults).mockReturnValue(new Promise(() => {}) as never);

    renderWithRouter('exec-123');
    expect(screen.getByText('Loading execution details...')).toBeInTheDocument();
  });

  it('renders execution not found', async () => {
    vi.mocked(executionApi.get).mockResolvedValue({ data: null } as never);
    vi.mocked(executionApi.getResults).mockResolvedValue({ data: [] } as never);

    renderWithRouter('invalid-id');

    expect(await screen.findByText('Execution not found')).toBeInTheDocument();
    expect(screen.getByText('Back to Executions')).toBeInTheDocument();
  });

  it('renders execution details with completed status', async () => {
    const mockExecution = {
      id: 'exec-12345678-abcd',
      scenario_id: 'test-scenario',
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
    };

    const mockResults = [
      {
        id: 'result-1',
        execution_id: 'exec-12345678-abcd',
        technique_id: 'T1082',
        agent_paw: 'agent-1',
        status: 'success',
        output: 'System info collected',
        detected: false,
        start_time: '2024-01-15T10:00:00Z',
        end_time: '2024-01-15T10:01:00Z',
      },
    ];

    vi.mocked(executionApi.get).mockResolvedValue({ data: mockExecution } as never);
    vi.mocked(executionApi.getResults).mockResolvedValue({ data: mockResults } as never);

    renderWithRouter('exec-12345678-abcd');

    // Check header
    expect(await screen.findByText('Execution Details')).toBeInTheDocument();
    expect(screen.getByText('exec-12345678-abcd')).toBeInTheDocument();
    expect(screen.getByText('completed')).toBeInTheDocument();

    // Check summary
    expect(screen.getByText('test-scenario')).toBeInTheDocument();
    expect(screen.getByText('Safe Mode')).toBeInTheDocument();
    expect(screen.getByText('75.5%')).toBeInTheDocument();

    // Check score breakdown
    expect(screen.getByText('5')).toBeInTheDocument(); // blocked
    expect(screen.getByText('3')).toBeInTheDocument(); // detected
    expect(screen.getByText('2')).toBeInTheDocument(); // successful
    expect(screen.getByText('10')).toBeInTheDocument(); // total

    // Check results table
    expect(screen.getByText('T1082')).toBeInTheDocument();
    expect(screen.getByText('agent-1')).toBeInTheDocument();
  });

  it('renders execution with running status', async () => {
    const mockExecution = {
      id: 'exec-running',
      scenario_id: 'running-scenario',
      status: 'running',
      started_at: '2024-01-15T12:00:00Z',
      safe_mode: false,
    };

    vi.mocked(executionApi.get).mockResolvedValue({ data: mockExecution } as never);
    vi.mocked(executionApi.getResults).mockResolvedValue({ data: [] } as never);

    renderWithRouter('exec-running');

    expect(await screen.findByText('running')).toBeInTheDocument();
    expect(screen.getByText('Full Mode')).toBeInTheDocument();
    expect(screen.getByText('0%')).toBeInTheDocument();
  });

  it('renders empty results state', async () => {
    const mockExecution = {
      id: 'exec-pending',
      scenario_id: 'pending-scenario',
      status: 'pending',
      started_at: '2024-01-15T14:00:00Z',
      safe_mode: true,
    };

    vi.mocked(executionApi.get).mockResolvedValue({ data: mockExecution } as never);
    vi.mocked(executionApi.getResults).mockResolvedValue({ data: [] } as never);

    renderWithRouter('exec-pending');

    expect(await screen.findByText('No results yet')).toBeInTheDocument();
    expect(screen.getByText('Results will appear here as techniques are executed')).toBeInTheDocument();
  });

  it('renders back button that navigates to executions', async () => {
    const mockExecution = {
      id: 'exec-test',
      scenario_id: 'test-scenario',
      status: 'completed',
      started_at: '2024-01-15T10:00:00Z',
      safe_mode: true,
      score: { overall: 100, blocked: 5, detected: 0, successful: 0, total: 5 },
    };

    vi.mocked(executionApi.get).mockResolvedValue({ data: mockExecution } as never);
    vi.mocked(executionApi.getResults).mockResolvedValue({ data: [] } as never);

    renderWithRouter('exec-test');

    expect(await screen.findByText('Back to Executions')).toBeInTheDocument();
  });

  it('renders result with failed status', async () => {
    const mockExecution = {
      id: 'exec-failed',
      scenario_id: 'failed-scenario',
      status: 'completed',
      started_at: '2024-01-15T10:00:00Z',
      safe_mode: true,
      score: { overall: 100, blocked: 1, detected: 0, successful: 0, total: 1 },
    };

    const mockResults = [
      {
        id: 'result-failed',
        execution_id: 'exec-failed',
        technique_id: 'T1059',
        agent_paw: 'agent-1',
        status: 'failed',
        output: '',
        detected: false,
        start_time: '2024-01-15T10:00:00Z',
        end_time: '2024-01-15T10:01:00Z',
      },
    ];

    vi.mocked(executionApi.get).mockResolvedValue({ data: mockExecution } as never);
    vi.mocked(executionApi.getResults).mockResolvedValue({ data: mockResults } as never);

    renderWithRouter('exec-failed');

    expect(await screen.findByText('Execution Failed')).toBeInTheDocument();
  });

  it('renders result with detected status', async () => {
    const mockExecution = {
      id: 'exec-detected',
      scenario_id: 'detected-scenario',
      status: 'completed',
      started_at: '2024-01-15T10:00:00Z',
      safe_mode: true,
      score: { overall: 50, blocked: 0, detected: 1, successful: 0, total: 1 },
    };

    const mockResults = [
      {
        id: 'result-detected',
        execution_id: 'exec-detected',
        technique_id: 'T1016',
        agent_paw: 'agent-1',
        status: 'detected',
        output: 'Network info collected but detected',
        detected: true,
        start_time: '2024-01-15T10:00:00Z',
        end_time: '2024-01-15T10:01:00Z',
      },
    ];

    vi.mocked(executionApi.get).mockResolvedValue({ data: mockExecution } as never);
    vi.mocked(executionApi.getResults).mockResolvedValue({ data: mockResults } as never);

    renderWithRouter('exec-detected');

    // Wait for results to load and check technique ID is rendered
    expect(await screen.findByText('T1016')).toBeInTheDocument();
  });

  it('renders result with pending status', async () => {
    const mockExecution = {
      id: 'exec-with-pending',
      scenario_id: 'pending-scenario',
      status: 'running',
      started_at: '2024-01-15T10:00:00Z',
      safe_mode: true,
    };

    const mockResults = [
      {
        id: 'result-pending',
        execution_id: 'exec-with-pending',
        technique_id: 'T1049',
        agent_paw: 'agent-1',
        status: 'pending',
        output: '',
        detected: false,
        start_time: '2024-01-15T10:00:00Z',
        end_time: '',
      },
    ];

    vi.mocked(executionApi.get).mockResolvedValue({ data: mockExecution } as never);
    vi.mocked(executionApi.getResults).mockResolvedValue({ data: mockResults } as never);

    renderWithRouter('exec-with-pending');

    expect(await screen.findByText('Pending')).toBeInTheDocument();
  });

  it('renders result output when available', async () => {
    const mockExecution = {
      id: 'exec-output',
      scenario_id: 'output-scenario',
      status: 'completed',
      started_at: '2024-01-15T10:00:00Z',
      safe_mode: true,
      score: { overall: 0, blocked: 0, detected: 0, successful: 1, total: 1 },
    };

    const mockResults = [
      {
        id: 'result-with-output',
        execution_id: 'exec-output',
        technique_id: 'T1082',
        agent_paw: 'agent-1',
        status: 'success',
        output: 'Linux server-01 5.4.0-generic',
        detected: false,
        start_time: '2024-01-15T10:00:00Z',
        end_time: '2024-01-15T10:01:00Z',
      },
    ];

    vi.mocked(executionApi.get).mockResolvedValue({ data: mockExecution } as never);
    vi.mocked(executionApi.getResults).mockResolvedValue({ data: mockResults } as never);

    renderWithRouter('exec-output');

    expect(await screen.findByText('View output')).toBeInTheDocument();
  });

  it('renders no output message when output is empty', async () => {
    const mockExecution = {
      id: 'exec-no-output',
      scenario_id: 'no-output-scenario',
      status: 'completed',
      started_at: '2024-01-15T10:00:00Z',
      safe_mode: true,
      score: { overall: 100, blocked: 1, detected: 0, successful: 0, total: 1 },
    };

    const mockResults = [
      {
        id: 'result-no-output',
        execution_id: 'exec-no-output',
        technique_id: 'T1059',
        agent_paw: 'agent-1',
        status: 'blocked',
        output: '',
        detected: false,
        start_time: '2024-01-15T10:00:00Z',
        end_time: '2024-01-15T10:01:00Z',
      },
    ];

    vi.mocked(executionApi.get).mockResolvedValue({ data: mockExecution } as never);
    vi.mocked(executionApi.getResults).mockResolvedValue({ data: mockResults } as never);

    renderWithRouter('exec-no-output');

    expect(await screen.findByText('No output')).toBeInTheDocument();
  });

  it('renders cancelled execution status', async () => {
    const mockExecution = {
      id: 'exec-cancelled',
      scenario_id: 'cancelled-scenario',
      status: 'cancelled',
      started_at: '2024-01-15T10:00:00Z',
      safe_mode: true,
    };

    vi.mocked(executionApi.get).mockResolvedValue({ data: mockExecution } as never);
    vi.mocked(executionApi.getResults).mockResolvedValue({ data: [] } as never);

    renderWithRouter('exec-cancelled');

    expect(await screen.findByText('cancelled')).toBeInTheDocument();
  });

  it('renders multiple results', async () => {
    const mockExecution = {
      id: 'exec-multi',
      scenario_id: 'multi-scenario',
      status: 'completed',
      started_at: '2024-01-15T10:00:00Z',
      safe_mode: true,
      score: { overall: 66.7, blocked: 1, detected: 1, successful: 1, total: 3 },
    };

    const mockResults = [
      {
        id: 'result-1',
        execution_id: 'exec-multi',
        technique_id: 'T1082',
        agent_paw: 'agent-1',
        status: 'success',
        output: 'Output 1',
        detected: false,
        start_time: '2024-01-15T10:00:00Z',
        end_time: '2024-01-15T10:01:00Z',
      },
      {
        id: 'result-2',
        execution_id: 'exec-multi',
        technique_id: 'T1016',
        agent_paw: 'agent-1',
        status: 'detected',
        output: 'Output 2',
        detected: true,
        start_time: '2024-01-15T10:01:00Z',
        end_time: '2024-01-15T10:02:00Z',
      },
      {
        id: 'result-3',
        execution_id: 'exec-multi',
        technique_id: 'T1059',
        agent_paw: 'agent-1',
        status: 'blocked',
        output: '',
        detected: false,
        start_time: '2024-01-15T10:02:00Z',
        end_time: '2024-01-15T10:03:00Z',
      },
    ];

    vi.mocked(executionApi.get).mockResolvedValue({ data: mockExecution } as never);
    vi.mocked(executionApi.getResults).mockResolvedValue({ data: mockResults } as never);

    renderWithRouter('exec-multi');

    expect(await screen.findByText('T1082')).toBeInTheDocument();
    expect(screen.getByText('T1016')).toBeInTheDocument();
    expect(screen.getByText('T1059')).toBeInTheDocument();
  });
});
