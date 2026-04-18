import { Alert, Button, Card, Descriptions, Typography } from 'antd';
import type { RegistrationBundle } from '@/services/clusterLifecycle';

type Props = {
  bundle?: RegistrationBundle | null;
  onCreateGuide?: () => void;
  loading?: boolean;
};

export const RegistrationGuideCard = ({ bundle, onCreateGuide, loading }: Props) => (
  <Card
    size="small"
    title="注册引导"
    extra={
      onCreateGuide ? (
        <Button type="primary" onClick={onCreateGuide} loading={loading}>
          生成注册指引
        </Button>
      ) : null
    }
  >
    {!bundle ? (
      <Alert
        type="info"
        showIcon
        message="尚未生成注册指引"
        description="填写集群基础信息后生成注册命令，交给平台工程团队在目标集群侧执行。"
      />
    ) : (
      <Descriptions column={1} size="small">
        <Descriptions.Item label="集群 ID">{bundle.clusterId}</Descriptions.Item>
        <Descriptions.Item label="注册令牌">
          <Typography.Text code>{bundle.registrationToken || '—'}</Typography.Text>
        </Descriptions.Item>
        <Descriptions.Item label="过期时间">{bundle.expiresAt || '—'}</Descriptions.Item>
        <Descriptions.Item label="接入命令">
          <Typography.Paragraph copyable code style={{ marginBottom: 0 }}>
            {bundle.commandSnippet || '—'}
          </Typography.Paragraph>
        </Descriptions.Item>
      </Descriptions>
    )}
  </Card>
);
