# Data Model: 企业级治理报表与产品化交付收尾

## 1. Permission Change Trail

- **Purpose**: 记录一次权限授予、变更、委派或回收的完整治理链路。
- **Core Fields**:
  - `id`
  - `subjectType`
  - `subjectRef`
  - `sourceIdentity`
  - `changeType`
  - `beforeState`
  - `afterState`
  - `authorizationBasis`
  - `approvalReference`
  - `scopeType`
  - `scopeRef`
  - `changedAt`
  - `changedBy`
  - `evidenceCompleteness`
- **Validation Rules**:
  - 必须同时保留变更前和变更后状态摘要。
  - 若存在审批或授权依据缺失，`evidenceCompleteness` 必须标记为不完整。
- **Relationships**:
  - 一个 `Permission Change Trail` 可关联多个 `Governance Risk Event`。
  - 一个 `Permission Change Trail` 可被多个 `Governance Report Package` 引用。

## 2. Key Operation Trace

- **Purpose**: 表示一次关键治理操作的追踪记录。
- **Core Fields**:
  - `id`
  - `actorType`
  - `actorRef`
  - `operationType`
  - `targetType`
  - `targetRef`
  - `contextSummary`
  - `riskLevel`
  - `outcome`
  - `occurredAt`
  - `relatedTrailId`
- **Validation Rules**:
  - 关键操作必须明确主体、目标和结果。
  - 高风险操作必须具备风险等级和上下文摘要。
- **Relationships**:
  - 多个 `Key Operation Trace` 可归属同一 `Governance Risk Event`。
  - `Key Operation Trace` 可被 `Governance Report Package` 聚合。

## 3. Cross-Team Authorization Snapshot

- **Purpose**: 保存某一时间点的跨团队授权分布情况。
- **Core Fields**:
  - `id`
  - `snapshotAt`
  - `sourceTeam`
  - `targetTeam`
  - `grantType`
  - `scopeSummary`
  - `temporality`
  - `delegationFlag`
  - `riskHint`
  - `trendLabel`
- **Validation Rules**:
  - 必须明确授权来源团队和目标团队。
  - 临时授权或委派授权必须显式标识时效或委派属性。
- **Relationships**:
  - 多个 `Cross-Team Authorization Snapshot` 可被同一 `Governance Coverage Snapshot` 引用。
  - `Cross-Team Authorization Snapshot` 可被 `Governance Report Package` 统计。

## 4. Governance Risk Event

- **Purpose**: 记录一次高风险访问、越权迹象或长期治理异常。
- **Core Fields**:
  - `id`
  - `riskType`
  - `severity`
  - `subjectSummary`
  - `scopeSummary`
  - `triggerReason`
  - `status`
  - `recommendedAction`
  - `owner`
  - `firstSeenAt`
  - `lastSeenAt`
- **Validation Rules**:
  - 必须有明确风险类型和严重级别。
  - 如果状态未关闭，必须有建议动作或责任人。
- **Relationships**:
  - 一个 `Governance Risk Event` 可关联多个 `Permission Change Trail` 与 `Key Operation Trace`。
  - 一个 `Governance Risk Event` 可进入 `Governance Action Item`。

## 5. Governance Coverage Snapshot

- **Purpose**: 表示某个时间点的治理覆盖率与状态分类结果。
- **Core Fields**:
  - `id`
  - `snapshotAt`
  - `coverageDomain`
  - `coverageRate`
  - `statusBreakdown`
  - `missingReasonSummary`
  - `confidenceLevel`
  - `trendSummary`
  - `owner`
- **Validation Rules**:
  - `coverageRate` 必须配合状态分类解释。
  - `confidenceLevel` 不能为空。
- **Relationships**:
  - 一个 `Governance Coverage Snapshot` 可关联多个 `Cross-Team Authorization Snapshot`。
  - 一个 `Governance Coverage Snapshot` 可被多个 `Governance Report Package` 引用。

## 6. Governance Report Package

- **Purpose**: 表示一次面向管理汇报、审计复核或客户交付生成的标准化报表包。
- **Core Fields**:
  - `id`
  - `reportType`
  - `title`
  - `audienceType`
  - `timeRange`
  - `summarySection`
  - `detailSection`
  - `attachmentCatalog`
  - `visibilityPolicy`
  - `generatedAt`
  - `generatedBy`
  - `status`
