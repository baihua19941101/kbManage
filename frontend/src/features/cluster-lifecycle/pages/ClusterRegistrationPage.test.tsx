import { screen, waitFor } from '@testing-library/react';
import { ClusterRegistrationPage } from '@/features/cluster-lifecycle/pages/ClusterRegistrationPage';
import { renderWithProviders } from '@/features/cluster-lifecycle/pages/testUtils';
import { useAuthStore } from '@/features/auth/store';
import { listClusterLifecycleRecords, registerCluster } from '@/services/clusterLifecycle';
import { installAntdDomShims } from '@/test/installAntdDomShims';

vi.mock('@/services/clusterLifecycle', async () => ({
  ...(await vi.importActual<typeof import('@/services/clusterLifecycle')>('@/services/clusterLifecycle')),
  listClusterLifecycleRecords: vi.fn(),
  registerCluster: vi.fn()
}));

describe('ClusterRegistrationPage', () => {
  beforeAll(() => installAntdDomShims());

  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(listClusterLifecycleRecords).mockResolvedValue({ items: [] });
    vi.mocked(registerCluster).mockResolvedValue({ clusterId: 'c1', commandSnippet: 'kubectl apply -f agent.yaml' });
  });

  it('shows unauthorized empty', () => {
    useAuthStore.setState({ isAuthenticated: true, accessToken: 't', refreshToken: 'r', user: { id: 'u', username: 'u', platformRoles: [] } });
    renderWithProviders(<ClusterRegistrationPage />);
    expect(screen.getByText('你暂无集群注册访问权限。')).toBeInTheDocument();
  });

  it('loads pending registration summary', async () => {
    useAuthStore.setState({ isAuthenticated: true, accessToken: 't', refreshToken: 'r', user: { id: 'u', username: 'u', platformRoles: ['platform-admin'] } });
    renderWithProviders(<ClusterRegistrationPage />);
    await waitFor(() => expect(listClusterLifecycleRecords).toHaveBeenCalled());
    expect(screen.getByText(/待完成接入/)).toBeInTheDocument();
  });
});
