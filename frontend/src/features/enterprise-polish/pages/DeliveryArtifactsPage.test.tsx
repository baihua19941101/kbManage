import { screen } from '@testing-library/react';
import { useAuthStore } from '@/features/auth/store';
import { DeliveryArtifactsPage } from '@/features/enterprise-polish/pages/DeliveryArtifactsPage';
import { renderWithProviders } from '@/features/enterprise-polish/pages/testUtils';
import { listDeliveryArtifacts } from '@/services/enterprisePolish';

vi.mock('@/services/enterprisePolish', async () => ({
  ...(await vi.importActual<typeof import('@/services/enterprisePolish')>('@/services/enterprisePolish')),
  listDeliveryArtifacts: vi.fn()
}));

describe('DeliveryArtifactsPage', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(listDeliveryArtifacts).mockResolvedValue({ items: [] });
  });
  it('renders title', async () => {
    useAuthStore.setState({ isAuthenticated: true, accessToken: 't', refreshToken: 'r', user: { id: 'u', username: 'u', platformRoles: ['platform-admin'] } });
    renderWithProviders(<DeliveryArtifactsPage />);
    expect(await screen.findByText('交付材料目录')).toBeInTheDocument();
  });
});
