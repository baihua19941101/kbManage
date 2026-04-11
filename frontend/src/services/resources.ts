import { fetchJSON } from '@/services/api/client';
import type {
  ListResourcesResponseDTO,
  ResourceInventoryDTO,
  ResourceListQueryDTO
} from '@/services/api/types';

export type ResourceListItem = {
  id: string;
  cluster: string;
  namespace: string;
  resourceType: string;
  name: string;
  status: string;
  labels: Record<string, string>;
  updatedAt: string;
};

const isRecord = (value: unknown): value is Record<string, unknown> =>
  typeof value === 'object' && value !== null;

const normalizeText = (value: unknown): string | undefined => {
  if (typeof value === 'string') {
    const trimmed = value.trim();
    return trimmed.length > 0 ? trimmed : undefined;
  }

  if (typeof value === 'number' || typeof value === 'bigint') {
    return String(value);
  }

  return undefined;
};

const pickText = (item: ResourceInventoryDTO, key: string): string | undefined => {
  const source = item as Record<string, unknown>;
  return normalizeText(source[key]);
};

const pickLabels = (item: ResourceInventoryDTO): Record<string, string> => {
  const source = item as Record<string, unknown>;
  const raw = source.labels;
  if (!isRecord(raw)) {
    return {};
  }

  const labels: Record<string, string> = {};
  for (const [key, value] of Object.entries(raw)) {
    const normalized = normalizeText(value);
    if (normalized) {
      labels[key] = normalized;
    }
  }
  return labels;
};

const mapStatus = (rawStatus: string | undefined): string => {
  const value = rawStatus?.trim().toLowerCase();
  if (!value) {
    return 'Unknown';
  }

  if (value === 'healthy' || value === 'running') {
    return 'Running';
  }
  if (value === 'degraded' || value === 'unhealthy') {
    return 'Degraded';
  }
  if (value === 'pending') {
    return 'Pending';
  }
  if (value === 'unknown') {
    return 'Unknown';
  }

  return `${value.slice(0, 1).toUpperCase()}${value.slice(1)}`;
};

const mapResourceItem = (item: ResourceInventoryDTO): ResourceListItem => {
  const cluster = pickText(item, 'clusterId') || pickText(item, 'cluster') || '-';
  const namespace = pickText(item, 'namespace') || '-';
  const resourceType = pickText(item, 'kind') || 'Unknown';
  const name = pickText(item, 'name') || 'unknown-resource';
  const status = mapStatus(pickText(item, 'status') || pickText(item, 'health'));

  const id = pickText(item, 'id') || `${cluster}/${namespace}/${resourceType}/${name}`;
  const updatedAt = pickText(item, 'updatedAt') || '-';

  return {
    id,
    cluster,
    namespace,
    resourceType,
    name,
    status,
    labels: pickLabels(item),
    updatedAt
  };
};

const buildQueryString = (query: ResourceListQueryDTO): string => {
  const search = new URLSearchParams();

  if (query.workspaceId) search.set('workspaceId', query.workspaceId);
  if (query.projectId) search.set('projectId', query.projectId);
  if (query.namespace) search.set('namespace', query.namespace);
  if (query.kind) search.set('kind', query.kind);
  if (query.keyword) search.set('keyword', query.keyword);
  if (query.health) search.set('health', query.health);
  if (typeof query.limit === 'number') search.set('limit', String(query.limit));
  if (typeof query.offset === 'number') search.set('offset', String(query.offset));

  const serialized = search.toString();
  return serialized ? `?${serialized}` : '';
};

export const listResources = async (
  query: ResourceListQueryDTO = {}
): Promise<ResourceListItem[]> => {
  const queryString = buildQueryString(query);
  const path = query.clusterId
    ? `/clusters/${encodeURIComponent(query.clusterId)}/resources${queryString}`
    : `/resources${queryString}`;
  const response = await fetchJSON<ListResourcesResponseDTO>(path, {
    method: 'GET'
  });
  const items = Array.isArray(response.items) ? response.items : [];

  return items.map(mapResourceItem);
};
