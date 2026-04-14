# Tasks: 多集群 GitOps 与应用发布中心

**Input**: Design documents from `/specs/004-gitops-and-release/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/, quickstart.md

**Tests**: 本特性明确要求测试任务，必须覆盖后端 contract/integration 测试以及前端 Vitest 页面测试。

**Organization**: 任务按用户故事组织，确保 US1、US2、US3 可以独立实现与独立验收。

**Constitutional Gates**: 功能分支校验、数据库备份证据、国内源配置、中文 PR、远程推送和明确的用户合并授权必须在实现与交付时满足。

## Format: `[ID] [P?] [Story] Description`

- **[P]**: 可并行执行（不同文件、无未完成依赖）
- **[Story]**: 对应的用户故事标签（US1、US2、US3）
- 描述中包含明确文件路径，便于直接执行

## Phase 0: Governance Gates

**Purpose**: 在真正开始实现 004 之前完成宪章要求的治理门槛

- [X] T001 在 `artifacts/004-gitops-and-release/branch-check.txt` 记录当前分支、远程仓库和“未经用户批准不得合并”的门槛
- [X] T002 在 `artifacts/004-gitops-and-release/backup-manifest.txt` 核对并补充 004 本轮数据库备份、凭据差异和恢复抽样验证结果
- [X] T003 在 `artifacts/004-gitops-and-release/mirror-and-remote-check.txt` 记录 `GOPROXY=https://goproxy.cn,direct`、`https://registry.npmmirror.com`、Git/Helm/OCI 联调镜像来源和 `git@github.com:baihua19941101/kbManage.git` PR 流程
- [X] T004 更新 `specs/004-gitops-and-release/spec.md`、`specs/004-gitops-and-release/plan.md` 和 `specs/004-gitops-and-release/quickstart.md` 的执行状态，注明已进入 tasks 阶段

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: 为 004 建立最小模块骨架、配置入口和路由占位

- [X] T005 创建后端 GitOps 模块骨架 `backend/internal/service/gitops/`、`backend/internal/api/handler/gitops_handler.go`、`backend/internal/api/router/gitops_routes.go`
- [X] T006 [P] 创建交付适配层基础目录 `backend/internal/integration/delivery/git/`、`backend/internal/integration/delivery/helm/`、`backend/internal/integration/delivery/diff/`
- [X] T007 [P] 创建前端模块骨架 `frontend/src/features/gitops/`、`frontend/src/services/gitops.ts` 和路由占位到 `frontend/src/app/router.tsx`
- [X] T008 在 `backend/config/config.example.yaml`、`backend/config/config.dev.yaml`、`frontend/.env.example`、`frontend/.env.development` 和 `README.md` 增加 `gitops.sources.*`、`gitops.sync.*`、`gitops.diff.*`、`gitops.release.*`、`gitops.audit.*` 配置说明

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: 搭建 004 的共享数据模型、权限语义、路由、缓存和执行基座，阻塞所有用户故事

**⚠️ CRITICAL**: US1、US2、US3 都必须在本阶段完成后才能开始

