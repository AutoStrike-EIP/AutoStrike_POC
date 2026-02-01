import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import Settings from './Settings';

// Mock react-hot-toast
vi.mock('react-hot-toast', () => ({
  default: {
    success: vi.fn(),
    error: vi.fn(),
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

function renderSettings() {
  return render(
    <MemoryRouter>
      <Settings />
    </MemoryRouter>
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

  it('renders all configuration sections', () => {
    renderSettings();
    expect(screen.getByText('Server Configuration')).toBeInTheDocument();
    expect(screen.getByText('Execution Settings')).toBeInTheDocument();
    expect(screen.getByText('Agent Settings')).toBeInTheDocument();
    expect(screen.getByText('TLS / mTLS Configuration')).toBeInTheDocument();
  });

  it('renders server URL input with default value', () => {
    renderSettings();
    const inputs = screen.getAllByRole('textbox');
    // First text input should be server URL
    expect(inputs[0]).toHaveValue('https://localhost:8443');
  });

  it('renders heartbeat interval input with default value', () => {
    renderSettings();
    const spinbuttons = screen.getAllByRole('spinbutton');
    // First number input should be heartbeat interval
    expect(spinbuttons[0]).toHaveValue(30);
  });

  it('renders stale timeout input with default value', () => {
    renderSettings();
    const spinbuttons = screen.getAllByRole('spinbutton');
    // Second number input should be stale timeout
    expect(spinbuttons[1]).toHaveValue(120);
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

  it('updates server URL on input change', () => {
    renderSettings();
    const inputs = screen.getAllByRole('textbox');
    const serverUrlInput = inputs[0];

    fireEvent.change(serverUrlInput, { target: { value: 'https://newserver:8443' } });

    expect(serverUrlInput).toHaveValue('https://newserver:8443');
  });

  it('updates heartbeat interval on input change', () => {
    renderSettings();
    const spinbuttons = screen.getAllByRole('spinbutton');
    const heartbeatInput = spinbuttons[0];

    fireEvent.change(heartbeatInput, { target: { value: '60' } });

    expect(heartbeatInput).toHaveValue(60);
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

  it('loads settings from localStorage on mount', () => {
    const savedSettings = {
      serverUrl: 'https://custom:9999',
      heartbeatInterval: 45,
      staleTimeout: 180,
    };
    localStorageMock.getItem.mockReturnValue(JSON.stringify(savedSettings));

    renderSettings();

    const inputs = screen.getAllByRole('textbox');
    expect(inputs[0]).toHaveValue('https://custom:9999');
  });

  it('toggles safe mode when button clicked', () => {
    renderSettings();
    const toggleButton = screen.getByRole('button', { name: '' });

    // Initially enabled (has bg-primary-600)
    expect(toggleButton).toHaveClass('bg-primary-600');

    fireEvent.click(toggleButton);

    // Now disabled (has bg-gray-200)
    expect(toggleButton).toHaveClass('bg-gray-200');
  });
});

describe('Settings Default Values', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    localStorageMock.clear();
    localStorageMock.getItem.mockReturnValue(null);
  });

  it('uses default heartbeat interval when invalid input', () => {
    renderSettings();
    const spinbuttons = screen.getAllByRole('spinbutton');
    const heartbeatInput = spinbuttons[0];

    fireEvent.change(heartbeatInput, { target: { value: '' } });

    expect(heartbeatInput).toHaveValue(30);
  });

  it('uses default stale timeout when invalid input', () => {
    renderSettings();
    const spinbuttons = screen.getAllByRole('spinbutton');
    const staleInput = spinbuttons[1];

    fireEvent.change(staleInput, { target: { value: '' } });

    expect(staleInput).toHaveValue(120);
  });

  it('handles malformed JSON in localStorage gracefully', () => {
    localStorageMock.getItem.mockReturnValue('invalid-json');

    // Should not throw
    expect(() => renderSettings()).not.toThrow();

    // Should use defaults
    const inputs = screen.getAllByRole('textbox');
    expect(inputs[0]).toHaveValue('https://localhost:8443');
  });

  it('updates CA certificate path on input change', () => {
    renderSettings();
    const inputs = screen.getAllByRole('textbox');
    const caCertInput = inputs[1]; // Second text input

    fireEvent.change(caCertInput, { target: { value: '/path/to/ca.crt' } });

    expect(caCertInput).toHaveValue('/path/to/ca.crt');
  });

  it('updates server certificate path on input change', () => {
    renderSettings();
    const inputs = screen.getAllByRole('textbox');
    const serverCertInput = inputs[2]; // Third text input

    fireEvent.change(serverCertInput, { target: { value: '/path/to/server.crt' } });

    expect(serverCertInput).toHaveValue('/path/to/server.crt');
  });

  it('updates server key path on input change', () => {
    renderSettings();
    const inputs = screen.getAllByRole('textbox');
    const serverKeyInput = inputs[3]; // Fourth text input

    fireEvent.change(serverKeyInput, { target: { value: '/path/to/server.key' } });

    expect(serverKeyInput).toHaveValue('/path/to/server.key');
  });

  it('updates stale timeout on valid input change', () => {
    renderSettings();
    const spinbuttons = screen.getAllByRole('spinbutton');
    const staleInput = spinbuttons[1];

    fireEvent.change(staleInput, { target: { value: '300' } });

    expect(staleInput).toHaveValue(300);
  });
});
