import { screen } from '@testing-library/react';
import { useAuthStore } from '@/features/auth/store';
import { RunbookCenterPage } from '@/features/sre-scale/pages/RunbookCenterPage';
import { renderWithProviders } from '@/features/sre-scale/pages/testUtils';
import { listRunbooks } from '@/services/sreScale';

vi.mock('@/services/sreScale', async () => ({
  ...(await vi.importActual<typeof import('@/services/sreScale')>('@/services/sreScale')),
  listRunbooks: vi.fn()
}));

describe('RunbookCenterPage', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(listRunbooks).mockResolvedValue({ items: [] });
  });

  it('renders runbook title', async () => {
    useAuthStore.setState({ isAuthenticated: true, accessToken: 't', refreshToken: 'r', user: { id: 'u', username: 'u', platformRoles: ['platform-admin'] } });
    renderWithProviders(<RunbookCenterPage />);
    expect(await screen.findByText('运行手册中心')).toBeInTheDocument();
  });
});
