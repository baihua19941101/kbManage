import { screen } from '@testing-library/react';
import { TemplateDistributionPage } from '@/features/platform-marketplace/pages/TemplateDistributionPage';
import { renderWithProviders } from '@/features/platform-marketplace/pages/testUtils';
import { useAuthStore } from '@/features/auth/store';
import { createTemplateRelease, listTemplateReleases, listTemplates, syncCatalogSource } from '@/services/platformMarketplace';
import { installAntdDomShims } from '@/test/installAntdDomShims';

vi.mock('@/services/platformMarketplace', async () => ({
  ...(await vi.importActual<typeof import('@/services/platformMarketplace')>('@/services/platformMarketplace')),
  listTemplates: vi.fn(),
  listTemplateReleases: vi.fn(),
  createTemplateRelease: vi.fn(),
  syncCatalogSource: vi.fn()
}));

describe('TemplateDistributionPage', () => {
  beforeAll(() => installAntdDomShims());

  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(listTemplates).mockResolvedValue({ items: [] });
    vi.mocked(listTemplateReleases).mockResolvedValue({ items: [] });
    vi.mocked(createTemplateRelease).mockResolvedValue({ id: 'r1', version: '1.0.0' });
    vi.mocked(syncCatalogSource).mockResolvedValue({ id: 's1', name: 'repo' });
  });

  it('shows unauthorized state', () => {
    useAuthStore.setState({ isAuthenticated: true, accessToken: 't', refreshToken: 'r', user: { id: 'u', username: 'u', platformRoles: [] } });
    renderWithProviders(<TemplateDistributionPage />);
    expect(screen.getByText('你暂无模板分发访问权限。')).toBeInTheDocument();
  });

  it('renders template selector', async () => {
    useAuthStore.setState({ isAuthenticated: true, accessToken: 't', refreshToken: 'r', user: { id: 'u', username: 'u', platformRoles: ['platform-admin'] } });
    vi.mocked(listTemplates).mockResolvedValue({
      items: [{ id: 'tpl-1', name: '标准 Redis', latestVersion: '7.0.0' }]
    });
    renderWithProviders(<TemplateDistributionPage />);
    expect(await screen.findByText('选择模板')).toBeInTheDocument();
  });
});
