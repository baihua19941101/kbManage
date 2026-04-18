import '@testing-library/jest-dom/vitest';
import { render, screen, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { MemoryRouter } from 'react-router-dom';
import { ClusterLifecycleAuditPage } from '@/features/audit/pages/ClusterLifecycleAuditPage';
import { useAuthStore } from '@/features/auth/store';
import { listLifecycleAuditEvents } from '@/services/clusterLifecycle';
import { installAntdDomShims } from '@/test/installAntdDomShims';

vi.mock('@/services/clusterLifecycle', async () => ({
  ...(await vi.importActual<typeof import('@/services/clusterLifecycle')>('@/services/clusterLifecycle')),
  listLifecycleAuditEvents: vi.fn()
}));

const renderPage = () =>
  render(
    <QueryClientProvider client={new QueryClient({ defaultOptions: { queries: { retry: false } } })}>
      <MemoryRouter>
        <ClusterLifecycleAuditPage />
      </MemoryRouter>
    </QueryClientProvider>
  );

describe('ClusterLifecycleAuditPage', () => {
  beforeAll(() => installAntdDomShims());
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(listLifecycleAuditEvents).mockResolvedValue({ items: [] });
  });

  it('shows unauthorized empty', () => {
    useAuthStore.setState({ isAuthenticated: true, accessToken: 't', refreshToken: 'r', user: { id: 'u', username: 'u', platformRoles: [] } });
    renderPage();
    expect(screen.getByText('你暂无生命周期审计访问权限。')).toBeInTheDocument();
  });

  it('queries audit events', async () => {
    useAuthStore.setState({ isAuthenticated: true, accessToken: 't', refreshToken: 'r', user: { id: 'u', username: 'u', platformRoles: ['audit-reader'] } });
    vi.mocked(listLifecycleAuditEvents).mockResolvedValue({ items: [{ id: 'a1', action: 'clusterlifecycle.create', outcome: 'succeeded', occurredAt: '2026-04-17T10:00:00Z' }] });
    renderPage();
    await waitFor(() => expect(listLifecycleAuditEvents).toHaveBeenCalled());
    expect(await screen.findByText('clusterlifecycle.create')).toBeInTheDocument();
  });
});
