import '@testing-library/jest-dom/vitest';
import { render, screen } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { MemoryRouter } from 'react-router-dom';
import { useAuthStore } from '@/features/auth/store';
import { WorkloadOperationsPage } from '@/features/workload-ops/pages/WorkloadOperationsPage';
import { installAntdDomShims } from '@/test/installAntdDomShims';

vi.mock('@/services/workloadOps', () => ({
  getWorkloadOpsContext: vi.fn().mockResolvedValue({
    clusterId: 1,
    namespace: 'default',
    resourceKind: 'Deployment',
    resourceName: 'demo-api',
    healthStatus: 'unknown',
    rolloutStatus: 'unknown'
  }),
  listWorkloadOpsInstances: vi.fn().mockResolvedValue({
    items: [
      {
        podName: 'demo-api-pod-0',
        containerName: 'app',
        phase: 'Running',
        ready: true,
        terminalAvailable: true
      }
    ]
  }),
  createTerminalSession: vi.fn().mockResolvedValue({
    id: 1,
    status: 'active',
    podName: 'demo-api-pod-0',
    containerName: 'app'
  }),
  closeTerminalSession: vi.fn().mockResolvedValue({})
}));

describe('WorkloadOperationsPage', () => {
  beforeAll(() => {
    installAntdDomShims();
  });

  it('renders context and instance rows', async () => {
    useAuthStore.setState({
      isAuthenticated: true,
      accessToken: 'token',
      refreshToken: 'refresh',
      user: { id: 'u1', username: 'u1', platformRoles: ['ops-operator'] }
    });

    const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
    render(
      <QueryClientProvider client={client}>
        <MemoryRouter
          initialEntries={[
            '/workload-ops?clusterId=1&namespace=default&resourceKind=Deployment&resourceName=demo-api'
          ]}
        >
          <WorkloadOperationsPage />
        </MemoryRouter>
      </QueryClientProvider>
    );

    expect(screen.getByText('工作负载运维')).toBeInTheDocument();
    expect(await screen.findByText('demo-api')).toBeInTheDocument();
    expect(await screen.findByText('demo-api-pod-0')).toBeInTheDocument();
  });
});
