import { screen } from '@testing-library/react';
import { CatalogSourcePage } from '@/features/platform-marketplace/pages/CatalogSourcePage';
import { renderWithProviders } from '@/features/platform-marketplace/pages/testUtils';
import { useAuthStore } from '@/features/auth/store';
import { createCatalogSource, listCatalogSources, syncCatalogSource } from '@/services/platformMarketplace';
import { installAntdDomShims } from '@/test/installAntdDomShims';

vi.mock('@/services/platformMarketplace', async () => ({
  ...(await vi.importActual<typeof import('@/services/platformMarketplace')>('@/services/platformMarketplace')),
  listCatalogSources: vi.fn(),
  createCatalogSource: vi.fn(),
  syncCatalogSource: vi.fn()
}));

describe('CatalogSourcePage', () => {
  beforeAll(() => installAntdDomShims());

  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(listCatalogSources).mockResolvedValue({ items: [] });
    vi.mocked(createCatalogSource).mockResolvedValue({ id: 's1', name: 'repo' });
    vi.mocked(syncCatalogSource).mockResolvedValue({ id: 's1', name: 'repo', syncStatus: 'pending' });
  });

  it('shows unauthorized state', () => {
    useAuthStore.setState({ isAuthenticated: true, accessToken: 't', refreshToken: 'r', user: { id: 'u', username: 'u', platformRoles: [] } });
    renderWithProviders(<CatalogSourcePage />);
    expect(screen.getByText('你暂无应用目录访问权限。')).toBeInTheDocument();
  });

  it('renders catalog source list', async () => {
    useAuthStore.setState({ isAuthenticated: true, accessToken: 't', refreshToken: 'r', user: { id: 'u', username: 'u', platformRoles: ['platform-admin'] } });
    vi.mocked(listCatalogSources).mockResolvedValue({
      items: [{ id: 'source-1', name: '平台标准 Helm 目录', sourceType: 'helm-repository', status: 'active', syncStatus: 'synced' }]
    });
    renderWithProviders(<CatalogSourcePage />);
    expect(await screen.findByText('平台标准 Helm 目录')).toBeInTheDocument();
  });
});
