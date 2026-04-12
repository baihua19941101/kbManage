import { fetchJSON } from '@/services/api/client';
import type { ObservabilityAlertDTO, ObservabilityHandlingRecordDTO } from '@/services/api/types';
import { buildObservabilityQuery } from '@/services/observability/query';

export type AlertListParams = {
  clusterId?: string;
  workspaceId?: string;
  projectId?: string;
  status?: 'firing' | 'acknowledged' | 'silenced' | 'resolved';
  severity?: 'info' | 'warning' | 'critical';
  resourceKind?: string;
  handlerId?: string;
  startAt?: string;
  endAt?: string;
  limit?: number;
};

export type AlertListResponse = {
  items: ObservabilityAlertDTO[];
};

export const listAlerts = async (params: AlertListParams = {}) => {
  const suffix = buildObservabilityQuery(params);
  return fetchJSON<AlertListResponse>(`/observability/alerts${suffix ? `?${suffix}` : ''}`);
};

export const getAlert = async (alertId: string | number) => {
  return fetchJSON<ObservabilityAlertDTO>(`/observability/alerts/${alertId}`);
};

export const acknowledgeAlert = async (
  alertId: string | number,
  note?: string
) => {
  return fetchJSON<ObservabilityAlertDTO>(`/observability/alerts/${alertId}/acknowledge`, {
    method: 'POST',
    body: JSON.stringify({ note: note ?? '' })
  });
};

export const createAlertHandlingRecord = async (
  alertId: string | number,
  payload: {
    actionType: string;
    content?: string;
  }
) => {
  return fetchJSON<ObservabilityHandlingRecordDTO>(
    `/observability/alerts/${alertId}/handling-records`,
    {
      method: 'POST',
      body: JSON.stringify(payload)
    }
  );
};
