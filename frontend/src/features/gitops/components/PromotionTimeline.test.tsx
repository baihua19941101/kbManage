import '@testing-library/jest-dom/vitest';
import { render, screen } from '@testing-library/react';
import { PromotionTimeline } from '@/features/gitops/components/PromotionTimeline';
import { installAntdDomShims } from '@/test/installAntdDomShims';

describe('PromotionTimeline', () => {
  beforeAll(() => {
    installAntdDomShims();
  });

  it('renders stage timeline with operation progress', () => {
    render(
      <PromotionTimeline
        stages={[
          {
            name: 'dev',
            orderIndex: 1,
            targetGroupId: 1,
            promotionMode: 'manual',
            paused: false
          },
          {
            name: 'prod',
            orderIndex: 2,
            targetGroupId: 2,
            promotionMode: 'manual',
            paused: false
          }
        ]}
        status={{
          environments: [
            {
              environment: 'dev',
              syncStatus: 'succeeded',
              driftStatus: 'in_sync',
              targetCount: 2,
              succeededCount: 2,
              failedCount: 0
            },
            {
              environment: 'prod',
              syncStatus: 'pending',
              driftStatus: 'unknown',
              targetCount: 2,
              succeededCount: 0,
              failedCount: 0
            }
          ]
        }}
        operation={{
          id: 301,
          operationType: 'promote',
          status: 'running',
          stages: [
            {
              environment: 'prod',
              status: 'running',
              targetCount: 2,
              succeededCount: 1,
              failedCount: 0
            }
          ]
        }}
      />
    );

    expect(screen.getByText('当前动作：promote / running')).toBeInTheDocument();
    expect(screen.getByText('dev')).toBeInTheDocument();
    expect(screen.getByText('prod')).toBeInTheDocument();
    expect(screen.getByText('阶段状态：succeeded')).toBeInTheDocument();
    expect(screen.getByText('阶段状态：running')).toBeInTheDocument();
    expect(screen.getByText('漂移：in_sync')).toBeInTheDocument();
  });

  it('shows empty state when no stages configured', () => {
    render(<PromotionTimeline stages={[]} />);

    expect(screen.getByText('暂无环境阶段')).toBeInTheDocument();
  });
});
