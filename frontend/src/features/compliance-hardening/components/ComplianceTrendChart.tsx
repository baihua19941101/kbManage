import ReactECharts from 'echarts-for-react';
import { Card, Empty, Typography } from 'antd';
import type { ComplianceTrendResponse } from '@/services/compliance';

export const ComplianceTrendChart = ({ data }: { data?: ComplianceTrendResponse }) => {
  const points = data?.points || [];

  if (points.length === 0) {
    return (
      <Card size="small" title="趋势图">
        <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description="暂无趋势数据。" />
      </Card>
    );
  }

  const option = {
    tooltip: {
      trigger: 'axis'
    },
    legend: {
      data: ['得分', '覆盖率', '整改完成率']
    },
    xAxis: {
      type: 'category',
      data: points.map((point) => point.windowStart || point.windowEnd || '-')
    },
    yAxis: {
      type: 'value',
      max: 100
    },
    series: [
      {
        name: '得分',
        type: 'line',
        smooth: true,
        data: points.map((point) => point.scoreAvg ?? 0)
      },
      {
        name: '覆盖率',
        type: 'line',
        smooth: true,
        data: points.map((point) => point.coverageRate ?? 0)
      },
      {
        name: '整改完成率',
        type: 'line',
        smooth: true,
        data: points.map((point) => point.remediationCompletionRate ?? 0)
      }
    ]
  };

  return (
    <Card
      size="small"
      title="趋势图"
      extra={
        data?.comparisonBasis?.baselineVersions?.length ? (
          <Typography.Text type="secondary">
            基线版本：{data.comparisonBasis.baselineVersions.join(', ')}
          </Typography.Text>
        ) : null
      }
    >
      <ReactECharts option={option} style={{ height: 320 }} />
    </Card>
  );
};
