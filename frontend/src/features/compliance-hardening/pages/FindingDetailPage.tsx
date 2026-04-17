import { useMemo, useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Alert, Card, Descriptions, Empty, Space, Tag, Typography } from 'antd';
import { useNavigate, useParams } from 'react-router-dom';
import { canReadCompliance, useAuthStore } from '@/features/auth/store';
import { EvidenceDrawer } from '@/features/compliance-hardening/components/EvidenceDrawer';
import { FindingTable } from '@/features/compliance-hardening/components/FindingTable';
import {
  findingResultColorMap,
  formatDateTime,
  remediationStatusColorMap,
  riskColorMap
} from '@/features/compliance-hardening/utils';
import { normalizeApiError } from '@/services/api/client';
import {
  getComplianceFinding,
  listComplianceFindings
} from '@/services/compliance';

export const FindingDetailPage = () => {
  const { findingId = '' } = useParams();
  const navigate = useNavigate();
  const user = useAuthStore((state) => state.user);
  const canRead = canReadCompliance(user);
  const [evidenceOpen, setEvidenceOpen] = useState(false);

  const findingDetailQuery = useQuery({
    queryKey: ['compliance', 'finding-detail', findingId],
    enabled: canRead && Boolean(findingId),
    queryFn: () => getComplianceFinding(findingId)
  });

  const relatedFindingsQuery = useQuery({
    queryKey: ['compliance', 'finding-related', findingId],
    enabled: canRead && Boolean(findingId),
    queryFn: async () => {
      const detail = await getComplianceFinding(findingId);
      return listComplianceFindings({
        clusterId: detail.clusterId,
        namespace: detail.namespace,
        riskLevel: detail.riskLevel
      });
    }
  });

  const finding = findingDetailQuery.data;
  const relatedFindings = useMemo(
    () => (relatedFindingsQuery.data?.items || []).filter((item) => item.id !== findingId),
    [findingId, relatedFindingsQuery.data?.items]
  );

  if (!canRead) {
    return <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description="你暂无失败项详情访问权限。" />;
  }

  return (
    <Space direction="vertical" size={16} style={{ width: '100%' }}>
      <div>
        <Typography.Title level={3} style={{ marginBottom: 8 }}>
          失败项详情
        </Typography.Title>
        <Typography.Text type="secondary">
          查看失败项、证据、整改和复检链路，并从相同上下文快速跳转其他失败项。
        </Typography.Text>
      </div>

      {findingDetailQuery.error ? (
        <Alert type="error" showIcon message="失败项详情加载失败" description={normalizeApiError(findingDetailQuery.error, '失败项详情加载失败，请稍后重试。')} />
      ) : null}

      <Card loading={findingDetailQuery.isLoading || findingDetailQuery.isFetching} size="small" title={finding?.controlTitle || finding?.controlId || '失败项详情'} extra={finding?.id ? <Typography.Text code>{finding.id}</Typography.Text> : null}>
        {finding ? (
          <Space direction="vertical" size={16} style={{ width: '100%' }}>
            <Descriptions bordered size="small" column={2}>
              <Descriptions.Item label="结果">
                {finding.result ? <Tag color={findingResultColorMap[finding.result]}>{finding.result}</Tag> : '—'}
              </Descriptions.Item>
              <Descriptions.Item label="风险">
                {finding.riskLevel ? <Tag color={riskColorMap[finding.riskLevel]}>{finding.riskLevel}</Tag> : '—'}
              </Descriptions.Item>
              <Descriptions.Item label="治理状态">
                {finding.remediationStatus ? <Tag color={remediationStatusColorMap[finding.remediationStatus]}>{finding.remediationStatus}</Tag> : '—'}
              </Descriptions.Item>
              <Descriptions.Item label="证据数">{finding.evidences?.length || 0}</Descriptions.Item>
              <Descriptions.Item label="资源上下文" span={2}>
                {[finding.clusterId, finding.namespace, finding.resourceKind, finding.resourceName].filter(Boolean).join(' / ') || '—'}
              </Descriptions.Item>
              <Descriptions.Item label="摘要" span={2}>{finding.summary || '—'}</Descriptions.Item>
              <Descriptions.Item label="整改任务">{finding.remediationTasks?.length || 0}</Descriptions.Item>
              <Descriptions.Item label="复检任务">{finding.rechecks?.length || 0}</Descriptions.Item>
            </Descriptions>
            <Space>
              <Tag>基线证据更新时间：{formatDateTime(finding.evidences?.[0]?.collectedAt)}</Tag>
              <Typography.Link onClick={() => setEvidenceOpen(true)}>查看证据</Typography.Link>
            </Space>
          </Space>
        ) : (
          <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description="未找到失败项详情。" />
        )}
      </Card>

      <Card size="small" title={`相近失败项（${relatedFindings.length}）`}>
        <FindingTable
          findings={relatedFindings}
          loading={relatedFindingsQuery.isLoading || relatedFindingsQuery.isFetching}
          onView={(item) => void navigate(`/compliance-hardening/findings/${item.id}`)}
        />
      </Card>

      <EvidenceDrawer open={evidenceOpen} evidences={finding?.evidences} onClose={() => setEvidenceOpen(false)} />
    </Space>
  );
};
