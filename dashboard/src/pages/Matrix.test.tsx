import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import Matrix from './Matrix';
import { api } from '../lib/api';

// Mock the API
vi.mock('../lib/api', () => ({
  api: {
    get: vi.fn(),
  },
}));

// Mock MitreMatrix component to simplify testing the page
vi.mock('../components/MitreMatrix', () => ({
  MitreMatrix: ({ techniques }: { techniques: unknown[] }) => (
    <div data-testid="mitre-matrix">MitreMatrix with {techniques.length} techniques</div>
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

function renderMatrix() {
  const testQueryClient = createTestQueryClient();
  return render(
    <QueryClientProvider client={testQueryClient}>
      <Matrix />
    </QueryClientProvider>
  );
}

const mockTechniques = [
  {
    id: 'T1082',
    name: 'System Information Discovery',
    description: 'Test description',
    tactic: 'discovery',
    platforms: ['windows', 'linux'],
    is_safe: true,
    detection: [],
  },
  {
    id: 'T1083',
    name: 'File and Directory Discovery',
    description: 'Test description',
    tactic: 'discovery',
    platforms: ['windows', 'linux'],
    is_safe: true,
    detection: [],
  },
  {
    id: 'T1059.001',
    name: 'PowerShell',
    description: 'Test description',
    tactic: 'execution',
    platforms: ['windows'],
    is_safe: false,
    detection: [],
  },
];

const mockCoverage = {
  reconnaissance: 0,
  resource_development: 0,
  initial_access: 0,
  execution: 1,
  persistence: 0,
  privilege_escalation: 0,
  defense_evasion: 0,
  credential_access: 0,
  discovery: 2,
  lateral_movement: 0,
  collection: 0,
  command_and_control: 0,
  exfiltration: 0,
  impact: 0,
};

describe('Matrix Page', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders loading state', () => {
    vi.mocked(api.get).mockImplementation(() => new Promise(() => {}) as never);

    renderMatrix();

    expect(screen.getByText('Loading MITRE ATT&CK matrix...')).toBeInTheDocument();
  });

  it('renders page title and description', async () => {
    vi.mocked(api.get).mockImplementation((url: string) => {
      if (url === '/techniques') {
        return Promise.resolve({ data: mockTechniques }) as never;
      }
      if (url === '/techniques/coverage') {
        return Promise.resolve({ data: mockCoverage }) as never;
      }
      return Promise.reject(new Error('Unknown endpoint')) as never;
    });

    renderMatrix();

    await waitFor(() => {
      expect(screen.getByText('MITRE ATT&CK Matrix')).toBeInTheDocument();
    });

    expect(screen.getByText('Interactive visualization of available attack techniques')).toBeInTheDocument();
  });

  it('renders technique count stat', async () => {
    vi.mocked(api.get).mockImplementation((url: string) => {
      if (url === '/techniques') {
        return Promise.resolve({ data: mockTechniques }) as never;
      }
      if (url === '/techniques/coverage') {
        return Promise.resolve({ data: mockCoverage }) as never;
      }
      return Promise.reject(new Error('Unknown endpoint')) as never;
    });

    renderMatrix();

    await waitFor(() => {
      expect(screen.getByText('3')).toBeInTheDocument();
    });

    expect(screen.getByText('Techniques')).toBeInTheDocument();
  });

  it('renders tactic coverage stat', async () => {
    vi.mocked(api.get).mockImplementation((url: string) => {
      if (url === '/techniques') {
        return Promise.resolve({ data: mockTechniques }) as never;
      }
      if (url === '/techniques/coverage') {
        return Promise.resolve({ data: mockCoverage }) as never;
      }
      return Promise.reject(new Error('Unknown endpoint')) as never;
    });

    renderMatrix();

    await waitFor(() => {
      expect(screen.getByText('2/14')).toBeInTheDocument();
    });

    expect(screen.getByText('Tactics Covered')).toBeInTheDocument();
  });

  it('renders MitreMatrix component with techniques', async () => {
    vi.mocked(api.get).mockImplementation((url: string) => {
      if (url === '/techniques') {
        return Promise.resolve({ data: mockTechniques }) as never;
      }
      if (url === '/techniques/coverage') {
        return Promise.resolve({ data: mockCoverage }) as never;
      }
      return Promise.reject(new Error('Unknown endpoint')) as never;
    });

    renderMatrix();

    await waitFor(() => {
      expect(screen.getByTestId('mitre-matrix')).toBeInTheDocument();
    });

    expect(screen.getByText('MitreMatrix with 3 techniques')).toBeInTheDocument();
  });

  it('renders empty state when no techniques', async () => {
    vi.mocked(api.get).mockImplementation((url: string) => {
      if (url === '/techniques') {
        return Promise.resolve({ data: [] }) as never;
      }
      if (url === '/techniques/coverage') {
        return Promise.resolve({ data: {} }) as never;
      }
      return Promise.reject(new Error('Unknown endpoint')) as never;
    });

    renderMatrix();

    await waitFor(() => {
      expect(screen.getByText('No techniques loaded')).toBeInTheDocument();
    });

    expect(screen.getByText('Import techniques to see the MITRE ATT&CK matrix')).toBeInTheDocument();
  });

  it('displays 0 techniques when empty', async () => {
    vi.mocked(api.get).mockImplementation((url: string) => {
      if (url === '/techniques') {
        return Promise.resolve({ data: [] }) as never;
      }
      if (url === '/techniques/coverage') {
        return Promise.resolve({ data: {} }) as never;
      }
      return Promise.reject(new Error('Unknown endpoint')) as never;
    });

    renderMatrix();

    await waitFor(() => {
      expect(screen.getByText('0')).toBeInTheDocument();
    });

    expect(screen.getByText('Techniques')).toBeInTheDocument();
  });

  it('displays 0/14 tactics when no coverage', async () => {
    vi.mocked(api.get).mockImplementation((url: string) => {
      if (url === '/techniques') {
        return Promise.resolve({ data: [] }) as never;
      }
      if (url === '/techniques/coverage') {
        return Promise.resolve({ data: {} }) as never;
      }
      return Promise.reject(new Error('Unknown endpoint')) as never;
    });

    renderMatrix();

    await waitFor(() => {
      expect(screen.getByText('0/14')).toBeInTheDocument();
    });
  });

  it('calls API endpoints on mount', async () => {
    vi.mocked(api.get).mockImplementation((url: string) => {
      if (url === '/techniques') {
        return Promise.resolve({ data: mockTechniques }) as never;
      }
      if (url === '/techniques/coverage') {
        return Promise.resolve({ data: mockCoverage }) as never;
      }
      return Promise.reject(new Error('Unknown endpoint')) as never;
    });

    renderMatrix();

    await waitFor(() => {
      expect(api.get).toHaveBeenCalledWith('/techniques');
      expect(api.get).toHaveBeenCalledWith('/techniques/coverage');
    });
  });

  it('handles coverage data with all zero counts', async () => {
    const zeroCoverage = {
      reconnaissance: 0,
      resource_development: 0,
      initial_access: 0,
      execution: 0,
      persistence: 0,
      privilege_escalation: 0,
      defense_evasion: 0,
      credential_access: 0,
      discovery: 0,
      lateral_movement: 0,
      collection: 0,
      command_and_control: 0,
      exfiltration: 0,
      impact: 0,
    };

    vi.mocked(api.get).mockImplementation((url: string) => {
      if (url === '/techniques') {
        return Promise.resolve({ data: mockTechniques }) as never;
      }
      if (url === '/techniques/coverage') {
        return Promise.resolve({ data: zeroCoverage }) as never;
      }
      return Promise.reject(new Error('Unknown endpoint')) as never;
    });

    renderMatrix();

    await waitFor(() => {
      expect(screen.getByText('0/14')).toBeInTheDocument();
    });
  });

  it('handles full coverage data', async () => {
    const fullCoverage = {
      reconnaissance: 5,
      resource_development: 3,
      initial_access: 7,
      execution: 10,
      persistence: 8,
      privilege_escalation: 6,
      defense_evasion: 12,
      credential_access: 4,
      discovery: 9,
      lateral_movement: 5,
      collection: 6,
      command_and_control: 8,
      exfiltration: 3,
      impact: 4,
    };

    vi.mocked(api.get).mockImplementation((url: string) => {
      if (url === '/techniques') {
        return Promise.resolve({ data: mockTechniques }) as never;
      }
      if (url === '/techniques/coverage') {
        return Promise.resolve({ data: fullCoverage }) as never;
      }
      return Promise.reject(new Error('Unknown endpoint')) as never;
    });

    renderMatrix();

    await waitFor(() => {
      expect(screen.getByText('14/14')).toBeInTheDocument();
    });
  });

  it('handles null coverage data gracefully', async () => {
    vi.mocked(api.get).mockImplementation((url: string) => {
      if (url === '/techniques') {
        return Promise.resolve({ data: mockTechniques }) as never;
      }
      if (url === '/techniques/coverage') {
        return Promise.resolve({ data: null }) as never;
      }
      return Promise.reject(new Error('Unknown endpoint')) as never;
    });

    renderMatrix();

    await waitFor(() => {
      expect(screen.getByText('0/14')).toBeInTheDocument();
    });
  });

  it('handles undefined techniques gracefully', async () => {
    vi.mocked(api.get).mockImplementation((url: string) => {
      if (url === '/techniques') {
        return Promise.resolve({ data: undefined }) as never;
      }
      if (url === '/techniques/coverage') {
        return Promise.resolve({ data: mockCoverage }) as never;
      }
      return Promise.reject(new Error('Unknown endpoint')) as never;
    });

    renderMatrix();

    await waitFor(() => {
      expect(screen.getByText('No techniques loaded')).toBeInTheDocument();
    });
  });
});
