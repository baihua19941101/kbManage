import { screen } from '@testing-library/react';
import { useAuthStore } from '@/features/auth/store';
import { DeliveryReadinessPage } from '@/features/enterprise-polish/pages/DeliveryReadinessPage';
import { renderWithProviders } from '@/features/enterprise-polish/pages/testUtils';
import { listDeliveryBundles, listDeliveryChecklist } from '@/services/enterprisePolish';

vi.mock('@/services/enterprisePolish', async () => ({
  ...(await vi.importActual<typeof import('@/services/enterprisePolish')>('@/services/enterprisePolish')),
  listDeliveryBundles: vi.fn(),
  listDeliveryChecklist: vi.fn()
}));

describe('DeliveryReadinessPage', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(listDeliveryBundles).mockResolvedValue({ items: [{ id: '1', name: 'bundle', missingItems: [] }] });
    vi.mocked(listDeliveryChecklist).mockResolvedValue({ items: [] });
  });
  it('renders title', async () => {
    useAuthStore.setState({ isAuthenticated: true, accessToken: 't', refreshToken: 'r', user: { id: 'u', username: 'u', platformRoles: ['platform-admin'] } });
    renderWithProviders(<DeliveryReadinessPage />);
    expect(await screen.findByText('交付就绪检查')).toBeInTheDocument();
  });
});
