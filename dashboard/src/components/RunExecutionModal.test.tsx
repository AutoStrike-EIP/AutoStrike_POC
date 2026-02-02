import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { RunExecutionModal } from './RunExecutionModal';
import { api } from '../lib/api';
import { Scenario } from '../types';

// Mock the API
vi.mock('../lib/api', () => ({
  api: {
    get: vi.fn(),
  },
}));

const createTestQueryClient = () =>
  new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
      },
    },
  });

const mockScenario: Scenario = {
  id: 'scenario-1',
  name: 'Test Scenario',
  description: 'A test scenario',
  phases: [
    { name: 'Phase 1', techniques: ['T1082', 'T1083'] },
    { name: 'Phase 2', techniques: ['T1057'] },
  ],
  tags: ['test'],
};

const mockAgents = [
  {
    paw: 'agent-001-uuid-test',
    hostname: 'workstation-1',
    platform: 'windows',
    status: 'online',
    last_seen: new Date().toISOString(),
    os_info: '',
    username: 'user1',
    privilege: 'user',
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString(),
  },
  {
    paw: 'agent-002-uuid-test',
    hostname: 'server-1',
    platform: 'linux',
    status: 'online',
    last_seen: new Date().toISOString(),
    os_info: '',
    username: 'root',
    privilege: 'admin',
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString(),
  },
  {
    paw: 'agent-003-uuid-test',
    hostname: 'offline-host',
    platform: 'windows',
    status: 'offline',
    last_seen: new Date().toISOString(),
    os_info: '',
    username: 'user2',
    privilege: 'user',
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString(),
  },
];

interface RenderModalProps {
  onConfirm?: (agentPaws: string[], safeMode: boolean) => void;
  onCancel?: () => void;
  isLoading?: boolean;
}

function renderModal(props: RenderModalProps = {}) {
  const testQueryClient = createTestQueryClient();
  const defaultProps = {
    scenario: mockScenario,
    onConfirm: vi.fn(),
    onCancel: vi.fn(),
    isLoading: false,
    ...props,
  };

  return {
    ...render(
      <QueryClientProvider client={testQueryClient}>
        <RunExecutionModal {...defaultProps} />
      </QueryClientProvider>
    ),
    ...defaultProps,
  };
}

