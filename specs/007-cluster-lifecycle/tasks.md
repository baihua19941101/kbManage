# Tasks: 多集群 Kubernetes 集群生命周期中心

**Input**: Design documents from `/specs/007-cluster-lifecycle/`  
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/openapi.yaml, quickstart.md

**Tests**: 本特性包含明确的独立验收标准、权限隔离、创建前校验和审计闭环要求，任务清单包含后端契约/集成测试与前端页面测试任务。  
**Organization**: 任务按用户故事分组，保证每个故事可独立实现、独立验证。

**Constitutional Gates**: 必须满足功能分支、数据库备份证据、国内依赖源配置、中文 PR、远程推送、用户同意后合并。

## Format: `[ID] [P?] [Story] Description`

- `[P]`: 可并行执行（不同文件、无前置依赖）
- `[Story]`: 任务归属用户故事（US1/US2/US3）
- 每条任务必须包含明确文件路径

## Phase 0: Governance Gates

**Purpose**: 完成宪章门槛与实施前证据准备

- [X] T001 在 `artifacts/007-cluster-lifecycle/branch-check.txt` 记录当前分支、禁止在 `main/master` 开发、`006` 已合并完成以及“未获用户同意不得合并”门槛
- [X] T002 执行 007 开发前数据库备份并在 `artifacts/007-cluster-lifecycle/backup-manifest.txt` 记录命令、时间戳、产物路径与恢复抽样验证结果
- [X] T003 在 `artifacts/007-cluster-lifecycle/mirror-and-remote-check.txt` 记录 `GOPROXY=https://goproxy.cn,direct`、`https://registry.npmmirror.com`、驱动/联调镜像来源与 `git@github.com:baihua19941101/kbManage.git` PR 流程

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: 建立 007 模块骨架、配置入口和导航占位

- [X] T004 创建后端集群生命周期模块骨架 `backend/internal/service/clusterlifecycle/`、`backend/internal/api/handler/cluster_lifecycle_handler.go`、`backend/internal/api/router/cluster_lifecycle_routes.go`
- [X] T005 [P] 创建驱动适配层目录与占位 `backend/internal/integration/clusterlifecycle/`、`backend/internal/integration/clusterlifecycle/driver/provider.go`、`backend/internal/integration/clusterlifecycle/validator/provider.go`
- [X] T006 [P] 创建前端模块骨架 `frontend/src/features/cluster-lifecycle/`、`frontend/src/services/clusterLifecycle.ts` 与路由占位到 `frontend/src/app/router.tsx`
- [X] T007 在 `backend/config/config.example.yaml`、`backend/config/config.dev.yaml`、`frontend/.env.example`、`frontend/.env.development`、`README.md` 增加 `clusterLifecycle.*` 配置说明
- [X] T008 在 `frontend/src/app/AuthorizedMenu.tsx`、`frontend/src/app/ProtectedRoute.tsx`、`frontend/src/features/auth/store.ts` 接入集群生命周期入口门禁占位

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: 完成所有用户故事共享且阻塞性的基础能力

**⚠️ CRITICAL**: US1/US2/US3 必须在本阶段完成后才可开始

