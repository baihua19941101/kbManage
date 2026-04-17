import { useMemo, useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Alert, Button, Card, Empty, Input, Select, Space, Typography } from 'antd';
import {
  canManageComplianceRemediation,
  canReadCompliance,
  useAuthStore
} from '@/features/auth/store';
import { RemediationTaskDrawer } from '@/features/compliance-hardening/components/RemediationTaskDrawer';
import { RemediationTaskTable } from '@/features/compliance-hardening/components/RemediationTaskTable';
import { useComplianceAction } from '@/features/compliance-hardening/hooks/useComplianceAction';
import { isAuthorizationError, normalizeApiError } from '@/services/api/client';
import { listRemediationTasks, type RemediationTask, type RemediationTaskStatus } from '@/services/compliance';

const statusOptions: Array<{ label: string; value: RemediationTaskStatus | '' }> = [
  { label: '全部状态', value: '' },
  { label: '待处理', value: 'todo' },
  { label: '处理中', value: 'in_progress' },
  { label: '阻塞', value: 'blocked' },
  { label: '完成', value: 'done' },
  { label: '取消', value: 'canceled' }
];

export const RemediationQueuePage = () => {
  const user = useAuthStore((state) => state.user);
  const canRead = canReadCompliance(user);
  const canManage = canManageComplianceRemediation(user);
  const [status, setStatus] = useState<RemediationTaskStatus | ''>('');
  const [owner, setOwner] = useState('');
  const [selectedTask, setSelectedTask] = useState<RemediationTask>();
  const [drawerOpen, setDrawerOpen] = useState(false);
  const { updateRemediationMutation } = useComplianceAction();

  const queryInput = useMemo(
    () => ({ status: status || undefined, owner: owner.trim() || undefined }),
    [owner, status]
  );

  const tasksQuery = useQuery({
    queryKey: ['compliance', 'remediation', queryInput],
    enabled: canRead,
    queryFn: () => listRemediationTasks(queryInput)
  });

  if (!canRead) {
    return <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description="你暂无整改工作台访问权限。" />;
  }

  const tasks = tasksQuery.data?.items || [];
  const permissionChanged = isAuthorizationError(tasksQuery.error);

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <div>
        <Typography.Title level={3} style={{ marginBottom: 8 }}>
          整改任务工作台
        </Typography.Title>
        <Typography.Text type="secondary">
          跟踪失败项整改任务的负责人、优先级、状态和到期时间。
        </Typography.Text>
      </div>

      {!canManage ? (
        <Alert type="info" showIcon message="当前为只读模式" description="你可以查看整改任务，但无法更新任务状态。" />
      ) : null}

      {permissionChanged ? (
        <Alert type="warning" showIcon message="权限已变更" description="当前会话已失去整改任务权限，请刷新或重新登录。" />
      ) : null}

      {tasksQuery.error && !permissionChanged ? (
        <Alert type="error" showIcon message="整改任务加载失败" description={normalizeApiError(tasksQuery.error, '整改任务加载失败，请稍后重试。')} />
      ) : null}

      <Card size="small" title="筛选">
        <Space wrap>
          <Select style={{ width: 180 }} options={statusOptions} value={status} onChange={setStatus} />
          <Input style={{ width: 220 }} placeholder="负责人" value={owner} onChange={(event) => setOwner(event.target.value)} />
          <Button
            disabled={!canManage || permissionChanged || !selectedTask}
            onClick={() => setDrawerOpen(true)}
          >
            更新选中任务
          </Button>
        </Space>
      </Card>

      <Card size="small" title={`整改任务（${tasks.length}）`}>
        <RemediationTaskTable
          tasks={tasks}
          loading={tasksQuery.isLoading || tasksQuery.isFetching}
          readonly={!canManage || permissionChanged}
          onEdit={(task) => {
            setSelectedTask(task);
            setDrawerOpen(true);
          }}
        />
      </Card>

      <RemediationTaskDrawer
        open={drawerOpen}
        task={selectedTask}
        readonly={!canManage || permissionChanged}
        loading={updateRemediationMutation.isPending}
        onClose={() => setDrawerOpen(false)}
        onCreate={() => undefined}
        onUpdate={(taskId, payload) => updateRemediationMutation.mutate({ taskId, payload })}
      />
    </Space>
  );
};
