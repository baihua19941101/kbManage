import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Alert, Button, Card, Space, Table } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { HAPolicyDrawer } from '@/features/sre-scale/components/HAPolicyDrawer';
import { PageHeader } from '@/features/sre-scale/components/PageHeader';
import { PermissionDenied } from '@/features/sre-scale/components/PermissionDenied';
import { StatusTag } from '@/features/sre-scale/components/status';
import { useSREPermissions } from '@/features/sre-scale/hooks/permissions';
import { useScaleEvidenceAction } from '@/features/sre-scale/hooks/useScaleEvidenceAction';
import { normalizeApiError } from '@/services/api/client';
import { listHAPolicies, listMaintenanceWindows, type HAPolicy, type MaintenanceWindow } from '@/services/sreScale';

const haColumns: ColumnsType<HAPolicy> = [
  { title: '策略', dataIndex: 'name', key: 'name' },
  { title: '模式', dataIndex: 'deploymentMode', key: 'deploymentMode' },
  { title: '副本', dataIndex: 'replicaExpectation', key: 'replicaExpectation' },
  { title: '接管状态', dataIndex: 'takeoverStatus', key: 'takeoverStatus', render: (value?: string) => <StatusTag value={value} /> },
  { title: '状态', dataIndex: 'status', key: 'status', render: (value?: string) => <StatusTag value={value} /> }
];

const windowColumns: ColumnsType<MaintenanceWindow> = [
  { title: '窗口', dataIndex: 'name', key: 'name' },
  { title: '类型', dataIndex: 'windowType', key: 'windowType' },
  { title: '范围', dataIndex: 'scope', key: 'scope' },
  { title: '状态', dataIndex: 'status', key: 'status', render: (value?: string) => <StatusTag value={value} /> }
];

export const HAControlPage = () => {
  const permissions = useSREPermissions();
  const action = useScaleEvidenceAction();
  const [drawerOpen, setDrawerOpen] = useState(false);
  const haQuery = useQuery({ queryKey: ['sreScale', 'haPolicies'], queryFn: () => listHAPolicies({}) });
  const windowQuery = useQuery({ queryKey: ['sreScale', 'maintenanceWindows'], queryFn: () => listMaintenanceWindows() });

  if (!permissions.canRead) {
    return <PermissionDenied description="你暂无平台 SRE 访问权限。" />;
  }

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <PageHeader
        title="高可用治理"
        description="配置控制面高可用策略、维护窗口和故障切换门槛。"
        actions={
          <Button type="primary" disabled={!permissions.canManageHA} onClick={() => setDrawerOpen(true)}>
            新增高可用策略
          </Button>
        }
      />
      {haQuery.error ? <Alert type="error" showIcon message={normalizeApiError(haQuery.error, '高可用策略加载失败')} /> : null}
      {windowQuery.error ? <Alert type="error" showIcon message={normalizeApiError(windowQuery.error, '维护窗口加载失败')} /> : null}
      <Card size="small" title="高可用策略">
        <Table rowKey="id" columns={haColumns} dataSource={haQuery.data?.items || []} pagination={{ pageSize: 5 }} />
      </Card>
      <Card size="small" title="维护窗口">
        <Table rowKey="id" columns={windowColumns} dataSource={windowQuery.data?.items || []} pagination={{ pageSize: 5 }} />
      </Card>
      <HAPolicyDrawer
        open={drawerOpen}
        submitting={action.createHAPolicyMutation.isPending}
        onClose={() => setDrawerOpen(false)}
        onSubmit={(payload) =>
          action.createHAPolicyMutation.mutate(payload, {
            onSuccess: () => setDrawerOpen(false)
          })
        }
      />
    </Space>
  );
};
