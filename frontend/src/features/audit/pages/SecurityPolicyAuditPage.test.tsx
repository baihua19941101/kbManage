import '@testing-library/jest-dom/vitest';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { MemoryRouter } from 'react-router-dom';
import { SecurityPolicyAuditPage } from '@/features/audit/pages/SecurityPolicyAuditPage';
import { useAuthStore } from '@/features/auth/store';
import { listSecurityPolicyAuditEvents } from '@/services/audit';
import { installAntdDomShims } from '@/test/installAntdDomShims';

vi.mock('@/services/audit', async () => {
  const actual = await vi.importActual<typeof import('@/services/audit')>('@/services/audit');
  return {
    ...actual,
    listSecurityPolicyAuditEvents: vi.fn()
  };
});

const renderWithClient = () => {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(
    <QueryClientProvider client={client}>
      <MemoryRouter>
        <SecurityPolicyAuditPage />
      </MemoryRouter>
    </QueryClientProvider>
  );
};

describe('SecurityPolicyAuditPage', () => {
  beforeAll(() => {
    installAntdDomShims();
  });

  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(listSecurityPolicyAuditEvents).mockResolvedValue({ items: [] });
  });

  it('shows unauthorized empty when user cannot read security policy audit', () => {
    useAuthStore.setState({
      isAuthenticated: true,
      accessToken: 'token',
      refreshToken: 'refresh',
      user: {
        id: 'u0',
        username: 'u0',
        platformRoles: []
      }
    });

    renderWithClient();

    expect(
      screen.getByText('你暂无策略审计访问权限，请联系管理员授予 securitypolicy:read 或审计角色。')
    ).toBeInTheDocument();
  });

  it('queries security policy audit events by actor/action/result', async () => {
    useAuthStore.setState({
      isAuthenticated: true,
      accessToken: 'token',
      refreshToken: 'refresh',
      user: {
        id: 'u1',
        username: 'auditor',
        platformRoles: ['audit-reader']
      }
    });

    vi.mocked(listSecurityPolicyAuditEvents).mockResolvedValue({
      items: [
        {
          id: 'evt-1',
          action: 'securitypolicy.exception.review',
          eventType: 'securitypolicy.exception.review',
          actorUserId: 'auditor',
          result: 'success',
          occurredAt: '2026-04-14T10:00:00Z'
        }
      ]
    });

    renderWithClient();

    await waitFor(() => {
      expect(listSecurityPolicyAuditEvents).toHaveBeenCalledWith({});
    });

    fireEvent.change(screen.getByPlaceholderText('例如：admin'), {
      target: { value: 'auditor' }
    });
    fireEvent.change(screen.getByPlaceholderText('例如：securitypolicy.exception.review'), {
      target: { value: 'securitypolicy.exception.review' }
    });
    fireEvent.change(screen.getByPlaceholderText('例如：success / failed / denied'), {
      target: { value: 'success' }
    });
    fireEvent.submit(screen.getByRole('textbox', { name: '操作者' }).closest('form') as HTMLFormElement);

    await waitFor(() => {
      expect(listSecurityPolicyAuditEvents).toHaveBeenLastCalledWith({
        actorUserId: 'auditor',
        action: 'securitypolicy.exception.review',
        result: 'success',
        from: undefined,
        to: undefined
      });
    });

    expect(await screen.findByText('securitypolicy.exception.review')).toBeInTheDocument();
  });
});