- [X] T009 新增 007 数据库迁移 `backend/migrations/0009_cluster_lifecycle_core.sql`，落库集群生命周期记录、驱动版本、能力矩阵、模板、升级计划、节点池快照与生命周期动作/审计表
- [X] T010 [P] 在 `backend/internal/domain/cluster_lifecycle.go` 定义 `ClusterLifecycleRecord`、`ClusterDriverVersion`、`CapabilityMatrixEntry`、`ClusterTemplate`、`LifecycleOperation`、`UpgradePlan`、`NodePoolProfile`、`LifecycleAuditEvent`
- [X] T011 [P] 创建仓储实现 `backend/internal/repository/cluster_lifecycle_repository.go`、`backend/internal/repository/cluster_driver_repository.go`、`backend/internal/repository/cluster_template_repository.go`、`backend/internal/repository/cluster_capability_repository.go`、`backend/internal/repository/upgrade_plan_repository.go`、`backend/internal/repository/node_pool_repository.go`
- [X] T012 [P] 在 `backend/internal/integration/clusterlifecycle/driver/provider.go`、`backend/internal/integration/clusterlifecycle/validator/provider.go` 定义驱动操作抽象、校验结果模型和错误归一化语义
- [X] T013 在 `backend/internal/service/clusterlifecycle/service.go`、`backend/internal/service/clusterlifecycle/scope_service.go`、`backend/internal/service/auth/scope_authorizer.go` 建立统一生命周期授权与范围过滤入口
- [X] T014 [P] 在 `backend/internal/repository/redis.go`、`backend/internal/service/clusterlifecycle/progress_cache.go`、`backend/internal/service/clusterlifecycle/validation_cache.go`、`backend/internal/service/clusterlifecycle/operation_lock.go` 建立动作进度缓存、校验缓存与互斥锁协调
- [X] T015 [P] 在 `backend/internal/service/auth/permission_service.go`、`backend/internal/api/middleware/authorization.go` 增加 007 权限语义 `clusterlifecycle:read`、`clusterlifecycle:import`、`clusterlifecycle:create`、`clusterlifecycle:upgrade`、`clusterlifecycle:manage-nodepool`、`clusterlifecycle:retire`、`clusterlifecycle:manage-driver`
- [X] T016 在 `backend/internal/api/router/cluster_lifecycle_routes.go`、`backend/internal/api/router/router.go` 注册 007 API 路由骨架
- [X] T017 [P] 在 `frontend/src/services/api/types.ts`、`frontend/src/services/api/client.ts`、`frontend/src/app/queryClient.ts` 增加 007 共享类型、查询 key 和错误归一化
- [X] T018 在 `backend/internal/service/audit/event_writer.go`、`backend/internal/service/audit/service.go` 预置 `clusterlifecycle.*` 审计动作类型与查询维度映射

**Checkpoint**: 基础能力完成，可开始用户故事实现

---

## Phase 3: User Story 1 - 导入与注册已有集群 (Priority: P1) 🎯 MVP

**Goal**: 提供已有集群导入、注册引导、接入状态跟踪和基础生命周期详情能力。  
**Independent Test**: 导入一个已有集群并注册一个待纳管集群后，管理员能够看到列表、详情、接入状态、版本、健康摘要与失败原因。

### Tests for User Story 1

- [X] T019 [P] [US1] 编写后端契约测试 `backend/tests/contract/cluster_lifecycle_clusters_contract_test.go`、`backend/tests/contract/cluster_lifecycle_import_contract_test.go`、`backend/tests/contract/cluster_lifecycle_register_contract_test.go`
- [X] T020 [P] [US1] 编写后端集成测试 `backend/tests/integration/cluster_lifecycle_import_test.go`、`backend/tests/integration/cluster_lifecycle_register_test.go`、`backend/tests/integration/cluster_lifecycle_scope_authorization_test.go`
- [X] T021 [P] [US1] 编写前端 Vitest 页面测试 `frontend/src/features/cluster-lifecycle/pages/ClusterLifecycleListPage.test.tsx`、`frontend/src/features/cluster-lifecycle/pages/ClusterRegistrationPage.test.tsx`、`frontend/src/features/cluster-lifecycle/pages/ClusterLifecycleDetailPage.test.tsx`

### Implementation for User Story 1

- [X] T022 [P] [US1] 实现导入与注册服务 `backend/internal/service/clusterlifecycle/import_service.go`、`backend/internal/service/clusterlifecycle/registration_service.go`
- [X] T023 [P] [US1] 实现集群生命周期查询服务 `backend/internal/service/clusterlifecycle/cluster_query_service.go`，覆盖列表、详情、状态聚合与失败原因展示
- [X] T024 [US1] 在 `backend/internal/api/handler/cluster_lifecycle_handler.go`、`backend/internal/api/router/cluster_lifecycle_routes.go` 落地 `/cluster-lifecycle/clusters`、`/cluster-lifecycle/clusters/import`、`/cluster-lifecycle/clusters/register`、`/cluster-lifecycle/clusters/{clusterId}`
- [X] T025 [US1] 在 `backend/internal/service/clusterlifecycle/scope_service.go`、`backend/internal/api/middleware/authorization.go`、`backend/internal/service/auth/scope_authorizer.go` 落地导入/注册读取路径的范围过滤和权限回收处理
- [X] T026 [P] [US1] 实现前端服务层 `frontend/src/services/clusterLifecycle.ts`，覆盖列表、导入、注册和详情查询接口
- [X] T027 [P] [US1] 实现集群生命周期列表与导入页面 `frontend/src/features/cluster-lifecycle/pages/ClusterLifecycleListPage.tsx`、`frontend/src/features/cluster-lifecycle/components/ImportClusterDrawer.tsx`、`frontend/src/features/cluster-lifecycle/components/ClusterLifecycleTable.tsx`
- [X] T028 [P] [US1] 实现注册与详情页面 `frontend/src/features/cluster-lifecycle/pages/ClusterRegistrationPage.tsx`、`frontend/src/features/cluster-lifecycle/pages/ClusterLifecycleDetailPage.tsx`、`frontend/src/features/cluster-lifecycle/components/RegistrationGuideCard.tsx`
- [X] T029 [US1] 在 `frontend/src/app/router.tsx`、`frontend/src/features/clusters/pages/ClusterOverviewPage.tsx` 打通从现有集群总览进入生命周期中心的导航入口

