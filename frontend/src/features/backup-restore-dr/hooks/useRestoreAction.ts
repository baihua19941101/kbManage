import { useMutation, useQueryClient } from '@tanstack/react-query';
import {
  backupRestoreQueryKeys,
  createMigrationPlan,
  createRestoreJob,
  runBackupPolicy,
  validateRestoreJob
} from '@/services/backupRestore';

export const useRestoreAction = () => {
  const queryClient = useQueryClient();

  const runBackupMutation = useMutation({
    mutationFn: runBackupPolicy,
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: backupRestoreQueryKeys.restorePoints() });
    }
  });

  const restoreMutation = useMutation({
    mutationFn: createRestoreJob,
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: backupRestoreQueryKeys.restoreJobs() });
    }
  });

  const precheckMutation = useMutation({
    mutationFn: validateRestoreJob
  });

  const migrationMutation = useMutation({
    mutationFn: createMigrationPlan
  });

  return {
    runBackupMutation,
    restoreMutation,
    precheckMutation,
    migrationMutation
  };
};
