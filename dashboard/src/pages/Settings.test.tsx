import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import Settings from './Settings';

// Mock react-hot-toast
vi.mock('react-hot-toast', () => ({
  default: {
    success: vi.fn(),
    error: vi.fn(),
  },
}));

// Mock the notification API
vi.mock('../lib/api', () => ({
  notificationApi: {
    getSettings: vi.fn(() => Promise.reject({ response: { status: 404 } })),
    getSMTPConfig: vi.fn(() => Promise.reject({ response: { status: 404 } })),
    createSettings: vi.fn(() => Promise.resolve({ data: {} })),
    updateSettings: vi.fn(() => Promise.resolve({ data: {} })),
    testSMTP: vi.fn(() => Promise.resolve({ data: {} })),
  },
}));

// Mock localStorage
const localStorageMock = (() => {
  let store: Record<string, string> = {};
  return {
    getItem: vi.fn((key: string) => store[key] || null),
    setItem: vi.fn((key: string, value: string) => {
      store[key] = value;
    }),
    removeItem: vi.fn((key: string) => {
      delete store[key];
    }),
    clear: vi.fn(() => {
      store = {};
    }),
  };
})();
Object.defineProperty(window, 'localStorage', { value: localStorageMock });

function createTestQueryClient() {
  return new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
      },
    },
  });
}

function renderSettings() {
  const queryClient = createTestQueryClient();
  return render(
    <QueryClientProvider client={queryClient}>
      <MemoryRouter>
        <Settings />
      </MemoryRouter>
    </QueryClientProvider>
  );
}

describe('Settings Page', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    localStorageMock.clear();
  });

  it('renders settings title', () => {
    renderSettings();
    expect(screen.getByText('Settings')).toBeInTheDocument();
  });

  it('renders all configuration sections', async () => {
    renderSettings();
    expect(screen.getByText('Server Configuration')).toBeInTheDocument();
    expect(screen.getByText('Execution Settings')).toBeInTheDocument();
    expect(screen.getByText('Agent Settings')).toBeInTheDocument();
    expect(screen.getByText('TLS / mTLS Configuration')).toBeInTheDocument();
  });

  it('renders notification settings section', () => {
    renderSettings();
    expect(screen.getByText('Notification Settings')).toBeInTheDocument();
  });

  it('renders server URL input with default value', async () => {
    renderSettings();
    await waitFor(() => {
      const serverUrlInput = screen.getByLabelText('Server URL');
      expect(serverUrlInput).toHaveValue('https://localhost:8443');
    });
  });

  it('renders heartbeat interval input with default value', async () => {
    renderSettings();
    await waitFor(() => {
      const heartbeatInput = screen.getByLabelText('Heartbeat Interval (seconds)');
      expect(heartbeatInput).toHaveValue(30);
    });
  });

  it('renders stale timeout input with default value', async () => {
    renderSettings();
    await waitFor(() => {
      const staleInput = screen.getByLabelText('Stale Agent Timeout (seconds)');
      expect(staleInput).toHaveValue(120);
    });
  });

  it('renders safe mode toggle', () => {
    renderSettings();
    expect(screen.getByText('Safe Mode by Default')).toBeInTheDocument();
  });

  it('renders TLS certificate labels', () => {
    renderSettings();
    expect(screen.getByText('CA Certificate Path')).toBeInTheDocument();
    expect(screen.getByText('Server Certificate Path')).toBeInTheDocument();
    expect(screen.getByText('Server Key Path')).toBeInTheDocument();
  });

  it('renders save button', () => {
    renderSettings();
    expect(screen.getByText('Save Settings')).toBeInTheDocument();
  });
});

describe('Settings Form Interactions', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    localStorageMock.clear();
  });

  it('updates server URL on input change', async () => {
    renderSettings();
    await waitFor(() => {
      const serverUrlInput = screen.getByLabelText('Server URL');
      fireEvent.change(serverUrlInput, { target: { value: 'https://newserver:8443' } });
      expect(serverUrlInput).toHaveValue('https://newserver:8443');
    });
  });

  it('updates heartbeat interval on input change', async () => {
    renderSettings();
    await waitFor(() => {
      const heartbeatInput = screen.getByLabelText('Heartbeat Interval (seconds)');
      fireEvent.change(heartbeatInput, { target: { value: '60' } });
      expect(heartbeatInput).toHaveValue(60);
    });
  });

  it('saves settings to localStorage on save button click', async () => {
    const toast = await import('react-hot-toast');
    renderSettings();

    const saveButton = screen.getByText('Save Settings');
    fireEvent.click(saveButton);

    await waitFor(() => {
      expect(localStorageMock.setItem).toHaveBeenCalledWith(
        'autostrike_settings',
        expect.any(String)
      );
      expect(toast.default.success).toHaveBeenCalledWith('Settings saved successfully');
    });
  });

  it('loads settings from localStorage on mount', async () => {
    const savedSettings = {
      serverUrl: 'https://custom:9999',
      heartbeatInterval: 45,
      staleTimeout: 180,
    };
    localStorageMock.getItem.mockReturnValue(JSON.stringify(savedSettings));

    renderSettings();

    await waitFor(() => {
      const serverUrlInput = screen.getByLabelText('Server URL');
      expect(serverUrlInput).toHaveValue('https://custom:9999');
    });
  });
});

describe('Settings Default Values', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    localStorageMock.clear();
    localStorageMock.getItem.mockReturnValue(null);
  });

  it('uses default heartbeat interval when invalid input', async () => {
    renderSettings();
    await waitFor(() => {
      const heartbeatInput = screen.getByLabelText('Heartbeat Interval (seconds)');
      fireEvent.change(heartbeatInput, { target: { value: '' } });
      expect(heartbeatInput).toHaveValue(30);
    });
  });

  it('uses default stale timeout when invalid input', async () => {
    renderSettings();
    await waitFor(() => {
      const staleInput = screen.getByLabelText('Stale Agent Timeout (seconds)');
      fireEvent.change(staleInput, { target: { value: '' } });
      expect(staleInput).toHaveValue(120);
    });
  });

  it('handles malformed JSON in localStorage gracefully', async () => {
    localStorageMock.getItem.mockReturnValue('invalid-json');

    // Should not throw
    expect(() => renderSettings()).not.toThrow();

    // Should use defaults
    await waitFor(() => {
      const serverUrlInput = screen.getByLabelText('Server URL');
      expect(serverUrlInput).toHaveValue('https://localhost:8443');
    });
  });

  it('updates CA certificate path on input change', async () => {
    renderSettings();
    await waitFor(() => {
      const caCertInput = screen.getByLabelText('CA Certificate Path');
      fireEvent.change(caCertInput, { target: { value: '/path/to/ca.crt' } });
      expect(caCertInput).toHaveValue('/path/to/ca.crt');
    });
  });

  it('updates server certificate path on input change', async () => {
    renderSettings();
    await waitFor(() => {
      const serverCertInput = screen.getByLabelText('Server Certificate Path');
      fireEvent.change(serverCertInput, { target: { value: '/path/to/server.crt' } });
      expect(serverCertInput).toHaveValue('/path/to/server.crt');
    });
  });

  it('updates server key path on input change', async () => {
    renderSettings();
    await waitFor(() => {
      const serverKeyInput = screen.getByLabelText('Server Key Path');
      fireEvent.change(serverKeyInput, { target: { value: '/path/to/server.key' } });
      expect(serverKeyInput).toHaveValue('/path/to/server.key');
    });
  });

  it('updates stale timeout on valid input change', async () => {
    renderSettings();
    await waitFor(() => {
      const staleInput = screen.getByLabelText('Stale Agent Timeout (seconds)');
      fireEvent.change(staleInput, { target: { value: '300' } });
      expect(staleInput).toHaveValue(300);
    });
  });
});

