import { useMutation, useQueryClient } from '@tanstack/react-query';
import { message } from 'antd';
import { normalizeApiError } from '@/services/api/client';
import {
  createTemplateRelease,
  platformMarketplaceQueryKeys,
  syncCatalogSource,
  type CreateTemplateReleasePayload
} from '@/services/platformMarketplace';

export const useTemplateDistributionAction = () => {
  const queryClient = useQueryClient();

  const invalidateAll = () => {
    void queryClient.invalidateQueries({ queryKey: platformMarketplaceQueryKeys.all });
  };

  const syncSourceMutation = useMutation({
    mutationFn: (sourceId: string) => syncCatalogSource(sourceId),
    onSuccess: () => {
      message.success('目录同步请求已提交');
      invalidateAll();
    },
    onError: (error) => {
      message.error(normalizeApiError(error, '目录同步失败'));
    }
  });

  const createReleaseMutation = useMutation({
    mutationFn: (input: { templateId: string; payload: CreateTemplateReleasePayload }) =>
      createTemplateRelease(input.templateId, input.payload),
    onSuccess: () => {
      message.success('模板发布已提交');
      invalidateAll();
    },
    onError: (error) => {
      message.error(normalizeApiError(error, '模板发布失败'));
    }
  });

  return {
    syncSourceMutation,
    createReleaseMutation
  };
};
