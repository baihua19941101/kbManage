import { screen } from '@testing-library/react';
import { ExtensionCompatibilityPage } from '@/features/platform-marketplace/pages/ExtensionCompatibilityPage';
import { renderWithProviders } from '@/features/platform-marketplace/pages/testUtils';
import { useAuthStore } from '@/features/auth/store';
import { getExtensionCompatibility, listExtensions } from '@/services/platformMarketplace';
import { installAntdDomShims } from '@/test/installAntdDomShims';

vi.mock('@/services/platformMarketplace', async () => ({
  ...(await vi.importActual<typeof import('@/services/platformMarketplace')>('@/services/platformMarketplace')),
  listExtensions: vi.fn(),
  getExtensionCompatibility: vi.fn()
}));

describe('ExtensionCompatibilityPage', () => {
  beforeAll(() => installAntdDomShims());

  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(listExtensions).mockResolvedValue({ items: [] });
    vi.mocked(getExtensionCompatibility).mockResolvedValue({
      id: 'c1',
      compatibilityStatus: 'compatible',
      blockedReasons: [],
      suggestedActions: []
    });
  });

  it('shows unauthorized state', () => {
    useAuthStore.setState({ isAuthenticated: true, accessToken: 't', refreshToken: 'r', user: { id: 'u', username: 'u', platformRoles: [] } });
    renderWithProviders(<ExtensionCompatibilityPage />);
    expect(screen.getByText('你暂无扩展兼容性访问权限。')).toBeInTheDocument();
  });

  it('renders compatibility page shell', async () => {
    useAuthStore.setState({ isAuthenticated: true, accessToken: 't', refreshToken: 'r', user: { id: 'u', username: 'u', platformRoles: ['platform-admin'] } });
    vi.mocked(listExtensions).mockResolvedValue({
      items: [{ id: 'ext-1', name: '服务网格观测', version: '2.1.0' }]
    });
    renderWithProviders(<ExtensionCompatibilityPage />);
    expect(await screen.findByText('选择扩展')).toBeInTheDocument();
  });
});
