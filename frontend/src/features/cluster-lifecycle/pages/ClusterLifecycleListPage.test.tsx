import { screen } from '@testing-library/react';
import { ClusterLifecycleListPage } from '@/features/cluster-lifecycle/pages/ClusterLifecycleListPage';
import { renderWithProviders } from '@/features/cluster-lifecycle/pages/testUtils';
import { useAuthStore } from '@/features/auth/store';
import { listClusterLifecycleRecords } from '@/services/clusterLifecycle';
import { installAntdDomShims } from '@/test/installAntdDomShims';

vi.mock('@/services/clusterLifecycle', async () => ({
  ...(await vi.importActual<typeof import('@/services/clusterLifecycle')>('@/services/clusterLifecycle')),
  listClusterLifecycleRecords: vi.fn(),
  importCluster: vi.fn()
}));

describe('ClusterLifecycleListPage', () => {
  beforeAll(() => installAntdDomShims());

  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(listClusterLifecycleRecords).mockResolvedValue({ items: [] });
  });

  it('shows unauthorized empty', () => {
    useAuthStore.setState({ isAuthenticated: true, accessToken: 't', refreshToken: 'r', user: { id: 'u', username: 'u', platformRoles: [] } });
    renderWithProviders(<ClusterLifecycleListPage />);
    expect(screen.getByText('你暂无集群生命周期中心访问权限。')).toBeInTheDocument();
  });

  it('renders cluster table', async () => {
    useAuthStore.setState({ isAuthenticated: true, accessToken: 't', refreshToken: 'r', user: { id: 'u', username: 'u', platformRoles: ['platform-admin'] } });
    vi.mocked(listClusterLifecycleRecords).mockResolvedValue({
      items: [{ id: 'c1', name: 'prod-a', status: 'active', healthStatus: 'healthy', registrationStatus: 'connected', lifecycleMode: 'imported' }]
    });
    renderWithProviders(<ClusterLifecycleListPage />);
    expect(await screen.findByText('prod-a')).toBeInTheDocument();
  });
});
