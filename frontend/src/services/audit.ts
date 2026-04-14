import { ApiError, fetchJSON } from '@/services/api/client';

export type AuditResult = 'success' | 'failed' | 'denied' | 'pending';

export type AuditEvent = {
  id: string;
  eventType?: string;
  action?: string;
  actorUserId?: string;
  clusterId?: string;
  scopeType?: string;
  scopeId?: string;
  resourceKind?: string;
  resourceNamespace?: string;
  resourceName?: string;
  result?: AuditResult | string;
  outcome?: string;
  occurredAt?: string;
  createdAt?: string;
};

export type AuditEventFilters = {
  from?: string;
  to?: string;
  actorUserId?: string;
  clusterId?: string;
  workspaceId?: string;
  projectId?: string;
  result?: string;
  eventType?: string;
  resource?: string;
  actionPrefix?: string;
};

export type ListAuditEventsResponse = {
  items: AuditEvent[];
};

export type AuditExportFormat = 'csv';

export type AuditExportRequest = AuditEventFilters & {
  from: string;
  to: string;
  format: AuditExportFormat;
};

export type AuditExportResponse = {
  taskId: string;
  status: 'pending' | 'running' | 'succeeded' | 'failed';
};

export type AuditExportTaskStatus = AuditExportResponse['status'];

export type AuditExportTask = {
  taskId: string;
  status: AuditExportTaskStatus;
  resultTotal?: number;
  downloadUrl?: string;
  errorMessage?: string;
  createdAt?: string;
  updatedAt?: string;
  completedAt?: string;
};

type UnknownRecord = Record<string, unknown>;

const toRecord = (value: unknown): UnknownRecord =>
  typeof value === 'object' && value !== null ? (value as UnknownRecord) : {};

const firstString = (record: UnknownRecord, keys: string[]): string | undefined => {
  for (const key of keys) {
    const value = record[key];
    if (typeof value === 'string') {
      const trimmed = value.trim();
      if (trimmed.length > 0) {
        return trimmed;
      }
    }
  }
  return undefined;
};

const firstNumber = (record: UnknownRecord, keys: string[]): number | undefined => {
  for (const key of keys) {
    const value = record[key];
    if (typeof value === 'number' && Number.isFinite(value)) {
      return value;
    }
    if (typeof value === 'string') {
      const parsed = Number(value);
      if (Number.isFinite(parsed)) {
        return parsed;
      }
    }
  }
  return undefined;
};

const normalizeExportStatus = (status: unknown): AuditExportTaskStatus | undefined => {
  if (typeof status !== 'string') {
    return undefined;
  }
  const normalized = status.trim().toLowerCase();
  return isAuditExportStatus(normalized) ? normalized : undefined;
};

const buildBackendQuery = (filters: AuditEventFilters): string => {
  const search = new URLSearchParams();

  if (filters.from) search.set('startAt', filters.from);
  if (filters.to) search.set('endAt', filters.to);
  if (filters.actorUserId) search.set('actorId', filters.actorUserId);
  if (filters.workspaceId) search.set('workspaceId', filters.workspaceId);
  if (filters.projectId) search.set('projectId', filters.projectId);
  if (filters.eventType) search.set('action', filters.eventType);
  if (filters.result) search.set('outcome', filters.result);
  if (filters.resource) search.set('resource', filters.resource);

  const query = search.toString();
  return query ? `?${query}` : '';
};

const isAuditExportStatus = (
  status: unknown
): status is AuditExportResponse['status'] =>
  status === 'pending' ||
  status === 'running' ||
  status === 'succeeded' ||
  status === 'failed';

