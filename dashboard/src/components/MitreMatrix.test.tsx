import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { MitreMatrix } from './MitreMatrix';
import { Technique } from '../types';

const mockTechniques: Technique[] = [
  {
    id: 'T1082',
    name: 'System Information Discovery',
    description: 'Adversaries may attempt to get detailed information about the operating system.',
    tactic: 'discovery',
    platforms: ['windows', 'linux'],
    is_safe: true,
    detection: [
      { source: 'Process Creation', indicator: 'systeminfo.exe execution' },
    ],
  },
  {
    id: 'T1083',
    name: 'File and Directory Discovery',
    description: 'Adversaries may enumerate files and directories.',
    tactic: 'discovery',
    platforms: ['windows', 'linux'],
    is_safe: true,
    detection: [],
  },
  {
    id: 'T1059.001',
    name: 'PowerShell',
    description: 'Adversaries may abuse PowerShell commands.',
    tactic: 'execution',
    platforms: ['windows'],
    is_safe: false,
    detection: [
      { source: 'Script Block', indicator: 'PowerShell execution' },
    ],
  },
  {
    id: 'T1566',
    name: 'Phishing',
    description: 'Adversaries may send phishing messages.',
    tactic: 'initial_access',
    platforms: ['windows', 'linux', 'macos'],
    is_safe: true,
    detection: [],
  },
];

