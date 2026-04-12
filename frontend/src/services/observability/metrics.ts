import { fetchJSON } from '@/services/api/client';
import type { ObservabilityMetricSeriesDTO } from '@/services/api/types';
import { buildObservabilityQuery } from '@/services/observability/query';

export type MetricSeriesParams = {
  clusterId?: string;
  namespace?: string;
  subjectType: 'cluster' | 'node' | 'namespace' | 'workload' | 'pod';
  subjectRef: string;
  metricKey: string;
  startAt?: string;
  endAt?: string;
  step?: string;
};

export const queryMetricSeries = async (params: MetricSeriesParams) => {
  const queryString = buildObservabilityQuery(params);
  return fetchJSON<ObservabilityMetricSeriesDTO>(
    `/observability/metrics/series${queryString ? `?${queryString}` : ''}`
  );
};
