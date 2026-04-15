import { useMemo } from 'react';
import { useQuery } from '@tanstack/react-query';
import { canReadPolicy, useAuthStore } from '@/features/auth/store';
import { isAuthorizationError } from '@/services/api/client';
import {
  listPolicyExceptions,
  listPolicyHits,
  type PolicyExceptionListQuery,
  type PolicyHitListQuery
} from '@/services/securityPolicy';

export const usePolicyRollout = (policyId?: string) => {
  const user = useAuthStore((state) => state.user);
  const canRead = canReadPolicy(user);

  const hitQuery = useMemo<PolicyHitListQuery>(() => ({ policyId }), [policyId]);
  const exceptionQuery = useMemo<PolicyExceptionListQuery>(() => ({ policyId }), [policyId]);

  const hitsQuery = useQuery({
    queryKey: ['securityPolicy', 'hits', policyId ?? 'all'],
    enabled: canRead && Boolean(policyId),
    queryFn: () => listPolicyHits(hitQuery)
  });

  const exceptionsQuery = useQuery({
    queryKey: ['securityPolicy', 'exceptions', policyId ?? 'all'],
    enabled: canRead && Boolean(policyId),
    queryFn: () => listPolicyExceptions(exceptionQuery)
  });

  const permissionChanged =
    isAuthorizationError(hitsQuery.error) || isAuthorizationError(exceptionsQuery.error);

  return {
    canRead,
    hitsQuery,
    exceptionsQuery,
    permissionChanged
  };
};
