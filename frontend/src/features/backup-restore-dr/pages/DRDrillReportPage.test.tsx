import { screen } from '@testing-library/react';
import { DRDrillReportPage } from '@/features/backup-restore-dr/pages/DRDrillReportPage';
import { renderWithProviders } from '@/features/backup-restore-dr/pages/testUtils';
import { useAuthStore } from '@/features/auth/store';
import { installAntdDomShims } from '@/test/installAntdDomShims';

describe('DRDrillReportPage', () => {
  beforeAll(() => installAntdDomShims());

  it('shows unauthorized state', () => {
    useAuthStore.setState({
      isAuthenticated: true,
      accessToken: 't',
      refreshToken: 'r',
      user: { id: 'u', username: 'u', platformRoles: [] }
    });
    renderWithProviders(<DRDrillReportPage />);
    expect(screen.getByText('你暂无灾备演练报告访问权限。')).toBeInTheDocument();
  });

  it('renders page shell for operator', () => {
    useAuthStore.setState({
      isAuthenticated: true,
      accessToken: 't',
      refreshToken: 'r',
      user: { id: 'u', username: 'u', platformRoles: ['ops-operator'] }
    });
    renderWithProviders(<DRDrillReportPage />);
    expect(screen.getByText('灾备演练报告')).toBeInTheDocument();
  });
});
