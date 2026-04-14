import { Button, Card, Empty, Input, InputNumber, Select, Space, Switch, Typography } from 'antd';
import type { GitOpsEnvironmentStage, GitOpsTargetGroupItem } from '@/services/gitops';

type EnvironmentStageEditorProps = {
  value: GitOpsEnvironmentStage[];
  targetGroups?: GitOpsTargetGroupItem[];
  onChange: (value: GitOpsEnvironmentStage[]) => void;
};

const buildEmptyStage = (index: number, defaultTargetGroupId?: number): GitOpsEnvironmentStage => ({
  name: `stage-${index}`,
  orderIndex: index,
  targetGroupId: defaultTargetGroupId ?? 1,
  promotionMode: 'manual',
  paused: false
});

export const EnvironmentStageEditor = ({
  value,
  targetGroups = [],
  onChange
}: EnvironmentStageEditorProps) => {
  const defaultTargetGroupId = targetGroups[0]?.id ? Number(targetGroups[0].id) : undefined;

  const handleItemChange = (
    index: number,
    patch: Partial<GitOpsEnvironmentStage>
  ) => {
    const next = value.map((item, itemIndex) =>
      itemIndex === index ? { ...item, ...patch } : item
    );
    onChange(next);
  };

  const handleRemove = (index: number) => {
    const next = value.filter((_, itemIndex) => itemIndex !== index);
    onChange(next);
  };

  const handleAdd = () => {
    const next = [
      ...value,
      buildEmptyStage(value.length + 1, defaultTargetGroupId)
    ];
    onChange(next);
  };

  return (
    <Card
      size="small"
      title="环境阶段"
      extra={
        <Button size="small" onClick={handleAdd}>
          新增阶段
        </Button>
      }
    >
      {value.length === 0 ? <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description="暂无环境阶段" /> : null}
      <Space direction="vertical" size={12} style={{ width: '100%' }}>
        {value.map((stage, index) => (
          <Card
            key={`${stage.name}-${index}`}
            size="small"
            bodyStyle={{ paddingBottom: 12 }}
            extra={
              <Button danger type="link" onClick={() => handleRemove(index)}>
                删除
              </Button>
            }
          >
            <Space direction="vertical" size={8} style={{ width: '100%' }}>
              <Input
                value={stage.name}
                placeholder="阶段名称"
                onChange={(event) => handleItemChange(index, { name: event.target.value })}
              />
              <Space wrap>
                <Typography.Text type="secondary">顺序</Typography.Text>
                <InputNumber
                  min={1}
                  value={stage.orderIndex}
                  onChange={(value) =>
                    handleItemChange(index, { orderIndex: Number(value || 1) })
                  }
                />
                <Typography.Text type="secondary">目标组</Typography.Text>
                <Select
                  style={{ minWidth: 180 }}
                  value={stage.targetGroupId}
                  options={targetGroups.map((item) => ({
                    label: `${item.name} (#${item.id})`,
                    value: Number(item.id)
                  }))}
                  onChange={(targetGroupId) => handleItemChange(index, { targetGroupId })}
                />
                <Typography.Text type="secondary">推进模式</Typography.Text>
                <Select
                  style={{ width: 120 }}
                  value={stage.promotionMode}
                  options={[
                    { label: 'manual', value: 'manual' },
                    { label: 'automatic', value: 'automatic' }
                  ]}
                  onChange={(promotionMode) =>
                    handleItemChange(index, { promotionMode })
                  }
                />
                <Typography.Text type="secondary">暂停</Typography.Text>
                <Switch
                  checked={Boolean(stage.paused)}
                  onChange={(paused) => handleItemChange(index, { paused })}
                />
              </Space>
            </Space>
          </Card>
        ))}
      </Space>
    </Card>
  );
};
