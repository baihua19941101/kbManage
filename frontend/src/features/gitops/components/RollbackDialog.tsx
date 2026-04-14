import { useEffect, useRef, useState } from 'react';
import { Alert, Button, Form, Input, Modal, Space, Tag, Typography } from 'antd';
import { canRollbackGitOps, useAuthStore } from '@/features/auth/store';
import { useDeliveryOperation } from '@/features/gitops/hooks/useDeliveryOperation';
import { isAuthorizationError, normalizeApiError } from '@/services/api/client';
import type { GitOpsOperationDTO } from '@/services/api/types';
import type { GitOpsReleaseRevision, ResourceId } from '@/services/gitops';

type RollbackDialogProps = {
  open: boolean;
  unitId: ResourceId;
  revision?: GitOpsReleaseRevision;
  onClose: () => void;
  onOperationChange?: (operation: GitOpsOperationDTO) => void;
};

type RollbackFormData = {
  environment?: string;
  reason?: string;
};

const statusColorMap: Record<string, string> = {
  pending: 'gold',
  running: 'blue',
  succeeded: 'green',
  partially_succeeded: 'orange',
  failed: 'red',
  canceled: 'default'
};

export const RollbackDialog = ({
  open,
  unitId,
  revision,
  onClose,
  onOperationChange
}: RollbackDialogProps) => {
  const [form] = Form.useForm<RollbackFormData>();
  const user = useAuthStore((state) => state.user);
  const canRollback = canRollbackGitOps(user);

  const { submit, reset, operation, isSubmitting, isPolling, error } = useDeliveryOperation({
    enabled: open,
    onOperationChange
  });
  const resetRef = useRef(reset);
  const [permissionRevoked, setPermissionRevoked] = useState(false);

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

    form.setFieldsValue({ reason: undefined, environment: undefined });
  }, [form, open]);

  const handleSubmit = async () => {
    if (!revision?.id) {
      return;
    }

    const values = await form.validateFields();
    await submit({
      unitId,
      payload: {
        actionType: 'rollback',
        targetReleaseId: Number(revision.id),
        environment: values.environment || undefined,
        reason: values.reason || undefined
      }
    });
  };

  const status = operation?.status;

  return (
    <Modal
      title="回滚发布"
      open={open}
      onCancel={onClose}
      footer={
        <Space>
          <Button onClick={onClose}>关闭</Button>
          <Button
            type="primary"
            onClick={() => void handleSubmit()}
            loading={isSubmitting}
            disabled={!revision || !canRollback || permissionRevoked}
          >
            确认回滚
          </Button>
        </Space>
      }
      destroyOnHidden
    >
      <Space direction="vertical" style={{ width: '100%' }}>
        {!canRollback || permissionRevoked ? (
          <Alert
            type="info"
            showIcon
            message={permissionRevoked ? '权限已回收' : '当前账号没有回滚权限'}
            description={
              permissionRevoked
                ? '检测到后端权限拒绝，当前回滚入口已锁定。'
                : '你可以查看发布历史，但不能提交回滚。'
            }
          />
        ) : null}

        <Typography.Text>
          回滚目标：发布 ID <strong>{revision?.id || '-'}</strong>
        </Typography.Text>
        <Typography.Text type="secondary">
          版本：{revision?.appVersion || '-'} / 配置：{revision?.configVersion || '-'}
        </Typography.Text>

        <Form form={form} layout="vertical">
          <Form.Item label="环境（可选）" name="environment">
            <Input placeholder="例如：prod" />
          </Form.Item>
          <Form.Item label="回滚原因（可选）" name="reason">
            <Input.TextArea rows={3} placeholder="建议记录回滚触发原因" />
          </Form.Item>
        </Form>

        {status ? (
          <Space wrap>
            <Tag color={statusColorMap[status] || 'default'}>状态：{status}</Tag>
            {isPolling ? <Tag color="processing">轮询中</Tag> : null}
          </Space>
        ) : null}

        {operation?.resultMessage ? <Alert type="success" showIcon message={operation.resultMessage} /> : null}
        {operation?.failureReason ? <Alert type="error" showIcon message={operation.failureReason} /> : null}
        {error ? (
          <Alert
            type="error"
            showIcon
            message="回滚提交失败"
            description={normalizeApiError(error, '回滚提交失败，请稍后重试。')}
          />
        ) : null}
      </Space>
    </Modal>
  );
};
