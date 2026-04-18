import { useMemo, useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Alert, Button, Card, Empty, Form, Input, Select, Space, Table } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import {
  canCreateClusterLifecycle,
  canReadClusterLifecycle,
  useAuthStore
} from '@/features/auth/store';
import { ClusterTemplateDrawer } from '@/features/cluster-lifecycle/components/ClusterTemplateDrawer';
import { PageHeader } from '@/features/cluster-lifecycle/components/PageHeader';
import { LifecycleStatusTag } from '@/features/cluster-lifecycle/components/status';
import { useLifecycleAction } from '@/features/cluster-lifecycle/hooks/useLifecycleAction';
import { normalizeApiError } from '@/services/api/client';
import {
  clusterLifecycleQueryKeys,
  hasBlockingValidationIssue,
  listClusterDrivers,
  listClusterTemplates,
  type ClusterTemplate,
  type CreateClusterRequest
} from '@/services/clusterLifecycle';

const templateColumns: ColumnsType<ClusterTemplate> = [
  { title: '模板', dataIndex: 'name', key: 'name' },
  { title: '基础设施', dataIndex: 'infrastructureType', key: 'infrastructureType' },
  { title: '驱动', dataIndex: 'driverKey', key: 'driverKey' },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
    render: (value?: string) => <LifecycleStatusTag value={value} />
  }
];

type ProvisionFormValues = CreateClusterRequest;

export const ClusterProvisionPage = () => {
  const [form] = Form.useForm<ProvisionFormValues>();
  const user = useAuthStore((state) => state.user);
  const canRead = canReadClusterLifecycle(user);
  const canCreate = canCreateClusterLifecycle(user);
  const [selectedTemplate, setSelectedTemplate] = useState<ClusterTemplate | null>(null);
  const [drawerOpen, setDrawerOpen] = useState(false);
  const { provisionMutation, validateTemplateMutation } = useLifecycleAction();

  const driversQuery = useQuery({
    queryKey: clusterLifecycleQueryKeys.drivers(),
    enabled: canRead,
    queryFn: () => listClusterDrivers()
  });
  const templatesQuery = useQuery({
    queryKey: clusterLifecycleQueryKeys.templates(),
    enabled: canRead,
    queryFn: () => listClusterTemplates()
  });

  const driverOptions = useMemo(
    () =>
      (driversQuery.data?.items || []).map((item) => ({
        label: `${item.displayName || item.driverKey} / ${item.version}`,
        value: item.id
      })),
    [driversQuery.data?.items]
  );

  if (!canRead) {
    return <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description="你暂无集群创建访问权限。" />;
  }

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <PageHeader
        title="模板化创建集群"
        description="基于驱动和模板创建新集群，并在提交前查看兼容性和阻断项。"
      />

      {templatesQuery.error ? (
        <Alert
          type="error"
          showIcon
          message="模板列表加载失败"
          description={normalizeApiError(templatesQuery.error, '模板列表加载失败，请稍后重试。')}
        />
      ) : null}

      {!canCreate ? (
        <Alert type="warning" showIcon message="当前账号没有创建集群的动作权限。" />
      ) : null}

      <Card size="small" title="创建请求">
        <Form
          form={form}
          layout="vertical"
          disabled={!canCreate}
          onFinish={(values) => provisionMutation.mutate(values)}
        >
          <Form.Item name="name" label="集群名称" rules={[{ required: true, message: '请输入集群名称' }]}>
            <Input placeholder="prod-hz-core-01" />
          </Form.Item>
          <Form.Item
            name="infrastructureType"
            label="基础设施类型"
            rules={[{ required: true, message: '请输入基础设施类型' }]}
          >
            <Input placeholder="cloud / baremetal / virtualized" />
          </Form.Item>
          <Form.Item name="driverRef" label="驱动版本" rules={[{ required: true, message: '请选择驱动版本' }]}>
            <Select options={driverOptions} loading={driversQuery.isLoading} />
          </Form.Item>
          <Form.Item name="templateId" label="模板 ID" rules={[{ required: true, message: '请输入模板 ID' }]}>
            <Input placeholder="template-prod-standard" />
          </Form.Item>
          <Form.Item name="workspaceId" label="工作空间 ID">
            <Input placeholder="workspace-platform" />
          </Form.Item>
          <Form.Item name="projectId" label="项目 ID">
            <Input placeholder="project-control-plane" />
          </Form.Item>
          <Space>
            <Button type="primary" htmlType="submit" loading={provisionMutation.isPending} disabled={!canCreate}>
              提交创建
            </Button>
            <Button
              onClick={() => {
                const templateId = form.getFieldValue('templateId');
                const driverRef = form.getFieldValue('driverRef');
                if (!templateId) {
                  return;
                }
                validateTemplateMutation.mutate({ templateId, payload: { driverRef } });
              }}
              loading={validateTemplateMutation.isPending}
            >
              执行创建前校验
            </Button>
          </Space>
        </Form>
      </Card>

      <Card size="small" title={`可用模板（${templatesQuery.data?.items.length ?? 0}）`}>
        <Table<ClusterTemplate>
          rowKey={(record) => record.id}
          columns={[
            ...templateColumns,
            {
              title: '操作',
              key: 'actions',
              render: (_, record) => (
                <Button
                  type="link"
                  onClick={() => {
                    setSelectedTemplate(record);
                    setDrawerOpen(true);
                  }}
                >
                  查看模板
                </Button>
              )
            }
          ]}
          loading={templatesQuery.isLoading || templatesQuery.isFetching}
          dataSource={templatesQuery.data?.items || []}
          pagination={{ pageSize: 6 }}
        />
      </Card>

      {validateTemplateMutation.data ? (
        <Alert
          type={hasBlockingValidationIssue(validateTemplateMutation.data) ? 'error' : 'success'}
          showIcon
          message={`创建前校验：${validateTemplateMutation.data.overallStatus || '未知'}`}
          description={`阻断项 ${(validateTemplateMutation.data.blockers || []).length} 个，风险提示 ${(validateTemplateMutation.data.warnings || []).length} 个。`}
        />
      ) : null}

      <ClusterTemplateDrawer
        open={drawerOpen}
        template={selectedTemplate}
        validation={
          selectedTemplate?.id === form.getFieldValue('templateId')
            ? validateTemplateMutation.data
            : null
        }
        loading={validateTemplateMutation.isPending}
        onClose={() => setDrawerOpen(false)}
      />
    </Space>
  );
};
