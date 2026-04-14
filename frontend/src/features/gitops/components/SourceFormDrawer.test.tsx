import '@testing-library/jest-dom/vitest';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { SourceFormDrawer } from '@/features/gitops/components/SourceFormDrawer';
import {
  createGitOpsSource,
  updateGitOpsSource,
  type GitOpsSourceFormData
} from '@/services/gitops';
import { installAntdDomShims } from '@/test/installAntdDomShims';

vi.mock('@/services/gitops', async () => {
  const actual = await vi.importActual<typeof import('@/services/gitops')>('@/services/gitops');
  return {
    ...actual,
    createGitOpsSource: vi.fn(),
    updateGitOpsSource: vi.fn()
  };
});

const renderWithClient = (ui: JSX.Element) => {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(<QueryClientProvider client={client}>{ui}</QueryClientProvider>);
};

describe('SourceFormDrawer', () => {
  beforeAll(() => {
    installAntdDomShims();
  });

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('submits create payload', async () => {
    vi.mocked(createGitOpsSource).mockResolvedValue({
      id: 'src-1',
      name: 'payments-git',
      sourceType: 'git',
      endpoint: 'https://git.example.com/payments.git',
      status: 'ready'
    });

    const onClose = vi.fn();
    const onSuccess = vi.fn();

    renderWithClient(<SourceFormDrawer open onClose={onClose} onSuccess={onSuccess} />);

    fireEvent.change(screen.getByPlaceholderText('例如：payments-git'), {
      target: { value: 'payments-git' }
    });
    fireEvent.change(screen.getByPlaceholderText('https://git.example.com/repo.git'), {
      target: { value: 'https://git.example.com/payments.git' }
    });
    fireEvent.change(screen.getByPlaceholderText('例如：main'), {
      target: { value: 'main' }
    });
    fireEvent.change(screen.getByPlaceholderText('例如：git-credential-prod'), {
      target: { value: 'git-credential-prod' }
    });
    fireEvent.change(screen.getByPlaceholderText('例如：1001'), {
      target: { value: '1001' }
    });
    fireEvent.change(screen.getByPlaceholderText('例如：2001（可选）'), {
      target: { value: '2001' }
    });
    fireEvent.click(screen.getByRole('button', { name: '保存来源' }));

    await waitFor(() => {
      const expectedPayload: GitOpsSourceFormData = {
        name: 'payments-git',
        sourceType: 'git',
        endpoint: 'https://git.example.com/payments.git',
        defaultRef: 'main',
        credentialRef: 'git-credential-prod',
        workspaceId: 1001,
        projectId: 2001
      };
      expect(createGitOpsSource).toHaveBeenCalledWith(expectedPayload);
    });

    await waitFor(() => {
      expect(onSuccess).toHaveBeenCalled();
      expect(onClose).toHaveBeenCalled();
    });
  });

  it('submits update payload when source exists', async () => {
    vi.mocked(updateGitOpsSource).mockResolvedValue({
      id: 'src-1',
      name: 'payments-git-updated',
      sourceType: 'git',
      endpoint: 'https://git.example.com/payments.git',
      status: 'ready'
    });

    renderWithClient(
      <SourceFormDrawer
        open
        onClose={vi.fn()}
        source={{
          id: 'src-1',
          name: 'payments-git',
          sourceType: 'git',
          endpoint: 'https://git.example.com/payments.git',
          defaultRef: 'main',
          credentialRef: 'git-credential-prod',
          workspaceId: 1001,
          projectId: 2001,
          status: 'ready'
        }}
      />
    );

    fireEvent.change(screen.getByPlaceholderText('例如：payments-git'), {
      target: { value: 'payments-git-updated' }
    });
    fireEvent.click(screen.getByRole('button', { name: '保存来源' }));

    await waitFor(() => {
      expect(updateGitOpsSource).toHaveBeenCalledWith('src-1', {
        name: 'payments-git-updated',
        sourceType: 'git',
        endpoint: 'https://git.example.com/payments.git',
        defaultRef: 'main',
        credentialRef: 'git-credential-prod',
        workspaceId: 1001,
        projectId: 2001
      });
    });
  });
});
