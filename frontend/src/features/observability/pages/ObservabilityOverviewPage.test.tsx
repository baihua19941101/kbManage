import '@testing-library/jest-dom/vitest';
import { render, screen } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { MemoryRouter } from 'react-router-dom';
import { useAuthStore } from '@/features/auth/store';
import { ObservabilityOverviewPage } from '@/features/observability/pages/ObservabilityOverviewPage';
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
    cards: [{ title: 'Healthy Clusters', value: 1, unit: 'clusters' }]
  })
}));

vi.mock('@/services/observability/metrics', () => ({
  queryMetricSeries: vi.fn().mockResolvedValue({
    metricKey: 'cpu_usage',
    subjectType: 'workload',
    subjectRef: 'mock-app',
    points: [
      { timestamp: '2026-01-01T00:00:00Z', value: 0.3 },
      { timestamp: '2026-01-01T00:05:00Z', value: 0.4 }
    ]
  })
}));

vi.mock('@/features/observability/components/MetricsTrendChart', () => ({
  MetricsTrendChart: () => <div data-testid="metrics-trend-chart" />
}));

describe('ObservabilityOverviewPage', () => {
  beforeAll(() => {
    installAntdDomShims();
  });

  beforeEach(() => {
    useAuthStore.setState({
      isAuthenticated: true,
      accessToken: 'token',
      refreshToken: 'refresh-token',
      user: {
        id: 'user-1',
        username: 'tester',
        platformRoles: ['platform-admin']
      }
    });
  });

  it('renders observability overview page', async () => {
    const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
    render(
      <QueryClientProvider client={client}>
        <MemoryRouter>
          <ObservabilityOverviewPage />
        </MemoryRouter>
      </QueryClientProvider>
    );

    expect(screen.getByText('可观测中心')).toBeInTheDocument();
    expect(await screen.findByText('Healthy Clusters')).toBeInTheDocument();
    expect(screen.getByTestId('metrics-trend-chart')).toBeInTheDocument();
  });
});
