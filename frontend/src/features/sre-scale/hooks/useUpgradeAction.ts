import { useMutation, useQueryClient } from '@tanstack/react-query';
import { message } from 'antd';
import { normalizeApiError } from '@/services/api/client';
import {
  createRollbackValidation,
  createUpgradePlan,
  runUpgradePrecheck
} from '@/services/sreScale';

export const useUpgradeAction = () => {
  const queryClient = useQueryClient();
  const invalidateAll = () => void queryClient.invalidateQueries({ queryKey: ['sreScale'] });

  return {
    precheckMutation: useMutation({
      mutationFn: (payload: Record<string, unknown>) => runUpgradePrecheck(payload),
      onSuccess: () => {
        message.success('升级前检查已完成');
        invalidateAll();
      },
      onError: (error) => message.error(normalizeApiError(error, '升级前检查失败'))
    }),
    createUpgradeMutation: useMutation({
      mutationFn: (payload: Record<string, unknown>) => createUpgradePlan(payload),
      onSuccess: () => {
        message.success('升级计划已创建');
        invalidateAll();
      },
      onError: (error) => message.error(normalizeApiError(error, '升级计划创建失败'))
    }),
    rollbackMutation: useMutation({
      mutationFn: (input: { upgradeId: string; payload: Record<string, unknown> }) =>
        createRollbackValidation(input.upgradeId, input.payload),
      onSuccess: () => {
        message.success('回退验证已登记');
        invalidateAll();
      },
      onError: (error) => message.error(normalizeApiError(error, '回退验证提交失败'))
    })
  };
};
