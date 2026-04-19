import { screen } from '@testing-library/react';
import { RestorePointPage } from '@/features/backup-restore-dr/pages/RestorePointPage';
import { renderWithProviders } from '@/features/backup-restore-dr/pages/testUtils';
import { useAuthStore } from '@/features/auth/store';
import { listRestorePoints } from '@/services/backupRestore';
import { installAntdDomShims } from '@/test/installAntdDomShims';

vi.mock('@/services/backupRestore', async () => ({
  ...(await vi.importActual<typeof import('@/services/backupRestore')>('@/services/backupRestore')),
  listRestorePoints: vi.fn()
}));

describe('RestorePointPage', () => {
  beforeAll(() => installAntdDomShims());

  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(listRestorePoints).mockResolvedValue({ items: [] });
  });

  it('shows unauthorized state', () => {
    useAuthStore.setState({
      isAuthenticated: true,
      accessToken: 't',
      refreshToken: 'r',
      user: { id: 'u', username: 'u', platformRoles: [] }
    });
    renderWithProviders(<RestorePointPage />);
    expect(screen.getByText('你暂无恢复点访问权限。')).toBeInTheDocument();
  });

  it('renders restore point rows', async () => {
    useAuthStore.setState({
      isAuthenticated: true,
      accessToken: 't',
      refreshToken: 'r',
      user: { id: 'u', username: 'u', platformRoles: ['platform-admin'] }
    });
    vi.mocked(listRestorePoints).mockResolvedValue({
      items: [{ id: 'rp-1', policyId: 'policy-1', result: 'succeeded', durationSeconds: 18 }]
    });
    renderWithProviders(<RestorePointPage />);
    expect(await screen.findByText('rp-1')).toBeInTheDocument();
  });
});
