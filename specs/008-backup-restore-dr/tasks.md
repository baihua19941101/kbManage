# Tasks: 平台级备份恢复与灾备中心

**Input**: Design documents from `/specs/008-backup-restore-dr/`  
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/openapi.yaml, quickstart.md

**Tests**: 本特性包含明确的独立验收标准、恢复前校验、跨集群迁移和灾备演练报告要求，任务清单包含后端契约/集成测试与前端页面测试任务。  
**Organization**: 任务按用户故事分组，保证每个故事可独立实现、独立验证。

**Constitutional Gates**: 必须满足功能分支、数据库备份证据、国内依赖源配置、中文 PR、远程推送、用户同意后合并。

## Format: `[ID] [P?] [Story] Description`

- `[P]`: 可并行执行（不同文件、无前置依赖）
- `[Story]`: 任务归属用户故事（US1/US2/US3）
- 每条任务必须包含明确文件路径

## Phase 0: Governance Gates

**Purpose**: 完成宪章门槛与实施前证据准备

- [X] T001 在 `artifacts/008-backup-restore-dr/branch-check.txt` 记录当前分支、禁止在 `main/master` 开发、`007` 已合并完成以及“未获用户同意不得合并”门槛
- [X] T002 执行 008 开发前数据库备份并在 `artifacts/008-backup-restore-dr/backup-manifest.txt` 记录命令、时间戳、产物路径与恢复抽样验证结果
- [X] T003 在 `artifacts/008-backup-restore-dr/mirror-and-remote-check.txt` 记录 `GOPROXY=https://goproxy.cn,direct`、`https://registry.npmmirror.com`、备份执行器/对象存储联调镜像来源与 `git@github.com:baihua19941101/kbManage.git` PR 流程

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: 建立 008 模块骨架、配置入口和导航占位

- [X] T004 创建后端备份恢复模块骨架 `backend/internal/service/backuprestore/`、`backend/internal/api/handler/backup_restore_handler.go`、`backend/internal/api/router/backup_restore_routes.go`
- [X] T005 [P] 创建备份执行适配层目录与占位 `backend/internal/integration/backuprestore/`、`backend/internal/integration/backuprestore/executor/provider.go`、`backend/internal/integration/backuprestore/validator/provider.go`
- [X] T006 [P] 创建前端模块骨架 `frontend/src/features/backup-restore-dr/`、`frontend/src/services/backupRestore.ts` 与路由占位到 `frontend/src/app/router.tsx`
- [X] T007 在 `backend/config/config.example.yaml`、`backend/config/config.dev.yaml`、`frontend/.env.example`、`frontend/.env.development`、`README.md` 增加 `backupRestore.*` 配置说明
- [X] T008 在 `frontend/src/app/AuthorizedMenu.tsx`、`frontend/src/app/ProtectedRoute.tsx`、`frontend/src/features/auth/store.ts` 接入备份恢复中心入口门禁占位

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: 完成所有用户故事共享且阻塞性的基础能力

**⚠️ CRITICAL**: US1/US2/US3 必须在本阶段完成后才可开始

