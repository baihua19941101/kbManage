import '@testing-library/jest-dom/vitest';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import { useAuthStore } from '@/features/auth/store';
import { DeliveryUnitDetailPage } from '@/features/gitops/pages/DeliveryUnitDetailPage';
import {
  getGitOpsDeliveryUnit,
  getGitOpsDeliveryUnitDiff,
  getGitOpsDeliveryUnitStatus,
  listGitOpsReleaseRevisions,
  listGitOpsTargetGroups
} from '@/services/gitops';
import { installAntdDomShims } from '@/test/installAntdDomShims';

vi.mock('@/services/gitops', async () => {
  const actual = await vi.importActual<typeof import('@/services/gitops')>('@/services/gitops');
  return {
    ...actual,
    getGitOpsDeliveryUnit: vi.fn(),
    getGitOpsDeliveryUnitStatus: vi.fn(),
    getGitOpsDeliveryUnitDiff: vi.fn(),
    listGitOpsReleaseRevisions: vi.fn(),
    listGitOpsTargetGroups: vi.fn()
  };
});

describe('DeliveryUnitDetailPage', () => {
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

  it('renders detail with diff, revisions and stage editing', async () => {
    vi.mocked(getGitOpsDeliveryUnit).mockResolvedValue({
      id: 'du-1',
      name: 'payment-api',
      sourceId: 'src-1',
      sourcePath: 'apps/payment',
      syncMode: 'manual',
      desiredRevision: 'main',
      desiredAppVersion: '1.2.3',
      desiredConfigVersion: '2026.04.13',
      environments: [
        {
          name: 'dev',
          orderIndex: 1,
          targetGroupId: 1,
          promotionMode: 'manual',
          paused: false
        },
        {
          name: 'prod',
          orderIndex: 2,
          targetGroupId: 2,
          promotionMode: 'manual',
          paused: false
        }
      ],
      overlays: [
        {
          overlayType: 'values',
          overlayRef: 'overlays/dev-values.yaml',
          effectiveScope: 'dev',
          precedence: 10
        }
      ]
    });

    vi.mocked(getGitOpsDeliveryUnitStatus).mockResolvedValue({
      deliveryStatus: 'ready',
      driftStatus: 'in_sync',
      lastSyncResult: 'succeeded',
      environments: [
        {
          environment: 'dev',
          syncStatus: 'succeeded',
          driftStatus: 'in_sync',
          targetCount: 2,
          succeededCount: 2,
          failedCount: 0
        }
      ]
    });

    vi.mocked(getGitOpsDeliveryUnitDiff).mockResolvedValue({
      summary: {
        added: 0,
        modified: 1,
        removed: 0,
        unavailable: 0
      },
      items: [
        {
          objectRef: 'Deployment/default/payment-api',
          environment: 'dev',
          diffType: 'modified',
          desiredSummary: 'replicas=3',
          liveSummary: 'replicas=2'
        }
      ]
    });

    vi.mocked(listGitOpsReleaseRevisions).mockResolvedValue({
      items: [
        {
          id: 1,
          sourceRevision: 'main@a1b2',
          appVersion: '1.2.3',
          configVersion: '2026.04.13',
          status: 'active',
          rollbackAvailable: true,
          createdAt: '2026-04-13T10:00:00Z'
        }
      ]
    });

    vi.mocked(listGitOpsTargetGroups).mockResolvedValue({
      items: [
        { id: 1, name: 'dev-group', workspaceId: 1001 },
        { id: 2, name: 'prod-group', workspaceId: 1001 }
      ]
    });

    const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
    render(
      <QueryClientProvider client={client}>
        <MemoryRouter initialEntries={['/gitops/delivery-units/du-1']}>
          <Routes>
            <Route path="/gitops/delivery-units/:unitId" element={<DeliveryUnitDetailPage />} />
          </Routes>
        </MemoryRouter>
      </QueryClientProvider>
    );

    expect(await screen.findByText('payment-api')).toBeInTheDocument();
    expect(await screen.findByText('overlays/dev-values.yaml')).toBeInTheDocument();
    expect(await screen.findByText('Deployment/default/payment-api')).toBeInTheDocument();
    expect(await screen.findByText('main@a1b2')).toBeInTheDocument();

    const stageInputs = await screen.findAllByPlaceholderText('阶段名称');
    fireEvent.change(stageInputs[0], { target: { value: 'dev-canary' } });

    await waitFor(() => {
      expect(screen.getByDisplayValue('dev-canary')).toBeInTheDocument();
    });
  });
});
