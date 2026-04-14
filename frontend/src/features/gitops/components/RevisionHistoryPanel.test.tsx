import '@testing-library/jest-dom/vitest';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { RevisionHistoryPanel } from '@/features/gitops/components/RevisionHistoryPanel';
import { listGitOpsReleaseRevisions } from '@/services/gitops';
import { installAntdDomShims } from '@/test/installAntdDomShims';

vi.mock('@/services/gitops', async () => {
  const actual = await vi.importActual<typeof import('@/services/gitops')>('@/services/gitops');
  return {
    ...actual,
    listGitOpsReleaseRevisions: vi.fn()
  };
});

const renderWithClient = (ui: JSX.Element) => {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(<QueryClientProvider client={client}>{ui}</QueryClientProvider>);
};

describe('RevisionHistoryPanel', () => {
  beforeAll(() => {
    installAntdDomShims();
  });

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('loads revision list and emits rollback selection', async () => {
    vi.mocked(listGitOpsReleaseRevisions).mockResolvedValue({
      items: [
        {
          id: 22,
          sourceRevision: 'main@a1b2',
          appVersion: '1.2.3',
          configVersion: '2026.04.13',
          status: 'historical',
          rollbackAvailable: true,
          createdAt: '2026-04-13T10:00:00Z'
        }
      ]
    });

    const onRollback = vi.fn();

    renderWithClient(<RevisionHistoryPanel unitId="du-1" onRollback={onRollback} />);

    expect(await screen.findByText('main@a1b2')).toBeInTheDocument();
    expect(screen.getByText('1.2.3')).toBeInTheDocument();

    fireEvent.click(screen.getByRole('button', { name: '回滚' }));

    await waitFor(() => {
      expect(onRollback).toHaveBeenCalledWith(expect.objectContaining({ id: 22 }));
    });
  });
});
