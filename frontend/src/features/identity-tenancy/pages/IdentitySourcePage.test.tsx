import { screen, waitFor } from '@testing-library/react';
import { useAuthStore } from '@/features/auth/store';
import { IdentitySourcePage } from '@/features/identity-tenancy/pages/IdentitySourcePage';
import { renderWithProviders } from '@/features/identity-tenancy/pages/testUtils';
import { listIdentitySources } from '@/services/identityTenancy';
import { installAntdDomShims } from '@/test/installAntdDomShims';

vi.mock('@/services/identityTenancy', async () => ({
  ...(await vi.importActual<typeof import('@/services/identityTenancy')>('@/services/identityTenancy')),
  listIdentitySources: vi.fn(),
  createIdentitySource: vi.fn()
}));

describe('IdentitySourcePage', () => {
  beforeAll(() => installAntdDomShims());

  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(listIdentitySources).mockResolvedValue({ items: [] });
  });

  it('shows unauthorized state', () => {
    useAuthStore.setState({
      isAuthenticated: true,
      accessToken: 't',
      refreshToken: 'r',
      user: { id: 'u', username: 'u', platformRoles: [] }
    });
    renderWithProviders(<IdentitySourcePage />);
    expect(screen.getByText('你暂无身份源治理访问权限。')).toBeInTheDocument();
  });

  it('renders identity source list', async () => {
    useAuthStore.setState({
      isAuthenticated: true,
      accessToken: 't',
      refreshToken: 'r',
      user: { id: 'u', username: 'u', platformRoles: ['platform-admin'] }
    });
    vi.mocked(listIdentitySources).mockResolvedValue({
      items: [{ id: 's1', name: '企业 OIDC', sourceType: 'oidc', loginMode: 'mixed', status: 'active' }]
    });
    renderWithProviders(<IdentitySourcePage />);
    await waitFor(() => expect(listIdentitySources).toHaveBeenCalled());
    expect(await screen.findByText('企业 OIDC')).toBeInTheDocument();
  });
});
