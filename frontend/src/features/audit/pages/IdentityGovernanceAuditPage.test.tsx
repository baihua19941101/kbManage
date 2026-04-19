import '@testing-library/jest-dom/vitest';
import { render, screen, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { MemoryRouter } from 'react-router-dom';
import { IdentityGovernanceAuditPage } from '@/features/audit/pages/IdentityGovernanceAuditPage';
import { useAuthStore } from '@/features/auth/store';
import { listIdentityGovernanceAuditEvents } from '@/services/identityTenancy';
import { installAntdDomShims } from '@/test/installAntdDomShims';

vi.mock('@/services/identityTenancy', async () => ({
  ...(await vi.importActual<typeof import('@/services/identityTenancy')>('@/services/identityTenancy')),
  listIdentityGovernanceAuditEvents: vi.fn()
}));

const renderPage = () =>
  render(
    <QueryClientProvider client={new QueryClient({ defaultOptions: { queries: { retry: false } } })}>
      <MemoryRouter>
        <IdentityGovernanceAuditPage />
      </MemoryRouter>
    </QueryClientProvider>
  );

describe('IdentityGovernanceAuditPage', () => {
  beforeAll(() => installAntdDomShims());
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(listIdentityGovernanceAuditEvents).mockResolvedValue({ items: [] });
  });

  it('shows unauthorized empty', () => {
    useAuthStore.setState({
      isAuthenticated: true,
      accessToken: 't',
      refreshToken: 'r',
      user: { id: 'u', username: 'u', platformRoles: [] }
    });
    renderPage();
    expect(screen.getByText('你暂无身份治理审计访问权限。')).toBeInTheDocument();
  });

  it('queries audit events', async () => {
    useAuthStore.setState({
      isAuthenticated: true,
      accessToken: 't',
      refreshToken: 'r',
      user: { id: 'u', username: 'u', platformRoles: ['audit-reader'] }
    });
    vi.mocked(listIdentityGovernanceAuditEvents).mockResolvedValue({
      items: [
        {
          id: 'a1',
          action: 'identitytenancy.assignment.create',
          outcome: 'succeeded',
          occurredAt: '2026-04-19T00:00:00Z'
        }
      ]
    });
    renderPage();
    await waitFor(() => expect(listIdentityGovernanceAuditEvents).toHaveBeenCalled());
    expect(await screen.findByText('identitytenancy.assignment.create')).toBeInTheDocument();
  });
});
