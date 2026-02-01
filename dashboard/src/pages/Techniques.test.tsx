import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import Techniques from './Techniques';
import { api } from '../lib/api';

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

function renderWithClient(ui: React.ReactElement) {
  const testQueryClient = createTestQueryClient();
  return render(
    <QueryClientProvider client={testQueryClient}>{ui}</QueryClientProvider>
  );
}

describe('Techniques Page', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders loading state', () => {
    vi.mocked(api.get).mockReturnValue(new Promise(() => {}) as never);

    renderWithClient(<Techniques />);
    expect(screen.getByText('Loading techniques...')).toBeInTheDocument();
  });

  it('renders techniques list', async () => {
    const mockTechniques = [
      {
        id: 'T1059.001',
        name: 'PowerShell',
        description: 'Adversaries may use PowerShell to execute commands',
        tactic: 'execution',
        platforms: ['windows'],
        is_safe: true,
      },
    ];
    vi.mocked(api.get).mockResolvedValue({ data: mockTechniques } as never);

    renderWithClient(<Techniques />);

    expect(await screen.findByText('T1059.001')).toBeInTheDocument();
    expect(screen.getByText('PowerShell')).toBeInTheDocument();
    expect(screen.getByText('execution')).toBeInTheDocument();
    expect(screen.getByText('windows')).toBeInTheDocument();
    // "Safe" appears both as column header and as badge - check we have both
    const safeElements = screen.getAllByText('Safe');
    expect(safeElements.length).toBeGreaterThanOrEqual(2);
  });

  it('renders unsafe technique badge', async () => {
    const mockTechniques = [
      {
        id: 'T1055',
        name: 'Process Injection',
        description: 'Adversaries may inject code into processes',
        tactic: 'defense_evasion',
        platforms: ['windows', 'linux'],
        is_safe: false,
      },
    ];
    vi.mocked(api.get).mockResolvedValue({ data: mockTechniques } as never);

    renderWithClient(<Techniques />);

    expect(await screen.findByText('Unsafe')).toBeInTheDocument();
    expect(screen.getByText('defense evasion')).toBeInTheDocument();
  });

  it('renders multiple platforms', async () => {
    const mockTechniques = [
      {
        id: 'T1082',
        name: 'System Information Discovery',
        description: 'Gather system information',
        tactic: 'discovery',
        platforms: ['windows', 'linux', 'darwin'],
        is_safe: true,
      },
    ];
    vi.mocked(api.get).mockResolvedValue({ data: mockTechniques } as never);

    renderWithClient(<Techniques />);

    expect(await screen.findByText('windows')).toBeInTheDocument();
    expect(screen.getByText('linux')).toBeInTheDocument();
    expect(screen.getByText('darwin')).toBeInTheDocument();
  });

  it('renders empty state when no techniques', async () => {
    vi.mocked(api.get).mockResolvedValue({ data: [] } as never);

    renderWithClient(<Techniques />);

    expect(await screen.findByText('No techniques loaded')).toBeInTheDocument();
    expect(screen.getByText('Import techniques from Atomic Red Team')).toBeInTheDocument();
  });

  it('renders page title and import button', async () => {
    vi.mocked(api.get).mockResolvedValue({ data: [] } as never);

    renderWithClient(<Techniques />);

    expect(await screen.findByText('Techniques')).toBeInTheDocument();
    expect(screen.getByText('Import Techniques')).toBeInTheDocument();
  });

  it('applies correct tactic colors', async () => {
    const mockTechniques = [
      {
        id: 'T1595',
        name: 'Active Scanning',
        description: 'Scan target',
        tactic: 'reconnaissance',
        platforms: ['windows'],
        is_safe: true,
      },
    ];
    vi.mocked(api.get).mockResolvedValue({ data: mockTechniques } as never);

    renderWithClient(<Techniques />);

    const tacticBadge = await screen.findByText('reconnaissance');
    expect(tacticBadge).toHaveClass('bg-purple-100');
  });

  it('falls back to default color for unknown tactic', async () => {
    const mockTechniques = [
      {
        id: 'T9999',
        name: 'Unknown',
        description: 'Unknown technique',
        tactic: 'unknown_tactic',
        platforms: ['windows'],
        is_safe: true,
      },
    ];
    vi.mocked(api.get).mockResolvedValue({ data: mockTechniques } as never);

    renderWithClient(<Techniques />);

    const tacticBadge = await screen.findByText('unknown tactic');
    expect(tacticBadge).toHaveClass('bg-gray-100');
  });
});