const mapAuditEvent = (input: unknown): AuditEvent => {
  const record = toRecord(input);
  const eventType = firstString(record, ['eventType', 'EventType', 'action', 'Action']);
  const result = firstString(record, ['result', 'Result', 'outcome', 'Outcome']);
  const occurredAt = firstString(record, ['occurredAt', 'OccurredAt', 'createdAt', 'CreatedAt']);

  return {
    id: firstString(record, ['id', 'ID']) || '',
    eventType,
    action: firstString(record, ['action', 'Action']),
    actorUserId: firstString(record, ['actorUserId', 'ActorUserId', 'actorId', 'ActorId']),
    clusterId: firstString(record, ['clusterId', 'ClusterId', 'clusterID', 'ClusterID']),
    scopeType: firstString(record, ['scopeType', 'ScopeType']),
    scopeId: firstString(record, ['scopeId', 'ScopeId', 'scopeID', 'ScopeID']),
    resourceKind: firstString(record, ['resourceKind', 'ResourceKind']),
    resourceNamespace: firstString(record, ['resourceNamespace', 'ResourceNamespace']),
    resourceName: firstString(record, ['resourceName', 'ResourceName']),
    result,
    outcome: firstString(record, ['outcome', 'Outcome']),
    occurredAt,
    createdAt: firstString(record, ['createdAt', 'CreatedAt', 'occurredAt', 'OccurredAt'])
  };
};

const mapAuditExportResponse = (input: unknown): AuditExportResponse => {
  const record = toRecord(input);
  const taskId = firstString(record, ['taskId', 'TaskId', 'taskID', 'TaskID']) || '';
  const status = normalizeExportStatus(record.status ?? record.Status);

  if (!taskId) {
    throw new ApiError(500, 'Invalid audit export response: missing taskId', {
      url: '/audits/exports'
    });
  }

  if (!status) {
    throw new ApiError(500, 'Invalid audit export response: invalid status', {
      url: '/audits/exports'
    });
  }

  return { taskId, status };
};

const mapAuditExportTask = (input: unknown, taskId: string): AuditExportTask => {
  const record = toRecord(input);
  const mappedTaskId = firstString(record, ['taskId', 'TaskId', 'taskID', 'TaskID']) || taskId;
  const status = normalizeExportStatus(record.status ?? record.Status);

  if (!status) {
    throw new ApiError(500, 'Invalid audit export task response: invalid status', {
      url: `/audits/exports/${taskId}`
    });
  }

  return {
    taskId: mappedTaskId,
    status,
    resultTotal: firstNumber(record, ['resultTotal', 'ResultTotal']),
    downloadUrl: firstString(record, ['downloadUrl', 'DownloadUrl', 'DownloadURL']),
    errorMessage: firstString(record, ['errorMessage', 'ErrorMessage']),
    createdAt: firstString(record, ['createdAt', 'CreatedAt']),
    updatedAt: firstString(record, ['updatedAt', 'UpdatedAt']),
    completedAt: firstString(record, ['completedAt', 'CompletedAt'])
  };
};

export const listAuditEvents = async (
  filters: AuditEventFilters = {}
): Promise<ListAuditEventsResponse> => {
  const response = await fetchJSON<{ items?: unknown[]; Items?: unknown[] }>(
    `/audits/events${buildBackendQuery(filters)}`,
    {
      method: 'GET'
    }
  );

  const items = Array.isArray(response.items)
    ? response.items
    : Array.isArray(response.Items)
      ? response.Items
      : [];

  const mappedItems = items.map(mapAuditEvent).filter((item) => item.id.length > 0);
  const actionPrefix = filters.actionPrefix?.trim().toLowerCase();
  if (!actionPrefix) {
    return { items: mappedItems };
  }
  return {
    items: mappedItems.filter((item) => (item.action || item.eventType || '').toLowerCase().startsWith(actionPrefix))
  };
};

export const exportAuditEvents = async (
  payload: AuditExportRequest
): Promise<AuditExportResponse> => {
  const body = {
    startAt: payload.from,
    endAt: payload.to,
    actorId: payload.actorUserId,
    action: payload.eventType,
    outcome: payload.result,
    format: payload.format
  };

  const response = await fetchJSON<unknown>('/audits/exports', {
    method: 'POST',
    body: JSON.stringify(body)
  });

  return mapAuditExportResponse(response);
};

export const getAuditExportTask = async (taskId: string): Promise<AuditExportTask> => {
  const trimmed = taskId.trim();
  if (!trimmed) {
    throw new Error('taskId is required');
  }

  const response = await fetchJSON<unknown>(`/audits/exports/${encodeURIComponent(trimmed)}`, {
    method: 'GET'
  });

  return mapAuditExportTask(response, trimmed);
};