**Checkpoint**: US1 完整可测，可作为 007 MVP 交付

---

## Phase 4: User Story 2 - 创建、升级与退役集群 (Priority: P1)

**Goal**: 提供模板化创建、创建前校验、升级计划、节点池调整、停用与退役闭环。  
**Independent Test**: 使用一个模板创建集群后，管理员能够执行创建前校验、升级计划、节点池扩缩和停用/退役流程，并看到完整状态与审计链路。

### Tests for User Story 2

- [X] T030 [P] [US2] 编写后端契约测试 `backend/tests/contract/cluster_lifecycle_create_contract_test.go`、`backend/tests/contract/cluster_lifecycle_validation_contract_test.go`、`backend/tests/contract/cluster_lifecycle_upgrade_contract_test.go`、`backend/tests/contract/cluster_lifecycle_node_pool_contract_test.go`、`backend/tests/contract/cluster_lifecycle_retire_contract_test.go`
- [X] T031 [P] [US2] 编写后端集成测试 `backend/tests/integration/cluster_lifecycle_create_flow_test.go`、`backend/tests/integration/cluster_lifecycle_upgrade_flow_test.go`、`backend/tests/integration/cluster_lifecycle_node_pool_conflict_test.go`、`backend/tests/integration/cluster_lifecycle_retire_flow_test.go`
- [X] T032 [P] [US2] 编写前端 Vitest 页面测试 `frontend/src/features/cluster-lifecycle/pages/ClusterProvisionPage.test.tsx`、`frontend/src/features/cluster-lifecycle/pages/ClusterUpgradePage.test.tsx`、`frontend/src/features/cluster-lifecycle/pages/NodePoolPage.test.tsx`、`frontend/src/features/cluster-lifecycle/pages/ClusterRetirementPage.test.tsx`

### Implementation for User Story 2

- [X] T033 [P] [US2] 实现创建与校验服务 `backend/internal/service/clusterlifecycle/provision_service.go`、`backend/internal/service/clusterlifecycle/validation_service.go`
- [X] T034 [P] [US2] 实现升级计划与执行服务 `backend/internal/service/clusterlifecycle/upgrade_service.go`、`backend/internal/worker/cluster_lifecycle_upgrade_worker.go`
- [X] T035 [P] [US2] 实现节点池管理与冲突锁服务 `backend/internal/service/clusterlifecycle/node_pool_service.go`、`backend/internal/service/clusterlifecycle/operation_lock.go`
- [X] T036 [P] [US2] 实现停用与退役服务 `backend/internal/service/clusterlifecycle/retirement_service.go`、`backend/internal/worker/cluster_lifecycle_retire_worker.go`
- [X] T037 [US2] 在 `backend/internal/api/handler/cluster_lifecycle_handler.go`、`backend/internal/api/router/cluster_lifecycle_routes.go` 落地 `/cluster-lifecycle/clusters` POST、`/cluster-lifecycle/clusters/{clusterId}/validate`、`/cluster-lifecycle/clusters/{clusterId}/upgrade-plans`、`/cluster-lifecycle/clusters/{clusterId}/upgrade-plans/{planId}/execute`、`/cluster-lifecycle/clusters/{clusterId}/node-pools/{nodePoolId}/scale`、`/cluster-lifecycle/clusters/{clusterId}/disable`、`/cluster-lifecycle/clusters/{clusterId}/retire`
- [X] T038 [US2] 在 `backend/internal/service/audit/event_writer.go`、`backend/internal/service/audit/service.go` 打通创建、校验、升级、节点池调整和退役的审计写入与查询聚合
- [X] T039 [P] [US2] 扩展前端服务与 hooks `frontend/src/services/clusterLifecycle.ts`、`frontend/src/features/cluster-lifecycle/hooks/useLifecycleAction.ts`
- [X] T040 [P] [US2] 实现创建与升级页面 `frontend/src/features/cluster-lifecycle/pages/ClusterProvisionPage.tsx`、`frontend/src/features/cluster-lifecycle/components/ClusterTemplateDrawer.tsx`、`frontend/src/features/cluster-lifecycle/pages/ClusterUpgradePage.tsx`、`frontend/src/features/cluster-lifecycle/components/UpgradePlanDrawer.tsx`
- [X] T041 [P] [US2] 实现节点池与退役页面 `frontend/src/features/cluster-lifecycle/pages/NodePoolPage.tsx`、`frontend/src/features/cluster-lifecycle/components/NodePoolScaleDrawer.tsx`、`frontend/src/features/cluster-lifecycle/pages/ClusterRetirementPage.tsx`、`frontend/src/features/cluster-lifecycle/components/RetireClusterDrawer.tsx`
- [X] T042 [US2] 在 `frontend/src/app/ProtectedRoute.tsx`、`frontend/src/features/cluster-lifecycle/pages/ClusterProvisionPage.tsx`、`frontend/src/features/cluster-lifecycle/pages/ClusterRetirementPage.tsx` 落地动作级权限门控、冲突动作空态与权限回收后的状态处理