- [X] T009 新增 004 数据库迁移 `backend/migrations/0006_gitops_release_core.sql`，落库交付来源、目标组、环境阶段、配置覆盖、交付单元、发布修订和交付动作表
- [X] T010 [P] 在 `backend/internal/domain/gitops.go` 定义 `DeliverySource`、`ClusterTargetGroup`、`EnvironmentStage`、`ConfigurationOverlay`、`ApplicationDeliveryUnit`、`ReleaseRevision`、`DeliveryOperation`
- [X] T011 [P] 创建仓储实现 `backend/internal/repository/delivery_source_repository.go`、`backend/internal/repository/cluster_target_group_repository.go`、`backend/internal/repository/delivery_unit_repository.go`、`backend/internal/repository/release_revision_repository.go`、`backend/internal/repository/delivery_operation_repository.go`
- [X] T012 [P] 在 `backend/internal/service/auth/permission_service.go`、`backend/internal/api/middleware/authorization.go` 定义 004 权限语义 `gitops:read`、`gitops:manage-source`、`gitops:sync`、`gitops:promote`、`gitops:rollback`、`gitops:override`
- [X] T013 [P] 在 `backend/internal/api/router/gitops_routes.go`、`backend/internal/api/router/router.go` 注册 004 API 路由骨架
- [X] T014 [P] 在 `backend/internal/repository/redis.go`、`backend/internal/service/gitops/progress_cache.go`、`backend/internal/service/gitops/diff_cache.go`、`backend/internal/service/gitops/lock_service.go` 建立动作进度、差异缓存和分布式锁基础设施
- [X] T015 [P] 在 `backend/internal/integration/delivery/git/client.go`、`backend/internal/integration/delivery/helm/client.go`、`backend/internal/integration/delivery/diff/comparator.go` 建立来源访问、发布源访问和差异比较抽象
- [X] T016 在 `backend/internal/service/gitops/service.go`、`backend/internal/service/gitops/scope_service.go` 建立统一 GitOps 授权、范围过滤和动作入口
- [X] T017 [P] 在 `frontend/src/services/api/types.ts`、`frontend/src/services/api/client.ts`、`frontend/src/app/queryClient.ts` 增加 004 共享类型、查询 key 和错误归一化
- [X] T018 在 `frontend/src/app/AuthorizedMenu.tsx`、`frontend/src/app/ProtectedRoute.tsx`、`frontend/src/features/auth/store.ts` 接入 GitOps 导航入口和基础权限门控

**Checkpoint**: 004 的共享模型、权限、路由、缓存和适配器骨架就绪，用户故事可以开始实现

---

## Phase 3: User Story 1 - 统一交付源与多集群目标建模 (Priority: P1) 🎯 MVP

**Goal**: 提供交付来源接入、目标组复用、环境分层、配置覆盖和交付单元建模，形成最小可用 GitOps 读路径闭环

**Independent Test**: 接入至少一个代码仓库和一个发布来源后，授权用户能够创建一个应用交付单元，为其绑定多个目标集群或集群组、环境层级和配置覆盖，并看到每个目标的期望状态、最近同步结果和当前漂移状态。

### Tests for User Story 1

- [X] T019 [P] [US1] 编写后端契约测试 `backend/tests/contract/gitops_sources_contract_test.go`、`backend/tests/contract/gitops_target_groups_contract_test.go`、`backend/tests/contract/gitops_delivery_units_contract_test.go`
- [X] T020 [P] [US1] 编写后端集成测试 `backend/tests/integration/gitops_sources_test.go`、`backend/tests/integration/gitops_delivery_unit_modeling_test.go`
- [X] T021 [P] [US1] 编写前端 Vitest 页面测试 `frontend/src/features/gitops/pages/GitOpsOverviewPage.test.tsx`、`frontend/src/features/gitops/pages/DeliveryUnitDetailPage.test.tsx`、`frontend/src/features/gitops/components/SourceFormDrawer.test.tsx`

### Implementation for User Story 1

