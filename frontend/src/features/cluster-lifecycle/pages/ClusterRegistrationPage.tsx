import { useMemo } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Alert, Button, Card, Empty, Form, Input, Space } from 'antd';
import {
  canImportClusterLifecycle,
  canReadClusterLifecycle,
  useAuthStore
} from '@/features/auth/store';
import { PageHeader } from '@/features/cluster-lifecycle/components/PageHeader';
import { RegistrationGuideCard } from '@/features/cluster-lifecycle/components/RegistrationGuideCard';
import { useLifecycleAction } from '@/features/cluster-lifecycle/hooks/useLifecycleAction';
import { normalizeApiError } from '@/services/api/client';
import {
  clusterLifecycleQueryKeys,
  listClusterLifecycleRecords,
  type RegisterClusterRequest
} from '@/services/clusterLifecycle';

export const ClusterRegistrationPage = () => {
  const [form] = Form.useForm<RegisterClusterRequest>();
  const user = useAuthStore((state) => state.user);
  const canRead = canReadClusterLifecycle(user);
  const canRegister = canImportClusterLifecycle(user);
  const { registerMutation } = useLifecycleAction();

  const pendingRegistrationsQuery = useQuery({
    queryKey: clusterLifecycleQueryKeys.clusters('registration:pending'),
    enabled: canRead,
    queryFn: () => listClusterLifecycleRecords({ status: 'pending' })
  });

  const pendingCount = useMemo(
    () => pendingRegistrationsQuery.data?.items.length ?? 0,
    [pendingRegistrationsQuery.data?.items]
  );

  if (!canRead) {
    return <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description="你暂无集群注册访问权限。" />;
  }

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <PageHeader
        title="注册新集群"
        description="为待纳管集群生成接入令牌和注册命令，跟踪接入状态。"
      />

      {pendingRegistrationsQuery.error ? (
        <Alert
          type="error"
          showIcon
          message="注册状态加载失败"
          description={normalizeApiError(
            pendingRegistrationsQuery.error,
            '注册状态加载失败，请稍后重试。'
          )}
        />
      ) : null}

      <Card size="small" title={`待完成接入（${pendingCount}）`}>
        <Form
          form={form}
          layout="vertical"
          disabled={!canRegister}
          onFinish={(values) => registerMutation.mutate(values)}
        >
          <Form.Item name="name" label="集群名称" rules={[{ required: true, message: '请输入集群名称' }]}>
            <Input placeholder="edge-cn-sh-01" />
          </Form.Item>
          <Form.Item
            name="infrastructureType"
            label="基础设施类型"
            rules={[{ required: true, message: '请输入基础设施类型' }]}
          >
            <Input placeholder="managed-kubernetes / cloud / baremetal" />
          </Form.Item>
          <Form.Item name="workspaceId" label="工作空间 ID">
            <Input placeholder="workspace-platform" />
          </Form.Item>
          <Form.Item name="projectId" label="项目 ID">
            <Input placeholder="project-edge" />
          </Form.Item>
          <Button
            type="primary"
            htmlType="submit"
            loading={registerMutation.isPending}
            disabled={!canRegister}
          >
            生成注册指引
          </Button>
        </Form>
      </Card>

      {!canRegister ? (
        <Alert type="warning" showIcon message="当前账号没有注册新集群的动作权限。" />
      ) : null}

      <RegistrationGuideCard bundle={registerMutation.data} loading={registerMutation.isPending} />
    </Space>
  );
};
