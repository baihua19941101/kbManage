import { Descriptions, Drawer, Empty, List, Space, Typography } from 'antd';
import type { DRDrillReport } from '@/services/backupRestore';

type DRDrillReportDrawerProps = {
  open: boolean;
  report?: DRDrillReport;
  onClose: () => void;
};

export const DRDrillReportDrawer = ({
  open,
  report,
  onClose
}: DRDrillReportDrawerProps) => (
  <Drawer title="演练报告" width={560} open={open} onClose={onClose} destroyOnClose>
    {!report ? (
      <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description="请先生成演练报告。" />
    ) : (
      <Space direction="vertical" size={16} style={{ width: '100%' }}>
        <Descriptions bordered column={1} size="small">
          <Descriptions.Item label="报告 ID">{report.id}</Descriptions.Item>
          <Descriptions.Item label="演练记录">{report.drillRecordId || '—'}</Descriptions.Item>
          <Descriptions.Item label="目标评估">{report.goalAssessment || '—'}</Descriptions.Item>
          <Descriptions.Item label="差距摘要">{report.gapSummary || '—'}</Descriptions.Item>
          <Descriptions.Item label="发布时间">{report.publishedAt || '—'}</Descriptions.Item>
          <Descriptions.Item label="发布人">{report.publishedBy || '—'}</Descriptions.Item>
        </Descriptions>
        <div>
          <Typography.Text strong>问题项</Typography.Text>
          <List
            size="small"
            dataSource={report.issuesFound}
            locale={{ emptyText: '无' }}
            renderItem={(item) => <List.Item>{item}</List.Item>}
          />
        </div>
        <div>
          <Typography.Text strong>改进建议</Typography.Text>
          <List
            size="small"
            dataSource={report.improvementActions}
            locale={{ emptyText: '无' }}
            renderItem={(item) => <List.Item>{item}</List.Item>}
          />
        </div>
      </Space>
    )}
  </Drawer>
);
