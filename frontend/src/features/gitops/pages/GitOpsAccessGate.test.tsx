import '@testing-library/jest-dom/vitest';
import { render, screen } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { MemoryRouter } from 'react-router-dom';
import { useAuthStore } from '@/features/auth/store';
import { GitOpsOverviewPage } from '@/features/gitops/pages/GitOpsOverviewPage';
import { ApiError } from '@/services/api/client';
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

describe('GitOpsAccessGate', () => {
  beforeAll(() => {
    installAntdDomShims();
  });

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('shows unauthorized empty state for users without gitops read roles', () => {
    useAuthStore.setState({
      isAuthenticated: true,
      accessToken: 'token',
      refreshToken: 'refresh-token',
      user: {
        id: 'gitops-no-read',
        username: 'gitops-no-read',
        platformRoles: []
      }
    });

    const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });

    render(
      <QueryClientProvider client={client}>
        <MemoryRouter>
          <GitOpsOverviewPage />
        </MemoryRouter>
      </QueryClientProvider>
    );

    expect(screen.getByText('你暂无 GitOps 访问权限，请联系管理员授予范围权限。')).toBeInTheDocument();
    expect(listGitOpsSources).not.toHaveBeenCalled();
    expect(listGitOpsTargetGroups).not.toHaveBeenCalled();
    expect(listGitOpsDeliveryUnits).not.toHaveBeenCalled();
  });

  it('shows permission changed warning when query returns authorization error', async () => {
    useAuthStore.setState({
      isAuthenticated: true,
      accessToken: 'token',
      refreshToken: 'refresh-token',
      user: {
        id: 'gitops-read',
        username: 'gitops-read',
        platformRoles: ['ops-operator']
      }
    });

    vi.mocked(listGitOpsSources).mockRejectedValue(
      new ApiError(403, 'forbidden', {
        url: '/api/v1/gitops/sources'
      })
    );
    vi.mocked(listGitOpsTargetGroups).mockResolvedValue({ items: [] });
    vi.mocked(listGitOpsDeliveryUnits).mockResolvedValue({ items: [] });

    const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });

    render(
      <QueryClientProvider client={client}>
        <MemoryRouter>
          <GitOpsOverviewPage />
        </MemoryRouter>
      </QueryClientProvider>
    );

    expect(await screen.findByText('权限已变更')).toBeInTheDocument();
  });
});
