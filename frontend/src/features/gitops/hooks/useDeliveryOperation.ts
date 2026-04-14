import { useCallback, useEffect, useRef, useState } from 'react';
import { useMutation, useQuery } from '@tanstack/react-query';
import { queryKeys } from '@/app/queryClient';
import type { GitOpsActionRequestDTO, GitOpsOperationDTO } from '@/services/api/types';
import { getGitOpsOperation, submitGitOpsAction, type ResourceId } from '@/services/gitops';

const RUNNING_STATUSES = new Set(['pending', 'running']);
const TERMINAL_STATUSES = new Set(['partially_succeeded', 'succeeded', 'failed', 'canceled']);
const DEFAULT_POLL_INTERVAL_MS = 2000;

const normalizeStatus = (status?: string) => (status || '').toLowerCase();

const buildOperationSnapshotKey = (operation: GitOpsOperationDTO) => {
  const stageSignature = (operation.stages || [])
    .map((stage) => `${stage.environment || 'unknown'}:${stage.status || 'unknown'}`)
    .join('|');

  return [
    operation.id,
    normalizeStatus(operation.status),
    operation.progressPercent ?? 'na',
    operation.resultMessage || operation.resultSummary || '',
    operation.failureReason || '',
    stageSignature
  ].join(':');
};

export type SubmitDeliveryOperationInput = {
  unitId: ResourceId;
  payload: GitOpsActionRequestDTO;
};

type UseDeliveryOperationOptions = {
  enabled?: boolean;
  pollIntervalMs?: number;
  onOperationChange?: (operation: GitOpsOperationDTO) => void;
};

export const useDeliveryOperation = (options: UseDeliveryOperationOptions = {}) => {
  const {
    enabled = true,
    pollIntervalMs = DEFAULT_POLL_INTERVAL_MS,
    onOperationChange
  } = options;
  const [operationId, setOperationId] = useState<ResourceId>();
  const [submittedOperation, setSubmittedOperation] = useState<GitOpsOperationDTO>();
  const lastNotifiedRef = useRef<string>('');

  const submitMutation = useMutation({
    mutationFn: ({ unitId, payload }: SubmitDeliveryOperationInput) => submitGitOpsAction(unitId, payload),
    onSuccess: (operation) => {
      setSubmittedOperation(operation);
      setOperationId(operation.id);
    }
  });

  const operationQuery = useQuery({
    queryKey: queryKeys.gitops.operation(operationId),
    enabled: enabled && operationId !== undefined,
    queryFn: () => getGitOpsOperation(operationId as ResourceId),
    refetchInterval: (query) => {
      const data = query.state.data as GitOpsOperationDTO | undefined;
      const status = normalizeStatus(data?.status || submittedOperation?.status);
      return RUNNING_STATUSES.has(status) ? pollIntervalMs : false;
    },
    meta: { suppressGlobalError: true }
  });

  const operation = operationQuery.data || submittedOperation;

  useEffect(() => {
    if (!operation) {
      return;
    }

    const notifyKey = buildOperationSnapshotKey(operation);
    if (notifyKey !== lastNotifiedRef.current) {
      lastNotifiedRef.current = notifyKey;
      onOperationChange?.(operation);
    }

    const status = normalizeStatus(operation.status);
    if (TERMINAL_STATUSES.has(status)) {
      setOperationId(undefined);
    }
  }, [operation, onOperationChange]);

  const submit = useCallback(
    (input: SubmitDeliveryOperationInput) => submitMutation.mutateAsync(input),
    [submitMutation]
  );

  const reset = useCallback(() => {
    setOperationId(undefined);
    setSubmittedOperation(undefined);
    lastNotifiedRef.current = '';
    submitMutation.reset();
  }, [submitMutation]);

  const status = normalizeStatus(operation?.status);

  return {
    submit,
    reset,
    operation,
    isSubmitting: submitMutation.isPending,
    isPolling: operationId !== undefined && operationQuery.isFetching,
    isRunning: RUNNING_STATUSES.has(status),
    isTerminal: TERMINAL_STATUSES.has(status),
    error: submitMutation.error || operationQuery.error || null
  };
};
