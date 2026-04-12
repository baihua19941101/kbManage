import { fetchJSON } from '@/services/api/client';
import type { ObservabilitySilenceWindowDTO } from '@/services/api/types';
import { buildObservabilityQuery } from '@/services/observability/query';

export type SilencePayload = {
  name: string;
  reason?: string;
  scopeSnapshot?: string;
  startsAt: string;
  endsAt: string;
};

type SilenceListResponse = {
  items: ObservabilitySilenceWindowDTO[];
};

export const listSilences = async (
  status?: 'scheduled' | 'active' | 'expired' | 'canceled'
) => {
  const suffix = buildObservabilityQuery({ status });
  return fetchJSON<SilenceListResponse>(`/observability/silences${suffix ? `?${suffix}` : ''}`);
};

export const createSilence = async (payload: SilencePayload) => {
  return fetchJSON<ObservabilitySilenceWindowDTO>('/observability/silences', {
    method: 'POST',
    body: JSON.stringify(payload)
  });
};

export const cancelSilence = async (silenceId: string | number) => {
  return fetchJSON<ObservabilitySilenceWindowDTO>(`/observability/silences/${silenceId}`, {
    method: 'DELETE'
  });
};
