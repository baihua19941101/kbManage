import { screen } from '@testing-library/react';
import { ExtensionCenterPage } from '@/features/platform-marketplace/pages/ExtensionCenterPage';
import { renderWithProviders } from '@/features/platform-marketplace/pages/testUtils';
import { useAuthStore } from '@/features/auth/store';
import { createExtensionPackage, disableExtension, enableExtension, listExtensions } from '@/services/platformMarketplace';
import { installAntdDomShims } from '@/test/installAntdDomShims';

vi.mock('@/services/platformMarketplace', async () => ({
  ...(await vi.importActual<typeof import('@/services/platformMarketplace')>('@/services/platformMarketplace')),
  listExtensions: vi.fn(),
  createExtensionPackage: vi.fn(),
  enableExtension: vi.fn(),
  disableExtension: vi.fn()
}));

describe('ExtensionCenterPage', () => {
  beforeAll(() => installAntdDomShims());

  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(listExtensions).mockResolvedValue({ items: [] });
    vi.mocked(createExtensionPackage).mockResolvedValue({ id: 'e1', name: 'mesh' });
    vi.mocked(enableExtension).mockResolvedValue({ id: 'e1', name: 'mesh', status: 'enabled' });
    vi.mocked(disableExtension).mockResolvedValue({ id: 'e1', name: 'mesh', status: 'disabled' });
  });

  it('shows unauthorized state', () => {
    useAuthStore.setState({ isAuthenticated: true, accessToken: 't', refreshToken: 'r', user: { id: 'u', username: 'u', platformRoles: [] } });
    renderWithProviders(<ExtensionCenterPage />);
    expect(screen.getByText('你暂无扩展中心访问权限。')).toBeInTheDocument();
  });

  it('renders extension list', async () => {
    useAuthStore.setState({ isAuthenticated: true, accessToken: 't', refreshToken: 'r', user: { id: 'u', username: 'u', platformRoles: ['platform-admin'] } });
    vi.mocked(listExtensions).mockResolvedValue({
      items: [{ id: 'ext-1', name: '服务网格观测', version: '2.1.0', status: 'enabled', compatibilityStatus: 'active', permissionSummary: '读取平台事件' }]
    });
    renderWithProviders(<ExtensionCenterPage />);
    expect(await screen.findByText('服务网格观测')).toBeInTheDocument();
  });
});
