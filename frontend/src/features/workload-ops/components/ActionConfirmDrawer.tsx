import { useMutation } from '@tanstack/react-query';
import { Alert, Button, Drawer, Space, Typography } from 'antd';
import { ApiError, normalizeApiError } from '@/services/api/client';
import { submitWorkloadAction } from '@/services/workloadOps';
import type { SubmitWorkloadActionRequestDTO } from '@/services/api/types';

type ActionConfirmDrawerProps = {
  open: boolean;
  payload: SubmitWorkloadActionRequestDTO;
  onClose: () => void;
};

export const ActionConfirmDrawer = ({ open, payload, onClose }: ActionConfirmDrawerProps) => {
  const mutation = useMutation({
    mutationFn: () => submitWorkloadAction(payload)
  });
  const permissionRevoked = mutation.error instanceof ApiError && (mutation.error.status === 401 || mutation.error.status === 403);

  return (
    <Drawer title="动作确认" open={open} onClose={onClose} width={480}>
      <Space direction="vertical" style={{ width: '100%' }}>
        <Typography.Text>
          动作类型：<strong>{payload.actionType}</strong>
        </Typography.Text>
        <Typography.Text>
          目标：{payload.namespace}/{payload.resourceName}
        </Typography.Text>
        {mutation.error ? (
          <Alert type={permissionRevoked ? 'warning' : 'error'} showIcon message={`动作提交失败：${normalizeApiError(mutation.error)}`} />
        ) : null}
        {mutation.data ? (
          <Alert type="success" showIcon message={`动作已提交，状态：${mutation.data.status}`} />
        ) : null}
        <Button type="primary" loading={mutation.isPending} disabled={permissionRevoked} onClick={() => mutation.mutate()}>
          确认提交
        </Button>
      </Space>
    </Drawer>
  );
};
