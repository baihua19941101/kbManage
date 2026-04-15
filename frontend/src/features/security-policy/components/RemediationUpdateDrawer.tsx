import { useEffect } from 'react';
import { useMutation } from '@tanstack/react-query';
import { Alert, Button, Drawer, Form, Input, Radio, Space, message } from 'antd';
import { normalizeApiError } from '@/services/api/client';
import { updatePolicyHitRemediation } from '@/services/securityPolicy';
import type {
  PolicyHitRecordDTO,
  PolicyRemediationStatus,
  UpdatePolicyRemediationRequestDTO
} from '@/services/api/types';

type RemediationUpdateDrawerProps = {
  open: boolean;
  violation?: PolicyHitRecordDTO;
  readonly?: boolean;
  onClose: () => void;
  onSuccess?: () => void;
};

type FormValues = {
  remediationStatus: PolicyRemediationStatus;
  comment?: string;
};

export const RemediationUpdateDrawer = ({
  open,
  violation,
  readonly,
  onClose,
  onSuccess
}: RemediationUpdateDrawerProps) => {
  const [form] = Form.useForm<FormValues>();

  useEffect(() => {
    if (!open) {
      return;
    }
    form.setFieldsValue({
      remediationStatus: violation?.remediationStatus || 'open',
      comment: undefined
    });
  }, [form, open, violation?.remediationStatus]);

  const mutation = useMutation({
    mutationFn: async (values: FormValues) => {
      if (!violation?.id) {
        throw new Error('缺少违规记录 ID，无法更新整改状态');
      }
      const payload: UpdatePolicyRemediationRequestDTO = {
        remediationStatus: values.remediationStatus,
        comment: values.comment?.trim() || undefined
      };
      return updatePolicyHitRemediation(violation.id, payload);
    },
    onSuccess: () => {
      message.success('整改状态更新已提交');
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
      title={violation ? `更新整改状态 - ${violation.id}` : '更新整改状态'}
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
            disabled={readonly || !violation}
            loading={mutation.isPending}
            onClick={() => form.submit()}
          >
            提交更新
          </Button>
        </Space>
      }
    >
      {!violation ? (
        <Alert type="info" showIcon message="请先选择一条违规记录" />
      ) : (
        <Space direction="vertical" size={12} style={{ width: '100%' }}>
          {readonly ? (
            <Alert
              type="info"
              showIcon
              message="当前为只读模式"
              description="你可以查看违规状态，但无法提交整改更新。"
            />
          ) : null}

          {mutation.error ? (
            <Alert
              type="error"
              showIcon
              message={normalizeApiError(mutation.error, '整改状态更新失败')}
            />
          ) : null}

          <Form<FormValues> form={form} layout="vertical" onFinish={(values) => mutation.mutate(values)}>
            <Form.Item
              label="整改状态"
              name="remediationStatus"
              rules={[{ required: true, message: '请选择整改状态' }]}
            >
              <Radio.Group>
                <Space direction="vertical">
                  <Radio value="open">open</Radio>
                  <Radio value="in_progress">in_progress</Radio>
                  <Radio value="mitigated">mitigated</Radio>
                  <Radio value="closed">closed</Radio>
                </Space>
              </Radio.Group>
            </Form.Item>

            <Form.Item label="处置备注（可选）" name="comment">
              <Input.TextArea
                autoSize={{ minRows: 3, maxRows: 5 }}
                placeholder="例如：已完成镜像修复并重新部署"
              />
            </Form.Item>
          </Form>
        </Space>
      )}
    </Drawer>
  );
};
