import { fetchJSON } from '@/services/api/client';
import type { ObservabilityEventDTO } from '@/services/api/types';
import { buildObservabilityQuery } from '@/services/observability/query';

export type EventQueryParams = {
  clusterId?: string;
  workspaceId?: string;
  projectId?: string;
  namespace?: string;
  resourceKind?: string;
  resourceName?: string;
  eventType?: string;
  startAt?: string;
  endAt?: string;
  limit?: number;
};

export type EventListResponse = {
  items: ObservabilityEventDTO[];
};

export const listObservabilityEvents = async (params: EventQueryParams = {}) => {
  const queryString = buildObservabilityQuery(params);
  return fetchJSON<EventListResponse>(`/observability/events${queryString ? `?${queryString}` : ''}`);
};
