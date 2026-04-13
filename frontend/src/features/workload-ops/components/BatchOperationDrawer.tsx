import { useMutation } from '@tanstack/react-query';
import { Alert, Button, Drawer, Space, Typography } from 'antd';
import { ApiError, normalizeApiError } from '@/services/api/client';
import { submitBatchOperation } from '@/services/workloadOps';
import type { SubmitBatchOperationRequestDTO } from '@/services/api/types';

type BatchOperationDrawerProps = {
  open: boolean;
  payload: SubmitBatchOperationRequestDTO;
  onClose: () => void;
  onSubmitted?: (batchId: number) => void;
};

export const BatchOperationDrawer = ({ open, payload, onClose, onSubmitted }: BatchOperationDrawerProps) => {
  const mutation = useMutation({
    mutationFn: () => submitBatchOperation(payload),
    onSuccess: (data) => {
      if (data.id) {
        onSubmitted?.(data.id);
      }
    }
  });
  const permissionRevoked = mutation.error instanceof ApiError && (mutation.error.status === 401 || mutation.error.status === 403);

  return (
    <Drawer title="批量动作提交" open={open} onClose={onClose} width={520}>
      <Space direction="vertical" style={{ width: '100%' }}>
        <Typography.Text>动作：{payload.actionType}</Typography.Text>
        <Typography.Text>目标数量：{payload.targets.length}</Typography.Text>
        {mutation.error ? <Alert type={permissionRevoked ? 'warning' : 'error'} showIcon message={normalizeApiError(mutation.error)} /> : null}
        {mutation.data ? (
          <Alert type="success" showIcon message={`批量任务已提交（ID: ${mutation.data.id}）`} />
        ) : null}
        <Button type="primary" loading={mutation.isPending} disabled={permissionRevoked} onClick={() => mutation.mutate()}>
          提交批量任务
        </Button>
      </Space>
    </Drawer>
  );
};