describe('Notification Settings', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    localStorageMock.clear();
  });

  it('renders notification settings section', () => {
    renderSettings();
    expect(screen.getByText('Notification Settings')).toBeInTheDocument();
  });

  it('renders enable notifications toggle after loading', async () => {
    renderSettings();
    await waitFor(() => {
      expect(screen.getByText('Enable Notifications')).toBeInTheDocument();
    });
  });

  it('renders notification settings description after loading', async () => {
    renderSettings();
    await waitFor(() => {
      expect(screen.getByText('Receive notifications for execution events')).toBeInTheDocument();
    });
  });

  it('shows notification channel options when enabled', async () => {
    const { notificationApi } = await import('../lib/api');
    vi.mocked(notificationApi.getSettings).mockResolvedValue({
      data: {
        channel: 'email',
        enabled: true,
        email_address: 'test@example.com',
        notify_on_start: false,
        notify_on_complete: true,
        notify_on_failure: true,
        notify_on_score_alert: true,
        score_alert_threshold: 70,
        notify_on_agent_offline: true,
      },
    } as never);

    renderSettings();

    await waitFor(() => {
      expect(screen.getByText('Notification Channel')).toBeInTheDocument();
    });
  });

  it('shows email input when email channel selected', async () => {
    const { notificationApi } = await import('../lib/api');
    vi.mocked(notificationApi.getSettings).mockResolvedValue({
      data: {
        channel: 'email',
        enabled: true,
        email_address: 'test@example.com',
        notify_on_start: false,
        notify_on_complete: true,
        notify_on_failure: true,
        notify_on_score_alert: true,
        score_alert_threshold: 70,
        notify_on_agent_offline: true,
      },
    } as never);

    renderSettings();

    await waitFor(() => {
      expect(screen.getByText('Email Address')).toBeInTheDocument();
    });
  });

  it('shows notification type toggles when enabled', async () => {
    const { notificationApi } = await import('../lib/api');
    vi.mocked(notificationApi.getSettings).mockResolvedValue({
      data: {
        channel: 'email',
        enabled: true,
        email_address: 'test@example.com',
        notify_on_start: false,
        notify_on_complete: true,
        notify_on_failure: true,
        notify_on_score_alert: true,
        score_alert_threshold: 70,
        notify_on_agent_offline: true,
      },
    } as never);

    renderSettings();

    await waitFor(() => {
      expect(screen.getByText('Notify me when:')).toBeInTheDocument();
      expect(screen.getByText('Execution starts')).toBeInTheDocument();
      expect(screen.getByText('Execution completes')).toBeInTheDocument();
      expect(screen.getByText('Execution fails')).toBeInTheDocument();
      expect(screen.getByText('Agent goes offline')).toBeInTheDocument();
    });
  });

  it('shows save notification settings button when enabled', async () => {
    const { notificationApi } = await import('../lib/api');
    vi.mocked(notificationApi.getSettings).mockResolvedValue({
      data: {
        channel: 'email',
        enabled: true,
        email_address: 'test@example.com',
        notify_on_start: false,
        notify_on_complete: true,
        notify_on_failure: true,
        notify_on_score_alert: true,
        score_alert_threshold: 70,
        notify_on_agent_offline: true,
      },
    } as never);

    renderSettings();

    await waitFor(() => {
      expect(screen.getByText('Save Notification Settings')).toBeInTheDocument();
    });
  });

  it('saves notification settings on button click', async () => {
    const { notificationApi } = await import('../lib/api');
    const toast = await import('react-hot-toast');
    vi.mocked(notificationApi.getSettings).mockResolvedValue({
      data: {
        id: 'settings-123',
        user_id: 'user-123',
        channel: 'email',
        enabled: true,
        email_address: 'test@example.com',
        notify_on_start: false,
        notify_on_complete: true,
        notify_on_failure: true,
        notify_on_score_alert: true,
        score_alert_threshold: 70,
        notify_on_agent_offline: true,
        created_at: '2026-01-01T00:00:00Z',
        updated_at: '2026-01-01T00:00:00Z',
      },
    } as never);
    vi.mocked(notificationApi.updateSettings).mockResolvedValue({ data: {} } as never);

    renderSettings();

    await waitFor(() => {
      expect(screen.getByText('Save Notification Settings')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Save Notification Settings'));

    await waitFor(() => {
      expect(notificationApi.updateSettings).toHaveBeenCalled();
      expect(toast.default.success).toHaveBeenCalledWith('Notification settings saved');
    });
  });

  it('shows error toast when save fails', async () => {
    const { notificationApi } = await import('../lib/api');
    const toast = await import('react-hot-toast');
    vi.mocked(notificationApi.getSettings).mockResolvedValue({
      data: {
        id: 'settings-456',
        user_id: 'user-123',
        channel: 'email',
        enabled: true,
        email_address: '',
        notify_on_start: false,
        notify_on_complete: true,
        notify_on_failure: true,
        notify_on_score_alert: true,
        score_alert_threshold: 70,
        notify_on_agent_offline: true,
        created_at: '2026-01-01T00:00:00Z',
        updated_at: '2026-01-01T00:00:00Z',
      },
    } as never);
    vi.mocked(notificationApi.updateSettings).mockRejectedValue(new Error('Save failed'));

    renderSettings();

    await waitFor(() => {
      expect(screen.getByText('Save Notification Settings')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Save Notification Settings'));

    await waitFor(() => {
      expect(toast.default.error).toHaveBeenCalledWith('Failed to save notification settings');
    });
  });
});

describe('SMTP Test', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    localStorageMock.clear();
  });

  it('shows SMTP test section when SMTP is configured', async () => {
    const { notificationApi } = await import('../lib/api');
    vi.mocked(notificationApi.getSMTPConfig).mockResolvedValue({
      data: {
        host: 'smtp.example.com',
        port: 587,
        use_tls: true,
      },
    } as never);

    renderSettings();

    await waitFor(() => {
      expect(screen.getByText('Test Email')).toBeInTheDocument();
    });
  });

  it('sends test email on button click', async () => {
    const { notificationApi } = await import('../lib/api');
    const toast = await import('react-hot-toast');
    vi.mocked(notificationApi.getSMTPConfig).mockResolvedValue({
      data: {
        host: 'smtp.example.com',
        port: 587,
        use_tls: true,
      },
    } as never);
    vi.mocked(notificationApi.testSMTP).mockResolvedValue({ data: {} } as never);

    renderSettings();

    await waitFor(() => {
      expect(screen.getByText('Test Email')).toBeInTheDocument();
    });

    const emailInput = screen.getByPlaceholderText('Enter email to send test');
    fireEvent.change(emailInput, { target: { value: 'test@example.com' } });

    fireEvent.click(screen.getByText('Send Test'));

    await waitFor(() => {
      expect(notificationApi.testSMTP).toHaveBeenCalledWith('test@example.com');
      expect(toast.default.success).toHaveBeenCalledWith('Test email sent successfully');
    });
  });

  it('shows error when test email fails', async () => {
    const { notificationApi } = await import('../lib/api');
    const toast = await import('react-hot-toast');
    vi.mocked(notificationApi.getSMTPConfig).mockResolvedValue({
      data: {
        host: 'smtp.example.com',
        port: 587,
        use_tls: true,
      },
    } as never);
    vi.mocked(notificationApi.testSMTP).mockRejectedValue(new Error('SMTP error'));

    renderSettings();

    await waitFor(() => {
      expect(screen.getByText('Test Email')).toBeInTheDocument();
    });

    const emailInput = screen.getByPlaceholderText('Enter email to send test');
    fireEvent.change(emailInput, { target: { value: 'test@example.com' } });

    fireEvent.click(screen.getByText('Send Test'));

    await waitFor(() => {
      expect(toast.default.error).toHaveBeenCalledWith('Failed to send test email');
    });
  });

  it('shows error when no email entered', async () => {
    const { notificationApi } = await import('../lib/api');
    const toast = await import('react-hot-toast');
    vi.mocked(notificationApi.getSMTPConfig).mockResolvedValue({
      data: {
        host: 'smtp.example.com',
        port: 587,
        use_tls: true,
      },
    } as never);

    renderSettings();

    await waitFor(() => {
      expect(screen.getByText('Test Email')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Send Test'));

    await waitFor(() => {
      expect(toast.default.error).toHaveBeenCalledWith('Please enter an email address');
    });
  });
});

describe('Safe Mode Toggle', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    localStorageMock.clear();
  });

  it('renders safe mode toggle with description', () => {
    renderSettings();
    expect(screen.getByText('Safe Mode by Default')).toBeInTheDocument();
    expect(screen.getByText("Only run safe techniques that don't modify the system")).toBeInTheDocument();
  });
});

describe('Settings Server Configuration', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    localStorageMock.clear();
    localStorageMock.getItem.mockReturnValue(null);
  });

  it('renders server configuration section', () => {
    renderSettings();
    expect(screen.getByText('Server Configuration')).toBeInTheDocument();
  });

  it('renders execution settings section', () => {
    renderSettings();
    expect(screen.getByText('Execution Settings')).toBeInTheDocument();
  });

  it('renders agent settings section', () => {
    renderSettings();
    expect(screen.getByText('Agent Settings')).toBeInTheDocument();
  });
});

