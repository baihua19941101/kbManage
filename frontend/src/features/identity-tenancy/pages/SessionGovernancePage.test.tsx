import { screen, waitFor } from '@testing-library/react';
import { useAuthStore } from '@/features/auth/store';
import { SessionGovernancePage } from '@/features/identity-tenancy/pages/SessionGovernancePage';
import { renderWithProviders } from '@/features/identity-tenancy/pages/testUtils';
import { listSessionRecords } from '@/services/identityTenancy';
import { installAntdDomShims } from '@/test/installAntdDomShims';

vi.mock('@/services/identityTenancy', async () => ({
  ...(await vi.importActual<typeof import('@/services/identityTenancy')>('@/services/identityTenancy')),
  listSessionRecords: vi.fn(),
  updatePreferredLoginMode: vi.fn(),
  revokeSessionRecord: vi.fn()
}));

describe('SessionGovernancePage', () => {
  beforeAll(() => installAntdDomShims());

  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(listSessionRecords).mockResolvedValue({ items: [] });
  });

  it('shows unauthorized state', () => {
    useAuthStore.setState({
      isAuthenticated: true,
      accessToken: 't',
      refreshToken: 'r',
      user: { id: 'u', username: 'u', platformRoles: [] }
    });
    renderWithProviders(<SessionGovernancePage />);
    expect(screen.getByText('你暂无会话治理访问权限。')).toBeInTheDocument();
  });

  it('renders session list', async () => {
    useAuthStore.setState({
      isAuthenticated: true,
      accessToken: 't',
      refreshToken: 'r',
      user: { id: 'u', username: 'u', platformRoles: ['platform-admin'] }
    });
    vi.mocked(listSessionRecords).mockResolvedValue({
      items: [{ id: 'ss1', username: 'alice', identitySourceId: 'oidc-main', status: 'active', riskLevel: 'low' }]
    });
    renderWithProviders(<SessionGovernancePage />);
    await waitFor(() => expect(listSessionRecords).toHaveBeenCalled());
    expect(await screen.findByText('alice')).toBeInTheDocument();
  });
});
