import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import App from './App';

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: false,
    },
  },
});

function renderWithProviders(initialRoute = '/') {
  return render(
    <QueryClientProvider client={queryClient}>
      <MemoryRouter initialEntries={[initialRoute]}>
        <App />
      </MemoryRouter>
    </QueryClientProvider>
  );
}

describe('App', () => {
  it('renders layout with AutoStrike title', () => {
    renderWithProviders('/dashboard');
    expect(screen.getByText('AutoStrike')).toBeInTheDocument();
  });

  it('renders navigation items', () => {
    renderWithProviders('/dashboard');
    expect(screen.getByRole('link', { name: 'Dashboard' })).toBeInTheDocument();
    expect(screen.getByRole('link', { name: 'Agents' })).toBeInTheDocument();
    expect(screen.getByRole('link', { name: 'Techniques' })).toBeInTheDocument();
    expect(screen.getByRole('link', { name: 'Scenarios' })).toBeInTheDocument();
    expect(screen.getByRole('link', { name: 'Executions' })).toBeInTheDocument();
    expect(screen.getByRole('link', { name: 'Settings' })).toBeInTheDocument();
  });

  it('redirects from / to /dashboard', () => {
    renderWithProviders('/');
    // After redirect, Dashboard nav link should exist
    expect(screen.getByRole('link', { name: 'Dashboard' })).toBeInTheDocument();
  });

  it('renders BAS Platform subtitle', () => {
    renderWithProviders('/dashboard');
    expect(screen.getByText('BAS Platform')).toBeInTheDocument();
  });
});

describe('Navigation', () => {
  it('highlights Dashboard link when on dashboard route', () => {
    renderWithProviders('/dashboard');
    const dashboardLink = screen.getByRole('link', { name: 'Dashboard' });
    expect(dashboardLink).toHaveClass('bg-primary-600');
  });

  it('highlights Agents link when on agents route', () => {
    renderWithProviders('/agents');
    const agentsLink = screen.getByRole('link', { name: 'Agents' });
    expect(agentsLink).toHaveClass('bg-primary-600');
  });

  it('highlights Settings link when on settings route', () => {
    renderWithProviders('/settings');
    const settingsLink = screen.getByRole('link', { name: 'Settings' });
    expect(settingsLink).toHaveClass('bg-primary-600');
  });
});
