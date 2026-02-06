import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter, Routes, Route } from 'react-router-dom';
import { ProtectedRoute } from './ProtectedRoute';
import { AuthProvider } from '../contexts/AuthContext';
import { authApi } from '../lib/api';

// Mock the API
vi.mock('../lib/api', () => ({
  authApi: {
    login: vi.fn(),
    me: vi.fn(),
    logout: vi.fn(),
    refresh: vi.fn(),
  },
  healthApi: {
    check: vi.fn(),
  },
}));

import { healthApi } from '../lib/api';

// Mock localStorage
const localStorageMock = {
  getItem: vi.fn(),
  setItem: vi.fn(),
  removeItem: vi.fn(),
  clear: vi.fn(),
};
Object.defineProperty(window, 'localStorage', { value: localStorageMock });

function renderWithRouter(
  initialEntries: string[] = ['/protected'],
  hasToken = false
) {
  localStorageMock.getItem.mockImplementation((key: string) => {
    if (key === 'token' && hasToken) return 'test-token';
    return null;
  });

  return render(
    <MemoryRouter initialEntries={initialEntries}>
      <AuthProvider>
        <Routes>
          <Route path="/login" element={<div>Login Page</div>} />
          <Route
            path="/protected"
            element={
              <ProtectedRoute>
                <div>Protected Content</div>
              </ProtectedRoute>
            }
          />
        </Routes>
      </AuthProvider>
    </MemoryRouter>
  );
}

describe('ProtectedRoute', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    localStorageMock.getItem.mockReturnValue(null);
  });

  it('shows loading spinner while checking auth', () => {
    localStorageMock.getItem.mockReturnValue('test-token');
    vi.mocked(authApi.me).mockImplementation(
      () => new Promise((resolve) => setTimeout(resolve, 1000))
    );

    renderWithRouter(['/protected'], true);

    expect(screen.getByText('Loading...')).toBeInTheDocument();
  });

  it('redirects to login when not authenticated', async () => {
    localStorageMock.getItem.mockReturnValue(null);

    renderWithRouter(['/protected'], false);

    await waitFor(() => {
      expect(screen.getByText('Login Page')).toBeInTheDocument();
    });
  });

  it('renders protected content when authenticated', async () => {
    localStorageMock.getItem.mockReturnValue('test-token');
    vi.mocked(authApi.me).mockResolvedValue({
      data: {
        id: 'user-1',
        username: 'testuser',
        email: 'test@example.com',
        role: 'admin',
      },
    } as never);

    renderWithRouter(['/protected'], true);

    await waitFor(() => {
      expect(screen.getByText('Protected Content')).toBeInTheDocument();
    });
  });

  it('redirects to login when token is invalid', async () => {
    localStorageMock.getItem.mockReturnValue('invalid-token');
    vi.mocked(authApi.me).mockRejectedValue(new Error('Unauthorized'));

    renderWithRouter(['/protected'], true);

    await waitFor(() => {
      expect(screen.getByText('Login Page')).toBeInTheDocument();
    });
  });
});

describe('ProtectedRoute with children', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders children correctly when authenticated', async () => {
    localStorageMock.getItem.mockReturnValue('test-token');
    vi.mocked(authApi.me).mockResolvedValue({
      data: {
        id: 'user-1',
        username: 'testuser',
        email: 'test@example.com',
        role: 'admin',
      },
    } as never);

    render(
      <MemoryRouter initialEntries={['/protected']}>
        <AuthProvider>
          <Routes>
            <Route
              path="/protected"
              element={
                <ProtectedRoute>
                  <div data-testid="custom-child">Custom Child Content</div>
                </ProtectedRoute>
              }
            />
          </Routes>
        </AuthProvider>
      </MemoryRouter>
    );

    await waitFor(() => {
      expect(screen.getByTestId('custom-child')).toBeInTheDocument();
      expect(screen.getByText('Custom Child Content')).toBeInTheDocument();
    });
  });
});

