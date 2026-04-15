import { useEffect } from 'react';
import { useMutation } from '@tanstack/react-query';
import { Alert, Button, Drawer, Form, Input, Space, Typography, message } from 'antd';
import { normalizeApiError } from '@/services/api/client';
import { createExceptionRequest } from '@/services/securityPolicy';
import type { CreateExceptionRequestDTO, PolicyHitRecordDTO } from '@/services/api/types';

type ExceptionRequestDrawerProps = {
  open: boolean;
  hit?: PolicyHitRecordDTO;
  readonly?: boolean;
  onClose: () => void;
  onSuccess?: () => void;
};

type FormValues = {
  reason: string;
  startsAt: string;
  expiresAt: string;
};

const normalizeDateTimeInput = (value: string): string => {
  const trimmed = value.trim();
  if (!trimmed) {
    return '';
  }

  const parsed = new Date(trimmed);
  if (Number.isNaN(parsed.getTime())) {
    throw new Error('时间格式不正确，请使用 ISO 时间或 datetime-local 格式');
  }
  return parsed.toISOString();
};

export const ExceptionRequestDrawer = ({
  open,
  hit,
  readonly,
  onClose,
  onSuccess
}: ExceptionRequestDrawerProps) => {
  const [form] = Form.useForm<FormValues>();

  useEffect(() => {
    if (!open) {
      return;
    }

    form.setFieldsValue({
      startsAt: new Date().toISOString(),
      expiresAt: new Date(Date.now() + 24 * 60 * 60 * 1000).toISOString()
    });
  }, [form, open]);

  const mutation = useMutation({
    mutationFn: async (values: FormValues) => {
      if (!hit?.id) {
        throw new Error('缺少命中记录 ID，无法提交例外申请');
      }

      const payload: CreateExceptionRequestDTO = {
        reason: values.reason.trim(),
        startsAt: normalizeDateTimeInput(values.startsAt),
        expiresAt: normalizeDateTimeInput(values.expiresAt)
      };

      return createExceptionRequest(hit.id, payload);
    },
    onSuccess: () => {
      message.success('例外申请已提交');
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
      title={hit ? `提交例外申请 - ${hit.resourceKind}/${hit.resourceName}` : '提交例外申请'}
      open={open}
      width={620}
      destroyOnClose
      getContainer={false}
      onClose={handleClose}
      footer={
        <Space>
          <Button onClick={handleClose}>取消</Button>
          <Button
            type="primary"
            disabled={!hit || readonly}
            loading={mutation.isPending}
            onClick={() => form.submit()}
          >
            提交申请
          </Button>
        </Space>
      }
    >
      {!hit ? (
        <Alert type="info" showIcon message="请先在页面中选择一条命中记录" />
      ) : (
        <Space direction="vertical" size={12} style={{ width: '100%' }}>
          {readonly ? (
            <Alert
              type="info"
              showIcon
              message="当前为只读模式"
              description="你可查看例外状态，但无法提交新的例外申请。"
            />
          ) : null}

          {mutation.error ? (
            <Alert
              type="error"
              showIcon
              message={normalizeApiError(mutation.error, '例外申请提交失败')}
            />
          ) : null}

          <Typography.Text type="secondary">
            命中对象：{hit.clusterId}/{hit.namespace}/{hit.resourceKind}/{hit.resourceName}
          </Typography.Text>

          <Form<FormValues> form={form} layout="vertical" onFinish={(values) => mutation.mutate(values)}>
            <Form.Item label="申请原因" name="reason" rules={[{ required: true, message: '请输入申请原因' }]}>
              <Input.TextArea autoSize={{ minRows: 3, maxRows: 6 }} placeholder="例如：业务高峰期间临时放行，待变更窗口修复" />
            </Form.Item>

            <Form.Item
              label="生效时间（ISO 或 datetime-local）"
              name="startsAt"
              rules={[{ required: true, message: '请输入生效时间' }]}
            >
              <Input placeholder="例如：2026-04-15T08:00:00Z" />
            </Form.Item>

            <Form.Item
              label="过期时间（ISO 或 datetime-local）"
              name="expiresAt"
              rules={[{ required: true, message: '请输入过期时间' }]}
            >
              <Input placeholder="例如：2026-04-16T08:00:00Z" />
            </Form.Item>
          </Form>
        </Space>
      )}
    </Drawer>
  );
};
