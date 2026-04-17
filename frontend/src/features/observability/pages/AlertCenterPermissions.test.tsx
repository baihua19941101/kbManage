import '@testing-library/jest-dom/vitest';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { MemoryRouter } from 'react-router-dom';
import { useAuthStore } from '@/features/auth/store';
import { AlertCenterPage } from '@/features/observability/pages/AlertCenterPage';
import { ApiError } from '@/services/api/client';
import {
  acknowledgeAlert,
  createAlertHandlingRecord,
  listAlerts
} from '@/services/observability/alerts';
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
  listAlerts: vi.fn(),
  acknowledgeAlert: vi.fn(),
  createAlertHandlingRecord: vi.fn()
}));

describe('AlertCenterPermissions', () => {
  beforeAll(() => {
    installAntdDomShims();
  });

  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(listAlerts).mockResolvedValue({
      items: [
        {
          id: 'inc-1',
          severity: 'warning',
          status: 'firing',
          summary: 'mock alert'
        }
      ]
    });
    vi.mocked(acknowledgeAlert).mockResolvedValue({
      id: 'inc-1',
      severity: 'warning',
      status: 'acknowledged',
      summary: 'mock alert'
    });
    vi.mocked(createAlertHandlingRecord).mockResolvedValue({
      id: 'record-1',
      incidentId: 'inc-1',
      actionType: 'note'
    });
  });

  it('keeps page in readonly mode and disables acknowledge actions', async () => {
    useAuthStore.setState({
      isAuthenticated: true,
      accessToken: 'token',
      refreshToken: 'refresh-token',
      user: {
        id: 'readonly-user',
        username: 'readonly-user',
        platformRoles: ['readonly']
      }
    });

    const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });

    render(
      <QueryClientProvider client={client}>
        <MemoryRouter>
          <AlertCenterPage />
        </MemoryRouter>
      </QueryClientProvider>
    );

    expect(await screen.findByText('当前为只读模式')).toBeInTheDocument();
    expect(await screen.findByText(/mock alert/i)).toBeInTheDocument();

    const acknowledgeButton = screen.getByRole('button', { name: /确\s*认/ });
    const handlingNote = screen.getByPlaceholderText('输入处理说明（用于确认/记录）');

    expect(acknowledgeButton).toBeDisabled();
    expect(handlingNote).toBeDisabled();
  });

  it('shows permission revoked warning and disables actions after acknowledge is rejected with 403', async () => {
    useAuthStore.setState({
      isAuthenticated: true,
      accessToken: 'token',
      refreshToken: 'refresh-token',
      user: {
        id: 'operator-user',
        username: 'operator-user',
        platformRoles: ['ops-operator']
      }
    });

    vi.mocked(acknowledgeAlert).mockRejectedValue(
      new ApiError(403, 'forbidden', {
        url: '/api/v1/observability/alerts/inc-1/acknowledge'
      })
    );
    vi.mocked(createAlertHandlingRecord).mockRejectedValue(
      new ApiError(403, 'forbidden', {
        url: '/api/v1/observability/alerts/inc-1/handling-records'
      })
    );

    const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });

    render(
      <QueryClientProvider client={client}>
        <MemoryRouter>
          <AlertCenterPage />
        </MemoryRouter>
      </QueryClientProvider>
    );

    const acknowledgeButton = await screen.findByRole('button', { name: /确\s*认/ });
    fireEvent.click(acknowledgeButton);

    expect(await screen.findByText('权限已回收')).toBeInTheDocument();
    await waitFor(() => {
      expect(screen.getByRole('button', { name: /确\s*认/ })).toBeDisabled();
    });
  });
});