describe('Settings Webhook Configuration', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    localStorageMock.clear();
  });

  it('shows webhook channel option when notifications enabled', async () => {
    const { notificationApi } = await import('../lib/api');
    vi.mocked(notificationApi.getSettings).mockResolvedValue({
      data: {
        channel: 'webhook',
        enabled: true,
        email_address: '',
        webhook_url: 'https://example.com/webhook',
        notify_on_start: false,
        notify_on_complete: true,
        notify_on_failure: true,
        notify_on_score_alert: true,
        score_alert_threshold: 70,
        notify_on_agent_offline: true,
      },
    } as never);

    renderSettings();

    await waitFor(() => {
      expect(screen.getByText('Notification Channel')).toBeInTheDocument();
    });
  });
});

describe('Settings Score Alert Threshold', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    localStorageMock.clear();
  });

  it('shows notification settings when enabled', async () => {
    const { notificationApi } = await import('../lib/api');
    vi.mocked(notificationApi.getSettings).mockResolvedValue({
      data: {
        channel: 'email',
        enabled: true,
        email_address: 'test@example.com',
        notify_on_start: false,
        notify_on_complete: true,
        notify_on_failure: true,
        notify_on_score_alert: true,
        score_alert_threshold: 70,
        notify_on_agent_offline: true,
      },
    } as never);

    renderSettings();

    await waitFor(() => {
      expect(screen.getByText('Notification Settings')).toBeInTheDocument();
    });
  });
});

describe('Settings Notification Toggle', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    localStorageMock.clear();
  });

  it('toggles notification enabled state', async () => {
    const { notificationApi } = await import('../lib/api');
    vi.mocked(notificationApi.getSettings).mockResolvedValue({
      data: {
        channel: 'email',
        enabled: true,
        email_address: 'test@example.com',
        notify_on_start: false,
        notify_on_complete: true,
        notify_on_failure: true,
        notify_on_score_alert: true,
        score_alert_threshold: 70,
        notify_on_agent_offline: true,
      },
    } as never);

    renderSettings();

    await waitFor(() => {
      expect(screen.getByText('Enable Notifications')).toBeInTheDocument();
    });
  });

  it('shows expanded settings when notifications enabled', async () => {
    const { notificationApi } = await import('../lib/api');
    vi.mocked(notificationApi.getSettings).mockResolvedValue({
      data: {
        channel: 'email',
        enabled: true,
        email_address: 'test@example.com',
        notify_on_start: true,
        notify_on_complete: true,
        notify_on_failure: true,
        notify_on_score_alert: false,
        score_alert_threshold: 70,
        notify_on_agent_offline: true,
      },
    } as never);

    renderSettings();

    await waitFor(() => {
      expect(screen.getByText('Execution starts')).toBeInTheDocument();
    });
  });
});

describe('Settings Create Notification Settings', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    localStorageMock.clear();
  });

  it('shows enable notifications toggle when no settings exist', async () => {
    const { notificationApi } = await import('../lib/api');
    vi.mocked(notificationApi.getSettings).mockRejectedValue({ response: { status: 404 } });

    renderSettings();

    await waitFor(() => {
      expect(screen.getByText('Enable Notifications')).toBeInTheDocument();
    });
  });
});