describe('MitreMatrix', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders all 14 tactic headers', () => {
    render(<MitreMatrix techniques={[]} />);

    expect(screen.getByText('Reconnaissance')).toBeInTheDocument();
    expect(screen.getByText('Resource Development')).toBeInTheDocument();
    expect(screen.getByText('Initial Access')).toBeInTheDocument();
    expect(screen.getByText('Execution')).toBeInTheDocument();
    expect(screen.getByText('Persistence')).toBeInTheDocument();
    expect(screen.getByText('Privilege Escalation')).toBeInTheDocument();
    expect(screen.getByText('Defense Evasion')).toBeInTheDocument();
    expect(screen.getByText('Credential Access')).toBeInTheDocument();
    expect(screen.getByText('Discovery')).toBeInTheDocument();
    expect(screen.getByText('Lateral Movement')).toBeInTheDocument();
    expect(screen.getByText('Collection')).toBeInTheDocument();
    expect(screen.getByText('C2')).toBeInTheDocument();
    expect(screen.getByText('Exfiltration')).toBeInTheDocument();
    expect(screen.getByText('Impact')).toBeInTheDocument();
  });

  it('renders techniques in correct tactic columns', () => {
    render(<MitreMatrix techniques={mockTechniques} />);

    // Discovery techniques
    expect(screen.getByText('T1082')).toBeInTheDocument();
    expect(screen.getByText('System Information Discovery')).toBeInTheDocument();
    expect(screen.getByText('T1083')).toBeInTheDocument();
    expect(screen.getByText('File and Directory Discovery')).toBeInTheDocument();

    // Execution technique
    expect(screen.getByText('T1059.001')).toBeInTheDocument();
    expect(screen.getByText('PowerShell')).toBeInTheDocument();

    // Initial Access technique
    expect(screen.getByText('T1566')).toBeInTheDocument();
    expect(screen.getByText('Phishing')).toBeInTheDocument();
  });

  it('renders technique count in headers', () => {
    render(<MitreMatrix techniques={mockTechniques} />);

    // Discovery has 2 techniques
    const discoveryCount = screen.getByText('Discovery').parentElement?.querySelector('.text-white\\/70');
    expect(discoveryCount?.textContent).toBe('2');
  });

  it('renders platform filter', () => {
    render(<MitreMatrix techniques={mockTechniques} />);

    expect(screen.getByText('Platform:')).toBeInTheDocument();
    expect(screen.getByRole('combobox')).toBeInTheDocument();
    expect(screen.getByText('All Platforms')).toBeInTheDocument();
  });

  it('filters techniques by platform', () => {
    render(<MitreMatrix techniques={mockTechniques} />);

    const platformSelect = screen.getByRole('combobox');

    // Filter by Windows only
    fireEvent.change(platformSelect, { target: { value: 'windows' } });

    // All techniques support Windows
    expect(screen.getByText('T1082')).toBeInTheDocument();
    expect(screen.getByText('T1059.001')).toBeInTheDocument();
  });

  it('shows technique count', () => {
    render(<MitreMatrix techniques={mockTechniques} />);

    expect(screen.getByText('4 techniques loaded')).toBeInTheDocument();
  });

  it('shows safe/unsafe indicators', () => {
    render(<MitreMatrix techniques={mockTechniques} />);

    // Check for safe (green) and unsafe (red) indicators
    const safeIndicators = document.querySelectorAll('.bg-green-500');
    const unsafeIndicators = document.querySelectorAll('.bg-red-500');

    expect(safeIndicators.length).toBeGreaterThan(0);
    expect(unsafeIndicators.length).toBeGreaterThan(0);
  });

  it('opens technique detail panel on click', () => {
    render(<MitreMatrix techniques={mockTechniques} />);

    const techniqueButton = screen.getByTitle('T1082: System Information Discovery');
    fireEvent.click(techniqueButton);

    // Modal should appear with technique details
    expect(screen.getByRole('heading', { name: 'System Information Discovery' })).toBeInTheDocument();
    expect(screen.getByText('Adversaries may attempt to get detailed information about the operating system.')).toBeInTheDocument();
    expect(screen.getByText('Safe')).toBeInTheDocument();
  });

  it('closes technique detail panel when close button clicked', () => {
    render(<MitreMatrix techniques={mockTechniques} />);

    // Open modal
    const techniqueButton = screen.getByTitle('T1082: System Information Discovery');
    fireEvent.click(techniqueButton);

    expect(screen.getByRole('heading', { name: 'System Information Discovery' })).toBeInTheDocument();

    // Close modal
    const closeButton = document.querySelector('.absolute.top-4.right-4');
    if (closeButton) {
      fireEvent.click(closeButton);
    }

    // Modal should be closed - look for dialog that's not visible
    expect(screen.queryByRole('heading', { name: 'System Information Discovery' })).not.toBeInTheDocument();
  });

  it('closes technique detail panel when overlay clicked', () => {
    render(<MitreMatrix techniques={mockTechniques} />);

    // Open modal
    const techniqueButton = screen.getByTitle('T1082: System Information Discovery');
    fireEvent.click(techniqueButton);

    // Click overlay to close
    const overlay = document.querySelector('.bg-gray-500.bg-opacity-75');
    if (overlay) {
      fireEvent.click(overlay);
    }

    expect(screen.queryByRole('heading', { name: 'System Information Discovery' })).not.toBeInTheDocument();
  });

  it('shows detection information in detail panel', () => {
    render(<MitreMatrix techniques={mockTechniques} />);

    const techniqueButton = screen.getByTitle('T1082: System Information Discovery');
    fireEvent.click(techniqueButton);

    expect(screen.getByText('Detection')).toBeInTheDocument();
    expect(screen.getByText('Process Creation:')).toBeInTheDocument();
    expect(screen.getByText('systeminfo.exe execution')).toBeInTheDocument();
  });

  it('shows platform badges in detail panel', () => {
    render(<MitreMatrix techniques={mockTechniques} />);

    const techniqueButton = screen.getByTitle('T1082: System Information Discovery');
    fireEvent.click(techniqueButton);

    // Find platform badges in the detail modal (they have specific badge classes)
    const modal = document.querySelector('.relative.bg-white.rounded-lg');
    expect(modal).toBeInTheDocument();

    const badges = modal?.querySelectorAll('.badge.bg-gray-100');
    expect(badges?.length).toBeGreaterThanOrEqual(2);
  });

  it('shows unsafe badge for unsafe technique', () => {
    render(<MitreMatrix techniques={mockTechniques} />);

    const techniqueButton = screen.getByTitle('T1059.001: PowerShell');
    fireEvent.click(techniqueButton);

    expect(screen.getByText('Unsafe')).toBeInTheDocument();
  });

  it('calls onTechniqueClick callback when provided', () => {
    const onTechniqueClick = vi.fn();
    render(<MitreMatrix techniques={mockTechniques} onTechniqueClick={onTechniqueClick} />);

    const techniqueButton = screen.getByTitle('T1082: System Information Discovery');
    fireEvent.click(techniqueButton);

    expect(onTechniqueClick).toHaveBeenCalledWith(mockTechniques[0]);
  });

  it('renders "No techniques" message for empty tactics', () => {
    render(<MitreMatrix techniques={mockTechniques} />);

    // Tactics without techniques should show "No techniques"
    const noTechMessages = screen.getAllByText('No techniques');
    expect(noTechMessages.length).toBeGreaterThan(0);
  });

  it('renders legend', () => {
    render(<MitreMatrix techniques={mockTechniques} />);

    expect(screen.getByText('Safe technique')).toBeInTheDocument();
    expect(screen.getByText('Potentially unsafe')).toBeInTheDocument();
  });

  it('handles techniques with hyphenated tactics', () => {
    // Test that component handles hyphenated tactics from backend
    const techniquesWithHyphenatedTactic: Technique[] = [
      {
        id: 'T1548',
        name: 'Abuse Elevation Control Mechanism',
        description: 'Test technique',
        tactic: 'privilege-escalation' as Technique['tactic'], // hyphenated form from backend
        platforms: ['windows'],
        is_safe: false,
        detection: [],
      },
    ];

    render(<MitreMatrix techniques={techniquesWithHyphenatedTactic} />);

    expect(screen.getByText('T1548')).toBeInTheDocument();
    expect(screen.getByText('Abuse Elevation Control Mechanism')).toBeInTheDocument();
  });

  it('sorts techniques by ID within each tactic', () => {
    const unsortedTechniques: Technique[] = [
      {
        id: 'T1083',
        name: 'File and Directory Discovery',
        description: 'Test',
        tactic: 'discovery',
        platforms: ['windows'],
        is_safe: true,
        detection: [],
      },
      {
        id: 'T1057',
        name: 'Process Discovery',
        description: 'Test',
        tactic: 'discovery',
        platforms: ['windows'],
        is_safe: true,
        detection: [],
      },
      {
        id: 'T1082',
        name: 'System Information Discovery',
        description: 'Test',
        tactic: 'discovery',
        platforms: ['windows'],
        is_safe: true,
        detection: [],
      },
    ];

    render(<MitreMatrix techniques={unsortedTechniques} />);

    const techniqueIds = screen.getAllByText(/^T\d+$/).map(el => el.textContent);
    // Should be sorted: T1057, T1082, T1083
    expect(techniqueIds).toContain('T1057');
    expect(techniqueIds).toContain('T1082');
    expect(techniqueIds).toContain('T1083');
  });

  it('renders available platforms in filter dropdown', () => {
    render(<MitreMatrix techniques={mockTechniques} />);

    const platformSelect = screen.getByRole('combobox');

    // Check for platform options
    expect(platformSelect).toContainHTML('linux');
    expect(platformSelect).toContainHTML('macos');
    expect(platformSelect).toContainHTML('windows');
  });

  it('closes technique detail panel when Escape key is pressed on overlay', () => {
    render(<MitreMatrix techniques={mockTechniques} />);

    // Open modal
    const techniqueButton = screen.getByTitle('T1082: System Information Discovery');
    fireEvent.click(techniqueButton);

    expect(screen.getByRole('heading', { name: 'System Information Discovery' })).toBeInTheDocument();

    // Press Escape on overlay
    const overlay = document.querySelector('.bg-gray-500.bg-opacity-75');
    if (overlay) {
      fireEvent.keyDown(overlay, { key: 'Escape' });
    }

    expect(screen.queryByRole('heading', { name: 'System Information Discovery' })).not.toBeInTheDocument();
  });

  it('does not close modal when non-Escape key is pressed on overlay', () => {
    render(<MitreMatrix techniques={mockTechniques} />);

    // Open modal
    const techniqueButton = screen.getByTitle('T1082: System Information Discovery');
    fireEvent.click(techniqueButton);

    expect(screen.getByRole('heading', { name: 'System Information Discovery' })).toBeInTheDocument();

    // Press Enter on overlay (should not close)
    const overlay = document.querySelector('.bg-gray-500.bg-opacity-75');
    if (overlay) {
      fireEvent.keyDown(overlay, { key: 'Enter' });
    }

    // Modal should still be open
    expect(screen.getByRole('heading', { name: 'System Information Discovery' })).toBeInTheDocument();
  });

  it('shows detection section only when detection data exists', () => {
    const techniqueWithNoDetection: Technique[] = [
      {
        id: 'T1083',
        name: 'File and Directory Discovery',
        description: 'Test technique',
        tactic: 'discovery',
        platforms: ['windows'],
        is_safe: true,
        detection: [],
      },
    ];

    render(<MitreMatrix techniques={techniqueWithNoDetection} />);

    const techniqueButton = screen.getByTitle('T1083: File and Directory Discovery');
    fireEvent.click(techniqueButton);

    // Detection section should not be visible when detection array is empty
    expect(screen.queryByText('Detection')).not.toBeInTheDocument();
  });

  it('filters out techniques for a platform that only some techniques support', () => {
    render(<MitreMatrix techniques={mockTechniques} />);

    const platformSelect = screen.getByRole('combobox');

    // Filter by macos - only T1566 (Phishing) supports macos
    fireEvent.change(platformSelect, { target: { value: 'macos' } });

    expect(screen.getByText('T1566')).toBeInTheDocument();
    expect(screen.getByText('Phishing')).toBeInTheDocument();

    // T1082 and T1059.001 do not support macos, should not be visible
    expect(screen.queryByText('T1082')).not.toBeInTheDocument();
    expect(screen.queryByText('T1059.001')).not.toBeInTheDocument();
  });

  it('resets to show all techniques when platform filter is changed back to all', () => {
    render(<MitreMatrix techniques={mockTechniques} />);

    const platformSelect = screen.getByRole('combobox');

    // Filter by macos first
    fireEvent.change(platformSelect, { target: { value: 'macos' } });
    expect(screen.queryByText('T1082')).not.toBeInTheDocument();

    // Reset to all
    fireEvent.change(platformSelect, { target: { value: 'all' } });

    // All techniques should be visible again
    expect(screen.getByText('T1082')).toBeInTheDocument();
    expect(screen.getByText('T1059.001')).toBeInTheDocument();
    expect(screen.getByText('T1566')).toBeInTheDocument();
  });
});

