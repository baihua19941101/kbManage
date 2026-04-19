import { useMutation, useQueryClient } from '@tanstack/react-query';
import {
  createDelegationGrant,
  createRoleAssignment,
  createRoleDefinition,
  identityTenancyQueryKeys
} from '@/services/identityTenancy';

export const useRoleGovernanceAction = () => {
  const queryClient = useQueryClient();

  const createRoleMutation = useMutation({
    mutationFn: createRoleDefinition,
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: identityTenancyQueryKeys.roles() });
    }
  });

  const createAssignmentMutation = useMutation({
    mutationFn: createRoleAssignment,
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: identityTenancyQueryKeys.assignments() });
    }
  });

  const createDelegationMutation = useMutation({
    mutationFn: createDelegationGrant,
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: identityTenancyQueryKeys.delegations() });
    }
  });

  return { createRoleMutation, createAssignmentMutation, createDelegationMutation };
};
