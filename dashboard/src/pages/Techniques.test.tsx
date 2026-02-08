import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import Techniques from './Techniques';
import { api, techniqueApi } from '../lib/api';
import toast from 'react-hot-toast';

// Mock the API
vi.mock('../lib/api', () => ({
  api: {
    get: vi.fn(),
  },
  techniqueApi: {
    import: vi.fn(),
  },
}));

// Mock react-hot-toast
vi.mock('react-hot-toast', () => ({
  default: {
    success: vi.fn(),
    error: vi.fn(),
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
    // "Safety" column header and "Safe" badge
    expect(screen.getByText('Safety')).toBeInTheDocument();
    expect(screen.getByText('Safe')).toBeInTheDocument();
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

describe('Techniques Executor Expand', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  const mockTechWithExecutors = [
    {
      id: 'T1059.001',
      name: 'PowerShell',
      description: 'Adversaries may use PowerShell',
      tactic: 'execution',
      platforms: ['windows'],
      is_safe: false,
      executors: [
        { name: 'Mimikatz', type: 'psh', platform: 'windows', command: 'Invoke-Mimikatz', timeout: 120, elevation_required: true, is_safe: false },
        { name: 'Encoded', type: 'psh', platform: 'windows', command: 'powershell -enc base64', cleanup: 'Remove-Item $env:TEMP\\test.ps1', timeout: 60, is_safe: true },
      ],
    },
  ];

  it('expands technique row to show executors on click', async () => {
    vi.mocked(api.get).mockResolvedValue({ data: mockTechWithExecutors } as never);

    renderWithClient(<Techniques />);

    const idButton = await screen.findByText('T1059.001');
    fireEvent.click(idButton);

    expect(screen.getByText('Executors (2)')).toBeInTheDocument();
    expect(screen.getByText('Mimikatz')).toBeInTheDocument();
    expect(screen.getByText('Encoded')).toBeInTheDocument();
  });

  it('collapses technique row on second click', async () => {
    vi.mocked(api.get).mockResolvedValue({ data: mockTechWithExecutors } as never);

    renderWithClient(<Techniques />);

    const idButton = await screen.findByText('T1059.001');
    fireEvent.click(idButton);
    expect(screen.getByText('Executors (2)')).toBeInTheDocument();

    fireEvent.click(idButton);
    expect(screen.queryByText('Executors (2)')).not.toBeInTheDocument();
  });

  it('expands executor to show command details', async () => {
    vi.mocked(api.get).mockResolvedValue({ data: mockTechWithExecutors } as never);

    renderWithClient(<Techniques />);

    const idButton = await screen.findByText('T1059.001');
    fireEvent.click(idButton);

    // Command should not be visible yet
    expect(screen.queryByText('Invoke-Mimikatz')).not.toBeInTheDocument();

    // Click on Mimikatz executor
    const mimikatzButton = screen.getByText('Mimikatz').closest('button');
    fireEvent.click(mimikatzButton!);

    // Command and details should be visible
    expect(screen.getByText('Invoke-Mimikatz')).toBeInTheDocument();
    expect(screen.getByText('Command')).toBeInTheDocument();
    expect(screen.getByText('Timeout: 120s')).toBeInTheDocument();
  });

  it('shows elevation badge on executor requiring elevation', async () => {
    vi.mocked(api.get).mockResolvedValue({ data: mockTechWithExecutors } as never);

    renderWithClient(<Techniques />);

    const idButton = await screen.findByText('T1059.001');
    fireEvent.click(idButton);

    // Mimikatz has elevation_required: true
    expect(screen.getByText('Elevation')).toBeInTheDocument();
  });

  it('shows cleanup command in expanded executor', async () => {
    vi.mocked(api.get).mockResolvedValue({ data: mockTechWithExecutors } as never);

    renderWithClient(<Techniques />);

    const idButton = await screen.findByText('T1059.001');
    fireEvent.click(idButton);

    // Expand Encoded executor (has cleanup)
    const encodedButton = screen.getByText('Encoded').closest('button');
    fireEvent.click(encodedButton!);

    expect(screen.getByText('Cleanup')).toBeInTheDocument();
    expect(screen.getByText('Remove-Item $env:TEMP\\test.ps1')).toBeInTheDocument();
  });

  it('switches between executors showing only one at a time', async () => {
    vi.mocked(api.get).mockResolvedValue({ data: mockTechWithExecutors } as never);

    renderWithClient(<Techniques />);

    const idButton = await screen.findByText('T1059.001');
    fireEvent.click(idButton);

    // Expand Mimikatz
    const mimikatzButton = screen.getByText('Mimikatz').closest('button');
    fireEvent.click(mimikatzButton!);
    expect(screen.getByText('Invoke-Mimikatz')).toBeInTheDocument();

    // Expand Encoded - Mimikatz should collapse
    const encodedButton = screen.getByText('Encoded').closest('button');
    fireEvent.click(encodedButton!);
    expect(screen.queryByText('Invoke-Mimikatz')).not.toBeInTheDocument();
    expect(screen.getByText('powershell -enc base64')).toBeInTheDocument();
  });

  it('shows "No executors available" for technique without executors', async () => {
    const techNoExecutors = [
      {
        id: 'T1566',
        name: 'Phishing',
        description: 'Test',
        tactic: 'initial_access',
        platforms: ['windows'],
        is_safe: true,
        executors: [],
      },
    ];
    vi.mocked(api.get).mockResolvedValue({ data: techNoExecutors } as never);

    renderWithClient(<Techniques />);

    const idButton = await screen.findByText('T1566');
    fireEvent.click(idButton);

    expect(screen.getByText('No executors available')).toBeInTheDocument();
  });

  it('resets expanded executor when collapsing technique row', async () => {
    vi.mocked(api.get).mockResolvedValue({ data: mockTechWithExecutors } as never);

    renderWithClient(<Techniques />);

    const idButton = await screen.findByText('T1059.001');
    fireEvent.click(idButton);

    // Expand executor
    const mimikatzButton = screen.getByText('Mimikatz').closest('button');
    fireEvent.click(mimikatzButton!);
    expect(screen.getByText('Invoke-Mimikatz')).toBeInTheDocument();

    // Collapse technique row
    fireEvent.click(idButton);

    // Re-expand technique row - executor should not be expanded
    fireEvent.click(idButton);
    expect(screen.queryByText('Invoke-Mimikatz')).not.toBeInTheDocument();
  });
});

describe('Techniques Import Modal', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('opens import modal when Import Techniques button is clicked', async () => {
    vi.mocked(api.get).mockResolvedValue({ data: [] } as never);

    renderWithClient(<Techniques />);

    await screen.findByText('Techniques');
    fireEvent.click(screen.getByText('Import Techniques'));

    // Modal title is a heading
    expect(screen.getByRole('heading', { name: 'Import Techniques' })).toBeInTheDocument();
    expect(screen.getByText(/Upload a JSON file/)).toBeInTheDocument();
  });

  it('closes import modal when X is clicked', async () => {
    vi.mocked(api.get).mockResolvedValue({ data: [] } as never);

    renderWithClient(<Techniques />);

    await screen.findByText('Techniques');
    fireEvent.click(screen.getByText('Import Techniques'));

    // Find the close button in the modal header
    const modalHeader = screen.getByRole('heading', { name: 'Import Techniques' }).closest('div');
    const closeButton = modalHeader?.parentElement?.querySelector('button');
    if (closeButton) {
      fireEvent.click(closeButton);
    }

    await waitFor(() => {
      expect(screen.queryByText(/Upload a JSON file/)).not.toBeInTheDocument();
    });
  });

  it('imports techniques from valid JSON file', async () => {
    vi.mocked(api.get).mockResolvedValue({ data: [] } as never);
    vi.mocked(techniqueApi.import).mockResolvedValue({
      data: { imported: 2, failed: 0 },
    } as never);

    renderWithClient(<Techniques />);

    await screen.findByText('Techniques');
    fireEvent.click(screen.getByText('Import Techniques'));

    const file = new File(
      [JSON.stringify([{ id: 'T1082', name: 'Test', tactic: 'discovery', platforms: ['windows'], is_safe: true }])],
      'test.json',
      { type: 'application/json' }
    );

    const input = document.querySelector('input[type="file"]') as HTMLInputElement;
    fireEvent.change(input, { target: { files: [file] } });

    await waitFor(() => {
      expect(techniqueApi.import).toHaveBeenCalled();
    });
  });

  it('shows import success result', async () => {
    vi.mocked(api.get).mockResolvedValue({ data: [] } as never);
    vi.mocked(techniqueApi.import).mockResolvedValue({
      data: { imported: 3, failed: 0 },
    } as never);

    renderWithClient(<Techniques />);

    await screen.findByText('Techniques');
    fireEvent.click(screen.getByText('Import Techniques'));

    const file = new File(
      [JSON.stringify([{ id: 'T1082', name: 'Test', tactic: 'discovery', platforms: ['windows'], is_safe: true }])],
      'test.json',
      { type: 'application/json' }
    );

    const input = document.querySelector('input[type="file"]') as HTMLInputElement;
    fireEvent.change(input, { target: { files: [file] } });

    await waitFor(() => {
      expect(screen.getByText('Import Successful')).toBeInTheDocument();
      expect(screen.getByText('3 imported, 0 failed')).toBeInTheDocument();
    });
  });

  it('shows partial import result with errors', async () => {
    vi.mocked(api.get).mockResolvedValue({ data: [] } as never);
    vi.mocked(techniqueApi.import).mockResolvedValue({
      data: { imported: 1, failed: 2, errors: ['T1083: Invalid format', 'T1057: Missing field'] },
    } as never);

    renderWithClient(<Techniques />);

    await screen.findByText('Techniques');
    fireEvent.click(screen.getByText('Import Techniques'));

    const file = new File(
      [JSON.stringify([{ id: 'T1082', name: 'Test' }])],
      'test.json',
      { type: 'application/json' }
    );

    const input = document.querySelector('input[type="file"]') as HTMLInputElement;
    fireEvent.change(input, { target: { files: [file] } });

    await waitFor(() => {
      expect(screen.getByText('Partial Import')).toBeInTheDocument();
      expect(screen.getByText('1 imported, 2 failed')).toBeInTheDocument();
      expect(screen.getByText('T1083: Invalid format')).toBeInTheDocument();
    });
  });

  it('shows import failed result', async () => {
    vi.mocked(api.get).mockResolvedValue({ data: [] } as never);
    vi.mocked(techniqueApi.import).mockResolvedValue({
      data: { imported: 0, failed: 1, errors: ['All failed'] },
    } as never);

    renderWithClient(<Techniques />);

    await screen.findByText('Techniques');
    fireEvent.click(screen.getByText('Import Techniques'));

    const file = new File(
      [JSON.stringify([{ id: 'T1082', name: 'Test' }])],
      'test.json',
      { type: 'application/json' }
    );

    const input = document.querySelector('input[type="file"]') as HTMLInputElement;
    fireEvent.change(input, { target: { files: [file] } });

    await waitFor(() => {
      expect(screen.getByText('Import Failed')).toBeInTheDocument();
    });
  });

  it('shows error toast for invalid JSON file', async () => {
    vi.mocked(api.get).mockResolvedValue({ data: [] } as never);

    renderWithClient(<Techniques />);

    await screen.findByText('Techniques');
    fireEvent.click(screen.getByText('Import Techniques'));

    const file = new File(['not valid json'], 'test.json', { type: 'application/json' });

    const input = document.querySelector('input[type="file"]') as HTMLInputElement;
    fireEvent.change(input, { target: { files: [file] } });

    await waitFor(() => {
      expect(toast.error).toHaveBeenCalledWith('Failed to parse file');
    });
  });

  it('shows error toast for non-JSON file', async () => {
    vi.mocked(api.get).mockResolvedValue({ data: [] } as never);

    renderWithClient(<Techniques />);

    await screen.findByText('Techniques');
    fireEvent.click(screen.getByText('Import Techniques'));

    const file = new File(['yaml content'], 'test.yaml', { type: 'text/yaml' });

    const input = document.querySelector('input[type="file"]') as HTMLInputElement;
    fireEvent.change(input, { target: { files: [file] } });

    await waitFor(() => {
      expect(toast.error).toHaveBeenCalledWith('Please use JSON format for importing techniques');
    });
  });

  it('shows error toast for empty techniques array', async () => {
    vi.mocked(api.get).mockResolvedValue({ data: [] } as never);

    renderWithClient(<Techniques />);

    await screen.findByText('Techniques');
    fireEvent.click(screen.getByText('Import Techniques'));

    const file = new File([JSON.stringify([])], 'test.json', { type: 'application/json' });

    const input = document.querySelector('input[type="file"]') as HTMLInputElement;
    fireEvent.change(input, { target: { files: [file] } });

    await waitFor(() => {
      expect(toast.error).toHaveBeenCalledWith('Invalid format: expected techniques array');
    });
  });

  it('handles wrapped techniques format', async () => {
    vi.mocked(api.get).mockResolvedValue({ data: [] } as never);
    vi.mocked(techniqueApi.import).mockResolvedValue({
      data: { imported: 1, failed: 0 },
    } as never);

    renderWithClient(<Techniques />);

    await screen.findByText('Techniques');
    fireEvent.click(screen.getByText('Import Techniques'));

    const file = new File(
      [JSON.stringify({ techniques: [{ id: 'T1082', name: 'Test', tactic: 'discovery', platforms: ['windows'], is_safe: true }] })],
      'test.json',
      { type: 'application/json' }
    );

    const input = document.querySelector('input[type="file"]') as HTMLInputElement;
    fireEvent.change(input, { target: { files: [file] } });

    await waitFor(() => {
      expect(techniqueApi.import).toHaveBeenCalled();
    });
  });

  it('allows Import More after result', async () => {
    vi.mocked(api.get).mockResolvedValue({ data: [] } as never);
    vi.mocked(techniqueApi.import).mockResolvedValue({
      data: { imported: 1, failed: 0 },
    } as never);

    renderWithClient(<Techniques />);

    await screen.findByText('Techniques');
    fireEvent.click(screen.getByText('Import Techniques'));

    const file = new File(
      [JSON.stringify([{ id: 'T1082', name: 'Test', tactic: 'discovery', platforms: ['windows'], is_safe: true }])],
      'test.json',
      { type: 'application/json' }
    );

    const input = document.querySelector('input[type="file"]') as HTMLInputElement;
    fireEvent.change(input, { target: { files: [file] } });

    await waitFor(() => {
      expect(screen.getByText('Import Successful')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Import More'));

    await waitFor(() => {
      expect(screen.getByText(/Upload a JSON file/)).toBeInTheDocument();
    });
  });

  it('closes modal after clicking Done', async () => {
    vi.mocked(api.get).mockResolvedValue({ data: [] } as never);
    vi.mocked(techniqueApi.import).mockResolvedValue({
      data: { imported: 1, failed: 0 },
    } as never);

    renderWithClient(<Techniques />);

    await screen.findByText('Techniques');
    fireEvent.click(screen.getByText('Import Techniques'));

    const file = new File(
      [JSON.stringify([{ id: 'T1082', name: 'Test', tactic: 'discovery', platforms: ['windows'], is_safe: true }])],
      'test.json',
      { type: 'application/json' }
    );

    const input = document.querySelector('input[type="file"]') as HTMLInputElement;
    fireEvent.change(input, { target: { files: [file] } });

    await waitFor(() => {
      expect(screen.getByText('Import Successful')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Done'));

    await waitFor(() => {
      expect(screen.queryByText('Import Successful')).not.toBeInTheDocument();
    });
  });

  it('handles import API error', async () => {
    vi.mocked(api.get).mockResolvedValue({ data: [] } as never);
    vi.mocked(techniqueApi.import).mockRejectedValue({
      response: { data: { error: 'Server error' } },
    } as never);

    renderWithClient(<Techniques />);

    await screen.findByText('Techniques');
    fireEvent.click(screen.getByText('Import Techniques'));

    const file = new File(
      [JSON.stringify([{ id: 'T1082', name: 'Test', tactic: 'discovery', platforms: ['windows'], is_safe: true }])],
      'test.json',
      { type: 'application/json' }
    );

    const input = document.querySelector('input[type="file"]') as HTMLInputElement;
    fireEvent.change(input, { target: { files: [file] } });

    await waitFor(() => {
      expect(toast.error).toHaveBeenCalledWith('Server error');
    });
  });

  it('handles import API error without message', async () => {
    vi.mocked(api.get).mockResolvedValue({ data: [] } as never);
    vi.mocked(techniqueApi.import).mockRejectedValue(new Error('Network error') as never);

    renderWithClient(<Techniques />);

    await screen.findByText('Techniques');
    fireEvent.click(screen.getByText('Import Techniques'));

    const file = new File(
      [JSON.stringify([{ id: 'T1082', name: 'Test', tactic: 'discovery', platforms: ['windows'], is_safe: true }])],
      'test.json',
      { type: 'application/json' }
    );

    const input = document.querySelector('input[type="file"]') as HTMLInputElement;
    fireEvent.change(input, { target: { files: [file] } });

    await waitFor(() => {
      expect(toast.error).toHaveBeenCalledWith('Failed to import techniques');
    });
  });

  it('does nothing when file input change fires with no file selected', async () => {
    vi.mocked(api.get).mockResolvedValue({ data: [] } as never);

    renderWithClient(<Techniques />);

    await screen.findByText('Techniques');
    fireEvent.click(screen.getByText('Import Techniques'));

    const input = document.querySelector('input[type="file"]') as HTMLInputElement;
    // Fire change event with no files
    fireEvent.change(input, { target: { files: [] } });

    // No toast error and no import call
    expect(techniqueApi.import).not.toHaveBeenCalled();
    expect(toast.error).not.toHaveBeenCalled();
  });

  it('closes import modal and resets importResult via the close callback', async () => {
    vi.mocked(api.get).mockResolvedValue({ data: [] } as never);
    vi.mocked(techniqueApi.import).mockResolvedValue({
      data: { imported: 1, failed: 0 },
    } as never);

    renderWithClient(<Techniques />);

    await screen.findByText('Techniques');
    fireEvent.click(screen.getByText('Import Techniques'));

    // Import a file to get a result
    const file = new File(
      [JSON.stringify([{ id: 'T1082', name: 'Test', tactic: 'discovery', platforms: ['windows'], is_safe: true }])],
      'test.json',
      { type: 'application/json' }
    );

    const input = document.querySelector('input[type="file"]') as HTMLInputElement;
    fireEvent.change(input, { target: { files: [file] } });

    await waitFor(() => {
      expect(screen.getByText('Import Successful')).toBeInTheDocument();
    });

    // Click Done to close modal (this calls closeImportModal which resets both showImportModal and importResult)
    fireEvent.click(screen.getByText('Done'));

    await waitFor(() => {
      expect(screen.queryByText('Import Successful')).not.toBeInTheDocument();
    });

    // Re-open the modal - importResult should be reset, showing the upload form
    fireEvent.click(screen.getByText('Import Techniques'));
    expect(screen.getByText(/Upload a JSON file/)).toBeInTheDocument();
  });
});
