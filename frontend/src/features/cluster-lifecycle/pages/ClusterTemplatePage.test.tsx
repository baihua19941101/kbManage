import { screen, waitFor } from '@testing-library/react';
import { ClusterTemplatePage } from '@/features/cluster-lifecycle/pages/ClusterTemplatePage';
import { renderWithProviders } from '@/features/cluster-lifecycle/pages/testUtils';
import { useAuthStore } from '@/features/auth/store';
import { listClusterDrivers, listClusterTemplates } from '@/services/clusterLifecycle';
import { installAntdDomShims } from '@/test/installAntdDomShims';

vi.mock('@/services/clusterLifecycle', async () => ({
  ...(await vi.importActual<typeof import('@/services/clusterLifecycle')>('@/services/clusterLifecycle')),
  listClusterDrivers: vi.fn(),
  listClusterTemplates: vi.fn(),
  createClusterTemplate: vi.fn()
}));

describe('ClusterTemplatePage', () => {
  beforeAll(() => installAntdDomShims());
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(listClusterDrivers).mockResolvedValue({ items: [{ id: 'd1', driverKey: 'rke2', version: '1.0.0' }] });
    vi.mocked(listClusterTemplates).mockResolvedValue({ items: [{ id: 't1', name: 'prod-template', driverKey: 'rke2', infrastructureType: 'cloud', status: 'active' }] });
  });

  it('shows unauthorized empty', () => {
    useAuthStore.setState({ isAuthenticated: true, accessToken: 't', refreshToken: 'r', user: { id: 'u', username: 'u', platformRoles: [] } });
    renderWithProviders(<ClusterTemplatePage />);
    expect(screen.getByText('你暂无模板管理访问权限。')).toBeInTheDocument();
  });

  it('renders template table', async () => {
    useAuthStore.setState({ isAuthenticated: true, accessToken: 't', refreshToken: 'r', user: { id: 'u', username: 'u', platformRoles: ['platform-admin'] } });
    renderWithProviders(<ClusterTemplatePage />);
    await waitFor(() => expect(listClusterTemplates).toHaveBeenCalled());
    expect(await screen.findByText('prod-template')).toBeInTheDocument();
  });
});
