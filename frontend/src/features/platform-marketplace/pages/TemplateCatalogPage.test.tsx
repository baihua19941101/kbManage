import { screen } from '@testing-library/react';
import { TemplateCatalogPage } from '@/features/platform-marketplace/pages/TemplateCatalogPage';
import { renderWithProviders } from '@/features/platform-marketplace/pages/testUtils';
import { useAuthStore } from '@/features/auth/store';
import { listTemplates } from '@/services/platformMarketplace';
import { installAntdDomShims } from '@/test/installAntdDomShims';

vi.mock('@/services/platformMarketplace', async () => ({
  ...(await vi.importActual<typeof import('@/services/platformMarketplace')>('@/services/platformMarketplace')),
  listTemplates: vi.fn()
}));

describe('TemplateCatalogPage', () => {
  beforeAll(() => installAntdDomShims());

  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(listTemplates).mockResolvedValue({ items: [] });
  });

  it('shows unauthorized state', () => {
    useAuthStore.setState({ isAuthenticated: true, accessToken: 't', refreshToken: 'r', user: { id: 'u', username: 'u', platformRoles: [] } });
    renderWithProviders(<TemplateCatalogPage />);
    expect(screen.getByText('你暂无模板中心访问权限。')).toBeInTheDocument();
  });

  it('renders template list', async () => {
    useAuthStore.setState({ isAuthenticated: true, accessToken: 't', refreshToken: 'r', user: { id: 'u', username: 'u', platformRoles: ['platform-admin'] } });
    vi.mocked(listTemplates).mockResolvedValue({
      items: [{ id: 'tpl-1', name: '标准 Nginx', category: 'web', status: 'active', latestVersion: '1.2.0' }]
    });
    renderWithProviders(<TemplateCatalogPage />);
    expect(await screen.findByText('标准 Nginx')).toBeInTheDocument();
  });
});