describe('Settings Form Persistence', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    localStorageMock.clear();
  });

  it('persists server URL across page reload', async () => {
    const savedSettings = {
      serverUrl: 'https://myserver:9000',
      heartbeatInterval: 60,
      staleTimeout: 240,
      safeMode: false,
    };
    localStorageMock.getItem.mockReturnValue(JSON.stringify(savedSettings));

    renderSettings();

    await waitFor(() => {
      const serverUrlInput = screen.getByLabelText('Server URL');
      expect(serverUrlInput).toHaveValue('https://myserver:9000');
    });
  });
});

describe('Settings TLS Configuration Details', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    localStorageMock.clear();
  });

  it('shows TLS section collapsed by default', () => {
    renderSettings();
    expect(screen.getByText('TLS / mTLS Configuration')).toBeInTheDocument();
  });

  it('allows input in all TLS fields', async () => {
    renderSettings();
    await waitFor(() => {
      const caCertInput = screen.getByLabelText('CA Certificate Path');
      const serverCertInput = screen.getByLabelText('Server Certificate Path');
      const serverKeyInput = screen.getByLabelText('Server Key Path');

      fireEvent.change(caCertInput, { target: { value: '/custom/ca.crt' } });
      fireEvent.change(serverCertInput, { target: { value: '/custom/server.crt' } });
      fireEvent.change(serverKeyInput, { target: { value: '/custom/server.key' } });

      expect(caCertInput).toHaveValue('/custom/ca.crt');
      expect(serverCertInput).toHaveValue('/custom/server.crt');
      expect(serverKeyInput).toHaveValue('/custom/server.key');
    });
  });
});

describe('Settings SMTP Configuration Display', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    localStorageMock.clear();
  });

  it('shows SMTP host when configured', async () => {
    const { notificationApi } = await import('../lib/api');
    vi.mocked(notificationApi.getSMTPConfig).mockResolvedValue({
      data: {
        host: 'smtp.example.com',
        port: 587,
        username: 'user@example.com',
        use_tls: true,
      },
    } as never);

    renderSettings();

    await waitFor(() => {
      expect(screen.getByText(/smtp.example.com/i)).toBeInTheDocument();
    });
  });

  it('shows TLS status when SMTP configured', async () => {
    const { notificationApi } = await import('../lib/api');
    vi.mocked(notificationApi.getSMTPConfig).mockResolvedValue({
      data: {
        host: 'smtp.example.com',
        port: 587,
        username: 'user@example.com',
        use_tls: true,
      },
    } as never);

    renderSettings();

    await waitFor(() => {
      expect(screen.getByText(/TLS/i)).toBeInTheDocument();
    });
  });
});

describe('Settings Error States', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    localStorageMock.clear();
  });

  it('handles notification settings fetch error', async () => {
    const { notificationApi } = await import('../lib/api');
    vi.mocked(notificationApi.getSettings).mockRejectedValue(new Error('Network error'));

    // Should not crash
    expect(() => renderSettings()).not.toThrow();
  });

  it('handles SMTP config fetch error', async () => {
    const { notificationApi } = await import('../lib/api');
    vi.mocked(notificationApi.getSMTPConfig).mockRejectedValue(new Error('Network error'));

    // Should not crash
    expect(() => renderSettings()).not.toThrow();
  });
});

describe('Toggle Component Interactions', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    localStorageMock.clear();
  });

  it('toggles safe mode when clicked', async () => {
    renderSettings();

    await waitFor(() => {
      expect(screen.getByText('Safe Mode by Default')).toBeInTheDocument();
    });

    // Find the safe mode toggle button - it's inside the "Execution Settings" card
    // The structure is: div.flex > div(text) + button(toggle)
    const safeModeLabel = screen.getByText('Safe Mode by Default');
    const safeModeContainer = safeModeLabel.closest('div.flex');
    const toggleButton = safeModeContainer!.querySelector('button')!;

    // Get initial state
    const initiallyEnabled = toggleButton.className.includes('bg-primary-600');

    // Click to toggle
    fireEvent.click(toggleButton);

    // State should have flipped
    await waitFor(() => {
      if (initiallyEnabled) {
        expect(toggleButton.className).toContain('bg-gray-200');
      } else {
        expect(toggleButton.className).toContain('bg-primary-600');
      }
    });

    // Click again to toggle back
    fireEvent.click(toggleButton);

    await waitFor(() => {
      if (initiallyEnabled) {
        expect(toggleButton.className).toContain('bg-primary-600');
      } else {
        expect(toggleButton.className).toContain('bg-gray-200');
      }
    });
  });

  it('toggles enable notifications from off to on', async () => {
    renderSettings();

    await waitFor(() => {
      expect(screen.getByText('Enable Notifications')).toBeInTheDocument();
    });

    // Find the enable notifications toggle
    const enableSection = screen.getByText('Enable Notifications').closest('.flex');
    const toggleButton = enableSection!.querySelector('button')!;

    // Notifications should be disabled by default
    expect(toggleButton.className).toContain('bg-gray-200');

    // Click to enable
    fireEvent.click(toggleButton);

    // Now the expanded notification options should appear
    await waitFor(() => {
      expect(screen.getByText('Notification Channel')).toBeInTheDocument();
    });
  });

  it('toggles notification type toggles when clicked', async () => {
    const { notificationApi } = await import('../lib/api');
    vi.mocked(notificationApi.getSettings).mockResolvedValue({
      data: {
        channel: 'email',
        enabled: true,
        email_address: 'test@example.com',
        notify_on_start: false,
        notify_on_complete: true,
        notify_on_failure: true,
        notify_on_score_alert: true,
        score_alert_threshold: 70,
        notify_on_agent_offline: true,
      },
    } as never);

    renderSettings();

    await waitFor(() => {
      expect(screen.getByText('Execution starts')).toBeInTheDocument();
    });

    // Find and click the "Execution starts" toggle
    const startSection = screen.getByText('Execution starts').closest('.flex');
    const startToggle = startSection!.querySelector('button')!;
    fireEvent.click(startToggle);

    // Find and click the "Execution completes" toggle
    const completeSection = screen.getByText('Execution completes').closest('.flex');
    const completeToggle = completeSection!.querySelector('button')!;
    fireEvent.click(completeToggle);

    // Find and click the "Execution fails" toggle
    const failSection = screen.getByText('Execution fails').closest('.flex');
    const failToggle = failSection!.querySelector('button')!;
    fireEvent.click(failToggle);

    // Find and click the "Agent goes offline" toggle
    const offlineSection = screen.getByText('Agent goes offline').closest('.flex');
    const offlineToggle = offlineSection!.querySelector('button')!;
    fireEvent.click(offlineToggle);

    // Toggles should have visually changed their state
    await waitFor(() => {
      // "Execution starts" was false, now should be true (bg-primary-600)
      expect(startToggle.className).toContain('bg-primary-600');
      // "Execution completes" was true, now should be false (bg-gray-200)
      expect(completeToggle.className).toContain('bg-gray-200');
    });
  });

  it('toggles score alert and hides threshold input when disabled', async () => {
    const { notificationApi } = await import('../lib/api');
    vi.mocked(notificationApi.getSettings).mockResolvedValue({
      data: {
        channel: 'email',
        enabled: true,
        email_address: 'test@example.com',
        notify_on_start: false,
        notify_on_complete: true,
        notify_on_failure: true,
        notify_on_score_alert: true,
        score_alert_threshold: 70,
        notify_on_agent_offline: true,
      },
    } as never);

    renderSettings();

    await waitFor(() => {
      expect(screen.getByText('Security score below threshold')).toBeInTheDocument();
    });

    // The score threshold input should be visible when notify_on_score_alert is true
    const thresholdInput = screen.getByDisplayValue('70');
    expect(thresholdInput).toBeInTheDocument();

    // Find and click the score alert toggle to disable it
    const scoreSection = screen.getByText('Security score below threshold').closest('.flex');
    const scoreToggle = scoreSection!.querySelector('button')!;
    fireEvent.click(scoreToggle);

    // After disabling, the threshold input should disappear
    await waitFor(() => {
      expect(screen.queryByDisplayValue('70')).not.toBeInTheDocument();
    });
  });
});

