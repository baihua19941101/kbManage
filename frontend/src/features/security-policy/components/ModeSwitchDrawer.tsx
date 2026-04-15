import { useEffect } from 'react';
import { useMutation } from '@tanstack/react-query';
import { Alert, Button, Drawer, Form, Input, Select, Space, message } from 'antd';
import { normalizeApiError } from '@/services/api/client';
import { switchPolicyMode } from '@/services/securityPolicy';
import type {
  SecurityPolicyDTO,
  SecurityPolicyEnforcementMode,
  SwitchPolicyModeRequestDTO
} from '@/services/api/types';

type ModeSwitchDrawerProps = {
  open: boolean;
  policy?: SecurityPolicyDTO;
  readonly?: boolean;
  onClose: () => void;
  onSuccess?: () => void;
};

type FormValues = {
  targetMode: SecurityPolicyEnforcementMode;
  assignmentIds?: string;
  reason?: string;
};

const splitCommaText = (value?: string): string[] | undefined => {
  const items =
    value
      ?.split(',')
      .map((item) => item.trim())
      .filter((item) => item.length > 0) ?? [];
  return items.length > 0 ? items : undefined;
};

export const ModeSwitchDrawer = ({
  open,
  policy,
  readonly,
  onClose,
  onSuccess
}: ModeSwitchDrawerProps) => {
  const [form] = Form.useForm<FormValues>();

  useEffect(() => {
    if (!open) {
      return;
    }

    form.setFieldsValue({
      targetMode: policy?.defaultEnforcementMode ?? 'warn'
    });
  }, [form, open, policy?.defaultEnforcementMode]);

  const mutation = useMutation({
    mutationFn: async (values: FormValues) => {
      if (!policy?.id) {
        throw new Error('缺少策略 ID，无法切换模式');
      }

      const payload: SwitchPolicyModeRequestDTO = {
        targetMode: values.targetMode,
        assignmentIds: splitCommaText(values.assignmentIds),
        reason: values.reason?.trim() || undefined
      };

      return switchPolicyMode(policy.id, payload);
    },
    onSuccess: () => {
      message.success('模式切换任务已提交');
      form.resetFields();
      onSuccess?.();
      onClose();
    }
  });

  const handleClose = () => {
    if (mutation.isPending) {
      return;
    }
    form.resetFields();
    onClose();
  };

  return (
    <Drawer
      title={policy ? `模式切换 - ${policy.name}` : '模式切换'}
      open={open}
      width={560}
      destroyOnClose
      getContainer={false}
      onClose={handleClose}
      footer={
        <Space>
          <Button onClick={handleClose}>取消</Button>
          <Button
            type="primary"
            disabled={readonly || !policy}
            loading={mutation.isPending}
            onClick={() => form.submit()}
          >
            提交切换
          </Button>
        </Space>
      }
    >
      {!policy ? (
        <Alert type="info" showIcon message="请先在页面中选择一个策略" />
      ) : (
        <Space direction="vertical" size={12} style={{ width: '100%' }}>
          {readonly ? (
            <Alert
              type="info"
              showIcon
              message="当前为只读模式"
              description="你可以查看命中和例外状态，但无法提交模式切换。"
            />
          ) : null}

          {mutation.error ? (
            <Alert
              type="error"
              showIcon
              message={normalizeApiError(mutation.error, '模式切换提交失败')}
            />
          ) : null}

          <Form<FormValues> form={form} layout="vertical" onFinish={(values) => mutation.mutate(values)}>
            <Form.Item
              label="目标模式"
              name="targetMode"
              rules={[{ required: true, message: '请选择目标模式' }]}
            >
              <Select
                options={[
                  { label: 'audit', value: 'audit' },
                  { label: 'alert', value: 'alert' },
                  { label: 'warn', value: 'warn' },
                  { label: 'enforce', value: 'enforce' }
                ]}
              />
            </Form.Item>

            <Form.Item label="分配 ID（逗号分隔，可选）" name="assignmentIds">
              <Input placeholder="例如：assign-1,assign-2" />
            </Form.Item>

            <Form.Item label="切换原因（可选）" name="reason">
              <Input.TextArea autoSize={{ minRows: 3, maxRows: 5 }} placeholder="例如：灰度验证通过，扩大 enforce 范围" />
            </Form.Item>
          </Form>
        </Space>
      )}
    </Drawer>
  );
};
