# 008 Quickstart Validation

## 场景 1：备份策略与恢复点

- 创建备份策略：通过 `/api/v1/backup-restore/policies` 成功返回 `201`
- 手动触发备份：通过 `/api/v1/backup-restore/policies/{policyId}/run` 成功返回 `202`
- 查询恢复点：通过 `/api/v1/backup-restore/restore-points` 与 `/api/v1/backup-restore/restore-points/{restorePointId}` 成功返回结果

## 场景 2：恢复与迁移

- 创建恢复任务：通过 `/api/v1/backup-restore/restore-jobs` 成功返回 `202`
- 恢复前校验：通过 `/api/v1/backup-restore/restore-jobs/{jobId}/validate` 返回校验摘要
- 创建迁移计划：通过 `/api/v1/backup-restore/migrations` 成功返回 `201`

## 场景 3：灾备演练与审计

- 创建演练计划：通过 `/api/v1/backup-restore/drills/plans` 成功返回 `201`
- 发起演练：通过 `/api/v1/backup-restore/drills/plans/{planId}/run` 成功返回 `202`
- 生成报告：通过 `/api/v1/backup-restore/drills/records/{recordId}/report` 成功返回 `201`
- 查询审计：通过 `/api/v1/audit/backup-restore/events` 成功返回事件列表

## 结论

- 008 首期范围内的备份策略、恢复点、恢复、迁移、演练、报告和审计链路均已具备最小可用实现。
- 当前可进入用户审查阶段；未完成项不在 008 主链路内。