**Checkpoint**: US2 可独立验证创建、升级、节点池管理和退役闭环

---

## Phase 5: User Story 3 - 能力矩阵与驱动扩展管理 (Priority: P2)

**Goal**: 提供驱动版本管理、模板管理和能力矩阵比较与兼容状态查看能力。  
**Independent Test**: 准备两个驱动版本和多个模板后，管理员能够查看能力矩阵、维护驱动/模板，并在模板校验时看到兼容结论。

### Tests for User Story 3

- [X] T043 [P] [US3] 编写后端契约测试 `backend/tests/contract/cluster_lifecycle_drivers_contract_test.go`、`backend/tests/contract/cluster_lifecycle_templates_contract_test.go`、`backend/tests/contract/cluster_lifecycle_capabilities_contract_test.go`、`backend/tests/contract/cluster_lifecycle_audit_contract_test.go`
- [X] T044 [P] [US3] 编写后端集成测试 `backend/tests/integration/cluster_lifecycle_driver_management_test.go`、`backend/tests/integration/cluster_lifecycle_template_validation_test.go`、`backend/tests/integration/cluster_lifecycle_audit_query_test.go`
- [X] T045 [P] [US3] 编写前端 Vitest 页面测试 `frontend/src/features/cluster-lifecycle/pages/ClusterDriverPage.test.tsx`、`frontend/src/features/cluster-lifecycle/pages/ClusterTemplatePage.test.tsx`、`frontend/src/features/cluster-lifecycle/pages/CapabilityMatrixPage.test.tsx`、`frontend/src/features/audit/pages/ClusterLifecycleAuditPage.test.tsx`

### Implementation for User Story 3

- [X] T046 [P] [US3] 实现驱动与能力矩阵服务 `backend/internal/service/clusterlifecycle/driver_service.go`、`backend/internal/service/clusterlifecycle/capability_service.go`
- [X] T047 [P] [US3] 实现模板管理与兼容性校验服务 `backend/internal/service/clusterlifecycle/template_service.go`
- [X] T048 [US3] 在 `backend/internal/api/handler/cluster_lifecycle_handler.go`、`backend/internal/api/router/cluster_lifecycle_routes.go` 落地 `/cluster-lifecycle/drivers`、`/cluster-lifecycle/drivers/{driverId}/capabilities`、`/cluster-lifecycle/templates`、`/cluster-lifecycle/templates/{templateId}/validate`
- [X] T049 [US3] 在 `backend/internal/api/handler/audit_handler.go`、`backend/internal/api/router/audit_routes.go` 聚合并暴露 `/audit/cluster-lifecycle/events` 查询链路
- [X] T050 [P] [US3] 扩展前端服务与 hooks `frontend/src/services/clusterLifecycle.ts`、`frontend/src/features/cluster-lifecycle/hooks/useCapabilityMatrix.ts`
- [X] T051 [P] [US3] 实现驱动与模板页面 `frontend/src/features/cluster-lifecycle/pages/ClusterDriverPage.tsx`、`frontend/src/features/cluster-lifecycle/components/DriverVersionDrawer.tsx`、`frontend/src/features/cluster-lifecycle/pages/ClusterTemplatePage.tsx`、`frontend/src/features/cluster-lifecycle/components/ClusterTemplateFormDrawer.tsx`
- [X] T052 [P] [US3] 实现能力矩阵与审计页面 `frontend/src/features/cluster-lifecycle/pages/CapabilityMatrixPage.tsx`、`frontend/src/features/cluster-lifecycle/components/CapabilityMatrixTable.tsx`、`frontend/src/features/audit/pages/ClusterLifecycleAuditPage.tsx`
- [X] T053 [US3] 在 `frontend/src/app/AuthorizedMenu.tsx`、`frontend/src/app/router.tsx`、`frontend/src/features/cluster-lifecycle/pages/CapabilityMatrixPage.tsx` 落地驱动管理权限门控、筛选持久化和未授权空态

