import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Alert, Card, List, Space } from 'antd';
import { PageHeader } from '@/features/sre-scale/components/PageHeader';
import { PermissionDenied } from '@/features/sre-scale/components/PermissionDenied';
import { RunbookDrawer } from '@/features/sre-scale/components/RunbookDrawer';
import { useSREPermissions } from '@/features/sre-scale/hooks/permissions';
import { normalizeApiError } from '@/services/api/client';
import { listRunbooks, type RunbookArticle } from '@/services/sreScale';

export const RunbookCenterPage = () => {
  const permissions = useSREPermissions();
  const [selected, setSelected] = useState<RunbookArticle>();
  const runbookQuery = useQuery({ queryKey: ['sreScale', 'runbooks'], queryFn: () => listRunbooks() });

  if (!permissions.canRead) {
    return <PermissionDenied description="你暂无运行手册访问权限。" />;
  }

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <PageHeader title="运行手册中心" description="查看平台异常处理、恢复步骤和验证建议。" />
      {runbookQuery.error ? <Alert type="error" showIcon message={normalizeApiError(runbookQuery.error, '运行手册加载失败')} /> : null}
      <Card size="small" title="手册列表">
        <List
          dataSource={runbookQuery.data?.items || []}
          renderItem={(item) => (
            <List.Item onClick={() => setSelected(item)} style={{ cursor: 'pointer' }}>
              {item.title} / {item.riskLevel || 'unknown'}
            </List.Item>
          )}
        />
      </Card>
      <RunbookDrawer open={Boolean(selected)} runbook={selected} onClose={() => setSelected(undefined)} />
    </Space>
  );
};
