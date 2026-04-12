import '@testing-library/jest-dom/vitest';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { MemoryRouter } from 'react-router-dom';
import { ClusterOverviewPage } from '@/features/clusters/pages/ClusterOverviewPage';
import { createCluster, listClusters, type Cluster } from '@/services/clusters';

vi.mock('@/services/clusters', () => ({
  listClusters: vi.fn(),
  createCluster: vi.fn()
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

describe('ClusterOverviewPage', () => {
  beforeAll(() => {
    installAntdDomShims();
  });

  beforeEach(() => {
    vi.clearAllMocks();
  });

  const renderPage = () => {
    const queryClient = new QueryClient({
      defaultOptions: {
        queries: { retry: false },
        mutations: { retry: false }
      }
    });

    return render(
      <QueryClientProvider client={queryClient}>
        <MemoryRouter>
          <ClusterOverviewPage />
        </MemoryRouter>
      </QueryClientProvider>
    );
  };

  it('loads clusters from api and submits onboarding through api', async () => {
    const initialClusters: Cluster[] = [
      { id: '1', name: 'prod-cn', status: 'healthy', namespaces: 12 },
      { id: '2', name: 'staging-us', status: 'unknown', namespaces: 5 }
    ];
    const createdCluster: Cluster = {
      id: '3',
      name: 'qa-cn',
      status: 'unknown',
      namespaces: 0
    };

    vi.mocked(listClusters)
      .mockResolvedValueOnce(initialClusters)
      .mockResolvedValueOnce([...initialClusters, createdCluster]);
    vi.mocked(createCluster).mockResolvedValue(createdCluster);

    renderPage();

    expect(await screen.findByText('prod-cn')).toBeInTheDocument();
    expect(screen.getByText('staging-us')).toBeInTheDocument();

    fireEvent.click(screen.getByRole('button', { name: '接入集群' }));
    expect(await screen.findByText('接入 Kubernetes 集群')).toBeInTheDocument();

    fireEvent.change(screen.getByLabelText('Cluster Name'), {
      target: { value: 'qa-cn' }
    });
    fireEvent.change(screen.getByLabelText('kubeconfig'), {
      target: {
        value:
          'apiVersion: v1\nclusters:\n- name: qa-cn\n  cluster:\n    server: https://k8s.example.com'
      }
    });
    fireEvent.click(screen.getByRole('button', { name: '提交接入' }));

    await waitFor(() => {
      expect(createCluster).toHaveBeenCalledWith(
        expect.objectContaining({
          name: 'qa-cn',
          credentialType: 'kubeconfig',
          credentialPayload:
            'apiVersion: v1\nclusters:\n- name: qa-cn\n  cluster:\n    server: https://k8s.example.com'
        }),
        expect.anything()
      );
    });

    await waitFor(() => {
      expect(listClusters).toHaveBeenCalledTimes(2);
    });

    expect(await screen.findByText('qa-cn')).toBeInTheDocument();
  });

  it('shows query error from clusters api on page', async () => {
    vi.mocked(listClusters).mockRejectedValue(new Error('query failed'));

    renderPage();

    expect(await screen.findByText('集群列表加载失败')).toBeInTheDocument();
    expect(screen.getByText('query failed')).toBeInTheDocument();
  });
});
