# Tasks: 平台 SRE 与规模化治理

**Input**: Design documents from `/specs/011-sre-scale/`  
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/openapi.yaml, quickstart.md

**Tests**: 本特性包含明确的独立验收标准、高可用治理、升级闭环和规模化容量治理要求，任务清单包含后端契约/集成测试与前端页面测试任务。  
**Organization**: 任务按用户故事分组，保证每个故事可独立实现、独立验证。

**Constitutional Gates**: 必须满足功能分支、数据库备份证据、国内依赖源配置、中文 PR、远程推送、用户同意后合并。

## Format: `[ID] [P?] [Story] Description`

- `[P]`: 可并行执行（不同文件、无前置依赖）
- `[Story]`: 任务归属用户故事（US1/US2/US3）
- 每条任务必须包含明确文件路径

## Phase 0: Governance Gates

**Purpose**: 完成宪章门槛与实施前证据准备

- [X] T001 在 `artifacts/011-sre-scale/branch-check.txt` 记录当前分支、禁止在 `main/master` 开发、`010` 已合并完成以及“未获用户同意不得合并”门槛
- [X] T002 执行 011 开发前数据库备份并在 `artifacts/011-sre-scale/backup-manifest.txt` 记录命令、时间戳、产物路径与恢复抽样验证结果
- [X] T003 在 `artifacts/011-sre-scale/mirror-and-remote-check.txt` 记录 `GOPROXY=https://goproxy.cn,direct`、`https://registry.npmmirror.com`、SRE 相关依赖镜像策略与 `git@github.com:baihua19941101/kbManage.git` PR 流程

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: 建立 011 模块骨架、配置入口和导航占位

- [X] T004 创建后端 SRE 模块骨架 `backend/internal/service/sre/`、`backend/internal/api/handler/sre_handler.go`、`backend/internal/api/router/sre_routes.go`
- [X] T005 [P] 创建平台健康、升级与容量证据适配层目录 `backend/internal/integration/sre/`、`backend/internal/integration/sre/health_provider.go`、`backend/internal/integration/sre/upgrade_validator.go`、`backend/internal/integration/sre/scale_analyzer.go`
- [X] T006 [P] 创建前端模块骨架 `frontend/src/features/sre-scale/`、`frontend/src/services/sreScale.ts` 与路由占位到 `frontend/src/app/router.tsx`
- [X] T007 在 `backend/config/config.example.yaml`、`backend/config/config.dev.yaml`、`frontend/.env.example`、`frontend/.env.development`、`README.md` 增加 `sreScale.*` 配置说明
- [X] T008 在 `frontend/src/app/AuthorizedMenu.tsx`、`frontend/src/app/ProtectedRoute.tsx`、`frontend/src/features/auth/store.ts` 接入 SRE 与规模化治理入口门禁占位

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: 完成所有用户故事共享且阻塞性的基础能力

**⚠️ CRITICAL**: US1/US2/US3 必须在本阶段完成后才可开始

- [X] T009 新增 011 数据库迁移 `backend/migrations/0013_sre_scale_core.sql`，落库高可用策略、维护窗口、平台健康快照、容量基线、升级计划、回退验证、运行手册、告警基线、规模化证据与审计表
- [X] T010 [P] 在 `backend/internal/domain/sre_scale.go` 定义 `HAPolicy`、`MaintenanceWindow`、`PlatformHealthSnapshot`、`CapacityBaseline`、`UpgradePlan`、`RollbackValidation`、`RunbookArticle`、`AlertBaseline`、`ScaleEvidence`
- [X] T011 [P] 创建仓储实现 `backend/internal/repository/ha_policy_repository.go`、`backend/internal/repository/maintenance_window_repository.go`、`backend/internal/repository/platform_health_snapshot_repository.go`、`backend/internal/repository/capacity_baseline_repository.go`、`backend/internal/repository/upgrade_plan_repository.go`、`backend/internal/repository/rollback_validation_repository.go`、`backend/internal/repository/runbook_article_repository.go`、`backend/internal/repository/alert_baseline_repository.go`、`backend/internal/repository/scale_evidence_repository.go`
- [X] T012 [P] 在 `backend/internal/integration/sre/health_provider.go`、`backend/internal/integration/sre/upgrade_validator.go`、`backend/internal/integration/sre/scale_analyzer.go` 定义健康聚合、升级前检查、容量预测与压测证据接入抽象
- [X] T013 在 `backend/internal/service/sre/service.go`、`backend/internal/service/sre/scope_service.go`、`backend/internal/service/auth/scope_authorizer.go` 建立 SRE 范围过滤、平台级资源访问和治理入口
- [X] T014 [P] 在 `backend/internal/repository/redis.go`、`backend/internal/service/sre/health_cache.go`、`backend/internal/service/sre/upgrade_coordinator.go`、`backend/internal/service/sre/scale_cache.go` 建立健康缓存、升级协调、容量预测缓存与幂等键
- [X] T015 [P] 在 `backend/internal/service/auth/permission_service.go`、`backend/internal/api/middleware/authorization.go` 增加 011 权限语义 `sre:read`、`sre:manage-ha`、`sre:manage-upgrade`、`sre:manage-scale`
- [X] T016 在 `backend/internal/api/router/sre_routes.go`、`backend/internal/api/router/router.go` 注册 011 API 路由骨架
- [X] T017 [P] 在 `frontend/src/services/api/types.ts`、`frontend/src/services/api/client.ts`、`frontend/src/app/queryClient.ts` 增加 011 共享类型、查询 key 和错误归一化
- [X] T018 在 `backend/internal/service/audit/event_writer.go`、`backend/internal/service/audit/service.go` 预置 `sre.*` 审计动作类型与查询维度映射

