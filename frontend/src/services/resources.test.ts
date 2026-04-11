import { fetchJSON } from '@/services/api/client';
import { listResources } from '@/services/resources';

vi.mock('@/services/api/client', () => ({
  fetchJSON: vi.fn()
}));

describe('listResources', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('maps contract fields and builds cluster-scoped query path', async () => {
    vi.mocked(fetchJSON).mockResolvedValue({
      items: [
        {
          id: 'res-1',
          clusterId: 'prod-cn',
          cluster: 'legacy-cluster',
          namespace: 'payments',
          kind: 'Deployment',
          name: 'payment-api',
          health: 'healthy',
          labels: {
            app: 'payment-api',
            revision: 3,
            empty: '   '
          },
          updatedAt: '2026-04-09T11:20:00Z'
        }
      ]
    });

    const result = await listResources({
      clusterId: 'prod-cn',
      namespace: 'payments',
      kind: 'Deployment',
      keyword: 'payment',
      health: 'healthy',
      limit: 20,
      offset: 5
    });

    expect(fetchJSON).toHaveBeenCalledWith(
      '/clusters/prod-cn/resources?namespace=payments&kind=Deployment&keyword=payment&health=healthy&limit=20&offset=5',
      { method: 'GET' }
    );
    expect(result).toEqual([
      {
        id: 'res-1',
        cluster: 'prod-cn',
        namespace: 'payments',
        resourceType: 'Deployment',
        name: 'payment-api',
        status: 'Running',
        labels: {
          app: 'payment-api',
          revision: '3'
        },
        updatedAt: '2026-04-09T11:20:00Z'
      }
    ]);
  });

  it('ignores legacy uppercase fields and keeps safe defaults', async () => {
    vi.mocked(fetchJSON).mockResolvedValue({
      items: [
        {
          ID: 'legacy-id',
          Cluster: 'legacy-cluster',
          Namespace: 'legacy-namespace',
          Kind: 'Pod',
          Name: 'legacy-name',
          Status: 'running',
          Labels: {
            app: 'legacy'
          },
          UpdatedAt: '2026-04-09T11:20:00Z'
        }
      ]
    });

    const result = await listResources();

    expect(result).toEqual([
      {
        id: '-/-/Unknown/unknown-resource',
        cluster: '-',
        namespace: '-',
        resourceType: 'Unknown',
        name: 'unknown-resource',
        status: 'Unknown',
        labels: {},
        updatedAt: '-'
      }
    ]);
  });

  it('accepts only items response field from contract', async () => {
    vi.mocked(fetchJSON).mockResolvedValue({
      Items: [
        {
          id: 'res-2',
          clusterId: 'prod-cn',
          namespace: 'payments',
          kind: 'Pod',
          name: 'payment-api-0',
          status: 'running',
          labels: {},
          updatedAt: '2026-04-09T11:20:00Z'
        }
      ]
    });

    const result = await listResources();

    expect(result).toEqual([]);
  });
});
