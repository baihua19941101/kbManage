import { useMutation, useQueryClient } from '@tanstack/react-query';
import {
  createOrganizationUnit,
  createTenantScopeMapping,
  identityTenancyQueryKeys
} from '@/services/identityTenancy';

export const useOrganizationAction = (selectedUnitId?: string) => {
  const queryClient = useQueryClient();

  const createUnitMutation = useMutation({
    mutationFn: createOrganizationUnit,
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: identityTenancyQueryKeys.organizations() });
    }
  });

  const createMappingMutation = useMutation({
    mutationFn: (payload: Parameters<typeof createTenantScopeMapping>[1]) => {
      if (!selectedUnitId) {
        throw new Error('未选择组织单元');
      }
      return createTenantScopeMapping(selectedUnitId, payload);
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({
        queryKey: identityTenancyQueryKeys.mappings(selectedUnitId)
      });
    }
  });

  return { createUnitMutation, createMappingMutation };
};
