import { useMutation, useQueryClient } from '@tanstack/react-query';
import { message } from 'antd';
import { normalizeApiError } from '@/services/api/client';
import { createHAPolicy, createMaintenanceWindow } from '@/services/sreScale';

export const useScaleEvidenceAction = () => {
  const queryClient = useQueryClient();
  const invalidateAll = () => void queryClient.invalidateQueries({ queryKey: ['sreScale'] });

  return {
    createHAPolicyMutation: useMutation({
      mutationFn: (payload: Record<string, unknown>) => createHAPolicy(payload),
      onSuccess: () => {
        message.success('高可用策略已创建');
        invalidateAll();
      },
      onError: (error) => message.error(normalizeApiError(error, '高可用策略创建失败'))
    }),
    createMaintenanceWindowMutation: useMutation({
      mutationFn: (payload: Record<string, unknown>) => createMaintenanceWindow(payload),
      onSuccess: () => {
        message.success('维护窗口已创建');
        invalidateAll();
      },
      onError: (error) => message.error(normalizeApiError(error, '维护窗口创建失败'))
    })
  };
};
