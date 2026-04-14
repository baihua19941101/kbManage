import { Alert, Card, Empty, Space, Tag, Timeline, Typography } from 'antd';
import type { GitOpsOperationDTO } from '@/services/api/types';
import type { GitOpsDeliveryUnitStatus, GitOpsEnvironmentStage } from '@/services/gitops';

type PromotionTimelineProps = {
  stages?: GitOpsEnvironmentStage[];
  status?: GitOpsDeliveryUnitStatus;
  operation?: GitOpsOperationDTO;
};

type StageSnapshot = {
  syncStatus?: string;
  driftStatus?: string;
  targetCount?: number;
  succeededCount?: number;
  failedCount?: number;
};

const stageColorMap: Record<string, string> = {
  running: 'blue',
  pending: 'gray',
  succeeded: 'green',
  partially_succeeded: 'orange',
  failed: 'red',
  paused: 'purple',
  active: 'green',
  historical: 'gray',
  in_sync: 'green',
  drifted: 'orange',
  unknown: 'gray'
};

const resolveStageColor = (status?: string) => stageColorMap[status || 'unknown'] || 'gray';

export const PromotionTimeline = ({ stages = [], status, operation }: PromotionTimelineProps) => {
  const stageStatusMap = new Map<string, StageSnapshot>();
  (status?.environments || []).forEach((item) => {
    const environment = item.environment || '';
    if (!environment) {
      return;
    }

    stageStatusMap.set(environment, {
      syncStatus: item.syncStatus,
      driftStatus: item.driftStatus,
      targetCount: item.targetCount,
      succeededCount: item.succeededCount,
      failedCount: item.failedCount
    });
  });

  const operationStageMap = new Map<string, NonNullable<GitOpsOperationDTO['stages']>[number]>();
  (operation?.stages || []).forEach((item) => {
    if (!item.environment) {
      return;
    }
    operationStageMap.set(item.environment, item);
  });

  const orderedStages = [...stages].sort((left, right) => left.orderIndex - right.orderIndex);

  return (
    <Card size="small" title="环境推进时间线">
      {operation ? (
        <Alert
          type="info"
          showIcon
          style={{ marginBottom: 12 }}
          message={`当前动作：${operation.operationType || operation.actionType || 'unknown'} / ${operation.status || 'pending'}`}
        />
      ) : null}

      {orderedStages.length === 0 ? (
        <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description="暂无环境阶段" />
      ) : (
        <Timeline
          items={orderedStages.map((stage) => {
            const stageStatus = stageStatusMap.get(stage.name);
            const stageOperation = operationStageMap.get(stage.name);
            const timelineState =
              stageOperation?.status ||
              stageStatus?.syncStatus ||
              (stage.paused ? 'paused' : 'pending');

            return {
              color: resolveStageColor(timelineState),
              children: (
                <Space direction="vertical" size={4}>
                  <Typography.Text strong>{stage.name}</Typography.Text>
                  <Space wrap>
                    <Tag>{stage.promotionMode}</Tag>
                    <Tag>目标组 #{stage.targetGroupId}</Tag>
                    <Tag color={resolveStageColor(timelineState)}>阶段状态：{timelineState}</Tag>
                    {stageStatus?.driftStatus ? (
                      <Tag color={resolveStageColor(stageStatus.driftStatus)}>
                        漂移：{stageStatus.driftStatus}
                      </Tag>
                    ) : null}
                    {typeof stageStatus?.targetCount === 'number' ? (
                      <Tag>
                        目标/成功/失败：
                        {stageStatus.targetCount}/{stageStatus.succeededCount ?? 0}/
                        {stageStatus.failedCount ?? 0}
                      </Tag>
                    ) : null}
                    {stage.paused ? <Tag color="purple">已暂停</Tag> : null}
                  </Space>
                </Space>
              )
            };
          })}
        />
      )}
    </Card>
  );
};
