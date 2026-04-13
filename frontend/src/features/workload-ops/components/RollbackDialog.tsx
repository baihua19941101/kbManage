import { useEffect, useState } from 'react';
import { useMutation } from '@tanstack/react-query';
import { Alert, Button, Modal, Space, Typography } from 'antd';
import { canRollbackWorkloadOps, useAuthStore } from '@/features/auth/store';
import { ApiError, normalizeApiError } from '@/services/api/client';
import { submitWorkloadAction } from '@/services/workloadOps';
import type { ReleaseRevisionDTO } from '@/services/api/types';

type RollbackDialogProps = {
  open: boolean;
  clusterId: number;
  namespace: string;
  resourceKind: 'Deployment' | 'StatefulSet' | 'DaemonSet';
  resourceName: string;
  revision?: ReleaseRevisionDTO;
  onClose: () => void;
};

export const RollbackDialog = ({
  open,
  clusterId,
  namespace,
  resourceKind,
  resourceName,
  revision,
  onClose
}: RollbackDialogProps) => {
  const user = useAuthStore((state) => state.user);
  const canRollback = canRollbackWorkloadOps(user);
  const [permissionMessage, setPermissionMessage] = useState<string>();

  const mutation = useMutation({
    mutationFn: () =>
      submitWorkloadAction({
        clusterId,
        namespace,
        resourceKind,
        resourceName,
        actionType: 'rollback',
        riskConfirmed: true,
        payload: { revision: revision?.revision }
      }),
    onSuccess: () => {
      setPermissionMessage(undefined);
    },
    onError: (error) => {
      if (error instanceof ApiError && (error.status === 401 || error.status === 403)) {
        setPermissionMessage(normalizeApiError(error, '权限已回收，无法继续执行回滚。'));
      }
    }
  });
  const permissionRevoked =
    mutation.error instanceof ApiError && (mutation.error.status === 401 || mutation.error.status === 403);
  const actionDisabled = !canRollback || Boolean(permissionMessage) || permissionRevoked;

  useEffect(() => {
    if (!open) {
      return;
    }
    setPermissionMessage(undefined);
  }, [open]);

  return (
    <Modal open={open} onCancel={onClose} footer={null} title="回滚确认">
      <Space direction="vertical" style={{ width: '100%' }}>
        {!canRollback ? (
          <Alert
            type="info"
            showIcon
            message="当前为只读模式"
            description="你没有回滚权限，可查看发布历史但不能提交回滚。"
          />
        ) : null}
        {permissionMessage ? (
          <Alert type="warning" showIcon message="权限已回收" description={permissionMessage} />
        ) : null}
        <Typography.Text>
          目标版本：<strong>{revision?.revision ?? '-'}</strong>
        </Typography.Text>
        {mutation.error ? (
          <Alert
            type={permissionRevoked ? 'warning' : 'error'}
            showIcon
            message={normalizeApiError(mutation.error)}
          />
        ) : null}
        {mutation.data ? (
          <Alert type="success" showIcon message={mutation.data.resultMessage || '回滚已提交'} />
        ) : null}
        <Button type="primary" loading={mutation.isPending} disabled={!revision || actionDisabled} onClick={() => mutation.mutate()}>
          确认回滚
        </Button>
      </Space>
    </Modal>
  );
};
