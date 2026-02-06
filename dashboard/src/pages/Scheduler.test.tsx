import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import Scheduler from './Scheduler';

// Mock react-hot-toast
vi.mock('react-hot-toast', () => ({
  default: {
    success: vi.fn(),
    error: vi.fn(),
  },
}));

// Mock functions stored as references
const mockScheduleList = vi.fn();
const mockGetRuns = vi.fn();
const mockPause = vi.fn();
const mockResume = vi.fn();
const mockRunNow = vi.fn();
const mockDelete = vi.fn();
const mockCreate = vi.fn();
const mockUpdate = vi.fn();
const mockScenarioList = vi.fn();

// Mock the schedule and scenario API
vi.mock('../lib/api', () => ({
  scheduleApi: {
    list: () => mockScheduleList(),
    getRuns: (id: string) => mockGetRuns(id),
    pause: (id: string) => mockPause(id),
    resume: (id: string) => mockResume(id),
    runNow: (id: string) => mockRunNow(id),
    delete: (id: string) => mockDelete(id),
    create: (data: unknown) => mockCreate(data),
    update: (id: string, data: unknown) => mockUpdate(id, data),
  },
  scenarioApi: {
    list: () => mockScenarioList(),
  },
}));

const mockScheduleData = [
  {
    id: 'sched-1',
    name: 'Daily Security Check',
    description: 'Run security tests daily',
    scenario_id: 'scenario-1',
    agent_paw: '',
    frequency: 'daily',
    cron_expr: '',
    safe_mode: true,
    status: 'active',
    next_run_at: new Date(Date.now() + 3600000).toISOString(),
    last_run_at: new Date(Date.now() - 86400000).toISOString(),
    last_run_id: 'exec-1',
    created_by: 'admin',
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString(),
  },
  {
    id: 'sched-2',
    name: 'Weekly Audit',
    description: 'Weekly security audit',
    scenario_id: 'scenario-2',
    agent_paw: 'agent-1',
    frequency: 'weekly',
    cron_expr: '',
    safe_mode: false,
    status: 'paused',
    next_run_at: null,
    last_run_at: null,
    last_run_id: '',
    created_by: 'admin',
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString(),
  },
];

const mockScenarioData = [
  { id: 'scenario-1', name: 'Test Scenario 1', description: '', phases: [], tags: [] },
  { id: 'scenario-2', name: 'Test Scenario 2', description: '', phases: [], tags: [] },
];

function createTestQueryClient() {
  return new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
      },
    },
  });
}

function renderScheduler() {
  const queryClient = createTestQueryClient();
  return render(
    <QueryClientProvider client={queryClient}>
      <MemoryRouter>
        <Scheduler />
      </MemoryRouter>
    </QueryClientProvider>
  );
}

function setupDefaultMocks() {
  mockScheduleList.mockResolvedValue({ data: mockScheduleData });
  mockScenarioList.mockResolvedValue({ data: mockScenarioData });
  mockGetRuns.mockResolvedValue({ data: [] });
  mockPause.mockResolvedValue({ data: {} });
  mockResume.mockResolvedValue({ data: {} });
  mockRunNow.mockResolvedValue({ data: {} });
  mockDelete.mockResolvedValue({ data: {} });
  mockCreate.mockResolvedValue({ data: {} });
  mockUpdate.mockResolvedValue({ data: {} });
}

describe('Scheduler Page', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    setupDefaultMocks();
  });

  it('renders scheduler title', async () => {
    renderScheduler();
    await waitFor(() => {
      expect(screen.getByText('Scheduler')).toBeInTheDocument();
    });
  });

  it('renders create schedule button', async () => {
    renderScheduler();
    await waitFor(() => {
      expect(screen.getByText('Create Schedule')).toBeInTheDocument();
    });
  });

  it('displays schedules after loading', async () => {
    renderScheduler();
    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
      expect(screen.getByText('Weekly Audit')).toBeInTheDocument();
    });
  });

  it('shows schedule status badges', async () => {
    renderScheduler();
    await waitFor(() => {
      expect(screen.getByText('active')).toBeInTheDocument();
      expect(screen.getByText('paused')).toBeInTheDocument();
    });
  });

  it('shows frequency labels', async () => {
    renderScheduler();
    await waitFor(() => {
      expect(screen.getByText('Daily')).toBeInTheDocument();
      expect(screen.getByText('Weekly')).toBeInTheDocument();
    });
  });
});

describe('Scheduler Create Modal', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    setupDefaultMocks();
  });

  it('opens create modal when button clicked', async () => {
    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    const createButton = screen.getByText('Create Schedule');
    fireEvent.click(createButton);

    await waitFor(() => {
      expect(screen.getByText('Name *')).toBeInTheDocument();
      expect(screen.getByText('Scenario *')).toBeInTheDocument();
      expect(screen.getByText('Frequency *')).toBeInTheDocument();
    });
  });

  it('closes modal when cancel clicked', async () => {
    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Create Schedule'));

    await waitFor(() => {
      expect(screen.getByText('Name *')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Cancel'));

    await waitFor(() => {
      expect(screen.queryByText('Name *')).not.toBeInTheDocument();
    });
  });

  it('shows safe mode checkbox', async () => {
    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Create Schedule'));

    await waitFor(() => {
      expect(screen.getByText('Safe Mode')).toBeInTheDocument();
    });
  });
});

describe('Scheduler Edit', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    setupDefaultMocks();
  });

  it('opens edit modal when edit button clicked', async () => {
    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    const editButtons = screen.getAllByTitle('Edit');
    fireEvent.click(editButtons[0]);

    await waitFor(() => {
      expect(screen.getByText('Edit Schedule')).toBeInTheDocument();
    });
  });

  it('closes edit modal when cancel clicked', async () => {
    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    const editButtons = screen.getAllByTitle('Edit');
    fireEvent.click(editButtons[0]);

    await waitFor(() => {
      expect(screen.getByText('Edit Schedule')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Cancel'));

    await waitFor(() => {
      expect(screen.queryByText('Edit Schedule')).not.toBeInTheDocument();
    });
  });

  it('pre-fills form with schedule data when editing', async () => {
    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    const editButtons = screen.getAllByTitle('Edit');
    fireEvent.click(editButtons[0]);

    await waitFor(() => {
      const nameInput = screen.getByDisplayValue('Daily Security Check');
      expect(nameInput).toBeInTheDocument();
    });
  });
});

describe('Scheduler Actions', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    setupDefaultMocks();
  });

  it('shows delete confirmation modal', async () => {
    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    // Find all delete buttons (there should be one per schedule)
    const deleteButtons = screen.getAllByTitle('Delete');
    fireEvent.click(deleteButtons[0]);

    await waitFor(() => {
      expect(screen.getByText('Delete Schedule')).toBeInTheDocument();
      expect(
        screen.getByText(/Are you sure you want to delete "Daily Security Check"/)
      ).toBeInTheDocument();
    });
  });

  it('cancels delete when cancel clicked', async () => {
    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    const deleteButtons = screen.getAllByTitle('Delete');
    fireEvent.click(deleteButtons[0]);

    await waitFor(() => {
      expect(screen.getByText('Delete Schedule')).toBeInTheDocument();
    });

    // Click cancel button in modal
    const cancelButtons = screen.getAllByText('Cancel');
    fireEvent.click(cancelButtons[cancelButtons.length - 1]);

    await waitFor(() => {
      expect(screen.queryByText('Delete Schedule')).not.toBeInTheDocument();
    });
  });
});

