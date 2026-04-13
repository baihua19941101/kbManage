import { useMemo, useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Alert, Button, Card, Descriptions, Empty, Space, Typography } from 'antd';
import { useLocation, useNavigate } from 'react-router-dom';
import {
  canAccessWorkloadOpsTerminal,
  canBatchWorkloadOps,
  canExecuteWorkloadOps,
  canReadWorkloadOps,
  canRollbackWorkloadOps,
  useAuthStore
} from '@/features/auth/store';
import { ActionConfirmDrawer } from '@/features/workload-ops/components/ActionConfirmDrawer';
import { BatchOperationDrawer } from '@/features/workload-ops/components/BatchOperationDrawer';
import { InstanceListPanel } from '@/features/workload-ops/components/InstanceListPanel';
import { RevisionHistoryPanel } from '@/features/workload-ops/components/RevisionHistoryPanel';
import { RollbackDialog } from '@/features/workload-ops/components/RollbackDialog';
import { TerminalSessionDrawer } from '@/features/workload-ops/components/TerminalSessionDrawer';
import type { ReleaseRevisionDTO, WorkloadInstanceDTO } from '@/services/api/types';
import { ApiError, buildScopeQueryKey, normalizeApiError } from '@/services/api/client';
import { queryKeys } from '@/app/queryClient';
import {
  getWorkloadOpsContext,
  listWorkloadOpsInstances,
  listWorkloadOpsRevisions
} from '@/services/workloadOps';

const isAuthorizationError = (error: unknown): boolean =>
  error instanceof ApiError && (error.status === 401 || error.status === 403);