describe('ProtectedRoute with requiredRole', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('allows access when user has exact required role', async () => {
    localStorageMock.getItem.mockReturnValue('test-token');
    vi.mocked(authApi.me).mockResolvedValue({
      data: {
        id: 'user-1',
        username: 'testuser',
        email: 'test@example.com',
        role: 'operator',
      },
    } as never);

    render(
      <MemoryRouter initialEntries={['/admin']}>
        <AuthProvider>
          <Routes>
            <Route
              path="/admin"
              element={
                <ProtectedRoute requiredRole="operator">
                  <div>Operator Content</div>
                </ProtectedRoute>
              }
            />
          </Routes>
        </AuthProvider>
      </MemoryRouter>
    );

    await waitFor(() => {
      expect(screen.getByText('Operator Content')).toBeInTheDocument();
    });
  });

  it('allows access when user has higher role than required', async () => {
    localStorageMock.getItem.mockReturnValue('test-token');
    vi.mocked(authApi.me).mockResolvedValue({
      data: {
        id: 'user-1',
        username: 'testuser',
        email: 'test@example.com',
        role: 'admin',
      },
    } as never);

    render(
      <MemoryRouter initialEntries={['/content']}>
        <AuthProvider>
          <Routes>
            <Route
              path="/content"
              element={
                <ProtectedRoute requiredRole="viewer">
                  <div>Protected Content</div>
                </ProtectedRoute>
              }
            />
          </Routes>
        </AuthProvider>
      </MemoryRouter>
    );

    await waitFor(() => {
      expect(screen.getByText('Protected Content')).toBeInTheDocument();
    });
  });

  it('denies access when user has lower role than required', async () => {
    localStorageMock.getItem.mockReturnValue('test-token');
    vi.mocked(authApi.me).mockResolvedValue({
      data: {
        id: 'user-1',
        username: 'testuser',
        email: 'test@example.com',
        role: 'viewer',
      },
    } as never);

    render(
      <MemoryRouter initialEntries={['/admin']}>
        <AuthProvider>
          <Routes>
            <Route
              path="/admin"
              element={
                <ProtectedRoute requiredRole="admin">
                  <div>Admin Only Content</div>
                </ProtectedRoute>
              }
            />
          </Routes>
        </AuthProvider>
      </MemoryRouter>
    );

    await waitFor(() => {
      expect(screen.getByText('403')).toBeInTheDocument();
      expect(screen.getByText('Access Denied')).toBeInTheDocument();
    });
  });

  it('shows return to dashboard link on access denied', async () => {
    localStorageMock.getItem.mockReturnValue('test-token');
    vi.mocked(authApi.me).mockResolvedValue({
      data: {
        id: 'user-1',
        username: 'testuser',
        email: 'test@example.com',
        role: 'viewer',
      },
    } as never);

    render(
      <MemoryRouter initialEntries={['/admin']}>
        <AuthProvider>
          <Routes>
            <Route
              path="/admin"
              element={
                <ProtectedRoute requiredRole="admin">
                  <div>Admin Only Content</div>
                </ProtectedRoute>
              }
            />
          </Routes>
        </AuthProvider>
      </MemoryRouter>
    );

    await waitFor(() => {
      expect(screen.getByText('Return to Dashboard')).toBeInTheDocument();
    });
  });
});

describe('ProtectedRoute with allowedRoles', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('allows access when user role is in allowedRoles list', async () => {
    localStorageMock.getItem.mockReturnValue('test-token');
    vi.mocked(authApi.me).mockResolvedValue({
      data: {
        id: 'user-1',
        username: 'testuser',
        email: 'test@example.com',
        role: 'analyst',
      },
    } as never);

    render(
      <MemoryRouter initialEntries={['/reports']}>
        <AuthProvider>
          <Routes>
            <Route
              path="/reports"
              element={
                <ProtectedRoute allowedRoles={['admin', 'analyst', 'rssi']}>
                  <div>Reports Content</div>
                </ProtectedRoute>
              }
            />
          </Routes>
        </AuthProvider>
      </MemoryRouter>
    );

    await waitFor(() => {
      expect(screen.getByText('Reports Content')).toBeInTheDocument();
    });
  });

  it('denies access when user role is not in allowedRoles list', async () => {
    localStorageMock.getItem.mockReturnValue('test-token');
    vi.mocked(authApi.me).mockResolvedValue({
      data: {
        id: 'user-1',
        username: 'testuser',
        email: 'test@example.com',
        role: 'viewer',
      },
    } as never);

    render(
      <MemoryRouter initialEntries={['/reports']}>
        <AuthProvider>
          <Routes>
            <Route
              path="/reports"
              element={
                <ProtectedRoute allowedRoles={['admin', 'analyst']}>
                  <div>Reports Content</div>
                </ProtectedRoute>
              }
            />
          </Routes>
        </AuthProvider>
      </MemoryRouter>
    );

    await waitFor(() => {
      expect(screen.getByText('403')).toBeInTheDocument();
      expect(screen.getByText('Access Denied')).toBeInTheDocument();
    });
  });

  it('shows permission denied message', async () => {
    localStorageMock.getItem.mockReturnValue('test-token');
    vi.mocked(authApi.me).mockResolvedValue({
      data: {
        id: 'user-1',
        username: 'testuser',
        email: 'test@example.com',
        role: 'viewer',
      },
    } as never);

    render(
      <MemoryRouter initialEntries={['/admin']}>
        <AuthProvider>
          <Routes>
            <Route
              path="/admin"
              element={
                <ProtectedRoute allowedRoles={['admin']}>
                  <div>Admin Content</div>
                </ProtectedRoute>
              }
            />
          </Routes>
        </AuthProvider>
      </MemoryRouter>
    );

    await waitFor(() => {
      expect(
        screen.getByText("You don't have permission to access this page.")
      ).toBeInTheDocument();
    });
  });
});

