import { useEffect, useMemo, useState } from 'react';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import { Alert, Button, Card, Descriptions, Empty, Space, Table, Tag, Typography } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { Link, useParams } from 'react-router-dom';
import { queryKeys } from '@/app/queryClient';
import { canReadGitOps, useAuthStore } from '@/features/auth/store';
import { DiffSummaryPanel } from '@/features/gitops/components/DiffSummaryPanel';
import { EnvironmentStageEditor } from '@/features/gitops/components/EnvironmentStageEditor';
import { OverlaySummaryPanel } from '@/features/gitops/components/OverlaySummaryPanel';
import { PromotionTimeline } from '@/features/gitops/components/PromotionTimeline';
import { ReleaseActionDrawer } from '@/features/gitops/components/ReleaseActionDrawer';
import { RevisionHistoryPanel } from '@/features/gitops/components/RevisionHistoryPanel';
import { RollbackDialog } from '@/features/gitops/components/RollbackDialog';
import { isAuthorizationError, normalizeApiError } from '@/services/api/client';
import type { GitOpsOperationDTO } from '@/services/api/types';
import {
  getGitOpsDeliveryUnit,
  getGitOpsDeliveryUnitStatus,
  listGitOpsTargetGroups,
  type GitOpsDeliveryUnitStatus,
  type GitOpsEnvironmentStage,
  type GitOpsReleaseRevision
} from '@/services/gitops';

type EnvironmentStatusRecord = NonNullable<GitOpsDeliveryUnitStatus['environments']>[number];

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
  in_sync: 'green',
  succeeded: 'green'
};

const terminalOperationStatuses = new Set(['partially_succeeded', 'succeeded', 'failed', 'canceled']);

const areStagesEqual = (left: GitOpsEnvironmentStage[], right: GitOpsEnvironmentStage[]) => {
  if (left.length !== right.length) {
    return false;
  }

  return left.every((stage, index) => {
    const other = right[index];
    return (
      stage.name === other.name &&
      stage.orderIndex === other.orderIndex &&
      stage.targetGroupId === other.targetGroupId &&
      stage.promotionMode === other.promotionMode &&
      Boolean(stage.paused) === Boolean(other.paused)
    );
  });
};

const environmentStatusColumns: ColumnsType<EnvironmentStatusRecord> = [
  {
    title: '环境',
    dataIndex: 'environment',
    key: 'environment'
  },
  {
    title: '同步状态',
    dataIndex: 'syncStatus',
    key: 'syncStatus',
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
    title: '目标数',
    dataIndex: 'targetCount',
    key: 'targetCount'
  },
  {
    title: '成功/失败',
    key: 'result',
    render: (_, record) => `${record.succeededCount ?? 0}/${record.failedCount ?? 0}`
  }
];

