import '@testing-library/jest-dom/vitest';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { MemoryRouter } from 'react-router-dom';
import { useAuthStore } from '@/features/auth/store';
import { PolicyRolloutPage } from '@/features/security-policy/pages/PolicyRolloutPage';
import {
  listPolicyExceptions,
  listPolicyHits,
  listSecurityPolicies,
  reviewExceptionRequest,
  switchPolicyMode
} from '@/services/securityPolicy';
import { installAntdDomShims } from '@/test/installAntdDomShims';

vi.mock('@/services/securityPolicy', async () => {
  const actual = await vi.importActual<typeof import('@/services/securityPolicy')>(
    '@/services/securityPolicy'
  );
  return {
    ...actual,
    listSecurityPolicies: vi.fn(),
    listPolicyHits: vi.fn(),
    listPolicyExceptions: vi.fn(),
    switchPolicyMode: vi.fn(),
    createExceptionRequest: vi.fn(),
    reviewExceptionRequest: vi.fn()
  };
});

const renderWithClient = () => {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(
    <QueryClientProvider client={client}>
      <MemoryRouter>
        <PolicyRolloutPage />
      </MemoryRouter>
    </QueryClientProvider>
  );
};

describe('PolicyRolloutPage', () => {
  beforeAll(() => {
    installAntdDomShims();
  });

  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(listSecurityPolicies).mockResolvedValue({ items: [] });
    vi.mocked(listPolicyHits).mockResolvedValue({ items: [] });
    vi.mocked(listPolicyExceptions).mockResolvedValue({ items: [] });
    vi.mocked(switchPolicyMode).mockResolvedValue({
      id: 'task-1',
      policyId: 'policy-1',
      operation: 'mode-switch',
      status: 'pending'
    });
    vi.mocked(reviewExceptionRequest).mockResolvedValue({
      id: 'ex-1',
      policyId: 'policy-1',
      hitRecordId: 'hit-1',
      status: 'approved',
      startsAt: '2026-04-15T08:00:00Z',
      expiresAt: '2026-04-16T08:00:00Z'
    });
  });

  it('shows unauthorized empty when user cannot read rollout data', () => {
    useAuthStore.setState({
      isAuthenticated: true,
      accessToken: 'token',
      refreshToken: 'refresh',
      user: {
        id: 'u0',
        username: 'u0',
        platformRoles: []
      }
    });

    renderWithClient();

    expect(
      screen.getByText('你暂无策略灰度与例外治理访问权限，请联系管理员授予 policy:read 或对应平台角色。')
    ).toBeInTheDocument();
  });

  it('renders rollout data and supports mode switch submission', async () => {
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

    vi.mocked(listSecurityPolicies).mockResolvedValue({
      items: [
        {
          id: 'policy-1',
          name: 'restrict-root',
          scopeLevel: 'platform',
          category: 'pod-security',
          defaultEnforcementMode: 'warn',
          riskLevel: 'high',
          status: 'active'
        }
      ]
    });

    vi.mocked(listPolicyHits).mockResolvedValue({
      items: [
        {
          id: 'hit-1',
          policyId: 'policy-1',
          clusterId: 'prod-cn-1',
          namespace: 'payments',
          resourceKind: 'Pod',
          resourceName: 'pay-api-6dd7',
          hitResult: 'warn',
          riskLevel: 'high',
          remediationStatus: 'open',
          detectedAt: '2026-04-14T08:00:00Z'
        }
      ]
    });

    vi.mocked(listPolicyExceptions).mockResolvedValue({
      items: [
        {
          id: 'ex-1',
          policyId: 'policy-1',
          hitRecordId: 'hit-1',
          reason: 'temporary override',
          status: 'pending',
          startsAt: '2026-04-15T08:00:00Z',
          expiresAt: '2026-04-16T08:00:00Z'
        }
      ]
    });

    renderWithClient();

    await waitFor(() => {
      expect(listPolicyHits).toHaveBeenCalledWith({ policyId: 'policy-1' });
      expect(listPolicyExceptions).toHaveBeenCalledWith({ policyId: 'policy-1' });
    });
    expect(await screen.findByText('策略命中（1）')).toBeInTheDocument();
    expect(screen.getByText('待审批')).toBeInTheDocument();

    fireEvent.click(screen.getByRole('button', { name: '模式切换' }));
    expect(await screen.findByText('模式切换 - restrict-root')).toBeInTheDocument();

    fireEvent.click(screen.getByRole('button', { name: '提交切换' }));
    await waitFor(() => {
      expect(switchPolicyMode).toHaveBeenCalledWith('policy-1', {
        targetMode: 'warn',
        assignmentIds: undefined,
        reason: undefined
      });
    });
  });
});
