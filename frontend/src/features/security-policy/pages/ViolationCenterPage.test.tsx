import '@testing-library/jest-dom/vitest';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { MemoryRouter } from 'react-router-dom';
import { useAuthStore } from '@/features/auth/store';
import { ViolationCenterPage } from '@/features/security-policy/pages/ViolationCenterPage';
import {
  listPolicyHits,
  listSecurityPolicies,
  updatePolicyHitRemediation
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
    updatePolicyHitRemediation: vi.fn()
  };
});

const renderWithClient = () => {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(
    <QueryClientProvider client={client}>
      <MemoryRouter>
        <ViolationCenterPage />
      </MemoryRouter>
    </QueryClientProvider>
  );
};

describe('ViolationCenterPage', () => {
  beforeAll(() => {
    installAntdDomShims();
  });

  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(listSecurityPolicies).mockResolvedValue({ items: [] });
    vi.mocked(listPolicyHits).mockResolvedValue({ items: [] });
    vi.mocked(updatePolicyHitRemediation).mockResolvedValue({
      id: 'hit-1',
      policyId: 'policy-1',
      hitResult: 'warn',
      riskLevel: 'high',
      remediationStatus: 'in_progress',
      detectedAt: '2026-04-14T08:00:00Z'
    });
  });

  it('shows unauthorized empty when user cannot read violations', () => {
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
      screen.getByText('你暂无违规中心访问权限，请联系管理员授予 policy:read 或对应平台角色。')
    ).toBeInTheDocument();
  });

  it('renders violations and submits remediation update', async () => {
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

    renderWithClient();

    expect(await screen.findByText('违规列表（1）')).toBeInTheDocument();
    expect(screen.getByText('风险级别分布')).toBeInTheDocument();

    fireEvent.click(screen.getByRole('button', { name: '更新整改' }));
    expect(await screen.findByText('更新整改状态 - hit-1')).toBeInTheDocument();

    fireEvent.click(screen.getByLabelText('in_progress'));
    fireEvent.change(screen.getByPlaceholderText('例如：已完成镜像修复并重新部署'), {
      target: { value: '修复镜像并重新发布' }
    });
    fireEvent.click(screen.getByRole('button', { name: '提交更新' }));

    await waitFor(() => {
      expect(updatePolicyHitRemediation).toHaveBeenCalledWith('hit-1', {
        remediationStatus: 'in_progress',
        comment: '修复镜像并重新发布'
      });
    });
  });
});