**Checkpoint**: 基础能力完成，可开始用户故事实现

---

## Phase 3: User Story 1 - 建立平台高可用与运维基线 (Priority: P1) 🎯 MVP

**Goal**: 提供高可用策略、维护窗口、平台健康总览与异常恢复记录能力。  
**Independent Test**: 配置高可用策略和维护窗口后，可独立查看平台组件健康、依赖状态、任务积压、容量风险、切换状态和恢复结果。

### Tests for User Story 1

- [X] T019 [P] [US1] 编写后端契约测试 `backend/tests/contract/sre_ha_policy_contract_test.go`、`backend/tests/contract/sre_health_overview_contract_test.go`
- [X] T020 [P] [US1] 编写后端集成测试 `backend/tests/integration/sre_ha_failover_flow_test.go`、`backend/tests/integration/sre_maintenance_window_flow_test.go`
- [X] T021 [P] [US1] 编写前端 Vitest 页面测试 `frontend/src/features/sre-scale/pages/HAControlPage.test.tsx`、`frontend/src/features/sre-scale/pages/HealthOverviewPage.test.tsx`

### Implementation for User Story 1

- [X] T022 [P] [US1] 实现高可用策略与维护窗口服务 `backend/internal/service/sre/ha_policy_service.go`、`backend/internal/service/sre/maintenance_window_service.go`
- [X] T023 [P] [US1] 实现平台健康总览与恢复摘要服务 `backend/internal/service/sre/health_overview_service.go`、`backend/internal/service/sre/recovery_summary_service.go`
- [X] T024 [US1] 在 `backend/internal/api/handler/sre_handler.go`、`backend/internal/api/router/sre_routes.go` 落地 `/sre/ha-policies`、`/sre/health/overview`、`/sre/maintenance-windows`
- [X] T025 [US1] 在 `backend/internal/service/sre/scope_service.go`、`backend/internal/api/middleware/authorization.go`、`backend/internal/service/auth/scope_authorizer.go` 落地平台高可用与健康视图的范围过滤
- [X] T026 [P] [US1] 实现前端服务层 `frontend/src/services/sreScale.ts`，覆盖高可用策略、维护窗口和健康总览接口
- [X] T027 [P] [US1] 实现高可用治理页面与表单 `frontend/src/features/sre-scale/pages/HAControlPage.tsx`、`frontend/src/features/sre-scale/components/HAPolicyDrawer.tsx`
- [X] T028 [P] [US1] 实现平台健康总览页面 `frontend/src/features/sre-scale/pages/HealthOverviewPage.tsx`、`frontend/src/features/sre-scale/components/HealthStatusCard.tsx`
- [X] T029 [US1] 在 `frontend/src/app/router.tsx`、`frontend/src/features/sre-scale/pages/HealthOverviewPage.tsx` 打通 SRE 总览导航、异常分类展示和维护窗口空态

**Checkpoint**: US1 完整可测，可作为 011 MVP 交付

---

## Phase 4: User Story 2 - 安全执行平台升级与回退验证 (Priority: P1)

**Goal**: 提供升级前检查、滚动升级、升级后验收和回退验证闭环。  
**Independent Test**: 准备一组可升级版本后，可独立完成升级前检查、发起升级计划、查看滚动阶段并登记回退验证结果。

### Tests for User Story 2

- [X] T030 [P] [US2] 编写后端契约测试 `backend/tests/contract/sre_upgrade_precheck_contract_test.go`、`backend/tests/contract/sre_upgrade_plan_contract_test.go`、`backend/tests/contract/sre_rollback_validation_contract_test.go`
- [X] T031 [P] [US2] 编写后端集成测试 `backend/tests/integration/sre_upgrade_rollout_flow_test.go`、`backend/tests/integration/sre_rollback_validation_flow_test.go`
- [X] T032 [P] [US2] 编写前端 Vitest 页面测试 `frontend/src/features/sre-scale/pages/UpgradeGovernancePage.test.tsx`、`frontend/src/features/sre-scale/pages/RollbackValidationPage.test.tsx`

