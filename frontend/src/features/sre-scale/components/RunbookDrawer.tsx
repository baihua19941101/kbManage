import { Drawer, Descriptions } from 'antd';
import type { RunbookArticle } from '@/services/sreScale';

export const RunbookDrawer = ({
  open,
  runbook,
  onClose
}: {
  open: boolean;
  runbook?: RunbookArticle;
  onClose: () => void;
}) => (
  <Drawer title="运行手册详情" open={open} width={480} onClose={onClose} destroyOnClose>
    <Descriptions column={1} bordered size="small">
      <Descriptions.Item label="标题">{runbook?.title || '—'}</Descriptions.Item>
      <Descriptions.Item label="场景">{runbook?.scenarioType || '—'}</Descriptions.Item>
      <Descriptions.Item label="风险等级">{runbook?.riskLevel || '—'}</Descriptions.Item>
      <Descriptions.Item label="检查摘要">{runbook?.checklistSummary || '—'}</Descriptions.Item>
      <Descriptions.Item label="恢复步骤">{runbook?.recoverySteps || '—'}</Descriptions.Item>
      <Descriptions.Item label="验证摘要">{runbook?.verificationSummary || '—'}</Descriptions.Item>
    </Descriptions>
  </Drawer>
);
