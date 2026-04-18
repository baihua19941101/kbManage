import { useMemo } from 'react';
import { useQuery } from '@tanstack/react-query';
import {
  clusterLifecycleQueryKeys,
  groupCapabilityMatrixByDomain,
  listCapabilityOwners,
  listDriverCapabilities
} from '@/services/clusterLifecycle';

export const useCapabilityMatrix = (driverId?: string) => {
  const query = useQuery({
    queryKey: clusterLifecycleQueryKeys.capabilityMatrix(driverId),
    enabled: Boolean(driverId),
    queryFn: () => listDriverCapabilities(driverId || '')
  });

  const grouped = useMemo(
    () => groupCapabilityMatrixByDomain(query.data?.items || []),
    [query.data?.items]
  );

  const owners = useMemo(() => listCapabilityOwners(query.data?.items || []), [query.data?.items]);

  return {
    ...query,
    grouped,
    owners
  };
};
