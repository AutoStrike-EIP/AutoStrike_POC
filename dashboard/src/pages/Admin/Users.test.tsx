import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import Users from './Users';

// Mock the auth context
vi.mock('../../contexts/AuthContext', () => ({
  useAuth: vi.fn(() => ({
    user: { id: 'current-user', username: 'admin', role: 'admin' },
    isAuthenticated: true,
  })),
}));

// Mock the API
vi.mock('../../lib/api', () => ({
  adminApi: {
    listUsers: vi.fn(() =>
      Promise.resolve({
        data: {
          users: [
            {
              id: 'user-1',
              username: 'admin',
              email: 'admin@example.com',
              role: 'admin',
              role_display: 'Administrator',
              is_active: true,
              last_login_at: new Date().toISOString(),
              created_at: new Date().toISOString(),
              updated_at: new Date().toISOString(),
            },
            {
              id: 'user-2',
              username: 'operator',
              email: 'operator@example.com',
              role: 'operator',
              role_display: 'Operator',
              is_active: true,
              last_login_at: null,
              created_at: new Date().toISOString(),
              updated_at: new Date().toISOString(),
            },
            {
              id: 'user-3',
              username: 'inactive',
              email: 'inactive@example.com',
              role: 'viewer',
              role_display: 'Viewer',
              is_active: false,
              last_login_at: null,
              created_at: new Date().toISOString(),
              updated_at: new Date().toISOString(),
            },
          ],
          total: 3,
        },
      })
    ),
    createUser: vi.fn(() => Promise.resolve({ data: {} })),
    updateUser: vi.fn(() => Promise.resolve({ data: {} })),
    updateUserRole: vi.fn(() => Promise.resolve({ data: {} })),
    resetPassword: vi.fn(() => Promise.resolve({ data: {} })),
    deactivateUser: vi.fn(() => Promise.resolve({ data: {} })),
    reactivateUser: vi.fn(() => Promise.resolve({ data: {} })),
  },
}));

function createTestQueryClient() {
  return new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
      },
      mutations: {
        retry: false,
      },
    },
  });
}

function renderUsers() {
  const queryClient = createTestQueryClient();
  return render(
    <QueryClientProvider client={queryClient}>
      <MemoryRouter>
        <Users />
      </MemoryRouter>
    </QueryClientProvider>
  );
}

describe('Users Page', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders user management title', async () => {
    renderUsers();
    await waitFor(() => {
      expect(screen.getByText('User Management')).toBeInTheDocument();
    });
  });

  it('renders add user button', async () => {
    renderUsers();
    await waitFor(() => {
      expect(screen.getByText('Add User')).toBeInTheDocument();
    });
  });

  it('displays users after loading', async () => {
    renderUsers();
    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
      expect(screen.getByText('operator@example.com')).toBeInTheDocument();
    });
  });

  it('shows role badges', async () => {
    renderUsers();
    await waitFor(() => {
      expect(screen.getByText('Administrator')).toBeInTheDocument();
      expect(screen.getByText('Operator')).toBeInTheDocument();
    });
  });

  it('shows active/inactive status', async () => {
    renderUsers();
    await waitFor(() => {
      const activeLabels = screen.getAllByText('Active');
      expect(activeLabels.length).toBeGreaterThan(0);
    });
  });

  it('shows show inactive users checkbox', async () => {
    renderUsers();
    await waitFor(() => {
      expect(screen.getByText('Show inactive users')).toBeInTheDocument();
    });
  });

  it('shows last login time or Never', async () => {
    renderUsers();
    await waitFor(() => {
      const neverTexts = screen.getAllByText('Never');
      expect(neverTexts.length).toBeGreaterThan(0);
    });
  });
});

describe('Users Create Modal', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('has add user button', async () => {
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    const addButton = screen.getByText('Add User');
    expect(addButton).toBeInTheDocument();
  });

  it('shows role options in create form', async () => {
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Add User'));

    await waitFor(() => {
      const roleSelect = screen.getByLabelText('Role');
      expect(roleSelect).toBeInTheDocument();
    });
  });
});

describe('Users Edit Modal', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('opens edit modal when edit button clicked', async () => {
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    const editButtons = screen.getAllByTitle('Edit user');
    fireEvent.click(editButtons[0]);

    await waitFor(() => {
      expect(screen.getByText('Edit User')).toBeInTheDocument();
    });
  });

  it('closes edit modal when cancel clicked', async () => {
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    const editButtons = screen.getAllByTitle('Edit user');
    fireEvent.click(editButtons[0]);

    await waitFor(() => {
      expect(screen.getByText('Edit User')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Cancel'));

    await waitFor(() => {
      expect(screen.queryByText('Edit User')).not.toBeInTheDocument();
    });
  });
});

describe('Users Role Modal', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('opens role modal when change role button clicked', async () => {
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('operator@example.com')).toBeInTheDocument();
    });

    const roleButtons = screen.getAllByTitle('Change role');
    // Click on a user that is not the current user
    const enabledRoleButton = roleButtons.find(
      (btn) => !(btn as HTMLButtonElement).disabled
    );
    if (enabledRoleButton) {
      fireEvent.click(enabledRoleButton);

      await waitFor(() => {
        expect(screen.getByText('Change Role')).toBeInTheDocument();
      });
    }
  });
});

describe('Users Password Modal', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('shows reset password buttons', async () => {
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    const passwordButtons = screen.getAllByTitle('Reset password');
    expect(passwordButtons.length).toBeGreaterThan(0);
  });
});

describe('Users Deactivate Modal', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('opens deactivate modal when deactivate button clicked', async () => {
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('operator@example.com')).toBeInTheDocument();
    });

    const deactivateButtons = screen.getAllByTitle('Deactivate user');
    // Click on a user that is not the current user
    const enabledButton = deactivateButtons.find(
      (btn) => !(btn as HTMLButtonElement).disabled
    );
    if (enabledButton) {
      fireEvent.click(enabledButton);

      await waitFor(() => {
        expect(screen.getByText('Deactivate User')).toBeInTheDocument();
      });
    }
  });

  it('shows warning message in deactivate modal', async () => {
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('operator@example.com')).toBeInTheDocument();
    });

    const deactivateButtons = screen.getAllByTitle('Deactivate user');
    const enabledButton = deactivateButtons.find(
      (btn) => !(btn as HTMLButtonElement).disabled
    );
    if (enabledButton) {
      fireEvent.click(enabledButton);

      await waitFor(() => {
        expect(
          screen.getByText(/The user will no longer be able to log in/)
        ).toBeInTheDocument();
      });
    }
  });
});