- [X] T022 [P] [US1] 实现交付来源服务 `backend/internal/service/gitops/source_service.go`，覆盖来源创建、更新、启停和连通性校验
- [X] T023 [P] [US1] 实现目标组与环境阶段服务 `backend/internal/service/gitops/target_group_service.go`、`backend/internal/service/gitops/environment_service.go`
- [X] T024 [P] [US1] 实现配置覆盖合成与最终生效配置预览 `backend/internal/service/gitops/overlay_service.go`
- [X] T025 [P] [US1] 实现交付单元建模与状态聚合读取 `backend/internal/service/gitops/delivery_unit_service.go`、`backend/internal/service/gitops/status_service.go`
- [X] T026 [US1] 在 `backend/internal/api/handler/gitops_handler.go` 和 `backend/internal/api/router/gitops_routes.go` 落地 `/gitops/sources`、`/gitops/target-groups`、`/gitops/delivery-units`、`/gitops/delivery-units/{unitId}/status`
- [X] T027 [P] [US1] 实现前端服务层 `frontend/src/services/gitops.ts`
- [X] T028 [P] [US1] 实现来源管理页和交付单元列表页 `frontend/src/features/gitops/pages/GitOpsOverviewPage.tsx`、`frontend/src/features/gitops/components/SourceFormDrawer.tsx`、`frontend/src/features/gitops/components/TargetGroupDrawer.tsx`
- [X] T029 [US1] 实现交付单元详情页、环境阶段编辑和配置覆盖展示 `frontend/src/features/gitops/pages/DeliveryUnitDetailPage.tsx`、`frontend/src/features/gitops/components/EnvironmentStageEditor.tsx`、`frontend/src/features/gitops/components/OverlaySummaryPanel.tsx`
- [X] T030 [US1] 在 `frontend/src/app/router.tsx`、`frontend/src/features/resources/pages/ResourcesPage.tsx` 打通到 004 GitOps 中心的入口和上下文跳转

**Checkpoint**: US1 完成后，平台应具备可独立演示的 GitOps 建模与状态查看 MVP

---

## Phase 4: User Story 2 - 发布生命周期与多环境推进 (Priority: P1)

**Goal**: 提供同步、安装、升级、暂停、恢复、回滚、卸载和环境推进能力，形成 GitOps 写路径闭环

**Independent Test**: 选择一个已建好的应用交付单元后，授权用户能够对其执行首次安装、版本升级、配置升级、暂停同步、恢复同步、按环境推进和版本回滚，并查看每一步的结果、失败原因和历史记录。

### Tests for User Story 2

- [X] T031 [P] [US2] 编写后端契约测试 `backend/tests/contract/gitops_actions_contract_test.go`、`backend/tests/contract/gitops_revisions_contract_test.go`、`backend/tests/contract/gitops_diff_contract_test.go`
- [X] T032 [P] [US2] 编写后端集成测试 `backend/tests/integration/gitops_sync_execution_test.go`、`backend/tests/integration/gitops_promotion_flow_test.go`、`backend/tests/integration/gitops_rollback_test.go`
- [X] T033 [P] [US2] 编写前端 Vitest 页面测试 `frontend/src/features/gitops/components/ReleaseActionDrawer.test.tsx`、`frontend/src/features/gitops/components/PromotionTimeline.test.tsx`、`frontend/src/features/gitops/components/RevisionHistoryPanel.test.tsx`

### Implementation for User Story 2

- [X] T034 [P] [US2] 实现交付动作执行器 `backend/internal/service/gitops/executor.go`，覆盖 `install`、`sync`、`resync`、`upgrade`、`pause`、`resume`、`uninstall`
- [X] T035 [P] [US2] 实现环境推进与阶段执行编排 `backend/internal/service/gitops/promotion_service.go`、`backend/internal/service/gitops/stage_execution_service.go`
- [X] T036 [P] [US2] 实现发布修订历史与回滚目标识别 `backend/internal/service/gitops/revision_service.go`
- [X] T037 [P] [US2] 实现差异与漂移读取链路 `backend/internal/service/gitops/diff_service.go`、`backend/internal/integration/delivery/diff/comparator.go`
- [X] T038 [US2] 在 `backend/internal/service/gitops/service.go`、`backend/internal/worker/delivery_operation_worker.go` 打通动作状态流转、部分成功归一化和幂等处理
- [X] T039 [US2] 在 `backend/internal/api/handler/gitops_handler.go`、`backend/internal/api/router/gitops_routes.go` 落地 `/gitops/delivery-units/{unitId}/diff`、`/gitops/delivery-units/{unitId}/actions`、`/gitops/delivery-units/{unitId}/releases`、`/gitops/operations/{operationId}`
- [X] T040 [P] [US2] 实现前端发布动作提交与轮询 `frontend/src/features/gitops/components/ReleaseActionDrawer.tsx`、`frontend/src/features/gitops/hooks/useDeliveryOperation.ts`
- [X] T041 [P] [US2] 实现差异/漂移面板和发布历史面板 `frontend/src/features/gitops/components/DiffSummaryPanel.tsx`、`frontend/src/features/gitops/components/RevisionHistoryPanel.tsx`
- [X] T042 [US2] 实现环境推进时间线与回滚交互 `frontend/src/features/gitops/components/PromotionTimeline.tsx`、`frontend/src/features/gitops/components/RollbackDialog.tsx`

