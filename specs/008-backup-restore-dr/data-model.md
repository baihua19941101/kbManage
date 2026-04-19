# Data Model: 平台级备份恢复与灾备中心

## 1. BackupPolicy

### Purpose

统一表示一类平台对象或业务范围的备份规则，用于持续生成可恢复恢复点。

### Fields

- `id`: 主键标识
- `name`: 策略名称
- `description`: 策略说明
- `scopeType`: `platform-metadata | access-config | audit-records | cluster-config | namespace`
- `scopeRef`: 作用范围引用
- `executionMode`: `manual | scheduled`
- `scheduleExpression`: 计划执行描述，可为空
- `retentionRule`: 保留规则说明
- `consistencyLevel`: `best-effort | application-consistent | platform-consistent`
- `status`: `draft | active | paused | deprecated`
- `ownerUserId`: 责任人
- `createdAt`: 创建时间
- `updatedAt`: 更新时间

### Relationships

- 一对多关联 `RestorePoint`
- 一对多关联 `BackupAuditEvent`

### Validation Rules

- 同一作用范围内策略名称必须唯一
- `executionMode=scheduled` 时必须存在计划执行描述
- `status=paused` 时不可生成新的自动恢复点

## 2. RestorePoint

### Purpose

表示某次备份执行形成的可恢复快照或恢复目录项。

### Fields

- `id`: 主键标识
- `policyId`: 来源策略
- `scopeSnapshot`: 范围快照
- `backupStartedAt`: 开始时间
- `backupCompletedAt`: 完成时间
- `durationSeconds`: 执行耗时
- `result`: `succeeded | partially_succeeded | failed | expired`
- `consistencySummary`: 一致性说明
- `failureReason`: 失败原因，可为空
- `storageRef`: 外部介质引用
- `expiresAt`: 过期时间
- `createdBy`: 触发人或系统

### Relationships

- 多对一关联 `BackupPolicy`
- 一对多关联 `RestoreJob`

### Validation Rules

- `result=failed` 时不可作为恢复来源
- `expiresAt` 到期后不可继续恢复

## 3. RestoreJob

### Purpose

表示一次原地恢复、跨集群恢复、环境迁移或定向恢复动作。

### Fields

- `id`: 主键标识
- `restorePointId`: 来源恢复点
- `jobType`: `in-place-restore | cross-cluster-restore | environment-migration | selective-restore`
- `sourceEnvironment`: 源环境引用
- `targetEnvironment`: 目标环境引用
- `scopeSelection`: 选定恢复范围
- `conflictSummary`: 冲突检查摘要
- `consistencyNotice`: 一致性说明
- `status`: `pending | validating | running | partially_succeeded | succeeded | failed | canceled | blocked`
- `resultSummary`: 结果摘要
- `failureReason`: 失败原因
- `requestedBy`: 发起人
- `startedAt`: 开始时间
- `completedAt`: 完成时间

### Relationships

- 多对一关联 `RestorePoint`
- 一对多关联 `BackupAuditEvent`

### Validation Rules

- 恢复范围必须属于恢复点覆盖范围
- `status=blocked` 时必须保留冲突或一致性原因

## 4. MigrationPlan

### Purpose

表示一次环境迁移的源目标映射和切换方案。

### Fields

- `id`: 主键标识
- `name`: 迁移计划名称
- `sourceClusterId`: 源集群
- `targetClusterId`: 目标集群
- `scopeSelection`: 迁移范围
- `mappingRules`: 目标映射规则
- `cutoverSteps`: 切换步骤摘要
- `status`: `draft | approved | running | succeeded | failed | canceled`
- `createdBy`: 创建人
- `createdAt`: 创建时间
- `updatedAt`: 更新时间

### Relationships

- 可关联一个或多个 `RestoreJob`

### Validation Rules

- 源环境和目标环境不能相同
- `status=running` 时不可编辑映射规则

## 5. DRDrillPlan

### Purpose

表示一份可周期执行的灾备演练计划。

### Fields

- `id`: 主键标识
- `name`: 演练计划名称
- `description`: 演练说明
- `scopeSelection`: 演练范围
- `rpoTargetMinutes`: RPO 目标
- `rtoTargetMinutes`: RTO 目标
- `roleAssignments`: 参与角色与责任说明
- `cutoverProcedure`: 切换步骤摘要
- `validationChecklistRef`: 验证清单引用
- `status`: `draft | active | paused | retired`
- `createdBy`: 创建人
- `createdAt`: 创建时间
- `updatedAt`: 更新时间

### Relationships

- 一对多关联 `DRDrillRecord`

### Validation Rules

- `rpoTargetMinutes` 和 `rtoTargetMinutes` 必须大于零
- 激活中的计划必须具备切换步骤和验证清单

## 6. DRDrillRecord

### Purpose

表示一次灾备演练的执行过程和结果。

### Fields

- `id`: 主键标识
- `planId`: 来源演练计划
- `startedAt`: 开始时间
- `completedAt`: 完成时间
- `actualRpoMinutes`: 实际 RPO
- `actualRtoMinutes`: 实际 RTO
- `status`: `pending | running | succeeded | failed | partially_succeeded | canceled`
- `stepResults`: 步骤结果快照
- `validationResults`: 验证清单结果
- `incidentNotes`: 异常点记录
- `executedBy`: 发起人

### Relationships

- 多对一关联 `DRDrillPlan`
- 一对一关联 `DRDrillReport`

### Validation Rules

- `status=succeeded` 时必须具备步骤和验证结果
- 若实际 RPO 或 RTO 超标，必须记录偏差说明

## 7. DRDrillReport

### Purpose

表示一次灾备演练后的总结报告和改进建议。

### Fields

- `id`: 主键标识
- `drillRecordId`: 来源演练记录
- `goalAssessment`: 目标达成情况
- `gapSummary`: 偏差说明
- `issuesFound`: 问题项列表
- `improvementActions`: 改进建议
- `publishedAt`: 发布时间
- `publishedBy`: 发布人

### Relationships

- 一对一关联 `DRDrillRecord`

### Validation Rules

- 报告发布前必须存在对应演练记录
- 报告必须至少包含目标达成情况和改进建议

## 8. BackupAuditEvent

### Purpose

表示围绕备份、恢复、迁移和演练产生的标准化审计记录。

### Fields

- `id`: 主键标识
- `action`: 动作类型
- `actorUserId`: 操作者
- `targetType`: `backup-policy | restore-point | restore-job | migration-plan | drill-plan | drill-record | drill-report`
- `targetRef`: 目标引用
- `scopeSnapshot`: 范围快照
- `outcome`: `succeeded | failed | blocked | canceled`
- `detailSnapshot`: 详情快照
- `occurredAt`: 发生时间

### Validation Rules

- 所有关键动作都必须生成对应审计事件
- `targetType` 与 `targetRef` 必须成对出现

## State Transitions

### BackupPolicy.status

- `draft -> active`
- `active -> paused`
- `paused -> active`
- `active -> deprecated`

### RestorePoint.result

- `succeeded -> expired`
- `partially_succeeded -> expired`

### RestoreJob.status

- `pending -> validating`
- `validating -> running`
- `validating -> blocked`
- `running -> succeeded`
- `running -> partially_succeeded`
- `running -> failed`
- `pending -> canceled`

### DRDrillPlan.status

- `draft -> active`
- `active -> paused`
- `paused -> active`
- `active -> retired`

### DRDrillRecord.status

- `pending -> running`
- `running -> succeeded`
- `running -> partially_succeeded`
- `running -> failed`
- `pending -> canceled`