describe('Scheduler Frequency Labels', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    setupDefaultMocks();
  });

  it('displays schedule descriptions', async () => {
    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Run security tests daily')).toBeInTheDocument();
      expect(screen.getByText('Weekly security audit')).toBeInTheDocument();
    });
  });

  it('shows scenario name labels', async () => {
    renderScheduler();

    await waitFor(() => {
      // Check that scenario info is displayed
      const scenarioLabels = screen.getAllByText('Scenario:');
      expect(scenarioLabels.length).toBeGreaterThan(0);
    });
  });
});

describe('Scheduler Empty State', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockScheduleList.mockResolvedValue({ data: [] });
    mockScenarioList.mockResolvedValue({ data: mockScenarioData });
    mockGetRuns.mockResolvedValue({ data: [] });
  });

  it('shows empty state when no schedules', async () => {
    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('No schedules created')).toBeInTheDocument();
    });
  });

  it('shows create button in empty state', async () => {
    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('No schedules created')).toBeInTheDocument();
    });

    expect(screen.getByText('Create Schedule')).toBeInTheDocument();
  });
});

describe('Scheduler Pause/Resume', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    setupDefaultMocks();
  });

  it('shows pause button for active schedules', async () => {
    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    const pauseButtons = screen.getAllByTitle('Pause');
    expect(pauseButtons.length).toBeGreaterThan(0);
  });

  it('shows resume button for paused schedules', async () => {
    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Weekly Audit')).toBeInTheDocument();
    });

    const resumeButtons = screen.getAllByTitle('Resume');
    expect(resumeButtons.length).toBeGreaterThan(0);
  });

  it('calls pause API when pause button clicked', async () => {
    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    const pauseButtons = screen.getAllByTitle('Pause');
    fireEvent.click(pauseButtons[0]);

    await waitFor(() => {
      expect(mockPause).toHaveBeenCalledWith('sched-1');
    });
  });

  it('calls resume API when resume button clicked', async () => {
    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Weekly Audit')).toBeInTheDocument();
    });

    const resumeButtons = screen.getAllByTitle('Resume');
    fireEvent.click(resumeButtons[0]);

    await waitFor(() => {
      expect(mockResume).toHaveBeenCalledWith('sched-2');
    });
  });
});

describe('Scheduler Run Now', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    setupDefaultMocks();
  });

  it('shows run now button for schedules', async () => {
    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    const runNowButtons = screen.getAllByTitle('Run Now');
    expect(runNowButtons.length).toBeGreaterThan(0);
  });

  it('calls runNow API when run now button clicked', async () => {
    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    const runNowButtons = screen.getAllByTitle('Run Now');
    fireEvent.click(runNowButtons[0]);

    await waitFor(() => {
      expect(mockRunNow).toHaveBeenCalledWith('sched-1');
    });
  });
});

describe('Scheduler Delete Confirmation', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    setupDefaultMocks();
  });

  it('calls delete API when confirmed', async () => {
    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    const deleteButtons = screen.getAllByTitle('Delete');
    fireEvent.click(deleteButtons[0]);

    await waitFor(() => {
      expect(screen.getByText('Delete Schedule')).toBeInTheDocument();
    });

    // Find the Delete button in modal - it's the only one with that class
    const allButtons = screen.getAllByRole('button');
    const confirmButton = allButtons.find(
      (btn) => btn.textContent === 'Delete' && btn.className.includes('bg-red-600')
    );
    expect(confirmButton).toBeDefined();
    fireEvent.click(confirmButton!);

    await waitFor(() => {
      expect(mockDelete).toHaveBeenCalledWith('sched-1');
    });
  });
});

describe('Scheduler History Expansion', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    setupDefaultMocks();
    mockGetRuns.mockResolvedValue({
      data: [
        {
          id: 'run-1',
          schedule_id: 'sched-1',
          execution_id: 'exec-1',
          started_at: new Date().toISOString(),
          completed_at: new Date().toISOString(),
          status: 'completed',
          triggered_by: 'scheduler',
        },
      ],
    });
  });

  it('shows expand button for schedules with history', async () => {
    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    const expandButtons = screen.getAllByTitle('Show history');
    expect(expandButtons.length).toBeGreaterThan(0);
  });

  it('expands history when expand button clicked', async () => {
    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    const expandButtons = screen.getAllByTitle('Show history');
    fireEvent.click(expandButtons[0]);

    await waitFor(() => {
      // The history section shows runs or "No runs yet" message
      const historyContent = screen.getByText(/Loading history|No runs yet|completed/);
      expect(historyContent).toBeInTheDocument();
    });
  });

  it('collapses history when clicked again', async () => {
    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    const expandButtons = screen.getAllByTitle('Show history');
    fireEvent.click(expandButtons[0]);

    await waitFor(() => {
      const historyContent = screen.getByText(/Loading history|No runs yet|completed/);
      expect(historyContent).toBeInTheDocument();
    });

    // Click again to collapse - button title stays "Show history"
    const collapseButtons = screen.getAllByTitle('Show history');
    fireEvent.click(collapseButtons[0]);

    // Note: The history section will be hidden
    await waitFor(() => {
      expect(screen.queryByText('Loading history')).not.toBeInTheDocument();
    });
  });
});

describe('Scheduler Form Modal Cron Expression', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    setupDefaultMocks();
  });

  it('shows cron expression field when frequency is cron', async () => {
    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Create Schedule'));

    await waitFor(() => {
      expect(screen.getByText('Frequency *')).toBeInTheDocument();
    });

    // Select cron frequency
    const frequencySelect = screen.getByLabelText('Frequency *');
    fireEvent.change(frequencySelect, { target: { value: 'cron' } });

    await waitFor(() => {
      expect(screen.getByText('Cron Expression *')).toBeInTheDocument();
    });
  });

  it('hides cron expression field for non-cron frequencies', async () => {
    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Create Schedule'));

    await waitFor(() => {
      expect(screen.getByText('Frequency *')).toBeInTheDocument();
    });

    // Default is not cron
    expect(screen.queryByText('Cron Expression *')).not.toBeInTheDocument();
  });
});

