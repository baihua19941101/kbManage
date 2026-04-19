import { screen, waitFor } from '@testing-library/react';
import { useAuthStore } from '@/features/auth/store';
import { AccessRiskPage } from '@/features/identity-tenancy/pages/AccessRiskPage';
import { renderWithProviders } from '@/features/identity-tenancy/pages/testUtils';
import { listAccessRisks } from '@/services/identityTenancy';
import { installAntdDomShims } from '@/test/installAntdDomShims';

vi.mock('@/services/identityTenancy', async () => ({
  ...(await vi.importActual<typeof import('@/services/identityTenancy')>('@/services/identityTenancy')),
  listAccessRisks: vi.fn()
}));

describe('AccessRiskPage', () => {
  beforeAll(() => installAntdDomShims());

  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(listAccessRisks).mockResolvedValue({ items: [] });
  });

  it('shows unauthorized state', () => {
    useAuthStore.setState({
      isAuthenticated: true,
      accessToken: 't',
      refreshToken: 'r',
      user: { id: 'u', username: 'u', platformRoles: [] }
    });
    renderWithProviders(<AccessRiskPage />);
    expect(screen.getByText('你暂无访问风险视图访问权限。')).toBeInTheDocument();
  });

  it('renders risk summary', async () => {
    useAuthStore.setState({
      isAuthenticated: true,
      accessToken: 't',
      refreshToken: 'r',
      user: { id: 'u', username: 'u', platformRoles: ['platform-admin'] }
    });
    vi.mocked(listAccessRisks).mockResolvedValue({
      items: [{ id: 'r1', subjectType: 'user', subjectRef: 'alice', severity: 'critical', summary: '跨租户权限扩散' }]
    });
    renderWithProviders(<AccessRiskPage />);
    await waitFor(() => expect(listAccessRisks).toHaveBeenCalled());
    expect(await screen.findByText('跨租户权限扩散')).toBeInTheDocument();
  });
});
