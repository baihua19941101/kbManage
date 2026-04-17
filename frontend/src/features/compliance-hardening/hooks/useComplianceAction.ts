import { useMutation, useQueryClient } from '@tanstack/react-query';
import { message } from 'antd';
import {
  createComplianceException,
  createRecheckTask,
  createRemediationTask,
  updateRemediationTask,
  reviewComplianceException,
  type CreateComplianceExceptionRequest,
  type CreateRecheckTaskRequest,
  type CreateRemediationTaskRequest,
  type ReviewComplianceExceptionRequest,
  type UpdateRemediationTaskRequest
} from '@/services/compliance';
import { normalizeErrorMessage } from '@/app/queryClient';

export const useComplianceAction = () => {
  const queryClient = useQueryClient();

  const invalidateGovernance = () => {
    void queryClient.invalidateQueries({ queryKey: ['compliance', 'findings'] });
    void queryClient.invalidateQueries({ queryKey: ['compliance', 'finding-detail'] });
    void queryClient.invalidateQueries({ queryKey: ['compliance', 'remediation'] });
    void queryClient.invalidateQueries({ queryKey: ['compliance', 'exceptions'] });
    void queryClient.invalidateQueries({ queryKey: ['compliance', 'rechecks'] });
    void queryClient.invalidateQueries({ queryKey: ['compliance', 'overview'] });
    void queryClient.invalidateQueries({ queryKey: ['compliance', 'trends'] });
  };

  const createRemediationMutation = useMutation({
    mutationFn: (input: { findingId: string; payload: CreateRemediationTaskRequest }) =>
      createRemediationTask(input.findingId, input.payload),
    onSuccess: () => {
      message.success('整改任务已创建');
      invalidateGovernance();
    },
    onError: (error) => {
      message.error(normalizeErrorMessage(error, '整改任务创建失败'));
    }
  });

  const updateRemediationMutation = useMutation({
    mutationFn: (input: { taskId: string; payload: UpdateRemediationTaskRequest }) =>
      updateRemediationTask(input.taskId, input.payload),
    onSuccess: () => {
      message.success('整改任务已更新');
      invalidateGovernance();
    },
    onError: (error) => {
      message.error(normalizeErrorMessage(error, '整改任务更新失败'));
    }
  });

  const createExceptionMutation = useMutation({
    mutationFn: (input: { findingId: string; payload: CreateComplianceExceptionRequest }) =>
      createComplianceException(input.findingId, input.payload),
    onSuccess: () => {
      message.success('例外申请已提交');
      invalidateGovernance();
    },
    onError: (error) => {
      message.error(normalizeErrorMessage(error, '例外申请提交失败'));
    }
  });

  const reviewExceptionMutation = useMutation({
    mutationFn: (input: { exceptionId: string; payload: ReviewComplianceExceptionRequest }) =>
      reviewComplianceException(input.exceptionId, input.payload),
    onSuccess: () => {
      message.success('例外审批已提交');
      invalidateGovernance();
    },
    onError: (error) => {
      message.error(normalizeErrorMessage(error, '例外审批失败'));
    }
  });

  const createRecheckMutation = useMutation({
    mutationFn: (input: { findingId: string; payload?: CreateRecheckTaskRequest }) =>
      createRecheckTask(input.findingId, input.payload),
    onSuccess: () => {
      message.success('复检任务已创建');
      invalidateGovernance();
    },
    onError: (error) => {
      message.error(normalizeErrorMessage(error, '复检任务创建失败'));
    }
  });

  return {
    createRemediationMutation,
    updateRemediationMutation,
    createExceptionMutation,
    reviewExceptionMutation,
    createRecheckMutation
  };
};
