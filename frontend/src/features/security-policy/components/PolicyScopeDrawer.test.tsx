import '@testing-library/jest-dom/vitest';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { PolicyScopeDrawer } from '@/features/security-policy/components/PolicyScopeDrawer';
import { createPolicyAssignment, listPolicyAssignments } from '@/services/securityPolicy';
import { installAntdDomShims } from '@/test/installAntdDomShims';

vi.mock('@/services/securityPolicy', async () => {
  const actual = await vi.importActual<typeof import('@/services/securityPolicy')>(
    '@/services/securityPolicy'
  );
  return {
    ...actual,
    listPolicyAssignments: vi.fn(),
    createPolicyAssignment: vi.fn()
  };
});

const renderWithClient = (ui: JSX.Element) => {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(<QueryClientProvider client={client}>{ui}</QueryClientProvider>);
};

describe('PolicyScopeDrawer', () => {
  beforeAll(() => {
    installAntdDomShims();
  });

  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(listPolicyAssignments).mockResolvedValue({
      items: [
        {
          id: 'as-1',
          policyId: 'policy-1',
          clusterRefs: ['prod-cn-1'],
          namespaceRefs: ['payments'],
          resourceKinds: ['Pod'],
          enforcementMode: 'audit',
          rolloutStage: 'pilot',
          status: 'active'
        }
      ]
    });
    vi.mocked(createPolicyAssignment).mockResolvedValue({
      id: 'task-1',
      policyId: 'policy-1',
      operation: 'assign',
      status: 'pending'
    });
  });

  it('submits assignment payload', async () => {
    renderWithClient(
      <PolicyScopeDrawer
        open
        policy={{
          id: 'policy-1',
          name: 'restrict-privileged',
          scopeLevel: 'platform',
          category: 'pod-security',
          defaultEnforcementMode: 'warn',
          riskLevel: 'high',
          status: 'active'
        }}
        onClose={vi.fn()}
      />
    );

    expect(await screen.findByText('当前分配')).toBeInTheDocument();

    fireEvent.change(screen.getByPlaceholderText('例如：prod-cn-1,prod-cn-2'), {
      target: { value: 'prod-cn-1,prod-cn-2' }
    });
    fireEvent.change(screen.getByPlaceholderText('例如：payments,checkout'), {
      target: { value: 'payments,checkout' }
    });
    fireEvent.change(screen.getByPlaceholderText('例如：Pod,Deployment'), {
      target: { value: 'Pod,Deployment' }
    });

    fireEvent.click(screen.getByRole('button', { name: '提交分配' }));

    await waitFor(() => {
      expect(createPolicyAssignment).toHaveBeenCalledWith('policy-1', {
        workspaceId: undefined,
        projectId: undefined,
        clusterRefs: ['prod-cn-1', 'prod-cn-2'],
        namespaceRefs: ['payments', 'checkout'],
        resourceKinds: ['Pod', 'Deployment'],
        enforcementMode: 'audit',
        rolloutStage: 'pilot'
      });
    });
  });
});
