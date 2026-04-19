import { screen } from '@testing-library/react';
import { BackupPolicyPage } from '@/features/backup-restore-dr/pages/BackupPolicyPage';
import { renderWithProviders } from '@/features/backup-restore-dr/pages/testUtils';
import { useAuthStore } from '@/features/auth/store';
import { listBackupPolicies } from '@/services/backupRestore';
import { installAntdDomShims } from '@/test/installAntdDomShims';

vi.mock('@/services/backupRestore', async () => ({
  ...(await vi.importActual<typeof import('@/services/backupRestore')>('@/services/backupRestore')),
  listBackupPolicies: vi.fn(),
  createBackupPolicy: vi.fn(),
  runBackupPolicy: vi.fn()
}));

describe('BackupPolicyPage', () => {
  beforeAll(() => installAntdDomShims());

  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(listBackupPolicies).mockResolvedValue({ items: [] });
  });

  it('shows unauthorized state', () => {
    useAuthStore.setState({
      isAuthenticated: true,
      accessToken: 't',
      refreshToken: 'r',
      user: { id: 'u', username: 'u', platformRoles: [] }
    });
    renderWithProviders(<BackupPolicyPage />);
    expect(screen.getByText('你暂无备份策略访问权限。')).toBeInTheDocument();
  });

  it('renders policy list', async () => {
    useAuthStore.setState({
      isAuthenticated: true,
      accessToken: 't',
      refreshToken: 'r',
      user: { id: 'u', username: 'u', platformRoles: ['platform-admin'] }
    });
    vi.mocked(listBackupPolicies).mockResolvedValue({
      items: [{ id: 'policy-1', name: '核心元数据备份', scopeType: 'platform-metadata', status: 'active' }]
    });
    renderWithProviders(<BackupPolicyPage />);
    expect(await screen.findByText('核心元数据备份')).toBeInTheDocument();
  });
});