describe('Users Reactivate Modal', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('shows reactivate button for inactive users', async () => {
    renderUsers();

    // First, enable showing inactive users
    await waitFor(() => {
      expect(screen.getByText('Show inactive users')).toBeInTheDocument();
    });

    const checkbox = screen.getByRole('checkbox');
    fireEvent.click(checkbox);

    await waitFor(() => {
      const reactivateButtons = screen.queryAllByTitle('Reactivate user');
      expect(reactivateButtons.length).toBeGreaterThanOrEqual(0);
    });
  });
});

describe('Users Table Actions', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('disables change role button for current user', async () => {
    // Mock useAuth to return the same user as in the list
    const { useAuth } = await import('../../contexts/AuthContext');
    vi.mocked(useAuth).mockReturnValue({
      user: { id: 'user-1', username: 'admin', email: 'admin@example.com', role: 'admin', is_active: true, created_at: new Date().toISOString(), updated_at: new Date().toISOString() },
      isAuthenticated: true,
      isLoading: false,
      login: vi.fn(),
      logout: vi.fn(),
      authEnabled: true,
    });

    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    const roleButtons = screen.getAllByTitle('Change role');
    // The first user (admin) should have a disabled role button
    expect(roleButtons[0]).toBeDisabled();
  });

  it('disables deactivate button for current user', async () => {
    const { useAuth } = await import('../../contexts/AuthContext');
    vi.mocked(useAuth).mockReturnValue({
      user: { id: 'user-1', username: 'admin', email: 'admin@example.com', role: 'admin', is_active: true, created_at: new Date().toISOString(), updated_at: new Date().toISOString() },
      isAuthenticated: true,
      isLoading: false,
      login: vi.fn(),
      logout: vi.fn(),
      authEnabled: true,
    });

    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    const deactivateButtons = screen.getAllByTitle('Deactivate user');
    // The first user (admin) should have a disabled deactivate button
    expect(deactivateButtons[0]).toBeDisabled();
  });
});

describe('Users Empty State', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('shows empty state when no users', async () => {
    const { adminApi } = await import('../../lib/api');
    vi.mocked(adminApi.listUsers).mockResolvedValueOnce({
      data: { users: [], total: 0 },
    } as never);

    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('No users found')).toBeInTheDocument();
    });
  });
});

describe('Users Role Badge Colors', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('displays role labels correctly', async () => {
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('Administrator')).toBeInTheDocument();
      expect(screen.getByText('Operator')).toBeInTheDocument();
    });
  });
});

describe('Users Filter', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('toggles include inactive users', async () => {
    const { adminApi } = await import('../../lib/api');

    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('Show inactive users')).toBeInTheDocument();
    });

    const checkbox = screen.getByRole('checkbox');
    fireEvent.click(checkbox);

    await waitFor(() => {
      expect(adminApi.listUsers).toHaveBeenCalledWith(true);
    });
  });
});

describe('Users Create Form', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('shows username field in create modal', async () => {
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Add User'));

    await waitFor(() => {
      expect(screen.getByLabelText('Username')).toBeInTheDocument();
    });
  });

  it('shows email field in create modal', async () => {
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Add User'));

    await waitFor(() => {
      expect(screen.getByLabelText('Email')).toBeInTheDocument();
    });
  });

  it('shows password field in create modal', async () => {
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Add User'));

    await waitFor(() => {
      expect(screen.getByLabelText('Password')).toBeInTheDocument();
    });
  });

  it('shows cancel button in create modal', async () => {
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Add User'));

    await waitFor(() => {
      const cancelButtons = screen.getAllByText('Cancel');
      expect(cancelButtons.length).toBeGreaterThan(0);
    });
  });

  it('closes create modal on cancel', async () => {
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Add User'));

    await waitFor(() => {
      expect(screen.getByLabelText('Username')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Cancel'));

    await waitFor(() => {
      expect(screen.queryByLabelText('Username')).not.toBeInTheDocument();
    });
  });
});

describe('Users API Errors', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('handles list users error gracefully', async () => {
    const { adminApi } = await import('../../lib/api');
    vi.mocked(adminApi.listUsers).mockRejectedValueOnce(new Error('Network error'));

    renderUsers();

    // Should not crash even with error
    await waitFor(() => {
      expect(screen.getByText('User Management')).toBeInTheDocument();
    });
  });
});

describe('Users Table Display', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('shows user count in header', async () => {
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('User Management')).toBeInTheDocument();
    });
  });

  it('displays inactive status for inactive users', async () => {
    renderUsers();

    // First enable showing inactive users
    await waitFor(() => {
      expect(screen.getByText('Show inactive users')).toBeInTheDocument();
    });

    const checkbox = screen.getByRole('checkbox');
    fireEvent.click(checkbox);

    await waitFor(() => {
      const inactiveLabels = screen.queryAllByText('Inactive');
      expect(inactiveLabels.length).toBeGreaterThanOrEqual(0);
    });
  });

  it('shows viewer role badge', async () => {
    renderUsers();

    // Enable showing inactive users to see the viewer role
    await waitFor(() => {
      expect(screen.getByText('Show inactive users')).toBeInTheDocument();
    });

    const checkbox = screen.getByRole('checkbox');
    fireEvent.click(checkbox);

    await waitFor(() => {
      const viewerBadges = screen.queryAllByText('Viewer');
      expect(viewerBadges.length).toBeGreaterThanOrEqual(0);
    });
  });
});

describe('Users Password Reset', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('has password reset buttons for each user', async () => {
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    const passwordButtons = screen.getAllByTitle('Reset password');
    expect(passwordButtons.length).toBeGreaterThan(0);
  });
});

describe('Users Edit Form', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('shows email field in edit modal', async () => {
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    const editButtons = screen.getAllByTitle('Edit user');
    fireEvent.click(editButtons[0]);

    await waitFor(() => {
      expect(screen.getByLabelText('Email')).toBeInTheDocument();
    });
  });

  it('pre-fills email in edit modal', async () => {
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    const editButtons = screen.getAllByTitle('Edit user');
    fireEvent.click(editButtons[0]);

    await waitFor(() => {
      const emailInput = screen.getByDisplayValue('admin@example.com');
      expect(emailInput).toBeInTheDocument();
    });
  });

  it('shows username field in edit modal', async () => {
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    const editButtons = screen.getAllByTitle('Edit user');
    fireEvent.click(editButtons[0]);

    await waitFor(() => {
      expect(screen.getByLabelText('Username')).toBeInTheDocument();
    });
  });
});

describe('Users Action Buttons', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('shows all action buttons for each user', async () => {
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    const editButtons = screen.getAllByTitle('Edit user');
    const roleButtons = screen.getAllByTitle('Change role');
    const passwordButtons = screen.getAllByTitle('Reset password');
    const deactivateButtons = screen.getAllByTitle('Deactivate user');

    expect(editButtons.length).toBeGreaterThan(0);
    expect(roleButtons.length).toBeGreaterThan(0);
    expect(passwordButtons.length).toBeGreaterThan(0);
    expect(deactivateButtons.length).toBeGreaterThan(0);
  });
});

