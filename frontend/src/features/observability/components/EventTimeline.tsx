import { Timeline, Typography } from 'antd';
import type { ObservabilityEventDTO } from '@/services/api/types';

type EventTimelineProps = {
  items: ObservabilityEventDTO[];
};

export const EventTimeline = ({ items }: EventTimelineProps) => {
  return (
    <Timeline
      items={items.map((item) => ({
        color: item.eventType === 'warning' ? 'red' : 'green',
        children: (
          <div>
            <Typography.Text strong>{item.reason ?? 'Event'}</Typography.Text>
            <br />
            <Typography.Text type="secondary">{item.message ?? '-'}</Typography.Text>
          </div>
        )
      }))}
    />
  );
};
