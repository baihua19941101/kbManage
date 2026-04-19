import { useMutation, useQueryClient } from '@tanstack/react-query';
import {
  backupRestoreQueryKeys,
  createDrillPlan,
  createDrillReport,
  runDrillPlan
} from '@/services/backupRestore';

export const useDrillAction = () => {
  const queryClient = useQueryClient();

  const createPlanMutation = useMutation({
    mutationFn: createDrillPlan,
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: backupRestoreQueryKeys.drillPlans() });
    }
  });

  const runDrillMutation = useMutation({
    mutationFn: runDrillPlan
  });

  const reportMutation = useMutation({
    mutationFn: createDrillReport
  });

  return {
    createPlanMutation,
    runDrillMutation,
    reportMutation
  };
};
