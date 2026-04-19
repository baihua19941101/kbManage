import { Descriptions, Drawer, Typography } from 'antd';
import type { SessionRecord } from '@/services/identityTenancy';

type SessionRiskDrawerProps = {
  open: boolean;
  record?: SessionRecord;
  onClose: () => void;
};

export const SessionRiskDrawer = ({ open, record, onClose }: SessionRiskDrawerProps) => (
  <Drawer title="会话风险详情" width={520} open={open} onClose={onClose} destroyOnClose>
    <Descriptions column={1} size="small" bordered>
      <Descriptions.Item label="用户">{record?.username || record?.userId || '未标记用户'}</Descriptions.Item>
      <Descriptions.Item label="身份来源">{record?.identitySourceId || '本地'}</Descriptions.Item>
      <Descriptions.Item label="登录方式">{record?.loginMethod || '未知'}</Descriptions.Item>
      <Descriptions.Item label="状态">{record?.status || '未知'}</Descriptions.Item>
      <Descriptions.Item label="风险等级">{record?.riskLevel || '未知'}</Descriptions.Item>
      <Descriptions.Item label="最后活跃时间">{record?.lastSeenAt || '未知'}</Descriptions.Item>
    </Descriptions>
    <Typography.Paragraph type="secondary" style={{ marginTop: 16 }}>
      {record?.riskSummary || '当前未返回更详细的风险摘要。'}
    </Typography.Paragraph>
  </Drawer>
);
