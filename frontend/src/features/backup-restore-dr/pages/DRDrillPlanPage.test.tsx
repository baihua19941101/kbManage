import { screen } from '@testing-library/react';
import { DRDrillPlanPage } from '@/features/backup-restore-dr/pages/DRDrillPlanPage';
import { renderWithProviders } from '@/features/backup-restore-dr/pages/testUtils';
import { useAuthStore } from '@/features/auth/store';
import { listDrillPlans } from '@/services/backupRestore';
import { installAntdDomShims } from '@/test/installAntdDomShims';

vi.mock('@/services/backupRestore', async () => ({
  ...(await vi.importActual<typeof import('@/services/backupRestore')>('@/services/backupRestore')),
  listDrillPlans: vi.fn(),
  createDrillPlan: vi.fn(),
  createDrillReport: vi.fn(),
  runDrillPlan: vi.fn()
}));

describe('DRDrillPlanPage', () => {
  beforeAll(() => installAntdDomShims());

  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(listDrillPlans).mockResolvedValue({ items: [] });
  });

  it('shows unauthorized state', () => {
    useAuthStore.setState({
      isAuthenticated: true,
      accessToken: 't',
      refreshToken: 'r',
      user: { id: 'u', username: 'u', platformRoles: [] }
    });
    renderWithProviders(<DRDrillPlanPage />);
    expect(screen.getByText('你暂无灾备演练计划访问权限。')).toBeInTheDocument();
  });

  it('renders drill plan list', async () => {
    useAuthStore.setState({
      isAuthenticated: true,
      accessToken: 't',
      refreshToken: 'r',
      user: { id: 'u', username: 'u', platformRoles: ['platform-admin'] }
    });
    vi.mocked(listDrillPlans).mockResolvedValue({
      items: [{ id: 'plan-1', name: '季度演练', status: 'scheduled', roleAssignments: [], cutoverProcedure: [], validationChecklist: [] }]
    });
    renderWithProviders(<DRDrillPlanPage />);
    expect(await screen.findByText('季度演练')).toBeInTheDocument();
  });
});
