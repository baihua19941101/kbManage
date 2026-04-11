# Tasks: 多集群 Kubernetes 可视化管理平台（真实后续任务清单）

**Input**: Existing implementation plus design documents from `/specs/001-k8s-ops-platform/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/openapi.yaml, quickstart.md

**Current Reality Check**:
- 当前仓库已有前后端骨架、基础路由、领域模型和部分页面；`spec.md` 已切换为 Follow-up 执行态，当前工作分支为专用 feature branch `001-k8s-ops-platform-followup`，后续任务聚焦剩余真实能力补齐。
- 后端 `go test ./...` 可通过，但现有契约/集成测试允许 `404` 与 `501`，更多是在校验接口形状而不是校验真实业务流程。
- 前端 `npm run lint` 与 `npm run build` 可通过，但 `npm run test -- --run` 结果为 “No test files found”。
- 集群接入、资源同步、工作空间/项目页面、操作中心、审计导出仍存在 `mock`、`stub` 或 `fallback` 实现，本清单仅覆盖这些剩余工作。

**Tests**: 接下来必须把后端 contract/integration 测试从 “允许占位实现” 收紧到 “验证真实行为”，并补齐前端关键页面测试。

**Organization**: 任务按“治理与基线 -> 共享基础 -> 用户故事”重新编排，只保留当前 001 的剩余工作。

## Format: `[ID] [P?] [Story] Description`

- **[P]**: 可并行执行（不同文件、无未完成依赖）
- **[Story]**: 用户故事标签，仅出现在用户故事阶段
- 每个任务都包含明确文件路径，保证可直接执行

## Phase 0: Governance & Re-Baseline

**Purpose**: 把当前实现从“框架初成”校正到符合 001 宪章和真实进度的可持续开发状态

- [x] T001 Record the dedicated feature-branch requirement and current repository state in artifacts/001-k8s-ops-platform/branch-check.txt and artifacts/001-k8s-ops-platform/current-baseline.md
- [x] T002 Refresh the pre-change MySQL backup evidence for the remaining schema work in artifacts/001-k8s-ops-platform/backup-manifest.txt and artifacts/001-k8s-ops-platform/mysql-backup-<timestamp>.sql
- [x] T003 Update execution status and next-step notes in specs/001-k8s-ops-platform/spec.md, specs/001-k8s-ops-platform/plan.md, and specs/001-k8s-ops-platform/quickstart.md, including fixed first-batch Kind/role matrix, high-risk confirmation flow, CSV-only audit export with masking, and performance evidence requirements

---

## Phase 1: Shared Quality & Foundation

**Purpose**: 先修复测试基线和共享调用链，避免继续在 mock/stub 上堆功能

- [x] T004 Create DB-backed backend test bootstrap helpers in backend/tests/testutil/app.go and backend/tests/testutil/auth.go
- [x] T005 [P] Tighten auth, cluster, access, operation, and audit contract tests to reject `404`/`501` fallback statuses in backend/tests/contract/auth_contract_test.go, backend/tests/contract/clusters_contract_test.go, backend/tests/contract/access_control_contract_test.go, backend/tests/contract/operations_contract_test.go, and backend/tests/contract/audit_contract_test.go
- [x] T006 [P] Replace permissive integration skeletons with seeded end-to-end assertions in backend/tests/integration/cluster_overview_test.go, backend/tests/integration/scope_authorization_test.go, backend/tests/integration/operation_execution_test.go, and backend/tests/integration/audit_query_test.go
- [x] T007 [P] Add first frontend Vitest coverage for login and cluster/resource screens in frontend/src/features/auth/pages/LoginPage.test.tsx, frontend/src/features/clusters/pages/ClusterOverviewPage.test.tsx, and frontend/src/features/resources/pages/ResourceListPage.test.tsx
- [x] T008 Implement shared frontend API client, auth header injection, and query error normalization in frontend/src/services/api/client.ts, frontend/src/features/auth/store.ts, and frontend/src/app/queryClient.ts

**Checkpoint**: 后续用户故事必须建立在“真实路由 + 真实测试 + 统一 API 客户端”之上

---

## Phase 2: User Story 1 - 多集群统一接入与资源总览 (Priority: P1) 🎯 MVP

**Goal**: 把集群接入、健康总览、资源索引和资源详情从 demo/mock 收口为可用能力

**Independent Test**: 接入至少 2 个集群后，管理员能够在统一视图中看到真实集群状态、同步时间、资源筛选结果和资源详情。

### Tests for User Story 1

- [x] T009 [P] [US1] Upgrade cluster and resource backend tests to require DB-backed routes and concrete payload assertions in backend/tests/contract/clusters_contract_test.go and backend/tests/integration/cluster_overview_test.go

### Implementation for User Story 1

- [x] T010 [P] [US1] Add missing cluster listing and overview response models in backend/internal/api/handler/cluster_handler.go, backend/internal/api/router/cluster_routes.go, and backend/internal/service/cluster/service.go
- [x] T011 [P] [US1] Implement kubeconfig parsing and real connectivity verification in backend/internal/service/cluster/service.go, backend/internal/kube/client/manager.go, and backend/internal/repository/cluster_credential_repository.go
- [x] T012 [P] [US1] Replace the noop resource sync with a real cluster inventory worker in backend/internal/kube/adapter/resource_indexer.go, backend/internal/worker/cluster_sync_worker.go, and backend/internal/repository/resource_inventory_repository.go (首期仅索引 `Deployment`、`StatefulSet`、`DaemonSet`、`Pod`、`Service`、`Ingress`、`Node`、`Namespace`)
- [x] T013 [US1] Persist sync status, last-success timestamp, and failure reason for clusters in backend/internal/domain/cluster.go, backend/migrations/0002_cluster_sync_status.sql, and backend/internal/repository/cluster_repository.go
- [x] T014 [US1] Expose resource detail and health-summary queries in backend/internal/api/handler/resource_handler.go, backend/internal/api/handler/cluster_handler.go, and backend/internal/api/router/cluster_routes.go
- [x] T015 [P] [US1] Create frontend cluster/resource API services and typings in frontend/src/services/clusters.ts, frontend/src/services/resources.ts, and frontend/src/services/api/types.ts
- [x] T016 [US1] Replace mock onboarding and local cluster state with backend-driven mutations and queries in frontend/src/features/clusters/components/ClusterOnboardDrawer.tsx and frontend/src/features/clusters/pages/ClusterOverviewPage.tsx
- [x] T017 [US1] Replace hardcoded resource list/detail data with real filtering and refresh flows in frontend/src/features/resources/pages/ResourceListPage.tsx, frontend/src/features/resources/components/ResourceFilters.tsx, and frontend/src/features/resources/components/ResourceDetailDrawer.tsx

**Checkpoint**: User Story 1 完成后，平台应具备真实的多集群接入、同步与资源浏览能力

---

## Phase 3: User Story 2 - 团队授权与资源隔离 (Priority: P1)

**Goal**: 把当前“页面存在但未真正控权”的状态升级为后端强制授权 + 前端按权限呈现

**Independent Test**: 创建至少两个工作空间并配置角色绑定后，不同用户只能看到并操作自己授权范围内的集群、资源和操作入口。

### Tests for User Story 2

- [x] T018 [P] [US2] Tighten access-control contract and integration tests around real scope isolation in backend/tests/contract/access_control_contract_test.go and backend/tests/integration/scope_authorization_test.go

### Implementation for User Story 2

- [x] T019 [P] [US2] Expand scope models and persistence for workspace-cluster/project bindings and role metadata in backend/internal/domain/scope.go, backend/migrations/0003_scope_bindings.sql, backend/internal/repository/workspace_cluster_repository.go, and backend/internal/repository/scope_role_repository.go (首批角色矩阵固定为 `platform-admin`、`ops-operator`、`audit-reader`、`readonly`)
- [x] T020 [US2] Enforce permission and scope checks in request middleware and route wiring in backend/internal/api/middleware/authorization.go, backend/internal/api/router/access_routes.go, and backend/internal/service/auth/permission_service.go
- [x] T021 [US2] Apply scope filtering to cluster, resource, operation, and audit queries in backend/internal/service/auth/scope_authorizer.go, backend/internal/service/cluster/service.go, backend/internal/service/operation/service.go, and backend/internal/service/audit/service.go
- [x] T022 [US2] Remove the auto-create role fallback and return explicit role-binding validation errors in backend/internal/api/handler/role_binding_handler.go and backend/internal/repository/scope_role_binding_repository.go
- [x] T023 [P] [US2] Create frontend workspace, project, and role-binding API services in frontend/src/services/workspaces.ts, frontend/src/services/projects.ts, and frontend/src/services/roleBindings.ts
- [x] T024 [US2] Replace local mock tables with backend-driven CRUD flows in frontend/src/features/workspaces/pages/WorkspacePage.tsx and frontend/src/features/projects/pages/ProjectPage.tsx
- [x] T025 [US2] Drive navigation visibility and action gating from real permissions in frontend/src/app/AuthorizedMenu.tsx, frontend/src/features/auth/components/RoleBindingForm.tsx, and frontend/src/app/App.tsx

**Checkpoint**: User Story 2 完成后，平台应具备最小安全闭环，MVP 交付建议至少覆盖 US1 + US2

---

## Phase 4: User Story 3 - 受控运维操作执行 (Priority: P2)

**Goal**: 把当前“模拟状态流转”的操作中心替换为真实执行、进度回传和失败反馈

**Independent Test**: 在授权范围内对工作负载执行扩缩容/重启，对节点执行维护动作时，平台能够展示确认、进度、结果和失败原因。

### Tests for User Story 3

- [x] T026 [P] [US3] Tighten operation contracts and execution-path assertions in backend/tests/contract/operations_contract_test.go and backend/tests/integration/operation_execution_test.go

### Implementation for User Story 3

- [x] T027 [P] [US3] Implement real Kubernetes executors for scale, restart, and node-maintenance in backend/internal/service/operation/executor.go and backend/internal/worker/operation_worker.go
- [x] T028 [US3] Persist operation progress messages, failure reason, and idempotent status transitions in backend/internal/domain/operation.go, backend/internal/repository/operation_repository.go, and backend/internal/service/operation/service.go
- [x] T029 [US3] Emit audit events for operation submit/start/success/failure in backend/internal/service/audit/event_writer.go and backend/internal/worker/operation_worker.go
- [x] T030 [P] [US3] Replace in-memory operation center data with backend polling in frontend/src/services/operations.ts and frontend/src/features/operations/pages/OperationCenterPage.tsx
- [x] T031 [US3] Connect the resource action panel and confirmation drawer to real backend payloads and error handling in frontend/src/features/resources/components/ResourceActionPanel.tsx and frontend/src/features/operations/components/OperationConfirmDrawer.tsx
- [x] T031A [US3] Ensure high-risk operation flow executes immediately after second confirmation and does not require approver assignment in backend/internal/service/operation/service.go, backend/internal/domain/operation.go, and frontend/src/features/operations/components/OperationConfirmDrawer.tsx

**Checkpoint**: User Story 3 完成后，平台应具备真实可追踪的受控运维操作链路

---

## Phase 5: User Story 4 - 审计追踪与操作复盘 (Priority: P3)

**Goal**: 去掉前端 fallback 审计数据，补齐真实查询、导出和保留策略

**Independent Test**: 审计人员能够按时间、操作者、集群、资源和结果检索真实记录，并生成可下载的导出任务。

### Tests for User Story 4

- [x] T032 [P] [US4] Tighten audit query and export tests around real filter dimensions and task lifecycle in backend/tests/contract/audit_contract_test.go and backend/tests/integration/audit_query_test.go

### Implementation for User Story 4

- [x] T033 [P] [US4] Expand audit query filters to cover cluster, workspace, project, resource, and result dimensions in backend/internal/repository/audit_repository.go, backend/internal/service/audit/service.go, and backend/internal/api/handler/audit_handler.go
- [x] T034 [US4] Implement export artifact generation, task status, and download flow in backend/internal/service/audit/service.go, backend/internal/worker/audit_export_worker.go, backend/internal/api/handler/audit_handler.go, and backend/internal/repository/audit_export_repository.go (首期仅支持 CSV 导出并执行敏感字段脱敏)
- [x] T035 [US4] Add the 180-day default audit retention cleanup job in backend/internal/worker/audit_retention_worker.go and backend/internal/repository/audit_repository.go
- [x] T036 [P] [US4] Remove frontend audit mock fallback and align payloads with backend contracts in frontend/src/services/audit.ts and frontend/src/services/api/types.ts
- [x] T037 [US4] Implement export polling, empty states, and failure handling in frontend/src/features/audit/pages/AuditEventPage.tsx, frontend/src/features/audit/components/AuditExportModal.tsx, and frontend/src/features/audit/components/AuditEventTable.tsx

**Checkpoint**: User Story 4 完成后，平台应形成完整的访问与操作审计闭环

---

## Final Phase: Polish & Delivery Readiness

**Purpose**: 收尾性能、文档、验证与交付证据

- [x] T038 [P] Add route-level code splitting to reduce the oversized frontend bundle in frontend/src/app/router.tsx and frontend/src/main.tsx
- [x] T039 [P] Refresh startup/configuration documentation and examples in README.md, backend/config/config.example.yaml, backend/config/config.dev.yaml, frontend/.env.example, and frontend/.env.development
- [x] T040 [P] Record the real verification baseline for `go test ./...`, `npm run lint`, `npm run build`, and frontend Vitest in artifacts/001-k8s-ops-platform/verification.md, and attach performance acceptance evidence (`test-environment` load report + reproducible experiment scripts)
- [x] T041 Prepare the updated Chinese PR summary, remaining risks, and delivery notes in artifacts/001-k8s-ops-platform/pr-summary.md and artifacts/001-k8s-ops-platform/pr-readiness.md

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 0: Governance & Re-Baseline**: 无前置依赖，必须先完成分支与基线校正
- **Phase 1: Shared Quality & Foundation**: 依赖 Phase 0，是所有后续故事的共同前提
- **Phase 2: US1**: 依赖 Phase 1，是“真实可用平台入口”的核心
- **Phase 3: US2**: 依赖 Phase 1，建议与 US1 紧密衔接以形成权限闭环
- **Phase 4: US3**: 依赖 US1 的真实资源视图与 US2 的真实授权边界
- **Phase 5: US4**: 依赖 Phase 1，且最佳效果建立在 US1-US3 真实事件流之上
- **Final Phase**: 依赖目标用户故事完成

### User Story Dependencies

- **US1**: 当前最关键缺口是把集群接入、同步和资源浏览从 mock/stub 变成真实链路
- **US2**: 当前最关键缺口是把权限模型从“定义存在”变成“请求强制校验 + 页面受控呈现”
- **US3**: 当前最关键缺口是把操作执行从模拟状态机变成真实 Kubernetes 动作
- **US4**: 当前最关键缺口是把审计从 fallback 数据变成真实查询与导出

### Parallel Opportunities

- Phase 1 中，`T005`、`T006`、`T007` 可并行推进
- US1 中，`T009`、`T010`、`T012`、`T015` 可并行推进
- US2 中，`T018`、`T019`、`T023` 可并行推进
- US3 中，`T026`、`T027`、`T030` 可并行推进
- US4 中，`T032`、`T033`、`T036` 可并行推进
- Final Phase 中，`T038`、`T039`、`T040` 可并行推进

---

## Recommended Delivery Strategy

### MVP Recovery First

1. 完成 Phase 0 与 Phase 1，修正治理偏差并建立真实测试基线
2. 完成 US1，把“能看 demo”推进到“能接真实集群并浏览真实资源”
3. 立刻完成 US2，补齐权限隔离，形成最小可交付闭环

### Incremental Delivery

1. **Increment 1**: US1 + US2，形成真实接入、浏览、授权闭环
2. **Increment 2**: US3，补齐受控运维执行
3. **Increment 3**: US4，补齐审计追踪与导出
4. 每个增量完成后更新中文 PR 说明、验证结果和风险清单

### Notes

- 本文件已替换掉原先“所有任务全部完成”的失真状态，只保留当前剩余工作
- 如果后续要继续对标 Rancher 的日志、终端、Helm、监控等能力，应为 `002+` 新 feature 单独走 `/speckit.specify -> /speckit.plan -> /speckit.tasks`
- 当前 001 的最优先交付范围是 **US1 + US2**，而不是继续在 mock 页面上横向扩功能
