import { Alert, Button, Descriptions, Drawer, Space, Typography } from 'antd';
import type { ClusterTemplate, ValidationResult } from '@/services/clusterLifecycle';

type Props = {
  open: boolean;
  template?: ClusterTemplate | null;
  validation?: ValidationResult | null;
  loading?: boolean;
  onClose: () => void;
};

export const ClusterTemplateDrawer = ({ open, template, validation, loading, onClose }: Props) => (
  <Drawer open={open} width={520} title="模板详情与校验结论" onClose={onClose} destroyOnClose>
    {!template ? (
      <Alert type="info" showIcon message="请选择一个模板查看详情。" />
    ) : (
      <Space direction="vertical" size={16} style={{ width: '100%' }}>
        <Descriptions column={1} size="small">
          <Descriptions.Item label="模板名称">{template.name}</Descriptions.Item>
          <Descriptions.Item label="基础设施">{template.infrastructureType || '—'}</Descriptions.Item>
          <Descriptions.Item label="驱动">{template.driverKey || '—'}</Descriptions.Item>
          <Descriptions.Item label="版本范围">{template.driverVersionRange || '—'}</Descriptions.Item>
          <Descriptions.Item label="依赖能力">
            {(template.requiredCapabilities || []).join(' / ') || '—'}
          </Descriptions.Item>
        </Descriptions>

        {validation ? (
          <Alert
            type={validation.blockers?.length ? 'error' : validation.warnings?.length ? 'warning' : 'success'}
            showIcon
            message={`校验状态：${validation.overallStatus || '未知'}`}
            description={
              <Space direction="vertical" size={4}>
                <Typography.Text>阻断项：{(validation.blockers || []).join('；') || '无'}</Typography.Text>
                <Typography.Text>风险提示：{(validation.warnings || []).join('；') || '无'}</Typography.Text>
                <Typography.Text>通过项：{(validation.passedChecks || []).join('；') || '无'}</Typography.Text>
              </Space>
            }
          />
        ) : (
          <Alert type="info" showIcon message="尚未执行创建前校验。" />
        )}

        <Button onClick={onClose} loading={loading}>
          关闭
        </Button>
      </Space>
    )}
  </Drawer>
);
