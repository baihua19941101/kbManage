import { screen } from '@testing-library/react';
import { ClusterLifecycleDetailPage } from '@/features/cluster-lifecycle/pages/ClusterLifecycleDetailPage';
import { renderWithProviders } from '@/features/cluster-lifecycle/pages/testUtils';
import { useAuthStore } from '@/features/auth/store';
import { getClusterLifecycleDetail } from '@/services/clusterLifecycle';
import { installAntdDomShims } from '@/test/installAntdDomShims';

vi.mock('@/services/clusterLifecycle', async () => ({
  ...(await vi.importActual<typeof import('@/services/clusterLifecycle')>('@/services/clusterLifecycle')),
  getClusterLifecycleDetail: vi.fn()
}));

describe('ClusterLifecycleDetailPage', () => {
  beforeAll(() => installAntdDomShims());
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(getClusterLifecycleDetail).mockResolvedValue({
      id: 'c1',
      name: 'prod-a',
      status: 'active',
      healthStatus: 'healthy',
      kubernetesVersion: 'v1.30.4',
      cluster: { id: 'c1', name: 'prod-a', status: 'active', healthStatus: 'healthy', kubernetesVersion: 'v1.30.4' },
      nodePools: [],
      upgradePlans: []
    });
  });

  it('shows unauthorized empty', () => {
    useAuthStore.setState({ isAuthenticated: true, accessToken: 't', refreshToken: 'r', user: { id: 'u', username: 'u', platformRoles: [] } });
    renderWithProviders(<ClusterLifecycleDetailPage />, ['/cluster-lifecycle/clusters/c1']);
    expect(screen.getByText('你暂无生命周期详情访问权限。')).toBeInTheDocument();
  });

  it('renders detail when clusterId exists', async () => {
    useAuthStore.setState({ isAuthenticated: true, accessToken: 't', refreshToken: 'r', user: { id: 'u', username: 'u', platformRoles: ['platform-admin'] } });
    renderWithProviders(<ClusterLifecycleDetailPage />, ['/cluster-lifecycle/clusters/c1?clusterId=c1']);
    expect(await screen.findByText('prod-a')).toBeInTheDocument();
  });
});
