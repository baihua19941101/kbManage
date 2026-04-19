import { screen, waitFor } from '@testing-library/react';
import { useAuthStore } from '@/features/auth/store';
import { RoleCatalogPage } from '@/features/identity-tenancy/pages/RoleCatalogPage';
import { renderWithProviders } from '@/features/identity-tenancy/pages/testUtils';
import { listRoleDefinitions } from '@/services/identityTenancy';
import { installAntdDomShims } from '@/test/installAntdDomShims';

vi.mock('@/services/identityTenancy', async () => ({
  ...(await vi.importActual<typeof import('@/services/identityTenancy')>('@/services/identityTenancy')),
  listRoleDefinitions: vi.fn(),
  createRoleDefinition: vi.fn()
}));

describe('RoleCatalogPage', () => {
  beforeAll(() => installAntdDomShims());

  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(listRoleDefinitions).mockResolvedValue({ items: [] });
  });

  it('shows unauthorized state', () => {
    useAuthStore.setState({
      isAuthenticated: true,
      accessToken: 't',
      refreshToken: 'r',
      user: { id: 'u', username: 'u', platformRoles: [] }
    });
    renderWithProviders(<RoleCatalogPage />);
    expect(screen.getByText('你暂无角色目录访问权限。')).toBeInTheDocument();
  });

  it('renders role catalog', async () => {
    useAuthStore.setState({
      isAuthenticated: true,
      accessToken: 't',
      refreshToken: 'r',
      user: { id: 'u', username: 'u', platformRoles: ['platform-admin'] }
    });
    vi.mocked(listRoleDefinitions).mockResolvedValue({
      items: [{ id: 'r1', name: '组织安全管理员', roleLevel: 'organization', delegable: true }]
    });
    renderWithProviders(<RoleCatalogPage />);
    await waitFor(() => expect(listRoleDefinitions).toHaveBeenCalled());
    expect(await screen.findByText('组织安全管理员')).toBeInTheDocument();
  });
});