describe('ProtectedRoute role hierarchy', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('admin has access to all role levels', async () => {
    localStorageMock.getItem.mockReturnValue('test-token');
    vi.mocked(authApi.me).mockResolvedValue({
      data: {
        id: 'user-1',
        username: 'admin',
        email: 'admin@example.com',
        role: 'admin',
      },
    } as never);

    render(
      <MemoryRouter initialEntries={['/viewer-page']}>
        <AuthProvider>
          <Routes>
            <Route
              path="/viewer-page"
              element={
                <ProtectedRoute requiredRole="viewer">
                  <div>Viewer Content</div>
                </ProtectedRoute>
              }
            />
          </Routes>
        </AuthProvider>
      </MemoryRouter>
    );

    await waitFor(() => {
      expect(screen.getByText('Viewer Content')).toBeInTheDocument();
    });
  });

  it('rssi has access to analyst-level content', async () => {
    localStorageMock.getItem.mockReturnValue('test-token');
    vi.mocked(authApi.me).mockResolvedValue({
      data: {
        id: 'user-1',
        username: 'rssi',
        email: 'rssi@example.com',
        role: 'rssi',
      },
    } as never);

    render(
      <MemoryRouter initialEntries={['/analyst-page']}>
        <AuthProvider>
          <Routes>
            <Route
              path="/analyst-page"
              element={
                <ProtectedRoute requiredRole="analyst">
                  <div>Analyst Content</div>
                </ProtectedRoute>
              }
            />
          </Routes>
        </AuthProvider>
      </MemoryRouter>
    );

    await waitFor(() => {
      expect(screen.getByText('Analyst Content')).toBeInTheDocument();
    });
  });

  it('operator cannot access rssi-level content', async () => {
    localStorageMock.getItem.mockReturnValue('test-token');
    vi.mocked(authApi.me).mockResolvedValue({
      data: {
        id: 'user-1',
        username: 'operator',
        email: 'operator@example.com',
        role: 'operator',
      },
    } as never);

    render(
      <MemoryRouter initialEntries={['/rssi-page']}>
        <AuthProvider>
          <Routes>
            <Route
              path="/rssi-page"
              element={
                <ProtectedRoute requiredRole="rssi">
                  <div>RSSI Content</div>
                </ProtectedRoute>
              }
            />
          </Routes>
        </AuthProvider>
      </MemoryRouter>
    );

    await waitFor(() => {
      expect(screen.getByText('403')).toBeInTheDocument();
      expect(screen.getByText('Access Denied')).toBeInTheDocument();
    });
  });
});

describe('ProtectedRoute without role requirements', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('allows any authenticated user when no role specified', async () => {
    localStorageMock.getItem.mockReturnValue('test-token');
    vi.mocked(authApi.me).mockResolvedValue({
      data: {
        id: 'user-1',
        username: 'viewer',
        email: 'viewer@example.com',
        role: 'viewer',
      },
    } as never);

    render(
      <MemoryRouter initialEntries={['/general']}>
        <AuthProvider>
          <Routes>
            <Route
              path="/general"
              element={
                <ProtectedRoute>
                  <div>General Content</div>
                </ProtectedRoute>
              }
            />
          </Routes>
        </AuthProvider>
      </MemoryRouter>
    );

    await waitFor(() => {
      expect(screen.getByText('General Content')).toBeInTheDocument();
    });
  });
});

