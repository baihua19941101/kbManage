import { useMemo, useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Alert, Card, Empty, Select, Space, Typography } from 'antd';
import { canReadClusterLifecycle, useAuthStore } from '@/features/auth/store';
import { CapabilityMatrixTable } from '@/features/cluster-lifecycle/components/CapabilityMatrixTable';
import { PageHeader } from '@/features/cluster-lifecycle/components/PageHeader';
import { useCapabilityMatrix } from '@/features/cluster-lifecycle/hooks/useCapabilityMatrix';
import { normalizeApiError } from '@/services/api/client';
import { clusterLifecycleQueryKeys, listClusterDrivers } from '@/services/clusterLifecycle';

export const CapabilityMatrixPage = () => {
  const user = useAuthStore((state) => state.user);
  const canRead = canReadClusterLifecycle(user);
  const [driverId, setDriverId] = useState<string>();

  const driversQuery = useQuery({
    queryKey: clusterLifecycleQueryKeys.drivers('capability-page'),
    enabled: canRead,
    queryFn: () => listClusterDrivers()
  });

  const options = useMemo(
    () =>
      (driversQuery.data?.items || []).map((item) => ({
        label: `${item.displayName || item.driverKey} / ${item.version}`,
        value: item.id
      })),
    [driversQuery.data?.items]
  );

  const matrixQuery = useCapabilityMatrix(driverId);

  if (!canRead) {
    return <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description="你暂无能力矩阵访问权限。" />;
  }

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <PageHeader
        title="能力矩阵"
        description="对比不同驱动在网络、存储、身份、监控、安全、备份和发布方面的兼容性。"
      />

      {driversQuery.error ? (
        <Alert
          type="error"
          showIcon
          message="驱动列表加载失败"
          description={normalizeApiError(driversQuery.error, '驱动列表加载失败，请稍后重试。')}
        />
      ) : null}

      <Card size="small" title="筛选驱动">
        <Select
          allowClear
          showSearch
          style={{ width: 360 }}
          placeholder="选择驱动版本"
          options={options}
          value={driverId}
          onChange={(value) => setDriverId(value)}
        />
      </Card>

      {!driverId ? (
        <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description="请选择一个驱动版本查看能力矩阵。" />
      ) : (
        <>
          {matrixQuery.error ? (
            <Alert
              type="error"
              showIcon
              message="能力矩阵加载失败"
              description={normalizeApiError(matrixQuery.error, '能力矩阵加载失败，请稍后重试。')}
            />
          ) : null}
          <Card size="small" title={`能力域（${matrixQuery.data?.items.length ?? 0}）`}>
            <CapabilityMatrixTable
              data={matrixQuery.data?.items || []}
              loading={matrixQuery.isLoading || matrixQuery.isFetching}
            />
          </Card>
          <Typography.Text type="secondary">
            当前共识别 {matrixQuery.owners.length} 个 owner，后续可在此基础上扩展跨驱动对比视图。
          </Typography.Text>
        </>
      )}
    </Space>
  );
};
