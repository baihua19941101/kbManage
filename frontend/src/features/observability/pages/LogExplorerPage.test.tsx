import '@testing-library/jest-dom/vitest';
import { render, screen } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { MemoryRouter } from 'react-router-dom';
import { LogExplorerPage } from '@/features/observability/pages/LogExplorerPage';
import { installAntdDomShims } from '@/test/installAntdDomShims';

Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: vi.fn().mockImplementation(() => ({
    matches: false,
    media: '',
    onchange: null,
    addListener: vi.fn(),
    removeListener: vi.fn(),
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
    dispatchEvent: vi.fn()
  }))
});

vi.mock('@/services/observability/logs', () => ({
  queryObservabilityLogs: vi.fn().mockResolvedValue({
    items: [
      {
        timestamp: '2026-01-01T00:00:00Z',
        clusterId: '1',
        namespace: 'default',
        workload: 'mock-app',
        pod: 'mock-app-1',
        container: 'main',
        message: 'probe succeeded'
      }
    ]
  })
}));

describe('LogExplorerPage', () => {
  beforeAll(() => {
    installAntdDomShims();
  });

  it('renders log explorer and row data', async () => {
    const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
    render(
      <QueryClientProvider client={client}>
        <MemoryRouter>
          <LogExplorerPage />
        </MemoryRouter>
      </QueryClientProvider>
    );

    expect(screen.getByText('日志检索')).toBeInTheDocument();
    expect(await screen.findByText(/probe succeeded/i)).toBeInTheDocument();
  });
});
