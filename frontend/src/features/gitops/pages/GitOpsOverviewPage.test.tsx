import '@testing-library/jest-dom/vitest';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { MemoryRouter } from 'react-router-dom';
import { useAuthStore } from '@/features/auth/store';
import { GitOpsOverviewPage } from '@/features/gitops/pages/GitOpsOverviewPage';
import {
  listGitOpsDeliveryUnits,
  listGitOpsSources,
  listGitOpsTargetGroups
} from '@/services/gitops';
import { installAntdDomShims } from '@/test/installAntdDomShims';

vi.mock('@/services/gitops', async () => {
  const actual = await vi.importActual<typeof import('@/services/gitops')>('@/services/gitops');
  return {
    ...actual,
    listGitOpsSources: vi.fn(),
    listGitOpsTargetGroups: vi.fn(),
    listGitOpsDeliveryUnits: vi.fn()
  };
});

describe('GitOpsOverviewPage', () => {
  beforeAll(() => {
    installAntdDomShims();
  });

  beforeEach(() => {
    vi.clearAllMocks();
    useAuthStore.setState({
      isAuthenticated: true,
      accessToken: 'token',
      refreshToken: 'refresh',
      user: {
        id: 'u1',
        username: 'alice',
        platformRoles: ['platform-admin']
      }
    });
  });

  it('loads source, target group and delivery unit data and opens source drawer', async () => {
    vi.mocked(listGitOpsSources).mockResolvedValue({
      items: [
        {
          id: 'src-1',
          name: 'payments-git',
          sourceType: 'git',
          endpoint: 'https://git.example.com/payments.git',
          status: 'ready'
        }
      ]
    });
    vi.mocked(listGitOpsTargetGroups).mockResolvedValue({
      items: [
        {
          id: 'tg-1',
          name: 'prod-cn-group',
          workspaceId: 1001,
          clusterRefs: [1, 2],
          status: 'active'
        }
      ]
    });
    vi.mocked(listGitOpsDeliveryUnits).mockResolvedValue({
      items: [
        {
          id: 'du-1',
          name: 'payment-api',
          sourceId: 'src-1',
          desiredState: 'main',
          actualState: 'main',
          driftStatus: 'in_sync',
          lastSyncResult: 'succeeded'
        }
      ]
    });

    const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
    render(
      <QueryClientProvider client={client}>
        <MemoryRouter initialEntries={['/gitops?keyword=payment']}>
          <GitOpsOverviewPage />
        </MemoryRouter>
      </QueryClientProvider>
    );

    expect(await screen.findByText('payments-git')).toBeInTheDocument();
    expect(await screen.findByText('prod-cn-group')).toBeInTheDocument();
    expect(await screen.findByText('payment-api')).toBeInTheDocument();

    await waitFor(() => {
      expect(listGitOpsDeliveryUnits).toHaveBeenCalledWith(expect.objectContaining({ keyword: 'payment' }));
    });

    fireEvent.click(screen.getByRole('button', { name: '新建来源' }));
    expect(await screen.findByText('新建交付来源')).toBeInTheDocument();
  });
});
