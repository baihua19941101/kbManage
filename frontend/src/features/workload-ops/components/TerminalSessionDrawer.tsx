import { useEffect, useState } from 'react';
import { useMutation } from '@tanstack/react-query';
import { Alert, Button, Descriptions, Drawer, Space, Typography } from 'antd';
import { canAccessWorkloadOpsTerminal, useAuthStore } from '@/features/auth/store';
import { ApiError, normalizeApiError } from '@/services/api/client';
import { closeTerminalSession, createTerminalSession } from '@/services/workloadOps';
import type { WorkloadInstanceDTO } from '@/services/api/types';

type TerminalSessionDrawerProps = {
  open: boolean;
  clusterId: number;
  namespace: string;
  workloadKind: 'Deployment' | 'StatefulSet' | 'DaemonSet';
  workloadName: string;
  target?: WorkloadInstanceDTO;
  onClose: () => void;
};

export const TerminalSessionDrawer = ({
  open,
  clusterId,
  namespace,
  workloadKind,
  workloadName,
  target,
  onClose
}: TerminalSessionDrawerProps) => {
  const user = useAuthStore((state) => state.user);
  const canUseTerminal = canAccessWorkloadOpsTerminal(user);
  const [permissionMessage, setPermissionMessage] = useState<string>();

  const createMutation = useMutation({
    mutationFn: () =>
      createTerminalSession({
        clusterId,
        namespace,
        podName: target?.podName,
        containerName: target?.containerName,
        workloadKind,
        workloadName
      }),
    onSuccess: () => {
      setPermissionMessage(undefined);
    },
    onError: (error) => {
      if (error instanceof ApiError && (error.status === 401 || error.status === 403)) {
        setPermissionMessage(normalizeApiError(error, '权限已回收，无法继续创建终端会话。'));
      }
    }
  });

  const closeMutation = useMutation({
    mutationFn: (sessionId: number) => closeTerminalSession(sessionId),
    onSuccess: () => {
      setPermissionMessage(undefined);
      onClose();
    },
    onError: (error) => {
      if (error instanceof ApiError && (error.status === 401 || error.status === 403)) {
        setPermissionMessage(normalizeApiError(error, '权限已回收，无法继续关闭终端会话。'));
      }
    }
  });

  useEffect(() => {
    if (!open) {
      return;
    }
    setPermissionMessage(undefined);
  }, [open]);

  const created = createMutation.data;
  const permissionRevoked =
    createMutation.error instanceof ApiError &&
    (createMutation.error.status === 401 || createMutation.error.status === 403);
  const actionDisabled = !canUseTerminal || Boolean(permissionMessage) || permissionRevoked;

  return (
    <Drawer title="容器终端会话" open={open} onClose={onClose} width={520} destroyOnHidden>
      <Space direction="vertical" style={{ width: '100%' }}>
        {!canUseTerminal ? (
          <Alert
            type="info"
            showIcon
            message="当前为只读模式"
            description="你没有容器终端权限，可查看实例信息但不能创建或关闭终端会话。"
          />
        ) : null}
        {permissionMessage ? (
          <Alert type="warning" showIcon message="权限已回收" description={permissionMessage} />
        ) : null}
        {target ? (
          <Descriptions column={1} bordered size="small">
            <Descriptions.Item label="Pod">{target.podName}</Descriptions.Item>
            <Descriptions.Item label="Container">{target.containerName ?? '-'}</Descriptions.Item>
            <Descriptions.Item label="Phase">{target.phase}</Descriptions.Item>
          </Descriptions>
        ) : (
          <Alert type="info" showIcon message="请选择一个实例后再创建终端会话" />
        )}

        {createMutation.error ? (
          <Alert
            type={permissionRevoked ? 'warning' : 'error'}
            showIcon
            message={`创建会话失败：${normalizeApiError(createMutation.error)}`}
          />
        ) : null}
        {created ? (
          <Alert
            type="success"
            showIcon
            message={`会话已创建（ID: ${created.id}）`}
            description={`状态：${created.status}`}
          />
        ) : null}

        <Space>
          <Button
            type="primary"
            disabled={!target || actionDisabled}
            loading={createMutation.isPending}
            onClick={() => createMutation.mutate()}
          >
            创建会话
          </Button>
          <Button
            disabled={!created?.id || actionDisabled}
            loading={closeMutation.isPending}
            onClick={() => {
              if (created?.id) {
                closeMutation.mutate(created.id);
              }
            }}
          >
            关闭会话
          </Button>
        </Space>

        <Typography.Paragraph type="secondary" style={{ marginBottom: 0 }}>
          首期仅提供受控会话管理与状态反馈，不展示完整终端流。
        </Typography.Paragraph>
      </Space>
    </Drawer>
  );
};