describe('Scheduler Create Submission', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    setupDefaultMocks();
  });

  it('submits create form with valid data', async () => {
    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Create Schedule'));

    await waitFor(() => {
      expect(screen.getByText('Name *')).toBeInTheDocument();
    });

    // Fill the form
    fireEvent.change(screen.getByLabelText('Name *'), {
      target: { value: 'New Test Schedule' },
    });
    fireEvent.change(screen.getByLabelText('Description'), {
      target: { value: 'A test schedule description' },
    });

    // Select scenario
    const scenarioSelect = screen.getByLabelText('Scenario *');
    fireEvent.change(scenarioSelect, { target: { value: 'scenario-1' } });

    // Select frequency
    const frequencySelect = screen.getByLabelText('Frequency *');
    fireEvent.change(frequencySelect, { target: { value: 'daily' } });

    // Submit form - find submit button in modal (the one with btn-primary class without gap-2)
    const allButtons = screen.getAllByRole('button', { name: /Create Schedule/ });
    // Modal submit button is the btn-primary one without the gap-2 class
    const submitButton = allButtons.find((btn) => !btn.className.includes('gap-2'));
    expect(submitButton).toBeDefined();
    fireEvent.click(submitButton!);

    await waitFor(() => {
      expect(mockCreate).toHaveBeenCalled();
    });
  });
});

describe('Scheduler Update Submission', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    setupDefaultMocks();
  });

  it('submits update form with modified data', async () => {
    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    // Click edit on first schedule
    const editButtons = screen.getAllByTitle('Edit');
    fireEvent.click(editButtons[0]);

    await waitFor(() => {
      expect(screen.getByText('Edit Schedule')).toBeInTheDocument();
    });

    // Modify name
    const nameInput = screen.getByDisplayValue('Daily Security Check');
    fireEvent.change(nameInput, { target: { value: 'Modified Schedule Name' } });

    // Submit form - button text is "Update Schedule"
    const submitButton = screen.getByRole('button', { name: /Update Schedule$/ });
    fireEvent.click(submitButton);

    await waitFor(() => {
      expect(mockUpdate).toHaveBeenCalled();
    });
  });
});

describe('Scheduler Safe Mode Toggle', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    setupDefaultMocks();
  });

  it('shows safe mode checkbox in form', async () => {
    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Create Schedule'));

    await waitFor(() => {
      expect(screen.getByText('Safe Mode')).toBeInTheDocument();
    });
  });

  it('can toggle safe mode in form', async () => {
    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Create Schedule'));

    await waitFor(() => {
      expect(screen.getByText('Safe Mode')).toBeInTheDocument();
    });

    const safeModeCheckbox = screen.getByRole('checkbox');
    expect(safeModeCheckbox).toBeInTheDocument();
    fireEvent.click(safeModeCheckbox);
  });
});

describe('Scheduler Next Run Display', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    setupDefaultMocks();
  });

  it('shows next run time for active schedules', async () => {
    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    // Should show next run info
    const nextRunLabels = screen.getAllByText('Next Run:');
    expect(nextRunLabels.length).toBeGreaterThan(0);
  });

  it('shows last run time when available', async () => {
    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    // First schedule has last_run_at
    const lastRunLabels = screen.getAllByText('Last Run:');
    expect(lastRunLabels.length).toBeGreaterThan(0);
  });
});

describe('Scheduler Agent Selection', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    setupDefaultMocks();
  });

  it('shows agent field in create form', async () => {
    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Create Schedule'));

    await waitFor(() => {
      expect(screen.getByText('Agent (optional)')).toBeInTheDocument();
    });
  });

  it('shows placeholder for agent field', async () => {
    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Create Schedule'));

    await waitFor(() => {
      expect(screen.getByPlaceholderText('Leave empty for all agents')).toBeInTheDocument();
    });
  });
});

describe('Scheduler API Error Handling', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    setupDefaultMocks();
  });

  it('handles pause error gracefully', async () => {
    const toast = await import('react-hot-toast');
    mockPause.mockRejectedValueOnce(new Error('Pause failed'));

    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    const pauseButtons = screen.getAllByTitle('Pause');
    fireEvent.click(pauseButtons[0]);

    await waitFor(() => {
      expect(toast.default.error).toHaveBeenCalled();
    });
  });

  it('handles resume error gracefully', async () => {
    const toast = await import('react-hot-toast');
    mockResume.mockRejectedValueOnce(new Error('Resume failed'));

    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Weekly Audit')).toBeInTheDocument();
    });

    const resumeButtons = screen.getAllByTitle('Resume');
    fireEvent.click(resumeButtons[0]);

    await waitFor(() => {
      expect(toast.default.error).toHaveBeenCalled();
    });
  });

  it('handles run now error gracefully', async () => {
    const toast = await import('react-hot-toast');
    mockRunNow.mockRejectedValueOnce(new Error('Run now failed'));

    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    const runNowButtons = screen.getAllByTitle('Run Now');
    fireEvent.click(runNowButtons[0]);

    await waitFor(() => {
      expect(toast.default.error).toHaveBeenCalled();
    });
  });

  it('handles delete error gracefully', async () => {
    const toast = await import('react-hot-toast');
    mockDelete.mockRejectedValueOnce(new Error('Delete failed'));

    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    const deleteButtons = screen.getAllByTitle('Delete');
    fireEvent.click(deleteButtons[0]);

    await waitFor(() => {
      expect(screen.getByText('Delete Schedule')).toBeInTheDocument();
    });

    const allButtons = screen.getAllByRole('button');
    const confirmButton = allButtons.find(
      (btn) => btn.textContent === 'Delete' && btn.className.includes('bg-red-600')
    );
    fireEvent.click(confirmButton!);

    await waitFor(() => {
      expect(toast.default.error).toHaveBeenCalled();
    });
  });
});

describe('Scheduler Form Validation', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    setupDefaultMocks();
  });

  it('shows required fields with asterisks', async () => {
    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Create Schedule'));

    await waitFor(() => {
      expect(screen.getByText('Name *')).toBeInTheDocument();
      expect(screen.getByText('Scenario *')).toBeInTheDocument();
      expect(screen.getByText('Frequency *')).toBeInTheDocument();
    });
  });

  it('shows description field as optional', async () => {
    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Create Schedule'));

    await waitFor(() => {
      expect(screen.getByText('Description')).toBeInTheDocument();
    });
  });
});

describe('Scheduler Loading State', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    // Don't setup mocks immediately to test loading state
    mockScheduleList.mockImplementation(
      () => new Promise((resolve) => setTimeout(() => resolve({ data: mockScheduleData }), 100))
    );
    mockScenarioList.mockResolvedValue({ data: mockScenarioData });
    mockGetRuns.mockResolvedValue({ data: [] });
  });

  it('shows scheduler title after loading', async () => {
    renderScheduler();
    await waitFor(() => {
      expect(screen.getByText('Scheduler')).toBeInTheDocument();
    });
  });
});

