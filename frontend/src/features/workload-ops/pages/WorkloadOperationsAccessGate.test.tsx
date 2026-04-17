import '@testing-library/jest-dom/vitest';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { MemoryRouter } from 'react-router-dom';
import { useAuthStore } from '@/features/auth/store';
import { WorkloadOperationsPage } from '@/features/workload-ops/pages/WorkloadOperationsPage';
import { ApiError } from '@/services/api/client';
import { installAntdDomShims } from '@/test/installAntdDomShims';
import {
  getWorkloadOpsContext,
  listWorkloadOpsInstances,
  listWorkloadOpsRevisions,
  submitWorkloadAction
} from '@/services/workloadOps';

vi.mock('@/services/workloadOps', () => ({
  getWorkloadOpsContext: vi.fn().mockResolvedValue({
    clusterId: 1,
    namespace: 'default',
    resourceKind: 'Deployment',
    resourceName: 'demo-api',
    healthStatus: 'healthy',
    rolloutStatus: 'running'
  }),
  listWorkloadOpsInstances: vi.fn().mockResolvedValue({
    items: []
  }),
  listWorkloadOpsRevisions: vi.fn().mockResolvedValue({
    items: []
  }),
  submitWorkloadAction: vi.fn().mockResolvedValue({
    id: 1,
    actionType: 'restart',
    riskLevel: 'medium',
    status: 'pending'
  }),
  createTerminalSession: vi.fn().mockResolvedValue({ id: 1, status: 'active' }),
  closeTerminalSession: vi.fn().mockResolvedValue({})
}));

describe('WorkloadOperationsAccessGate', () => {
  beforeAll(() => {
    installAntdDomShims();
  });

  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(getWorkloadOpsContext).mockResolvedValue({
      clusterId: 1,
      namespace: 'default',
      resourceKind: 'Deployment',
      resourceName: 'demo-api',
      healthStatus: 'healthy',
      rolloutStatus: 'running'
    });
    vi.mocked(listWorkloadOpsInstances).mockResolvedValue({ items: [] });
    vi.mocked(listWorkloadOpsRevisions).mockResolvedValue({ items: [] });
    vi.mocked(submitWorkloadAction).mockResolvedValue({
      id: 1,
      actionType: 'restart',
      riskLevel: 'medium',
      status: 'pending'
    });
  });

  it('shows unauthorized empty state for users without workloadops read roles', async () => {
    useAuthStore.setState({
      isAuthenticated: true,
      accessToken: 'token',
      refreshToken: 'refresh',
      user: { id: 'u1', username: 'u1', platformRoles: [] }
    });

    const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
    render(
      <QueryClientProvider client={client}>
        <MemoryRouter>
          <WorkloadOperationsPage />
        </MemoryRouter>
      </QueryClientProvider>
    );

    expect(screen.getByText('你暂无工作负载运维访问权限，请联系管理员授予工作空间/项目范围。')).toBeInTheDocument();
  });

  it('shows readonly alert and disables actions for readonly users', async () => {
    useAuthStore.setState({
      isAuthenticated: true,
      accessToken: 'token',
      refreshToken: 'refresh',
      user: { id: 'readonly-u', username: 'readonly-u', platformRoles: ['readonly'] }
    });

    const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
    render(
      <QueryClientProvider client={client}>
        <MemoryRouter>
          <WorkloadOperationsPage />
        </MemoryRouter>
      </QueryClientProvider>
    );

    expect(await screen.findByText('当前为只读模式')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: '提交重启动作' })).toBeDisabled();
    expect(screen.getByRole('button', { name: '提交批量重启' })).toBeDisabled();
  });

  it('shows permission changed warning when query returns authorization error', async () => {
    vi.mocked(getWorkloadOpsContext).mockRejectedValue(new ApiError(403, 'forbidden', { url: '/api/v1/workload-ops/resources/context' }));
    vi.mocked(listWorkloadOpsInstances).mockResolvedValue({ items: [] });
    vi.mocked(listWorkloadOpsRevisions).mockResolvedValue({ items: [] });

    useAuthStore.setState({
      isAuthenticated: true,
      accessToken: 'token',
      refreshToken: 'refresh',
      user: { id: 'u2', username: 'u2', platformRoles: ['ops-operator'] }
    });

    const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
    render(
      <QueryClientProvider client={client}>
        <MemoryRouter initialEntries={['/workload-ops?clusterId=1&namespace=default&resourceKind=Deployment&resourceName=demo']}>
          <WorkloadOperationsPage />
        </MemoryRouter>
      </QueryClientProvider>
    );

    expect(await screen.findByText('权限已变更，当前动作入口已锁定。')).toBeInTheDocument();
  });

  it('locks rollback action after permission revoked error', async () => {
    vi.mocked(listWorkloadOpsRevisions).mockResolvedValue({
      items: [
        {
          revision: 3,
          sourceKind: 'replicaset',
          sourceName: 'demo-api-rs-3',
          isCurrent: false,
          rollbackAvailable: true
        }
      ]
    });
    vi.mocked(submitWorkloadAction).mockRejectedValue(
      new ApiError(403, 'forbidden', { url: '/api/v1/workload-ops/actions' })
    );

    useAuthStore.setState({
      isAuthenticated: true,
      accessToken: 'token',
      refreshToken: 'refresh',
      user: { id: 'ops-u', username: 'ops-u', platformRoles: ['ops-operator'] }
    });

    const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
    render(
      <QueryClientProvider client={client}>
        <MemoryRouter>
          <WorkloadOperationsPage />
        </MemoryRouter>
      </QueryClientProvider>
    );

    fireEvent.click(await screen.findByRole('button', { name: '回滚' }));
    fireEvent.click(await screen.findByRole('button', { name: '确认回滚' }));

    expect(await screen.findByText('权限已回收')).toBeInTheDocument();
    await waitFor(() => {
      expect(screen.getByRole('button', { name: '确认回滚' })).toBeDisabled();
    });
  });
});