export const WorkloadOperationsPage = () => {
  const navigate = useNavigate();
  const user = useAuthStore((state) => state.user);
  const canRead = canReadWorkloadOps(user);
  const canExecute = canExecuteWorkloadOps(user);
  const canBatch = canBatchWorkloadOps(user);
  const canRollback = canRollbackWorkloadOps(user);
  const canUseTerminal = canAccessWorkloadOpsTerminal(user);
  const { search } = useLocation();
  const [targetInstance, setTargetInstance] = useState<WorkloadInstanceDTO>();
  const [terminalOpen, setTerminalOpen] = useState(false);
  const [actionOpen, setActionOpen] = useState(false);
  const [batchOpen, setBatchOpen] = useState(false);
  const [rollbackOpen, setRollbackOpen] = useState(false);
  const [selectedRevision, setSelectedRevision] = useState<ReleaseRevisionDTO>();
  const params = useMemo(() => new URLSearchParams(search), [search]);

  const clusterId = Number(params.get('clusterId') || 1);
  const namespace = params.get('namespace') || 'default';
  const resourceKind = (params.get('resourceKind') || 'Deployment') as
    | 'Deployment'
    | 'StatefulSet'
    | 'DaemonSet';
  const resourceName = params.get('resourceName') || 'demo-app';

  const scopeKey = buildScopeQueryKey([clusterId, namespace, resourceKind, resourceName]);

  const contextQuery = useQuery({
    queryKey: queryKeys.workloadOps.context(scopeKey),
    enabled: canRead,
    queryFn: () =>
      getWorkloadOpsContext({
        clusterId,
        namespace,
        resourceKind,
        resourceName
      })
  });

  const instancesQuery = useQuery({
    queryKey: queryKeys.workloadOps.instances(scopeKey),
    enabled: canRead,
    queryFn: () =>
      listWorkloadOpsInstances({
        clusterId,
        namespace,
        resourceKind,
        resourceName
      })
  });
  const revisionsQuery = useQuery({
    queryKey: queryKeys.workloadOps.revisions(scopeKey),
    enabled: canRead,
    queryFn: () =>
      listWorkloadOpsRevisions({
        clusterId,
        namespace,
        resourceKind,
        resourceName
      })
  });

  const context = contextQuery.data;
  const items = instancesQuery.data?.items ?? [];
  const revisions = revisionsQuery.data?.items ?? [];
  const queryPermissionError =
    isAuthorizationError(contextQuery.error) ||
    isAuthorizationError(instancesQuery.error) ||
    isAuthorizationError(revisionsQuery.error);

  if (!canRead) {
    return (
      <Empty description="你暂无工作负载运维访问权限，请联系管理员授予工作空间/项目范围。" />
    );
  }

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <Typography.Title level={3} style={{ marginBottom: 0 }}>
        工作负载运维
      </Typography.Title>
      {contextQuery.error ? (
        <Alert type="error" showIcon message={normalizeApiError(contextQuery.error, '运维上下文加载失败')} />
      ) : null}
      {instancesQuery.error ? (
        <Alert type="error" showIcon message={normalizeApiError(instancesQuery.error, '实例列表加载失败')} />
      ) : null}
      {queryPermissionError ? (
        <Alert type="warning" showIcon message="权限已变更，当前动作入口已锁定。" />
      ) : null}
      <Space wrap>
        <Button type="primary" disabled={!canExecute || queryPermissionError} onClick={() => setActionOpen(true)}>
          提交重启动作
        </Button>
        <Button disabled={!canBatch || queryPermissionError} onClick={() => setBatchOpen(true)}>
          提交批量重启
        </Button>
      </Space>
      {!canExecute ? (
        <Alert
          type="info"
          showIcon
          message="当前为只读模式"
          description="你可以查看工作负载上下文、实例和发布历史，但无法执行重启、批量动作、终端和回滚。"
        />
      ) : null}
      {context ? (
        <Card size="small" title="资源上下文">
          <Descriptions column={1} size="small" bordered>
            <Descriptions.Item label="集群 ID">{context.clusterId}</Descriptions.Item>
            <Descriptions.Item label="命名空间">{context.namespace}</Descriptions.Item>
            <Descriptions.Item label="资源类型">{context.resourceKind}</Descriptions.Item>
            <Descriptions.Item label="资源名称">{context.resourceName}</Descriptions.Item>
            <Descriptions.Item label="健康状态">{context.healthStatus}</Descriptions.Item>
            <Descriptions.Item label="发布状态">{context.rolloutStatus}</Descriptions.Item>
          </Descriptions>
        </Card>
      ) : (
        <Empty description="暂无运维上下文数据" />
      )}
      <InstanceListPanel
        loading={instancesQuery.isLoading || instancesQuery.isFetching}
        items={items}
        terminalDisabled={!canUseTerminal || queryPermissionError}
        onOpenTerminal={(item) => {
          setTargetInstance(item);
          setTerminalOpen(true);
        }}
      />
      <RevisionHistoryPanel
        loading={revisionsQuery.isLoading || revisionsQuery.isFetching}
        items={revisions}
        onRollback={(item) => {
          if (!canRollback || queryPermissionError) {
            return;
          }
          setSelectedRevision(item);
          setRollbackOpen(true);
        }}
      />
      <TerminalSessionDrawer
        open={terminalOpen}
        clusterId={clusterId}
        namespace={namespace}
        workloadKind={resourceKind}
        workloadName={resourceName}
        target={targetInstance}
        onClose={() => setTerminalOpen(false)}
      />
      <ActionConfirmDrawer
        open={actionOpen && canExecute && !queryPermissionError}
        payload={{
          clusterId,
          namespace,
          resourceKind,
          resourceName,
          actionType: 'restart',
          riskConfirmed: true
        }}
        onClose={() => setActionOpen(false)}
      />
      <BatchOperationDrawer
        open={batchOpen && canBatch && !queryPermissionError}
        payload={{
          actionType: 'restart',
          riskConfirmed: true,
          targets: [
            { clusterId, namespace, resourceKind, resourceName },
            { clusterId, namespace, resourceKind, resourceName: `${resourceName}-canary` }
          ]
        }}
        onClose={() => setBatchOpen(false)}
        onSubmitted={(batchId) => {
          void navigate(`/workload-ops/batches?batchId=${batchId}`);
        }}
      />
      <RollbackDialog
        open={rollbackOpen && canRollback && !queryPermissionError}
        clusterId={clusterId}
        namespace={namespace}
        resourceKind={resourceKind}
        resourceName={resourceName}
        revision={selectedRevision}
        onClose={() => setRollbackOpen(false)}
      />
    </Space>
  );
};
