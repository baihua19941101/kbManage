import { useMemo, useState } from 'react';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import { Alert, Button, Card, Empty, Space, Table, Tag, Typography } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { Link, useLocation } from 'react-router-dom';
import { queryKeys } from '@/app/queryClient';
import {
  canManageGitOpsSource,
  canReadGitOps,
  useAuthStore
} from '@/features/auth/store';
import { SourceFormDrawer } from '@/features/gitops/components/SourceFormDrawer';
import { TargetGroupDrawer } from '@/features/gitops/components/TargetGroupDrawer';
import { isAuthorizationError, normalizeApiError } from '@/services/api/client';
import {
  listGitOpsDeliveryUnits,
  listGitOpsSources,
  listGitOpsTargetGroups,
  type GitOpsDeliveryUnitItem,
  type GitOpsSourceItem,
  type GitOpsTargetGroupItem
} from '@/services/gitops';

const normalizeSearchValue = (value: string | null): string | undefined => {
  const normalized = value?.trim();
  return normalized && normalized.length > 0 ? normalized : undefined;
};

const statusColorMap: Record<string, string> = {
  ready: 'green',
  enabled: 'green',
  active: 'green',
  pending: 'gold',
  progressing: 'blue',
  degraded: 'red',
  failed: 'red',
  disabled: 'default',
  stale: 'gold',
  out_of_sync: 'orange',
  paused: 'purple',
  unknown: 'default',
  drifted: 'orange',
  in_sync: 'green'
};

const sourceColumns: ColumnsType<GitOpsSourceItem> = [
  {
    title: '来源名称',
    dataIndex: 'name',
    key: 'name'
  },
  {
    title: '类型',
    dataIndex: 'sourceType',
    key: 'sourceType'
  },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
    render: (status?: string) => (
      <Tag color={statusColorMap[status || 'unknown'] || 'default'}>{status || 'unknown'}</Tag>
    )
  },
  {
    title: '地址',
    dataIndex: 'endpoint',
    key: 'endpoint'
  }
];

const targetGroupColumns: ColumnsType<GitOpsTargetGroupItem> = [
  {
    title: '目标组',
    dataIndex: 'name',
    key: 'name'
  },
  {
    title: '集群数',
    dataIndex: 'clusterRefs',
    key: 'clusterRefs',
    render: (clusterRefs?: number[]) => clusterRefs?.length ?? 0
  },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
    render: (status?: string) => (
      <Tag color={statusColorMap[status || 'unknown'] || 'default'}>{status || 'unknown'}</Tag>
    )
  }
];

const deliveryUnitColumns: ColumnsType<GitOpsDeliveryUnitItem> = [
  {
    title: '交付单元',
    dataIndex: 'name',
    key: 'name'
  },
  {
    title: '交付状态',
    dataIndex: 'deliveryStatus',
    key: 'deliveryStatus',
    render: (status?: string) => (
      <Tag color={statusColorMap[status || 'unknown'] || 'default'}>{status || 'unknown'}</Tag>
    )
  },
  {
    title: '漂移状态',
    dataIndex: 'driftStatus',
    key: 'driftStatus',
    render: (status?: string) => (
      <Tag color={statusColorMap[status || 'unknown'] || 'default'}>{status || 'unknown'}</Tag>
    )
  },
  {
    title: '最近同步结果',
    dataIndex: 'lastSyncResult',
    key: 'lastSyncResult'
  },
  {
    title: '操作',
    key: 'actions',
    render: (_, unit) => (
      <Link to={`/gitops/delivery-units/${unit.id}`}>查看详情</Link>
    )
  }
];