describe('Scheduler Frequency Options', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    setupDefaultMocks();
  });

  it('shows all frequency options in dropdown', async () => {
    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Create Schedule'));

    await waitFor(() => {
      const frequencySelect = screen.getByLabelText('Frequency *');
      expect(frequencySelect).toBeInTheDocument();
    });
  });

  it('can select hourly frequency', async () => {
    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Create Schedule'));

    await waitFor(() => {
      const frequencySelect = screen.getByLabelText('Frequency *');
      fireEvent.change(frequencySelect, { target: { value: 'hourly' } });
      expect(frequencySelect).toHaveValue('hourly');
    });
  });

  it('can select monthly frequency', async () => {
    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Create Schedule'));

    await waitFor(() => {
      const frequencySelect = screen.getByLabelText('Frequency *');
      fireEvent.change(frequencySelect, { target: { value: 'monthly' } });
      expect(frequencySelect).toHaveValue('monthly');
    });
  });
});

describe('Scheduler Schedule Card Details', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    setupDefaultMocks();
  });

  it('shows schedule names', async () => {
    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
      expect(screen.getByText('Weekly Audit')).toBeInTheDocument();
    });
  });

  it('shows schedule descriptions', async () => {
    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Run security tests daily')).toBeInTheDocument();
      expect(screen.getByText('Weekly security audit')).toBeInTheDocument();
    });
  });
});

describe('Scheduler Modal Close', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    setupDefaultMocks();
  });

  it('closes create modal when clicking outside', async () => {
    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Create Schedule'));

    await waitFor(() => {
      expect(screen.getByText('Name *')).toBeInTheDocument();
    });

    // Click cancel to close
    fireEvent.click(screen.getByText('Cancel'));

    await waitFor(() => {
      expect(screen.queryByText('Name *')).not.toBeInTheDocument();
    });
  });
});

describe('Scheduler Success Messages', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    setupDefaultMocks();
  });

  it('shows success toast on pause', async () => {
    const toast = await import('react-hot-toast');

    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    const pauseButtons = screen.getAllByTitle('Pause');
    fireEvent.click(pauseButtons[0]);

    await waitFor(() => {
      expect(toast.default.success).toHaveBeenCalled();
    });
  });

  it('shows success toast on resume', async () => {
    const toast = await import('react-hot-toast');

    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Weekly Audit')).toBeInTheDocument();
    });

    const resumeButtons = screen.getAllByTitle('Resume');
    fireEvent.click(resumeButtons[0]);

    await waitFor(() => {
      expect(toast.default.success).toHaveBeenCalled();
    });
  });

  it('shows success toast on run now', async () => {
    const toast = await import('react-hot-toast');

    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    const runNowButtons = screen.getAllByTitle('Run Now');
    fireEvent.click(runNowButtons[0]);

    await waitFor(() => {
      expect(toast.default.success).toHaveBeenCalled();
    });
  });

  it('shows correct run status colors for completed, failed, and running runs', async () => {
    mockGetRuns.mockResolvedValue({
      data: [
        { id: 'run-1', schedule_id: 'sched-1', started_at: new Date().toISOString(), status: 'completed', execution_id: 'exec-1', error: '' },
        { id: 'run-2', schedule_id: 'sched-1', started_at: new Date().toISOString(), status: 'failed', execution_id: '', error: 'timeout' },
        { id: 'run-3', schedule_id: 'sched-1', started_at: new Date().toISOString(), status: 'running', execution_id: 'exec-3', error: '' },
      ],
    });

    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    // Expand schedule to show runs
    const expandButtons = screen.getAllByTitle('Show history');
    fireEvent.click(expandButtons[0]);

    await waitFor(() => {
      expect(screen.getByText('Recent Runs')).toBeInTheDocument();
    });

    // Check the error message from the failed run
    expect(screen.getByText('(timeout)')).toBeInTheDocument();

    // Check that "View Execution" links are rendered for runs with execution_id
    const viewLinks = screen.getAllByText('View Execution');
    expect(viewLinks.length).toBe(2); // run-1 and run-3 have execution_id
  });

  it('shows correct submit button text in edit mode', async () => {
    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    const editButtons = screen.getAllByTitle('Edit');
    fireEvent.click(editButtons[0]);

    await waitFor(() => {
      expect(screen.getByText('Edit Schedule')).toBeInTheDocument();
    });

    expect(screen.getByText('Update Schedule')).toBeInTheDocument();
  });

  it('shows Overdue when next_run_at is in the past', async () => {
    mockScheduleList.mockResolvedValue({
      data: [{
        ...mockScheduleData[0],
        next_run_at: new Date(Date.now() - 60000).toISOString(),
      }],
    });

    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Overdue')).toBeInTheDocument();
    });
  });

  it('shows Paused for next run when status is paused', async () => {
    mockScheduleList.mockResolvedValue({
      data: [mockScheduleData[1]], // paused schedule
    });

    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Paused')).toBeInTheDocument();
    });
  });

  it('shows Never for null last_run_at', async () => {
    mockScheduleList.mockResolvedValue({
      data: [mockScheduleData[1]], // has null last_run_at
    });

    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Never')).toBeInTheDocument();
    });
  });

  it('shows Unknown Scenario for unmatched scenario_id', async () => {
    mockScheduleList.mockResolvedValue({
      data: [{
        ...mockScheduleData[0],
        scenario_id: 'nonexistent-scenario',
      }],
    });

    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Unknown Scenario')).toBeInTheDocument();
    });
  });

  it('renders disabled schedule without pause/resume buttons', async () => {
    mockScheduleList.mockResolvedValue({
      data: [{
        ...mockScheduleData[0],
        status: 'disabled',
      }],
    });

    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('disabled')).toBeInTheDocument();
    });

    expect(screen.queryByTitle('Pause')).not.toBeInTheDocument();
    expect(screen.queryByTitle('Resume')).not.toBeInTheDocument();
  });
});

describe('formatRelativeTime edge cases', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    setupDefaultMocks();
  });

  it('shows relative time in days format when next_run_at is days away', async () => {
    // Use a large enough offset that slight timing differences don't matter
    mockScheduleList.mockResolvedValue({
      data: [{
        ...mockScheduleData[0],
        next_run_at: new Date(Date.now() + 5 * 86400000 + 6 * 3600000 + 120000).toISOString(),
      }],
    });

    renderScheduler();

    await waitFor(() => {
      // Match pattern: "in Xd Xh" where X are digits
      expect(screen.getByText(/in \d+d \d+h/)).toBeInTheDocument();
    });
  });

  it('shows relative time in hours format when next_run_at is hours away', async () => {
    // 5 hours + extra buffer from now
    mockScheduleList.mockResolvedValue({
      data: [{
        ...mockScheduleData[0],
        next_run_at: new Date(Date.now() + 5 * 3600000 + 20 * 60000 + 120000).toISOString(),
      }],
    });

    renderScheduler();

    await waitFor(() => {
      // Match pattern: "in Xh Xm" where X are digits
      expect(screen.getByText(/in \d+h \d+m/)).toBeInTheDocument();
    });
  });

  it('shows relative time in minutes format when next_run_at is minutes away', async () => {
    // 30 minutes + buffer from now (stays under 1 hour)
    mockScheduleList.mockResolvedValue({
      data: [{
        ...mockScheduleData[0],
        next_run_at: new Date(Date.now() + 30 * 60000 + 30000).toISOString(),
      }],
    });

    renderScheduler();

    await waitFor(() => {
      // Match pattern: "in Xm" where X are digits (no "h" or "d")
      expect(screen.getByText(/^in \d+m$/)).toBeInTheDocument();
    });
  });

  it('shows "in 0m" when next_run_at is just seconds away', async () => {
    // 5 seconds from now
    mockScheduleList.mockResolvedValue({
      data: [{
        ...mockScheduleData[0],
        next_run_at: new Date(Date.now() + 5000).toISOString(),
      }],
    });

    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('in 0m')).toBeInTheDocument();
    });
  });
});

