import { useQuery } from '@tanstack/react-query';
import { queryKeys } from '@/app/queryClient';
import { listDeliveryChecklist } from '@/services/enterprisePolish';

export const useDeliveryBundleChecklist = (bundleId?: string) =>
  useQuery({
    queryKey: queryKeys.enterprisePolish.deliveryChecklist(bundleId),
    enabled: Boolean(bundleId),
    queryFn: () => listDeliveryChecklist(bundleId || '')
  });