- [X] T009 新增 008 数据库迁移 `backend/migrations/0010_backup_restore_dr_core.sql`，落库备份策略、恢复点、恢复任务、迁移计划、演练计划、演练记录、演练报告与备份恢复审计表
- [X] T010 [P] 在 `backend/internal/domain/backup_restore.go` 定义 `BackupPolicy`、`RestorePoint`、`RestoreJob`、`MigrationPlan`、`DRDrillPlan`、`DRDrillRecord`、`DRDrillReport`、`BackupAuditEvent`
- [X] T011 [P] 创建仓储实现 `backend/internal/repository/backup_policy_repository.go`、`backend/internal/repository/restore_point_repository.go`、`backend/internal/repository/restore_job_repository.go`、`backend/internal/repository/migration_plan_repository.go`、`backend/internal/repository/dr_drill_plan_repository.go`、`backend/internal/repository/dr_drill_record_repository.go`、`backend/internal/repository/dr_drill_report_repository.go`
- [X] T012 [P] 在 `backend/internal/integration/backuprestore/executor/provider.go`、`backend/internal/integration/backuprestore/validator/provider.go` 定义备份执行抽象、恢复前校验结果模型和错误归一化语义
- [X] T013 在 `backend/internal/service/backuprestore/service.go`、`backend/internal/service/backuprestore/scope_service.go`、`backend/internal/service/auth/scope_authorizer.go` 建立统一备份恢复授权与范围过滤入口
- [X] T014 [P] 在 `backend/internal/repository/redis.go`、`backend/internal/service/backuprestore/progress_cache.go`、`backend/internal/service/backuprestore/precheck_cache.go`、`backend/internal/service/backuprestore/operation_lock.go` 建立备份/恢复/演练进度缓存、校验缓存与互斥锁协调
- [X] T015 [P] 在 `backend/internal/service/auth/permission_service.go`、`backend/internal/api/middleware/authorization.go` 增加 008 权限语义 `backuprestore:read`、`backuprestore:manage-policy`、`backuprestore:backup`、`backuprestore:restore`、`backuprestore:migrate`、`backuprestore:drill`
- [X] T016 在 `backend/internal/api/router/backup_restore_routes.go`、`backend/internal/api/router/router.go` 注册 008 API 路由骨架
- [X] T017 [P] 在 `frontend/src/services/api/types.ts`、`frontend/src/services/api/client.ts`、`frontend/src/app/queryClient.ts` 增加 008 共享类型、查询 key 和错误归一化
- [X] T018 在 `backend/internal/service/audit/event_writer.go`、`backend/internal/service/audit/service.go` 预置 `backuprestore.*` 审计动作类型与查询维度映射

**Checkpoint**: 基础能力完成，可开始用户故事实现

---

## Phase 3: User Story 1 - 定义平台备份策略并生成恢复点 (Priority: P1) 🎯 MVP

**Goal**: 提供受保护对象范围管理、备份策略、恢复点目录和备份结果可视化能力。  
**Independent Test**: 为至少两类对象配置备份策略并执行备份后，管理员能够看到恢复点、覆盖范围、保留规则、耗时、结果和失败原因。

### Tests for User Story 1

- [X] T019 [P] [US1] 编写后端契约测试 `backend/tests/contract/backup_policy_contract_test.go`、`backend/tests/contract/restore_point_contract_test.go`、`backend/tests/contract/manual_backup_run_contract_test.go`
- [X] T020 [P] [US1] 编写后端集成测试 `backend/tests/integration/backup_policy_lifecycle_test.go`、`backend/tests/integration/manual_backup_run_test.go`、`backend/tests/integration/restore_point_visibility_scope_test.go`
- [X] T021 [P] [US1] 编写前端 Vitest 页面测试 `frontend/src/features/backup-restore-dr/pages/BackupPolicyPage.test.tsx`、`frontend/src/features/backup-restore-dr/pages/RestorePointPage.test.tsx`

### Implementation for User Story 1

- [X] T022 [P] [US1] 实现备份策略与恢复点服务 `backend/internal/service/backuprestore/policy_service.go`、`backend/internal/service/backuprestore/restore_point_service.go`
- [X] T023 [P] [US1] 实现手动备份执行服务 `backend/internal/service/backuprestore/backup_run_service.go`、`backend/internal/worker/backup_run_worker.go`
- [X] T024 [US1] 在 `backend/internal/api/handler/backup_restore_handler.go`、`backend/internal/api/router/backup_restore_routes.go` 落地 `/backup-restore/policies`、`/backup-restore/policies/{policyId}/run`、`/backup-restore/restore-points`、`/backup-restore/restore-points/{restorePointId}`
- [X] T025 [US1] 在 `backend/internal/service/backuprestore/scope_service.go`、`backend/internal/api/middleware/authorization.go`、`backend/internal/service/auth/scope_authorizer.go` 落地策略、恢复点查询和手动备份路径的范围过滤
- [X] T026 [P] [US1] 实现前端服务层 `frontend/src/services/backupRestore.ts`，覆盖策略、手动备份和恢复点查询接口
- [X] T027 [P] [US1] 实现备份策略页面与表单 `frontend/src/features/backup-restore-dr/pages/BackupPolicyPage.tsx`、`frontend/src/features/backup-restore-dr/components/BackupPolicyDrawer.tsx`
- [X] T028 [P] [US1] 实现恢复点列表与详情页面 `frontend/src/features/backup-restore-dr/pages/RestorePointPage.tsx`、`frontend/src/features/backup-restore-dr/components/RestorePointDetailDrawer.tsx`
- [X] T029 [US1] 在 `frontend/src/app/router.tsx`、`frontend/src/features/resources/pages/ResourcesPage.tsx` 打通从平台资源视图进入备份恢复中心的导航入口

