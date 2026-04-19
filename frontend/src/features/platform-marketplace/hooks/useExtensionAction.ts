import { useMutation, useQueryClient } from '@tanstack/react-query';
import { message } from 'antd';
import { normalizeApiError } from '@/services/api/client';
import {
  createExtensionPackage,
  disableExtension,
  enableExtension,
  platformMarketplaceQueryKeys,
  type CreateExtensionPackagePayload
} from '@/services/platformMarketplace';

export const useExtensionAction = () => {
  const queryClient = useQueryClient();

  const invalidateAll = () => {
    void queryClient.invalidateQueries({ queryKey: platformMarketplaceQueryKeys.all });
  };

  const createExtensionMutation = useMutation({
    mutationFn: (payload: CreateExtensionPackagePayload) => createExtensionPackage(payload),
    onSuccess: () => {
      message.success('扩展已注册');
      invalidateAll();
    },
    onError: (error) => {
      message.error(normalizeApiError(error, '扩展注册失败'));
    }
  });

  const enableExtensionMutation = useMutation({
    mutationFn: (extensionId: string) => enableExtension(extensionId),
    onSuccess: () => {
      message.success('扩展启用请求已提交');
      invalidateAll();
    },
    onError: (error) => {
      message.error(normalizeApiError(error, '扩展启用失败'));
    }
  });

  const disableExtensionMutation = useMutation({
    mutationFn: (extensionId: string) => disableExtension(extensionId),
    onSuccess: () => {
      message.success('扩展停用请求已提交');
      invalidateAll();
    },
    onError: (error) => {
      message.error(normalizeApiError(error, '扩展停用失败'));
    }
  });

  return {
    createExtensionMutation,
    enableExtensionMutation,
    disableExtensionMutation
  };
};