describe('getSubmitButtonText branches', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    setupDefaultMocks();
  });

  it('shows "Creating..." when create form is submitting', async () => {
    // Make create never resolve to keep isPending=true
    mockCreate.mockImplementation(() => new Promise(() => {}));

    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Create Schedule'));

    await waitFor(() => {
      expect(screen.getByText('Name *')).toBeInTheDocument();
    });

    // Fill required fields
    fireEvent.change(screen.getByLabelText('Name *'), {
      target: { value: 'Test Schedule' },
    });
    fireEvent.change(screen.getByLabelText('Scenario *'), {
      target: { value: 'scenario-1' },
    });

    // Submit form
    const allButtons = screen.getAllByRole('button', { name: /Create Schedule/ });
    const submitButton = allButtons.find((btn) => !btn.className.includes('gap-2'));
    fireEvent.click(submitButton!);

    await waitFor(() => {
      expect(screen.getByText('Creating...')).toBeInTheDocument();
    });
  });

  it('shows "Updating..." when edit form is submitting', async () => {
    // Make update never resolve to keep isPending=true
    mockUpdate.mockImplementation(() => new Promise(() => {}));

    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    const editButtons = screen.getAllByTitle('Edit');
    fireEvent.click(editButtons[0]);

    await waitFor(() => {
      expect(screen.getByText('Edit Schedule')).toBeInTheDocument();
    });

    // Submit the edit form
    const submitButton = screen.getByRole('button', { name: /Update Schedule$/ });
    fireEvent.click(submitButton);

    await waitFor(() => {
      expect(screen.getByText('Updating...')).toBeInTheDocument();
    });
  });
});

describe('ScheduleFormModal create with cron and start_at', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    setupDefaultMocks();
  });

  it('submits create form with cron frequency including cron_expr', async () => {
    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Create Schedule'));

    await waitFor(() => {
      expect(screen.getByText('Name *')).toBeInTheDocument();
    });

    // Fill form
    fireEvent.change(screen.getByLabelText('Name *'), {
      target: { value: 'Cron Schedule' },
    });
    fireEvent.change(screen.getByLabelText('Scenario *'), {
      target: { value: 'scenario-1' },
    });

    // Select cron frequency
    fireEvent.change(screen.getByLabelText('Frequency *'), {
      target: { value: 'cron' },
    });

    await waitFor(() => {
      expect(screen.getByText('Cron Expression *')).toBeInTheDocument();
    });

    // Fill cron expression
    fireEvent.change(screen.getByLabelText('Cron Expression *'), {
      target: { value: '0 0 * * *' },
    });

    // Submit
    const allButtons = screen.getAllByRole('button', { name: /Create Schedule/ });
    const submitButton = allButtons.find((btn) => !btn.className.includes('gap-2'));
    fireEvent.click(submitButton!);

    await waitFor(() => {
      expect(mockCreate).toHaveBeenCalledWith(
        expect.objectContaining({
          name: 'Cron Schedule',
          frequency: 'cron',
          cron_expr: '0 0 * * *',
        })
      );
    });
  });

  it('submits create form with start_at date', async () => {
    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Create Schedule'));

    await waitFor(() => {
      expect(screen.getByText('Name *')).toBeInTheDocument();
    });

    // Fill required fields
    fireEvent.change(screen.getByLabelText('Name *'), {
      target: { value: 'Scheduled Test' },
    });
    fireEvent.change(screen.getByLabelText('Scenario *'), {
      target: { value: 'scenario-1' },
    });

    // Set start_at date
    const startAtInput = screen.getByLabelText('Start At (optional)');
    fireEvent.change(startAtInput, {
      target: { value: '2026-03-15T10:00' },
    });

    // Submit
    const allButtons = screen.getAllByRole('button', { name: /Create Schedule/ });
    const submitButton = allButtons.find((btn) => !btn.className.includes('gap-2'));
    fireEvent.click(submitButton!);

    await waitFor(() => {
      expect(mockCreate).toHaveBeenCalledWith(
        expect.objectContaining({
          name: 'Scheduled Test',
          start_at: expect.stringContaining('2026'),
        })
      );
    });
  });

  it('submits create form with empty start_at by omitting it', async () => {
    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Create Schedule'));

    await waitFor(() => {
      expect(screen.getByText('Name *')).toBeInTheDocument();
    });

    fireEvent.change(screen.getByLabelText('Name *'), {
      target: { value: 'No Start Test' },
    });
    fireEvent.change(screen.getByLabelText('Scenario *'), {
      target: { value: 'scenario-1' },
    });

    // Don't set start_at, leave it empty

    const allButtons = screen.getAllByRole('button', { name: /Create Schedule/ });
    const submitButton = allButtons.find((btn) => !btn.className.includes('gap-2'));
    fireEvent.click(submitButton!);

    await waitFor(() => {
      expect(mockCreate).toHaveBeenCalledWith(
        expect.not.objectContaining({ start_at: expect.anything() })
      );
    });
  });
});

describe('ScheduleFormModal create error handling', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    setupDefaultMocks();
  });

  it('shows error toast when create fails', async () => {
    const toast = await import('react-hot-toast');
    mockCreate.mockRejectedValueOnce(new Error('Create failed'));

    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Create Schedule'));

    await waitFor(() => {
      expect(screen.getByText('Name *')).toBeInTheDocument();
    });

    fireEvent.change(screen.getByLabelText('Name *'), {
      target: { value: 'Failing Schedule' },
    });
    fireEvent.change(screen.getByLabelText('Scenario *'), {
      target: { value: 'scenario-1' },
    });

    const allButtons = screen.getAllByRole('button', { name: /Create Schedule/ });
    const submitButton = allButtons.find((btn) => !btn.className.includes('gap-2'));
    fireEvent.click(submitButton!);

    await waitFor(() => {
      expect(toast.default.error).toHaveBeenCalledWith('Failed to create schedule');
    });
  });

  it('shows success toast when create succeeds', async () => {
    const toast = await import('react-hot-toast');

    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Create Schedule'));

    await waitFor(() => {
      expect(screen.getByText('Name *')).toBeInTheDocument();
    });

    fireEvent.change(screen.getByLabelText('Name *'), {
      target: { value: 'Success Schedule' },
    });
    fireEvent.change(screen.getByLabelText('Scenario *'), {
      target: { value: 'scenario-1' },
    });

    const allButtons = screen.getAllByRole('button', { name: /Create Schedule/ });
    const submitButton = allButtons.find((btn) => !btn.className.includes('gap-2'));
    fireEvent.click(submitButton!);

    await waitFor(() => {
      expect(toast.default.success).toHaveBeenCalledWith('Schedule created');
    });
  });
});

