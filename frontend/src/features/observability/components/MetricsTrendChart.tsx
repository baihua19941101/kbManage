import { useEffect, useMemo, useRef } from 'react';
import * as echarts from 'echarts';
import type { ObservabilityMetricSeriesDTO } from '@/services/api/types';

type MetricsTrendChartProps = {
  series: ObservabilityMetricSeriesDTO;
};

export const MetricsTrendChart = ({ series }: MetricsTrendChartProps) => {
  const ref = useRef<HTMLDivElement | null>(null);
  const option = useMemo(
    () => ({
      tooltip: { trigger: 'axis' },
      xAxis: {
        type: 'category',
        data: series.points.map((p) => p.timestamp)
      },
      yAxis: { type: 'value' },
      series: [
        {
          data: series.points.map((p) => p.value),
          type: 'line',
          smooth: true
        }
      ]
    }),
    [series]
  );

  useEffect(() => {
    if (!ref.current) {
      return;
    }
    const chart = echarts.init(ref.current);
    chart.setOption(option);
    const onResize = () => chart.resize();
    window.addEventListener('resize', onResize);
    return () => {
      window.removeEventListener('resize', onResize);
      chart.dispose();
    };
  }, [option]);

  return <div ref={ref} style={{ height: 280, width: '100%' }} />;
};
