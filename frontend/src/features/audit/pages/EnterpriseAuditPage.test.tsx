import { screen } from '@testing-library/react';
import { useAuthStore } from '@/features/auth/store';
import { EnterpriseAuditPage } from '@/features/audit/pages/EnterpriseAuditPage';
import { renderWithProviders } from '@/features/enterprise-polish/pages/testUtils';
import { listEnterpriseAuditEvents } from '@/services/enterprisePolish';

vi.mock('@/services/enterprisePolish', async () => ({
  ...(await vi.importActual<typeof import('@/services/enterprisePolish')>('@/services/enterprisePolish')),
  listEnterpriseAuditEvents: vi.fn()
}));

describe('EnterpriseAuditPage', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(listEnterpriseAuditEvents).mockResolvedValue({ items: [] });
  });
  it('renders title', async () => {
    useAuthStore.setState({ isAuthenticated: true, accessToken: 't', refreshToken: 'r', user: { id: 'u', username: 'u', platformRoles: ['platform-admin'] } });
    renderWithProviders(<EnterpriseAuditPage />);
    expect(await screen.findByText('企业治理审计')).toBeInTheDocument();
  });
});
