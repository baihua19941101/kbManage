import { useEffect } from 'react';
import { useMutation } from '@tanstack/react-query';
import { Alert, Button, Drawer, Form, Input, Select, Space, Tag, message } from 'antd';
import { normalizeApiError } from '@/services/api/client';
import { reviewExceptionRequest } from '@/services/securityPolicy';
import type {
  PolicyExceptionRequestDTO,
  ReviewExceptionDecision,
  ReviewExceptionRequestDTO
} from '@/services/api/types';

type ExceptionReviewDrawerProps = {
  open: boolean;
  exception?: PolicyExceptionRequestDTO;
  readonly?: boolean;
  onClose: () => void;
  onSuccess?: () => void;
};

type FormValues = {
  decision: ReviewExceptionDecision;
  comment?: string;
};

const statusColorMap: Record<string, string> = {
  pending: 'gold',
  approved: 'blue',
  rejected: 'red',
  active: 'green',
  expired: 'default',
  revoked: 'volcano'
};

export const ExceptionReviewDrawer = ({
  open,
  exception,
  readonly,
  onClose,
  onSuccess
}: ExceptionReviewDrawerProps) => {
  const [form] = Form.useForm<FormValues>();

  useEffect(() => {
    if (!open) {
      return;
    }
    form.setFieldsValue({ decision: 'approve' });
  }, [form, open]);

  const mutation = useMutation({
    mutationFn: async (values: FormValues) => {
      if (!exception?.id) {
        throw new Error('缺少例外申请 ID，无法审批');
      }

      const payload: ReviewExceptionRequestDTO = {
        decision: values.decision,
        comment: values.comment?.trim() || undefined
      };
      return reviewExceptionRequest(exception.id, payload);
    },
    onSuccess: () => {
      message.success('例外审批已提交');
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
      title={exception ? `例外审批 - ${exception.id}` : '例外审批'}
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
            disabled={!exception || readonly}
            loading={mutation.isPending}
            onClick={() => form.submit()}
          >
            提交审批
          </Button>
        </Space>
      }
    >
      {!exception ? (
        <Alert type="info" showIcon message="请先在页面中选择一条例外申请" />
      ) : (
        <Space direction="vertical" size={12} style={{ width: '100%' }}>
          {readonly ? (
            <Alert
              type="info"
              showIcon
              message="当前为只读模式"
              description="你可以查看例外状态，但无法执行审批动作。"
            />
          ) : null}

          {mutation.error ? (
            <Alert
              type="error"
              showIcon
              message={normalizeApiError(mutation.error, '例外审批失败')}
            />
          ) : null}

          <div>
            当前状态：
            <Tag color={statusColorMap[exception.status] || 'default'}>{exception.status}</Tag>
          </div>

          <Form<FormValues> form={form} layout="vertical" onFinish={(values) => mutation.mutate(values)}>
            <Form.Item
              label="审批动作"
              name="decision"
              rules={[{ required: true, message: '请选择审批动作' }]}
            >
              <Select
                options={[
                  { label: 'approve', value: 'approve' },
                  { label: 'reject', value: 'reject' },
                  { label: 'revoke', value: 'revoke' }
                ]}
              />
            </Form.Item>

            <Form.Item label="审批备注（可选）" name="comment">
              <Input.TextArea autoSize={{ minRows: 3, maxRows: 5 }} placeholder="例如：同意临时例外，要求 24h 内完成整改" />
            </Form.Item>
          </Form>
        </Space>
      )}
    </Drawer>
  );
};
