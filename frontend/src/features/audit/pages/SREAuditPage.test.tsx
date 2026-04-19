import { screen } from '@testing-library/react';
import { useAuthStore } from '@/features/auth/store';
import { SREAuditPage } from '@/features/audit/pages/SREAuditPage';
import { renderWithProviders } from '@/features/sre-scale/pages/testUtils';
import { listSREAuditEvents } from '@/services/sreScale';

vi.mock('@/services/sreScale', async () => ({
  ...(await vi.importActual<typeof import('@/services/sreScale')>('@/services/sreScale')),
  listSREAuditEvents: vi.fn()
}));

describe('SREAuditPage', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(listSREAuditEvents).mockResolvedValue({ items: [] });
  });

  it('renders audit title', async () => {
    useAuthStore.setState({ isAuthenticated: true, accessToken: 't', refreshToken: 'r', user: { id: 'u', username: 'u', platformRoles: ['platform-admin'] } });
    renderWithProviders(<SREAuditPage />);
    expect(await screen.findByText('平台 SRE 审计')).toBeInTheDocument();
  });
});
