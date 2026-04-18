import { screen } from '@testing-library/react';
import { ClusterRetirementPage } from '@/features/cluster-lifecycle/pages/ClusterRetirementPage';
import { renderWithProviders } from '@/features/cluster-lifecycle/pages/testUtils';
import { useAuthStore } from '@/features/auth/store';
import { getClusterLifecycleDetail } from '@/services/clusterLifecycle';
import { installAntdDomShims } from '@/test/installAntdDomShims';

vi.mock('@/services/clusterLifecycle', async () => ({
  ...(await vi.importActual<typeof import('@/services/clusterLifecycle')>('@/services/clusterLifecycle')),
  getClusterLifecycleDetail: vi.fn(),
  disableCluster: vi.fn(),
  retireCluster: vi.fn()
}));

describe('ClusterRetirementPage', () => {
  beforeAll(() => installAntdDomShims());
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(getClusterLifecycleDetail).mockResolvedValue({
      id: 'c1',
      name: 'prod-a',
      status: 'active',
      cluster: { id: 'c1', name: 'prod-a', status: 'active' },
      nodePools: [],
      upgradePlans: []
    });
  });

  it('shows unauthorized empty', () => {
    useAuthStore.setState({ isAuthenticated: true, accessToken: 't', refreshToken: 'r', user: { id: 'u', username: 'u', platformRoles: [] } });
    renderWithProviders(<ClusterRetirementPage />, ['/cluster-lifecycle/retire?clusterId=c1']);
    expect(screen.getByText('你暂无停用/退役访问权限。')).toBeInTheDocument();
  });

  it('renders cluster summary', async () => {
    useAuthStore.setState({ isAuthenticated: true, accessToken: 't', refreshToken: 'r', user: { id: 'u', username: 'u', platformRoles: ['platform-admin'] } });
    renderWithProviders(<ClusterRetirementPage />, ['/cluster-lifecycle/retire?clusterId=c1']);
    expect(await screen.findByText('prod-a')).toBeInTheDocument();
  });
});
