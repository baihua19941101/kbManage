import { useMemo, useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Alert, Card, Empty, Select, Space, Table, Typography } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { PageHeader } from '@/features/platform-marketplace/components/PageHeader';
import { PermissionDenied } from '@/features/platform-marketplace/components/PermissionDenied';
import { StatusTag } from '@/features/platform-marketplace/components/status';
import { usePlatformMarketplacePermissions } from '@/features/platform-marketplace/hooks/permissions';
import { normalizeApiError } from '@/services/api/client';
import {
  isTemplateOffline,
  listTemplates,
  platformMarketplaceQueryKeys,
  templateQueryScope,
  type ApplicationTemplate
} from '@/services/platformMarketplace';

const columns: ColumnsType<ApplicationTemplate> = [
  { title: '模板名称', dataIndex: 'name', key: 'name' },
  { title: '分类', dataIndex: 'category', key: 'category', render: (value?: string) => value || '—' },
  { title: '来源', dataIndex: 'sourceName', key: 'sourceName', render: (value?: string) => value || '—' },
  {
    title: '最新版本',
    dataIndex: 'latestVersion',
    key: 'latestVersion',
    render: (value?: string) => value || '—'
  },
  { title: '状态', dataIndex: 'status', key: 'status', render: (value?: string) => <StatusTag value={value} /> }
];

export const TemplateCatalogPage = () => {
  const permissions = usePlatformMarketplacePermissions();
  const [category, setCategory] = useState<string>();
  const [selectedTemplateId, setSelectedTemplateId] = useState<string>();
  const templatesQuery = useQuery({
    queryKey: platformMarketplaceQueryKeys.templates(templateQueryScope({ category })),
    enabled: permissions.canRead,
    queryFn: () => listTemplates({ category })
  });

  const categories = useMemo(
    () =>
      Array.from(
        new Set((templatesQuery.data?.items || []).map((item) => item.category).filter(Boolean))
      ).map((item) => ({
        label: item,
        value: item
      })),
    [templatesQuery.data?.items]
  );
  const items = templatesQuery.data?.items || [];
  const selected = items.find((item) => item.id === selectedTemplateId) || items[0];

  if (!permissions.canRead) {
    return <PermissionDenied description="你暂无模板中心访问权限。" />;
  }

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <PageHeader
        title="模板中心"
        description="查看模板分类、版本、依赖、参数表单和部署约束，识别可发布与已下线资产。"
      />

      {templatesQuery.error ? (
        <Alert
          type="error"
          showIcon
          message="模板中心加载失败"
          description={normalizeApiError(templatesQuery.error, '模板中心加载失败，请稍后重试。')}
        />
      ) : null}

      <Card size="small" title="筛选模板">
        <Select
          allowClear
          style={{ width: 280 }}
          placeholder="按分类筛选"
          options={categories}
          value={category}
          onChange={(value) => setCategory(value)}
        />
      </Card>

      <Card size="small" title={`模板列表（${items.length}）`}>
        <Table<ApplicationTemplate>
          rowKey={(record) => record.id}
          columns={columns}
          dataSource={items}
          loading={templatesQuery.isLoading || templatesQuery.isFetching}
          pagination={{ pageSize: 6 }}
          onRow={(record) => ({ onClick: () => setSelectedTemplateId(record.id) })}
        />
      </Card>

      {selected ? (
        <Card size="small" title={`模板详情预览：${selected.name}`}>
          <Typography.Paragraph>范围：{selected.scopeSummary || '未提供范围摘要'}</Typography.Paragraph>
          <Typography.Paragraph>依赖：{selected.dependencySummary || '未声明依赖'}</Typography.Paragraph>
          <Typography.Paragraph>
            发布说明：{selected.releaseNoteSummary || '当前版本未提供发布说明'}
          </Typography.Paragraph>
          {isTemplateOffline(selected) ? (
            <Alert
              type="warning"
              showIcon
              message="该模板已处于下线或历史保留状态"
              description="新安装入口应由主线程在全局路由接线后根据后端约束进行封禁。"
            />
          ) : null}
        </Card>
      ) : (
        <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description="当前没有可展示的模板。" />
      )}
    </Space>
  );
};
