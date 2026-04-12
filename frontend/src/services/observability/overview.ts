import { fetchJSON } from '@/services/api/client';
import type { ObservabilityOverviewDTO } from '@/services/api/types';
import { buildObservabilityQuery } from '@/services/observability/query';

export type ObservabilityOverviewQuery = {
  clusterId?: string;
  startAt?: string;
  endAt?: string;
};

export const getObservabilityOverview = async (
  query: ObservabilityOverviewQuery = {}
): Promise<ObservabilityOverviewDTO> => {
  const suffix = buildObservabilityQuery({
    clusterIds: query.clusterId,
    startAt: query.startAt,
    endAt: query.endAt
  });
  return fetchJSON<ObservabilityOverviewDTO>(
    `/observability/overview${suffix ? `?${suffix}` : ''}`
  );
};
