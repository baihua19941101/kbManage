import { useState } from 'react';
import { useMutation } from '@tanstack/react-query';
import { Alert, Button, Drawer, Form, Input, message } from 'antd';
import { normalizeErrorMessage } from '@/app/queryClient';
import { createCluster } from '@/services/clusters';

type ClusterOnboardFormValues = {
  name: string;
  kubeconfig: string;
  description?: string;
};

type ClusterOnboardDrawerProps = {
  open: boolean;
  onClose: () => void;
  onSuccess: (clusterName: string) => void;
};

export const ClusterOnboardDrawer = ({ open, onClose, onSuccess }: ClusterOnboardDrawerProps) => {
  const [form] = Form.useForm<ClusterOnboardFormValues>();
  const [submitError, setSubmitError] = useState<string | null>(null);

  const onboardMutation = useMutation({
    mutationFn: createCluster,
    onSuccess: (result) => {
      message.success(`集群 ${result.name} 已提交接入`);
      onSuccess(result.name);
      form.resetFields();
      setSubmitError(null);
      onClose();
    },
    onError: (error) => {
      setSubmitError(normalizeErrorMessage(error, '接入失败，请稍后重试'));
    }
  });

  const onSubmit = async () => {
    try {
      const values = await form.validateFields();
      setSubmitError(null);
      onboardMutation.mutate({
        name: values.name,
        credentialType: 'kubeconfig',
        credentialPayload: values.kubeconfig,
        description: values.description
      });
    } catch {
      return;
    }
  };

  return (
    <Drawer
      title="接入 Kubernetes 集群"
      width={520}
      open={open}
      onClose={onClose}
      extra={
        <Button
          type="primary"
          loading={onboardMutation.isPending}
          onClick={() => void onSubmit()}
        >
          提交接入
        </Button>
      }
    >
      {submitError ? (
        <Alert type="error" showIcon message={submitError} style={{ marginBottom: 16 }} />
      ) : null}
      <Form
        layout="vertical"
        form={form}
        initialValues={{
          name: '',
          kubeconfig: '',
          description: ''
        }}
      >
        <Form.Item
          label="Cluster Name"
          name="name"
          rules={[{ required: true, message: '请输入集群名称' }]}
        >
          <Input placeholder="例如 prod-cn" />
        </Form.Item>
        <Form.Item label="描述" name="description">
          <Input placeholder="可选，描述集群用途" />
        </Form.Item>
        <Form.Item
          label="kubeconfig"
          name="kubeconfig"
          rules={[
            { required: true, message: '请粘贴 kubeconfig 内容' },
            { min: 20, message: 'kubeconfig 内容长度过短' }
          ]}
        >
          <Input.TextArea rows={12} placeholder="apiVersion: v1..." />
        </Form.Item>
      </Form>
    </Drawer>
  );
};
