import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import Layout from './Layout';

// Mock localStorage
const localStorageMock = {
  getItem: vi.fn(),
  setItem: vi.fn(),
  removeItem: vi.fn(),
  clear: vi.fn(),
};
Object.defineProperty(window, 'localStorage', { value: localStorageMock });

function renderLayout(path = '/dashboard') {
  return render(
    <MemoryRouter initialEntries={[path]}>
      <Layout>
        <div data-testid="child-content">Test Content</div>
      </Layout>
    </MemoryRouter>
  );
}

describe('Layout', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders AutoStrike brand name', () => {
    renderLayout();
    expect(screen.getByText('AutoStrike')).toBeInTheDocument();
  });

  it('renders BAS Platform subtitle', () => {
    renderLayout();
    expect(screen.getByText('BAS Platform')).toBeInTheDocument();
  });

  it('renders all navigation links', () => {
    renderLayout();

    expect(screen.getByText('Dashboard')).toBeInTheDocument();
    expect(screen.getByText('Agents')).toBeInTheDocument();
    expect(screen.getByText('Techniques')).toBeInTheDocument();
    expect(screen.getByText('Scenarios')).toBeInTheDocument();
    expect(screen.getByText('Executions')).toBeInTheDocument();
    expect(screen.getByText('Settings')).toBeInTheDocument();
  });

  it('renders children content', () => {
    renderLayout();
    expect(screen.getByTestId('child-content')).toBeInTheDocument();
    expect(screen.getByText('Test Content')).toBeInTheDocument();
  });

  it('displays default username when not in localStorage', () => {
    localStorageMock.getItem.mockReturnValue(null);
    renderLayout();
    expect(screen.getByText('Admin')).toBeInTheDocument();
  });

  it('displays username from localStorage', () => {
    localStorageMock.getItem.mockImplementation((key: string) => {
      if (key === 'username') return 'TestUser';
      return null;
    });
    renderLayout();
    expect(screen.getByText('TestUser')).toBeInTheDocument();
  });

  it('displays default email when not in localStorage', () => {
    localStorageMock.getItem.mockReturnValue(null);
    renderLayout();
    expect(screen.getByText('admin@autostrike.local')).toBeInTheDocument();
  });

  it('displays email from localStorage', () => {
    localStorageMock.getItem.mockImplementation((key: string) => {
      if (key === 'email') return 'test@example.com';
      return null;
    });
    renderLayout();
    expect(screen.getByText('test@example.com')).toBeInTheDocument();
  });

  it('displays first letter of username in avatar', () => {
    localStorageMock.getItem.mockImplementation((key: string) => {
      if (key === 'username') return 'John';
      return null;
    });
    renderLayout();
    expect(screen.getByText('J')).toBeInTheDocument();
  });
});

describe('Layout Navigation Active State', () => {
  it('highlights active route', () => {
    renderLayout('/dashboard');
    const dashboardLink = screen.getByText('Dashboard').closest('a');
    expect(dashboardLink).toHaveClass('bg-primary-600');
  });

  it('does not highlight inactive routes', () => {
    renderLayout('/dashboard');
    const agentsLink = screen.getByText('Agents').closest('a');
    expect(agentsLink).not.toHaveClass('bg-primary-600');
    expect(agentsLink).toHaveClass('text-gray-300');
  });
});
