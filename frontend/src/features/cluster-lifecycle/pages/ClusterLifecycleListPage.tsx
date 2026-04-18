import { useMemo, useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Alert, Button, Card, Empty, Form, Input, Select, Space } from 'antd';
import {
  canImportClusterLifecycle,
  canReadClusterLifecycle,
  useAuthStore
} from '@/features/auth/store';
import { ImportClusterDrawer } from '@/features/cluster-lifecycle/components/ImportClusterDrawer';
import { ClusterLifecycleTable } from '@/features/cluster-lifecycle/components/ClusterLifecycleTable';
import { LifecycleSummaryCards } from '@/features/cluster-lifecycle/components/LifecycleSummaryCards';
import { PageHeader } from '@/features/cluster-lifecycle/components/PageHeader';
import { normalizeApiError } from '@/services/api/client';
import {
  clusterLifecycleQueryKeys,
  listClusterLifecycleRecords,
  type ClusterLifecycleListQuery
} from '@/services/clusterLifecycle';
import { useLifecycleAction } from '@/features/cluster-lifecycle/hooks/useLifecycleAction';

const statusOptions = [
  { label: '全部状态', value: '' },
  { label: '运行中', value: 'active' },
  { label: '待处理', value: 'pending' },
  { label: '升级中', value: 'upgrading' },
  { label: '退役中', value: 'retiring' },
  { label: '失败', value: 'failed' }
];

type FilterValues = ClusterLifecycleListQuery;

export const ClusterLifecycleListPage = () => {
  const [form] = Form.useForm<FilterValues>();
  const user = useAuthStore((state) => state.user);
  const canRead = canReadClusterLifecycle(user);
  const canImport = canImportClusterLifecycle(user);
  const [drawerOpen, setDrawerOpen] = useState(false);
  const [filters, setFilters] = useState<ClusterLifecycleListQuery>({});
  const { importMutation } = useLifecycleAction();

  const queryKey = useMemo(
    () =>
      clusterLifecycleQueryKeys.clusters(
        [filters.keyword, filters.status, filters.infrastructureType, filters.driverKey]
          .filter(Boolean)
          .join(':')
      ),
    [filters]
  );

  const listQuery = useQuery({
    queryKey,
    enabled: canRead,
    queryFn: () => listClusterLifecycleRecords(filters)
  });

  if (!canRead) {
    return (
      <Empty
        image={Empty.PRESENTED_IMAGE_SIMPLE}
        description="你暂无集群生命周期中心访问权限。"
      />
    );
  }

  const clusters = listQuery.data?.items || [];

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <PageHeader
        title="集群生命周期中心"
        description="统一查看导入、注册、创建、升级、停用和退役中的多集群资产。"
        actions={
          <>
            <Button type="primary" disabled={!canImport} onClick={() => setDrawerOpen(true)}>
              导入已有集群
            </Button>
            <Button href="/cluster-lifecycle/register">注册新集群</Button>
            <Button href="/cluster-lifecycle/provision">模板化创建</Button>
          </>
        }
      />

      {listQuery.error ? (
        <Alert
          type="error"
          showIcon
          message="生命周期列表加载失败"
          description={normalizeApiError(listQuery.error, '生命周期列表加载失败，请稍后重试。')}
        />
      ) : null}

      <Card size="small" title="筛选条件">
        <Form
          form={form}
          layout="inline"
          onFinish={(values) =>
            setFilters({
              keyword: values.keyword?.trim() || undefined,
              status: values.status || undefined,
              infrastructureType: values.infrastructureType?.trim() || undefined,
              driverKey: values.driverKey?.trim() || undefined
            })
          }
        >
          <Form.Item name="keyword">
            <Input allowClear placeholder="搜索集群名 / 集群 ID" style={{ width: 220 }} />
          </Form.Item>
          <Form.Item name="status">
            <Select options={statusOptions} style={{ width: 160 }} />
          </Form.Item>
          <Form.Item name="infrastructureType">
            <Input allowClear placeholder="基础设施类型" style={{ width: 180 }} />
          </Form.Item>
          <Form.Item name="driverKey">
            <Input allowClear placeholder="驱动键" style={{ width: 160 }} />
          </Form.Item>
          <Form.Item>
            <Space>
              <Button type="primary" htmlType="submit" loading={listQuery.isFetching}>
                查询
              </Button>
              <Button
                onClick={() => {
                  form.resetFields();
                  setFilters({});
                }}
              >
                重置
              </Button>
            </Space>
          </Form.Item>
        </Form>
      </Card>

      <LifecycleSummaryCards clusters={clusters} />

      <Card size="small" title={`集群清单（${clusters.length}）`}>
        <ClusterLifecycleTable
          data={clusters}
          loading={listQuery.isLoading || listQuery.isFetching}
        />
      </Card>

      <ImportClusterDrawer
        open={drawerOpen}
        submitting={importMutation.isPending}
        onClose={() => setDrawerOpen(false)}
        onSubmit={(payload) =>
          importMutation.mutate(payload, { onSuccess: () => setDrawerOpen(false) })
        }
      />
    </Space>
  );
};
