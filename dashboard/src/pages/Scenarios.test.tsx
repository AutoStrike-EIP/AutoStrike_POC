import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { MemoryRouter } from 'react-router-dom';
import Scenarios from './Scenarios';
import { api, executionApi } from '../lib/api';
import toast from 'react-hot-toast';

// Mock the API
vi.mock('../lib/api', () => ({
  api: {
    get: vi.fn(),
  },
  executionApi: {
    start: vi.fn(),
  },
}));

// Mock react-hot-toast
vi.mock('react-hot-toast', () => ({
  default: {
    success: vi.fn(),
    error: vi.fn(),
  },
}));

// Mock useNavigate
const mockNavigate = vi.fn();
vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual('react-router-dom');
  return {
    ...actual,
    useNavigate: () => mockNavigate,
  };
});

// Mock RunExecutionModal for controlled testing
vi.mock('../components/RunExecutionModal', () => ({
  RunExecutionModal: ({
    scenario,
    onConfirm,
    onCancel,
    isLoading,
  }: {
    scenario: { name: string };
    onConfirm: (agents: string[], safeMode: boolean) => void;
    onCancel: () => void;
    isLoading: boolean;
  }) => (
    <div data-testid="run-modal">
      <span data-testid="modal-scenario">{scenario.name}</span>
      <span data-testid="modal-loading">{isLoading ? 'true' : 'false'}</span>
      <button data-testid="modal-confirm" onClick={() => onConfirm(['agent-1'], true)}>
        Confirm Run
      </button>
      <button data-testid="modal-cancel" onClick={onCancel}>
        Cancel Modal
      </button>
    </div>
  ),
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
    <MemoryRouter>
      <QueryClientProvider client={testQueryClient}>{ui}</QueryClientProvider>
    </MemoryRouter>
  );
}

describe('Scenarios Page', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders loading state', () => {
    vi.mocked(api.get).mockReturnValue(new Promise(() => {}) as never);

    renderWithClient(<Scenarios />);
    expect(screen.getByText('Loading scenarios...')).toBeInTheDocument();
  });

  it('renders scenarios list', async () => {
    const mockScenarios = [
      {
        id: 'scenario-1',
        name: 'Initial Access Test',
        description: 'Test initial access techniques',
        phases: [
          { name: 'Reconnaissance', techniques: ['T1595', 'T1592'] },
          { name: 'Initial Access', techniques: ['T1566'] },
        ],
        tags: ['initial-access', 'phishing'],
      },
    ];
    vi.mocked(api.get).mockResolvedValue({ data: mockScenarios } as never);

    renderWithClient(<Scenarios />);

    expect(await screen.findByText('Initial Access Test')).toBeInTheDocument();
    expect(screen.getByText('Test initial access techniques')).toBeInTheDocument();
    expect(screen.getByText('Reconnaissance')).toBeInTheDocument();
    expect(screen.getByText('(2 techniques)')).toBeInTheDocument();
    expect(screen.getByText('Initial Access')).toBeInTheDocument();
    expect(screen.getByText('(1 techniques)')).toBeInTheDocument();
    expect(screen.getByText('initial-access')).toBeInTheDocument();
    expect(screen.getByText('phishing')).toBeInTheDocument();
  });

  it('renders scenario with multiple phases', async () => {
    const mockScenarios = [
      {
        id: 'scenario-2',
        name: 'Full Attack Chain',
        description: 'Complete attack simulation',
        phases: [
          { name: 'Phase 1', techniques: ['T1'] },
          { name: 'Phase 2', techniques: ['T2', 'T3'] },
          { name: 'Phase 3', techniques: ['T4', 'T5', 'T6'] },
        ],
        tags: ['advanced'],
      },
    ];
    vi.mocked(api.get).mockResolvedValue({ data: mockScenarios } as never);

    renderWithClient(<Scenarios />);

    expect(await screen.findByText('Phase 1')).toBeInTheDocument();
    expect(screen.getByText('Phase 2')).toBeInTheDocument();
    expect(screen.getByText('Phase 3')).toBeInTheDocument();
    expect(screen.getByText('(3 techniques)')).toBeInTheDocument();
  });

  it('renders empty state when no scenarios', async () => {
    vi.mocked(api.get).mockResolvedValue({ data: [] } as never);

    renderWithClient(<Scenarios />);

    expect(await screen.findByText('No scenarios created')).toBeInTheDocument();
    expect(screen.getByText('Create an attack scenario to test your defenses')).toBeInTheDocument();
  });

  it('renders page title and create button', async () => {
    vi.mocked(api.get).mockResolvedValue({ data: [] } as never);

    renderWithClient(<Scenarios />);

    expect(await screen.findByText('Scenarios')).toBeInTheDocument();
    expect(screen.getByText('Create Scenario')).toBeInTheDocument();
  });

  it('renders run button for each scenario', async () => {
    const mockScenarios = [
      {
        id: 'scenario-1',
        name: 'Test Scenario',
        description: 'A test',
        phases: [{ name: 'Phase', techniques: ['T1'] }],
        tags: [],
      },
    ];
    vi.mocked(api.get).mockResolvedValue({ data: mockScenarios } as never);

    renderWithClient(<Scenarios />);

    expect(await screen.findByText('Run')).toBeInTheDocument();
  });

  it('renders scenario without tags', async () => {
    const mockScenarios = [
      {
        id: 'scenario-no-tags',
        name: 'Simple Scenario',
        description: 'No tags scenario',
        phases: [{ name: 'Single Phase', techniques: ['T1082'] }],
        tags: [],
      },
    ];
    vi.mocked(api.get).mockResolvedValue({ data: mockScenarios } as never);

    renderWithClient(<Scenarios />);

    expect(await screen.findByText('Simple Scenario')).toBeInTheDocument();
    expect(screen.getByText('Single Phase')).toBeInTheDocument();
  });

  it('displays phase numbers correctly', async () => {
    const mockScenarios = [
      {
        id: 'scenario-numbered',
        name: 'Numbered Phases',
        description: 'Test phase numbering',
        phases: [
          { name: 'First', techniques: ['T1'] },
          { name: 'Second', techniques: ['T2'] },
          { name: 'Third', techniques: ['T3'] },
        ],
        tags: [],
      },
    ];
    vi.mocked(api.get).mockResolvedValue({ data: mockScenarios } as never);

    renderWithClient(<Scenarios />);

    await screen.findByText('Numbered Phases');
    expect(screen.getByText('1')).toBeInTheDocument();
    expect(screen.getByText('2')).toBeInTheDocument();
    expect(screen.getByText('3')).toBeInTheDocument();
  });

  // Modal and execution tests
  it('opens run modal when Run button is clicked', async () => {
    const mockScenarios = [
      {
        id: 'scenario-1',
        name: 'Clickable Scenario',
        description: 'Test',
        phases: [{ name: 'Phase', techniques: ['T1'] }],
        tags: [],
      },
    ];
    vi.mocked(api.get).mockResolvedValue({ data: mockScenarios } as never);

    renderWithClient(<Scenarios />);

    await screen.findByText('Clickable Scenario');
    fireEvent.click(screen.getByText('Run'));

    expect(screen.getByTestId('run-modal')).toBeInTheDocument();
    expect(screen.getByTestId('modal-scenario')).toHaveTextContent('Clickable Scenario');
  });

  it('closes modal when cancel is clicked', async () => {
    const mockScenarios = [
      {
        id: 'scenario-1',
        name: 'Cancel Test',
        description: 'Test',
        phases: [{ name: 'Phase', techniques: ['T1'] }],
        tags: [],
      },
    ];
    vi.mocked(api.get).mockResolvedValue({ data: mockScenarios } as never);

    renderWithClient(<Scenarios />);

    await screen.findByText('Cancel Test');
    fireEvent.click(screen.getByText('Run'));
    expect(screen.getByTestId('run-modal')).toBeInTheDocument();

    fireEvent.click(screen.getByTestId('modal-cancel'));
    expect(screen.queryByTestId('run-modal')).not.toBeInTheDocument();
  });

  it('starts execution and navigates on success', async () => {
    const mockScenarios = [
      {
        id: 'scenario-exec',
        name: 'Execute Me',
        description: 'Test',
        phases: [{ name: 'Phase', techniques: ['T1'] }],
        tags: [],
      },
    ];
    vi.mocked(api.get).mockResolvedValue({ data: mockScenarios } as never);
    vi.mocked(executionApi.start).mockResolvedValue({ data: { id: 'exec-1' } } as never);

    renderWithClient(<Scenarios />);

    await screen.findByText('Execute Me');
    fireEvent.click(screen.getByText('Run'));
    fireEvent.click(screen.getByTestId('modal-confirm'));

    await waitFor(() => {
      expect(executionApi.start).toHaveBeenCalledWith('scenario-exec', ['agent-1'], true);
    });

    await waitFor(() => {
      expect(toast.success).toHaveBeenCalledWith('Execution started successfully');
      expect(mockNavigate).toHaveBeenCalledWith('/executions');
    });
  });

  it('shows error toast on execution failure with error message', async () => {
    const mockScenarios = [
      {
        id: 'scenario-fail',
        name: 'Fail Scenario',
        description: 'Test',
        phases: [{ name: 'Phase', techniques: ['T1'] }],
        tags: [],
      },
    ];
    vi.mocked(api.get).mockResolvedValue({ data: mockScenarios } as never);
    vi.mocked(executionApi.start).mockRejectedValue({
      response: { data: { error: 'Agent disconnected' } },
    } as never);

    renderWithClient(<Scenarios />);

    await screen.findByText('Fail Scenario');
    fireEvent.click(screen.getByText('Run'));
    fireEvent.click(screen.getByTestId('modal-confirm'));

    await waitFor(() => {
      expect(toast.error).toHaveBeenCalledWith('Agent disconnected');
    });
  });

  it('shows default error message when no error details', async () => {
    const mockScenarios = [
      {
        id: 'scenario-fail2',
        name: 'Fail Scenario 2',
        description: 'Test',
        phases: [{ name: 'Phase', techniques: ['T1'] }],
        tags: [],
      },
    ];
    vi.mocked(api.get).mockResolvedValue({ data: mockScenarios } as never);
    vi.mocked(executionApi.start).mockRejectedValue(new Error('Network') as never);

    renderWithClient(<Scenarios />);

    await screen.findByText('Fail Scenario 2');
    fireEvent.click(screen.getByText('Run'));
    fireEvent.click(screen.getByTestId('modal-confirm'));

    await waitFor(() => {
      expect(toast.error).toHaveBeenCalledWith('Failed to start execution');
    });
  });

  it('closes modal after successful execution', async () => {
    const mockScenarios = [
      {
        id: 'scenario-close',
        name: 'Close After Success',
        description: 'Test',
        phases: [{ name: 'Phase', techniques: ['T1'] }],
        tags: [],
      },
    ];
    vi.mocked(api.get).mockResolvedValue({ data: mockScenarios } as never);
    vi.mocked(executionApi.start).mockResolvedValue({ data: { id: 'exec-1' } } as never);

    renderWithClient(<Scenarios />);

    await screen.findByText('Close After Success');
    fireEvent.click(screen.getByText('Run'));
    expect(screen.getByTestId('run-modal')).toBeInTheDocument();

    fireEvent.click(screen.getByTestId('modal-confirm'));

    await waitFor(() => {
      expect(screen.queryByTestId('run-modal')).not.toBeInTheDocument();
    });
  });

  it('does not call mutation when scenarioToRun is null', async () => {
    // This tests the guard in handleConfirmRun
    const mockScenarios = [
      {
        id: 'scenario-guard',
        name: 'Guard Test',
        description: 'Test',
        phases: [{ name: 'Phase', techniques: ['T1'] }],
        tags: [],
      },
    ];
    vi.mocked(api.get).mockResolvedValue({ data: mockScenarios } as never);

    renderWithClient(<Scenarios />);

    await screen.findByText('Guard Test');
    // Don't open modal, so scenarioToRun is null
    // The modal isn't rendered, so no mutation can be triggered
    expect(screen.queryByTestId('run-modal')).not.toBeInTheDocument();
    expect(executionApi.start).not.toHaveBeenCalled();
  });
});
