import { screen } from '@testing-library/react';
import { useAuthStore } from '@/features/auth/store';
import { HAControlPage } from '@/features/sre-scale/pages/HAControlPage';
import { renderWithProviders } from '@/features/sre-scale/pages/testUtils';
import { listHAPolicies, listMaintenanceWindows } from '@/services/sreScale';
import { installAntdDomShims } from '@/test/installAntdDomShims';

vi.mock('@/services/sreScale', async () => ({
  ...(await vi.importActual<typeof import('@/services/sreScale')>('@/services/sreScale')),
  listHAPolicies: vi.fn(),
  listMaintenanceWindows: vi.fn()
}));

describe('HAControlPage', () => {
  beforeAll(() => installAntdDomShims());
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(listHAPolicies).mockResolvedValue({ items: [] });
    vi.mocked(listMaintenanceWindows).mockResolvedValue({ items: [] });
  });

  it('renders page title', async () => {
    useAuthStore.setState({ isAuthenticated: true, accessToken: 't', refreshToken: 'r', user: { id: 'u', username: 'u', platformRoles: ['platform-admin'] } });
    renderWithProviders(<HAControlPage />);
    expect(await screen.findByText('高可用治理')).toBeInTheDocument();
  });
});