### Implementation for User Story 2

- [X] T033 [P] [US2] 实现升级前检查与升级计划服务 `backend/internal/service/sre/upgrade_precheck_service.go`、`backend/internal/service/sre/upgrade_plan_service.go`
- [X] T034 [P] [US2] 实现滚动阶段协调与回退验证服务 `backend/internal/service/sre/upgrade_rollout_service.go`、`backend/internal/service/sre/rollback_validation_service.go`
- [X] T035 [US2] 在 `backend/internal/api/handler/sre_handler.go`、`backend/internal/api/router/sre_routes.go` 落地 `/sre/upgrades/prechecks`、`/sre/upgrades`、`/sre/upgrades/{upgradeId}/rollback-validations`
- [X] T036 [US2] 在 `backend/internal/service/audit/event_writer.go`、`backend/internal/service/audit/service.go` 打通升级前检查、升级执行、暂停、回退验证和验收动作的审计写入与查询聚合
- [X] T037 [P] [US2] 扩展前端服务与 hooks `frontend/src/services/sreScale.ts`、`frontend/src/features/sre-scale/hooks/useUpgradeAction.ts`
- [X] T038 [P] [US2] 实现升级治理页面与表单 `frontend/src/features/sre-scale/pages/UpgradeGovernancePage.tsx`、`frontend/src/features/sre-scale/components/UpgradePlanDrawer.tsx`
- [X] T039 [P] [US2] 实现回退验证与升级后验收页面 `frontend/src/features/sre-scale/pages/RollbackValidationPage.tsx`、`frontend/src/features/sre-scale/components/RollbackValidationDrawer.tsx`
- [X] T040 [US2] 在 `frontend/src/app/router.tsx`、`frontend/src/features/sre-scale/pages/UpgradeGovernancePage.tsx` 落地阻断项展示、阶段进度说明和升级暂停/回退提示

**Checkpoint**: US2 可独立验证平台升级与回退治理闭环

---

## Phase 5: User Story 3 - 管理规模化性能与容量治理 (Priority: P2)

**Goal**: 提供容量基线、趋势分析、压测证据、自诊断、运行手册与告警基线治理能力。  
**Independent Test**: 导入容量基线、压测样本和预测数据后，可独立查看容量风险、瓶颈摘要、可信度说明，并关联运行手册与告警基线。

### Tests for User Story 3

- [X] T041 [P] [US3] 编写后端契约测试 `backend/tests/contract/sre_capacity_baseline_contract_test.go`、`backend/tests/contract/sre_scale_evidence_contract_test.go`、`backend/tests/contract/sre_runbook_contract_test.go`
- [X] T042 [P] [US3] 编写后端集成测试 `backend/tests/integration/sre_capacity_forecast_flow_test.go`、`backend/tests/integration/sre_runbook_linkage_flow_test.go`
- [X] T043 [P] [US3] 编写前端 Vitest 页面测试 `frontend/src/features/sre-scale/pages/CapacityGovernancePage.test.tsx`、`frontend/src/features/sre-scale/pages/RunbookCenterPage.test.tsx`、`frontend/src/features/audit/pages/SREAuditPage.test.tsx`

### Implementation for User Story 3

- [X] T044 [P] [US3] 实现容量基线与规模化证据服务 `backend/internal/service/sre/capacity_baseline_service.go`、`backend/internal/service/sre/scale_evidence_service.go`
- [X] T045 [P] [US3] 实现运行手册、告警基线与自诊断摘要服务 `backend/internal/service/sre/runbook_service.go`、`backend/internal/service/sre/alert_baseline_service.go`、`backend/internal/service/sre/self_diagnosis_service.go`
- [X] T046 [US3] 在 `backend/internal/api/handler/sre_handler.go`、`backend/internal/api/router/sre_routes.go` 落地 `/sre/capacity/baselines`、`/sre/scale-evidence`、`/sre/runbooks`
- [X] T047 [US3] 在 `backend/internal/api/handler/audit_handler.go`、`backend/internal/api/router/audit_routes.go` 聚合并暴露 `/audit/sre/events` 查询链路
- [X] T048 [P] [US3] 扩展前端服务与 hooks `frontend/src/services/sreScale.ts`、`frontend/src/features/sre-scale/hooks/useScaleEvidenceAction.ts`
- [X] T049 [P] [US3] 实现容量与性能治理页面 `frontend/src/features/sre-scale/pages/CapacityGovernancePage.tsx`、`frontend/src/features/sre-scale/components/CapacityTrendChart.tsx`
- [X] T050 [P] [US3] 实现运行手册与自诊断页面 `frontend/src/features/sre-scale/pages/RunbookCenterPage.tsx`、`frontend/src/features/sre-scale/components/RunbookDrawer.tsx`、`frontend/src/features/sre-scale/components/SelfDiagnosisCard.tsx`
- [X] T051 [P] [US3] 实现 SRE 审计页面 `frontend/src/features/audit/pages/SREAuditPage.tsx`、`frontend/src/features/audit/pages/SREAuditPage.test.tsx`
- [X] T052 [US3] 在 `frontend/src/app/AuthorizedMenu.tsx`、`frontend/src/app/router.tsx`、`frontend/src/features/sre-scale/pages/CapacityGovernancePage.tsx` 落地可信度提示、运行手册关联和容量不足告警空态

