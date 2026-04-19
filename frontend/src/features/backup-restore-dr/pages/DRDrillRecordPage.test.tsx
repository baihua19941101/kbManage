import { screen } from '@testing-library/react';
import { DRDrillRecordPage } from '@/features/backup-restore-dr/pages/DRDrillRecordPage';
import { renderWithProviders } from '@/features/backup-restore-dr/pages/testUtils';
import { useAuthStore } from '@/features/auth/store';
import { getDrillRecordDetail } from '@/services/backupRestore';
import { installAntdDomShims } from '@/test/installAntdDomShims';

vi.mock('@/services/backupRestore', async () => ({
  ...(await vi.importActual<typeof import('@/services/backupRestore')>('@/services/backupRestore')),
  getDrillRecordDetail: vi.fn()
}));

describe('DRDrillRecordPage', () => {
  beforeAll(() => installAntdDomShims());

  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(getDrillRecordDetail).mockResolvedValue({
      id: 'record-001',
      status: 'completed',
      stepResults: ['切换入口完成'],
      validationResults: ['订单读写验证通过']
    });
  });

  it('shows unauthorized state', () => {
    useAuthStore.setState({
      isAuthenticated: true,
      accessToken: 't',
      refreshToken: 'r',
      user: { id: 'u', username: 'u', platformRoles: [] }
    });
    renderWithProviders(<DRDrillRecordPage />);
    expect(screen.getByText('你暂无灾备演练记录访问权限。')).toBeInTheDocument();
  });

  it('renders drill record details', async () => {
    useAuthStore.setState({
      isAuthenticated: true,
      accessToken: 't',
      refreshToken: 'r',
      user: { id: 'u', username: 'u', platformRoles: ['platform-admin'] }
    });
    renderWithProviders(<DRDrillRecordPage />);
    expect(await screen.findByText('切换入口完成')).toBeInTheDocument();
  });
});
