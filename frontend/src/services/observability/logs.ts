import { fetchJSON } from '@/services/api/client';
import type { ObservabilityLogEntryDTO } from '@/services/api/types';
import { buildObservabilityQuery } from '@/services/observability/query';

export type LogQueryParams = {
  clusterId?: string;
  workspaceId?: string;
  projectId?: string;
  namespace?: string;
  resourceKind?: string;
  resourceName?: string;
  workload?: string;
  pod?: string;
  container?: string;
  keyword?: string;
  startAt?: string;
  endAt?: string;
  limit?: number;
};

export type LogQueryResponse = {
  queryId: string;
  status: string;
  items: ObservabilityLogEntryDTO[];
  dataFreshness?: string;
};

export const queryObservabilityLogs = async (params: LogQueryParams = {}) => {
  const queryString = buildObservabilityQuery(params);
  return fetchJSON<LogQueryResponse>(
    `/observability/logs/query${queryString ? `?${queryString}` : ''}`
  );
};
