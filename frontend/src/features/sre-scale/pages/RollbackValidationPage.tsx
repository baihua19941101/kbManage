import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Alert, Button, Card, Select, Space } from 'antd';
import { PageHeader } from '@/features/sre-scale/components/PageHeader';
import { PermissionDenied } from '@/features/sre-scale/components/PermissionDenied';
import { RollbackValidationDrawer } from '@/features/sre-scale/components/RollbackValidationDrawer';
import { useSREPermissions } from '@/features/sre-scale/hooks/permissions';
import { useUpgradeAction } from '@/features/sre-scale/hooks/useUpgradeAction';
import { normalizeApiError } from '@/services/api/client';
import { listUpgradePlans } from '@/services/sreScale';

export const RollbackValidationPage = () => {
  const permissions = useSREPermissions();
  const action = useUpgradeAction();
  const [drawerOpen, setDrawerOpen] = useState(false);
  const [upgradeId, setUpgradeId] = useState<string>();
  const plansQuery = useQuery({ queryKey: ['sreScale', 'upgrades', 'rollback'], queryFn: () => listUpgradePlans() });

  if (!permissions.canRead) {
    return <PermissionDenied description="你暂无回退验证访问权限。" />;
  }

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <PageHeader
        title="回退验证"
        description="登记升级回退可行性和剩余风险。"
        actions={
          <Button type="primary" disabled={!permissions.canManageUpgrade || !upgradeId} onClick={() => setDrawerOpen(true)}>
            新增回退验证
          </Button>
        }
      />
      {plansQuery.error ? <Alert type="error" showIcon message={normalizeApiError(plansQuery.error, '升级计划加载失败')} /> : null}
      <Card size="small" title="选择升级计划">
        <Select
          style={{ width: 360 }}
          placeholder="选择升级计划"
          options={(plansQuery.data?.items || []).map((item) => ({ label: item.name, value: item.id }))}
          value={upgradeId}
          onChange={(value) => setUpgradeId(value)}
        />
      </Card>
      <RollbackValidationDrawer
        open={drawerOpen}
        submitting={action.rollbackMutation.isPending}
        onClose={() => setDrawerOpen(false)}
        onSubmit={(payload) => {
          if (!upgradeId) return;
          action.rollbackMutation.mutate(
            { upgradeId, payload },
            {
              onSuccess: () => setDrawerOpen(false)
            }
          );
        }}
      />
    </Space>
  );
};
