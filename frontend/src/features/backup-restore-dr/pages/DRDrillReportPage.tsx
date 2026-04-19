import { useState } from 'react';
import { Alert, Button, Card, Input, Space } from 'antd';
import { useParams } from 'react-router-dom';
import { DRDrillReportDrawer } from '@/features/backup-restore-dr/components/DRDrillReportDrawer';
import { PageHeader } from '@/features/backup-restore-dr/components/PageHeader';
import { PermissionDenied } from '@/features/backup-restore-dr/components/PermissionDenied';
import { useDrillAction } from '@/features/backup-restore-dr/hooks/useDrillAction';
import { useBackupRestorePermissions } from '@/features/backup-restore-dr/hooks/permissions';
import { normalizeApiError } from '@/services/api/client';

export const DRDrillReportPage = () => {
  const permissions = useBackupRestorePermissions();
  const { reportMutation } = useDrillAction();
  const params = useParams<{ recordId?: string }>();
  const [recordId, setRecordId] = useState(params.recordId || 'record-001');
  const [drawerOpen, setDrawerOpen] = useState(false);

  if (!permissions.canRead) {
    return <PermissionDenied description="你暂无灾备演练报告访问权限。" />;
  }

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <PageHeader
        title="灾备演练报告"
        description="根据演练记录生成目标达成情况、差距说明和改进建议。"
      />

      {!permissions.canDrill ? (
        <Alert type="info" showIcon message="你当前只有查看权限，报告生成动作已被禁用。" />
      ) : null}
      {reportMutation.error ? (
        <Alert
          type="error"
          showIcon
          message="演练报告生成失败"
          description={normalizeApiError(reportMutation.error, '演练报告生成失败，请稍后重试。')}
        />
      ) : null}

      <Card size="small" title="生成报告">
        <Space wrap>
          <Input
            value={recordId}
            onChange={(event) => setRecordId(event.target.value)}
            placeholder="输入演练记录 ID"
            style={{ width: 260 }}
          />
          <Button
            type="primary"
            disabled={!permissions.canDrill}
            loading={reportMutation.isPending}
            onClick={() =>
              reportMutation.mutate(recordId.trim(), {
                onSuccess: () => setDrawerOpen(true)
              })
            }
          >
            生成报告
          </Button>
        </Space>
      </Card>

      <DRDrillReportDrawer
        open={drawerOpen}
        report={reportMutation.data}
        onClose={() => setDrawerOpen(false)}
      />
    </Space>
  );
};