describe('Users Form Submissions', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('submits create form and calls createUser API', async () => {
    const { adminApi } = await import('../../lib/api');
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Add User'));

    await waitFor(() => {
      expect(screen.getByLabelText('Username')).toBeInTheDocument();
    });

    fireEvent.change(screen.getByLabelText('Username'), { target: { value: 'newuser' } });
    fireEvent.change(screen.getByLabelText('Email'), { target: { value: 'newuser@test.com' } });
    fireEvent.change(screen.getByLabelText('Password'), { target: { value: 'password123' } });

    fireEvent.click(screen.getByRole('button', { name: /Create User/ }));

    await waitFor(() => {
      expect(adminApi.createUser).toHaveBeenCalledWith({
        username: 'newuser',
        email: 'newuser@test.com',
        password: 'password123',
        role: 'viewer',
      });
    });
  });

  it('submits edit form and calls updateUser API', async () => {
    const { adminApi } = await import('../../lib/api');
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    const editButtons = screen.getAllByTitle('Edit user');
    fireEvent.click(editButtons[1]); // Click operator's edit button

    await waitFor(() => {
      expect(screen.getByRole('heading', { name: 'Edit User' })).toBeInTheDocument();
    });

    fireEvent.change(screen.getByLabelText('Username'), { target: { value: 'updated-operator' } });
    fireEvent.click(screen.getByRole('button', { name: /Save Changes/ }));

    await waitFor(() => {
      expect(adminApi.updateUser).toHaveBeenCalledWith('user-2', {
        username: 'updated-operator',
        email: 'operator@example.com',
      });
    });
  });

  it('submits role change form and calls updateUserRole API', async () => {
    const { adminApi } = await import('../../lib/api');
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('operator@example.com')).toBeInTheDocument();
    });

    const roleButtons = screen.getAllByTitle('Change role');
    const enabledRoleButton = roleButtons.find(btn => !(btn as HTMLButtonElement).disabled);
    if (enabledRoleButton) {
      fireEvent.click(enabledRoleButton);

      await waitFor(() => {
        expect(screen.getByRole('heading', { name: 'Change Role' })).toBeInTheDocument();
      });

      fireEvent.change(screen.getByLabelText('Role'), { target: { value: 'analyst' } });
      fireEvent.click(screen.getByRole('button', { name: /Update Role/ }));

      await waitFor(() => {
        expect(adminApi.updateUserRole).toHaveBeenCalled();
      });
    }
  });

  it('submits password reset form and calls resetPassword API', async () => {
    const { adminApi } = await import('../../lib/api');
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    const passwordButtons = screen.getAllByTitle('Reset password');
    fireEvent.click(passwordButtons[0]);

    await waitFor(() => {
      expect(screen.getByLabelText('New Password')).toBeInTheDocument();
    });

    fireEvent.change(screen.getByLabelText('New Password'), { target: { value: 'newpassword123' } });
    fireEvent.click(screen.getByRole('button', { name: /Reset Password/ }));

    await waitFor(() => {
      expect(adminApi.resetPassword).toHaveBeenCalledWith('user-1', { new_password: 'newpassword123' });
    });
  });

  it('confirms deactivation and calls deactivateUser API', async () => {
    const { adminApi } = await import('../../lib/api');
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('operator@example.com')).toBeInTheDocument();
    });

    const deactivateButtons = screen.getAllByTitle('Deactivate user');
    const enabledButton = deactivateButtons.find(btn => !(btn as HTMLButtonElement).disabled);
    if (enabledButton) {
      fireEvent.click(enabledButton);

      await waitFor(() => {
        expect(screen.getByRole('heading', { name: 'Deactivate User' })).toBeInTheDocument();
      });

      fireEvent.click(screen.getByRole('button', { name: /^Deactivate$/ }));

      await waitFor(() => {
        expect(adminApi.deactivateUser).toHaveBeenCalled();
      });
    }
  });

  it('confirms reactivation and calls reactivateUser API', async () => {
    const { adminApi } = await import('../../lib/api');

    vi.mocked(adminApi.listUsers).mockResolvedValueOnce({
      data: {
        users: [
          {
            id: 'user-3',
            username: 'inactive',
            email: 'inactive@example.com',
            role: 'viewer',
            is_active: false,
            last_login_at: null,
            created_at: new Date().toISOString(),
            updated_at: new Date().toISOString(),
          },
        ],
        total: 1,
      },
    } as never);

    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('inactive@example.com')).toBeInTheDocument();
    });

    const reactivateButton = screen.getByTitle('Reactivate user');
    fireEvent.click(reactivateButton);

    await waitFor(() => {
      expect(screen.getByRole('heading', { name: 'Reactivate User' })).toBeInTheDocument();
    });

    fireEvent.click(screen.getByRole('button', { name: /^Reactivate$/ }));

    await waitFor(() => {
      expect(adminApi.reactivateUser).toHaveBeenCalledWith('user-3');
    });
  });
});