**Checkpoint**: US2 完成后，平台应具备真实可追踪的 GitOps 发布生命周期与多环境推进链路

---

## Phase 5: User Story 3 - 权限隔离与发布审计闭环 (Priority: P2)

**Goal**: 对来源管理、交付单元、环境推进、发布动作和回滚统一执行授权校验，并形成完整审计闭环

**Independent Test**: 为两个不同工作空间或环境范围的用户分别授权后，他们只能看到并操作各自范围内的交付来源、应用交付单元、环境推进和发布历史；审计人员能够按时间、操作者、对象和结果检索完整发布记录。

### Tests for User Story 3

- [X] T043 [P] [US3] 编写后端契约测试 `backend/tests/contract/gitops_access_control_contract_test.go`、`backend/tests/contract/gitops_audit_contract_test.go`
- [X] T044 [P] [US3] 编写后端集成测试 `backend/tests/integration/gitops_scope_authorization_test.go`、`backend/tests/integration/gitops_audit_test.go`
- [X] T045 [P] [US3] 编写前端 Vitest 页面测试 `frontend/src/features/gitops/pages/GitOpsAccessGate.test.tsx`

### Implementation for User Story 3

- [X] T046 [P] [US3] 在 `backend/internal/service/auth/scope_authorizer.go`、`backend/internal/service/gitops/scope_service.go` 落地来源、交付单元、环境和目标组的范围校验
- [X] T047 [P] [US3] 在 `backend/internal/service/audit/event_writer.go`、`backend/internal/service/audit/service.go` 增加 `gitops.*` 审计动作类型，覆盖来源校验、同步、推进、回滚、暂停、恢复和卸载
- [X] T048 [US3] 在 `backend/internal/api/middleware/authorization.go`、`backend/internal/api/router/gitops_routes.go` 区分只读、来源管理、同步、环境推进、回滚和配置覆盖变更权限
- [X] T049 [P] [US3] 在 `frontend/src/app/AuthorizedMenu.tsx`、`frontend/src/features/auth/store.ts` 实现 004 导航和动作级门控
- [X] T050 [US3] 在 `frontend/src/features/gitops/pages/DeliveryUnitDetailPage.tsx`、`frontend/src/features/gitops/components/ReleaseActionDrawer.tsx`、`frontend/src/features/gitops/components/RollbackDialog.tsx` 实现未授权空态、环境级禁用态和权限回收处理
- [X] T051 [US3] 实现 GitOps 审计查询页面与筛选交互 `frontend/src/features/audit/pages/GitOpsAuditPage.tsx`、`frontend/src/services/audit.ts`

**Checkpoint**: US3 完成后，004 的所有 GitOps 与发布能力都应运行在现有租户隔离和审计模型之下

---

## Final Phase: Polish & Cross-Cutting Concerns

**Purpose**: 收尾文档、验证、性能与交付准备