**Checkpoint**: US1 完整可测，可作为 008 MVP 交付

---

## Phase 4: User Story 2 - 执行恢复与跨集群迁移 (Priority: P1)

**Goal**: 提供原地恢复、跨集群恢复、环境迁移、定向恢复和恢复前校验闭环。  
**Independent Test**: 选择一个恢复点执行原地恢复、跨集群恢复或定向恢复后，操作者能够看到范围、目标环境、耗时、结果、失败原因和一致性说明。

### Tests for User Story 2

- [X] T030 [P] [US2] 编写后端契约测试 `backend/tests/contract/restore_job_contract_test.go`、`backend/tests/contract/restore_precheck_contract_test.go`、`backend/tests/contract/migration_plan_contract_test.go`
- [X] T031 [P] [US2] 编写后端集成测试 `backend/tests/integration/in_place_restore_flow_test.go`、`backend/tests/integration/cross_cluster_restore_flow_test.go`、`backend/tests/integration/selective_restore_scope_test.go`、`backend/tests/integration/migration_plan_flow_test.go`
- [X] T032 [P] [US2] 编写前端 Vitest 页面测试 `frontend/src/features/backup-restore-dr/pages/RestoreJobPage.test.tsx`、`frontend/src/features/backup-restore-dr/pages/MigrationPlanPage.test.tsx`

### Implementation for User Story 2

- [X] T033 [P] [US2] 实现恢复与迁移服务 `backend/internal/service/backuprestore/restore_service.go`、`backend/internal/service/backuprestore/migration_service.go`
- [X] T034 [P] [US2] 实现恢复前校验与一致性说明服务 `backend/internal/service/backuprestore/precheck_service.go`、`backend/internal/service/backuprestore/consistency_service.go`
- [X] T035 [P] [US2] 实现恢复/迁移动作工作器 `backend/internal/worker/restore_job_worker.go`、`backend/internal/worker/migration_job_worker.go`
- [X] T036 [US2] 在 `backend/internal/api/handler/backup_restore_handler.go`、`backend/internal/api/router/backup_restore_routes.go` 落地 `/backup-restore/restore-jobs`、`/backup-restore/restore-jobs/{jobId}/validate`、`/backup-restore/migrations`
- [X] T037 [US2] 在 `backend/internal/service/audit/event_writer.go`、`backend/internal/service/audit/service.go` 打通恢复、迁移、校验动作的审计写入与查询聚合
- [X] T038 [P] [US2] 扩展前端服务与 hooks `frontend/src/services/backupRestore.ts`、`frontend/src/features/backup-restore-dr/hooks/useRestoreAction.ts`
- [X] T039 [P] [US2] 实现恢复任务页面与表单 `frontend/src/features/backup-restore-dr/pages/RestoreJobPage.tsx`、`frontend/src/features/backup-restore-dr/components/RestoreJobDrawer.tsx`
- [X] T040 [P] [US2] 实现迁移计划页面与详情 `frontend/src/features/backup-restore-dr/pages/MigrationPlanPage.tsx`、`frontend/src/features/backup-restore-dr/components/MigrationPlanDrawer.tsx`
- [X] T041 [US2] 在 `frontend/src/app/ProtectedRoute.tsx`、`frontend/src/features/backup-restore-dr/pages/RestoreJobPage.tsx`、`frontend/src/features/backup-restore-dr/pages/MigrationPlanPage.tsx` 落地动作级权限门控、冲突阻断空态和一致性提示

