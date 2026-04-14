import { useEffect, useRef, useState } from 'react';
import {
  Alert,
  Button,
  Drawer,
  Form,
  Input,
  InputNumber,
  Progress,
  Select,
  Space,
  Tag,
  Typography
} from 'antd';
import {
  canPromoteGitOps,
  canRollbackGitOps,
  canSyncGitOps,
  useAuthStore
} from '@/features/auth/store';
import { useDeliveryOperation } from '@/features/gitops/hooks/useDeliveryOperation';
import { isAuthorizationError, normalizeApiError } from '@/services/api/client';
import type { GitOpsActionRequestDTO, GitOpsOperationDTO } from '@/services/api/types';
import type { ResourceId } from '@/services/gitops';

type ReleaseActionDrawerProps = {
  open: boolean;
  unitId: ResourceId;
  unitName?: string;
  onClose: () => void;
  onOperationChange?: (operation: GitOpsOperationDTO) => void;
};

type ReleaseActionForm = {
  actionType: GitOpsActionRequestDTO['actionType'];
  environment?: string;
  targetReleaseId?: number;
  targetAppVersion?: string;
  targetConfigVersion?: string;
  reason?: string;
};

const ACTION_OPTIONS: Array<{ value: GitOpsActionRequestDTO['actionType']; label: string }> = [
  { value: 'install', label: 'install' },
  { value: 'sync', label: 'sync' },
  { value: 'resync', label: 'resync' },
  { value: 'upgrade', label: 'upgrade' },
  { value: 'promote', label: 'promote' },
  { value: 'rollback', label: 'rollback' },
  { value: 'pause', label: 'pause' },
  { value: 'resume', label: 'resume' },
  { value: 'uninstall', label: 'uninstall' }
];

const statusColorMap: Record<string, string> = {
  pending: 'gold',
  running: 'blue',
  succeeded: 'green',
  partially_succeeded: 'orange',
  failed: 'red',
  canceled: 'default'
};

const canSubmitAction = (
  actionType: GitOpsActionRequestDTO['actionType'],
  user: ReturnType<typeof useAuthStore.getState>['user']
) => {
  if (actionType === 'promote') {
    return canPromoteGitOps(user);
  }
  if (actionType === 'rollback') {
    return canRollbackGitOps(user);
  }
  return canSyncGitOps(user);
};

const isVersionAction = (actionType: GitOpsActionRequestDTO['actionType']) =>
  actionType === 'install' || actionType === 'upgrade';

const isRollbackAction = (actionType: GitOpsActionRequestDTO['actionType']) =>
  actionType === 'rollback';

export const ReleaseActionDrawer = ({
  open,
  unitId,
  unitName,
  onClose,
  onOperationChange
}: ReleaseActionDrawerProps) => {
  const [form] = Form.useForm<ReleaseActionForm>();
  const user = useAuthStore((state) => state.user);

  const { submit, reset, operation, isSubmitting, isPolling, error } = useDeliveryOperation({
    enabled: open,
    onOperationChange
  });
  const resetRef = useRef(reset);
  const [permissionRevoked, setPermissionRevoked] = useState(false);

  const actionType = Form.useWatch('actionType', form) || 'sync';
  const actionAllowed = canSubmitAction(actionType, user) && !permissionRevoked;

  useEffect(() => {
    resetRef.current = reset;
  }, [reset]);

  useEffect(() => {
    if (error && isAuthorizationError(error)) {
      setPermissionRevoked(true);
    }
  }, [error]);

  useEffect(() => {
    if (!open) {
      form.resetFields();
      resetRef.current();
      setPermissionRevoked(false);
      return;
    }

    form.setFieldsValue({ actionType: 'sync' });
  }, [form, open]);

  const handleSubmit = async () => {
    const values = await form.validateFields();
    const payload: GitOpsActionRequestDTO = {
      actionType: values.actionType,
      environment: values.environment || undefined,
      targetReleaseId: values.targetReleaseId || undefined,
      targetAppVersion: values.targetAppVersion || undefined,
      targetConfigVersion: values.targetConfigVersion || undefined,
      reason: values.reason || undefined
    };

    await submit({ unitId, payload });
  };

  const status = operation?.status || 'pending';
  const progressPercent = operation?.progressPercent ?? (status === 'running' ? 30 : 0);

  return (
    <Drawer
      title="发布动作"
      width={500}
      open={open}
      onClose={onClose}
      destroyOnHidden
      extra={
        <Space>
          <Button onClick={onClose}>关闭</Button>
          <Button
            type="primary"
            onClick={() => void handleSubmit()}
            loading={isSubmitting}
            disabled={!actionAllowed}
          >
            提交动作
          </Button>
        </Space>
      }
    >
      <Space direction="vertical" size={12} style={{ width: '100%' }}>
        <Typography.Text type="secondary">交付单元：{unitName || unitId}</Typography.Text>

        {!actionAllowed ? (
          <Alert
            type="info"
            showIcon
            message={permissionRevoked ? '权限已回收' : '当前动作未授权'}
            description={
              permissionRevoked
                ? '检测到后端权限拒绝，当前动作入口已锁定。'
                : '你当前没有该动作权限，请联系管理员。'
            }
          />
        ) : null}

        <Form form={form} layout="vertical" initialValues={{ actionType: 'sync' }}>
          <Form.Item label="动作类型" name="actionType" rules={[{ required: true, message: '请选择动作类型' }]}>
            <Select options={ACTION_OPTIONS} />
          </Form.Item>
          <Form.Item label="环境" name="environment">
            <Input placeholder="例如：dev / staging / prod" />
          </Form.Item>

          {isRollbackAction(actionType) ? (
            <Form.Item
              label="目标发布版本 ID"
              name="targetReleaseId"
              rules={[{ required: true, message: '回滚动作必须填写目标发布版本 ID' }]}
            >
              <InputNumber style={{ width: '100%' }} min={1} placeholder="例如：42" />
            </Form.Item>
          ) : null}

          {isVersionAction(actionType) ? (
            <>
              <Form.Item label="目标应用版本" name="targetAppVersion">
                <Input placeholder="例如：1.2.3" />
              </Form.Item>
              <Form.Item label="目标配置版本" name="targetConfigVersion">
                <Input placeholder="例如：2026.04.13" />
              </Form.Item>
            </>
          ) : null}

          <Form.Item label="原因说明" name="reason">
            <Input.TextArea rows={3} placeholder="可选，建议填写动作原因" />
          </Form.Item>
        </Form>

        {operation ? (
          <>
            <Space wrap>
              <Tag color={statusColorMap[status] || 'default'}>状态：{status}</Tag>
              <Tag>动作：{operation.operationType || operation.actionType || actionType}</Tag>
              {isPolling ? <Tag color="processing">轮询中</Tag> : null}
            </Space>
            {typeof operation.progressPercent === 'number' || status === 'running' ? (
              <Progress percent={progressPercent} size="small" status={status === 'failed' ? 'exception' : 'active'} />
            ) : null}
          </>
        ) : null}

        {operation?.resultMessage ? <Alert type="success" showIcon message={operation.resultMessage} /> : null}
        {operation?.failureReason ? <Alert type="error" showIcon message={operation.failureReason} /> : null}
        {error ? (
          <Alert
            type="error"
            showIcon
            message="动作提交失败"
            description={normalizeApiError(error, '动作提交失败，请稍后重试。')}
          />
        ) : null}
      </Space>
    </Drawer>
  );
};