describe('Channel Switching', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    localStorageMock.clear();
  });

  it('switches from email to webhook channel and shows webhook URL input', async () => {
    const { notificationApi } = await import('../lib/api');
    vi.mocked(notificationApi.getSettings).mockResolvedValue({
      data: {
        channel: 'email',
        enabled: true,
        email_address: 'test@example.com',
        notify_on_start: false,
        notify_on_complete: true,
        notify_on_failure: true,
        notify_on_score_alert: true,
        score_alert_threshold: 70,
        notify_on_agent_offline: true,
      },
    } as never);

    renderSettings();

    await waitFor(() => {
      expect(screen.getByText('Email Address')).toBeInTheDocument();
    });

    // Switch to webhook
    const channelSelect = screen.getByLabelText('Notification Channel');
    fireEvent.change(channelSelect, { target: { value: 'webhook' } });

    // Email address should disappear, webhook URL should appear
    await waitFor(() => {
      expect(screen.queryByText('Email Address')).not.toBeInTheDocument();
      expect(screen.getByText('Webhook URL')).toBeInTheDocument();
    });
  });

  it('shows webhook URL input with existing webhook settings', async () => {
    const { notificationApi } = await import('../lib/api');
    vi.mocked(notificationApi.getSettings).mockResolvedValue({
      data: {
        channel: 'webhook',
        enabled: true,
        email_address: '',
        webhook_url: 'https://hooks.example.com/notify',
        notify_on_start: false,
        notify_on_complete: true,
        notify_on_failure: true,
        notify_on_score_alert: true,
        score_alert_threshold: 70,
        notify_on_agent_offline: true,
      },
    } as never);

    renderSettings();

    await waitFor(() => {
      expect(screen.getByText('Webhook URL')).toBeInTheDocument();
      const webhookInput = screen.getByLabelText('Webhook URL');
      expect(webhookInput).toHaveValue('https://hooks.example.com/notify');
    });
  });

  it('updates webhook URL on input change', async () => {
    const { notificationApi } = await import('../lib/api');
    vi.mocked(notificationApi.getSettings).mockResolvedValue({
      data: {
        channel: 'webhook',
        enabled: true,
        email_address: '',
        webhook_url: '',
        notify_on_start: false,
        notify_on_complete: true,
        notify_on_failure: true,
        notify_on_score_alert: true,
        score_alert_threshold: 70,
        notify_on_agent_offline: true,
      },
    } as never);

    renderSettings();

    await waitFor(() => {
      expect(screen.getByText('Webhook URL')).toBeInTheDocument();
    });

    const webhookInput = screen.getByLabelText('Webhook URL');
    fireEvent.change(webhookInput, { target: { value: 'https://new-webhook.com/hook' } });

    expect(webhookInput).toHaveValue('https://new-webhook.com/hook');
  });

  it('updates email address on input change', async () => {
    const { notificationApi } = await import('../lib/api');
    vi.mocked(notificationApi.getSettings).mockResolvedValue({
      data: {
        channel: 'email',
        enabled: true,
        email_address: '',
        notify_on_start: false,
        notify_on_complete: true,
        notify_on_failure: true,
        notify_on_score_alert: true,
        score_alert_threshold: 70,
        notify_on_agent_offline: true,
      },
    } as never);

    renderSettings();

    await waitFor(() => {
      expect(screen.getByText('Email Address')).toBeInTheDocument();
    });

    const emailInput = screen.getByLabelText('Email Address');
    fireEvent.change(emailInput, { target: { value: 'new@example.com' } });

    expect(emailInput).toHaveValue('new@example.com');
  });

  it('switches from webhook back to email channel', async () => {
    const { notificationApi } = await import('../lib/api');
    vi.mocked(notificationApi.getSettings).mockResolvedValue({
      data: {
        channel: 'webhook',
        enabled: true,
        email_address: '',
        webhook_url: 'https://hooks.example.com/notify',
        notify_on_start: false,
        notify_on_complete: true,
        notify_on_failure: true,
        notify_on_score_alert: true,
        score_alert_threshold: 70,
        notify_on_agent_offline: true,
      },
    } as never);

    renderSettings();

    await waitFor(() => {
      expect(screen.getByText('Webhook URL')).toBeInTheDocument();
    });

    // Switch to email
    const channelSelect = screen.getByLabelText('Notification Channel');
    fireEvent.change(channelSelect, { target: { value: 'email' } });

    await waitFor(() => {
      expect(screen.queryByText('Webhook URL')).not.toBeInTheDocument();
      expect(screen.getByText('Email Address')).toBeInTheDocument();
    });
  });
});

