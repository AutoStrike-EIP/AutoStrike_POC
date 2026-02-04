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
});
