import { useEffect } from 'react';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { Button, Drawer, Form, Input, Select, Space, message } from 'antd';
import {
  createComplianceBaseline,
  updateComplianceBaseline,
  type ComplianceBaseline,
  type ComplianceRecordStatus,
  type ComplianceStandardType,
  type CreateComplianceBaselineRequest,
  type UpdateComplianceBaselineRequest
} from '@/services/compliance';
import { normalizeErrorMessage } from '@/app/queryClient';

const standardOptions: { label: string; value: ComplianceStandardType }[] = [
  { label: 'CIS', value: 'cis' },
  { label: 'STIG', value: 'stig' },
  { label: '平台基线', value: 'platform-baseline' }
];

const statusOptions: { label: string; value: ComplianceRecordStatus }[] = [
  { label: '草稿', value: 'draft' },
  { label: '启用', value: 'active' },
  { label: '停用', value: 'disabled' },
  { label: '归档', value: 'archived' }
];

type FormValues = {
  name: string;
  standardType: ComplianceStandardType;
  version: string;
  description?: string;
  status?: ComplianceRecordStatus;
};

type BaselineFormDrawerProps = {
  open: boolean;
  baseline?: ComplianceBaseline;
  readonly?: boolean;
  onClose: () => void;
};

export const BaselineFormDrawer = ({ open, baseline, readonly, onClose }: BaselineFormDrawerProps) => {
  const [form] = Form.useForm<FormValues>();
  const queryClient = useQueryClient();

  useEffect(() => {
    if (!open) {
      return;
    }

    form.setFieldsValue({
      name: baseline?.name || '',
      standardType: baseline?.standardType || 'cis',
      version: baseline?.version || '',
      description: baseline?.description || '',
      status: baseline?.status || 'draft'
    });
  }, [baseline, form, open]);

  const mutation = useMutation({
    mutationFn: async (values: FormValues) => {
      if (baseline?.id) {
        const payload: UpdateComplianceBaselineRequest = {
          name: values.name,
          description: values.description,
          status: values.status
        };
        return updateComplianceBaseline(baseline.id, payload);
      }

      const payload: CreateComplianceBaselineRequest = {
        name: values.name,
        standardType: values.standardType,
        version: values.version,
        description: values.description
      };
      return createComplianceBaseline(payload);
    },
    onSuccess: () => {
      message.success(baseline?.id ? '基线已更新' : '基线已创建');
      void queryClient.invalidateQueries({ queryKey: ['compliance', 'baselines'] });
      onClose();
    },
    onError: (error) => {
      message.error(normalizeErrorMessage(error, '基线保存失败'));
    }
  });

  return (
    <Drawer
      title={baseline?.id ? '编辑基线' : '新建基线'}
      width={420}
      open={open}
      onClose={onClose}
      destroyOnClose
    >
      <Form<FormValues> form={form} layout="vertical" onFinish={(values) => mutation.mutate(values)}>
        <Form.Item label="基线名称" name="name" rules={[{ required: true, message: '请输入基线名称' }]}>
          <Input disabled={readonly} placeholder="例如：CIS Kubernetes 1.31" />
        </Form.Item>
        <Form.Item label="标准类型" name="standardType" rules={[{ required: true }]}> 
          <Select disabled={readonly || Boolean(baseline?.id)} options={standardOptions} />
        </Form.Item>
        <Form.Item label="版本" name="version" rules={[{ required: true, message: '请输入版本' }]}>
          <Input disabled={readonly || Boolean(baseline?.id)} placeholder="例如：v1.31 L1" />
        </Form.Item>
        {baseline?.id ? (
          <Form.Item label="状态" name="status">
            <Select disabled={readonly} options={statusOptions} />
          </Form.Item>
        ) : null}
        <Form.Item label="说明" name="description">
          <Input.TextArea disabled={readonly} rows={4} placeholder="记录适用范围、规则来源或加固口径。" />
        </Form.Item>
        <Space>
          <Button onClick={onClose}>取消</Button>
          <Button type="primary" htmlType="submit" loading={mutation.isPending} disabled={readonly}>
            保存
          </Button>
        </Space>
      </Form>
    </Drawer>
  );
};
