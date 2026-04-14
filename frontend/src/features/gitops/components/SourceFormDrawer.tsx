import { useEffect } from 'react';
import { useMutation } from '@tanstack/react-query';
import { Alert, Button, Drawer, Form, Input, Select, Space, message } from 'antd';
import { normalizeApiError } from '@/services/api/client';
import {
  createGitOpsSource,
  updateGitOpsSource,
  type GitOpsSourceFormData,
  type GitOpsSourceItem
} from '@/services/gitops';

type SourceFormDrawerProps = {
  open: boolean;
  onClose: () => void;
  source?: GitOpsSourceItem;
  onSuccess?: (source: GitOpsSourceItem) => void;
};

type SourceFormValues = {
  name: string;
  sourceType: string;
  endpoint: string;
  defaultRef?: string;
  credentialRef?: string;
  workspaceId: string | number;
  projectId?: string | number;
};

const toNumber = (value?: string | number): number | undefined => {
  if (value === undefined || value === null || value === '') {
    return undefined;
  }

  const parsed = Number(value);
  return Number.isNaN(parsed) ? undefined : parsed;
};

export const SourceFormDrawer = ({ open, onClose, source, onSuccess }: SourceFormDrawerProps) => {
  const [form] = Form.useForm<SourceFormValues>();
  const isEdit = Boolean(source);

  useEffect(() => {
    if (!open) {
      return;
    }

    form.setFieldsValue({
      name: source?.name,
      sourceType: source?.sourceType ?? 'git',
      endpoint: source?.endpoint,
      defaultRef: source?.defaultRef,
      credentialRef: source?.credentialRef,
      workspaceId: source?.workspaceId ? String(source.workspaceId) : undefined,
      projectId: source?.projectId ? String(source.projectId) : undefined
    });
  }, [open, source, form]);

  const mutation = useMutation({
    mutationFn: async (values: SourceFormValues) => {
      const payload: GitOpsSourceFormData = {
        name: values.name.trim(),
        sourceType: values.sourceType,
        endpoint: values.endpoint.trim(),
        defaultRef: values.defaultRef?.trim() || undefined,
        credentialRef: values.credentialRef?.trim() || undefined,
        workspaceId: toNumber(values.workspaceId) ?? 0,
        projectId: toNumber(values.projectId)
      };

      if (source?.id !== undefined && source.id !== null) {
        return updateGitOpsSource(source.id, payload);
      }

      return createGitOpsSource(payload);
    },
    onSuccess: (savedSource) => {
      message.success(isEdit ? '交付来源已更新' : '交付来源已创建');
      onSuccess?.(savedSource);
      form.resetFields();
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
      title={isEdit ? '编辑交付来源' : '新建交付来源'}
      open={open}
      width={560}
      destroyOnClose
      getContainer={false}
      onClose={handleClose}
      footer={
        <Space>
          <Button onClick={handleClose}>取消</Button>
          <Button type="primary" loading={mutation.isPending} onClick={() => form.submit()}>
            保存来源
          </Button>
        </Space>
      }
    >
      <Space direction="vertical" size={12} style={{ width: '100%' }}>
        {mutation.error ? (
          <Alert type="error" showIcon message={normalizeApiError(mutation.error, '保存来源失败')} />
        ) : null}
        <Form<SourceFormValues> form={form} layout="vertical" onFinish={(values) => mutation.mutate(values)}>
          <Form.Item label="来源名称" name="name" rules={[{ required: true, message: '请输入来源名称' }]}>
            <Input placeholder="例如：payments-git" />
          </Form.Item>
          <Form.Item label="来源类型" name="sourceType" rules={[{ required: true, message: '请选择来源类型' }]}>
            <Select
              options={[
                { label: 'git', value: 'git' },
                { label: 'package', value: 'package' }
              ]}
            />
          </Form.Item>
          <Form.Item label="来源地址" name="endpoint" rules={[{ required: true, message: '请输入来源地址' }]}>
            <Input placeholder="https://git.example.com/repo.git" />
          </Form.Item>
          <Form.Item label="默认分支/版本" name="defaultRef">
            <Input placeholder="例如：main" />
          </Form.Item>
          <Form.Item label="凭据引用" name="credentialRef">
            <Input placeholder="例如：git-credential-prod" />
          </Form.Item>
          <Form.Item label="工作空间 ID" name="workspaceId" rules={[{ required: true, message: '请输入工作空间 ID' }]}>
            <Input placeholder="例如：1001" />
          </Form.Item>
          <Form.Item label="项目 ID" name="projectId">
            <Input placeholder="例如：2001（可选）" />
          </Form.Item>
        </Form>
      </Space>
    </Drawer>
  );
};
