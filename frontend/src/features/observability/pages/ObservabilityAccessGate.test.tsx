import '@testing-library/jest-dom/vitest';
import { render, screen } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { MemoryRouter } from 'react-router-dom';
import { useAuthStore } from '@/features/auth/store';
import { ObservabilityOverviewPage } from '@/features/observability/pages/ObservabilityOverviewPage';
import { ApiError } from '@/services/api/client';
import { getObservabilityOverview } from '@/services/observability/overview';
import { queryMetricSeries } from '@/services/observability/metrics';
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
  getObservabilityOverview: vi.fn()
}));

vi.mock('@/services/observability/metrics', () => ({
  queryMetricSeries: vi.fn()
}));

vi.mock('@/features/observability/components/MetricsTrendChart', () => ({
  MetricsTrendChart: () => <div data-testid="metrics-trend-chart" />
}));

describe('ObservabilityAccessGate', () => {
  beforeAll(() => {
    installAntdDomShims();
  });

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('shows unauthorized empty state for users without observability read roles', async () => {
    useAuthStore.setState({
      isAuthenticated: true,
      accessToken: 'token',
      refreshToken: 'refresh-token',
      user: {
        id: 'user-no-ob',
        username: 'no-observer',
        platformRoles: []
      }
    });

    const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });

    render(
      <QueryClientProvider client={client}>
        <MemoryRouter>
          <ObservabilityOverviewPage />
        </MemoryRouter>
      </QueryClientProvider>
    );

    expect(
      screen.getByText('你暂无可观测访问权限，请联系管理员授予工作空间/项目范围。')
    ).toBeInTheDocument();
    expect(getObservabilityOverview).not.toHaveBeenCalled();
    expect(queryMetricSeries).not.toHaveBeenCalled();
  });

  it('shows permission changed warning when query returns authorization error', async () => {
    useAuthStore.setState({
      isAuthenticated: true,
      accessToken: 'token',
      refreshToken: 'refresh-token',
      user: {
        id: 'user-readonly',
        username: 'readonly-user',
        platformRoles: ['readonly']
      }
    });

    vi.mocked(getObservabilityOverview).mockRejectedValue(
      new ApiError(403, 'forbidden', {
        url: '/api/v1/observability/overview'
      })
    );
    vi.mocked(queryMetricSeries).mockResolvedValue({
      metricKey: 'cpu_usage',
      subjectType: 'cluster',
      subjectRef: 'all-clusters',
      points: []
    });

    const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });

    render(
      <QueryClientProvider client={client}>
        <MemoryRouter>
          <ObservabilityOverviewPage />
        </MemoryRouter>
      </QueryClientProvider>
    );

    expect(await screen.findByText('权限已变更')).toBeInTheDocument();
  });
});
