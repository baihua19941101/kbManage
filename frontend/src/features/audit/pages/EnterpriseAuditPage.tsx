import { useQuery } from '@tanstack/react-query';
import { Alert, Card, List, Space, Typography } from 'antd';
import { canReadEnterpriseAudit, useAuthStore } from '@/features/auth/store';
import { normalizeApiError } from '@/services/api/client';
import { listEnterpriseAuditEvents } from '@/services/enterprisePolish';

export const EnterpriseAuditPage = () => {
  const user = useAuthStore((state) => state.user);
  const canRead = canReadEnterpriseAudit(user);
  const auditQuery = useQuery({
    queryKey: ['enterprisePolish', 'audit'],
    queryFn: () => listEnterpriseAuditEvents({})
  });
  if (!canRead) {
    return <Alert type="info" showIcon message="你暂无企业治理审计访问权限。" />;
  }
  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <Typography.Title level={3} style={{ marginBottom: 0 }}>
        企业治理审计
      </Typography.Title>
      {auditQuery.error ? <Alert type="error" showIcon message={normalizeApiError(auditQuery.error, '企业治理审计加载失败')} /> : null}
      <Card>
        <List
          loading={auditQuery.isLoading || auditQuery.isFetching}
          dataSource={auditQuery.data?.items ?? []}
          renderItem={(item) => (
            <List.Item>
              {item.action || '未命名动作'} / {item.targetType || '未标记类型'} / {item.targetRef || '未标记对象'} / {item.outcome || '未知'}
            </List.Item>
          )}
        />
      </Card>
    </Space>
  );
};
