import { Descriptions, Drawer, Empty, Typography } from 'antd';
import type { InstallationRecord } from '@/services/platformMarketplace';

type UpgradeEntryDrawerProps = {
  open: boolean;
  record?: InstallationRecord;
  onClose: () => void;
};

export const UpgradeEntryDrawer = ({ open, record, onClose }: UpgradeEntryDrawerProps) => (
  <Drawer title="升级入口" open={open} width={420} onClose={onClose} destroyOnClose>
    {record ? (
      <>
        <Descriptions column={1} size="small" bordered>
          <Descriptions.Item label="模板">{record.templateName || '—'}</Descriptions.Item>
          <Descriptions.Item label="当前版本">{record.currentVersion || '—'}</Descriptions.Item>
          <Descriptions.Item label="目标版本">{record.latestVersion || '—'}</Descriptions.Item>
          <Descriptions.Item label="状态">{record.status || '—'}</Descriptions.Item>
          <Descriptions.Item label="下线状态">{record.offlineState || '—'}</Descriptions.Item>
        </Descriptions>
        <Typography.Paragraph style={{ marginTop: 16 }}>
          {record.changeSummary || '当前记录未提供版本差异摘要，主线程接入后端后可补充精细化升级说明。'}
        </Typography.Paragraph>
      </>
    ) : (
      <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description="请选择一条安装记录查看升级信息。" />
    )}
  </Drawer>
);
