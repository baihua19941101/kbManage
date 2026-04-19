# Data Model: 平台 SRE 与规模化治理

## 1. HA Policy

- **Purpose**: 定义平台控制面的高可用策略和故障切换规则。
- **Core Fields**:
  - `id`
  - `name`
  - `controlPlaneScope`
  - `deploymentMode`
  - `replicaExpectation`
  - `failoverTriggerPolicy`
  - `failoverCooldown`
  - `takeoverStatus`
  - `lastFailoverAt`
  - `lastRecoveryResult`
  - `status`
  - `owner`
  - `updatedAt`
- **Validation Rules**:
  - `deploymentMode` 必须属于受支持的高可用模式。
  - `replicaExpectation` 必须大于 1 才能声明为高可用。
  - `failoverTriggerPolicy` 必须明确触发条件和保护阈值。
- **Relationships**:
  - 一个 `HA Policy` 可关联多个 `Platform Health Snapshot`。
  - 一个 `HA Policy` 可关联多个 `Maintenance Window` 作为限制条件。

## 2. Maintenance Window

- **Purpose**: 定义允许执行维护、升级和受限运维动作的时间窗口。
- **Core Fields**:
  - `id`
  - `name`
  - `windowType`
  - `scope`
  - `startAt`
  - `endAt`
  - `allowedOperations`
  - `restrictedOperations`
  - `status`
  - `exceptionReason`
  - `approvalRecord`
  - `postCheckStatus`
  - `updatedAt`
- **Validation Rules**:
  - 结束时间必须晚于开始时间。
  - 例外执行时必须填写原因和审批记录。
- **Relationships**:
  - 一个 `Maintenance Window` 可关联多个 `Upgrade Plan`。
  - 一个 `Maintenance Window` 可被 `HA Policy` 和 `Platform Health Snapshot` 引用。

## 3. Platform Health Snapshot

- **Purpose**: 描述某一时间点的平台运行健康总览。
- **Core Fields**:
  - `id`
  - `snapshotAt`
  - `componentHealthSummary`
  - `dependencyHealthSummary`
  - `taskBacklogSummary`
  - `capacityRiskLevel`
  - `throttlingStatus`
  - `recoverySummary`
  - `maintenanceStatus`
  - `overallStatus`
  - `recommendedActions`
- **Validation Rules**:
  - 必须同时包含平台组件和依赖状态摘要。
  - `overallStatus` 必须可追溯到子项摘要。
- **Relationships**:
  - 多个 `Platform Health Snapshot` 归属于同一平台控制面范围。
  - `Platform Health Snapshot` 可关联 `Runbook Article` 和 `Alert Baseline`。

## 4. Capacity Baseline

- **Purpose**: 保存性能基线、容量阈值和趋势判断。
- **Core Fields**:
  - `id`
  - `name`
  - `resourceDimension`
  - `baselineRange`
  - `thresholds`
  - `growthTrend`
  - `forecastWindow`
  - `forecastResult`
  - `confidenceLevel`
  - `status`
  - `updatedAt`
- **Validation Rules**:
  - 必须至少定义一个阈值区间。
  - `confidenceLevel` 不能为空。
- **Relationships**:
  - 一个 `Capacity Baseline` 可关联多个 `Scale Evidence`。
  - 一个 `Capacity Baseline` 可被 `Platform Health Snapshot` 引用。

## 5. Upgrade Plan

- **Purpose**: 表示一次平台升级治理闭环。
- **Core Fields**:
  - `id`
  - `name`
  - `currentVersion`
  - `targetVersion`
  - `compatibilitySummary`
  - `precheckResult`
  - `rolloutStrategy`
  - `executionStage`
  - `executionProgress`
  - `acceptanceResult`
  - `rollbackReadiness`
  - `status`
  - `createdBy`
  - `updatedAt`
- **Validation Rules**:
  - `currentVersion` 与 `targetVersion` 不能为空且不能相同。
  - 未通过 `precheckResult` 时不得进入执行阶段。
- **Relationships**:
  - 一个 `Upgrade Plan` 可关联一个 `Maintenance Window`。
  - 一个 `Upgrade Plan` 可关联多个 `Rollback Validation` 记录。

## 6. Rollback Validation

- **Purpose**: 记录某次升级回退方案的验证结果。
- **Core Fields**:
  - `id`
  - `upgradePlanId`
  - `validationScope`
  - `preconditions`
  - `result`
  - `remainingRisk`
  - `validatedAt`
  - `validatedBy`
- **Validation Rules**:
  - 必须明确验证结论。
  - 如果结果不是通过，必须记录剩余风险。
- **Relationships**:
  - 多个 `Rollback Validation` 归属于同一 `Upgrade Plan`。

## 7. Runbook Article

- **Purpose**: 提供平台异常处理、升级、恢复和维护的运行手册。
- **Core Fields**:
  - `id`
  - `title`
  - `scenarioType`
  - `applicableComponents`
  - `riskLevel`
  - `checklistSummary`
  - `recoverySteps`
  - `verificationSummary`
  - `status`
  - `updatedAt`
- **Validation Rules**:
  - 至少定义适用场景和恢复步骤。
- **Relationships**:
  - 一个 `Runbook Article` 可关联多个 `Platform Health Snapshot` 或 `Alert Baseline`。

## 8. Alert Baseline

- **Purpose**: 维护平台组件和运维场景的告警基线。
- **Core Fields**:
  - `id`
  - `name`
  - `componentScope`
  - `signalType`
  - `baselineCondition`
  - `severity`
  - `recommendedRunbookId`
  - `status`
  - `updatedAt`
- **Validation Rules**:
  - 必须指定组件范围和信号类型。
- **Relationships**:
  - 一个 `Alert Baseline` 可以关联一个首选 `Runbook Article`。

## 9. Scale Evidence

- **Purpose**: 保存压测结果、瓶颈分析、自诊断和容量预测的证据。
- **Core Fields**:
  - `id`
  - `evidenceType`
  - `scope`
  - `sampleWindow`
  - `summary`
  - `bottleneckSummary`
  - `forecastSummary`
  - `confidenceLevel`
  - `recoveryObservation`
  - `status`
  - `capturedAt`
- **Validation Rules**:
  - `evidenceType` 必须属于压测、预测、自诊断或恢复观测之一。
  - `confidenceLevel` 必须明确。
- **Relationships**:
  - `Scale Evidence` 可关联 `Capacity Baseline`、`Platform Health Snapshot` 和 `Runbook Article`。

## 10. State Notes

- `HA Policy.status`: `draft -> active -> degraded -> recovering`
- `Maintenance Window.status`: `scheduled -> active -> completed -> exception`
- `Upgrade Plan.status`: `draft -> precheck-failed|ready -> rolling -> accepted|rollback-required -> closed`
- `Rollback Validation.result`: `passed | warning | failed`
- `Scale Evidence.status`: `captured -> analyzed -> archived`
