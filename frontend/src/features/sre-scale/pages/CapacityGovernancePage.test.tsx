import { screen } from '@testing-library/react';
import { useAuthStore } from '@/features/auth/store';
import { CapacityGovernancePage } from '@/features/sre-scale/pages/CapacityGovernancePage';
import { renderWithProviders } from '@/features/sre-scale/pages/testUtils';
import { listCapacityBaselines, listScaleEvidence } from '@/services/sreScale';

vi.mock('@/services/sreScale', async () => ({
  ...(await vi.importActual<typeof import('@/services/sreScale')>('@/services/sreScale')),
  listCapacityBaselines: vi.fn(),
  listScaleEvidence: vi.fn()
}));

describe('CapacityGovernancePage', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(listCapacityBaselines).mockResolvedValue({ items: [] });
    vi.mocked(listScaleEvidence).mockResolvedValue({ items: [] });
  });

  it('renders capacity title', async () => {
    useAuthStore.setState({ isAuthenticated: true, accessToken: 't', refreshToken: 'r', user: { id: 'u', username: 'u', platformRoles: ['platform-admin'] } });
    renderWithProviders(<CapacityGovernancePage />);
    expect(await screen.findByText('容量与性能治理')).toBeInTheDocument();
  });
});
