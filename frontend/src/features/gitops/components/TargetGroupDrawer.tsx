import { useEffect } from 'react';
import { useMutation } from '@tanstack/react-query';
import { Alert, Button, Drawer, Form, Input, Space, message } from 'antd';
import { normalizeApiError } from '@/services/api/client';
import {
  createGitOpsTargetGroup,
  updateGitOpsTargetGroup,
  type GitOpsTargetGroupFormData,
  type GitOpsTargetGroupItem
} from '@/services/gitops';

type TargetGroupDrawerProps = {
  open: boolean;
  onClose: () => void;
  targetGroup?: GitOpsTargetGroupItem;
  onSuccess?: (targetGroup: GitOpsTargetGroupItem) => void;
};

type TargetGroupFormValues = {
  name: string;
  workspaceId: string | number;
  projectId?: string | number;
  clusterRefsText?: string;
  selectorSummary?: string;
};

const toNumber = (value?: string | number): number | undefined => {
  if (value === undefined || value === null || value === '') {
    return undefined;
  }

  const parsed = Number(value);
  return Number.isNaN(parsed) ? undefined : parsed;
};

const parseClusterRefs = (raw?: string): number[] | undefined => {
  if (!raw) {
    return undefined;
  }

  const refs = raw
    .split(',')
    .map((item) => Number(item.trim()))
    .filter((item) => !Number.isNaN(item));

  return refs.length > 0 ? refs : undefined;
};

export const TargetGroupDrawer = ({
  open,
  onClose,
  targetGroup,
  onSuccess
}: TargetGroupDrawerProps) => {
  const [form] = Form.useForm<TargetGroupFormValues>();
  const isEdit = Boolean(targetGroup);

  useEffect(() => {
    if (!open) {
      return;
    }

    form.setFieldsValue({
      name: targetGroup?.name,
      workspaceId: targetGroup?.workspaceId ? String(targetGroup.workspaceId) : undefined,
      projectId: targetGroup?.projectId ? String(targetGroup.projectId) : undefined,
      clusterRefsText: targetGroup?.clusterRefs?.join(','),
      selectorSummary: targetGroup?.selectorSummary
    });
  }, [open, targetGroup, form]);

  const mutation = useMutation({
    mutationFn: async (values: TargetGroupFormValues) => {
      const payload: GitOpsTargetGroupFormData = {
        name: values.name.trim(),
        workspaceId: toNumber(values.workspaceId) ?? 0,
        projectId: toNumber(values.projectId),
        clusterRefs: parseClusterRefs(values.clusterRefsText),
        selectorSummary: values.selectorSummary?.trim() || undefined
      };

      if (targetGroup?.id !== undefined && targetGroup.id !== null) {
        return updateGitOpsTargetGroup(targetGroup.id, payload);
      }

      return createGitOpsTargetGroup(payload);
    },
    onSuccess: (savedGroup) => {
      message.success(isEdit ? '目标组已更新' : '目标组已创建');
      onSuccess?.(savedGroup);
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
      title={isEdit ? '编辑目标组' : '新建目标组'}
      open={open}
      width={520}
      destroyOnClose
      getContainer={false}
      onClose={handleClose}
      footer={
        <Space>
          <Button onClick={handleClose}>取消</Button>
          <Button type="primary" loading={mutation.isPending} onClick={() => form.submit()}>
            保存目标组
          </Button>
        </Space>
      }
    >
      <Space direction="vertical" size={12} style={{ width: '100%' }}>
        {mutation.error ? (
          <Alert type="error" showIcon message={normalizeApiError(mutation.error, '保存目标组失败')} />
        ) : null}
        <Form<TargetGroupFormValues> form={form} layout="vertical" onFinish={(values) => mutation.mutate(values)}>
          <Form.Item label="目标组名称" name="name" rules={[{ required: true, message: '请输入目标组名称' }]}>
            <Input placeholder="例如：prod-cn-group" />
          </Form.Item>
          <Form.Item label="工作空间 ID" name="workspaceId" rules={[{ required: true, message: '请输入工作空间 ID' }]}>
            <Input placeholder="例如：1001" />
          </Form.Item>
          <Form.Item label="项目 ID" name="projectId">
            <Input placeholder="例如：2001（可选）" />
          </Form.Item>
          <Form.Item label="集群 ID 列表" name="clusterRefsText">
            <Input placeholder="例如：1,2,3" />
          </Form.Item>
          <Form.Item label="选择器摘要" name="selectorSummary">
            <Input.TextArea placeholder="例如：region=cn,env=prod" autoSize={{ minRows: 2, maxRows: 4 }} />
          </Form.Item>
        </Form>
      </Space>
    </Drawer>
  );
};
