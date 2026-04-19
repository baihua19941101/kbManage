import { screen, waitFor } from '@testing-library/react';
import { useAuthStore } from '@/features/auth/store';
import { OrganizationModelPage } from '@/features/identity-tenancy/pages/OrganizationModelPage';
import { renderWithProviders } from '@/features/identity-tenancy/pages/testUtils';
import { listOrganizationUnits } from '@/services/identityTenancy';
import { installAntdDomShims } from '@/test/installAntdDomShims';

vi.mock('@/services/identityTenancy', async () => ({
  ...(await vi.importActual<typeof import('@/services/identityTenancy')>('@/services/identityTenancy')),
  listOrganizationUnits: vi.fn(),
  createOrganizationUnit: vi.fn(),
  createTenantScopeMapping: vi.fn()
}));

describe('OrganizationModelPage', () => {
  beforeAll(() => installAntdDomShims());

  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(listOrganizationUnits).mockResolvedValue({ items: [] });
  });

  it('shows unauthorized state', () => {
    useAuthStore.setState({
      isAuthenticated: true,
      accessToken: 't',
      refreshToken: 'r',
      user: { id: 'u', username: 'u', platformRoles: [] }
    });
    renderWithProviders(<OrganizationModelPage />);
    expect(screen.getByText('你暂无组织模型访问权限。')).toBeInTheDocument();
  });

  it('renders organization units', async () => {
    useAuthStore.setState({
      isAuthenticated: true,
      accessToken: 't',
      refreshToken: 'r',
      user: { id: 'u', username: 'u', platformRoles: ['platform-admin'] }
    });
    vi.mocked(listOrganizationUnits).mockResolvedValue({
      items: [{ id: 'org-1', name: '零售事业部', unitType: 'organization', status: 'active' }]
    });
    renderWithProviders(<OrganizationModelPage />);
    await waitFor(() => expect(listOrganizationUnits).toHaveBeenCalled());
    expect(await screen.findByText('零售事业部')).toBeInTheDocument();
  });
});
