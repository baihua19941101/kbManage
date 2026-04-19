import { useMemo, useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Alert, Button, Card, Empty, Select, Space, Table } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { PageHeader } from '@/features/platform-marketplace/components/PageHeader';
import { PermissionDenied } from '@/features/platform-marketplace/components/PermissionDenied';
import { StatusTag } from '@/features/platform-marketplace/components/status';
import { TemplateReleaseDrawer } from '@/features/platform-marketplace/components/TemplateReleaseDrawer';
import { usePlatformMarketplacePermissions } from '@/features/platform-marketplace/hooks/permissions';
import { useTemplateDistributionAction } from '@/features/platform-marketplace/hooks/useTemplateDistributionAction';
import { normalizeApiError } from '@/services/api/client';
import {
  listTemplateReleases,
  listTemplates,
  platformMarketplaceQueryKeys,
  summarizeTargetScope,
  type ApplicationTemplate,
  type TemplateReleaseScope
} from '@/services/platformMarketplace';

const columns: ColumnsType<TemplateReleaseScope> = [
  { title: '模板', dataIndex: 'templateName', key: 'templateName', render: (value?: string) => value || '—' },
  { title: '版本', dataIndex: 'version', key: 'version', render: (value?: string) => value || '—' },
  {
    title: '目标范围',
    key: 'target',
    render: (_, record) => summarizeTargetScope(record.targetType, record.targetRef)
  },
  { title: '状态', dataIndex: 'status', key: 'status', render: (value?: string) => <StatusTag value={value} /> }
];

export const TemplateDistributionPage = () => {
  const permissions = usePlatformMarketplacePermissions();
  const action = useTemplateDistributionAction();
  const [templateId, setTemplateId] = useState<string>();
  const [drawerOpen, setDrawerOpen] = useState(false);

  const templatesQuery = useQuery({
    queryKey: platformMarketplaceQueryKeys.templates('distribution'),
    enabled: permissions.canRead,
    queryFn: () => listTemplates({})
  });

  const releaseQuery = useQuery({
    queryKey: platformMarketplaceQueryKeys.templateReleases(templateId),
    enabled: permissions.canRead && Boolean(templateId),
    queryFn: () => listTemplateReleases(templateId || '')
  });

  const templateOptions = useMemo(
    () =>
      (templatesQuery.data?.items || []).map((item: ApplicationTemplate) => ({
        label: `${item.name} / ${item.latestVersion || '未标记版本'}`,
        value: item.id
      })),
    [templatesQuery.data?.items]
  );

  if (!permissions.canRead) {
    return <PermissionDenied description="你暂无模板分发访问权限。" />;
  }

  const releases = releaseQuery.data?.items || [];

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <PageHeader
        title="模板分发"
        description="将标准模板发布到工作空间、项目或集群范围，并跟踪范围约束与版本说明。"
        actions={
          <Button type="primary" disabled={!permissions.canPublishTemplate || !templateId} onClick={() => setDrawerOpen(true)}>
            发布模板
          </Button>
        }
      />

      {templatesQuery.error ? (
        <Alert
          type="error"
          showIcon
          message="模板列表加载失败"
          description={normalizeApiError(templatesQuery.error, '模板列表加载失败，请稍后重试。')}
        />
      ) : null}
      {releaseQuery.error ? (
        <Alert
          type="error"
          showIcon
          message="模板分发记录加载失败"
          description={normalizeApiError(releaseQuery.error, '模板分发记录加载失败，请稍后重试。')}
        />
      ) : null}

      <Card size="small" title="选择模板">
        <Select
          allowClear
          showSearch
          style={{ width: 360 }}
          placeholder="选择一个模板查看分发记录"
          options={templateOptions}
          value={templateId}
          onChange={(value) => setTemplateId(value)}
        />
      </Card>

      {!templateId ? (
        <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description="请选择一个模板查看分发记录。" />
      ) : (
        <Card size="small" title={`分发记录（${releases.length}）`}>
          <Table<TemplateReleaseScope>
            rowKey={(record) => record.id}
            columns={columns}
            dataSource={releases}
            loading={releaseQuery.isLoading || releaseQuery.isFetching}
            pagination={{ pageSize: 6 }}
          />
        </Card>
      )}

      <TemplateReleaseDrawer
        open={drawerOpen}
        submitting={action.createReleaseMutation.isPending}
        onClose={() => setDrawerOpen(false)}
        onSubmit={(payload) => {
          if (!templateId) {
            return;
          }
          action.createReleaseMutation.mutate(
            { templateId, payload },
            {
              onSuccess: () => setDrawerOpen(false)
            }
          );
        }}
      />
    </Space>
  );
};
