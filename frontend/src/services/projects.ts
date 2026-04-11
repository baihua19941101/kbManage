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

  if (typeof value === 'bigint') {
    return String(value);
  }

  return undefined;
};

const pickText = (record: Record<string, unknown>, keys: string[]): string | undefined =>
  toText(pick(record, keys));

const pickItems = (payload: unknown): unknown[] => {
  if (Array.isArray(payload)) {
    return payload;
  }

  if (!isRecord(payload)) {
    return [];
  }

  const items = pick(payload, ['items', 'Items', 'list', 'List', 'data', 'Data']);
  return Array.isArray(items) ? items : [];
};

const normalizeCode = (name: string): string =>
  name
    .trim()
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/^-+|-+$/g, '') || 'project';

export type WorkspaceOption = {
  id: string;
  name: string;
};

export type ProjectItem = {
  id: string;
  workspaceId: string;
  workspaceName?: string;
  name: string;
  owner: string;
};

type CreateProjectPayload = {
  name: string;
  owner?: string;
};

const mapWorkspace = (value: unknown, index = 0): WorkspaceOption => {
  const record = isRecord(value) ? value : {};
  const name = pickText(record, ['name', 'Name', 'code', 'Code']) || `workspace-${index + 1}`;
  const id =
    pickText(record, [
      'id',
      'ID',
      'workspaceId',
      'WorkspaceId',
      'workspaceID',
      'WorkspaceID'
    ]) || name;

  return { id, name };
};

const mapProject = (value: unknown, index = 0, workspaceIdFallback?: string): ProjectItem => {
  const record = isRecord(value) ? value : {};
  const name = pickText(record, ['name', 'Name', 'code', 'Code']) || `project-${index + 1}`;
  const id =
    pickText(record, ['id', 'ID', 'projectId', 'ProjectId', 'projectID', 'ProjectID']) ||
    `${workspaceIdFallback || 'workspace'}-${name}-${index + 1}`;

  const workspaceId =
    pickText(record, [
      'workspaceId',
      'WorkspaceId',
      'workspaceID',
      'WorkspaceID',
      'workspace_id',
      'workspaceid'
    ]) ||
    workspaceIdFallback ||
    '';

  return {
    id,
    workspaceId,
    workspaceName: pickText(record, ['workspace', 'Workspace', 'workspaceName', 'WorkspaceName']),
    name,
    owner:
      pickText(record, [
        'owner',
        'Owner',
        'ownerName',
        'OwnerName',
        'creator',
        'Creator',
        'description',
        'Description'
      ]) || '-'
  };
};

export const listWorkspaces = async (): Promise<WorkspaceOption[]> => {
  const payload = await fetchJSON<unknown>('/workspaces', { method: 'GET' });
  return pickItems(payload).map((item, index) => mapWorkspace(item, index));
};

export const listProjectsByWorkspace = async (
  workspaceId: string
): Promise<ProjectItem[]> => {
  const payload = await fetchJSON<unknown>(
    `/workspaces/${encodeURIComponent(workspaceId)}/projects`,
    { method: 'GET' }
  );
  return pickItems(payload).map((item, index) => mapProject(item, index, workspaceId));
};

export const createProject = async (
  workspaceId: string,
  payload: CreateProjectPayload
): Promise<ProjectItem> => {
  const normalizedName = payload.name.trim();
  const description = payload.owner?.trim() || '';

  const response = await fetchJSON<unknown>(
    `/workspaces/${encodeURIComponent(workspaceId)}/projects`,
    {
      method: 'POST',
      body: JSON.stringify({
        name: normalizedName,
        code: normalizeCode(normalizedName),
        description
      })
    }
  );

  return mapProject(response, 0, workspaceId);
};
