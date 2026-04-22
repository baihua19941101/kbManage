import { screen } from '@testing-library/react';
import { useAuthStore } from '@/features/auth/store';
import { GovernanceReportsPage } from '@/features/enterprise-polish/pages/GovernanceReportsPage';
import { renderWithProviders } from '@/features/enterprise-polish/pages/testUtils';
import { listGovernanceReports } from '@/services/enterprisePolish';

vi.mock('@/services/enterprisePolish', async () => ({
  ...(await vi.importActual<typeof import('@/services/enterprisePolish')>('@/services/enterprisePolish')),
  listGovernanceReports: vi.fn()
}));

describe('GovernanceReportsPage', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(listGovernanceReports).mockResolvedValue({ items: [] });
  });
  it('renders title', async () => {
    useAuthStore.setState({ isAuthenticated: true, accessToken: 't', refreshToken: 'r', user: { id: 'u', username: 'u', platformRoles: ['platform-admin'] } });
    renderWithProviders(<GovernanceReportsPage />);
    expect(await screen.findByText('治理报表中心')).toBeInTheDocument();
  });
});
