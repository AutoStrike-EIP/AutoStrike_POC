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

  it('renders result with unknown status using default badge and raw label', async () => {
    const mockExecution = {
      id: 'exec-unknown-result',
      scenario_id: 'unknown-scenario',
      status: 'completed',
      started_at: '2024-01-15T10:00:00Z',
      safe_mode: true,
      score: { overall: 0, blocked: 0, detected: 0, successful: 0, total: 1 },
    };

    const mockResults = [
      {
        id: 'result-unknown',
        execution_id: 'exec-unknown-result',
        technique_id: 'T1082',
        agent_paw: 'agent-1',
        status: 'unknown_status',
        output: '',
        detected: false,
        start_time: '2024-01-15T10:00:00Z',
        end_time: '2024-01-15T10:01:00Z',
      },
    ];

    vi.mocked(executionApi.get).mockResolvedValue({ data: mockExecution } as never);
    vi.mocked(executionApi.getResults).mockResolvedValue({ data: mockResults } as never);

    renderWithRouter('exec-unknown-result');

    // getStatusLabel default case returns the raw status string
    expect(await screen.findByText('unknown_status')).toBeInTheDocument();
  });

  it('renders execution with unknown status using default badge', async () => {
    const mockExecution = {
      id: 'exec-unknown-exec',
      scenario_id: 'unknown-exec-scenario',
      status: 'some_unrecognized_status',
      started_at: '2024-01-15T10:00:00Z',
      safe_mode: true,
    };

    vi.mocked(executionApi.get).mockResolvedValue({ data: mockExecution } as never);
    vi.mocked(executionApi.getResults).mockResolvedValue({ data: [] } as never);

    renderWithRouter('exec-unknown-exec');

    // getExecutionStatusBadge default case applies the gray fallback badge
    const badge = await screen.findByText('some_unrecognized_status');
    expect(badge).toBeInTheDocument();
    expect(badge.className).toContain('bg-gray-100');
  });

  it('renders execution with failed status using danger badge', async () => {
    const mockExecution = {
      id: 'exec-failed-status',
      scenario_id: 'failed-exec-scenario',
      status: 'failed',
      started_at: '2024-01-15T10:00:00Z',
      safe_mode: true,
    };

    vi.mocked(executionApi.get).mockResolvedValue({ data: mockExecution } as never);
    vi.mocked(executionApi.getResults).mockResolvedValue({ data: [] } as never);

    renderWithRouter('exec-failed-status');

    // getExecutionStatusBadge 'failed' case returns 'badge-danger'
    const badge = await screen.findByText('failed');
    expect(badge).toBeInTheDocument();
    expect(badge.className).toContain('badge-danger');
  });

  it('renders result with "successful" status alias (Attack Succeeded label)', async () => {
    const mockExecution = {
      id: 'exec-successful-alias',
      scenario_id: 'successful-alias-scenario',
      status: 'completed',
      started_at: '2024-01-15T10:00:00Z',
      safe_mode: true,
      score: { overall: 0, blocked: 0, detected: 0, successful: 1, total: 1 },
    };

    const mockResults = [
      {
        id: 'result-successful',
        execution_id: 'exec-successful-alias',
        technique_id: 'T1082',
        agent_paw: 'agent-1',
        status: 'successful',
        output: 'System enumerated',
        detected: false,
        start_time: '2024-01-15T10:00:00Z',
        end_time: '2024-01-15T10:01:00Z',
      },
    ];

    vi.mocked(executionApi.get).mockResolvedValue({ data: mockExecution } as never);
    vi.mocked(executionApi.getResults).mockResolvedValue({ data: mockResults } as never);

    renderWithRouter('exec-successful-alias');

    // getStatusLabel 'successful' case returns 'Attack Succeeded'
    expect(await screen.findByText('Attack Succeeded')).toBeInTheDocument();
  });

  it('renders result with "running" status (Running label with warning badge)', async () => {
    const mockExecution = {
      id: 'exec-with-running-result',
      scenario_id: 'running-result-scenario',
      status: 'running',
      started_at: '2024-01-15T10:00:00Z',
      safe_mode: true,
    };

    const mockResults = [
      {
        id: 'result-running',
        execution_id: 'exec-with-running-result',
        technique_id: 'T1057',
        agent_paw: 'agent-1',
        status: 'running',
        output: '',
        detected: false,
        start_time: '2024-01-15T10:00:00Z',
        end_time: '',
      },
    ];

    vi.mocked(executionApi.get).mockResolvedValue({ data: mockExecution } as never);
    vi.mocked(executionApi.getResults).mockResolvedValue({ data: mockResults } as never);

    renderWithRouter('exec-with-running-result');

    // getStatusLabel 'running' case returns 'Running' as result label
    expect(await screen.findByText('Running')).toBeInTheDocument();
  });

  it('renders blocked result with "Blocked" label in status badge', async () => {
    const mockExecution = {
      id: 'exec-blocked-label',
      scenario_id: 'blocked-label-scenario',
      status: 'completed',
      started_at: '2024-01-15T10:00:00Z',
      safe_mode: true,
      score: { overall: 100, blocked: 1, detected: 0, successful: 0, total: 1 },
    };

    const mockResults = [
      {
        id: 'result-blocked-label',
        execution_id: 'exec-blocked-label',
        technique_id: 'T1070',
        agent_paw: 'agent-2',
        status: 'blocked',
        output: 'Blocked by endpoint protection',
        detected: false,
        start_time: '2024-01-15T10:00:00Z',
        end_time: '2024-01-15T10:01:00Z',
      },
    ];

    vi.mocked(executionApi.get).mockResolvedValue({ data: mockExecution } as never);
    vi.mocked(executionApi.getResults).mockResolvedValue({ data: mockResults } as never);

    renderWithRouter('exec-blocked-label');

    // Wait for the technique to render, then find the 'Blocked' status label in the badge
    expect(await screen.findByText('T1070')).toBeInTheDocument();
    // There are two "Blocked" texts: one in score breakdown, one as result status label
    const blockedElements = screen.getAllByText('Blocked');
    expect(blockedElements.length).toBeGreaterThanOrEqual(2);
    // The result status badge should have the badge-success class
    const statusBadge = blockedElements.find(el => el.className.includes('badge'));
    expect(statusBadge).toBeDefined();
    expect(statusBadge!.className).toContain('badge-success');
  });

  it('renders execution without score (no score breakdown section)', async () => {
    const mockExecution = {
      id: 'exec-no-score',
      scenario_id: 'no-score-scenario',
      status: 'completed',
      started_at: '2024-01-15T10:00:00Z',
      safe_mode: true,
      // no score property at all
    };

    vi.mocked(executionApi.get).mockResolvedValue({ data: mockExecution } as never);
    vi.mocked(executionApi.getResults).mockResolvedValue({ data: [] } as never);

    renderWithRouter('exec-no-score');

    expect(await screen.findByText('Execution Details')).toBeInTheDocument();
    // Without score, the breakdown section should not render
    expect(screen.queryByText('Security Score Breakdown')).not.toBeInTheDocument();
    // The overall score should show '0%' via the || '0' fallback
    expect(screen.getByText('0%')).toBeInTheDocument();
  });

  it('renders execution with score.overall = 0 using fallback display', async () => {
    const mockExecution = {
      id: 'exec-zero-score',
      scenario_id: 'zero-score-scenario',
      status: 'completed',
      started_at: '2024-01-15T10:00:00Z',
      safe_mode: true,
      score: { overall: 0, blocked: 0, detected: 0, successful: 1, total: 1 },
    };

    vi.mocked(executionApi.get).mockResolvedValue({ data: mockExecution } as never);
    vi.mocked(executionApi.getResults).mockResolvedValue({ data: [] } as never);

    renderWithRouter('exec-zero-score');

    expect(await screen.findByText('Execution Details')).toBeInTheDocument();
    // score.overall is 0, so toFixed(1) returns "0.0" which is falsy-ish but
    // actually "0.0" is truthy, so it should render "0.0%"
    // However if overall is exactly 0, 0.toFixed(1) = "0.0" which is truthy
    expect(screen.getByText('0.0%')).toBeInTheDocument();
    // Score breakdown section should render since score exists
    expect(screen.getByText('Security Score Breakdown')).toBeInTheDocument();
  });

  it('renders result with null output showing "No output" message', async () => {
    const mockExecution = {
      id: 'exec-null-output',
      scenario_id: 'null-output-scenario',
      status: 'completed',
      started_at: '2024-01-15T10:00:00Z',
      safe_mode: true,
      score: { overall: 50, blocked: 0, detected: 1, successful: 0, total: 1 },
    };

    const mockResults = [
      {
        id: 'result-null-output',
        execution_id: 'exec-null-output',
        technique_id: 'T1087',
        agent_paw: 'agent-1',
        status: 'detected',
        output: null,
        detected: true,
        start_time: '2024-01-15T10:00:00Z',
        end_time: '2024-01-15T10:01:00Z',
      },
    ];

    vi.mocked(executionApi.get).mockResolvedValue({ data: mockExecution } as never);
    vi.mocked(executionApi.getResults).mockResolvedValue({ data: mockResults } as never);

    renderWithRouter('exec-null-output');

    // When output is null (falsy), the "No output" fallback should render
    expect(await screen.findByText('No output')).toBeInTheDocument();
    expect(screen.queryByText('View output')).not.toBeInTheDocument();
  });

  it('renders "Attack Succeeded" label for "success" result status', async () => {
    const mockExecution = {
      id: 'exec-attack-succeeded',
      scenario_id: 'attack-scenario',
      status: 'completed',
      started_at: '2024-01-15T10:00:00Z',
      safe_mode: true,
      score: { overall: 0, blocked: 0, detected: 0, successful: 1, total: 1 },
    };

    const mockResults = [
      {
        id: 'result-attack-succeeded',
        execution_id: 'exec-attack-succeeded',
        technique_id: 'T1082',
        agent_paw: 'agent-1',
        status: 'success',
        output: 'Attack output data',
        detected: false,
        start_time: '2024-01-15T10:00:00Z',
        end_time: '2024-01-15T10:01:00Z',
      },
    ];

    vi.mocked(executionApi.get).mockResolvedValue({ data: mockExecution } as never);
    vi.mocked(executionApi.getResults).mockResolvedValue({ data: mockResults } as never);

    renderWithRouter('exec-attack-succeeded');

    // getStatusLabel 'success' case returns 'Attack Succeeded'
    expect(await screen.findByText('Attack Succeeded')).toBeInTheDocument();
  });

  it('renders "Detected" label for detected result status in badge', async () => {
    const mockExecution = {
      id: 'exec-detected-label',
      scenario_id: 'detected-label-scenario',
      status: 'completed',
      started_at: '2024-01-15T10:00:00Z',
      safe_mode: true,
      score: { overall: 50, blocked: 0, detected: 1, successful: 0, total: 1 },
    };

    const mockResults = [
      {
        id: 'result-detected-label',
        execution_id: 'exec-detected-label',
        technique_id: 'T1016',
        agent_paw: 'agent-1',
        status: 'detected',
        output: 'Network scan detected',
        detected: true,
        start_time: '2024-01-15T10:00:00Z',
        end_time: '2024-01-15T10:01:00Z',
      },
    ];

    vi.mocked(executionApi.get).mockResolvedValue({ data: mockExecution } as never);
    vi.mocked(executionApi.getResults).mockResolvedValue({ data: mockResults } as never);

    renderWithRouter('exec-detected-label');

    // Wait for technique to render
    expect(await screen.findByText('T1016')).toBeInTheDocument();
    // There are two "Detected" texts: one in score breakdown, one as result status label
    const detectedElements = screen.getAllByText('Detected');
    expect(detectedElements.length).toBeGreaterThanOrEqual(2);
    // The result status badge should have the badge-success class
    const statusBadge = detectedElements.find(el => el.className.includes('badge'));
    expect(statusBadge).toBeDefined();
    expect(statusBadge!.className).toContain('badge-success');
  });

  it('renders all result statuses in a single execution with multiple techniques', async () => {
    const mockExecution = {
      id: 'exec-all-statuses',
      scenario_id: 'all-statuses-scenario',
      status: 'completed',
      started_at: '2024-01-15T10:00:00Z',
      safe_mode: false,
      score: { overall: 50, blocked: 1, detected: 1, successful: 2, total: 5 },
    };

    const mockResults = [
      {
        id: 'r-success',
        execution_id: 'exec-all-statuses',
        technique_id: 'T1082',
        agent_paw: 'agent-1',
        status: 'success',
        output: 'success output',
        detected: false,
        start_time: '2024-01-15T10:00:00Z',
        end_time: '2024-01-15T10:01:00Z',
      },
      {
        id: 'r-successful',
        execution_id: 'exec-all-statuses',
        technique_id: 'T1083',
        agent_paw: 'agent-1',
        status: 'successful',
        output: 'successful output',
        detected: false,
        start_time: '2024-01-15T10:01:00Z',
        end_time: '2024-01-15T10:02:00Z',
      },
      {
        id: 'r-blocked',
        execution_id: 'exec-all-statuses',
        technique_id: 'T1059',
        agent_paw: 'agent-1',
        status: 'blocked',
        output: '',
        detected: false,
        start_time: '2024-01-15T10:02:00Z',
        end_time: '2024-01-15T10:03:00Z',
      },
      {
        id: 'r-detected',
        execution_id: 'exec-all-statuses',
        technique_id: 'T1016',
        agent_paw: 'agent-1',
        status: 'detected',
        output: null,
        detected: true,
        start_time: '2024-01-15T10:03:00Z',
        end_time: '2024-01-15T10:04:00Z',
      },
      {
        id: 'r-running',
        execution_id: 'exec-all-statuses',
        technique_id: 'T1057',
        agent_paw: 'agent-1',
        status: 'running',
        output: '',
        detected: false,
        start_time: '2024-01-15T10:04:00Z',
        end_time: '',
      },
    ];

    vi.mocked(executionApi.get).mockResolvedValue({ data: mockExecution } as never);
    vi.mocked(executionApi.getResults).mockResolvedValue({ data: mockResults } as never);

    renderWithRouter('exec-all-statuses');

    // Verify all technique IDs rendered
    expect(await screen.findByText('T1082')).toBeInTheDocument();
    expect(screen.getByText('T1083')).toBeInTheDocument();
    expect(screen.getByText('T1059')).toBeInTheDocument();
    expect(screen.getByText('T1016')).toBeInTheDocument();
    expect(screen.getByText('T1057')).toBeInTheDocument();

    // Verify status labels for all branches
    // Both 'success' and 'successful' map to 'Attack Succeeded'
    const attackSucceeded = screen.getAllByText('Attack Succeeded');
    expect(attackSucceeded).toHaveLength(2);

    // "Blocked" appears in both score breakdown section and result status badge
    const blockedElements = screen.getAllByText('Blocked');
    expect(blockedElements.length).toBeGreaterThanOrEqual(2);

    // "Detected" appears in both score breakdown section and result status badge
    const detectedElements = screen.getAllByText('Detected');
    expect(detectedElements.length).toBeGreaterThanOrEqual(2);

    expect(screen.getByText('Running')).toBeInTheDocument();

    // Full Mode badge since safe_mode is false
    expect(screen.getByText('Full Mode')).toBeInTheDocument();
  });

  it('renders execution with pending execution status using warning badge', async () => {
    const mockExecution = {
      id: 'exec-pending-badge',
      scenario_id: 'pending-badge-scenario',
      status: 'pending',
      started_at: '2024-01-15T14:00:00Z',
      safe_mode: true,
    };

    vi.mocked(executionApi.get).mockResolvedValue({ data: mockExecution } as never);
    vi.mocked(executionApi.getResults).mockResolvedValue({ data: [] } as never);

    renderWithRouter('exec-pending-badge');

    // getExecutionStatusBadge 'pending' case returns 'badge-warning'
    const badge = await screen.findByText('pending');
    expect(badge).toBeInTheDocument();
    expect(badge.className).toContain('badge-warning');
  });

  it('renders execution with running status using warning badge', async () => {
    const mockExecution = {
      id: 'exec-running-badge',
      scenario_id: 'running-badge-scenario',
      status: 'running',
      started_at: '2024-01-15T12:00:00Z',
      safe_mode: true,
    };

    vi.mocked(executionApi.get).mockResolvedValue({ data: mockExecution } as never);
    vi.mocked(executionApi.getResults).mockResolvedValue({ data: [] } as never);

    renderWithRouter('exec-running-badge');

    // getExecutionStatusBadge 'running' case returns 'badge-warning'
    const badge = await screen.findByText('running');
    expect(badge).toBeInTheDocument();
    expect(badge.className).toContain('badge-warning');
  });

  it('renders execution with completed status using success badge', async () => {
    const mockExecution = {
      id: 'exec-completed-badge',
      scenario_id: 'completed-badge-scenario',
      status: 'completed',
      started_at: '2024-01-15T10:00:00Z',
      safe_mode: true,
      score: { overall: 85.0, blocked: 8, detected: 1, successful: 1, total: 10 },
    };

    vi.mocked(executionApi.get).mockResolvedValue({ data: mockExecution } as never);
    vi.mocked(executionApi.getResults).mockResolvedValue({ data: [] } as never);

    renderWithRouter('exec-completed-badge');

    // getExecutionStatusBadge 'completed' case returns 'badge-success'
    const badge = await screen.findByText('completed');
    expect(badge).toBeInTheDocument();
    expect(badge.className).toContain('badge-success');
  });

  it('renders execution with cancelled status using danger badge', async () => {
    const mockExecution = {
      id: 'exec-cancelled-badge',
      scenario_id: 'cancelled-badge-scenario',
      status: 'cancelled',
      started_at: '2024-01-15T10:00:00Z',
      safe_mode: true,
    };

    vi.mocked(executionApi.get).mockResolvedValue({ data: mockExecution } as never);
    vi.mocked(executionApi.getResults).mockResolvedValue({ data: [] } as never);

    renderWithRouter('exec-cancelled-badge');

    // getExecutionStatusBadge 'cancelled' case returns 'badge-danger'
    const badge = await screen.findByText('cancelled');
    expect(badge).toBeInTheDocument();
    expect(badge.className).toContain('badge-danger');
  });
});
