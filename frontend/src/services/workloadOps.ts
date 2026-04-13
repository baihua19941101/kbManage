import { fetchJSON } from '@/services/api/client';
import type {
  BatchOperationTaskDTO,
  ReleaseRevisionDTO,
  SubmitBatchOperationRequestDTO,
  SubmitWorkloadActionRequestDTO,
  TerminalSessionDTO,
  WorkloadActionDTO,
  WorkloadInstanceDTO,
  WorkloadOperationsViewDTO
} from '@/services/api/types';

export type WorkloadResourceQuery = {
  clusterId: number;
  namespace: string;
  resourceKind: 'Deployment' | 'StatefulSet' | 'DaemonSet';
  resourceName: string;
};

export type CreateTerminalSessionRequest = {
  clusterId: number;
  namespace: string;
  podName?: string;
  containerName?: string;
  workloadKind: 'Deployment' | 'StatefulSet' | 'DaemonSet';
  workloadName: string;
  cols?: number;
  rows?: number;
};

const toQueryString = (query: Record<string, string | number | undefined>) => {
  const params = new URLSearchParams();
  Object.entries(query).forEach(([key, value]) => {
    if (value === undefined || value === null || value === '') {
      return;
    }
    params.set(key, String(value));
  });
  return params.toString();
};

export const getWorkloadOpsContext = async (query: WorkloadResourceQuery) => {
  const qs = toQueryString(query);
  return fetchJSON<WorkloadOperationsViewDTO>(`/workload-ops/resources/context?${qs}`);
};

export const listWorkloadOpsInstances = async (query: WorkloadResourceQuery) => {
  const qs = toQueryString(query);
  return fetchJSON<{ items: WorkloadInstanceDTO[] }>(`/workload-ops/resources/instances?${qs}`);
};

export const listWorkloadOpsRevisions = async (query: WorkloadResourceQuery) => {
  const qs = toQueryString(query);
  return fetchJSON<{ items: ReleaseRevisionDTO[] }>(`/workload-ops/resources/revisions?${qs}`);
};

export const submitWorkloadAction = async (payload: SubmitWorkloadActionRequestDTO) => {
  return fetchJSON<WorkloadActionDTO>('/workload-ops/actions', {
    method: 'POST',
    body: JSON.stringify(payload)
  });
};

export const getWorkloadAction = async (actionId: number) => {
  return fetchJSON<WorkloadActionDTO>(`/workload-ops/actions/${actionId}`);
};

export const submitBatchOperation = async (payload: SubmitBatchOperationRequestDTO) => {
  return fetchJSON<BatchOperationTaskDTO>('/workload-ops/batches', {
    method: 'POST',
    body: JSON.stringify(payload)
  });
};

export const getBatchOperation = async (batchId: number) => {
  return fetchJSON<BatchOperationTaskDTO>(`/workload-ops/batches/${batchId}`);
};

export const createTerminalSession = async (
  payload: CreateTerminalSessionRequest
) => {
  return fetchJSON<TerminalSessionDTO>('/workload-ops/terminal/sessions', {
    method: 'POST',
    body: JSON.stringify(payload)
  });
};

export const getTerminalSession = async (sessionId: number) => {
  return fetchJSON<TerminalSessionDTO>(`/workload-ops/terminal/sessions/${sessionId}`);
};

export const closeTerminalSession = async (sessionId: number) => {
  return fetchJSON<Record<string, never>>(`/workload-ops/terminal/sessions/${sessionId}`, {
    method: 'DELETE'
  });
};
