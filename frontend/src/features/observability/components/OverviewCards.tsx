import { Card, Col, Row, Statistic } from 'antd';
import type { ObservabilityOverviewCardDTO } from '@/services/api/types';

type OverviewCardsProps = {
  cards: ObservabilityOverviewCardDTO[];
};

export const OverviewCards = ({ cards }: OverviewCardsProps) => {
  return (
    <Row gutter={[16, 16]}>
      {cards.map((card) => (
        <Col xs={24} md={12} lg={8} key={card.title}>
          <Card>
            <Statistic title={card.title} value={card.value} suffix={card.unit} />
          </Card>
        </Col>
      ))}
    </Row>
  );
};