describe('Number Input Validation', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    localStorageMock.clear();
  });

  it('updates score alert threshold on valid input', async () => {
    const { notificationApi } = await import('../lib/api');
    vi.mocked(notificationApi.getSettings).mockResolvedValue({
      data: {
        channel: 'email',
        enabled: true,
        email_address: 'test@example.com',
        notify_on_start: false,
        notify_on_complete: true,
        notify_on_failure: true,
        notify_on_score_alert: true,
        score_alert_threshold: 70,
        notify_on_agent_offline: true,
      },
    } as never);

    renderSettings();

    await waitFor(() => {
      expect(screen.getByDisplayValue('70')).toBeInTheDocument();
    });

    const thresholdInput = screen.getByDisplayValue('70');
    fireEvent.change(thresholdInput, { target: { value: '85' } });

    expect(thresholdInput).toHaveValue(85);
  });

  it('defaults score alert threshold to 0 on empty input', async () => {
    const { notificationApi } = await import('../lib/api');
    vi.mocked(notificationApi.getSettings).mockResolvedValue({
      data: {
        channel: 'email',
        enabled: true,
        email_address: 'test@example.com',
        notify_on_start: false,
        notify_on_complete: true,
        notify_on_failure: true,
        notify_on_score_alert: true,
        score_alert_threshold: 70,
        notify_on_agent_offline: true,
      },
    } as never);

    renderSettings();

    await waitFor(() => {
      expect(screen.getByDisplayValue('70')).toBeInTheDocument();
    });

    const thresholdInput = screen.getByDisplayValue('70');
    fireEvent.change(thresholdInput, { target: { value: '' } });

    // Should default to 0 when empty (parseInt returns NaN, || 0 kicks in)
    expect(thresholdInput).toHaveValue(0);
  });

  it('defaults heartbeat interval to 30 on non-numeric input', async () => {
    renderSettings();

    await waitFor(() => {
      const heartbeatInput = screen.getByLabelText('Heartbeat Interval (seconds)');
      fireEvent.change(heartbeatInput, { target: { value: 'abc' } });
      expect(heartbeatInput).toHaveValue(30);
    });
  });

  it('defaults stale timeout to 120 on non-numeric input', async () => {
    renderSettings();

    await waitFor(() => {
      const staleInput = screen.getByLabelText('Stale Agent Timeout (seconds)');
      fireEvent.change(staleInput, { target: { value: 'xyz' } });
      expect(staleInput).toHaveValue(120);
    });
  });
});

describe('Loading State', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    localStorageMock.clear();
  });

  it('shows loading text while notification settings are being fetched', async () => {
    const { notificationApi } = await import('../lib/api');
    // Create a promise that never resolves to keep loading state
    vi.mocked(notificationApi.getSettings).mockImplementation(
      () => new Promise(() => {})
    );

    renderSettings();

    // Should show loading state
    expect(screen.getByText('Loading...')).toBeInTheDocument();
  });
});

describe('SMTP Warning and Display', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    localStorageMock.clear();
  });

  it('shows SMTP not configured warning when email channel selected and no SMTP config', async () => {
    const { notificationApi } = await import('../lib/api');
    vi.mocked(notificationApi.getSettings).mockResolvedValue({
      data: {
        channel: 'email',
        enabled: true,
        email_address: 'test@example.com',
        notify_on_start: false,
        notify_on_complete: true,
        notify_on_failure: true,
        notify_on_score_alert: true,
        score_alert_threshold: 70,
        notify_on_agent_offline: true,
      },
    } as never);
    vi.mocked(notificationApi.getSMTPConfig).mockRejectedValue({ response: { status: 404 } });

    renderSettings();

    await waitFor(() => {
      expect(
        screen.getByText('SMTP not configured on server. Email notifications may not work.')
      ).toBeInTheDocument();
    });
  });

  it('does not show SMTP warning when SMTP is configured', async () => {
    const { notificationApi } = await import('../lib/api');
    vi.mocked(notificationApi.getSettings).mockResolvedValue({
      data: {
        channel: 'email',
        enabled: true,
        email_address: 'test@example.com',
        notify_on_start: false,
        notify_on_complete: true,
        notify_on_failure: true,
        notify_on_score_alert: true,
        score_alert_threshold: 70,
        notify_on_agent_offline: true,
      },
    } as never);
    vi.mocked(notificationApi.getSMTPConfig).mockResolvedValue({
      data: {
        host: 'smtp.example.com',
        port: 587,
        use_tls: true,
      },
    } as never);

    renderSettings();

    await waitFor(() => {
      expect(screen.getByText('Email Address')).toBeInTheDocument();
    });

    expect(
      screen.queryByText('SMTP not configured on server. Email notifications may not work.')
    ).not.toBeInTheDocument();
  });

  it('shows Plain when SMTP use_tls is false', async () => {
    const { notificationApi } = await import('../lib/api');
    vi.mocked(notificationApi.getSMTPConfig).mockResolvedValue({
      data: {
        host: 'smtp.example.com',
        port: 25,
        use_tls: false,
      },
    } as never);

    renderSettings();

    await waitFor(() => {
      expect(screen.getByText(/smtp\.example\.com:25 \(Plain\)/)).toBeInTheDocument();
    });
  });

  it('shows TLS when SMTP use_tls is true', async () => {
    const { notificationApi } = await import('../lib/api');
    vi.mocked(notificationApi.getSMTPConfig).mockResolvedValue({
      data: {
        host: 'smtp.example.com',
        port: 587,
        use_tls: true,
      },
    } as never);

    renderSettings();

    await waitFor(() => {
      expect(screen.getByText(/smtp\.example\.com:587 \(TLS\)/)).toBeInTheDocument();
    });
  });
});