**Checkpoint**: US2 可独立验证恢复、迁移和恢复前校验闭环

---

## Phase 5: User Story 3 - 管理灾备演练与验证报告 (Priority: P2)

**Goal**: 提供灾备演练计划、执行记录、RPO/RTO 达成情况、验证清单和演练报告能力。  
**Independent Test**: 创建一份演练计划并完成一次演练后，负责人能够查看演练记录、目标达成情况、步骤结果、验证清单完成度和演练报告。

### Tests for User Story 3

- [X] T042 [P] [US3] 编写后端契约测试 `backend/tests/contract/dr_drill_plan_contract_test.go`、`backend/tests/contract/dr_drill_record_contract_test.go`、`backend/tests/contract/dr_drill_report_contract_test.go`、`backend/tests/contract/backup_audit_contract_test.go`
- [X] T043 [P] [US3] 编写后端集成测试 `backend/tests/integration/dr_drill_plan_flow_test.go`、`backend/tests/integration/dr_drill_execution_test.go`、`backend/tests/integration/dr_report_generation_test.go`、`backend/tests/integration/backup_audit_query_test.go`
- [X] T044 [P] [US3] 编写前端 Vitest 页面测试 `frontend/src/features/backup-restore-dr/pages/DRDrillPlanPage.test.tsx`、`frontend/src/features/backup-restore-dr/pages/DRDrillRecordPage.test.tsx`、`frontend/src/features/backup-restore-dr/pages/DRDrillReportPage.test.tsx`、`frontend/src/features/audit/pages/BackupRestoreAuditPage.test.tsx`

### Implementation for User Story 3

- [X] T045 [P] [US3] 实现演练计划与记录服务 `backend/internal/service/backuprestore/drill_plan_service.go`、`backend/internal/service/backuprestore/drill_record_service.go`
- [X] T046 [P] [US3] 实现演练报告与目标评估服务 `backend/internal/service/backuprestore/drill_report_service.go`、`backend/internal/service/backuprestore/rpo_rto_service.go`
- [X] T047 [US3] 在 `backend/internal/api/handler/backup_restore_handler.go`、`backend/internal/api/router/backup_restore_routes.go` 落地 `/backup-restore/drills/plans`、`/backup-restore/drills/plans/{planId}/run`、`/backup-restore/drills/records/{recordId}`、`/backup-restore/drills/records/{recordId}/report`
- [X] T048 [US3] 在 `backend/internal/api/handler/audit_handler.go`、`backend/internal/api/router/audit_routes.go` 聚合并暴露 `/audit/backup-restore/events` 查询链路
- [X] T049 [P] [US3] 扩展前端服务与 hooks `frontend/src/services/backupRestore.ts`、`frontend/src/features/backup-restore-dr/hooks/useDrillAction.ts`
- [X] T050 [P] [US3] 实现演练计划与记录页面 `frontend/src/features/backup-restore-dr/pages/DRDrillPlanPage.tsx`、`frontend/src/features/backup-restore-dr/components/DRDrillPlanDrawer.tsx`、`frontend/src/features/backup-restore-dr/pages/DRDrillRecordPage.tsx`
- [X] T051 [P] [US3] 实现演练报告与审计页面 `frontend/src/features/backup-restore-dr/pages/DRDrillReportPage.tsx`、`frontend/src/features/backup-restore-dr/components/DRDrillReportDrawer.tsx`、`frontend/src/features/audit/pages/BackupRestoreAuditPage.tsx`
- [X] T052 [US3] 在 `frontend/src/app/AuthorizedMenu.tsx`、`frontend/src/app/router.tsx`、`frontend/src/features/backup-restore-dr/pages/DRDrillPlanPage.tsx` 落地演练管理权限门控、未授权空态和 RPO/RTO 偏差提示

