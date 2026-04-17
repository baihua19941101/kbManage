import { useMutation, useQueryClient } from '@tanstack/react-query';
import { message } from 'antd';
import {
  createComplianceArchiveExport,
  type CreateArchiveExportTaskRequest
} from '@/services/compliance';
import { normalizeErrorMessage } from '@/app/queryClient';

export const useArchiveExport = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (payload: CreateArchiveExportTaskRequest) => createComplianceArchiveExport(payload),
    onSuccess: () => {
      message.success('归档导出任务已创建');
      void queryClient.invalidateQueries({ queryKey: ['compliance', 'archive-exports'] });
    },
    onError: (error) => {
      message.error(normalizeErrorMessage(error, '归档导出任务创建失败'));
    }
  });
};
