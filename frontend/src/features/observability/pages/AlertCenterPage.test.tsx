import '@testing-library/jest-dom/vitest';
import { render, screen } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { MemoryRouter } from 'react-router-dom';
import { useAuthStore } from '@/features/auth/store';
import { AlertCenterPage } from '@/features/observability/pages/AlertCenterPage';
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

vi.mock('@/services/observability/alerts', () => ({
  listAlerts: vi.fn().mockResolvedValue({
    items: [
      {
        id: 'inc-1',
        severity: 'warning',
        status: 'firing',
        summary: 'mock alert'
      }
    ]
  }),
  acknowledgeAlert: vi.fn().mockResolvedValue({}),
  createAlertHandlingRecord: vi.fn().mockResolvedValue({})
}));

describe('AlertCenterPage', () => {
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

  it('renders alert center and list items', async () => {
    const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
    render(
      <QueryClientProvider client={client}>
        <MemoryRouter>
          <AlertCenterPage />
        </MemoryRouter>
      </QueryClientProvider>
    );

    expect(screen.getByText('告警中心')).toBeInTheDocument();
    expect(await screen.findByText(/mock alert/i)).toBeInTheDocument();
  });
});
