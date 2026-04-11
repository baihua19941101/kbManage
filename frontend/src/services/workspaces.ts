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

const normalizeWorkspace = (value: unknown, index = 0): Workspace => {
  const record = isRecord(value) ? value : {};
  const name = toText(pick(record, ['name', 'Name'])) || `workspace-${index + 1}`;
  const id =
    toText(pick(record, ['id', 'ID', 'workspaceId', 'WorkspaceID'])) ||
    `${name}-${index + 1}`;

  return {
    id,
    name,
    description: toText(pick(record, ['description', 'Description'])) || '-'
  };
};

const normalizeListPayload = (value: unknown): Workspace[] => {
  if (Array.isArray(value)) {
    return value.map((item, index) => normalizeWorkspace(item, index));
  }

  if (!isRecord(value)) {
    return [];
  }

  const items = pick(value, ['items', 'Items', 'data', 'Data', 'list', 'List']);
  if (Array.isArray(items)) {
    return items.map((item, index) => normalizeWorkspace(item, index));
  }

  return [];
};

export type Workspace = {
  id: string;
  name: string;
  description: string;
};

export type CreateWorkspacePayload = {
  name: string;
  description?: string;
};

export const listWorkspaces = async (): Promise<Workspace[]> => {
  const payload = await fetchJSON<unknown>('/workspaces', {
    method: 'GET'
  });
  return normalizeListPayload(payload);
};

export const createWorkspace = async (
  payload: CreateWorkspacePayload
): Promise<Workspace> => {
  const name = payload.name.trim();
  const description = payload.description?.trim();

  const response = await fetchJSON<unknown>('/workspaces', {
    method: 'POST',
    body: JSON.stringify({
      name,
      ...(description ? { description } : {})
    })
  });

  return normalizeWorkspace(response);
};