describe('RunExecutionModal', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders modal with scenario information', async () => {
    vi.mocked(api.get).mockResolvedValue({ data: mockAgents } as never);

    renderModal();

    expect(screen.getByText('Run Scenario')).toBeInTheDocument();
    expect(screen.getByText('Test Scenario')).toBeInTheDocument();
    expect(screen.getByText('2 phases, 3 techniques')).toBeInTheDocument();
  });

  it('renders loading state while fetching agents', () => {
    vi.mocked(api.get).mockReturnValue(new Promise(() => {}) as never);

    renderModal();

    expect(screen.getByText('Loading agents...')).toBeInTheDocument();
  });

  it('renders online agents only', async () => {
    vi.mocked(api.get).mockResolvedValue({ data: mockAgents } as never);

    renderModal();

    await waitFor(() => {
      expect(screen.getByText('workstation-1')).toBeInTheDocument();
      expect(screen.getByText('server-1')).toBeInTheDocument();
    });

    // Offline agent should not be rendered
    expect(screen.queryByText('offline-host')).not.toBeInTheDocument();
  });

  it('renders empty state when no online agents', async () => {
    const offlineAgents = mockAgents.map(a => ({ ...a, status: 'offline' }));
    vi.mocked(api.get).mockResolvedValue({ data: offlineAgents } as never);

    renderModal();

    await waitFor(() => {
      expect(screen.getByText(/No online agents available/)).toBeInTheDocument();
    });
  });

  it('toggles agent selection', async () => {
    vi.mocked(api.get).mockResolvedValue({ data: mockAgents } as never);

    renderModal();

    await waitFor(() => {
      expect(screen.getByText('workstation-1')).toBeInTheDocument();
    });

    const checkbox = screen.getAllByRole('checkbox')[0];
    expect(checkbox).not.toBeChecked();

    fireEvent.click(checkbox);
    expect(checkbox).toBeChecked();

    fireEvent.click(checkbox);
    expect(checkbox).not.toBeChecked();
  });

  it('selects all agents', async () => {
    vi.mocked(api.get).mockResolvedValue({ data: mockAgents } as never);

    renderModal();

    await waitFor(() => {
      expect(screen.getByText('Select All')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Select All'));

    const checkboxes = screen.getAllByRole('checkbox');
    // First 2 are agent checkboxes, third is safe mode
    expect(checkboxes[0]).toBeChecked();
    expect(checkboxes[1]).toBeChecked();
  });

  it('deselects all agents', async () => {
    vi.mocked(api.get).mockResolvedValue({ data: mockAgents } as never);

    renderModal();

    await waitFor(() => {
      expect(screen.getByText('Select All')).toBeInTheDocument();
    });

    // Select all first
    fireEvent.click(screen.getByText('Select All'));

    await waitFor(() => {
      expect(screen.getByText('Deselect All')).toBeInTheDocument();
    });

    // Then deselect all
    fireEvent.click(screen.getByText('Deselect All'));

    const checkboxes = screen.getAllByRole('checkbox');
    expect(checkboxes[0]).not.toBeChecked();
    expect(checkboxes[1]).not.toBeChecked();
  });

  it('toggles safe mode', async () => {
    vi.mocked(api.get).mockResolvedValue({ data: mockAgents } as never);

    renderModal();

    await waitFor(() => {
      expect(screen.getByText('Safe Mode')).toBeInTheDocument();
    });

    const safeModeCheckbox = screen.getByRole('checkbox', { name: /Safe Mode/i });
    expect(safeModeCheckbox).toBeChecked(); // Safe mode is on by default

    fireEvent.click(safeModeCheckbox);
    expect(safeModeCheckbox).not.toBeChecked();
  });

  it('shows warning when no agent selected', async () => {
    vi.mocked(api.get).mockResolvedValue({ data: mockAgents } as never);

    renderModal();

    await waitFor(() => {
      expect(screen.getByText(/Please select at least one agent/)).toBeInTheDocument();
    });
  });

  it('disables run button when no agent selected', async () => {
    vi.mocked(api.get).mockResolvedValue({ data: mockAgents } as never);

    renderModal();

    await waitFor(() => {
      expect(screen.getByText('workstation-1')).toBeInTheDocument();
    });

    const runButton = screen.getByRole('button', { name: /Run on 0 agents/i });
    expect(runButton).toBeDisabled();
  });

  it('enables run button when agent is selected', async () => {
    vi.mocked(api.get).mockResolvedValue({ data: mockAgents } as never);

    renderModal();

    await waitFor(() => {
      expect(screen.getByText('workstation-1')).toBeInTheDocument();
    });

    const checkbox = screen.getAllByRole('checkbox')[0];
    fireEvent.click(checkbox);

    const runButton = screen.getByRole('button', { name: /Run on 1 agent$/i });
    expect(runButton).toBeEnabled();
  });

  it('calls onConfirm with selected agents and safe mode', async () => {
    vi.mocked(api.get).mockResolvedValue({ data: mockAgents } as never);

    const onConfirm = vi.fn();
    renderModal({ onConfirm });

    await waitFor(() => {
      expect(screen.getByText('workstation-1')).toBeInTheDocument();
    });

    // Select first agent
    const checkbox = screen.getAllByRole('checkbox')[0];
    fireEvent.click(checkbox);

    // Click run
    const runButton = screen.getByRole('button', { name: /Run on 1 agent$/i });
    fireEvent.click(runButton);

    expect(onConfirm).toHaveBeenCalledWith(['agent-001-uuid-test'], true);
  });

  it('calls onConfirm with safe mode off', async () => {
    vi.mocked(api.get).mockResolvedValue({ data: mockAgents } as never);

    const onConfirm = vi.fn();
    renderModal({ onConfirm });

    await waitFor(() => {
      expect(screen.getByText('workstation-1')).toBeInTheDocument();
    });

    // Select first agent
    const checkbox = screen.getAllByRole('checkbox')[0];
    fireEvent.click(checkbox);

    // Turn off safe mode
    const safeModeCheckbox = screen.getByRole('checkbox', { name: /Safe Mode/i });
    fireEvent.click(safeModeCheckbox);

    // Click run
    const runButton = screen.getByRole('button', { name: /Run on 1 agent$/i });
    fireEvent.click(runButton);

    expect(onConfirm).toHaveBeenCalledWith(['agent-001-uuid-test'], false);
  });

  it('calls onCancel when cancel button clicked', async () => {
    vi.mocked(api.get).mockResolvedValue({ data: mockAgents } as never);

    const onCancel = vi.fn();
    renderModal({ onCancel });

    await waitFor(() => {
      expect(screen.getByText('Cancel')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Cancel'));

    expect(onCancel).toHaveBeenCalled();
  });

  it('calls onCancel when overlay clicked', async () => {
    vi.mocked(api.get).mockResolvedValue({ data: mockAgents } as never);

    const onCancel = vi.fn();
    renderModal({ onCancel });

    await waitFor(() => {
      expect(screen.getByText('workstation-1')).toBeInTheDocument();
    });

    // Click the overlay (aria-hidden div)
    const overlay = document.querySelector('[aria-hidden="true"]');
    if (overlay) {
      fireEvent.click(overlay);
    }

    expect(onCancel).toHaveBeenCalled();
  });

  it('shows loading state when isLoading is true', async () => {
    vi.mocked(api.get).mockResolvedValue({ data: mockAgents } as never);

    renderModal({ isLoading: true });

    await waitFor(() => {
      expect(screen.getByText('workstation-1')).toBeInTheDocument();
    });

    // Select an agent to enable the button
    const checkbox = screen.getAllByRole('checkbox')[0];
    fireEvent.click(checkbox);

    const runButton = screen.getByRole('button', { name: /Starting.../i });
    expect(runButton).toBeDisabled();
  });

  it('disables cancel button when isLoading', async () => {
    vi.mocked(api.get).mockResolvedValue({ data: mockAgents } as never);

    renderModal({ isLoading: true });

    await waitFor(() => {
      expect(screen.getByText('Cancel')).toBeInTheDocument();
    });

    const cancelButton = screen.getByText('Cancel');
    expect(cancelButton).toBeDisabled();
  });

  it('displays correct button text for multiple agents', async () => {
    vi.mocked(api.get).mockResolvedValue({ data: mockAgents } as never);

    renderModal();

    await waitFor(() => {
      expect(screen.getByText('Select All')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Select All'));

    const runButton = screen.getByRole('button', { name: /Run on 2 agents/i });
    expect(runButton).toBeInTheDocument();
  });

  it('displays agent platform information', async () => {
    vi.mocked(api.get).mockResolvedValue({ data: mockAgents } as never);

    renderModal();

    await waitFor(() => {
      expect(screen.getByText('workstation-1')).toBeInTheDocument();
      expect(screen.getByText('server-1')).toBeInTheDocument();
    });

    // Check that platform info exists in the rendered output
    const agentLabels = document.querySelectorAll('.text-xs.text-gray-500');
    const platformTexts = Array.from(agentLabels).map(el => el.textContent);
    expect(platformTexts.some(t => t?.includes('windows'))).toBe(true);
    expect(platformTexts.some(t => t?.includes('linux'))).toBe(true);
  });
});
