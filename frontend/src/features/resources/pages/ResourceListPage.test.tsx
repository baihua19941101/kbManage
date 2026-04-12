import '@testing-library/jest-dom/vitest';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import { ResourceListPage } from '@/features/resources/pages/ResourceListPage';
import { listResources, type ResourceListItem } from '@/services/resources';

vi.mock('@/services/resources', () => ({
  listResources: vi.fn()
}));

const installAntdDomShims = () => {
  const nativeGetComputedStyle = window.getComputedStyle.bind(window);
  Object.defineProperty(window, 'getComputedStyle', {
    writable: true,
    value: ((elt: Element) => nativeGetComputedStyle(elt)) as typeof window.getComputedStyle
  });

  if (!window.matchMedia) {
    Object.defineProperty(window, 'matchMedia', {
      writable: true,
      value: vi.fn().mockImplementation((query: string) => ({
        matches: false,
        media: query,
        onchange: null,
        addListener: vi.fn(),
        removeListener: vi.fn(),
        addEventListener: vi.fn(),
        removeEventListener: vi.fn(),
        dispatchEvent: vi.fn()
      }))
    });
  }

  if (!window.scrollTo) {
    Object.defineProperty(window, 'scrollTo', {
      writable: true,
      value: vi.fn()
    });
  }

  if (!globalThis.ResizeObserver) {
    class ResizeObserverMock {
      observe() {}

      unobserve() {}

      disconnect() {}
    }
    vi.stubGlobal('ResizeObserver', ResizeObserverMock);
  }

  if (!window.requestAnimationFrame) {
    Object.defineProperty(window, 'requestAnimationFrame', {
      writable: true,
      value: (callback: FrameRequestCallback) => window.setTimeout(() => callback(0), 16)
    });
  }

  if (!window.cancelAnimationFrame) {
    Object.defineProperty(window, 'cancelAnimationFrame', {
      writable: true,
      value: (id: number) => window.clearTimeout(id)
    });
  }
};

describe('ResourceListPage', () => {
  beforeAll(() => {
    installAntdDomShims();
  });

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('filters table by keyword and opens resource detail drawer', async () => {
    const resources: ResourceListItem[] = [
      {
        id: 'res-1',
        cluster: 'prod-cn',
        namespace: 'payments',
        resourceType: 'Deployment',
        name: 'payment-api',
        status: 'Running',
        labels: { app: 'payment-api', env: 'prod' },
        updatedAt: '2026-04-09 11:20'
      },
      {
        id: 'res-2',
        cluster: 'prod-cn',
        namespace: 'gateway',
        resourceType: 'Service',
        name: 'edge-gateway',
        status: 'Running',
        labels: { app: 'gateway', env: 'prod' },
        updatedAt: '2026-04-09 11:10'
      },
      {
        id: 'res-3',
        cluster: 'staging-us',
        namespace: 'payments',
        resourceType: 'Pod',
        name: 'payment-api-66d8b87f74-z2vfh',
        status: 'Pending',
        labels: { app: 'payment-api', env: 'staging' },
        updatedAt: '2026-04-09 10:58'
      }
    ];

    vi.mocked(listResources).mockImplementation(async (query = {}) => {
      if (query.keyword === 'z2vfh') {
        return [resources[2]];
      }

      return resources;
    });

    render(
      <MemoryRouter>
        <ResourceListPage />
      </MemoryRouter>
    );

    await waitFor(() => {
      expect(listResources).toHaveBeenCalledWith({});
    });
    await waitFor(() => {
      expect(screen.getByText('payment-api')).toBeInTheDocument();
    });
    expect(screen.getByText('edge-gateway')).toBeInTheDocument();

    fireEvent.change(screen.getByPlaceholderText('名称/标签关键字'), {
      target: { value: 'z2vfh' }
    });

    await waitFor(() => {
      expect(listResources).toHaveBeenLastCalledWith({ keyword: 'z2vfh' });
    });
    await waitFor(() => {
      expect(screen.getByText('payment-api-66d8b87f74-z2vfh')).toBeInTheDocument();
    });
    expect(screen.queryByText('edge-gateway')).not.toBeInTheDocument();

    fireEvent.click(screen.getByRole('button', { name: '查看详情' }));

    await waitFor(() => {
      expect(
        screen.getByText('资源详情：payment-api-66d8b87f74-z2vfh')
      ).toBeInTheDocument();
    });
  });
});
