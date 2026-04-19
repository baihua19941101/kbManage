import { useQuery } from '@tanstack/react-query';
import { Alert, Card, List, Space, Typography } from 'antd';
import { canReadIdentityTenancyAudit, useAuthStore } from '@/features/auth/store';
import { normalizeApiError } from '@/services/api/client';
import { listIdentityGovernanceAuditEvents } from '@/services/identityTenancy';

export const IdentityGovernanceAuditPage = () => {
  const user = useAuthStore((state) => state.user);
  const canRead = canReadIdentityTenancyAudit(user);
  const auditQuery = useQuery({
    queryKey: ['identity-tenancy', 'audit-events'],
    queryFn: () => listIdentityGovernanceAuditEvents({})
  });

  if (!canRead) {
    return <Alert type="info" showIcon message="你暂无身份治理审计访问权限。" />;
  }

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <Typography.Title level={3} style={{ marginBottom: 0 }}>
        身份治理审计
      </Typography.Title>
      {auditQuery.error ? (
        <Alert
          type="error"
          showIcon
          message={normalizeApiError(auditQuery.error, '身份治理审计加载失败')}
        />
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
                  目标 {item.targetType || '未标记类型'} / {item.targetRef || '未标记对象'} / 结果{' '}
                  {item.outcome || '未知'}
                </Typography.Text>
              </Space>
            </List.Item>
          )}
        />
      </Card>
    </Space>
  );
};