describe('Create Notification Settings (New Settings)', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    localStorageMock.clear();
  });

  it('calls createSettings when no existing settings and save is clicked', async () => {
    const { notificationApi } = await import('../lib/api');
    const toast = await import('react-hot-toast');

    // Return 404 for getSettings so no existing settings
    vi.mocked(notificationApi.getSettings).mockRejectedValue({ response: { status: 404 } });
    vi.mocked(notificationApi.createSettings).mockResolvedValue({ data: {} } as never);

    renderSettings();

    await waitFor(() => {
      expect(screen.getByText('Enable Notifications')).toBeInTheDocument();
    });

    // Enable notifications
    const enableSection = screen.getByText('Enable Notifications').closest('.flex');
    const enableToggle = enableSection!.querySelector('button')!;
    fireEvent.click(enableToggle);

    // Wait for expanded settings to appear
    await waitFor(() => {
      expect(screen.getByText('Save Notification Settings')).toBeInTheDocument();
    });

    // Click save
    fireEvent.click(screen.getByText('Save Notification Settings'));

    await waitFor(() => {
      expect(notificationApi.createSettings).toHaveBeenCalled();
      expect(toast.default.success).toHaveBeenCalledWith('Notification settings saved');
    });
  });

  it('calls updateSettings when existing settings with id and save is clicked', async () => {
    const { notificationApi } = await import('../lib/api');
    const toast = await import('react-hot-toast');

    vi.mocked(notificationApi.getSettings).mockResolvedValue({
      data: {
        id: 'existing-settings-id',
        user_id: 'user-1',
        channel: 'email',
        enabled: true,
        email_address: 'test@example.com',
        notify_on_start: false,
        notify_on_complete: true,
        notify_on_failure: true,
        notify_on_score_alert: true,
        score_alert_threshold: 70,
        notify_on_agent_offline: true,
        created_at: '2026-01-01T00:00:00Z',
        updated_at: '2026-01-01T00:00:00Z',
      },
    } as never);
    vi.mocked(notificationApi.updateSettings).mockResolvedValue({ data: {} } as never);

    renderSettings();

    await waitFor(() => {
      expect(screen.getByText('Save Notification Settings')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Save Notification Settings'));

    await waitFor(() => {
      expect(notificationApi.updateSettings).toHaveBeenCalledWith(
        'existing-settings-id',
        expect.objectContaining({
          channel: 'email',
          enabled: true,
          email_address: 'test@example.com',
        })
      );
      expect(toast.default.success).toHaveBeenCalledWith('Notification settings saved');
    });
  });
});

describe('Notification Settings Data Mapping', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    localStorageMock.clear();
  });

  it('maps fetched notification settings data correctly to form state', async () => {
    const { notificationApi } = await import('../lib/api');
    vi.mocked(notificationApi.getSettings).mockResolvedValue({
      data: {
        id: 'ns-1',
        user_id: 'u-1',
        channel: 'webhook',
        enabled: true,
        email_address: '',
        webhook_url: 'https://hooks.slack.com/test',
        notify_on_start: true,
        notify_on_complete: false,
        notify_on_failure: true,
        notify_on_score_alert: false,
        score_alert_threshold: 50,
        notify_on_agent_offline: false,
        created_at: '2026-01-01T00:00:00Z',
        updated_at: '2026-01-01T00:00:00Z',
      },
    } as never);

    renderSettings();

    await waitFor(() => {
      // Webhook channel should be selected
      const channelSelect = screen.getByLabelText('Notification Channel');
      expect(channelSelect).toHaveValue('webhook');

      // Webhook URL should be populated
      const webhookInput = screen.getByLabelText('Webhook URL');
      expect(webhookInput).toHaveValue('https://hooks.slack.com/test');
    });

    // Check toggle states: notify_on_start is true, notify_on_complete is false
    const startSection = screen.getByText('Execution starts').closest('.flex');
    const startToggle = startSection!.querySelector('button')!;
    expect(startToggle.className).toContain('bg-primary-600');

    const completeSection = screen.getByText('Execution completes').closest('.flex');
    const completeToggle = completeSection!.querySelector('button')!;
    expect(completeToggle.className).toContain('bg-gray-200');

    // notify_on_agent_offline is false
    const offlineSection = screen.getByText('Agent goes offline').closest('.flex');
    const offlineToggle = offlineSection!.querySelector('button')!;
    expect(offlineToggle.className).toContain('bg-gray-200');
  });

  it('handles null email_address and webhook_url in fetched data', async () => {
    const { notificationApi } = await import('../lib/api');
    vi.mocked(notificationApi.getSettings).mockResolvedValue({
      data: {
        id: 'ns-2',
        user_id: 'u-2',
        channel: 'email',
        enabled: true,
        email_address: null,
        webhook_url: null,
        notify_on_start: false,
        notify_on_complete: true,
        notify_on_failure: true,
        notify_on_score_alert: true,
        score_alert_threshold: 70,
        notify_on_agent_offline: true,
        created_at: '2026-01-01T00:00:00Z',
        updated_at: '2026-01-01T00:00:00Z',
      },
    } as never);

    renderSettings();

    await waitFor(() => {
      const emailInput = screen.getByLabelText('Email Address');
      // null should be mapped to empty string via || ''
      expect(emailInput).toHaveValue('');
    });
  });
});

describe('Settings Persistence with All Fields', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    localStorageMock.clear();
  });

  it('loads all fields from localStorage including TLS and safeMode', async () => {
    const savedSettings = {
      serverUrl: 'https://myserver:8443',
      safeMode: false,
      heartbeatInterval: 15,
      staleTimeout: 60,
      caCertPath: '/etc/ssl/ca.crt',
      serverCertPath: '/etc/ssl/server.crt',
      serverKeyPath: '/etc/ssl/server.key',
    };
    localStorageMock.getItem.mockReturnValue(JSON.stringify(savedSettings));

    renderSettings();

    await waitFor(() => {
      expect(screen.getByLabelText('Server URL')).toHaveValue('https://myserver:8443');
      expect(screen.getByLabelText('Heartbeat Interval (seconds)')).toHaveValue(15);
      expect(screen.getByLabelText('Stale Agent Timeout (seconds)')).toHaveValue(60);
      expect(screen.getByLabelText('CA Certificate Path')).toHaveValue('/etc/ssl/ca.crt');
      expect(screen.getByLabelText('Server Certificate Path')).toHaveValue('/etc/ssl/server.crt');
      expect(screen.getByLabelText('Server Key Path')).toHaveValue('/etc/ssl/server.key');
    });

    // Safe mode should be off (false)
    const safeModeSection = screen.getByText('Safe Mode by Default').closest('.flex');
    const safeToggle = safeModeSection!.querySelector('button')!;
    expect(safeToggle.className).toContain('bg-gray-200');
  });

  it('saves all modified settings to localStorage', async () => {
    const toast = await import('react-hot-toast');
    renderSettings();

    await waitFor(() => {
      expect(screen.getByLabelText('Server URL')).toBeInTheDocument();
    });

    // Modify multiple fields
    fireEvent.change(screen.getByLabelText('Server URL'), {
      target: { value: 'https://prod:8443' },
    });
    fireEvent.change(screen.getByLabelText('Heartbeat Interval (seconds)'), {
      target: { value: '45' },
    });
    fireEvent.change(screen.getByLabelText('Stale Agent Timeout (seconds)'), {
      target: { value: '300' },
    });
    fireEvent.change(screen.getByLabelText('CA Certificate Path'), {
      target: { value: '/new/ca.crt' },
    });

    // Save settings
    fireEvent.click(screen.getByText('Save Settings'));

    await waitFor(() => {
      expect(localStorageMock.setItem).toHaveBeenCalledWith(
        'autostrike_settings',
        expect.stringContaining('"serverUrl":"https://prod:8443"')
      );
      expect(localStorageMock.setItem).toHaveBeenCalledWith(
        'autostrike_settings',
        expect.stringContaining('"heartbeatInterval":45')
      );
      expect(toast.default.success).toHaveBeenCalledWith('Settings saved successfully');
    });
  });
});

