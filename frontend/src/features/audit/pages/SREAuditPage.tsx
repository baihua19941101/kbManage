import { useQuery } from '@tanstack/react-query';
import { Alert, Card, List, Space, Typography } from 'antd';
import { canReadSREAudit, useAuthStore } from '@/features/auth/store';
import { normalizeApiError } from '@/services/api/client';
import { listSREAuditEvents } from '@/services/sreScale';

export const SREAuditPage = () => {
  const user = useAuthStore((state) => state.user);
  const canRead = canReadSREAudit(user);
  const auditQuery = useQuery({
    queryKey: ['sreScale', 'audit-events'],
    queryFn: () => listSREAuditEvents({})
  });

  if (!canRead) {
    return <Alert type="info" showIcon message="你暂无 SRE 审计访问权限。" />;
  }

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <Typography.Title level={3} style={{ marginBottom: 0 }}>
        平台 SRE 审计
      </Typography.Title>
      {auditQuery.error ? (
        <Alert type="error" showIcon message={normalizeApiError(auditQuery.error, 'SRE 审计加载失败')} />
      ) : null}
      <Card>
        <List
          loading={auditQuery.isLoading || auditQuery.isFetching}
          dataSource={auditQuery.data?.items ?? []}
          renderItem={(item) => (
            <List.Item>
              {item.action || '未命名动作'} / {item.targetType || '未标记类型'} / {item.outcome || '未知'}
            </List.Item>
          )}
        />
      </Card>
    </Space>
  );
};
