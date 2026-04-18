import { screen, waitFor } from '@testing-library/react';
import { CapabilityMatrixPage } from '@/features/cluster-lifecycle/pages/CapabilityMatrixPage';
import { renderWithProviders } from '@/features/cluster-lifecycle/pages/testUtils';
import { useAuthStore } from '@/features/auth/store';
import { listClusterDrivers } from '@/services/clusterLifecycle';
import { installAntdDomShims } from '@/test/installAntdDomShims';

vi.mock('@/services/clusterLifecycle', async () => ({
  ...(await vi.importActual<typeof import('@/services/clusterLifecycle')>('@/services/clusterLifecycle')),
  listClusterDrivers: vi.fn(),
  listDriverCapabilities: vi.fn()
}));

describe('CapabilityMatrixPage', () => {
  beforeAll(() => installAntdDomShims());
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(listClusterDrivers).mockResolvedValue({ items: [{ id: 'd1', driverKey: 'rke2', version: '1.0.0' }] });
  });

  it('shows unauthorized empty', () => {
    useAuthStore.setState({ isAuthenticated: true, accessToken: 't', refreshToken: 'r', user: { id: 'u', username: 'u', platformRoles: [] } });
    renderWithProviders(<CapabilityMatrixPage />);
    expect(screen.getByText('你暂无能力矩阵访问权限。')).toBeInTheDocument();
  });

  it('renders driver selector', async () => {
    useAuthStore.setState({ isAuthenticated: true, accessToken: 't', refreshToken: 'r', user: { id: 'u', username: 'u', platformRoles: ['platform-admin'] } });
    renderWithProviders(<CapabilityMatrixPage />);
    await waitFor(() => expect(listClusterDrivers).toHaveBeenCalled());
    expect(await screen.findByText('选择驱动版本')).toBeInTheDocument();
  });
});