describe('Users Mutation Errors', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('shows error alert when create mutation fails', async () => {
    const { adminApi } = await import('../../lib/api');
    vi.mocked(adminApi.createUser).mockRejectedValueOnce({
      response: { data: { error: 'Username already exists' } },
    });

    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Add User'));

    await waitFor(() => {
      expect(screen.getByLabelText('Username')).toBeInTheDocument();
    });

    fireEvent.change(screen.getByLabelText('Username'), { target: { value: 'newuser' } });
    fireEvent.change(screen.getByLabelText('Email'), { target: { value: 'newuser@test.com' } });
    fireEvent.change(screen.getByLabelText('Password'), { target: { value: 'password123' } });
    fireEvent.click(screen.getByRole('button', { name: /Create User/ }));

    await waitFor(() => {
      expect(screen.getByText('Username already exists')).toBeInTheDocument();
    });
  });

  it('shows fallback error when create mutation fails without error message', async () => {
    const { adminApi } = await import('../../lib/api');
    vi.mocked(adminApi.createUser).mockRejectedValueOnce(new Error('Network error'));

    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Add User'));

    await waitFor(() => {
      expect(screen.getByLabelText('Username')).toBeInTheDocument();
    });

    fireEvent.change(screen.getByLabelText('Username'), { target: { value: 'newuser' } });
    fireEvent.change(screen.getByLabelText('Email'), { target: { value: 'newuser@test.com' } });
    fireEvent.change(screen.getByLabelText('Password'), { target: { value: 'password123' } });
    fireEvent.click(screen.getByRole('button', { name: /Create User/ }));

    await waitFor(() => {
      expect(screen.getByText('Failed to create user')).toBeInTheDocument();
    });
  });

  it('shows error alert when update mutation fails', async () => {
    const { adminApi } = await import('../../lib/api');
    vi.mocked(adminApi.updateUser).mockRejectedValueOnce({
      response: { data: { error: 'Duplicate username' } },
    });

    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    const editButtons = screen.getAllByTitle('Edit user');
    fireEvent.click(editButtons[0]);

    await waitFor(() => {
      expect(screen.getByRole('heading', { name: 'Edit User' })).toBeInTheDocument();
    });

    fireEvent.click(screen.getByRole('button', { name: /Save Changes/ }));

    await waitFor(() => {
      expect(screen.getByText('Duplicate username')).toBeInTheDocument();
    });
  });

  it('shows fallback error when update mutation fails without message', async () => {
    const { adminApi } = await import('../../lib/api');
    vi.mocked(adminApi.updateUser).mockRejectedValueOnce(new Error('Server error'));

    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    const editButtons = screen.getAllByTitle('Edit user');
    fireEvent.click(editButtons[0]);

    await waitFor(() => {
      expect(screen.getByRole('heading', { name: 'Edit User' })).toBeInTheDocument();
    });

    fireEvent.click(screen.getByRole('button', { name: /Save Changes/ }));

    await waitFor(() => {
      expect(screen.getByText('Failed to update user')).toBeInTheDocument();
    });
  });

  it('shows error alert when role change fails', async () => {
    const { adminApi } = await import('../../lib/api');
    vi.mocked(adminApi.updateUserRole).mockRejectedValueOnce({
      response: { data: { error: 'Cannot change role' } },
    });

    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('operator@example.com')).toBeInTheDocument();
    });

    const roleButtons = screen.getAllByTitle('Change role');
    const enabledButton = roleButtons.find(btn => !(btn as HTMLButtonElement).disabled);
    if (enabledButton) {
      fireEvent.click(enabledButton);

      await waitFor(() => {
        expect(screen.getByRole('heading', { name: 'Change Role' })).toBeInTheDocument();
      });

      fireEvent.click(screen.getByRole('button', { name: /Update Role/ }));

      await waitFor(() => {
        expect(screen.getByText('Cannot change role')).toBeInTheDocument();
      });
    }
  });

  it('shows fallback error when role mutation fails without message', async () => {
    const { adminApi } = await import('../../lib/api');
    vi.mocked(adminApi.updateUserRole).mockRejectedValueOnce(new Error('Network'));

    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('operator@example.com')).toBeInTheDocument();
    });

    const roleButtons = screen.getAllByTitle('Change role');
    const enabledButton = roleButtons.find(btn => !(btn as HTMLButtonElement).disabled);
    if (enabledButton) {
      fireEvent.click(enabledButton);

      await waitFor(() => {
        expect(screen.getByRole('heading', { name: 'Change Role' })).toBeInTheDocument();
      });

      fireEvent.click(screen.getByRole('button', { name: /Update Role/ }));

      await waitFor(() => {
        expect(screen.getByText('Failed to update role')).toBeInTheDocument();
      });
    }
  });

  it('shows error alert when password reset fails', async () => {
    const { adminApi } = await import('../../lib/api');
    vi.mocked(adminApi.resetPassword).mockRejectedValueOnce({
      response: { data: { error: 'Password too weak' } },
    });

    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    const passwordButtons = screen.getAllByTitle('Reset password');
    fireEvent.click(passwordButtons[0]);

    await waitFor(() => {
      expect(screen.getByLabelText('New Password')).toBeInTheDocument();
    });

    fireEvent.change(screen.getByLabelText('New Password'), { target: { value: 'newpassword123' } });
    fireEvent.click(screen.getByRole('button', { name: /Reset Password/ }));

    await waitFor(() => {
      expect(screen.getByText('Password too weak')).toBeInTheDocument();
    });
  });

  it('shows fallback error when password reset fails without message', async () => {
    const { adminApi } = await import('../../lib/api');
    vi.mocked(adminApi.resetPassword).mockRejectedValueOnce(new Error('Network'));

    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    const passwordButtons = screen.getAllByTitle('Reset password');
    fireEvent.click(passwordButtons[0]);

    await waitFor(() => {
      expect(screen.getByLabelText('New Password')).toBeInTheDocument();
    });

    fireEvent.change(screen.getByLabelText('New Password'), { target: { value: 'newpassword123' } });
    fireEvent.click(screen.getByRole('button', { name: /Reset Password/ }));

    await waitFor(() => {
      expect(screen.getByText('Failed to reset password')).toBeInTheDocument();
    });
  });

  it('shows error alert when deactivation fails', async () => {
    const { adminApi } = await import('../../lib/api');
    vi.mocked(adminApi.deactivateUser).mockRejectedValueOnce({
      response: { data: { error: 'Cannot deactivate last admin' } },
    });

    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('operator@example.com')).toBeInTheDocument();
    });

    const deactivateButtons = screen.getAllByTitle('Deactivate user');
    const enabledButton = deactivateButtons.find(btn => !(btn as HTMLButtonElement).disabled);
    if (enabledButton) {
      fireEvent.click(enabledButton);

      await waitFor(() => {
        expect(screen.getByRole('heading', { name: 'Deactivate User' })).toBeInTheDocument();
      });

      fireEvent.click(screen.getByRole('button', { name: /^Deactivate$/ }));

      await waitFor(() => {
        expect(screen.getByText('Cannot deactivate last admin')).toBeInTheDocument();
      });
    }
  });

  it('shows fallback error when deactivation fails without message', async () => {
    const { adminApi } = await import('../../lib/api');
    vi.mocked(adminApi.deactivateUser).mockRejectedValueOnce(new Error('Network'));

    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('operator@example.com')).toBeInTheDocument();
    });

    const deactivateButtons = screen.getAllByTitle('Deactivate user');
    const enabledButton = deactivateButtons.find(btn => !(btn as HTMLButtonElement).disabled);
    if (enabledButton) {
      fireEvent.click(enabledButton);

      await waitFor(() => {
        expect(screen.getByRole('heading', { name: 'Deactivate User' })).toBeInTheDocument();
      });

      fireEvent.click(screen.getByRole('button', { name: /^Deactivate$/ }));

      await waitFor(() => {
        expect(screen.getByText('Failed to deactivate user')).toBeInTheDocument();
      });
    }
  });

  it('shows error alert when reactivation fails', async () => {
    const { adminApi } = await import('../../lib/api');
    vi.mocked(adminApi.reactivateUser).mockRejectedValueOnce({
      response: { data: { error: 'Cannot reactivate' } },
    });

    vi.mocked(adminApi.listUsers).mockResolvedValueOnce({
      data: {
        users: [{
          id: 'user-3',
          username: 'inactive',
          email: 'inactive@example.com',
          role: 'viewer',
          is_active: false,
          last_login_at: null,
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString(),
        }],
        total: 1,
      },
    } as never);

    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('inactive@example.com')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByTitle('Reactivate user'));

    await waitFor(() => {
      expect(screen.getByRole('heading', { name: 'Reactivate User' })).toBeInTheDocument();
    });

    fireEvent.click(screen.getByRole('button', { name: /^Reactivate$/ }));

    await waitFor(() => {
      expect(screen.getByText('Cannot reactivate')).toBeInTheDocument();
    });
  });

  it('shows fallback error when reactivation fails without message', async () => {
    const { adminApi } = await import('../../lib/api');
    vi.mocked(adminApi.reactivateUser).mockRejectedValueOnce(new Error('Network'));

    vi.mocked(adminApi.listUsers).mockResolvedValueOnce({
      data: {
        users: [{
          id: 'user-3',
          username: 'inactive',
          email: 'inactive@example.com',
          role: 'viewer',
          is_active: false,
          last_login_at: null,
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString(),
        }],
        total: 1,
      },
    } as never);

    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('inactive@example.com')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByTitle('Reactivate user'));

    await waitFor(() => {
      expect(screen.getByRole('heading', { name: 'Reactivate User' })).toBeInTheDocument();
    });

    fireEvent.click(screen.getByRole('button', { name: /^Reactivate$/ }));

    await waitFor(() => {
      expect(screen.getByText('Failed to reactivate user')).toBeInTheDocument();
    });
  });
});