describe('MitreMatrix Multi-Tactic Support', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('displays a multi-tactic technique in multiple columns', () => {
    const multiTacticTechniques: Technique[] = [
      {
        id: 'T1078',
        name: 'Valid Accounts',
        description: 'Adversaries may use valid accounts.',
        tactic: 'initial_access',
        tactics: ['initial_access', 'persistence', 'privilege_escalation', 'defense_evasion'],
        platforms: ['windows', 'linux'],
        is_safe: false,
        detection: [],
      },
    ];

    render(<MitreMatrix techniques={multiTacticTechniques} />);

    // T1078 should appear in all 4 tactic columns
    const techniqueButtons = screen.getAllByTitle('T1078: Valid Accounts');
    expect(techniqueButtons.length).toBe(4);

    // Check the counts in the headers
    const initialAccessCount = screen.getByText('Initial Access').parentElement?.querySelector('.text-white\\/70');
    expect(initialAccessCount?.textContent).toBe('1');

    const persistenceCount = screen.getByText('Persistence').parentElement?.querySelector('.text-white\\/70');
    expect(persistenceCount?.textContent).toBe('1');

    const privEscCount = screen.getByText('Privilege Escalation').parentElement?.querySelector('.text-white\\/70');
    expect(privEscCount?.textContent).toBe('1');

    const defenseEvasionCount = screen.getByText('Defense Evasion').parentElement?.querySelector('.text-white\\/70');
    expect(defenseEvasionCount?.textContent).toBe('1');
  });

  it('falls back to single tactic when tactics array is absent', () => {
    const singleTacticTechniques: Technique[] = [
      {
        id: 'T1082',
        name: 'System Information Discovery',
        description: 'Test',
        tactic: 'discovery',
        // No tactics array
        platforms: ['windows'],
        is_safe: true,
        detection: [],
      },
    ];

    render(<MitreMatrix techniques={singleTacticTechniques} />);

    // Should appear only in discovery column
    const techniqueButtons = screen.getAllByTitle('T1082: System Information Discovery');
    expect(techniqueButtons.length).toBe(1);

    const discoveryCount = screen.getByText('Discovery').parentElement?.querySelector('.text-white\\/70');
    expect(discoveryCount?.textContent).toBe('1');
  });

  it('shows all tactics in the detail panel for a multi-tactic technique', () => {
    const multiTacticTechniques: Technique[] = [
      {
        id: 'T1078',
        name: 'Valid Accounts',
        description: 'Adversaries may use valid accounts.',
        tactic: 'initial_access',
        tactics: ['initial_access', 'persistence'],
        platforms: ['windows'],
        is_safe: false,
        detection: [],
      },
    ];

    render(<MitreMatrix techniques={multiTacticTechniques} />);

    // Click on the technique
    const techniqueButtons = screen.getAllByTitle('T1078: Valid Accounts');
    fireEvent.click(techniqueButtons[0]);

    // Detail panel should show all tactics
    expect(screen.getByText('initial access, persistence')).toBeInTheDocument();
  });

  it('counts multi-tactic techniques correctly in each column header', () => {
    const techniques: Technique[] = [
      {
        id: 'T1078',
        name: 'Valid Accounts',
        description: 'Test',
        tactic: 'initial_access',
        tactics: ['initial_access', 'persistence'],
        platforms: ['windows'],
        is_safe: false,
        detection: [],
      },
      {
        id: 'T1053.005',
        name: 'Scheduled Task',
        description: 'Test',
        tactic: 'persistence',
        tactics: ['persistence', 'execution'],
        platforms: ['windows'],
        is_safe: false,
        detection: [],
      },
    ];

    render(<MitreMatrix techniques={techniques} />);

    // Persistence should have 2 techniques (T1078 + T1053.005)
    const persistenceCount = screen.getByText('Persistence').parentElement?.querySelector('.text-white\\/70');
    expect(persistenceCount?.textContent).toBe('2');

    // Initial Access should have 1 (T1078)
    const initialAccessCount = screen.getByText('Initial Access').parentElement?.querySelector('.text-white\\/70');
    expect(initialAccessCount?.textContent).toBe('1');

    // Execution should have 1 (T1053.005)
    const executionCount = screen.getByText('Execution').parentElement?.querySelector('.text-white\\/70');
    expect(executionCount?.textContent).toBe('1');
  });
});
