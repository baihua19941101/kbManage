import { fetchJSON } from '@/services/api/client';

const isRecord = (value: unknown): value is Record<string, unknown> =>
  typeof value === 'object' && value !== null;

const pick = (record: Record<string, unknown>, keys: string[]): unknown => {
  for (const key of keys) {
    if (key in record) {
      return record[key];
    }
  }
  return undefined;
};

const toText = (value: unknown): string | undefined => {
  if (typeof value === 'string') {
    const trimmed = value.trim();
    return trimmed.length > 0 ? trimmed : undefined;
  }

  if (typeof value === 'number' && Number.isFinite(value)) {
    return String(value);
  }

  return undefined;
};

const toNumber = (value: unknown): number | undefined => {
  if (typeof value === 'number' && Number.isFinite(value)) {
    return value;
  }

  if (typeof value === 'string') {
    const parsed = Number(value.trim());
    return Number.isFinite(parsed) ? parsed : undefined;
  }

  return undefined;
};

const extractApiServerFromKubeConfig = (kubeConfig: string): string | undefined => {
  const match = kubeConfig.match(/^\s*server:\s*([^\s]+)\s*$/im);
  if (!match) {
    return undefined;
  }

  const value = match[1]?.trim();
  return value && value.length > 0 ? value : undefined;
};

const normalizeCluster = (value: unknown, index = 0): Cluster => {
  const record = isRecord(value) ? value : {};
  const name = toText(pick(record, ['name', 'Name'])) || `cluster-${index + 1}`;
  const id =
    toText(pick(record, ['id', 'ID', 'clusterId', 'ClusterID'])) || `${name}-${index + 1}`;

  return {
    id,
    name,
    status: toText(pick(record, ['status', 'Status'])) || 'unknown',
    apiServer: toText(pick(record, ['apiServer', 'APIServer'])),
    namespaces:
      toNumber(
        pick(record, [
          'namespaces',
          'Namespaces',
          'namespaceCount',
          'NamespaceCount',
          'namespaceTotal',
          'NamespaceTotal'
        ])
      ) || 0
  };
};

const normalizeListPayload = (value: unknown): Cluster[] => {
  if (Array.isArray(value)) {
    return value.map((item, index) => normalizeCluster(item, index));
  }

  if (!isRecord(value)) {
    return [];
  }

  const items = pick(value, ['items', 'Items']);
  if (!Array.isArray(items)) {
    return [];
  }

  return items.map((item, index) => normalizeCluster(item, index));
};

export type Cluster = {
  id: string;
  name: string;
  status: string;
  apiServer?: string;
  namespaces: number;
};

export type CreateClusterPayload = {
  name: string;
  credentialPayload: string;
  credentialType?: 'kubeconfig' | 'token' | 'service-account';
  kubeConfig?: string;
  description?: string;
  apiServer?: string;
};

export const listClusters = async (): Promise<Cluster[]> => {
  const payload = await fetchJSON<unknown>('/clusters');
  return normalizeListPayload(payload);
};

export const createCluster = async (payload: CreateClusterPayload): Promise<Cluster> => {
  const name = payload.name.trim();
  const credentialPayload = payload.credentialPayload || payload.kubeConfig || '';
  const apiServer =
    payload.apiServer?.trim() || extractApiServerFromKubeConfig(credentialPayload) || '';
  const credentialType = payload.credentialType?.trim() || 'kubeconfig';

  const response = await fetchJSON<unknown>('/clusters', {
    method: 'POST',
    body: JSON.stringify({
      name,
      apiServer,
      description: payload.description?.trim() || '',
      credentialType,
      credentialPayload
    })
  });

  return normalizeCluster(response);
};
