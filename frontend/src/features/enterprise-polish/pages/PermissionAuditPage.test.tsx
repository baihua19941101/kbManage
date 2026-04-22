import { screen } from '@testing-library/react';
import { useAuthStore } from '@/features/auth/store';
import { PermissionAuditPage } from '@/features/enterprise-polish/pages/PermissionAuditPage';
import { renderWithProviders } from '@/features/enterprise-polish/pages/testUtils';
import { listPermissionTrails } from '@/services/enterprisePolish';

vi.mock('@/services/enterprisePolish', async () => ({
  ...(await vi.importActual<typeof import('@/services/enterprisePolish')>('@/services/enterprisePolish')),
  listPermissionTrails: vi.fn()
}));

describe('PermissionAuditPage', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(listPermissionTrails).mockResolvedValue({ items: [] });
  });
  it('renders title', async () => {
    useAuthStore.setState({ isAuthenticated: true, accessToken: 't', refreshToken: 'r', user: { id: 'u', username: 'u', platformRoles: ['platform-admin'] } });
    renderWithProviders(<PermissionAuditPage />);
    expect(await screen.findByText('权限变更审计')).toBeInTheDocument();
  });
});
