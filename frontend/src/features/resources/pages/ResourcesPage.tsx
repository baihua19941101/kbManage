import { useMemo } from 'react';
import { Button, Card, Space, Typography } from 'antd';
import { useLocation, useNavigate } from 'react-router-dom';
import { ResourceListPage } from '@/features/resources/pages/ResourceListPage';

const normalizeValue = (value: string | null): string | undefined => {
  const normalized = value?.trim();
  return normalized && normalized.length > 0 ? normalized : undefined;
};

export const ResourcesPage = () => {
  const navigate = useNavigate();
  const { search } = useLocation();

  const resourceKeyword = useMemo(() => {
    const params = new URLSearchParams(search);
    return (
      normalizeValue(params.get('resourceName')) ||
      normalizeValue(params.get('name')) ||
      normalizeValue(params.get('keyword'))
    );
  }, [search]);

  return (
    <Space direction="vertical" size="middle" style={{ width: '100%' }}>
      <Card size="small">
        <Space direction="vertical" size={8} style={{ width: '100%' }}>
          <Typography.Text type="secondary">
            可从资源视角跳转到 GitOps 发布中心，按资源上下文快速定位交付单元。
          </Typography.Text>
          <Space wrap>
            <Button type="primary" onClick={() => void navigate('/gitops')}>
              进入 GitOps 发布中心
            </Button>
            <Button onClick={() => void navigate('/backup-restore')}>
              进入备份恢复中心
            </Button>
            <Button
              disabled={!resourceKeyword}
              onClick={() => {
                if (!resourceKeyword) {
                  return;
                }
                void navigate(`/gitops?keyword=${encodeURIComponent(resourceKeyword)}`);
              }}
              >
                按资源上下文跳转
              </Button>
            <Button
              disabled={!resourceKeyword}
              onClick={() => {
                if (!resourceKeyword) {
                  return;
                }
                void navigate(`/backup-restore/restore-points?keyword=${encodeURIComponent(resourceKeyword)}`);
              }}
            >
              按资源上下文恢复
            </Button>
          </Space>
        </Space>
      </Card>

      <ResourceListPage />
    </Space>
  );
};
