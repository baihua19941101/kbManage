import '@testing-library/jest-dom/vitest';
import { render, screen, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { MemoryRouter } from 'react-router-dom';
import { PlatformMarketplaceAuditPage } from '@/features/audit/pages/PlatformMarketplaceAuditPage';
import { useAuthStore } from '@/features/auth/store';
import { listPlatformMarketplaceAuditEvents } from '@/services/platformMarketplace';
import { installAntdDomShims } from '@/test/installAntdDomShims';

vi.mock('@/services/platformMarketplace', async () => ({
  ...(await vi.importActual<typeof import('@/services/platformMarketplace')>('@/services/platformMarketplace')),
  listPlatformMarketplaceAuditEvents: vi.fn()
}));

const renderPage = () =>
  render(
    <QueryClientProvider client={new QueryClient({ defaultOptions: { queries: { retry: false } } })}>
      <MemoryRouter>
        <PlatformMarketplaceAuditPage />
      </MemoryRouter>
    </QueryClientProvider>
  );

describe('PlatformMarketplaceAuditPage', () => {
  beforeAll(() => installAntdDomShims());
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(listPlatformMarketplaceAuditEvents).mockResolvedValue({ items: [] });
  });

  it('shows unauthorized empty', () => {
    useAuthStore.setState({
      isAuthenticated: true,
      accessToken: 't',
      refreshToken: 'r',
      user: { id: 'u', username: 'u', platformRoles: [] }
    });
    renderPage();
    expect(screen.getByText('你暂无市场审计访问权限。')).toBeInTheDocument();
  });

  it('queries audit events', async () => {
    useAuthStore.setState({
      isAuthenticated: true,
      accessToken: 't',
      refreshToken: 'r',
      user: { id: 'u', username: 'u', platformRoles: ['audit-reader'] }
    });
    vi.mocked(listPlatformMarketplaceAuditEvents).mockResolvedValue({
      items: [{ id: 'a1', action: 'platformmarketplace.template.publish', outcome: 'succeeded' }]
    });
    renderPage();
    await waitFor(() => expect(listPlatformMarketplaceAuditEvents).toHaveBeenCalled());
    expect(await screen.findByText('platformmarketplace.template.publish')).toBeInTheDocument();
  });
});
