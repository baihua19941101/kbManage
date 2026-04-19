import { screen } from '@testing-library/react';
import { useAuthStore } from '@/features/auth/store';
import { RollbackValidationPage } from '@/features/sre-scale/pages/RollbackValidationPage';
import { renderWithProviders } from '@/features/sre-scale/pages/testUtils';
import { listUpgradePlans } from '@/services/sreScale';

vi.mock('@/services/sreScale', async () => ({
  ...(await vi.importActual<typeof import('@/services/sreScale')>('@/services/sreScale')),
  listUpgradePlans: vi.fn()
}));

describe('RollbackValidationPage', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(listUpgradePlans).mockResolvedValue({ items: [] });
  });

  it('renders rollback title', async () => {
    useAuthStore.setState({ isAuthenticated: true, accessToken: 't', refreshToken: 'r', user: { id: 'u', username: 'u', platformRoles: ['platform-admin'] } });
    renderWithProviders(<RollbackValidationPage />);
    expect(await screen.findByText('回退验证')).toBeInTheDocument();
  });
});
