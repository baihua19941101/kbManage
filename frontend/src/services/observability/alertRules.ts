import { fetchJSON } from '@/services/api/client';
import type { ObservabilityAlertRuleDTO } from '@/services/api/types';
import { buildObservabilityQuery } from '@/services/observability/query';

export type AlertRulePayload = {
  name: string;
  description?: string;
  severity?: 'info' | 'warning' | 'critical';
  scopeSnapshot?: string;
  conditionExpression: string;
  evaluationWindow?: string;
  notificationStrategy?: string;
  status?: 'enabled' | 'disabled';
};

type AlertRuleListResponse = {
  items: ObservabilityAlertRuleDTO[];
};

export const listAlertRules = async (status?: 'enabled' | 'disabled') => {
  const suffix = buildObservabilityQuery({ status });
  return fetchJSON<AlertRuleListResponse>(
    `/observability/alert-rules${suffix ? `?${suffix}` : ''}`
  );
};

export const createAlertRule = async (payload: AlertRulePayload) => {
  return fetchJSON<ObservabilityAlertRuleDTO>('/observability/alert-rules', {
    method: 'POST',
    body: JSON.stringify(payload)
  });
};

export const updateAlertRule = async (ruleId: string | number, payload: AlertRulePayload) => {
  return fetchJSON<ObservabilityAlertRuleDTO>(`/observability/alert-rules/${ruleId}`, {
    method: 'PUT',
    body: JSON.stringify(payload)
  });
};

export const deleteAlertRule = async (ruleId: string | number) => {
  return fetchJSON<Record<string, never>>(`/observability/alert-rules/${ruleId}`, {
    method: 'DELETE'
  });
};