export const DeliveryUnitDetailPage = () => {
  const { unitId } = useParams<{ unitId: string }>();
  const queryClient = useQueryClient();
  const user = useAuthStore((state) => state.user);
  const canRead = canReadGitOps(user);
  const [releaseActionOpen, setReleaseActionOpen] = useState(false);
  const [rollbackOpen, setRollbackOpen] = useState(false);
  const [selectedRevision, setSelectedRevision] = useState<GitOpsReleaseRevision>();
  const [latestOperation, setLatestOperation] = useState<GitOpsOperationDTO>();

  const detailQuery = useQuery({
    queryKey: queryKeys.gitops.deliveryUnits(`detail:${unitId || 'unknown'}`),
    enabled: canRead && Boolean(unitId),
    queryFn: () => getGitOpsDeliveryUnit(unitId as string)
  });

  const statusQuery = useQuery({
    queryKey: queryKeys.gitops.deliveryUnits(`status:${unitId || 'unknown'}`),
    enabled: canRead && Boolean(unitId),
    queryFn: () => getGitOpsDeliveryUnitStatus(unitId as string)
  });

  const targetGroupsQuery = useQuery({
    queryKey: queryKeys.gitops.sources(`target-groups:detail:${unitId || 'unknown'}`),
    enabled: canRead,
    queryFn: () => listGitOpsTargetGroups()
  });

  const [stages, setStages] = useState<GitOpsEnvironmentStage[]>(detailQuery.data?.environments || []);

  useEffect(() => {
    const nextStages = detailQuery.data?.environments || [];
    setStages((currentStages) => (areStagesEqual(currentStages, nextStages) ? currentStages : nextStages));
  }, [detailQuery.data?.environments]);

  const permissionChanged = useMemo(
    () =>
      isAuthorizationError(detailQuery.error) ||
      isAuthorizationError(statusQuery.error) ||
      isAuthorizationError(targetGroupsQuery.error),
    [detailQuery.error, statusQuery.error, targetGroupsQuery.error]
  );

  const handleOperationChange = (operation: GitOpsOperationDTO) => {
    setLatestOperation(operation);

    const normalizedStatus = (operation.status || '').toLowerCase();
    if (!unitId || !terminalOperationStatuses.has(normalizedStatus)) {
      return;
    }

    void queryClient.invalidateQueries({
      predicate: (query) => {
        const key = query.queryKey;
        return (
          Array.isArray(key) &&
          key[0] === 'gitops' &&
          key[1] === 'deliveryUnits' &&
          typeof key[2] === 'string' &&
          key[2].includes(String(unitId))
        );
      }
    });
  };

  if (!canRead) {
    return (
      <Empty
        image={Empty.PRESENTED_IMAGE_SIMPLE}
        description="你暂无 GitOps 访问权限，请联系管理员授予范围权限。"
      />
    );
  }

  if (!unitId) {
    return <Alert type="error" showIcon message="缺少交付单元 ID，无法加载详情。" />;
  }

  const detail = detailQuery.data;
  const status = statusQuery.data;

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <Space wrap>
        <Link to="/gitops">
          <Button>返回概览</Button>
        </Link>
        <Typography.Title level={3} style={{ marginBottom: 0 }}>
          {detail?.name || unitId}
        </Typography.Title>
      </Space>

      <Space wrap>
        <Button
          type="primary"
          disabled={permissionChanged}
          onClick={() => setReleaseActionOpen(true)}
        >
          提交发布动作
        </Button>
        <Button
          onClick={() => {
            void detailQuery.refetch();
            void statusQuery.refetch();
          }}
        >
          刷新状态
        </Button>
      </Space>

      {permissionChanged ? (
        <Alert
          type="warning"
          showIcon
          message="权限已变更"
          description="当前账号可能缺少详情读取权限，请刷新页面或重新登录后重试。"
        />
      ) : null}

      {detailQuery.error && !isAuthorizationError(detailQuery.error) ? (
        <Alert
          type="error"
          showIcon
          message="交付单元详情加载失败"
          description={normalizeApiError(detailQuery.error, '交付单元详情加载失败，请稍后重试。')}
        />
      ) : null}

      {statusQuery.error && !isAuthorizationError(statusQuery.error) ? (
        <Alert
          type="error"
          showIcon
          message="状态聚合加载失败"
          description={normalizeApiError(statusQuery.error, '状态聚合加载失败，请稍后重试。')}
        />
      ) : null}

      <Card size="small" title="交付单元信息" loading={detailQuery.isLoading || detailQuery.isFetching}>
        <Descriptions size="small" column={2} bordered>
          <Descriptions.Item label="来源 ID">{detail?.sourceId || '-'}</Descriptions.Item>
          <Descriptions.Item label="来源路径">{detail?.sourcePath || '-'}</Descriptions.Item>
          <Descriptions.Item label="同步模式">{detail?.syncMode || '-'}</Descriptions.Item>
          <Descriptions.Item label="期望版本">{detail?.desiredRevision || '-'}</Descriptions.Item>
          <Descriptions.Item label="应用版本">{detail?.desiredAppVersion || '-'}</Descriptions.Item>
          <Descriptions.Item label="配置版本">{detail?.desiredConfigVersion || '-'}</Descriptions.Item>
        </Descriptions>
      </Card>

      <Card size="small" title="状态聚合" loading={statusQuery.isLoading || statusQuery.isFetching}>
        <Space wrap>
          <Tag color={statusColorMap[status?.deliveryStatus || 'unknown'] || 'default'}>
            交付：{status?.deliveryStatus || 'unknown'}
          </Tag>
          <Tag color={statusColorMap[status?.driftStatus || 'unknown'] || 'default'}>
            漂移：{status?.driftStatus || 'unknown'}
          </Tag>
          <Tag color={statusColorMap[status?.lastSyncResult || 'unknown'] || 'default'}>
            最近结果：{status?.lastSyncResult || 'unknown'}
          </Tag>
        </Space>
        <Table<EnvironmentStatusRecord>
          style={{ marginTop: 12 }}
          rowKey={(record, index) => `${record.environment || 'unknown'}-${index}`}
          columns={environmentStatusColumns}
          dataSource={status?.environments || []}
          pagination={false}
        />
      </Card>

      <DiffSummaryPanel unitId={unitId} />

      <RevisionHistoryPanel
        unitId={unitId}
        onRollback={(revision) => {
          setSelectedRevision(revision);
          setRollbackOpen(true);
        }}
      />

      <PromotionTimeline stages={detail?.environments || []} status={status} operation={latestOperation} />

      <EnvironmentStageEditor
        value={stages}
        targetGroups={targetGroupsQuery.data?.items || []}
        onChange={setStages}
      />

      <OverlaySummaryPanel overlays={detail?.overlays || []} />

      <ReleaseActionDrawer
        open={releaseActionOpen}
        unitId={unitId}
        unitName={detail?.name}
        onClose={() => setReleaseActionOpen(false)}
        onOperationChange={handleOperationChange}
      />

      <RollbackDialog
        open={rollbackOpen}
        unitId={unitId}
        revision={selectedRevision}
        onClose={() => setRollbackOpen(false)}
        onOperationChange={handleOperationChange}
      />
    </Space>
  );
};
