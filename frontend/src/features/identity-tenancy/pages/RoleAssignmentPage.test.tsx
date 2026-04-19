import { screen, waitFor } from '@testing-library/react';
import { useAuthStore } from '@/features/auth/store';
import { RoleAssignmentPage } from '@/features/identity-tenancy/pages/RoleAssignmentPage';
import { renderWithProviders } from '@/features/identity-tenancy/pages/testUtils';
import { listRoleAssignments, listRoleDefinitions } from '@/services/identityTenancy';
import { installAntdDomShims } from '@/test/installAntdDomShims';

vi.mock('@/services/identityTenancy', async () => ({
  ...(await vi.importActual<typeof import('@/services/identityTenancy')>('@/services/identityTenancy')),
  listRoleAssignments: vi.fn(),
  listRoleDefinitions: vi.fn(),
  createRoleAssignment: vi.fn(),
  createRoleDefinition: vi.fn(),
  createDelegationGrant: vi.fn()
}));

describe('RoleAssignmentPage', () => {
  beforeAll(() => installAntdDomShims());

  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(listRoleAssignments).mockResolvedValue({ items: [] });
    vi.mocked(listRoleDefinitions).mockResolvedValue({ items: [] });
  });

  it('shows unauthorized state', () => {
    useAuthStore.setState({
      isAuthenticated: true,
      accessToken: 't',
      refreshToken: 'r',
      user: { id: 'u', username: 'u', platformRoles: [] }
    });
    renderWithProviders(<RoleAssignmentPage />);
    expect(screen.getByText('你暂无授权分配访问权限。')).toBeInTheDocument();
  });

  it('renders assignments', async () => {
    useAuthStore.setState({
      isAuthenticated: true,
      accessToken: 't',
      refreshToken: 'r',
      user: { id: 'u', username: 'u', platformRoles: ['platform-admin'] }
    });
    vi.mocked(listRoleAssignments).mockResolvedValue({
      items: [{ id: 'a1', subjectType: 'user', subjectRef: 'alice', roleDefinitionName: '项目管理员' }]
    });
    renderWithProviders(<RoleAssignmentPage />);
    await waitFor(() => expect(listRoleAssignments).toHaveBeenCalled());
    expect(await screen.findByText(/alice/)).toBeInTheDocument();
  });
});
