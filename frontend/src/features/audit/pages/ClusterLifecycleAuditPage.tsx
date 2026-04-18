import { useQuery } from '@tanstack/react-query';
import { Alert, Card, List, Space, Typography } from 'antd';
import { canReadClusterLifecycleAudit, useAuthStore } from '@/features/auth/store';
import { normalizeApiError } from '@/services/api/client';
import { listLifecycleAuditEvents } from '@/services/clusterLifecycle';

export const ClusterLifecycleAuditPage = () => {
  const user = useAuthStore((state) => state.user);
  const canRead = canReadClusterLifecycleAudit(user);
  const auditQuery = useQuery({
    queryKey: ['cluster-lifecycle', 'audit-events'],
    queryFn: () => listLifecycleAuditEvents({})
  });

  if (!canRead) {
    return <Alert type="info" showIcon message="你暂无集群生命周期审计访问权限。" />;
  }

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <Typography.Title level={3} style={{ marginBottom: 0 }}>
        集群生命周期审计
      </Typography.Title>
      {auditQuery.error ? (
        <Alert type="error" showIcon message={normalizeApiError(auditQuery.error, '生命周期审计加载失败')} />
      ) : null}
      <Card>
        <List
          loading={auditQuery.isLoading || auditQuery.isFetching}
          dataSource={auditQuery.data?.items ?? []}
          renderItem={(item) => (
            <List.Item>
              <Space direction="vertical" size={0}>
                <Typography.Text>{item.action || '未命名动作'}</Typography.Text>
                <Typography.Text type="secondary">
                  目标 {item.targetType || '未标记类型'} / {item.targetRef || '未标记对象'} / 结果 {item.outcome || '未知'}
                </Typography.Text>
              </Space>
            </List.Item>
          )}
        />
      </Card>
    </Space>
  );
};
