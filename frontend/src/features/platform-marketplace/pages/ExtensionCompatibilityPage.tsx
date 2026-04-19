import { useMemo, useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Alert, Card, Empty, Select, Space } from 'antd';
import { CompatibilitySummaryCard } from '@/features/platform-marketplace/components/CompatibilitySummaryCard';
import { PageHeader } from '@/features/platform-marketplace/components/PageHeader';
import { PermissionDenied } from '@/features/platform-marketplace/components/PermissionDenied';
import { usePlatformMarketplacePermissions } from '@/features/platform-marketplace/hooks/permissions';
import { normalizeApiError } from '@/services/api/client';
import {
  getExtensionCompatibility,
  isExtensionBlocked,
  listExtensions,
  platformMarketplaceQueryKeys
} from '@/services/platformMarketplace';

export const ExtensionCompatibilityPage = () => {
  const permissions = usePlatformMarketplacePermissions();
  const [extensionId, setExtensionId] = useState<string>();
  const extensionsQuery = useQuery({
    queryKey: platformMarketplaceQueryKeys.extensions('compatibility'),
    enabled: permissions.canRead,
    queryFn: () => listExtensions({})
  });
  const compatibilityQuery = useQuery({
    queryKey: platformMarketplaceQueryKeys.extensionCompatibility(extensionId),
    enabled: permissions.canRead && Boolean(extensionId),
    queryFn: () => getExtensionCompatibility(extensionId || '')
  });

  const options = useMemo(
    () =>
      (extensionsQuery.data?.items || []).map((item) => ({
        label: `${item.name} / ${item.version || '未标记版本'}`,
        value: item.id
      })),
    [extensionsQuery.data?.items]
  );

  if (!permissions.canRead) {
    return <PermissionDenied description="你暂无扩展兼容性访问权限。" />;
  }

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <PageHeader
        title="扩展兼容性"
        description="查看扩展兼容结论、阻断原因、权限影响和建议动作，辅助启停决策。"
      />

      {extensionsQuery.error ? (
        <Alert
          type="error"
          showIcon
          message="扩展列表加载失败"
          description={normalizeApiError(extensionsQuery.error, '扩展列表加载失败，请稍后重试。')}
        />
      ) : null}
      {compatibilityQuery.error ? (
        <Alert
          type="error"
          showIcon
          message="兼容性信息加载失败"
          description={normalizeApiError(compatibilityQuery.error, '兼容性信息加载失败，请稍后重试。')}
        />
      ) : null}

      <Card size="small" title="选择扩展">
        <Select
          allowClear
          showSearch
          style={{ width: 360 }}
          placeholder="选择扩展查看兼容性"
          options={options}
          value={extensionId}
          onChange={(value) => setExtensionId(value)}
        />
      </Card>

      {!extensionId ? (
        <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description="请选择一个扩展查看兼容性结论。" />
      ) : (
        <>
          {isExtensionBlocked(compatibilityQuery.data) ? (
            <Alert
              type="warning"
              showIcon
              message="当前扩展存在阻断条件"
              description="启用前请先解决兼容性或权限边界问题。"
            />
          ) : null}
          <CompatibilitySummaryCard
            compatibility={compatibilityQuery.data}
            loading={compatibilityQuery.isLoading || compatibilityQuery.isFetching}
          />
        </>
      )}
    </Space>
  );
};