describe('ProtectedRoute with auth disabled', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders children when auth is disabled even without a user', async () => {
    // No token in localStorage
    localStorageMock.getItem.mockReturnValue(null);
    // Health check returns auth_enabled: false
    vi.mocked(healthApi.check).mockResolvedValue({
      data: { status: 'ok', auth_enabled: false },
    } as never);

    render(
      <MemoryRouter initialEntries={['/protected']}>
        <AuthProvider>
          <Routes>
            <Route path="/login" element={<div>Login Page</div>} />
            <Route
              path="/protected"
              element={
                <ProtectedRoute>
                  <div>Protected Content</div>
                </ProtectedRoute>
              }
            />
          </Routes>
        </AuthProvider>
      </MemoryRouter>
    );

    await waitFor(() => {
      expect(screen.getByText('Protected Content')).toBeInTheDocument();
    });
  });

  it('skips role checks when auth is disabled', async () => {
    // No token in localStorage
    localStorageMock.getItem.mockReturnValue(null);
    // Health check returns auth_enabled: false
    vi.mocked(healthApi.check).mockResolvedValue({
      data: { status: 'ok', auth_enabled: false },
    } as never);

    render(
      <MemoryRouter initialEntries={['/admin']}>
        <AuthProvider>
          <Routes>
            <Route path="/login" element={<div>Login Page</div>} />
            <Route
              path="/admin"
              element={
                <ProtectedRoute requiredRole="admin">
                  <div>Admin Content</div>
                </ProtectedRoute>
              }
            />
          </Routes>
        </AuthProvider>
      </MemoryRouter>
    );

    await waitFor(() => {
      expect(screen.getByText('Admin Content')).toBeInTheDocument();
    });
  });
});

describe('ProtectedRoute with both requiredRole and allowedRoles', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    // Ensure auth is enabled so role checks are exercised
    vi.mocked(healthApi.check).mockResolvedValue({
      data: { status: 'ok', auth_enabled: true },
    } as never);
  });

  it('allowedRoles takes precedence over requiredRole', async () => {
    localStorageMock.getItem.mockReturnValue('test-token');
    vi.mocked(authApi.me).mockResolvedValue({
      data: {
        id: 'user-1',
        username: 'analyst',
        email: 'analyst@example.com',
        role: 'analyst',
      },
    } as never);

    render(
      <MemoryRouter initialEntries={['/special']}>
        <AuthProvider>
          <Routes>
            <Route
              path="/special"
              element={
                <ProtectedRoute requiredRole="admin" allowedRoles={['analyst']}>
                  <div>Special Content</div>
                </ProtectedRoute>
              }
            />
          </Routes>
        </AuthProvider>
      </MemoryRouter>
    );

    // analyst would be denied by requiredRole='admin' (hierarchy: analyst=2 < admin=5)
    // but allowedRoles=['analyst'] takes precedence and grants access
    await waitFor(() => {
      expect(screen.getByText('Special Content')).toBeInTheDocument();
    });
  });

  it('denies access when user role is not in allowedRoles even if requiredRole would allow', async () => {
    localStorageMock.getItem.mockReturnValue('test-token');
    vi.mocked(authApi.me).mockResolvedValue({
      data: {
        id: 'user-1',
        username: 'admin',
        email: 'admin@example.com',
        role: 'admin',
      },
    } as never);

    render(
      <MemoryRouter initialEntries={['/special']}>
        <AuthProvider>
          <Routes>
            <Route
              path="/special"
              element={
                <ProtectedRoute requiredRole="viewer" allowedRoles={['analyst']}>
                  <div>Special Content</div>
                </ProtectedRoute>
              }
            />
          </Routes>
        </AuthProvider>
      </MemoryRouter>
    );

    // admin would be allowed by requiredRole='viewer' (hierarchy: admin=5 > viewer=1)
    // but allowedRoles=['analyst'] takes precedence and denies access since admin is not in the list
    await waitFor(() => {
      expect(screen.getByText('403')).toBeInTheDocument();
      expect(screen.getByText('Access Denied')).toBeInTheDocument();
    });
  });
});

