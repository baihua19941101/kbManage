import { useState } from 'react';
import { Button, Drawer, Form, Input, message } from 'antd';

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

const mockOnboardCluster = async (payload: ClusterOnboardFormValues) => {
  await new Promise((resolve) => {
    setTimeout(resolve, 500);
  });
  return {
    clusterId: `cluster-${Date.now()}`,
    name: payload.name
  };
};

export const ClusterOnboardDrawer = ({ open, onClose, onSuccess }: ClusterOnboardDrawerProps) => {
  const [form] = Form.useForm<ClusterOnboardFormValues>();
  const [submitting, setSubmitting] = useState(false);

  const onSubmit = async () => {
    try {
      const values = await form.validateFields();
      setSubmitting(true);
      const result = await mockOnboardCluster(values);
      message.success(`集群 ${result.name} 已提交接入`);
      onSuccess(result.name);
      form.resetFields();
      onClose();
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <Drawer
      title="接入 Kubernetes 集群"
      width={520}
      open={open}
      onClose={onClose}
      extra={
        <Button type="primary" loading={submitting} onClick={() => void onSubmit()}>
          提交接入
        </Button>
      }
    >
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
