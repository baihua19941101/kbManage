# 006 Quickstart 验证

1. 创建基线：已通过 contract 覆盖 `/api/v1/compliance/baselines`
2. 创建扫描配置：已通过 contract 覆盖 `/api/v1/compliance/scan-profiles`
3. 触发扫描并查看失败项/证据：已通过 contract + integration 覆盖 `/api/v1/compliance/scans`、`/api/v1/compliance/findings`
4. 创建整改任务、例外与复检：已通过 contract + integration 覆盖 `/api/v1/compliance/remediation-tasks`、`/api/v1/compliance/exceptions`、`/api/v1/compliance/rechecks`
5. 查看总览、趋势、归档与审计：已通过 contract 覆盖 `/api/v1/compliance/overview`、`/api/v1/compliance/trends`、`/api/v1/compliance/archive-exports`、`/api/v1/audit/compliance/events`
