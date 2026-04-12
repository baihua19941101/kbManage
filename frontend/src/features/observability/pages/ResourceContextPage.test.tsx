import '@testing-library/jest-dom/vitest';
import { render, screen } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { MemoryRouter } from 'react-router-dom';
import { ResourceContextPage } from '@/features/observability/pages/ResourceContextPage';
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

vi.mock('@/services/observability/overview', () => ({
  getObservabilityOverview: vi.fn().mockResolvedValue({
    cards: [{ title: 'Healthy Clusters', value: 1 }]
  })
}));

vi.mock('@/services/observability/events', () => ({
  listObservabilityEvents: vi.fn().mockResolvedValue({
    items: [
      {
        reason: 'BackOff',
        message: 'mock event',
        eventType: 'warning',
        namespace: 'default',
        involvedKind: 'Deployment',
        involvedName: 'mock-app'
      }
    ]
  })
}));

vi.mock('@/services/observability/logs', () => ({
  queryObservabilityLogs: vi.fn().mockResolvedValue({
    items: []
  })
}));

vi.mock('@/services/observability/metrics', () => ({
  queryMetricSeries: vi.fn().mockResolvedValue({
    metricKey: 'cpu_usage',
    subjectType: 'workload',
    subjectRef: 'mock-app',
    points: []
  })
}));

vi.mock('@/features/observability/components/MetricsTrendChart', () => ({
  MetricsTrendChart: () => <div data-testid="metrics-trend-chart" />
}));

describe('ResourceContextPage', () => {
  beforeAll(() => {
    installAntdDomShims();
  });

  it('renders resource context and timeline', async () => {
    const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
    render(
      <QueryClientProvider client={client}>
        <MemoryRouter
          initialEntries={[
            '/observability/context?clusterId=1&namespace=default&resourceKind=Deployment&resourceName=mock-app'
          ]}
        >
          <ResourceContextPage />
        </MemoryRouter>
      </QueryClientProvider>
    );

    expect(screen.getByText('资源可观测上下文')).toBeInTheDocument();
    expect(await screen.findByText('BackOff')).toBeInTheDocument();
    expect(screen.getByTestId('metrics-trend-chart')).toBeInTheDocument();
  });
});
