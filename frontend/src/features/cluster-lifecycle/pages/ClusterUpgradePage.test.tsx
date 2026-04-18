import { screen } from '@testing-library/react';
import { ClusterUpgradePage } from '@/features/cluster-lifecycle/pages/ClusterUpgradePage';
import { renderWithProviders } from '@/features/cluster-lifecycle/pages/testUtils';
import { useAuthStore } from '@/features/auth/store';
import { getClusterLifecycleDetail } from '@/services/clusterLifecycle';
import { installAntdDomShims } from '@/test/installAntdDomShims';

vi.mock('@/services/clusterLifecycle', async () => ({
  ...(await vi.importActual<typeof import('@/services/clusterLifecycle')>('@/services/clusterLifecycle')),
  getClusterLifecycleDetail: vi.fn(),
  createUpgradePlan: vi.fn(),
  executeUpgradePlan: vi.fn()
}));

describe('ClusterUpgradePage', () => {
  beforeAll(() => installAntdDomShims());
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(getClusterLifecycleDetail).mockResolvedValue({
      id: 'c1',
      name: 'prod-a',
      kubernetesVersion: 'v1.29.0',
      cluster: { id: 'c1', name: 'prod-a', kubernetesVersion: 'v1.29.0' },
      nodePools: [],
      upgradePlans: []
    });
  });

  it('shows unauthorized empty', () => {
    useAuthStore.setState({ isAuthenticated: true, accessToken: 't', refreshToken: 'r', user: { id: 'u', username: 'u', platformRoles: [] } });
    renderWithProviders(<ClusterUpgradePage />, ['/cluster-lifecycle/upgrade?clusterId=c1']);
    expect(screen.getByText('你暂无集群升级访问权限。')).toBeInTheDocument();
  });

  it('renders target cluster info', async () => {
    useAuthStore.setState({ isAuthenticated: true, accessToken: 't', refreshToken: 'r', user: { id: 'u', username: 'u', platformRoles: ['platform-admin'] } });
    renderWithProviders(<ClusterUpgradePage />, ['/cluster-lifecycle/upgrade?clusterId=c1']);
    expect(await screen.findByText(/目标集群/)).toBeInTheDocument();
  });
});
