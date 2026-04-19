import { screen } from '@testing-library/react';
import { InstallationRecordPage } from '@/features/platform-marketplace/pages/InstallationRecordPage';
import { renderWithProviders } from '@/features/platform-marketplace/pages/testUtils';
import { useAuthStore } from '@/features/auth/store';
import { listInstallations } from '@/services/platformMarketplace';
import { installAntdDomShims } from '@/test/installAntdDomShims';

vi.mock('@/services/platformMarketplace', async () => ({
  ...(await vi.importActual<typeof import('@/services/platformMarketplace')>('@/services/platformMarketplace')),
  listInstallations: vi.fn()
}));

describe('InstallationRecordPage', () => {
  beforeAll(() => installAntdDomShims());

  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(listInstallations).mockResolvedValue({ items: [] });
  });

  it('shows unauthorized state', () => {
    useAuthStore.setState({ isAuthenticated: true, accessToken: 't', refreshToken: 'r', user: { id: 'u', username: 'u', platformRoles: [] } });
    renderWithProviders(<InstallationRecordPage />);
    expect(screen.getByText('你暂无安装记录访问权限。')).toBeInTheDocument();
  });

  it('renders installation record table', async () => {
    useAuthStore.setState({ isAuthenticated: true, accessToken: 't', refreshToken: 'r', user: { id: 'u', username: 'u', platformRoles: ['platform-admin'] } });
    vi.mocked(listInstallations).mockResolvedValue({
      items: [{ id: 'ins-1', templateName: '标准 Nginx', currentVersion: '1.0.0', latestVersion: '1.1.0', status: 'active' }]
    });
    renderWithProviders(<InstallationRecordPage />);
    expect(await screen.findByText('标准 Nginx')).toBeInTheDocument();
  });
});