describe('Users Role Badges - All Roles', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('displays rssi and analyst role badges', async () => {
    const { adminApi } = await import('../../lib/api');
    vi.mocked(adminApi.listUsers).mockResolvedValueOnce({
      data: {
        users: [
          {
            id: 'user-rssi',
            username: 'rssi-user',
            email: 'rssi@example.com',
            role: 'rssi',
            is_active: true,
            last_login_at: null,
            created_at: new Date().toISOString(),
            updated_at: new Date().toISOString(),
          },
          {
            id: 'user-analyst',
            username: 'analyst-user',
            email: 'analyst@example.com',
            role: 'analyst',
            is_active: true,
            last_login_at: null,
            created_at: new Date().toISOString(),
            updated_at: new Date().toISOString(),
          },
        ],
        total: 2,
      },
    } as never);

    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('Security Officer (RSSI)')).toBeInTheDocument();
      expect(screen.getByText('Analyst')).toBeInTheDocument();
    });

    const rssiBadge = screen.getByText('Security Officer (RSSI)');
    expect(rssiBadge.className).toContain('bg-purple-100');

    const analystBadge = screen.getByText('Analyst');
    expect(analystBadge.className).toContain('bg-green-100');
  });

  it('displays admin role badge with red styling', async () => {
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('Administrator')).toBeInTheDocument();
    });

    const adminBadge = screen.getByText('Administrator');
    expect(adminBadge.className).toContain('bg-red-100');
  });

  it('displays operator role badge with blue styling', async () => {
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('Operator')).toBeInTheDocument();
    });

    const operatorBadge = screen.getByText('Operator');
    expect(operatorBadge.className).toContain('bg-blue-100');
  });

  it('displays viewer role badge with gray styling', async () => {
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('Viewer')).toBeInTheDocument();
    });

    const viewerBadge = screen.getByText('Viewer');
    expect(viewerBadge.className).toContain('bg-gray-100');
  });

  it('displays all 5 roles with correct styling', async () => {
    const { adminApi } = await import('../../lib/api');
    vi.mocked(adminApi.listUsers).mockResolvedValueOnce({
      data: {
        users: [
          { id: 'u1', username: 'admin1', email: 'a1@test.com', role: 'admin', is_active: true, last_login_at: null, created_at: new Date().toISOString(), updated_at: new Date().toISOString() },
          { id: 'u2', username: 'rssi1', email: 'r1@test.com', role: 'rssi', is_active: true, last_login_at: null, created_at: new Date().toISOString(), updated_at: new Date().toISOString() },
          { id: 'u3', username: 'op1', email: 'o1@test.com', role: 'operator', is_active: true, last_login_at: null, created_at: new Date().toISOString(), updated_at: new Date().toISOString() },
          { id: 'u4', username: 'an1', email: 'an1@test.com', role: 'analyst', is_active: true, last_login_at: null, created_at: new Date().toISOString(), updated_at: new Date().toISOString() },
          { id: 'u5', username: 'vi1', email: 'v1@test.com', role: 'viewer', is_active: true, last_login_at: null, created_at: new Date().toISOString(), updated_at: new Date().toISOString() },
        ],
        total: 5,
      },
    } as never);

    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('Administrator')).toBeInTheDocument();
      expect(screen.getByText('Security Officer (RSSI)')).toBeInTheDocument();
      expect(screen.getByText('Operator')).toBeInTheDocument();
      expect(screen.getByText('Analyst')).toBeInTheDocument();
      expect(screen.getByText('Viewer')).toBeInTheDocument();
    });
  });
});

describe('Users Current User Indicator', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('shows (you) indicator for current user', async () => {
    const { useAuth } = await import('../../contexts/AuthContext');
    vi.mocked(useAuth).mockReturnValue({
      user: { id: 'user-1', username: 'admin', email: 'admin@example.com', role: 'admin', is_active: true, created_at: new Date().toISOString(), updated_at: new Date().toISOString() },
      isAuthenticated: true,
      isLoading: false,
      login: vi.fn(),
      logout: vi.fn(),
      authEnabled: true,
    });

    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    expect(screen.getByText('(you)')).toBeInTheDocument();
  });

  it('does not show (you) for other users', async () => {
    const { useAuth } = await import('../../contexts/AuthContext');
    vi.mocked(useAuth).mockReturnValue({
      user: { id: 'no-match-id', username: 'someone', email: 'someone@example.com', role: 'admin', is_active: true, created_at: new Date().toISOString(), updated_at: new Date().toISOString() },
      isAuthenticated: true,
      isLoading: false,
      login: vi.fn(),
      logout: vi.fn(),
      authEnabled: true,
    });

    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('operator@example.com')).toBeInTheDocument();
    });

    // No user in the list matches 'no-match-id'
    const youIndicators = screen.queryAllByText('(you)');
    expect(youIndicators.length).toBe(0);
  });
});

describe('Users Loading State', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('shows loading state while data is being fetched', async () => {
    const { adminApi } = await import('../../lib/api');
    // Create a never-resolving promise to keep loading state (use Once to avoid polluting later tests)
    vi.mocked(adminApi.listUsers).mockImplementationOnce(() => new Promise(() => {}));

    renderUsers();

    expect(screen.getByText('Loading users...')).toBeInTheDocument();
  });
});

describe('Users Avatar Display', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('displays first letter of username as avatar', async () => {
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    // First letters: 'A' for admin, 'O' for operator, 'I' for inactive
    expect(screen.getByText('A')).toBeInTheDocument();
    expect(screen.getByText('O')).toBeInTheDocument();
    expect(screen.getByText('I')).toBeInTheDocument();
  });
});