describe('ProtectedRoute loading state details', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('shows spinner with correct structure during loading', () => {
    localStorageMock.getItem.mockReturnValue('test-token');
    vi.mocked(authApi.me).mockImplementation(
      () => new Promise((resolve) => setTimeout(resolve, 5000))
    );

    const { container } = renderWithRouter(['/protected'], true);

    // Verify the loading spinner element exists with the animate-spin class
    const spinner = container.querySelector('.animate-spin');
    expect(spinner).toBeInTheDocument();

    // Verify the loading text
    expect(screen.getByText('Loading...')).toBeInTheDocument();

    // Verify the container has the expected background classes
    const wrapper = container.querySelector('.min-h-screen');
    expect(wrapper).toBeInTheDocument();
  });

  it('does not show protected content while loading', () => {
    localStorageMock.getItem.mockReturnValue('test-token');
    vi.mocked(authApi.me).mockImplementation(
      () => new Promise((resolve) => setTimeout(resolve, 5000))
    );

    renderWithRouter(['/protected'], true);

    expect(screen.queryByText('Protected Content')).not.toBeInTheDocument();
    expect(screen.getByText('Loading...')).toBeInTheDocument();
  });
});

describe('ProtectedRoute redirect preserves location', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('redirects to login when not authenticated and preserves route', async () => {
    localStorageMock.getItem.mockReturnValue(null);

    renderWithRouter(['/protected'], false);

    await waitFor(() => {
      expect(screen.getByText('Login Page')).toBeInTheDocument();
    });

    // Verify we navigated away from protected content
    expect(screen.queryByText('Protected Content')).not.toBeInTheDocument();
  });
});

describe('ProtectedRoute 403 page details', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('shows return to dashboard link with correct href', async () => {
    localStorageMock.getItem.mockReturnValue('test-token');
    vi.mocked(authApi.me).mockResolvedValue({
      data: {
        id: 'user-1',
        username: 'viewer',
        email: 'viewer@example.com',
        role: 'viewer',
      },
    } as never);

    render(
      <MemoryRouter initialEntries={['/admin']}>
        <AuthProvider>
          <Routes>
            <Route
              path="/admin"
              element={
                <ProtectedRoute requiredRole="admin">
                  <div>Admin Content</div>
                </ProtectedRoute>
              }
            />
          </Routes>
        </AuthProvider>
      </MemoryRouter>
    );

    await waitFor(() => {
      const link = screen.getByText('Return to Dashboard');
      expect(link).toBeInTheDocument();
      expect(link.getAttribute('href')).toBe('/dashboard');
    });
  });

  it('shows all 403 page elements together', async () => {
    localStorageMock.getItem.mockReturnValue('test-token');
    vi.mocked(authApi.me).mockResolvedValue({
      data: {
        id: 'user-1',
        username: 'viewer',
        email: 'viewer@example.com',
        role: 'viewer',
      },
    } as never);

    render(
      <MemoryRouter initialEntries={['/admin']}>
        <AuthProvider>
          <Routes>
            <Route
              path="/admin"
              element={
                <ProtectedRoute requiredRole="admin">
                  <div>Admin Content</div>
                </ProtectedRoute>
              }
            />
          </Routes>
        </AuthProvider>
      </MemoryRouter>
    );

    await waitFor(() => {
      expect(screen.getByText('403')).toBeInTheDocument();
      expect(screen.getByText('Access Denied')).toBeInTheDocument();
      expect(screen.getByText("You don't have permission to access this page.")).toBeInTheDocument();
      expect(screen.getByText('Return to Dashboard')).toBeInTheDocument();
    });

    // Verify admin content is NOT rendered
    expect(screen.queryByText('Admin Content')).not.toBeInTheDocument();
  });
});

