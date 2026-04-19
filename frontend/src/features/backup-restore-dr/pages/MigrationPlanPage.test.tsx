import { screen } from '@testing-library/react';
import { MigrationPlanPage } from '@/features/backup-restore-dr/pages/MigrationPlanPage';
import { renderWithProviders } from '@/features/backup-restore-dr/pages/testUtils';
import { useAuthStore } from '@/features/auth/store';
import { installAntdDomShims } from '@/test/installAntdDomShims';

describe('MigrationPlanPage', () => {
  beforeAll(() => installAntdDomShims());

  it('shows unauthorized state', () => {
    useAuthStore.setState({
      isAuthenticated: true,
      accessToken: 't',
      refreshToken: 'r',
      user: { id: 'u', username: 'u', platformRoles: [] }
    });
    renderWithProviders(<MigrationPlanPage />);
    expect(screen.getByText('你暂无迁移计划访问权限。')).toBeInTheDocument();
  });

  it('renders page shell for operator', () => {
    useAuthStore.setState({
      isAuthenticated: true,
      accessToken: 't',
      refreshToken: 'r',
      user: { id: 'u', username: 'u', platformRoles: ['ops-operator'] }
    });
    renderWithProviders(<MigrationPlanPage />);
    expect(screen.getByText('迁移计划中心')).toBeInTheDocument();
  });
});
