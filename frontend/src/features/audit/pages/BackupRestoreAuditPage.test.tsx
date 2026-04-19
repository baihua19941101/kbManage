import '@testing-library/jest-dom/vitest';
import { render, screen, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { MemoryRouter } from 'react-router-dom';
import { BackupRestoreAuditPage } from '@/features/audit/pages/BackupRestoreAuditPage';
import { useAuthStore } from '@/features/auth/store';
import { listBackupRestoreAuditEvents } from '@/services/backupRestore';
import { installAntdDomShims } from '@/test/installAntdDomShims';

vi.mock('@/services/backupRestore', async () => ({
  ...(await vi.importActual<typeof import('@/services/backupRestore')>('@/services/backupRestore')),
  listBackupRestoreAuditEvents: vi.fn()
}));

const renderPage = () =>
  render(
    <QueryClientProvider client={new QueryClient({ defaultOptions: { queries: { retry: false } } })}>
      <MemoryRouter>
        <BackupRestoreAuditPage />
      </MemoryRouter>
    </QueryClientProvider>
  );

describe('BackupRestoreAuditPage', () => {
  beforeAll(() => installAntdDomShims());
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(listBackupRestoreAuditEvents).mockResolvedValue({ items: [] });
  });

  it('shows unauthorized empty', () => {
    useAuthStore.setState({
      isAuthenticated: true,
      accessToken: 't',
      refreshToken: 'r',
      user: { id: 'u', username: 'u', platformRoles: [] }
    });
    renderPage();
    expect(screen.getByText('你暂无备份恢复审计访问权限。')).toBeInTheDocument();
  });

  it('queries audit events', async () => {
    useAuthStore.setState({
      isAuthenticated: true,
      accessToken: 't',
      refreshToken: 'r',
      user: { id: 'u', username: 'u', platformRoles: ['audit-reader'] }
    });
    vi.mocked(listBackupRestoreAuditEvents).mockResolvedValue({
      items: [{ id: 'a1', action: 'backuprestore.backup.run', outcome: 'succeeded', occurredAt: '2026-04-18T10:00:00Z' }]
    });
    renderPage();
    await waitFor(() => expect(listBackupRestoreAuditEvents).toHaveBeenCalled());
    expect(await screen.findByText('backuprestore.backup.run')).toBeInTheDocument();
  });
});