describe('ProtectedRoute with empty allowedRoles', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('falls through to hierarchy check when allowedRoles is empty array', async () => {
    localStorageMock.getItem.mockReturnValue('test-token');
    vi.mocked(authApi.me).mockResolvedValue({
      data: {
        id: 'user-1',
        username: 'admin',
        email: 'admin@example.com',
        role: 'admin',
      },
    } as never);

    render(
      <MemoryRouter initialEntries={['/page']}>
        <AuthProvider>
          <Routes>
            <Route
              path="/page"
              element={
                <ProtectedRoute requiredRole="operator" allowedRoles={[]}>
                  <div>Fallthrough Content</div>
                </ProtectedRoute>
              }
            />
          </Routes>
        </AuthProvider>
      </MemoryRouter>
    );

    // Empty allowedRoles should not block; falls through to requiredRole hierarchy check
    // admin (5) >= operator (3), so access is granted
    await waitFor(() => {
      expect(screen.getByText('Fallthrough Content')).toBeInTheDocument();
    });
  });

  it('denies via hierarchy when allowedRoles is empty and user role is insufficient', async () => {
    localStorageMock.getItem.mockReturnValue('test-token');
    vi.mocked(authApi.me).mockResolvedValue({
      data: {
        id: 'user-1',
        username: 'viewer',
        email: 'viewer@example.com',
        role: 'viewer',
      },
    } as never);

    render(
      <MemoryRouter initialEntries={['/page']}>
        <AuthProvider>
          <Routes>
            <Route
              path="/page"
              element={
                <ProtectedRoute requiredRole="operator" allowedRoles={[]}>
                  <div>Fallthrough Content</div>
                </ProtectedRoute>
              }
            />
          </Routes>
        </AuthProvider>
      </MemoryRouter>
    );

    // Empty allowedRoles falls through to hierarchy; viewer (1) < operator (3), so denied
    await waitFor(() => {
      expect(screen.getByText('403')).toBeInTheDocument();
      expect(screen.getByText('Access Denied')).toBeInTheDocument();
    });
  });
});

describe('ProtectedRoute auth disabled with allowedRoles', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('skips allowedRoles check when auth is disabled', async () => {
    localStorageMock.getItem.mockReturnValue(null);
    vi.mocked(healthApi.check).mockResolvedValue({
      data: { status: 'ok', auth_enabled: false },
    } as never);

    render(
      <MemoryRouter initialEntries={['/reports']}>
        <AuthProvider>
          <Routes>
            <Route path="/login" element={<div>Login Page</div>} />
            <Route
              path="/reports"
              element={
                <ProtectedRoute allowedRoles={['analyst']}>
                  <div>Reports Content</div>
                </ProtectedRoute>
              }
            />
          </Routes>
        </AuthProvider>
      </MemoryRouter>
    );

    await waitFor(() => {
      expect(screen.getByText('Reports Content')).toBeInTheDocument();
    });
  });

  it('skips combined requiredRole and allowedRoles check when auth is disabled', async () => {
    localStorageMock.getItem.mockReturnValue(null);
    vi.mocked(healthApi.check).mockResolvedValue({
      data: { status: 'ok', auth_enabled: false },
    } as never);

    render(
      <MemoryRouter initialEntries={['/special']}>
        <AuthProvider>
          <Routes>
            <Route path="/login" element={<div>Login Page</div>} />
            <Route
              path="/special"
              element={
                <ProtectedRoute requiredRole="admin" allowedRoles={['rssi']}>
                  <div>Special Content</div>
                </ProtectedRoute>
              }
            />
          </Routes>
        </AuthProvider>
      </MemoryRouter>
    );

    await waitFor(() => {
      expect(screen.getByText('Special Content')).toBeInTheDocument();
    });
  });
});