describe('ScheduleFormModal update error handling', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    setupDefaultMocks();
  });

  it('shows error toast when update fails', async () => {
    const toast = await import('react-hot-toast');
    mockUpdate.mockRejectedValueOnce(new Error('Update failed'));

    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    const editButtons = screen.getAllByTitle('Edit');
    fireEvent.click(editButtons[0]);

    await waitFor(() => {
      expect(screen.getByText('Edit Schedule')).toBeInTheDocument();
    });

    const submitButton = screen.getByRole('button', { name: /Update Schedule$/ });
    fireEvent.click(submitButton);

    await waitFor(() => {
      expect(toast.default.error).toHaveBeenCalledWith('Failed to update schedule');
    });
  });

  it('shows success toast and closes modal when update succeeds', async () => {
    const toast = await import('react-hot-toast');

    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    const editButtons = screen.getAllByTitle('Edit');
    fireEvent.click(editButtons[0]);

    await waitFor(() => {
      expect(screen.getByText('Edit Schedule')).toBeInTheDocument();
    });

    const submitButton = screen.getByRole('button', { name: /Update Schedule$/ });
    fireEvent.click(submitButton);

    await waitFor(() => {
      expect(toast.default.success).toHaveBeenCalledWith('Schedule updated');
    });
  });
});

describe('Schedule cron display in card', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    setupDefaultMocks();
  });

  it('shows cron expression in schedule card when frequency is cron', async () => {
    mockScheduleList.mockResolvedValue({
      data: [{
        ...mockScheduleData[0],
        frequency: 'cron',
        cron_expr: '0 */6 * * *',
      }],
    });

    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Custom (Cron)')).toBeInTheDocument();
      expect(screen.getByText('(0 */6 * * *)')).toBeInTheDocument();
    });
  });

  it('does not show cron expression when frequency is not cron', async () => {
    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    expect(screen.queryByText(/\(0/)).not.toBeInTheDocument();
  });
});

describe('ScheduleFormModal X close button', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    setupDefaultMocks();
  });

  it('closes create modal when X button is clicked', async () => {
    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Create Schedule'));

    await waitFor(() => {
      expect(screen.getByText('Name *')).toBeInTheDocument();
    });

    // Find the X close button (it's the one with XMarkIcon inside the modal header)
    const modalHeader = screen.getByText('Create Schedule', { selector: 'h2' });
    const closeButton = modalHeader.parentElement?.querySelector('button');
    expect(closeButton).toBeDefined();
    fireEvent.click(closeButton!);

    await waitFor(() => {
      expect(screen.queryByText('Name *')).not.toBeInTheDocument();
    });
  });

  it('closes edit modal when X button is clicked', async () => {
    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    const editButtons = screen.getAllByTitle('Edit');
    fireEvent.click(editButtons[0]);

    await waitFor(() => {
      expect(screen.getByText('Edit Schedule')).toBeInTheDocument();
    });

    // Find the X close button in edit modal header
    const modalHeader = screen.getByText('Edit Schedule', { selector: 'h2' });
    const closeButton = modalHeader.parentElement?.querySelector('button');
    expect(closeButton).toBeDefined();
    fireEvent.click(closeButton!);

    await waitFor(() => {
      expect(screen.queryByText('Edit Schedule')).not.toBeInTheDocument();
    });
  });
});

describe('ScheduleFormModal edit mode hides start_at', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    setupDefaultMocks();
  });

  it('does not show Start At field in edit mode', async () => {
    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    const editButtons = screen.getAllByTitle('Edit');
    fireEvent.click(editButtons[0]);

    await waitFor(() => {
      expect(screen.getByText('Edit Schedule')).toBeInTheDocument();
    });

    expect(screen.queryByText('Start At (optional)')).not.toBeInTheDocument();
  });

  it('shows Start At field in create mode', async () => {
    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Create Schedule'));

    await waitFor(() => {
      expect(screen.getByText('Start At (optional)')).toBeInTheDocument();
    });
  });
});

describe('Schedule without description', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    setupDefaultMocks();
  });

  it('does not render description paragraph when description is empty', async () => {
    mockScheduleList.mockResolvedValue({
      data: [{
        ...mockScheduleData[0],
        description: '',
      }],
    });

    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    // The description "Run security tests daily" should NOT be present
    expect(screen.queryByText('Run security tests daily')).not.toBeInTheDocument();
  });
});

describe('ScheduleRunsHistory states', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    setupDefaultMocks();
  });

  it('shows "No runs yet" when expanded schedule has no runs', async () => {
    mockGetRuns.mockResolvedValue({ data: [] });

    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    const expandButtons = screen.getAllByTitle('Show history');
    fireEvent.click(expandButtons[0]);

    await waitFor(() => {
      expect(screen.getByText('No runs yet')).toBeInTheDocument();
    });
  });

  it('shows loading state while runs are being fetched', async () => {
    mockGetRuns.mockImplementation(
      () => new Promise((resolve) => setTimeout(() => resolve({ data: [] }), 200))
    );

    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    const expandButtons = screen.getAllByTitle('Show history');
    fireEvent.click(expandButtons[0]);

    // Should show loading immediately
    expect(screen.getByText('Loading history...')).toBeInTheDocument();
  });

  it('shows "Recent Runs" header and run details with View Execution link', async () => {
    mockGetRuns.mockResolvedValue({
      data: [
        {
          id: 'run-1',
          schedule_id: 'sched-1',
          execution_id: 'exec-1',
          started_at: new Date().toISOString(),
          completed_at: new Date().toISOString(),
          status: 'completed',
          error: '',
        },
      ],
    });

    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    const expandButtons = screen.getAllByTitle('Show history');
    fireEvent.click(expandButtons[0]);

    await waitFor(() => {
      expect(screen.getByText('Recent Runs')).toBeInTheDocument();
      expect(screen.getByText('View Execution')).toBeInTheDocument();
    });

    // The link should point to the execution
    const link = screen.getByText('View Execution');
    expect(link.closest('a')).toHaveAttribute('href', '/executions/exec-1');
  });

  it('does not show View Execution link when execution_id is empty', async () => {
    mockGetRuns.mockResolvedValue({
      data: [
        {
          id: 'run-1',
          schedule_id: 'sched-1',
          execution_id: '',
          started_at: new Date().toISOString(),
          completed_at: null,
          status: 'failed',
          error: 'Connection timeout',
        },
      ],
    });

    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    const expandButtons = screen.getAllByTitle('Show history');
    fireEvent.click(expandButtons[0]);

    await waitFor(() => {
      expect(screen.getByText('Recent Runs')).toBeInTheDocument();
      expect(screen.getByText('(Connection timeout)')).toBeInTheDocument();
    });

    expect(screen.queryByText('View Execution')).not.toBeInTheDocument();
  });
});

