import { Card, Col, Row, Statistic } from 'antd';
import type { ClusterLifecycleRecord } from '@/services/clusterLifecycle';
import { extractLifecycleSummary } from '@/services/clusterLifecycle';

export const LifecycleSummaryCards = ({
  clusters
}: {
  clusters: ClusterLifecycleRecord[];
}) => {
  const summary = extractLifecycleSummary(clusters);

  return (
    <Row gutter={[16, 16]}>
      <Col xs={24} md={12} xl={6}>
        <Card size="small">
          <Statistic title="集群总数" value={summary.total} />
        </Card>
      </Col>
      <Col xs={24} md={12} xl={6}>
        <Card size="small">
          <Statistic title="运行中" value={summary.active} />
        </Card>
      </Col>
      <Col xs={24} md={12} xl={6}>
        <Card size="small">
          <Statistic title="待处理" value={summary.pending} />
        </Card>
      </Col>
      <Col xs={24} md={12} xl={6}>
        <Card size="small">
          <Statistic title="退役处理中" value={summary.retiring} />
        </Card>
      </Col>
    </Row>
  );
};
