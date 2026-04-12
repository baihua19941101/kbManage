import '@testing-library/jest-dom/vitest';
import { render, screen } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { MemoryRouter } from 'react-router-dom';
import { SilenceWindowPage } from '@/features/observability/pages/SilenceWindowPage';
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

vi.mock('@/services/observability/silences', () => ({
  listSilences: vi.fn().mockResolvedValue({
    items: [
      {
        id: 'silence-1',
        name: 'release window',
        status: 'active',
        startsAt: '2026-01-01T00:00:00Z',
        endsAt: '2026-01-01T01:00:00Z'
      }
    ]
  }),
  createSilence: vi.fn().mockResolvedValue({}),
  cancelSilence: vi.fn().mockResolvedValue({})
}));

vi.mock('@/services/observability/notificationTargets', () => ({
  listNotificationTargets: vi.fn().mockResolvedValue({
    items: [{ id: 'target-1', name: 'oncall', targetType: 'webhook', status: 'active' }]
  }),
  createNotificationTarget: vi.fn().mockResolvedValue({})
}));

describe('SilenceWindowPage', () => {
  beforeAll(() => {
    installAntdDomShims();
  });

  it('renders silence and notification sections', async () => {
    const client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
    render(
      <QueryClientProvider client={client}>
        <MemoryRouter>
          <SilenceWindowPage />
        </MemoryRouter>
      </QueryClientProvider>
    );

    expect(screen.getByText('通知目标与静默窗口')).toBeInTheDocument();
    expect(await screen.findByText(/release window/i)).toBeInTheDocument();
    expect(await screen.findByText(/oncall/i)).toBeInTheDocument();
  });
});
