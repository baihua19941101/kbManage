import { useMemo, useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Alert, Card, Empty, Select, Space, Typography } from 'antd';
import {
  canReadCompliance,
  useAuthStore
} from '@/features/auth/store';
import { RecheckTaskTable } from '@/features/compliance-hardening/components/RecheckTaskTable';
import { normalizeApiError } from '@/services/api/client';
import { listRecheckTasks, type RecheckStatus } from '@/services/compliance';

const statusOptions: Array<{ label: string; value: RecheckStatus | '' }> = [
  { label: '全部状态', value: '' },
  { label: '待执行', value: 'pending' },
  { label: '执行中', value: 'running' },
  { label: '通过', value: 'passed' },
  { label: '失败', value: 'failed' },
  { label: '取消', value: 'canceled' }
];

export const RecheckCenterPage = () => {
  const user = useAuthStore((state) => state.user);
  const canRead = canReadCompliance(user);
  const [status, setStatus] = useState<RecheckStatus | ''>('');

  const queryInput = useMemo(() => ({ status: status || undefined }), [status]);
  const rechecksQuery = useQuery({
    queryKey: ['compliance', 'rechecks', queryInput],
    enabled: canRead,
    queryFn: () => listRecheckTasks(queryInput)
  });

  if (!canRead) {
    return <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description="你暂无复检中心访问权限。" />;
  }

  const tasks = rechecksQuery.data?.items || [];

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <div>
        <Typography.Title level={3} style={{ marginBottom: 8 }}>
          复检中心
        </Typography.Title>
        <Typography.Text type="secondary">
          查看整改完成或例外到期后触发的复检任务及结果回写状态。
        </Typography.Text>
      </div>

      {rechecksQuery.error ? (
        <Alert type="error" showIcon message="复检任务加载失败" description={normalizeApiError(rechecksQuery.error, '复检任务加载失败，请稍后重试。')} />
      ) : null}

      <Card size="small" title="筛选">
        <Select style={{ width: 180 }} options={statusOptions} value={status} onChange={setStatus} />
      </Card>

      <Card size="small" title={`复检任务（${tasks.length}）`}>
        <RecheckTaskTable tasks={tasks} loading={rechecksQuery.isLoading || rechecksQuery.isFetching} />
      </Card>
    </Space>
  );
};
