import '@testing-library/jest-dom/vitest';
import { fireEvent, render, screen } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { useAuthStore } from '@/features/auth/store';
import { RollbackDialog } from '@/features/workload-ops/components/RollbackDialog';
import { installAntdDomShims } from '@/test/installAntdDomShims';

vi.mock('@/services/workloadOps', () => ({
  submitWorkloadAction: vi.fn().mockResolvedValue({
    id: 100,
    status: 'succeeded',
    resultMessage: 'rolled back workload default/demo-api to revision 2'
  })
}));

describe('RollbackDialog', () => {
  beforeAll(() => {
    installAntdDomShims();
  });

  it('submits rollback and shows success message', async () => {
    useAuthStore.setState({
      isAuthenticated: true,
      accessToken: 'token',
      refreshToken: 'refresh',
      user: { id: 'ops-user', username: 'ops-user', platformRoles: ['ops-operator'] }
    });

    const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
    render(
      <QueryClientProvider client={client}>
        <RollbackDialog
          open
          clusterId={1}
          namespace="default"
          resourceKind="Deployment"
          resourceName="demo-api"
          revision={{
            revision: 2,
            sourceKind: 'replicaset',
            sourceName: 'demo-api-rs-2',
            isCurrent: false,
            rollbackAvailable: true
          }}
          onClose={() => {}}
        />
      </QueryClientProvider>
    );

    fireEvent.click(screen.getByRole('button', { name: '确认回滚' }));
    expect(await screen.findByText(/revision 2/)).toBeInTheDocument();
  });
});