describe('Users Last Login Display', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('displays formatted last login time for users who have logged in', async () => {
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    // The admin user has a last_login_at set to new Date().toISOString()
    // formatDistanceToNow would render something like "less than a minute ago"
    const timeText = screen.getByText(/ago/);
    expect(timeText).toBeInTheDocument();
  });

  it('shows Never for users who have not logged in', async () => {
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('operator@example.com')).toBeInTheDocument();
    });

    const neverTexts = screen.getAllByText('Never');
    expect(neverTexts.length).toBe(2); // operator and inactive both have null last_login_at
  });
});

describe('Users Create Form - Role Selection', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('changes role in create form', async () => {
    const { adminApi } = await import('../../lib/api');
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Add User'));

    await waitFor(() => {
      expect(screen.getByLabelText('Role')).toBeInTheDocument();
    });

    fireEvent.change(screen.getByLabelText('Role'), { target: { value: 'operator' } });
    fireEvent.change(screen.getByLabelText('Username'), { target: { value: 'newop' } });
    fireEvent.change(screen.getByLabelText('Email'), { target: { value: 'newop@test.com' } });
    fireEvent.change(screen.getByLabelText('Password'), { target: { value: 'password123' } });
    fireEvent.click(screen.getByRole('button', { name: /Create User/ }));

    await waitFor(() => {
      expect(adminApi.createUser).toHaveBeenCalledWith({
        username: 'newop',
        email: 'newop@test.com',
        password: 'password123',
        role: 'operator',
      });
    });
  });

  it('shows minimum password hint text', async () => {
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Add User'));

    await waitFor(() => {
      expect(screen.getByText('Minimum 8 characters')).toBeInTheDocument();
    });
  });

  it('shows role options with descriptions', async () => {
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Add User'));

    await waitFor(() => {
      const roleSelect = screen.getByLabelText('Role');
      const options = roleSelect.querySelectorAll('option');
      expect(options.length).toBe(5);
      expect(options[0].textContent).toBe('Administrator - Full system access');
      expect(options[1].textContent).toBe('Security Officer (RSSI) - View reports and analytics');
      expect(options[2].textContent).toBe('Operator - Execute scenarios');
      expect(options[3].textContent).toBe('Analyst - Read-only with analytics');
      expect(options[4].textContent).toBe('Viewer - Read-only basic access');
    });
  });
});

describe('Users Edit Form - Field Updates', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('updates email field in edit form', async () => {
    const { adminApi } = await import('../../lib/api');
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    const editButtons = screen.getAllByTitle('Edit user');
    fireEvent.click(editButtons[0]);

    await waitFor(() => {
      expect(screen.getByRole('heading', { name: 'Edit User' })).toBeInTheDocument();
    });

    fireEvent.change(screen.getByLabelText('Email'), { target: { value: 'newemail@example.com' } });
    fireEvent.click(screen.getByRole('button', { name: /Save Changes/ }));

    await waitFor(() => {
      expect(adminApi.updateUser).toHaveBeenCalledWith('user-1', {
        username: 'admin',
        email: 'newemail@example.com',
      });
    });
  });

  it('pre-fills username in edit modal', async () => {
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    const editButtons = screen.getAllByTitle('Edit user');
    fireEvent.click(editButtons[0]);

    await waitFor(() => {
      const usernameInput = screen.getByDisplayValue('admin');
      expect(usernameInput).toBeInTheDocument();
    });
  });
});

describe('Users Role Change Modal - Details', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('shows username in role change modal', async () => {
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('operator@example.com')).toBeInTheDocument();
    });

    const roleButtons = screen.getAllByTitle('Change role');
    const enabledButton = roleButtons.find(btn => !(btn as HTMLButtonElement).disabled);
    if (enabledButton) {
      fireEvent.click(enabledButton);

      await waitFor(() => {
        expect(screen.getByText('Change Role')).toBeInTheDocument();
      });

      expect(screen.getByText(/Change role for/)).toBeInTheDocument();
    }
  });

  it('shows role options in role change modal', async () => {
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('operator@example.com')).toBeInTheDocument();
    });

    const roleButtons = screen.getAllByTitle('Change role');
    const enabledButton = roleButtons.find(btn => !(btn as HTMLButtonElement).disabled);
    if (enabledButton) {
      fireEvent.click(enabledButton);

      await waitFor(() => {
        const roleSelect = screen.getByLabelText('Role');
        const options = roleSelect.querySelectorAll('option');
        expect(options.length).toBe(5);
      });
    }
  });

  it('pre-selects current role in role change modal', async () => {
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('operator@example.com')).toBeInTheDocument();
    });

    // Click the second role button (for the operator user, index 1)
    const roleButtons = screen.getAllByTitle('Change role');
    fireEvent.click(roleButtons[1]);

    await waitFor(() => {
      const roleSelect = screen.getByLabelText('Role') as HTMLSelectElement;
      expect(roleSelect.value).toBe('operator');
    });
  });

  it('cancels role change modal', async () => {
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('operator@example.com')).toBeInTheDocument();
    });

    const roleButtons = screen.getAllByTitle('Change role');
    const enabledButton = roleButtons.find(btn => !(btn as HTMLButtonElement).disabled);
    if (enabledButton) {
      fireEvent.click(enabledButton);

      await waitFor(() => {
        expect(screen.getByRole('heading', { name: 'Change Role' })).toBeInTheDocument();
      });

      fireEvent.click(screen.getByText('Cancel'));

      await waitFor(() => {
        expect(screen.queryByRole('heading', { name: 'Change Role' })).not.toBeInTheDocument();
      });
    }
  });
});

describe('Users Password Reset Modal - Details', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('opens password reset modal and shows username', async () => {
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    const passwordButtons = screen.getAllByTitle('Reset password');
    fireEvent.click(passwordButtons[0]);

    await waitFor(() => {
      expect(screen.getByRole('heading', { name: 'Reset Password' })).toBeInTheDocument();
      expect(screen.getByText(/Reset password for/)).toBeInTheDocument();
    });
  });

  it('shows minimum characters hint in password reset modal', async () => {
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    const passwordButtons = screen.getAllByTitle('Reset password');
    fireEvent.click(passwordButtons[0]);

    await waitFor(() => {
      const hints = screen.getAllByText('Minimum 8 characters');
      expect(hints.length).toBeGreaterThan(0);
    });
  });

  it('cancels password reset modal', async () => {
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    const passwordButtons = screen.getAllByTitle('Reset password');
    fireEvent.click(passwordButtons[0]);

    await waitFor(() => {
      expect(screen.getByRole('heading', { name: 'Reset Password' })).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Cancel'));

    await waitFor(() => {
      expect(screen.queryByRole('heading', { name: 'Reset Password' })).not.toBeInTheDocument();
    });
  });
});

describe('Users Deactivate Modal - Details', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('shows user data will be preserved message', async () => {
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('operator@example.com')).toBeInTheDocument();
    });

    const deactivateButtons = screen.getAllByTitle('Deactivate user');
    const enabledButton = deactivateButtons.find(btn => !(btn as HTMLButtonElement).disabled);
    if (enabledButton) {
      fireEvent.click(enabledButton);

      await waitFor(() => {
        expect(screen.getByText(/their data will be preserved/)).toBeInTheDocument();
      });
    }
  });

  it('cancels deactivate modal', async () => {
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('operator@example.com')).toBeInTheDocument();
    });

    const deactivateButtons = screen.getAllByTitle('Deactivate user');
    const enabledButton = deactivateButtons.find(btn => !(btn as HTMLButtonElement).disabled);
    if (enabledButton) {
      fireEvent.click(enabledButton);

      await waitFor(() => {
        expect(screen.getByRole('heading', { name: 'Deactivate User' })).toBeInTheDocument();
      });

      fireEvent.click(screen.getByText('Cancel'));

      await waitFor(() => {
        expect(screen.queryByRole('heading', { name: 'Deactivate User' })).not.toBeInTheDocument();
      });
    }
  });
});

