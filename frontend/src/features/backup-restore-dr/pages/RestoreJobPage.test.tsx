import { screen } from '@testing-library/react';
import { RestoreJobPage } from '@/features/backup-restore-dr/pages/RestoreJobPage';
import { renderWithProviders } from '@/features/backup-restore-dr/pages/testUtils';
import { useAuthStore } from '@/features/auth/store';
import { listRestoreJobs, listRestorePoints } from '@/services/backupRestore';
import { installAntdDomShims } from '@/test/installAntdDomShims';

vi.mock('@/services/backupRestore', async () => ({
  ...(await vi.importActual<typeof import('@/services/backupRestore')>('@/services/backupRestore')),
  listRestoreJobs: vi.fn(),
  listRestorePoints: vi.fn(),
  createRestoreJob: vi.fn(),
  validateRestoreJob: vi.fn(),
  createMigrationPlan: vi.fn(),
  runBackupPolicy: vi.fn()
}));

describe('RestoreJobPage', () => {
  beforeAll(() => installAntdDomShims());

  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(listRestoreJobs).mockResolvedValue({ items: [] });
    vi.mocked(listRestorePoints).mockResolvedValue({ items: [] });
  });

  it('shows unauthorized state', () => {
    useAuthStore.setState({
      isAuthenticated: true,
      accessToken: 't',
      refreshToken: 'r',
      user: { id: 'u', username: 'u', platformRoles: [] }
    });
    renderWithProviders(<RestoreJobPage />);
    expect(screen.getByText('你暂无恢复任务访问权限。')).toBeInTheDocument();
  });

  it('renders restore jobs', async () => {
    useAuthStore.setState({
      isAuthenticated: true,
      accessToken: 't',
      refreshToken: 'r',
      user: { id: 'u', username: 'u', platformRoles: ['platform-admin'] }
    });
    vi.mocked(listRestoreJobs).mockResolvedValue({
      items: [{ id: 'job-1', jobType: 'cross-cluster', targetEnvironment: 'dr-cn', status: 'running' }]
    });
    renderWithProviders(<RestoreJobPage />);
    expect(await screen.findByText('job-1')).toBeInTheDocument();
  });
});