export const GitOpsOverviewPage = () => {
  const { search } = useLocation();
  const queryClient = useQueryClient();
  const user = useAuthStore((state) => state.user);
  const canRead = canReadGitOps(user);
  const canManage = canManageGitOpsSource(user);
  const [sourceDrawerOpen, setSourceDrawerOpen] = useState(false);
  const [targetGroupDrawerOpen, setTargetGroupDrawerOpen] = useState(false);

  const searchParams = useMemo(() => new URLSearchParams(search), [search]);
  const keyword = normalizeSearchValue(searchParams.get('keyword'));
  const workspaceId = normalizeSearchValue(searchParams.get('workspaceId'));
  const projectId = normalizeSearchValue(searchParams.get('projectId'));
  const listQuery = useMemo(
    () => ({ keyword, workspaceId, projectId }),
    [keyword, workspaceId, projectId]
  );

  const refreshOverview = () => {
    void queryClient.invalidateQueries({ queryKey: queryKeys.gitops.sources() });
    void queryClient.invalidateQueries({ queryKey: queryKeys.gitops.deliveryUnits() });
  };

  const sourcesQuery = useQuery({
    queryKey: queryKeys.gitops.sources(`overview:${keyword || 'all'}`),
    enabled: canRead,
    queryFn: () => listGitOpsSources(listQuery)
  });

  const targetGroupsQuery = useQuery({
    queryKey: queryKeys.gitops.sources(`target-groups:${workspaceId || 'all'}`),
    enabled: canRead,
    queryFn: () => listGitOpsTargetGroups(listQuery)
  });

  const unitsQuery = useQuery({
    queryKey: queryKeys.gitops.deliveryUnits(`overview:${keyword || 'all'}`),
    enabled: canRead,
    queryFn: () => listGitOpsDeliveryUnits(listQuery)
  });

  if (!canRead) {
    return (
      <Empty
        image={Empty.PRESENTED_IMAGE_SIMPLE}
        description="你暂无 GitOps 访问权限，请联系管理员授予范围权限。"
      />
    );
  }

  const sourceCount = sourcesQuery.data?.items.length ?? 0;
  const targetGroupCount = targetGroupsQuery.data?.items.length ?? 0;
  const unitCount = unitsQuery.data?.items.length ?? 0;
  const sourceAuthError = isAuthorizationError(sourcesQuery.error);
  const targetGroupAuthError = isAuthorizationError(targetGroupsQuery.error);
  const unitAuthError = isAuthorizationError(unitsQuery.error);

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <div>
        <Typography.Title level={3} style={{ marginBottom: 8 }}>
          GitOps 与发布中心
        </Typography.Title>
        <Typography.Text type="secondary">
          管理交付来源、目标组和交付单元，查看状态聚合并进入详情。
        </Typography.Text>
      </div>

      {keyword ? (
        <Alert
          type="info"
          showIcon
          message={`已按资源上下文筛选：${keyword}`}
          description="该筛选通常来自资源页跳转，用于快速定位交付单元。"
        />
      ) : null}

      <Space wrap>
        <Button type="primary" disabled={!canManage} onClick={() => setSourceDrawerOpen(true)}>
          新建来源
        </Button>
        <Button disabled={!canManage} onClick={() => setTargetGroupDrawerOpen(true)}>
          新建目标组
        </Button>
      </Space>

      {sourceAuthError || targetGroupAuthError || unitAuthError ? (
        <Alert
          type="warning"
          showIcon
          message="权限已变更"
          description="当前账号可能缺少 GitOps 所需权限，请刷新页面或重新登录后重试。"
        />
      ) : null}

      {sourcesQuery.error && !sourceAuthError ? (
        <Alert
          type="error"
          showIcon
          message="来源列表加载失败"
          description={normalizeApiError(sourcesQuery.error, '来源列表加载失败，请稍后重试。')}
        />
      ) : null}

      {unitsQuery.error && !unitAuthError ? (
        <Alert
          type="error"
          showIcon
          message="交付单元加载失败"
          description={normalizeApiError(unitsQuery.error, '交付单元加载失败，请稍后重试。')}
        />
      ) : null}

      <Card size="small" title={`交付来源（${sourceCount}）`}>
        <Table<GitOpsSourceItem>
          rowKey={(record) => String(record.id)}
          loading={sourcesQuery.isLoading || sourcesQuery.isFetching}
          columns={sourceColumns}
          dataSource={sourcesQuery.data?.items || []}
          pagination={{ pageSize: 5 }}
        />
      </Card>

      <Card size="small" title={`目标组（${targetGroupCount}）`}>
        <Table<GitOpsTargetGroupItem>
          rowKey={(record) => String(record.id)}
          loading={targetGroupsQuery.isLoading || targetGroupsQuery.isFetching}
          columns={targetGroupColumns}
          dataSource={targetGroupsQuery.data?.items || []}
          pagination={{ pageSize: 5 }}
        />
      </Card>

      <Card size="small" title={`交付单元（${unitCount}）`}>
        <Table<GitOpsDeliveryUnitItem>
          rowKey={(record) => String(record.id)}
          loading={unitsQuery.isLoading || unitsQuery.isFetching}
          columns={deliveryUnitColumns}
          dataSource={unitsQuery.data?.items || []}
          pagination={{ pageSize: 8 }}
        />
      </Card>

      <SourceFormDrawer
        open={sourceDrawerOpen}
        onClose={() => setSourceDrawerOpen(false)}
        onSuccess={() => refreshOverview()}
      />
      <TargetGroupDrawer
        open={targetGroupDrawerOpen}
        onClose={() => setTargetGroupDrawerOpen(false)}
        onSuccess={() => refreshOverview()}
      />
    </Space>
  );
};
