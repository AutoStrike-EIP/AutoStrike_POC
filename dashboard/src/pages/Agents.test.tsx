import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import Agents from './Agents';
import { api } from '../lib/api';

// Mock the API
vi.mock('../lib/api', () => ({
  api: {
    get: vi.fn(),
  },
}));

// Mock date-fns
vi.mock('date-fns', () => ({
  formatDistanceToNow: vi.fn(() => '5 minutes ago'),
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

describe('Agents Page', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders loading state', () => {
    vi.mocked(api.get).mockReturnValue(new Promise(() => {}) as never);

    renderWithClient(<Agents />);
    expect(screen.getByText('Loading agents...')).toBeInTheDocument();
  });

  it('renders agents list', async () => {
    const mockAgents = [
      {
        paw: 'agent-123-456-789',
        hostname: 'DESKTOP-TEST',
        username: 'testuser',
        platform: 'windows',
        executors: ['psh', 'cmd'],
        status: 'online',
        last_seen: '2024-01-15T12:00:00Z',
      },
    ];
    vi.mocked(api.get).mockResolvedValue({ data: mockAgents } as never);

    renderWithClient(<Agents />);

    expect(await screen.findByText('DESKTOP-TEST')).toBeInTheDocument();
    expect(screen.getByText('testuser')).toBeInTheDocument();
    expect(screen.getByText('online')).toBeInTheDocument();
    expect(screen.getByText('windows')).toBeInTheDocument();
    expect(screen.getByText('psh')).toBeInTheDocument();
    expect(screen.getByText('cmd')).toBeInTheDocument();
  });

  it('renders offline agent with correct badge', async () => {
    const mockAgents = [
      {
        paw: 'agent-offline',
        hostname: 'OFFLINE-PC',
        username: 'user',
        platform: 'linux',
        executors: ['bash'],
        status: 'offline',
        last_seen: '2024-01-10T12:00:00Z',
      },
    ];
    vi.mocked(api.get).mockResolvedValue({ data: mockAgents } as never);

    renderWithClient(<Agents />);

    expect(await screen.findByText('offline')).toBeInTheDocument();
  });

  it('renders empty state when no agents', async () => {
    vi.mocked(api.get).mockResolvedValue({ data: [] } as never);

    renderWithClient(<Agents />);

    expect(await screen.findByText('No agents connected')).toBeInTheDocument();
    expect(screen.getByText('Deploy an agent to get started')).toBeInTheDocument();
  });

  it('renders page title and add button', async () => {
    vi.mocked(api.get).mockResolvedValue({ data: [] } as never);

    renderWithClient(<Agents />);

    expect(await screen.findByText('Agents')).toBeInTheDocument();
    expect(screen.getByText('Add Agent')).toBeInTheDocument();
  });

  it('displays truncated PAW', async () => {
    const mockAgents = [
      {
        paw: 'abcdefgh-1234-5678-9012-ijklmnopqrst',
        hostname: 'TEST-PC',
        username: 'user',
        platform: 'windows',
        executors: ['cmd'],
        status: 'online',
        last_seen: '2024-01-15T12:00:00Z',
      },
    ];
    vi.mocked(api.get).mockResolvedValue({ data: mockAgents } as never);

    renderWithClient(<Agents />);

    expect(await screen.findByText('PAW: abcdefgh...')).toBeInTheDocument();
  });
});
