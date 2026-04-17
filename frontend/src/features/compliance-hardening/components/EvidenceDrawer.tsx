import { Drawer, Empty, List, Space, Tag, Typography } from 'antd';
import type { EvidenceRecord } from '@/services/compliance';
import { formatDateTime } from '@/features/compliance-hardening/utils';

type EvidenceDrawerProps = {
  open: boolean;
  evidences?: EvidenceRecord[];
  onClose: () => void;
};

export const EvidenceDrawer = ({ open, evidences, onClose }: EvidenceDrawerProps) => {
  return (
    <Drawer title="证据详情" width={520} open={open} onClose={onClose} destroyOnClose>
      {evidences && evidences.length > 0 ? (
        <List
          dataSource={evidences}
          renderItem={(item) => (
            <List.Item key={item.id}>
              <Space direction="vertical" size={4} style={{ width: '100%' }}>
                <Typography.Text strong>{item.summary || item.evidenceType || item.id}</Typography.Text>
                <Typography.Text type="secondary">来源：{item.sourceRef || '—'}</Typography.Text>
                <Space wrap>
                  {item.confidence ? <Tag>{item.confidence}</Tag> : null}
                  {item.redactionStatus ? <Tag>{item.redactionStatus}</Tag> : null}
                  <Tag>{formatDateTime(item.collectedAt)}</Tag>
                </Space>
                {item.artifactRef ? (
                  <Typography.Link href={item.artifactRef} target="_blank">
                    打开证据附件
                  </Typography.Link>
                ) : null}
              </Space>
            </List.Item>
          )}
        />
      ) : (
        <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description="当前失败项暂无可查看证据。" />
      )}
    </Drawer>
  );
};