describe('ProtectedRoute complete role hierarchy coverage', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    // Ensure auth is enabled so role checks are exercised
    vi.mocked(healthApi.check).mockResolvedValue({
      data: { status: 'ok', auth_enabled: true },
    } as never);
  });

  it('viewer has access to viewer-level content (exact match)', async () => {
    localStorageMock.getItem.mockReturnValue('test-token');
    vi.mocked(authApi.me).mockResolvedValue({
      data: {
        id: 'user-1',
        username: 'viewer',
        email: 'viewer@example.com',
        role: 'viewer',
      },
    } as never);

    render(
      <MemoryRouter initialEntries={['/view']}>
        <AuthProvider>
          <Routes>
            <Route
              path="/view"
              element={
                <ProtectedRoute requiredRole="viewer">
                  <div>Viewer Content</div>
                </ProtectedRoute>
              }
            />
          </Routes>
        </AuthProvider>
      </MemoryRouter>
    );

    await waitFor(() => {
      expect(screen.getByText('Viewer Content')).toBeInTheDocument();
    });
  });

  it('analyst has access to analyst-level content (exact match)', async () => {
    localStorageMock.getItem.mockReturnValue('test-token');
    vi.mocked(authApi.me).mockResolvedValue({
      data: {
        id: 'user-1',
        username: 'analyst',
        email: 'analyst@example.com',
        role: 'analyst',
      },
    } as never);

    render(
      <MemoryRouter initialEntries={['/analysis']}>
        <AuthProvider>
          <Routes>
            <Route
              path="/analysis"
              element={
                <ProtectedRoute requiredRole="analyst">
                  <div>Analyst Content</div>
                </ProtectedRoute>
              }
            />
          </Routes>
        </AuthProvider>
      </MemoryRouter>
    );

    await waitFor(() => {
      expect(screen.getByText('Analyst Content')).toBeInTheDocument();
    });
  });

  it('analyst cannot access operator-level content', async () => {
    localStorageMock.getItem.mockReturnValue('test-token');
    vi.mocked(authApi.me).mockResolvedValue({
      data: {
        id: 'user-1',
        username: 'analyst',
        email: 'analyst@example.com',
        role: 'analyst',
      },
    } as never);

    render(
      <MemoryRouter initialEntries={['/ops']}>
        <AuthProvider>
          <Routes>
            <Route
              path="/ops"
              element={
                <ProtectedRoute requiredRole="operator">
                  <div>Operator Content</div>
                </ProtectedRoute>
              }
            />
          </Routes>
        </AuthProvider>
      </MemoryRouter>
    );

    await waitFor(() => {
      expect(screen.getByText('403')).toBeInTheDocument();
    });
  });

  it('viewer cannot access analyst-level content', async () => {
    localStorageMock.getItem.mockReturnValue('test-token');
    vi.mocked(authApi.me).mockResolvedValue({
      data: {
        id: 'user-1',
        username: 'viewer',
        email: 'viewer@example.com',
        role: 'viewer',
      },
    } as never);

    render(
      <MemoryRouter initialEntries={['/analysis']}>
        <AuthProvider>
          <Routes>
            <Route
              path="/analysis"
              element={
                <ProtectedRoute requiredRole="analyst">
                  <div>Analyst Content</div>
                </ProtectedRoute>
              }
            />
          </Routes>
        </AuthProvider>
      </MemoryRouter>
    );

    await waitFor(() => {
      expect(screen.getByText('403')).toBeInTheDocument();
    });
  });

  it('rssi has access to operator-level content via hierarchy', async () => {
    localStorageMock.getItem.mockReturnValue('test-token');
    vi.mocked(authApi.me).mockResolvedValue({
      data: {
        id: 'user-1',
        username: 'rssi',
        email: 'rssi@example.com',
        role: 'rssi',
      },
    } as never);

    render(
      <MemoryRouter initialEntries={['/ops']}>
        <AuthProvider>
          <Routes>
            <Route
              path="/ops"
              element={
                <ProtectedRoute requiredRole="operator">
                  <div>Operator Content</div>
                </ProtectedRoute>
              }
            />
          </Routes>
        </AuthProvider>
      </MemoryRouter>
    );

    await waitFor(() => {
      expect(screen.getByText('Operator Content')).toBeInTheDocument();
    });
  });

  it('rssi cannot access admin-level content', async () => {
    localStorageMock.getItem.mockReturnValue('test-token');
    vi.mocked(authApi.me).mockResolvedValue({
      data: {
        id: 'user-1',
        username: 'rssi',
        email: 'rssi@example.com',
        role: 'rssi',
      },
    } as never);

    render(
      <MemoryRouter initialEntries={['/admin']}>
        <AuthProvider>
          <Routes>
            <Route
              path="/admin"
              element={
                <ProtectedRoute requiredRole="admin">
                  <div>Admin Content</div>
                </ProtectedRoute>
              }
            />
          </Routes>
        </AuthProvider>
      </MemoryRouter>
    );

    await waitFor(() => {
      expect(screen.getByText('403')).toBeInTheDocument();
    });
  });

  it('admin has access to admin-level content (exact match at top)', async () => {
    localStorageMock.getItem.mockReturnValue('test-token');
    vi.mocked(authApi.me).mockResolvedValue({
      data: {
        id: 'user-1',
        username: 'admin',
        email: 'admin@example.com',
        role: 'admin',
      },
    } as never);

    render(
      <MemoryRouter initialEntries={['/admin']}>
        <AuthProvider>
          <Routes>
            <Route
              path="/admin"
              element={
                <ProtectedRoute requiredRole="admin">
                  <div>Admin Content</div>
                </ProtectedRoute>
              }
            />
          </Routes>
        </AuthProvider>
      </MemoryRouter>
    );

    await waitFor(() => {
      expect(screen.getByText('Admin Content')).toBeInTheDocument();
    });
  });
});

