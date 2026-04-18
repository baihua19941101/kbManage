import { Button, Drawer, Form, Input, Select, Space } from 'antd';
import type { ImportClusterRequest } from '@/services/clusterLifecycle';

const infrastructureOptions = [
  { label: '裸金属 / 自建', value: 'baremetal' },
  { label: 'VMware / 虚拟化', value: 'virtualized' },
  { label: '公有云', value: 'cloud' },
  { label: '托管 Kubernetes', value: 'managed-kubernetes' }
];

type Props = {
  open: boolean;
  submitting?: boolean;
  onClose: () => void;
  onSubmit: (payload: ImportClusterRequest) => void;
};

export const ImportClusterDrawer = ({ open, submitting, onClose, onSubmit }: Props) => {
  const [form] = Form.useForm<ImportClusterRequest>();

  return (
    <Drawer
      title="导入已有集群"
      width={520}
      open={open}
      onClose={() => {
        form.resetFields();
        onClose();
      }}
      destroyOnClose
    >
      <Form form={form} layout="vertical" onFinish={onSubmit} initialValues={{ infrastructureType: 'managed-kubernetes' }}>
        <Form.Item name="name" label="集群名称" rules={[{ required: true, message: '请输入集群名称' }]}>
          <Input placeholder="例如：prod-cn-hz-01" />
        </Form.Item>
        <Form.Item
          name="infrastructureType"
          label="基础设施类型"
          rules={[{ required: true, message: '请选择基础设施类型' }]}
        >
          <Select options={infrastructureOptions} />
        </Form.Item>
        <Form.Item
          name="accessEndpoint"
          label="接入地址"
          rules={[{ required: true, message: '请输入接入地址' }]}
        >
          <Input placeholder="https://10.0.0.8:6443" />
        </Form.Item>
        <Form.Item name="credentialRef" label="凭据引用">
          <Input placeholder="secret://cluster-lifecycle/prod-kubeconfig" />
        </Form.Item>
        <Form.Item name="workspaceId" label="工作空间 ID">
          <Input placeholder="workspace-shared" />
        </Form.Item>
        <Form.Item name="projectId" label="项目 ID">
          <Input placeholder="project-platform" />
        </Form.Item>
        <Space>
          <Button onClick={onClose}>取消</Button>
          <Button type="primary" htmlType="submit" loading={submitting}>
            提交导入
          </Button>
        </Space>
      </Form>
    </Drawer>
  );
};
