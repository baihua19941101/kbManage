import '@testing-library/jest-dom/vitest';
import { fireEvent, render, screen, within } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { MemoryRouter } from 'react-router-dom';
import { useAuthStore } from '@/features/auth/store';
import { PolicyCenterPage } from '@/features/security-policy/pages/PolicyCenterPage';
import { listPolicyAssignments, listSecurityPolicies } from '@/services/securityPolicy';
import { installAntdDomShims } from '@/test/installAntdDomShims';

vi.mock('@/services/securityPolicy', async () => {
  const actual = await vi.importActual<typeof import('@/services/securityPolicy')>(
    '@/services/securityPolicy'
  );
  return {
    ...actual,
    listSecurityPolicies: vi.fn(),
    listPolicyAssignments: vi.fn(),
    createSecurityPolicy: vi.fn(),
    updateSecurityPolicy: vi.fn(),
    createPolicyAssignment: vi.fn()
  };
});

const renderWithClient = () => {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(
    <QueryClientProvider client={client}>
      <MemoryRouter>
        <PolicyCenterPage />
      </MemoryRouter>
    </QueryClientProvider>
  );
};

describe('PolicyCenterPage', () => {
  beforeAll(() => {
    installAntdDomShims();
  });

  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(listSecurityPolicies).mockResolvedValue({ items: [] });
    vi.mocked(listPolicyAssignments).mockResolvedValue({ items: [] });
  });

  it('shows unauthorized empty when user cannot read policy domain', () => {
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
      screen.getByText('你暂无安全策略访问权限，请联系管理员授予 policy:read 或对应平台角色。')
    ).toBeInTheDocument();
  });

  it('renders policy list and opens scope drawer', async () => {
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
          name: 'restrict-privileged',
          scopeLevel: 'platform',
          category: 'pod-security',
          defaultEnforcementMode: 'warn',
          riskLevel: 'high',
          status: 'active',
          updatedAt: '2026-04-14T10:00:00Z'
        }
      ]
    });

    renderWithClient();

    expect(await screen.findByText('restrict-privileged')).toBeInTheDocument();
    const row = screen.getByText('restrict-privileged').closest('tr');
    expect(row).not.toBeNull();
    fireEvent.click(within(row as HTMLElement).getByRole('button', { name: '分配策略' }));
    expect(await screen.findByText('分配策略 - restrict-privileged')).toBeInTheDocument();
  });
});
