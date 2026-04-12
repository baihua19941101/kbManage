import { fetchJSON } from '@/services/api/client';
import type { ObservabilityNotificationTargetDTO } from '@/services/api/types';

export type NotificationTargetPayload = {
  name: string;
  targetType: string;
  configRef?: string;
  scopeSnapshot?: string;
  status?: 'active' | 'disabled';
};

type NotificationTargetListResponse = {
  items: ObservabilityNotificationTargetDTO[];
};

export const listNotificationTargets = async () => {
  return fetchJSON<NotificationTargetListResponse>('/observability/notification-targets');
};

export const createNotificationTarget = async (payload: NotificationTargetPayload) => {
  return fetchJSON<ObservabilityNotificationTargetDTO>('/observability/notification-targets', {
    method: 'POST',
    body: JSON.stringify(payload)
  });
};

export const updateNotificationTarget = async (
  targetId: string | number,
  payload: NotificationTargetPayload
) => {
  return fetchJSON<ObservabilityNotificationTargetDTO>(
    `/observability/notification-targets/${targetId}`,
    {
      method: 'PUT',
      body: JSON.stringify(payload)
    }
  );
};

export const deleteNotificationTarget = async (targetId: string | number) => {
  return fetchJSON<Record<string, never>>(`/observability/notification-targets/${targetId}`, {
    method: 'DELETE'
  });
};