- **Validation Rules**:
  - `reportType` 必须属于管理汇报、审计复核或客户交付之一。
  - 报表必须包含摘要区块和适用对象说明。
- **Relationships**:
  - 一个 `Governance Report Package` 可聚合多个 `Permission Change Trail`、`Governance Risk Event` 和 `Governance Coverage Snapshot`。
  - 一个 `Governance Report Package` 可关联多个 `Export Record`。

## 7. Export Record

- **Purpose**: 记录一次报表、清单或交付材料导出行为。
- **Core Fields**:
  - `id`
  - `packageId`
  - `exportType`
  - `audienceScope`
  - `contentLevel`
  - `exportedAt`
  - `exportedBy`
  - `result`
  - `auditReference`
- **Validation Rules**:
  - 导出必须明确对象范围和内容层级。
  - 导出失败时必须记录结果说明。
- **Relationships**:
  - 多个 `Export Record` 归属于同一 `Governance Report Package`。

## 8. Delivery Artifact

- **Purpose**: 表示一份可纳入产品化交付包的标准材料。
- **Core Fields**:
  - `id`
  - `artifactType`
  - `title`
  - `versionScope`
  - `environmentScope`
  - `ownerRole`
  - `updatedAt`
  - `status`
  - `applicabilityNote`
- **Validation Rules**:
  - 必须标记适用版本和适用环境。
  - 必须明确责任角色。
- **Relationships**:
  - 一个 `Delivery Artifact` 可归属多个 `Delivery Readiness Bundle`。

## 9. Delivery Readiness Bundle

- **Purpose**: 表示面向某个客户、环境或交付场景的一套标准化交付包。
- **Core Fields**:
  - `id`
  - `name`
  - `targetEnvironment`
  - `targetAudience`
  - `artifactSummary`
  - `checklistStatus`
  - `missingItems`
  - `readinessConclusion`
  - `updatedAt`
- **Validation Rules**:
  - 必须至少包含材料摘要和就绪结论。
  - 如存在缺失项，必须列出缺失内容。
- **Relationships**:
  - 一个 `Delivery Readiness Bundle` 可关联多个 `Delivery Artifact`。
  - 一个 `Delivery Readiness Bundle` 可关联一个 `Delivery Checklist`。

## 10. Delivery Checklist

- **Purpose**: 记录产品化交付与客户验收的标准检查项执行情况。
- **Core Fields**:
  - `id`
  - `bundleId`
  - `checkItem`
  - `category`
  - `owner`
  - `evidenceRequirement`
  - `status`
  - `completedAt`
  - `remark`
- **Validation Rules**:
  - 每个检查项必须指定责任人和证据要求。
  - 已完成项必须可追溯到完成时间或说明。
- **Relationships**:
  - 多个 `Delivery Checklist` 项归属于同一 `Delivery Readiness Bundle`。
  - 缺失项可进入 `Governance Action Item`。

## 11. Governance Action Item

- **Purpose**: 汇总高风险访问、覆盖缺失和交付缺项形成的统一治理待办。
- **Core Fields**:
  - `id`
  - `sourceType`
  - `sourceRef`
  - `title`
  - `priority`
  - `owner`
  - `dueAt`
  - `status`
  - `resolutionSummary`
- **Validation Rules**:
  - 必须明确来源类型和优先级。
  - 已关闭待办必须具备处理结论。
- **Relationships**:
  - `Governance Action Item` 可来源于 `Governance Risk Event`、`Governance Coverage Snapshot` 或 `Delivery Checklist`。

## 12. State Notes

- `Governance Risk Event.status`: `open -> reviewing -> mitigating -> closed`
- `Governance Report Package.status`: `draft -> generating -> ready -> exported -> archived`
- `Delivery Artifact.status`: `draft -> active -> superseded -> archived`
- `Delivery Readiness Bundle.readinessConclusion`: `not-ready | conditionally-ready | ready`
- `Governance Action Item.status`: `open -> in-progress -> blocked -> resolved -> closed`
