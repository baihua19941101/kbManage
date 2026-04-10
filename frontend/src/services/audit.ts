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
  result?: string;
  eventType?: string;
};

export type ListAuditEventsResponse = {
  items: AuditEvent[];
};

export type AuditExportFormat = 'csv' | 'json';

export type AuditExportRequest = AuditEventFilters & {
  from: string;
  to: string;
  format: AuditExportFormat;
};

export type AuditExportResponse = {
  taskId: string;
  status: 'queued' | 'mocked';
};

class ApiError extends Error {
  status: number;

  constructor(status: number, message: string) {
    super(message);
    this.status = status;
  }
}

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api/v1';

const wait = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms));

const isFallbackStatus = (status: number) =>
  [404, 405, 500, 501, 502, 503, 504].includes(status);

const shouldUseFallback = (error: unknown): boolean => {
  if (error instanceof ApiError) {
    return isFallbackStatus(error.status);
  }

  return error instanceof TypeError;
};

const fetchJSON = async <T>(path: string, init?: RequestInit): Promise<T> => {
  const response = await fetch(`${API_BASE_URL}${path}`, {
    headers: {
      'Content-Type': 'application/json',
      ...(init?.headers || {})
    },
    ...init
  });

  if (!response.ok) {
    const text = await response.text();
    throw new ApiError(
      response.status,
      text || `Request failed with status ${response.status}`
    );
  }

  if (response.status === 204) {
    return {} as T;
  }

  const contentType = response.headers.get('content-type') || '';
  if (!contentType.includes('application/json')) {
    return {} as T;
  }

  return (await response.json()) as T;
};

const buildBackendQuery = (filters: AuditEventFilters): string => {
  const search = new URLSearchParams();

  if (filters.from) search.set('startAt', filters.from);
  if (filters.to) search.set('endAt', filters.to);
  if (filters.actorUserId) search.set('actorId', filters.actorUserId);
  if (filters.eventType) search.set('action', filters.eventType);
  if (filters.result) search.set('outcome', filters.result);

  const query = search.toString();
  return query ? `?${query}` : '';
};

const mockAuditEvents: AuditEvent[] = [
  {
    id: 'ae-seed-1',
    eventType: 'operation.restart',
    actorUserId: 'alice',
    clusterId: 'prod-cn',
    scopeType: 'project',
    scopeId: 'payments',
    resourceKind: 'Deployment',
    resourceNamespace: 'payments',
    resourceName: 'payment-api',
    result: 'success',
    occurredAt: '2026-04-09T03:21:00.000Z'
  },
  {
    id: 'ae-seed-2',
    eventType: 'rbac.role_binding.update',
    actorUserId: 'admin',
    clusterId: 'platform',
    scopeType: 'workspace',
    scopeId: 'dev-team',
    result: 'success',
    occurredAt: '2026-04-09T03:10:00.000Z'
  }
];

const normalizeForSearch = (value?: string) => value?.trim().toLowerCase() || '';

const applyMockFilters = (filters: AuditEventFilters): AuditEvent[] => {
  const from = filters.from ? +new Date(filters.from) : Number.NEGATIVE_INFINITY;
  const to = filters.to ? +new Date(filters.to) : Number.POSITIVE_INFINITY;
  const actor = normalizeForSearch(filters.actorUserId);
  const result = normalizeForSearch(filters.result);
  const eventType = normalizeForSearch(filters.eventType);

  return [...mockAuditEvents]
    .filter((item) => {
      const occurredAt = +new Date(item.occurredAt || 0);
      const matchTime = occurredAt >= from && occurredAt <= to;
      const matchActor = !actor || normalizeForSearch(item.actorUserId).includes(actor);
      const matchResult = !result || normalizeForSearch(item.result).includes(result);
      const matchEventType =
        !eventType || normalizeForSearch(item.eventType).includes(eventType);

      return matchTime && matchActor && matchResult && matchEventType;
    })
    .sort((a, b) => +new Date(b.occurredAt || 0) - +new Date(a.occurredAt || 0));
};

export const listAuditEvents = async (
  filters: AuditEventFilters = {}
): Promise<ListAuditEventsResponse> => {
  try {
    const response = await fetchJSON<{ items: AuditEvent[] }>(
      `/audits/events${buildBackendQuery(filters)}`,
      {
        method: 'GET'
      }
    );

    return {
      items: (response.items || []).map((item) => ({
        ...item,
        eventType: item.eventType || item.action,
        result: item.result || item.outcome,
        occurredAt: item.occurredAt || item.createdAt
      }))
    };
  } catch (error) {
    if (!shouldUseFallback(error)) {
      throw error;
    }

    await wait(120);
    return {
      items: applyMockFilters(filters)
    };
  }
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

  try {
    const response = await fetchJSON<Partial<AuditExportResponse>>('/audits/exports', {
      method: 'POST',
      body: JSON.stringify(body)
    });

    return {
      taskId:
        typeof response.taskId === 'string' && response.taskId.length > 0
          ? response.taskId
          : `audit-export-${Date.now().toString(36)}`,
      status: 'queued'
    };
  } catch (error) {
    if (!shouldUseFallback(error)) {
      throw error;
    }

    await wait(120);
    return {
      taskId: `audit-export-mock-${Date.now().toString(36)}`,
      status: 'mocked'
    };
  }
};