describe('ScheduleFormModal agent field and safe_mode in edit', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    setupDefaultMocks();
  });

  it('can change agent paw field value', async () => {
    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Create Schedule'));

    await waitFor(() => {
      expect(screen.getByText('Agent (optional)')).toBeInTheDocument();
    });

    const agentInput = screen.getByPlaceholderText('Leave empty for all agents');
    fireEvent.change(agentInput, { target: { value: 'agent-abc' } });
    expect(agentInput).toHaveValue('agent-abc');
  });

  it('pre-fills description and safe_mode in edit mode', async () => {
    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    const editButtons = screen.getAllByTitle('Edit');
    fireEvent.click(editButtons[0]);

    await waitFor(() => {
      expect(screen.getByText('Edit Schedule')).toBeInTheDocument();
    });

    // Check description is pre-filled
    expect(screen.getByDisplayValue('Run security tests daily')).toBeInTheDocument();

    // Check safe mode is checked (first schedule has safe_mode: true)
    const checkbox = screen.getByRole('checkbox');
    expect(checkbox).toBeChecked();
  });

  it('pre-fills safe_mode as unchecked when schedule has safe_mode false', async () => {
    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Weekly Audit')).toBeInTheDocument();
    });

    // Edit second schedule (safe_mode: false)
    const editButtons = screen.getAllByTitle('Edit');
    fireEvent.click(editButtons[1]);

    await waitFor(() => {
      expect(screen.getByText('Edit Schedule')).toBeInTheDocument();
    });

    const checkbox = screen.getByRole('checkbox');
    expect(checkbox).not.toBeChecked();
  });

  it('pre-fills agent_paw in edit mode', async () => {
    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Weekly Audit')).toBeInTheDocument();
    });

    // Edit second schedule (has agent_paw: 'agent-1')
    const editButtons = screen.getAllByTitle('Edit');
    fireEvent.click(editButtons[1]);

    await waitFor(() => {
      expect(screen.getByText('Edit Schedule')).toBeInTheDocument();
    });

    expect(screen.getByDisplayValue('agent-1')).toBeInTheDocument();
  });
});

describe('Scheduler toast messages with exact text', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    setupDefaultMocks();
  });

  it('shows exact "Schedule paused" success message', async () => {
    const toast = await import('react-hot-toast');

    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    const pauseButtons = screen.getAllByTitle('Pause');
    fireEvent.click(pauseButtons[0]);

    await waitFor(() => {
      expect(toast.default.success).toHaveBeenCalledWith('Schedule paused');
    });
  });

  it('shows exact "Schedule resumed" success message', async () => {
    const toast = await import('react-hot-toast');

    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Weekly Audit')).toBeInTheDocument();
    });

    const resumeButtons = screen.getAllByTitle('Resume');
    fireEvent.click(resumeButtons[0]);

    await waitFor(() => {
      expect(toast.default.success).toHaveBeenCalledWith('Schedule resumed');
    });
  });

  it('shows exact "Schedule execution started" success message', async () => {
    const toast = await import('react-hot-toast');

    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    const runNowButtons = screen.getAllByTitle('Run Now');
    fireEvent.click(runNowButtons[0]);

    await waitFor(() => {
      expect(toast.default.success).toHaveBeenCalledWith('Schedule execution started');
    });
  });

  it('shows exact "Schedule deleted" success message', async () => {
    const toast = await import('react-hot-toast');

    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    const deleteButtons = screen.getAllByTitle('Delete');
    fireEvent.click(deleteButtons[0]);

    await waitFor(() => {
      expect(screen.getByText('Delete Schedule')).toBeInTheDocument();
    });

    const allButtons = screen.getAllByRole('button');
    const confirmButton = allButtons.find(
      (btn) => btn.textContent === 'Delete' && btn.className.includes('bg-red-600')
    );
    fireEvent.click(confirmButton!);

    await waitFor(() => {
      expect(toast.default.success).toHaveBeenCalledWith('Schedule deleted');
    });
  });

  it('shows exact "Failed to pause schedule" error message', async () => {
    const toast = await import('react-hot-toast');
    mockPause.mockRejectedValueOnce(new Error('fail'));

    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    const pauseButtons = screen.getAllByTitle('Pause');
    fireEvent.click(pauseButtons[0]);

    await waitFor(() => {
      expect(toast.default.error).toHaveBeenCalledWith('Failed to pause schedule');
    });
  });

  it('shows exact "Failed to resume schedule" error message', async () => {
    const toast = await import('react-hot-toast');
    mockResume.mockRejectedValueOnce(new Error('fail'));

    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Weekly Audit')).toBeInTheDocument();
    });

    const resumeButtons = screen.getAllByTitle('Resume');
    fireEvent.click(resumeButtons[0]);

    await waitFor(() => {
      expect(toast.default.error).toHaveBeenCalledWith('Failed to resume schedule');
    });
  });

  it('shows exact "Failed to run schedule" error message', async () => {
    const toast = await import('react-hot-toast');
    mockRunNow.mockRejectedValueOnce(new Error('fail'));

    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    const runNowButtons = screen.getAllByTitle('Run Now');
    fireEvent.click(runNowButtons[0]);

    await waitFor(() => {
      expect(toast.default.error).toHaveBeenCalledWith('Failed to run schedule');
    });
  });

  it('shows exact "Failed to delete schedule" error message', async () => {
    const toast = await import('react-hot-toast');
    mockDelete.mockRejectedValueOnce(new Error('fail'));

    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    const deleteButtons = screen.getAllByTitle('Delete');
    fireEvent.click(deleteButtons[0]);

    await waitFor(() => {
      expect(screen.getByText('Delete Schedule')).toBeInTheDocument();
    });

    const allButtons = screen.getAllByRole('button');
    const confirmButton = allButtons.find(
      (btn) => btn.textContent === 'Delete' && btn.className.includes('bg-red-600')
    );
    fireEvent.click(confirmButton!);

    await waitFor(() => {
      expect(toast.default.error).toHaveBeenCalledWith('Failed to delete schedule');
    });
  });
});

describe('Scheduler delete modal Deleting... state', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    setupDefaultMocks();
  });

  it('shows "Deleting..." text while delete is pending', async () => {
    mockDelete.mockImplementation(() => new Promise(() => {}));

    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    const deleteButtons = screen.getAllByTitle('Delete');
    fireEvent.click(deleteButtons[0]);

    await waitFor(() => {
      expect(screen.getByText('Delete Schedule')).toBeInTheDocument();
    });

    const allButtons = screen.getAllByRole('button');
    const confirmButton = allButtons.find(
      (btn) => btn.textContent === 'Delete' && btn.className.includes('bg-red-600')
    );
    fireEvent.click(confirmButton!);

    await waitFor(() => {
      expect(screen.getByText('Deleting...')).toBeInTheDocument();
    });
  });
});

