import { screen, waitFor } from '@testing-library/react';
import { useAuthStore } from '@/features/auth/store';
import { TenantMappingPage } from '@/features/identity-tenancy/pages/TenantMappingPage';
import { renderWithProviders } from '@/features/identity-tenancy/pages/testUtils';
import { listOrganizationUnits } from '@/services/identityTenancy';
import { installAntdDomShims } from '@/test/installAntdDomShims';

vi.mock('@/services/identityTenancy', async () => ({
  ...(await vi.importActual<typeof import('@/services/identityTenancy')>('@/services/identityTenancy')),
  listOrganizationUnits: vi.fn(),
  listTenantScopeMappings: vi.fn(),
  createOrganizationUnit: vi.fn(),
  createTenantScopeMapping: vi.fn()
}));

describe('TenantMappingPage', () => {
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
    renderWithProviders(<TenantMappingPage />);
    expect(screen.getByText('你暂无租户边界映射访问权限。')).toBeInTheDocument();
  });

  it('renders unit selector', async () => {
    useAuthStore.setState({
      isAuthenticated: true,
      accessToken: 't',
      refreshToken: 'r',
      user: { id: 'u', username: 'u', platformRoles: ['platform-admin'] }
    });
    vi.mocked(listOrganizationUnits).mockResolvedValue({
      items: [{ id: 'team-1', name: '平台团队', unitType: 'team', status: 'active' }]
    });
    renderWithProviders(<TenantMappingPage />);
    await waitFor(() => expect(listOrganizationUnits).toHaveBeenCalled());
    expect(await screen.findByText('选择组织单元')).toBeInTheDocument();
  });
});
