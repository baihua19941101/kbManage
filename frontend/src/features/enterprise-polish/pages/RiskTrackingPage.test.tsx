import { screen } from '@testing-library/react';
import { useAuthStore } from '@/features/auth/store';
import { RiskTrackingPage } from '@/features/enterprise-polish/pages/RiskTrackingPage';
import { renderWithProviders } from '@/features/enterprise-polish/pages/testUtils';
import { listKeyOperations } from '@/services/enterprisePolish';

vi.mock('@/services/enterprisePolish', async () => ({
  ...(await vi.importActual<typeof import('@/services/enterprisePolish')>('@/services/enterprisePolish')),
  listKeyOperations: vi.fn()
}));

describe('RiskTrackingPage', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(listKeyOperations).mockResolvedValue({ items: [] });
  });
  it('renders title', async () => {
    useAuthStore.setState({ isAuthenticated: true, accessToken: 't', refreshToken: 'r', user: { id: 'u', username: 'u', platformRoles: ['platform-admin'] } });
    renderWithProviders(<RiskTrackingPage />);
    expect(await screen.findByText('高风险访问追踪')).toBeInTheDocument();
  });
});
