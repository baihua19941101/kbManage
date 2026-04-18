import { screen, waitFor } from '@testing-library/react';
import { NodePoolPage } from '@/features/cluster-lifecycle/pages/NodePoolPage';
import { renderWithProviders } from '@/features/cluster-lifecycle/pages/testUtils';
import { useAuthStore } from '@/features/auth/store';
import { listNodePools } from '@/services/clusterLifecycle';
import { installAntdDomShims } from '@/test/installAntdDomShims';

vi.mock('@/services/clusterLifecycle', async () => ({
  ...(await vi.importActual<typeof import('@/services/clusterLifecycle')>('@/services/clusterLifecycle')),
  listNodePools: vi.fn(),
  scaleNodePool: vi.fn()
}));

describe('NodePoolPage', () => {
  beforeAll(() => installAntdDomShims());
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(listNodePools).mockResolvedValue({ items: [{ id: 'np1', name: 'worker-a', desiredCount: 3, currentCount: 2, status: 'active' }] });
  });

  it('shows unauthorized empty', () => {
    useAuthStore.setState({ isAuthenticated: true, accessToken: 't', refreshToken: 'r', user: { id: 'u', username: 'u', platformRoles: [] } });
    renderWithProviders(<NodePoolPage />, ['/cluster-lifecycle/node-pools?clusterId=c1']);
    expect(screen.getByText('你暂无节点池管理访问权限。')).toBeInTheDocument();
  });

  it('renders node pools', async () => {
    useAuthStore.setState({ isAuthenticated: true, accessToken: 't', refreshToken: 'r', user: { id: 'u', username: 'u', platformRoles: ['platform-admin'] } });
    renderWithProviders(<NodePoolPage />, ['/cluster-lifecycle/node-pools?clusterId=c1']);
    await waitFor(() => expect(listNodePools).toHaveBeenCalled());
    expect(await screen.findByText('worker-a')).toBeInTheDocument();
  });
});