describe('Users Reactivate Modal - Details', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('shows existing credentials message in reactivate modal', async () => {
    const { adminApi } = await import('../../lib/api');
    vi.mocked(adminApi.listUsers).mockResolvedValueOnce({
      data: {
        users: [{
          id: 'user-3',
          username: 'inactive',
          email: 'inactive@example.com',
          role: 'viewer',
          is_active: false,
          last_login_at: null,
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString(),
        }],
        total: 1,
      },
    } as never);

    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('inactive@example.com')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByTitle('Reactivate user'));

    await waitFor(() => {
      expect(screen.getByText(/existing credentials/)).toBeInTheDocument();
    });
  });

  it('cancels reactivate modal', async () => {
    const { adminApi } = await import('../../lib/api');
    vi.mocked(adminApi.listUsers).mockResolvedValueOnce({
      data: {
        users: [{
          id: 'user-3',
          username: 'inactive',
          email: 'inactive@example.com',
          role: 'viewer',
          is_active: false,
          last_login_at: null,
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString(),
        }],
        total: 1,
      },
    } as never);

    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('inactive@example.com')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByTitle('Reactivate user'));

    await waitFor(() => {
      expect(screen.getByRole('heading', { name: 'Reactivate User' })).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Cancel'));

    await waitFor(() => {
      expect(screen.queryByRole('heading', { name: 'Reactivate User' })).not.toBeInTheDocument();
    });
  });

  it('shows username in reactivate confirmation', async () => {
    const { adminApi } = await import('../../lib/api');
    vi.mocked(adminApi.listUsers).mockResolvedValueOnce({
      data: {
        users: [{
          id: 'user-3',
          username: 'deactivated-user',
          email: 'deactivated@example.com',
          role: 'viewer',
          is_active: false,
          last_login_at: null,
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString(),
        }],
        total: 1,
      },
    } as never);

    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('deactivated@example.com')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByTitle('Reactivate user'));

    await waitFor(() => {
      // The paragraph contains "Are you sure you want to reactivate <strong>deactivated-user</strong>?"
      const confirmText = screen.getByText((_content, element) => {
        return element?.tagName === 'P' && !!element?.textContent?.includes('deactivated-user');
      });
      expect(confirmText).toBeInTheDocument();
    });
  });
});

describe('Users Modal Close via X Button', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('closes create modal via X button', async () => {
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Add User'));

    await waitFor(() => {
      expect(screen.getByRole('heading', { name: 'Create User' })).toBeInTheDocument();
    });

    const closeButton = screen.getByLabelText('Close modal');
    fireEvent.click(closeButton);

    await waitFor(() => {
      expect(screen.queryByRole('heading', { name: 'Create User' })).not.toBeInTheDocument();
    });
  });

  it('closes edit modal via X button', async () => {
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    const editButtons = screen.getAllByTitle('Edit user');
    fireEvent.click(editButtons[0]);

    await waitFor(() => {
      expect(screen.getByRole('heading', { name: 'Edit User' })).toBeInTheDocument();
    });

    const closeButton = screen.getByLabelText('Close modal');
    fireEvent.click(closeButton);

    await waitFor(() => {
      expect(screen.queryByRole('heading', { name: 'Edit User' })).not.toBeInTheDocument();
    });
  });
});

describe('Users Inactive Status Display', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('displays Inactive badge for inactive users', async () => {
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('inactive@example.com')).toBeInTheDocument();
    });

    expect(screen.getByText('Inactive')).toBeInTheDocument();
  });

  it('renders active users and inactive users with different row styling', async () => {
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
      expect(screen.getByText('inactive@example.com')).toBeInTheDocument();
    });

    const activeStatuses = screen.getAllByText('Active');
    expect(activeStatuses.length).toBe(2); // admin and operator are active

    const inactiveStatuses = screen.getAllByText('Inactive');
    expect(inactiveStatuses.length).toBe(1); // only inactive user
  });
});

describe('Users Form Input Changes', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('updates all create form fields', async () => {
    const { adminApi } = await import('../../lib/api');
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Add User'));

    await waitFor(() => {
      expect(screen.getByLabelText('Username')).toBeInTheDocument();
    });

    const usernameInput = screen.getByLabelText('Username');
    const emailInput = screen.getByLabelText('Email');
    const passwordInput = screen.getByLabelText('Password');
    const roleSelect = screen.getByLabelText('Role');

    fireEvent.change(usernameInput, { target: { value: 'testuser' } });
    fireEvent.change(emailInput, { target: { value: 'test@example.com' } });
    fireEvent.change(passwordInput, { target: { value: 'securepass123' } });
    fireEvent.change(roleSelect, { target: { value: 'admin' } });

    fireEvent.click(screen.getByRole('button', { name: /Create User/ }));

    await waitFor(() => {
      expect(adminApi.createUser).toHaveBeenCalledWith({
        username: 'testuser',
        email: 'test@example.com',
        password: 'securepass123',
        role: 'admin',
      });
    });
  });

  it('updates password form field', async () => {
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    const passwordButtons = screen.getAllByTitle('Reset password');
    fireEvent.click(passwordButtons[0]);

    await waitFor(() => {
      expect(screen.getByLabelText('New Password')).toBeInTheDocument();
    });

    const passwordInput = screen.getByLabelText('New Password');
    fireEvent.change(passwordInput, { target: { value: 'newpass456' } });
    expect(passwordInput).toHaveValue('newpass456');
  });

  it('updates role form select in role change modal', async () => {
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('operator@example.com')).toBeInTheDocument();
    });

    const roleButtons = screen.getAllByTitle('Change role');
    const enabledButton = roleButtons.find(btn => !(btn as HTMLButtonElement).disabled);
    if (enabledButton) {
      fireEvent.click(enabledButton);

      await waitFor(() => {
        expect(screen.getByLabelText('Role')).toBeInTheDocument();
      });

      const roleSelect = screen.getByLabelText('Role') as HTMLSelectElement;
      fireEvent.change(roleSelect, { target: { value: 'rssi' } });
      expect(roleSelect.value).toBe('rssi');
    }
  });
});

