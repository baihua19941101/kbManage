import '@testing-library/jest-dom/vitest';
import { render, screen } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { MemoryRouter } from 'react-router-dom';
import { BatchOperationPage } from '@/features/workload-ops/pages/BatchOperationPage';
import { installAntdDomShims } from '@/test/installAntdDomShims';

vi.mock('@/services/workloadOps', () => ({
  getBatchOperation: vi.fn().mockResolvedValue({
    id: 12,
    status: 'succeeded',
    succeededTargets: 2,
    failedTargets: 0,
    items: [{ resourceRef: 'Deployment/default/demo-api', status: 'succeeded', resultMessage: 'ok' }]
  })
}));

describe('BatchOperationPage', () => {
  beforeAll(() => {
    installAntdDomShims();
  });

  it('renders batch operation summary and rows', async () => {
    const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
    render(
      <QueryClientProvider client={client}>
        <MemoryRouter initialEntries={['/workload-ops/batches?batchId=12']}>
          <BatchOperationPage />
        </MemoryRouter>
      </QueryClientProvider>
    );

    expect(screen.getByText('批量任务结果')).toBeInTheDocument();
    expect(await screen.findByText(/任务ID: 12/)).toBeInTheDocument();
    expect(await screen.findByText('Deployment/default/demo-api')).toBeInTheDocument();
  });
});
