import { useMutation, useQueryClient } from '@tanstack/react-query';
import { message } from 'antd';
import { normalizeErrorMessage } from '@/app/queryClient';
import {
  clusterLifecycleQueryKeys,
  createCluster,
  createClusterDriver,
  createClusterTemplate,
  createUpgradePlan,
  disableCluster,
  executeUpgradePlan,
  importCluster,
  registerCluster,
  retireCluster,
  scaleNodePool,
  validateClusterTemplate,
  type CreateClusterRequest,
  type CreateDriverRequest,
  type CreateTemplateRequest,
  type CreateUpgradePlanRequest,
  type DisableClusterRequest,
  type ImportClusterRequest,
  type RegisterClusterRequest,
  type RetireClusterRequest,
  type ScaleNodePoolRequest,
  type TemplateValidationRequest
} from '@/services/clusterLifecycle';

export const useLifecycleAction = () => {
  const queryClient = useQueryClient();

  const invalidateAll = () => {
    void queryClient.invalidateQueries({ queryKey: clusterLifecycleQueryKeys.all });
  };

  const importMutation = useMutation({
    mutationFn: (payload: ImportClusterRequest) => importCluster(payload),
    onSuccess: () => {
      message.success('导入请求已提交');
      invalidateAll();
    },
    onError: (error) => {
      message.error(normalizeErrorMessage(error, '导入请求提交失败'));
    }
  });

  const registerMutation = useMutation({
    mutationFn: (payload: RegisterClusterRequest) => registerCluster(payload),
    onSuccess: () => {
      message.success('注册指引已生成');
      invalidateAll();
    },
    onError: (error) => {
      message.error(normalizeErrorMessage(error, '注册指引生成失败'));
    }
  });

  const provisionMutation = useMutation({
    mutationFn: (payload: CreateClusterRequest) => createCluster(payload),
    onSuccess: () => {
      message.success('创建请求已提交');
      invalidateAll();
    },
    onError: (error) => {
      message.error(normalizeErrorMessage(error, '集群创建失败'));
    }
  });

  const upgradePlanMutation = useMutation({
    mutationFn: (input: { clusterId: string; payload: CreateUpgradePlanRequest }) =>
      createUpgradePlan(input.clusterId, input.payload),
    onSuccess: () => {
      message.success('升级计划已创建');
      invalidateAll();
    },
    onError: (error) => {
      message.error(normalizeErrorMessage(error, '升级计划创建失败'));
    }
  });

  const executeUpgradeMutation = useMutation({
    mutationFn: (input: { clusterId: string; planId: string }) =>
      executeUpgradePlan(input.clusterId, input.planId),
    onSuccess: () => {
      message.success('升级执行已受理');
      invalidateAll();
    },
    onError: (error) => {
      message.error(normalizeErrorMessage(error, '升级执行失败'));
    }
  });

  const scaleNodePoolMutation = useMutation({
    mutationFn: (input: { clusterId: string; nodePoolId: string; payload: ScaleNodePoolRequest }) =>
      scaleNodePool(input.clusterId, input.nodePoolId, input.payload),
    onSuccess: () => {
      message.success('节点池扩缩请求已提交');
      invalidateAll();
    },
    onError: (error) => {
      message.error(normalizeErrorMessage(error, '节点池扩缩失败'));
    }
  });

  const disableMutation = useMutation({
    mutationFn: (input: { clusterId: string; payload: DisableClusterRequest }) =>
      disableCluster(input.clusterId, input.payload),
    onSuccess: () => {
      message.success('停用流程已发起');
      invalidateAll();
    },
    onError: (error) => {
      message.error(normalizeErrorMessage(error, '停用流程发起失败'));
    }
  });

  const retireMutation = useMutation({
    mutationFn: (input: { clusterId: string; payload: RetireClusterRequest }) =>
      retireCluster(input.clusterId, input.payload),
    onSuccess: () => {
      message.success('退役流程已发起');
      invalidateAll();
    },
    onError: (error) => {
      message.error(normalizeErrorMessage(error, '退役流程发起失败'));
    }
  });

  const createDriverMutation = useMutation({
    mutationFn: (payload: CreateDriverRequest) => createClusterDriver(payload),
    onSuccess: () => {
      message.success('驱动版本已保存');
      invalidateAll();
    },
    onError: (error) => {
      message.error(normalizeErrorMessage(error, '驱动版本保存失败'));
    }
  });

  const createTemplateMutation = useMutation({
    mutationFn: (payload: CreateTemplateRequest) => createClusterTemplate(payload),
    onSuccess: () => {
      message.success('模板已保存');
      invalidateAll();
    },
    onError: (error) => {
      message.error(normalizeErrorMessage(error, '模板保存失败'));
    }
  });

  const validateTemplateMutation = useMutation({
    mutationFn: (input: { templateId: string; payload: TemplateValidationRequest }) =>
      validateClusterTemplate(input.templateId, input.payload),
    onError: (error) => {
      message.error(normalizeErrorMessage(error, '模板兼容性校验失败'));
    }
  });

  return {
    importMutation,
    registerMutation,
    provisionMutation,
    upgradePlanMutation,
    executeUpgradeMutation,
    scaleNodePoolMutation,
    disableMutation,
    retireMutation,
    createDriverMutation,
    createTemplateMutation,
    validateTemplateMutation
  };
};