describe('ProtectedRoute allowedRoles edge cases', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    // Ensure auth is enabled so role checks are exercised
    vi.mocked(healthApi.check).mockResolvedValue({
      data: { status: 'ok', auth_enabled: true },
    } as never);
  });

  it('allows access when user is one of multiple allowed roles', async () => {
    localStorageMock.getItem.mockReturnValue('test-token');
    vi.mocked(authApi.me).mockResolvedValue({
      data: {
        id: 'user-1',
        username: 'operator',
        email: 'operator@example.com',
        role: 'operator',
      },
    } as never);

    render(
      <MemoryRouter initialEntries={['/multi']}>
        <AuthProvider>
          <Routes>
            <Route
              path="/multi"
              element={
                <ProtectedRoute allowedRoles={['operator', 'analyst', 'rssi']}>
                  <div>Multi Role Content</div>
                </ProtectedRoute>
              }
            />
          </Routes>
        </AuthProvider>
      </MemoryRouter>
    );

    await waitFor(() => {
      expect(screen.getByText('Multi Role Content')).toBeInTheDocument();
    });
  });

  it('denies access with single-item allowedRoles when role does not match', async () => {
    localStorageMock.getItem.mockReturnValue('test-token');
    vi.mocked(authApi.me).mockResolvedValue({
      data: {
        id: 'user-1',
        username: 'operator',
        email: 'operator@example.com',
        role: 'operator',
      },
    } as never);

    render(
      <MemoryRouter initialEntries={['/single']}>
        <AuthProvider>
          <Routes>
            <Route
              path="/single"
              element={
                <ProtectedRoute allowedRoles={['rssi']}>
                  <div>RSSI Only Content</div>
                </ProtectedRoute>
              }
            />
          </Routes>
        </AuthProvider>
      </MemoryRouter>
    );

    await waitFor(() => {
      expect(screen.getByText('403')).toBeInTheDocument();
    });
  });

  it('allows rssi when rssi is in allowedRoles with admin', async () => {
    localStorageMock.getItem.mockReturnValue('test-token');
    vi.mocked(authApi.me).mockResolvedValue({
      data: {
        id: 'user-1',
        username: 'rssi',
        email: 'rssi@example.com',
        role: 'rssi',
      },
    } as never);

    render(
      <MemoryRouter initialEntries={['/management']}>
        <AuthProvider>
          <Routes>
            <Route
              path="/management"
              element={
                <ProtectedRoute allowedRoles={['admin', 'rssi']}>
                  <div>Management Content</div>
                </ProtectedRoute>
              }
            />
          </Routes>
        </AuthProvider>
      </MemoryRouter>
    );

    await waitFor(() => {
      expect(screen.getByText('Management Content')).toBeInTheDocument();
    });
  });

  it('denies viewer when only admin and rssi are in allowedRoles', async () => {
    localStorageMock.getItem.mockReturnValue('test-token');
    vi.mocked(authApi.me).mockResolvedValue({
      data: {
        id: 'user-1',
        username: 'viewer',
        email: 'viewer@example.com',
        role: 'viewer',
      },
    } as never);

    render(
      <MemoryRouter initialEntries={['/management']}>
        <AuthProvider>
          <Routes>
            <Route
              path="/management"
              element={
                <ProtectedRoute allowedRoles={['admin', 'rssi']}>
                  <div>Management Content</div>
                </ProtectedRoute>
              }
            />
          </Routes>
        </AuthProvider>
      </MemoryRouter>
    );

    await waitFor(() => {
      expect(screen.getByText('403')).toBeInTheDocument();
      expect(screen.getByText('Access Denied')).toBeInTheDocument();
    });
  });
});