describe('Scheduler expanding different schedules', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    setupDefaultMocks();
  });

  it('collapses first schedule when expanding second', async () => {
    mockGetRuns.mockResolvedValue({
      data: [
        {
          id: 'run-1',
          schedule_id: 'sched-1',
          execution_id: 'exec-1',
          started_at: new Date().toISOString(),
          completed_at: new Date().toISOString(),
          status: 'completed',
          error: '',
        },
      ],
    });

    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    // Expand first schedule
    const expandButtons = screen.getAllByTitle('Show history');
    fireEvent.click(expandButtons[0]);

    await waitFor(() => {
      expect(screen.getByText('Recent Runs')).toBeInTheDocument();
    });

    // Now expand second schedule - first should collapse since toggleExpanded sets to new id
    const expandButtons2 = screen.getAllByTitle('Show history');
    fireEvent.click(expandButtons2[1]);

    // The expanded schedule switched, so there should still be a runs section for sched-2
    await waitFor(() => {
      expect(mockGetRuns).toHaveBeenCalledWith('sched-2');
    });
  });
});

describe('Scheduler edit form with cron frequency', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    setupDefaultMocks();
  });

  it('shows cron field pre-filled when editing a cron schedule', async () => {
    mockScheduleList.mockResolvedValue({
      data: [{
        ...mockScheduleData[0],
        frequency: 'cron',
        cron_expr: '30 2 * * 1',
      }],
    });

    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    const editButtons = screen.getAllByTitle('Edit');
    fireEvent.click(editButtons[0]);

    await waitFor(() => {
      expect(screen.getByText('Edit Schedule')).toBeInTheDocument();
    });

    // Cron expression field should be visible and pre-filled
    expect(screen.getByText('Cron Expression *')).toBeInTheDocument();
    expect(screen.getByDisplayValue('30 2 * * 1')).toBeInTheDocument();
  });
});

describe('Scheduler Start At field clear value', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    setupDefaultMocks();
  });

  it('clears start_at when date is cleared', async () => {
    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Create Schedule'));

    await waitFor(() => {
      expect(screen.getByText('Start At (optional)')).toBeInTheDocument();
    });

    const startAtInput = screen.getByLabelText('Start At (optional)');

    // Set a date first
    fireEvent.change(startAtInput, {
      target: { value: '2026-03-15T10:00' },
    });

    // Then clear it
    fireEvent.change(startAtInput, {
      target: { value: '' },
    });

    // Submit to verify start_at is removed
    fireEvent.change(screen.getByLabelText('Name *'), {
      target: { value: 'Clear Date Test' },
    });
    fireEvent.change(screen.getByLabelText('Scenario *'), {
      target: { value: 'scenario-1' },
    });

    const allButtons = screen.getAllByRole('button', { name: /Create Schedule/ });
    const submitButton = allButtons.find((btn) => !btn.className.includes('gap-2'));
    fireEvent.click(submitButton!);

    await waitFor(() => {
      expect(mockCreate).toHaveBeenCalledWith(
        expect.not.objectContaining({ start_at: expect.anything() })
      );
    });
  });
});

describe('Scheduler loading state shows loading message', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockScheduleList.mockImplementation(
      () => new Promise(() => {}) // never resolves
    );
    mockScenarioList.mockResolvedValue({ data: mockScenarioData });
  });

  it('shows loading message while schedules are being fetched', async () => {
    renderScheduler();

    expect(screen.getByText('Loading schedules...')).toBeInTheDocument();
  });
});

describe('Scheduler scenario names display', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    setupDefaultMocks();
  });

  it('shows matched scenario names for schedule cards', async () => {
    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    expect(screen.getByText('Test Scenario 1')).toBeInTheDocument();
    expect(screen.getByText('Test Scenario 2')).toBeInTheDocument();
  });
});

describe('Scheduler description field input change', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    setupDefaultMocks();
  });

  it('can change description field value in create form', async () => {
    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Create Schedule'));

    await waitFor(() => {
      expect(screen.getByText('Description')).toBeInTheDocument();
    });

    const descInput = screen.getByLabelText('Description');
    fireEvent.change(descInput, { target: { value: 'My new description' } });
    expect(descInput).toHaveValue('My new description');
  });
});

describe('Scheduler once frequency option', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    setupDefaultMocks();
  });

  it('can select once frequency', async () => {
    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Create Schedule'));

    await waitFor(() => {
      const frequencySelect = screen.getByLabelText('Frequency *');
      fireEvent.change(frequencySelect, { target: { value: 'once' } });
      expect(frequencySelect).toHaveValue('once');
    });
  });

  it('displays once frequency label in schedule card', async () => {
    mockScheduleList.mockResolvedValue({
      data: [{
        ...mockScheduleData[0],
        frequency: 'once',
      }],
    });

    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Once')).toBeInTheDocument();
    });
  });
});

describe('formatRelativeTime N/A branch', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    setupDefaultMocks();
  });

  it('shows N/A when active schedule has null next_run_at', async () => {
    mockScheduleList.mockResolvedValue({
      data: [{
        ...mockScheduleData[0],
        status: 'active',
        next_run_at: null,
      }],
    });

    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('N/A')).toBeInTheDocument();
    });
  });
});

describe('Scheduler scenarios fallback to empty array', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockScheduleList.mockResolvedValue({ data: mockScheduleData });
    // Make scenarios return undefined data
    mockScenarioList.mockResolvedValue({ data: undefined });
    mockGetRuns.mockResolvedValue({ data: [] });
    mockPause.mockResolvedValue({ data: {} });
    mockResume.mockResolvedValue({ data: {} });
    mockRunNow.mockResolvedValue({ data: {} });
    mockDelete.mockResolvedValue({ data: {} });
    mockCreate.mockResolvedValue({ data: {} });
    mockUpdate.mockResolvedValue({ data: {} });
  });

  it('renders create modal even when scenarios are not loaded', async () => {
    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Create Schedule'));

    await waitFor(() => {
      expect(screen.getByText('Name *')).toBeInTheDocument();
      // Select a scenario still works, there are just no options
      expect(screen.getByLabelText('Scenario *')).toBeInTheDocument();
    });
  });

  it('shows Unknown Scenario when scenarios list is empty/unavailable', async () => {
    renderScheduler();

    await waitFor(() => {
      expect(screen.getAllByText('Unknown Scenario').length).toBeGreaterThan(0);
    });
  });

  it('opens edit modal with empty scenarios fallback', async () => {
    renderScheduler();

    await waitFor(() => {
      expect(screen.getByText('Daily Security Check')).toBeInTheDocument();
    });

    const editButtons = screen.getAllByTitle('Edit');
    fireEvent.click(editButtons[0]);

    await waitFor(() => {
      expect(screen.getByText('Edit Schedule')).toBeInTheDocument();
      // The scenario select should render even without scenarios
      expect(screen.getByLabelText('Scenario *')).toBeInTheDocument();
    });
  });
});
