import { Descriptions, Drawer, Empty, Space, Typography } from 'antd';
import type { RestorePoint } from '@/services/backupRestore';
import { backupRestoreScopeSummary } from '@/services/backupRestore';
import { StatusTag } from '@/features/backup-restore-dr/components/status';

type RestorePointDetailDrawerProps = {
  open: boolean;
  restorePoint?: RestorePoint;
  onClose: () => void;
};

export const RestorePointDetailDrawer = ({
  open,
  restorePoint,
  onClose
}: RestorePointDetailDrawerProps) => (
  <Drawer title="恢复点详情" width={560} open={open} onClose={onClose} destroyOnClose>
    {!restorePoint ? (
      <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description="请选择一个恢复点查看详情。" />
    ) : (
      <Space direction="vertical" size={16} style={{ width: '100%' }}>
        <Descriptions column={1} size="small" bordered>
          <Descriptions.Item label="恢复点 ID">{restorePoint.id}</Descriptions.Item>
          <Descriptions.Item label="策略 ID">{restorePoint.policyId || '—'}</Descriptions.Item>
          <Descriptions.Item label="执行结果">
            <StatusTag value={restorePoint.result} />
          </Descriptions.Item>
          <Descriptions.Item label="备份耗时">
            {typeof restorePoint.durationSeconds === 'number'
              ? `${restorePoint.durationSeconds} 秒`
              : '—'}
          </Descriptions.Item>
          <Descriptions.Item label="开始时间">{restorePoint.backupStartedAt || '—'}</Descriptions.Item>
          <Descriptions.Item label="完成时间">
            {restorePoint.backupCompletedAt || '—'}
          </Descriptions.Item>
          <Descriptions.Item label="过期时间">{restorePoint.expiresAt || '—'}</Descriptions.Item>
          <Descriptions.Item label="存储引用">{restorePoint.storageRef || '—'}</Descriptions.Item>
          <Descriptions.Item label="范围摘要">
            {backupRestoreScopeSummary(restorePoint.scopeSnapshot)}
          </Descriptions.Item>
        </Descriptions>
        <div>
          <Typography.Text strong>一致性说明</Typography.Text>
          <Typography.Paragraph style={{ marginBottom: 0, marginTop: 8 }}>
            {restorePoint.consistencySummary || '暂无一致性说明。'}
          </Typography.Paragraph>
        </div>
        <div>
          <Typography.Text strong>失败原因</Typography.Text>
          <Typography.Paragraph style={{ marginBottom: 0, marginTop: 8 }}>
            {restorePoint.failureReason || '无'}
          </Typography.Paragraph>
        </div>
      </Space>
    )}
  </Drawer>
);