**Checkpoint**: US3 完成后形成平台稳定性、升级治理与容量治理审计闭环

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: 收敛质量、文档、验证证据与 PR 交付材料

- [X] T053 [P] 收敛命名与共享类型，在 `backend/internal/service/sre/`、`backend/internal/integration/sre/`、`frontend/src/features/sre-scale/`、`frontend/src/services/sreScale.ts` 清理重复字段与错误文案
- [X] T054 [P] 刷新配置与启动文档，在 `README.md`、`backend/config/config.example.yaml`、`backend/config/config.dev.yaml`、`frontend/.env.example`、`frontend/.env.development` 补齐 011 说明
- [X] T055 [P] 记录验证基线到 `artifacts/011-sre-scale/verification.md`、`artifacts/011-sre-scale/quickstart-validation.md`、`artifacts/011-sre-scale/repro-sre-scale-smoke.sh`
- [X] T056 在 `artifacts/011-sre-scale/pr-summary.md`、`artifacts/011-sre-scale/pr-readiness.md` 准备中文 PR 摘要、备份证据、测试说明、风险清单与用户合并确认项

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 0 (Governance Gates)**: 无依赖，必须先完成
- **Phase 1 (Setup)**: 依赖 Phase 0
- **Phase 2 (Foundational)**: 依赖 Phase 1，阻塞所有用户故事
- **Phase 3/4/5 (User Stories)**: 均依赖 Phase 2；US1 作为 MVP 优先，US2 在 US1 主干稳定后推进，US3 在 US1/US2 的平台治理语义稳定后推进
- **Phase 6 (Polish)**: 依赖已完成的用户故事范围
- **Release / Merge**: 依赖远程推送、PR 更新、评审完成与用户明确同意

### User Story Dependencies

- **US1 (P1)**: 无用户故事前置依赖，Foundational 完成后可立即开始
- **US2 (P1)**: 依赖 US1 产出的高可用、维护窗口与平台健康治理语义，才能形成升级闭环
- **US3 (P2)**: 依赖 Foundational 的容量、证据与审计模型，以及 US1/US2 的平台运行状态与升级上下文

### Parallel Opportunities

- **Phase 1**: T005/T006 可并行
- **Phase 2**: T010/T011/T012/T014/T015/T017 可并行
- **US1**: T019/T020/T021 并行，T022/T023 并行，T026/T027/T028 并行
- **US2**: T030/T031/T032 并行，T033/T034 并行，T037/T038/T039 并行
- **US3**: T041/T042/T043 并行，T044/T045 并行，T048/T049/T050/T051 并行

---

## Parallel Example: User Story 1

```bash
# 并行测试任务
Task: "T019 [US1] backend/tests/contract/sre_ha_policy_contract_test.go"
Task: "T020 [US1] backend/tests/integration/sre_ha_failover_flow_test.go"
Task: "T021 [US1] frontend/src/features/sre-scale/pages/HAControlPage.test.tsx"

# 并行实现任务
Task: "T022 [US1] backend/internal/service/sre/ha_policy_service.go"
Task: "T023 [US1] backend/internal/service/sre/health_overview_service.go"
Task: "T026 [US1] frontend/src/services/sreScale.ts"
```

---

## Implementation Strategy

### MVP First (US1 Only)

1. 完成 Phase 0-2
2. 完成 US1（Phase 3）
3. 按 Independent Test 验证 US1
4. 产出阶段性演示或 PR 更新

### Incremental Delivery

1. 先交付 US1：高可用策略、维护窗口、平台健康总览与恢复摘要
2. 再交付 US2：升级前检查、滚动升级、升级后验收与回退验证
3. 最后交付 US3：容量基线、趋势、压测证据、运行手册、自诊断与 SRE 审计
4. 最终执行 Phase 6 文档与验证收尾

### Notes

- `[P]` 任务代表可并行，但仍需满足前置依赖
- 每个用户故事都可独立验收
- 所有提交说明与 PR 摘要必须为中文
- 未获用户明确同意前禁止合并
