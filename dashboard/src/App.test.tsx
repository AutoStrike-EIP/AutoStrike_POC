import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import App from './App';
import { AuthProvider } from './contexts/AuthContext';
import { ThemeProvider } from './contexts/ThemeContext';
import { authApi } from './lib/api';

// Mock the API
vi.mock('./lib/api', () => ({
  api: {
    get: vi.fn(),
  },
  authApi: {
    login: vi.fn(),
    me: vi.fn(),
    logout: vi.fn(),
    refresh: vi.fn(),
  },
}));

// Mock localStorage
const localStorageMock = {
  getItem: vi.fn(),
  setItem: vi.fn(),
  removeItem: vi.fn(),
  clear: vi.fn(),
  length: 0,
  key: vi.fn(),
};
Object.defineProperty(window, 'localStorage', { value: localStorageMock });

// Mock matchMedia for ThemeContext
const matchMediaMock = vi.fn().mockReturnValue({
  matches: false,
  addEventListener: vi.fn(),
  removeEventListener: vi.fn(),
});

const createQueryClient = () =>
  new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
      },
    },
  });

function renderWithProviders(initialRoute = '/') {
  const queryClient = createQueryClient();

  // Mock authenticated user
  localStorageMock.getItem.mockImplementation((key: string) => {
    if (key === 'token') return 'test-token';
    return null;
  });

  vi.mocked(authApi.me).mockResolvedValue({
    data: {
      id: 'user-1',
      username: 'TestUser',
      email: 'test@example.com',
      role: 'admin',
    },
  } as never);

  return render(
    <ThemeProvider>
      <QueryClientProvider client={queryClient}>
        <MemoryRouter initialEntries={[initialRoute]}>
          <AuthProvider>
            <App />
          </AuthProvider>
        </MemoryRouter>
      </QueryClientProvider>
    </ThemeProvider>
  );
}

describe('App', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.stubGlobal('matchMedia', matchMediaMock);
  });

  afterEach(() => {
    vi.unstubAllGlobals();
  });

  it('renders layout with AutoStrike title', async () => {
    renderWithProviders('/dashboard');
    await waitFor(() => {
      expect(screen.getByText('AutoStrike')).toBeInTheDocument();
    });
  });

  it('renders navigation items', async () => {
    renderWithProviders('/dashboard');
    await waitFor(() => {
      expect(screen.getByRole('link', { name: 'Dashboard' })).toBeInTheDocument();
    });
    expect(screen.getByRole('link', { name: 'Agents' })).toBeInTheDocument();
    expect(screen.getByRole('link', { name: 'Techniques' })).toBeInTheDocument();
    expect(screen.getByRole('link', { name: 'Scenarios' })).toBeInTheDocument();
    expect(screen.getByRole('link', { name: 'Executions' })).toBeInTheDocument();
    expect(screen.getByRole('link', { name: 'Settings' })).toBeInTheDocument();
  });

  it('redirects from / to /dashboard', async () => {
    renderWithProviders('/');
    await waitFor(() => {
      expect(screen.getByRole('link', { name: 'Dashboard' })).toBeInTheDocument();
    });
  });

  it('renders BAS Platform subtitle', async () => {
    renderWithProviders('/dashboard');
    await waitFor(() => {
      expect(screen.getByText('BAS Platform')).toBeInTheDocument();
    });
  });
});

describe('Navigation', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.stubGlobal('matchMedia', matchMediaMock);
  });

  afterEach(() => {
    vi.unstubAllGlobals();
  });

  it('highlights Dashboard link when on dashboard route', async () => {
    renderWithProviders('/dashboard');
    await waitFor(() => {
      const dashboardLink = screen.getByRole('link', { name: 'Dashboard' });
      expect(dashboardLink).toHaveClass('bg-primary-600');
    });
  });

  it('highlights Agents link when on agents route', async () => {
    renderWithProviders('/agents');
    await waitFor(() => {
      const agentsLink = screen.getByRole('link', { name: 'Agents' });
      expect(agentsLink).toHaveClass('bg-primary-600');
    });
  });

  it('highlights Settings link when on settings route', async () => {
    renderWithProviders('/settings');
    await waitFor(() => {
      const settingsLink = screen.getByRole('link', { name: 'Settings' });
      expect(settingsLink).toHaveClass('bg-primary-600');
    });
  });
});

describe('Login Page', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.stubGlobal('matchMedia', matchMediaMock);
    localStorageMock.getItem.mockReturnValue(null);
  });

  afterEach(() => {
    vi.unstubAllGlobals();
  });

  it('renders login page for unauthenticated users', async () => {
    const queryClient = createQueryClient();

    render(
      <ThemeProvider>
        <QueryClientProvider client={queryClient}>
          <MemoryRouter initialEntries={['/login']}>
            <AuthProvider>
              <App />
            </AuthProvider>
          </MemoryRouter>
        </QueryClientProvider>
      </ThemeProvider>
    );

    await waitFor(() => {
      expect(screen.getByText('Sign in to your account')).toBeInTheDocument();
    });
  });
});
