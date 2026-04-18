import { screen, waitFor } from '@testing-library/react';
import { ClusterDriverPage } from '@/features/cluster-lifecycle/pages/ClusterDriverPage';
import { renderWithProviders } from '@/features/cluster-lifecycle/pages/testUtils';
import { useAuthStore } from '@/features/auth/store';
import { listClusterDrivers } from '@/services/clusterLifecycle';
import { installAntdDomShims } from '@/test/installAntdDomShims';

vi.mock('@/services/clusterLifecycle', async () => ({
  ...(await vi.importActual<typeof import('@/services/clusterLifecycle')>('@/services/clusterLifecycle')),
  listClusterDrivers: vi.fn(),
  createClusterDriver: vi.fn()
}));

describe('ClusterDriverPage', () => {
  beforeAll(() => installAntdDomShims());
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(listClusterDrivers).mockResolvedValue({ items: [{ id: 'd1', driverKey: 'rke2', version: '1.0.0', providerType: 'cloud' }] });
  });

  it('shows unauthorized empty', () => {
    useAuthStore.setState({ isAuthenticated: true, accessToken: 't', refreshToken: 'r', user: { id: 'u', username: 'u', platformRoles: [] } });
    renderWithProviders(<ClusterDriverPage />);
    expect(screen.getByText('你暂无驱动管理访问权限。')).toBeInTheDocument();
  });

  it('renders driver table', async () => {
    useAuthStore.setState({ isAuthenticated: true, accessToken: 't', refreshToken: 'r', user: { id: 'u', username: 'u', platformRoles: ['platform-admin'] } });
    renderWithProviders(<ClusterDriverPage />);
    await waitFor(() => expect(listClusterDrivers).toHaveBeenCalled());
    expect(await screen.findByText('1.0.0')).toBeInTheDocument();
  });
});
