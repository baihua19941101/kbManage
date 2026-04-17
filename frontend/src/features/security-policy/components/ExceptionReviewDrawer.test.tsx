import '@testing-library/jest-dom/vitest';
import type { ReactElement } from 'react';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { ExceptionReviewDrawer } from '@/features/security-policy/components/ExceptionReviewDrawer';
import { reviewExceptionRequest } from '@/services/securityPolicy';
import { installAntdDomShims } from '@/test/installAntdDomShims';

vi.mock('@/services/securityPolicy', async () => {
  const actual = await vi.importActual<typeof import('@/services/securityPolicy')>(
    '@/services/securityPolicy'
  );
  return {
    ...actual,
    reviewExceptionRequest: vi.fn()
  };
});

const renderWithClient = (ui: ReactElement) => {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(<QueryClientProvider client={client}>{ui}</QueryClientProvider>);
};

describe('ExceptionReviewDrawer', () => {
  beforeAll(() => {
    installAntdDomShims();
  });

  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(reviewExceptionRequest).mockResolvedValue({
      id: 'ex-1',
      policyId: 'policy-1',
      hitRecordId: 'hit-1',
      reason: 'temporary override',
      status: 'approved',
      startsAt: '2026-04-15T08:00:00Z',
      expiresAt: '2026-04-16T08:00:00Z',
      reviewComment: 'ok'
    });
  });

  it('submits review payload', async () => {
    renderWithClient(
      <ExceptionReviewDrawer
        open
        exception={{
          id: 'ex-1',
          policyId: 'policy-1',
          hitRecordId: 'hit-1',
          reason: 'temporary override',
          status: 'pending',
          startsAt: '2026-04-15T08:00:00Z',
          expiresAt: '2026-04-16T08:00:00Z'
        }}
        onClose={vi.fn()}
      />
    );

    fireEvent.change(screen.getByPlaceholderText('例如：同意临时例外，要求 24h 内完成整改'), {
      target: { value: 'risk accepted temporarily' }
    });

    fireEvent.click(screen.getByRole('button', { name: '提交审批' }));

    await waitFor(() => {
      expect(reviewExceptionRequest).toHaveBeenCalledWith('ex-1', {
        decision: 'approve',
        comment: 'risk accepted temporarily'
      });
    });
  });
});
