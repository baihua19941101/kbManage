import { screen } from '@testing-library/react';
import { useAuthStore } from '@/features/auth/store';
import { HealthOverviewPage } from '@/features/sre-scale/pages/HealthOverviewPage';
import { renderWithProviders } from '@/features/sre-scale/pages/testUtils';
import { getHealthOverview } from '@/services/sreScale';

vi.mock('@/services/sreScale', async () => ({
  ...(await vi.importActual<typeof import('@/services/sreScale')>('@/services/sreScale')),
  getHealthOverview: vi.fn()
}));

describe('HealthOverviewPage', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(getHealthOverview).mockResolvedValue({ id: 'h1', overallStatus: 'healthy', recommendedActions: ['继续观察'] });
  });

  it('renders overview title', async () => {
    useAuthStore.setState({ isAuthenticated: true, accessToken: 't', refreshToken: 'r', user: { id: 'u', username: 'u', platformRoles: ['platform-admin'] } });
    renderWithProviders(<HealthOverviewPage />);
    expect(await screen.findByText('平台健康总览')).toBeInTheDocument();
  });
});