**Checkpoint**: US3 完成后形成备份恢复、迁移与灾备演练审计闭环

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: 收敛质量、文档、验证证据与 PR 交付材料

- [X] T053 [P] 收敛命名与共享类型，在 `backend/internal/service/backuprestore/`、`backend/internal/integration/backuprestore/`、`frontend/src/features/backup-restore-dr/`、`frontend/src/services/backupRestore.ts` 清理重复字段与错误文案
- [X] T054 [P] 刷新配置与启动文档，在 `README.md`、`backend/config/config.example.yaml`、`backend/config/config.dev.yaml`、`frontend/.env.example`、`frontend/.env.development` 补齐 008 说明
- [X] T055 [P] 记录验证基线到 `artifacts/008-backup-restore-dr/verification.md`、`artifacts/008-backup-restore-dr/quickstart-validation.md`、`artifacts/008-backup-restore-dr/repro-backup-restore-smoke.sh`
- [X] T056 在 `artifacts/008-backup-restore-dr/pr-summary.md`、`artifacts/008-backup-restore-dr/pr-readiness.md` 准备中文 PR 摘要、备份证据、测试说明、风险清单与用户合并确认项

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 0 (Governance Gates)**: 无依赖，必须先完成
- **Phase 1 (Setup)**: 依赖 Phase 0
- **Phase 2 (Foundational)**: 依赖 Phase 1，阻塞所有用户故事
- **Phase 3/4/5 (User Stories)**: 均依赖 Phase 2；US1 作为 MVP 优先，US2 在 US1 主干稳定后推进，US3 在 US1/US2 的基础数据与动作语义完成后推进
- **Phase 6 (Polish)**: 依赖已完成的用户故事范围
- **Release / Merge**: 依赖远程推送、PR 更新、评审完成与用户明确同意

### User Story Dependencies

- **US1 (P1)**: 无用户故事前置依赖，Foundational 完成后可立即开始
- **US2 (P1)**: 依赖 US1 产出的恢复点目录与保护范围对象，才能形成恢复、迁移与定向恢复闭环
- **US3 (P2)**: 依赖 Foundational 的演练与审计基础模型，以及 US1/US2 的恢复点和恢复动作语义

### Parallel Opportunities

- **Phase 1**: T005/T006 可并行
- **Phase 2**: T010/T011/T012/T014/T015/T017 可并行
- **US1**: T019/T020/T021 并行，T022/T023 并行，T026/T027/T028 并行
- **US2**: T030/T031/T032 并行，T033/T034/T035 并行，T038/T039/T040 并行
- **US3**: T042/T043/T044 并行，T045/T046 并行，T049/T050/T051 并行

---

## Parallel Example: User Story 1

```bash
# 并行测试任务
Task: "T019 [US1] backend/tests/contract/backup_policy_contract_test.go"
Task: "T020 [US1] backend/tests/integration/backup_policy_lifecycle_test.go"
Task: "T021 [US1] frontend/src/features/backup-restore-dr/pages/BackupPolicyPage.test.tsx"

# 并行实现任务
Task: "T022 [US1] backend/internal/service/backuprestore/policy_service.go"
Task: "T023 [US1] backend/internal/service/backuprestore/backup_run_service.go"
Task: "T026 [US1] frontend/src/services/backupRestore.ts"
```

---

## Implementation Strategy

### MVP First (US1 Only)

1. 完成 Phase 0-2
2. 完成 US1（Phase 3）
3. 按 Independent Test 验证 US1
4. 产出阶段性演示或 PR 更新

### Incremental Delivery

1. 先交付 US1：备份策略、手动备份和恢复点目录
2. 再交付 US2：原地恢复、跨集群恢复、环境迁移、定向恢复和恢复前校验
3. 最后交付 US3：灾备演练计划、演练记录、报告和备份恢复审计
4. 最终执行 Phase 6 文档与验证收尾

### Notes

- `[P]` 任务代表可并行，但仍需满足前置依赖
- 每个用户故事都可独立验收
- 所有提交说明与 PR 摘要必须为中文
- 未获用户明确同意前禁止合并