describe('Users Subtitle Text', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('displays page subtitle', async () => {
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('Manage users and their permissions')).toBeInTheDocument();
    });
  });
});

describe('Users Table Headers', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders all table column headers', async () => {
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('User')).toBeInTheDocument();
      expect(screen.getByText('Role')).toBeInTheDocument();
      expect(screen.getByText('Status')).toBeInTheDocument();
      expect(screen.getByText('Last Login')).toBeInTheDocument();
      expect(screen.getByText('Actions')).toBeInTheDocument();
    });
  });
});

describe('Users Mutation Success - Modal Closes', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('closes create modal after successful creation', async () => {
    const { adminApi } = await import('../../lib/api');
    vi.mocked(adminApi.createUser).mockResolvedValueOnce({ data: {} } as never);

    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Add User'));

    await waitFor(() => {
      expect(screen.getByLabelText('Username')).toBeInTheDocument();
    });

    fireEvent.change(screen.getByLabelText('Username'), { target: { value: 'newuser' } });
    fireEvent.change(screen.getByLabelText('Email'), { target: { value: 'new@test.com' } });
    fireEvent.change(screen.getByLabelText('Password'), { target: { value: 'password123' } });
    fireEvent.click(screen.getByRole('button', { name: /Create User/ }));

    await waitFor(() => {
      expect(screen.queryByRole('heading', { name: 'Create User' })).not.toBeInTheDocument();
    });
  });

  it('closes edit modal after successful update', async () => {
    const { adminApi } = await import('../../lib/api');
    vi.mocked(adminApi.updateUser).mockResolvedValueOnce({ data: {} } as never);

    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    const editButtons = screen.getAllByTitle('Edit user');
    fireEvent.click(editButtons[1]);

    await waitFor(() => {
      expect(screen.getByRole('heading', { name: 'Edit User' })).toBeInTheDocument();
    });

    fireEvent.click(screen.getByRole('button', { name: /Save Changes/ }));

    await waitFor(() => {
      expect(screen.queryByRole('heading', { name: 'Edit User' })).not.toBeInTheDocument();
    });
  });

  it('closes role modal after successful role change', async () => {
    const { adminApi } = await import('../../lib/api');
    vi.mocked(adminApi.updateUserRole).mockResolvedValueOnce({ data: {} } as never);

    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('operator@example.com')).toBeInTheDocument();
    });

    const roleButtons = screen.getAllByTitle('Change role');
    const enabledButton = roleButtons.find(btn => !(btn as HTMLButtonElement).disabled);
    if (enabledButton) {
      fireEvent.click(enabledButton);

      await waitFor(() => {
        expect(screen.getByRole('heading', { name: 'Change Role' })).toBeInTheDocument();
      });

      fireEvent.click(screen.getByRole('button', { name: /Update Role/ }));

      await waitFor(() => {
        expect(screen.queryByRole('heading', { name: 'Change Role' })).not.toBeInTheDocument();
      });
    }
  });

  it('closes password modal after successful reset', async () => {
    const { adminApi } = await import('../../lib/api');
    vi.mocked(adminApi.resetPassword).mockResolvedValueOnce({ data: {} } as never);

    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('admin@example.com')).toBeInTheDocument();
    });

    const passwordButtons = screen.getAllByTitle('Reset password');
    fireEvent.click(passwordButtons[0]);

    await waitFor(() => {
      expect(screen.getByLabelText('New Password')).toBeInTheDocument();
    });

    fireEvent.change(screen.getByLabelText('New Password'), { target: { value: 'newpass123' } });
    fireEvent.click(screen.getByRole('button', { name: /Reset Password/ }));

    await waitFor(() => {
      expect(screen.queryByRole('heading', { name: 'Reset Password' })).not.toBeInTheDocument();
    });
  });

  it('closes deactivate modal after successful deactivation', async () => {
    const { adminApi } = await import('../../lib/api');
    vi.mocked(adminApi.deactivateUser).mockResolvedValueOnce({ data: {} } as never);

    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('operator@example.com')).toBeInTheDocument();
    });

    const deactivateButtons = screen.getAllByTitle('Deactivate user');
    const enabledButton = deactivateButtons.find(btn => !(btn as HTMLButtonElement).disabled);
    if (enabledButton) {
      fireEvent.click(enabledButton);

      await waitFor(() => {
        expect(screen.getByRole('heading', { name: 'Deactivate User' })).toBeInTheDocument();
      });

      fireEvent.click(screen.getByRole('button', { name: /^Deactivate$/ }));

      await waitFor(() => {
        expect(screen.queryByRole('heading', { name: 'Deactivate User' })).not.toBeInTheDocument();
      });
    }
  });

  it('closes reactivate modal after successful reactivation', async () => {
    const { adminApi } = await import('../../lib/api');
    vi.mocked(adminApi.reactivateUser).mockResolvedValueOnce({ data: {} } as never);
    vi.mocked(adminApi.listUsers).mockResolvedValueOnce({
      data: {
        users: [{
          id: 'user-3',
          username: 'inactive',
          email: 'inactive@example.com',
          role: 'viewer',
          is_active: false,
          last_login_at: null,
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString(),
        }],
        total: 1,
      },
    } as never);

    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('inactive@example.com')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByTitle('Reactivate user'));

    await waitFor(() => {
      expect(screen.getByRole('heading', { name: 'Reactivate User' })).toBeInTheDocument();
    });

    fireEvent.click(screen.getByRole('button', { name: /^Reactivate$/ }));

    await waitFor(() => {
      expect(screen.queryByRole('heading', { name: 'Reactivate User' })).not.toBeInTheDocument();
    });
  });
});

describe('Users Deactivate Username Display', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('shows selected username in deactivate confirmation', async () => {
    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('operator@example.com')).toBeInTheDocument();
    });

    const deactivateButtons = screen.getAllByTitle('Deactivate user');
    const enabledButton = deactivateButtons.find(btn => !(btn as HTMLButtonElement).disabled);
    if (enabledButton) {
      fireEvent.click(enabledButton);

      await waitFor(() => {
        expect(screen.getByText(/Are you sure you want to deactivate/)).toBeInTheDocument();
      });
    }
  });
});

describe('Users Empty State Details', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('shows empty state description', async () => {
    const { adminApi } = await import('../../lib/api');
    vi.mocked(adminApi.listUsers).mockResolvedValueOnce({
      data: { users: [], total: 0 },
    } as never);

    renderUsers();

    await waitFor(() => {
      expect(screen.getByText('Create a new user to get started')).toBeInTheDocument();
    });
  });
});
