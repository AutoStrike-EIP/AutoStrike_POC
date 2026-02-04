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

// Mock the schedule and scenario API
vi.mock('../lib/api', () => ({
  scheduleApi: {
    list: vi.fn(() =>
      Promise.resolve({
        data: [
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
        ],
      })
    ),
    getRuns: vi.fn(() => Promise.resolve({ data: [] })),
    pause: vi.fn(() => Promise.resolve({ data: {} })),
    resume: vi.fn(() => Promise.resolve({ data: {} })),
    runNow: vi.fn(() => Promise.resolve({ data: {} })),
    delete: vi.fn(() => Promise.resolve({ data: {} })),
    create: vi.fn(() => Promise.resolve({ data: {} })),
    update: vi.fn(() => Promise.resolve({ data: {} })),
  },
  scenarioApi: {
    list: vi.fn(() =>
      Promise.resolve({
        data: [
          { id: 'scenario-1', name: 'Test Scenario 1', description: '', phases: [], tags: [] },
          { id: 'scenario-2', name: 'Test Scenario 2', description: '', phases: [], tags: [] },
        ],
      })
    ),
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

describe('Scheduler Page', () => {
  beforeEach(() => {
    vi.clearAllMocks();
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