- [X] T052 [P] 清理共享类型与命名，在 `backend/internal/service/gitops/`、`backend/internal/integration/delivery/`、`frontend/src/features/gitops/`、`frontend/src/services/gitops.ts` 做收敛
- [X] T053 [P] 刷新配置与启动文档，在 `README.md`、`backend/config/config.example.yaml`、`backend/config/config.dev.yaml`、`frontend/.env.example`、`frontend/.env.development` 补齐 004 说明
- [X] T054 [P] 记录验证基线到 `artifacts/004-gitops-and-release/verification.md`、`artifacts/004-gitops-and-release/quickstart-validation.md`、`artifacts/004-gitops-and-release/repro-gitops-smoke.sh`
- [X] T055 在 `artifacts/004-gitops-and-release/pr-summary.md`、`artifacts/004-gitops-and-release/pr-readiness.md` 准备中文 PR 摘要、备份证据、测试说明和风险清单

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 0: Governance Gates**: 无依赖，但 004 真正开始实现前必须完成
- **Phase 1: Setup**: 依赖 Phase 0，负责建立模块骨架和配置入口
- **Phase 2: Foundational**: 依赖 Phase 1，阻塞所有用户故事
- **Phase 3: US1**: 依赖 Phase 2，是 004 的 MVP
- **Phase 4: US2**: 依赖 Phase 2，最佳实践是在 US1 的交付单元建模和状态读取稳定后推进
- **Phase 5: US3**: 依赖 US1 和 US2，因为权限隔离和审计需要覆盖来源、交付单元、动作和回滚链路
- **Final Phase**: 依赖目标用户故事完成

### User Story Dependencies

- **US1**: 仅依赖共享基础设施，是最小可演示范围
- **US2**: 依赖共享基础设施；会复用 US1 的交付单元、环境阶段和状态读取基座
- **US3**: 依赖 US1、US2，因为权限隔离和审计需要覆盖已存在的来源、差异、动作、推进和回滚链路

### Parallel Opportunities

- Phase 1 中，`T006`、`T007`、`T008` 可并行
- Phase 2 中，`T010`、`T011`、`T012`、`T013`、`T014`、`T015`、`T017` 可并行
- US1 中，`T019`、`T020`、`T021` 可并行；`T022`、`T023`、`T024`、`T025` 可并行；`T027`、`T028`、`T029` 可并行
- US2 中，`T031`、`T032`、`T033` 可并行；`T034`、`T035`、`T036`、`T037` 可并行；`T040`、`T041`、`T042` 可并行
- US3 中，`T043`、`T044`、`T045` 可并行；`T046`、`T047`、`T049` 可并行
- Final Phase 中，`T052`、`T053`、`T054` 可并行

---

## Parallel Example: User Story 1

```bash
# 并行启动 US1 的测试任务
Task: "T019 [US1] 后端契约测试"
Task: "T020 [US1] 后端集成测试"
Task: "T021 [US1] 前端 Vitest 页面测试"

# 并行实现 US1 的后端能力
Task: "T022 [US1] 交付来源服务"
Task: "T023 [US1] 目标组与环境阶段服务"
Task: "T024 [US1] 配置覆盖合成服务"
Task: "T025 [US1] 交付单元建模与状态聚合"

# 并行实现 US1 的前端视图
Task: "T027 [US1] 前端服务层"
Task: "T028 [US1] 来源管理页和交付单元列表页"
Task: "T029 [US1] 交付单元详情页与环境配置编辑"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. 完成 Phase 0: Governance Gates
2. 完成 Phase 1: Setup
3. 完成 Phase 2: Foundational
4. 完成 Phase 3: US1
5. **STOP and VALIDATE**：按 US1 独立验收标准验证来源接入、目标建模、环境分层、配置覆盖和状态聚合

### Incremental Delivery

1. 完成 Governance + Setup + Foundational，形成 004 共享底座
2. 交付 US1，形成可演示的 GitOps 建模与状态查看 MVP
3. 交付 US2，补齐同步、发布、推进、回滚和暂停/恢复写路径闭环
4. 交付 US3，补齐 GitOps 权限隔离与发布审计闭环
5. 完成 Final Phase，准备 PR 证据与交付文档

### Notes

- 首期范围严格限定在 GitOps 持续交付与应用发布生命周期管理
- 通用 CI 流水线编排、制品仓库管理、终端运维、策略准入和合规扫描禁止混入本任务清单
- 如果后续需要扩展审批流、制品生命周期治理或策略门禁，应进入新的 feature，而不是继续扩大 004 范围
