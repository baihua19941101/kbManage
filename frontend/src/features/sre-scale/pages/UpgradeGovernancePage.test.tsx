import { screen } from '@testing-library/react';
import { useAuthStore } from '@/features/auth/store';
import { UpgradeGovernancePage } from '@/features/sre-scale/pages/UpgradeGovernancePage';
import { renderWithProviders } from '@/features/sre-scale/pages/testUtils';
import { listUpgradePlans } from '@/services/sreScale';

vi.mock('@/services/sreScale', async () => ({
  ...(await vi.importActual<typeof import('@/services/sreScale')>('@/services/sreScale')),
  listUpgradePlans: vi.fn()
}));

describe('UpgradeGovernancePage', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(listUpgradePlans).mockResolvedValue({ items: [] });
  });

  it('renders upgrade title', async () => {
    useAuthStore.setState({ isAuthenticated: true, accessToken: 't', refreshToken: 'r', user: { id: 'u', username: 'u', platformRoles: ['platform-admin'] } });
    renderWithProviders(<UpgradeGovernancePage />);
    expect(await screen.findByText('升级治理')).toBeInTheDocument();
  });
});