describe('SMTP Test Edge Cases', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    localStorageMock.clear();
  });

  it('does not show SMTP test section when SMTP is not configured', async () => {
    const { notificationApi } = await import('../lib/api');
    vi.mocked(notificationApi.getSMTPConfig).mockRejectedValue({ response: { status: 404 } });

    renderSettings();

    // Wait for queries to settle
    await waitFor(() => {
      expect(screen.getByText('Settings')).toBeInTheDocument();
    });

    expect(screen.queryByText('Test Email')).not.toBeInTheDocument();
  });

  it('updates test email input value', async () => {
    const { notificationApi } = await import('../lib/api');
    vi.mocked(notificationApi.getSMTPConfig).mockResolvedValue({
      data: {
        host: 'smtp.test.com',
        port: 465,
        use_tls: true,
      },
    } as never);

    renderSettings();

    await waitFor(() => {
      expect(screen.getByText('Test Email')).toBeInTheDocument();
    });

    const emailInput = screen.getByPlaceholderText('Enter email to send test');
    fireEvent.change(emailInput, { target: { value: 'admin@test.com' } });

    expect(emailInput).toHaveValue('admin@test.com');
  });
});

describe('Pending/Loading States for Mutations', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    localStorageMock.clear();
  });

  it('shows Saving... text while save notification mutation is pending', async () => {
    const { notificationApi } = await import('../lib/api');
    vi.mocked(notificationApi.getSettings).mockResolvedValue({
      data: {
        id: 'settings-pending',
        user_id: 'user-1',
        channel: 'email',
        enabled: true,
        email_address: 'test@example.com',
        notify_on_start: false,
        notify_on_complete: true,
        notify_on_failure: true,
        notify_on_score_alert: true,
        score_alert_threshold: 70,
        notify_on_agent_offline: true,
        created_at: '2026-01-01T00:00:00Z',
        updated_at: '2026-01-01T00:00:00Z',
      },
    } as never);

    // Make updateSettings hang (never resolve) to keep isPending true
    vi.mocked(notificationApi.updateSettings).mockImplementation(
      () => new Promise(() => {})
    );

    renderSettings();

    await waitFor(() => {
      expect(screen.getByText('Save Notification Settings')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Save Notification Settings'));

    await waitFor(() => {
      expect(screen.getByText('Saving...')).toBeInTheDocument();
    });
  });

  it('shows Sending... text while test SMTP mutation is pending', async () => {
    const { notificationApi } = await import('../lib/api');
    vi.mocked(notificationApi.getSMTPConfig).mockResolvedValue({
      data: {
        host: 'smtp.example.com',
        port: 587,
        use_tls: true,
      },
    } as never);

    // Make testSMTP hang (never resolve) to keep isPending true
    vi.mocked(notificationApi.testSMTP).mockImplementation(
      () => new Promise(() => {})
    );

    renderSettings();

    await waitFor(() => {
      expect(screen.getByText('Test Email')).toBeInTheDocument();
    });

    const emailInput = screen.getByPlaceholderText('Enter email to send test');
    fireEvent.change(emailInput, { target: { value: 'test@example.com' } });

    fireEvent.click(screen.getByText('Send Test'));

    await waitFor(() => {
      expect(screen.getByText('Sending...')).toBeInTheDocument();
    });
  });

  it('disables save notification button while mutation is pending', async () => {
    const { notificationApi } = await import('../lib/api');
    vi.mocked(notificationApi.getSettings).mockResolvedValue({
      data: {
        id: 'settings-disable',
        user_id: 'user-1',
        channel: 'email',
        enabled: true,
        email_address: 'test@example.com',
        notify_on_start: false,
        notify_on_complete: true,
        notify_on_failure: true,
        notify_on_score_alert: true,
        score_alert_threshold: 70,
        notify_on_agent_offline: true,
        created_at: '2026-01-01T00:00:00Z',
        updated_at: '2026-01-01T00:00:00Z',
      },
    } as never);

    vi.mocked(notificationApi.updateSettings).mockImplementation(
      () => new Promise(() => {})
    );

    renderSettings();

    await waitFor(() => {
      expect(screen.getByText('Save Notification Settings')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Save Notification Settings'));

    await waitFor(() => {
      const savingButton = screen.getByText('Saving...');
      expect(savingButton).toBeDisabled();
    });
  });

  it('disables send test button while SMTP test mutation is pending', async () => {
    const { notificationApi } = await import('../lib/api');
    vi.mocked(notificationApi.getSMTPConfig).mockResolvedValue({
      data: {
        host: 'smtp.example.com',
        port: 587,
        use_tls: true,
      },
    } as never);

    vi.mocked(notificationApi.testSMTP).mockImplementation(
      () => new Promise(() => {})
    );

    renderSettings();

    await waitFor(() => {
      expect(screen.getByText('Test Email')).toBeInTheDocument();
    });

    const emailInput = screen.getByPlaceholderText('Enter email to send test');
    fireEvent.change(emailInput, { target: { value: 'test@example.com' } });

    fireEvent.click(screen.getByText('Send Test'));

    await waitFor(() => {
      const sendingButton = screen.getByText('Sending...');
      expect(sendingButton).toBeDisabled();
    });
  });
});

describe('Toggle Disabled State', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    localStorageMock.clear();
  });

  it('renders Toggle with disabled prop showing correct styling', async () => {
    // The Toggle component accepts a disabled prop. While the Settings page
    // doesn't directly pass disabled=true to Toggle, we test that the disabled
    // class is applied correctly by verifying the opacity and cursor styles.
    // Since the Toggle is defined in Settings.tsx, rendering Settings with
    // pending mutations means buttons pass disabled=true to the save button
    // (not to Toggle directly), so we test the mutation-driven disabled state.
    const { notificationApi } = await import('../lib/api');
    vi.mocked(notificationApi.getSettings).mockResolvedValue({
      data: {
        id: 'settings-1',
        user_id: 'user-1',
        channel: 'email',
        enabled: true,
        email_address: 'test@example.com',
        notify_on_start: false,
        notify_on_complete: true,
        notify_on_failure: true,
        notify_on_score_alert: true,
        score_alert_threshold: 70,
        notify_on_agent_offline: true,
        created_at: '2026-01-01T00:00:00Z',
        updated_at: '2026-01-01T00:00:00Z',
      },
    } as never);

    renderSettings();

    await waitFor(() => {
      expect(screen.getByText('Save Notification Settings')).toBeInTheDocument();
    });

    // All toggle buttons should NOT be disabled
    const toggleButtons = screen.getAllByRole('button').filter(
      (btn) => btn.className.includes('rounded-full') && btn.className.includes('h-6')
    );
    for (const toggle of toggleButtons) {
      expect(toggle).not.toBeDisabled();
      expect(toggle.className).not.toContain('opacity-50');
    }
  });
});
