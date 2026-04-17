import '@testing-library/jest-dom/vitest';
import type { ReactElement } from 'react';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { useAuthStore } from '@/features/auth/store';
import { ReleaseActionDrawer } from '@/features/gitops/components/ReleaseActionDrawer';
import { getGitOpsOperation, submitGitOpsAction } from '@/services/gitops';
import { installAntdDomShims } from '@/test/installAntdDomShims';

vi.mock('@/services/gitops', async () => {
  const actual = await vi.importActual<typeof import('@/services/gitops')>('@/services/gitops');
  return {
    ...actual,
    submitGitOpsAction: vi.fn(),
    getGitOpsOperation: vi.fn()
  };
});

const renderWithClient = (ui: ReactElement) => {
  const client = new QueryClient({
    defaultOptions: {
      queries: { retry: false },
      mutations: { retry: false }
    }
  });

  return render(<QueryClientProvider client={client}>{ui}</QueryClientProvider>);
};

describe('ReleaseActionDrawer', () => {
  beforeAll(() => {
    installAntdDomShims();
  });

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('submits action and renders operation status', async () => {
    useAuthStore.setState({
      isAuthenticated: true,
      accessToken: 'token',
      refreshToken: 'refresh',
      user: { id: 'ops-user', username: 'ops-user', platformRoles: ['ops-operator'] }
    });

    vi.mocked(submitGitOpsAction).mockResolvedValue({
      id: 101,
      operationType: 'sync',
      actionType: 'sync',
      status: 'succeeded',
      resultMessage: '动作执行成功'
    });
    vi.mocked(getGitOpsOperation).mockResolvedValue({
      id: 101,
      operationType: 'sync',
      actionType: 'sync',
      status: 'succeeded',
      resultMessage: '动作执行成功'
    });

    const onOperationChange = vi.fn();

    renderWithClient(
      <ReleaseActionDrawer
        open
        unitId="du-1"
        unitName="payment-api"
        onClose={vi.fn()}
        onOperationChange={onOperationChange}
      />
    );

    fireEvent.click(screen.getByRole('button', { name: '提交动作' }));

    await waitFor(() => {
      expect(submitGitOpsAction).toHaveBeenCalledWith(
        'du-1',
        expect.objectContaining({ actionType: 'sync' })
      );
    });

    expect(await screen.findByText('动作执行成功')).toBeInTheDocument();

    await waitFor(() => {
      expect(onOperationChange).toHaveBeenCalled();
    });
  });

  it('shows unauthorized state when current role cannot execute action', async () => {
    useAuthStore.setState({
      isAuthenticated: true,
      accessToken: 'token',
      refreshToken: 'refresh',
      user: { id: 'readonly-user', username: 'readonly-user', platformRoles: ['readonly'] }
    });

    renderWithClient(
      <ReleaseActionDrawer open unitId="du-1" unitName="payment-api" onClose={vi.fn()} />
    );

    expect(await screen.findByText('当前动作未授权')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: '提交动作' })).toBeDisabled();
  });
});
