import { useQuery } from '@tanstack/react-query';
import { Alert, Card, Descriptions, Empty, List, Space } from 'antd';
import { useParams } from 'react-router-dom';
import { PageHeader } from '@/features/platform-marketplace/components/PageHeader';
import { PermissionDenied } from '@/features/platform-marketplace/components/PermissionDenied';
import { StatusTag } from '@/features/platform-marketplace/components/status';
import { usePlatformMarketplacePermissions } from '@/features/platform-marketplace/hooks/permissions';
import { normalizeApiError } from '@/services/api/client';
import {
  getTemplateDetail,
  platformMarketplaceQueryKeys
} from '@/services/platformMarketplace';

export const TemplateDetailPage = () => {
  const permissions = usePlatformMarketplacePermissions();
  const { templateId } = useParams<{ templateId: string }>();
  const detailQuery = useQuery({
    queryKey: platformMarketplaceQueryKeys.templateDetail(templateId),
    enabled: permissions.canRead && Boolean(templateId),
    queryFn: () => getTemplateDetail(templateId || '')
  });

  if (!permissions.canRead) {
    return <PermissionDenied description="你暂无模板详情访问权限。" />;
  }

  if (!templateId) {
    return <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description="请选择一个模板查看详情。" />;
  }

  const detail = detailQuery.data;

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <PageHeader
        title="模板详情"
        description="查看模板版本、参数摘要、依赖与部署约束，为发布和升级判断提供依据。"
      />

      {detailQuery.error ? (
        <Alert
          type="error"
          showIcon
          message="模板详情加载失败"
          description={normalizeApiError(detailQuery.error, '模板详情加载失败，请稍后重试。')}
        />
      ) : null}

      {detail ? (
        <>
          <Card size="small">
            <Descriptions column={2} size="small" bordered>
              <Descriptions.Item label="模板名称">{detail.name}</Descriptions.Item>
              <Descriptions.Item label="状态">
                <StatusTag value={detail.status} />
              </Descriptions.Item>
              <Descriptions.Item label="分类">{detail.category || '—'}</Descriptions.Item>
              <Descriptions.Item label="最新版本">{detail.latestVersion || '—'}</Descriptions.Item>
              <Descriptions.Item label="适用范围" span={2}>
                {detail.scopeSummary || '—'}
              </Descriptions.Item>
              <Descriptions.Item label="依赖摘要" span={2}>
                {detail.dependencySummary || '—'}
              </Descriptions.Item>
              <Descriptions.Item label="发布说明" span={2}>
                {detail.releaseNoteSummary || '—'}
              </Descriptions.Item>
            </Descriptions>
          </Card>
          <Card size="small" title="版本列表">
            <List
              locale={{ emptyText: '当前没有版本信息。' }}
              dataSource={detail.templateVersions || []}
              renderItem={(item) => (
                <List.Item>
                  <List.Item.Meta
                    title={`${item.version} / ${item.status || '未知状态'}`}
                    description={`${item.dependencySummary || '无依赖摘要'} | ${item.constraintSummary || '无部署约束'}`}
                  />
                </List.Item>
              )}
            />
          </Card>
        </>
      ) : (
        <Card loading={detailQuery.isLoading || detailQuery.isFetching} />
      )}
    </Space>
  );
};