**Checkpoint**: US3 完成后形成驱动、模板、能力矩阵与审计闭环

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: 收敛质量、文档、验证证据与 PR 交付材料

- [X] T054 [P] 收敛命名与共享类型，在 `backend/internal/service/clusterlifecycle/`、`backend/internal/integration/clusterlifecycle/`、`frontend/src/features/cluster-lifecycle/`、`frontend/src/services/clusterLifecycle.ts` 清理重复字段与错误文案
- [X] T055 [P] 刷新配置与启动文档，在 `README.md`、`backend/config/config.example.yaml`、`backend/config/config.dev.yaml`、`frontend/.env.example`、`frontend/.env.development` 补齐 007 说明
- [X] T056 [P] 记录验证基线到 `artifacts/007-cluster-lifecycle/verification.md`、`artifacts/007-cluster-lifecycle/quickstart-validation.md`、`artifacts/007-cluster-lifecycle/repro-cluster-lifecycle-smoke.sh`
- [X] T057 在 `artifacts/007-cluster-lifecycle/pr-summary.md`、`artifacts/007-cluster-lifecycle/pr-readiness.md` 准备中文 PR 摘要、备份证据、测试说明、风险清单与用户合并确认项

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
- **US2 (P1)**: 依赖 US1 产出的集群主记录与接入状态对象，才能形成创建、升级与退役闭环
- **US3 (P2)**: 依赖 Foundational 的驱动/模板基础模型和 US2 的创建前校验语义；能力矩阵页面可在 US2 主干稳定后独立推进

### Parallel Opportunities

- **Phase 1**: T005/T006 可并行
- **Phase 2**: T010/T011/T012/T014/T015/T017 可并行
- **US1**: T019/T020/T021 并行，T022/T023 并行，T026/T027/T028 并行
- **US2**: T030/T031/T032 并行，T033/T034/T035/T036 并行，T039/T040/T041 并行
- **US3**: T043/T044/T045 并行，T046/T047 并行，T050/T051/T052 并行

---

## Parallel Example: User Story 1

```bash
# 并行测试任务
Task: "T019 [US1] backend/tests/contract/cluster_lifecycle_clusters_contract_test.go"
Task: "T020 [US1] backend/tests/integration/cluster_lifecycle_import_test.go"
Task: "T021 [US1] frontend/src/features/cluster-lifecycle/pages/ClusterLifecycleListPage.test.tsx"

# 并行实现任务
Task: "T022 [US1] backend/internal/service/clusterlifecycle/import_service.go"
Task: "T023 [US1] backend/internal/service/clusterlifecycle/cluster_query_service.go"
Task: "T026 [US1] frontend/src/services/clusterLifecycle.ts"
```

---

## Implementation Strategy

### MVP First (US1 Only)

1. 完成 Phase 0-2
2. 完成 US1（Phase 3）
3. 按 Independent Test 验证 US1
4. 产出阶段性演示或 PR 更新

### Incremental Delivery

1. 先交付 US1：导入与注册已有集群、生命周期详情与接入状态跟踪
2. 再交付 US2：模板化创建、创建前校验、升级计划、节点池管理和退役闭环
3. 最后交付 US3：驱动版本、模板管理、能力矩阵和生命周期审计汇报
4. 最终执行 Phase 6 文档与验证收尾

### Notes

- `[P]` 任务代表可并行，但仍需满足前置依赖
- 每个用户故事都可独立验收
- 所有提交说明与 PR 摘要必须为中文
- 未获用户明确同意前禁止合并
