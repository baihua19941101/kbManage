import '@testing-library/jest-dom/vitest';
import { render, screen } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { MemoryRouter } from 'react-router-dom';
import { useAuthStore } from '@/features/auth/store';
import { AlertRulePage } from '@/features/observability/pages/AlertRulePage';
import { installAntdDomShims } from '@/test/installAntdDomShims';

Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: vi.fn().mockImplementation(() => ({
    matches: false,
    media: '',
    onchange: null,
    addListener: vi.fn(),
    removeListener: vi.fn(),
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
    dispatchEvent: vi.fn()
  }))
});

vi.mock('@/services/observability/alertRules', () => ({
  listAlertRules: vi.fn().mockResolvedValue({
    items: [
      {
        id: 'rule-1',
        name: 'cpu high',
        severity: 'critical',
        status: 'enabled',
        conditionExpression: 'cpu_usage > 80'
      }
    ]
  }),
  createAlertRule: vi.fn().mockResolvedValue({}),
  deleteAlertRule: vi.fn().mockResolvedValue({})
}));

describe('AlertRulePage', () => {
  beforeAll(() => {
    installAntdDomShims();
  });

  beforeEach(() => {
    useAuthStore.setState({
      isAuthenticated: true,
      accessToken: 'token',
      refreshToken: 'refresh-token',
      user: {
        id: 'user-1',
        username: 'tester',
        platformRoles: ['platform-admin']
      }
    });
  });

  it('renders alert rule page and table rows', async () => {
    const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
    render(
      <QueryClientProvider client={client}>
        <MemoryRouter>
          <AlertRulePage />
        </MemoryRouter>
      </QueryClientProvider>
    );

    expect(screen.getByText('告警规则治理')).toBeInTheDocument();
    expect(await screen.findByText(/cpu high/i)).toBeInTheDocument();
  });
});
