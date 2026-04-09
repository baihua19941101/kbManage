# Tasks: 多集群 Kubernetes 可视化管理平台

**Input**: Design documents from `/specs/001-k8s-ops-platform/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/openapi.yaml, quickstart.md

**Tests**: `plan.md` 已明确要求契约测试与集成测试，因此每个用户故事均包含对应测试任务。

**Organization**: 任务按用户故事分组，确保每个故事都可以独立实现、独立验证。

**Constitutional Gates**: 功能分支校验、数据库备份证据、国内依赖源配置、中文 PR 产物、GitHub 推送和合并前用户授权均为强制项。

## Format: `[ID] [P?] [Story] Description`

- **[P]**: 可并行执行（不同文件、无未完成依赖）
- **[Story]**: 用户故事标签，仅出现在用户故事阶段
- 每个任务都包含明确文件路径，保证可直接执行

## Phase 0: Governance Gates

**Purpose**: 在编码前完成宪章要求的治理动作并留下证据

- [ ] T001 Record current feature branch and remote status in artifacts/001-k8s-ops-platform/branch-check.txt
- [ ] T002 Create the pre-development MySQL backup artifact and manifest in artifacts/001-k8s-ops-platform/mysql-backup-<timestamp>.sql and artifacts/001-k8s-ops-platform/backup-manifest.txt
- [ ] T003 Configure domestic Go and npm mirrors in backend/.env.example and frontend/.npmrc
- [ ] T004 Record PR workflow, merge approval rule, and commit-message standard in artifacts/001-k8s-ops-platform/pr-readiness.md

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: 初始化前后端项目骨架和本地开发命令

- [ ] T005 Create the backend service entrypoint and module scaffold in backend/cmd/server/main.go and backend/go.mod
- [ ] T006 Create the frontend application scaffold in frontend/package.json and frontend/src/main.tsx
- [ ] T007 Configure frontend build, lint, and test tooling in frontend/vite.config.ts, frontend/vitest.config.ts, frontend/eslint.config.js, and frontend/tsconfig.json
- [ ] T008 [P] Add shared environment templates in backend/.env.example and frontend/.env.example
- [ ] T009 [P] Add repository-level development commands in Makefile

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: 完成所有用户故事都会依赖的基础设施和共享能力

**⚠️ CRITICAL**: 本阶段完成前不得开始任何用户故事开发

- [ ] T010 Create the initial MySQL schema migration for auth, RBAC, workspace, cluster, and audit core tables in backend/migrations/0001_platform_core.sql
- [ ] T011 [P] Implement backend config, GORM, and Redis bootstrap in backend/internal/repository/bootstrap.go and backend/internal/repository/redis.go
- [ ] T012 [P] Define shared domain models in backend/internal/domain/user.go, backend/internal/domain/scope.go, backend/internal/domain/cluster.go, backend/internal/domain/audit.go, and backend/internal/domain/operation.go
- [ ] T013 [P] Implement shared repositories for users, roles, sessions, workspaces, projects, and audits in backend/internal/repository/user_repository.go, backend/internal/repository/platform_role_repository.go, backend/internal/repository/session_repository.go, backend/internal/repository/workspace_repository.go, backend/internal/repository/project_repository.go, and backend/internal/repository/audit_repository.go
- [ ] T014 [P] Implement password hashing, token issuance, and refresh-session service in backend/internal/service/auth/password_service.go and backend/internal/service/auth/token_service.go
- [ ] T015 [P] Implement platform RBAC and scope authorization engine in backend/internal/service/auth/permission_service.go and backend/internal/service/auth/scope_authorizer.go
- [ ] T016 [P] Implement Gin router, request context, and unified error middleware in backend/internal/api/router/router.go, backend/internal/api/middleware/auth.go, backend/internal/api/middleware/request_id.go, and backend/internal/api/middleware/error_handler.go
- [ ] T017 [P] Implement shared audit event writer in backend/internal/service/audit/event_writer.go
- [ ] T018 [P] Implement Kubernetes client manager and informer cache bootstrap in backend/internal/kube/client/manager.go and backend/internal/kube/cache/sync_manager.go
- [ ] T019 [P] Implement frontend app shell, route guard, auth store, and query client in frontend/src/app/App.tsx, frontend/src/app/router.tsx, frontend/src/features/auth/store.ts, and frontend/src/app/queryClient.ts
- [ ] T020 Implement login page and refresh-session flow in frontend/src/features/auth/pages/LoginPage.tsx and frontend/src/services/auth.ts

**Checkpoint**: Foundation ready - 用户故事开发可以按优先级推进

---

## Phase 3: User Story 1 - 多集群统一接入与资源总览 (Priority: P1) 🎯 MVP

**Goal**: 提供多集群接入、健康总览、资源列表筛选和资源详情查看

**Independent Test**: 接入至少 2 个集群后，管理员能够在统一视图中按集群、命名空间和资源类型浏览资源，并识别异常集群与异常资源。

### Tests for User Story 1

- [ ] T021 [P] [US1] Add cluster and resource contract tests in backend/tests/contract/clusters_contract_test.go
- [ ] T022 [P] [US1] Add multi-cluster overview integration test in backend/tests/integration/cluster_overview_test.go

### Implementation for User Story 1

- [ ] T023 [P] [US1] Implement cluster credential repository and encryption helper in backend/internal/repository/cluster_credential_repository.go and backend/internal/service/cluster/credential_cipher.go
- [ ] T024 [P] [US1] Implement cluster and resource inventory repositories in backend/internal/repository/cluster_repository.go and backend/internal/repository/resource_inventory_repository.go
- [ ] T025 [US1] Implement cluster onboarding, connectivity verification, and health sync service in backend/internal/service/cluster/service.go
- [ ] T026 [US1] Implement resource inventory indexing adapter in backend/internal/kube/adapter/resource_indexer.go
- [ ] T027 [US1] Implement cluster and resource handlers/routes in backend/internal/api/handler/cluster_handler.go, backend/internal/api/handler/resource_handler.go, and backend/internal/api/router/cluster_routes.go
- [ ] T028 [P] [US1] Build cluster overview and onboarding UI in frontend/src/features/clusters/pages/ClusterOverviewPage.tsx and frontend/src/features/clusters/components/ClusterOnboardDrawer.tsx
- [ ] T029 [US1] Implement resource list, filters, and detail drawer in frontend/src/features/resources/pages/ResourceListPage.tsx, frontend/src/features/resources/components/ResourceFilters.tsx, and frontend/src/features/resources/components/ResourceDetailDrawer.tsx

**Checkpoint**: User Story 1 完成后，应可独立展示多集群接入与统一资源总览

---

## Phase 4: User Story 2 - 团队授权与资源隔离 (Priority: P1)

**Goal**: 提供平台级 RBAC、工作空间/项目管理和作用域隔离

**Independent Test**: 创建至少两个独立的工作空间并分别授权给不同角色后，每个用户只能看到并操作自己授权范围内的资源。

### Tests for User Story 2

- [ ] T030 [P] [US2] Add workspace, project, and role-binding contract tests in backend/tests/contract/access_control_contract_test.go
- [ ] T031 [P] [US2] Add workspace isolation integration test in backend/tests/integration/scope_authorization_test.go

### Implementation for User Story 2

- [ ] T032 [P] [US2] Implement workspace, project, and cluster-binding repositories in backend/internal/repository/workspace_repository.go, backend/internal/repository/project_repository.go, and backend/internal/repository/workspace_cluster_repository.go
- [ ] T033 [P] [US2] Implement scope role and binding repositories in backend/internal/repository/scope_role_repository.go and backend/internal/repository/scope_role_binding_repository.go
- [ ] T034 [US2] Implement workspace and project management services in backend/internal/service/workspace/service.go and backend/internal/service/project/service.go
- [ ] T035 [US2] Implement workspace, project, and role-binding handlers/routes in backend/internal/api/handler/workspace_handler.go, backend/internal/api/handler/project_handler.go, backend/internal/api/handler/role_binding_handler.go, and backend/internal/api/router/access_routes.go
- [ ] T036 [P] [US2] Build workspace and project management pages in frontend/src/features/workspaces/pages/WorkspacePage.tsx and frontend/src/features/projects/pages/ProjectPage.tsx
- [ ] T037 [US2] Implement role-binding UI and authorization-aware navigation in frontend/src/features/auth/components/RoleBindingForm.tsx and frontend/src/app/AuthorizedMenu.tsx

**Checkpoint**: User Story 2 完成后，应可独立验证平台级 RBAC 与作用域授权隔离

---

## Phase 5: User Story 3 - 受控运维操作执行 (Priority: P2)

**Goal**: 提供高风险确认、操作执行、进度跟踪和失败反馈

**Independent Test**: 在授权范围内选择一个工作负载和一个节点，执行常见运维动作后，平台能展示确认、执行过程与最终结果。

### Tests for User Story 3

- [ ] T038 [P] [US3] Add operation contract tests in backend/tests/contract/operations_contract_test.go
- [ ] T039 [P] [US3] Add controlled operation integration test in backend/tests/integration/operation_execution_test.go

### Implementation for User Story 3

- [ ] T040 [P] [US3] Implement operation request repository and idempotency lock support in backend/internal/repository/operation_repository.go and backend/internal/service/operation/idempotency_service.go
- [ ] T041 [P] [US3] Implement operation queue and worker orchestration in backend/internal/worker/operation_worker.go and backend/internal/service/operation/queue_service.go
- [ ] T042 [US3] Implement operation confirmation, risk evaluation, and execution service in backend/internal/service/operation/service.go
- [ ] T043 [US3] Implement operation handlers/routes in backend/internal/api/handler/operation_handler.go and backend/internal/api/router/operation_routes.go
- [ ] T044 [P] [US3] Build operation center and confirmation drawer in frontend/src/features/operations/pages/OperationCenterPage.tsx and frontend/src/features/operations/components/OperationConfirmDrawer.tsx
- [ ] T045 [US3] Integrate resource actions with operation workflow in frontend/src/features/resources/components/ResourceActionPanel.tsx and frontend/src/services/operations.ts

**Checkpoint**: User Story 3 完成后，应可独立验证受控运维操作流程

---

## Phase 6: User Story 4 - 审计追踪与操作复盘 (Priority: P3)

**Goal**: 提供审计检索、导出和跨动作复盘视图

**Independent Test**: 在平台产生多条访问和运维记录后，审计人员能筛选出指定时间段、指定操作者和指定资源的完整记录。

### Tests for User Story 4

- [ ] T046 [P] [US4] Add audit query and export contract tests in backend/tests/contract/audit_contract_test.go
- [ ] T047 [P] [US4] Add audit search and export integration test in backend/tests/integration/audit_query_test.go

### Implementation for User Story 4

- [ ] T048 [P] [US4] Implement audit query repository and export job repository in backend/internal/repository/audit_repository.go and backend/internal/repository/audit_export_repository.go
- [ ] T049 [US4] Implement audit query, retention, and export service in backend/internal/service/audit/service.go and backend/internal/worker/audit_export_worker.go
- [ ] T050 [US4] Implement audit handlers/routes in backend/internal/api/handler/audit_handler.go and backend/internal/api/router/audit_routes.go
- [ ] T051 [P] [US4] Build audit query and export UI in frontend/src/features/audit/pages/AuditEventPage.tsx and frontend/src/features/audit/components/AuditExportModal.tsx
- [ ] T052 [US4] Implement audit event table and filter client in frontend/src/features/audit/components/AuditEventTable.tsx and frontend/src/services/audit.ts

**Checkpoint**: User Story 4 完成后，应可独立验证审计查询与导出能力

---

## Final Phase: Polish & Cross-Cutting Concerns

**Purpose**: 完成跨故事质量加固、文档验证与交付准备

- [ ] T053 [P] Validate development steps against specs/001-k8s-ops-platform/quickstart.md
- [ ] T054 Harden password, credential, and permission edge handling in backend/internal/service/auth/password_service.go and backend/internal/service/cluster/credential_cipher.go
- [ ] T055 [P] Sync finalized API contract and frontend service typings in specs/001-k8s-ops-platform/contracts/openapi.yaml and frontend/src/services/api/types.ts
- [ ] T056 [P] Prepare the Chinese PR summary with backup evidence, test results, and risk notes in artifacts/001-k8s-ops-platform/pr-summary.md
- [ ] T057 Update GitHub push and PR evidence in artifacts/001-k8s-ops-platform/pr-readiness.md
- [ ] T058 Record explicit user approval before merge in artifacts/001-k8s-ops-platform/merge-approval.md

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 0: Governance Gates**: 无前置依赖，必须最先完成
- **Phase 1: Setup**: 依赖 Governance Gates
- **Phase 2: Foundational**: 依赖 Setup，阻塞所有用户故事
- **Phase 3: US1**: 依赖 Foundational，为 MVP 基础
- **Phase 4: US2**: 依赖 Foundational，可与 US1 并行推进，但推荐紧随 US1 完成
- **Phase 5: US3**: 依赖 US1 的资源总览能力和 US2 的权限边界能力
- **Phase 6: US4**: 依赖 Foundational 和共享审计事件能力；为获得完整审计样本，推荐在 US1-US3 后执行
- **Final Phase**: 依赖所有目标用户故事完成

### User Story Dependencies

- **US1**: 可在 Foundational 完成后单独实现，构成 MVP 主体
- **US2**: 可在 Foundational 完成后独立实现，但最终验证需要与 US1 的资源可见性联动
- **US3**: 依赖 US1 提供资源入口和 US2 提供授权边界
- **US4**: 依赖共享审计写入能力，最佳顺序是在 US1-US3 后完成，以便覆盖真实事件流

### Parallel Opportunities

- Governance 阶段中，`T003` 和 `T004` 可在备份任务之外并行处理
- Setup 阶段中，`T006`、`T008`、`T009` 可并行推进
- Foundational 阶段中，`T011` 至 `T019` 大部分为不同文件的并行任务
- US1 可并行：`T021`、`T022`、`T023`、`T024`、`T028`
- US2 可并行：`T030`、`T031`、`T032`、`T033`、`T036`
- US3 可并行：`T038`、`T039`、`T040`、`T041`、`T044`
- US4 可并行：`T046`、`T047`、`T048`、`T051`
- Final Phase 可并行：`T053`、`T055`、`T056`

---

## Parallel Example: User Story 1

```bash
Task: "T021 [US1] Add cluster and resource contract tests in backend/tests/contract/clusters_contract_test.go"
Task: "T022 [US1] Add multi-cluster overview integration test in backend/tests/integration/cluster_overview_test.go"
Task: "T023 [US1] Implement cluster credential repository and encryption helper in backend/internal/repository/cluster_credential_repository.go and backend/internal/service/cluster/credential_cipher.go"
Task: "T024 [US1] Implement cluster and resource inventory repositories in backend/internal/repository/cluster_repository.go and backend/internal/repository/resource_inventory_repository.go"
Task: "T028 [US1] Build cluster overview and onboarding UI in frontend/src/features/clusters/pages/ClusterOverviewPage.tsx and frontend/src/features/clusters/components/ClusterOnboardDrawer.tsx"
```

## Parallel Example: User Story 2

```bash
Task: "T030 [US2] Add workspace, project, and role-binding contract tests in backend/tests/contract/access_control_contract_test.go"
Task: "T031 [US2] Add workspace isolation integration test in backend/tests/integration/scope_authorization_test.go"
Task: "T032 [US2] Implement workspace, project, and cluster-binding repositories in backend/internal/repository/workspace_repository.go, backend/internal/repository/project_repository.go, and backend/internal/repository/workspace_cluster_repository.go"
Task: "T033 [US2] Implement scope role and binding repositories in backend/internal/repository/scope_role_repository.go and backend/internal/repository/scope_role_binding_repository.go"
Task: "T036 [US2] Build workspace and project management pages in frontend/src/features/workspaces/pages/WorkspacePage.tsx and frontend/src/features/projects/pages/ProjectPage.tsx"
```

## Parallel Example: User Story 3

```bash
Task: "T038 [US3] Add operation contract tests in backend/tests/contract/operations_contract_test.go"
Task: "T039 [US3] Add controlled operation integration test in backend/tests/integration/operation_execution_test.go"
Task: "T040 [US3] Implement operation request repository and idempotency lock support in backend/internal/repository/operation_repository.go and backend/internal/service/operation/idempotency_service.go"
Task: "T041 [US3] Implement operation queue and worker orchestration in backend/internal/worker/operation_worker.go and backend/internal/service/operation/queue_service.go"
Task: "T044 [US3] Build operation center and confirmation drawer in frontend/src/features/operations/pages/OperationCenterPage.tsx and frontend/src/features/operations/components/OperationConfirmDrawer.tsx"
```

## Parallel Example: User Story 4

```bash
Task: "T046 [US4] Add audit query and export contract tests in backend/tests/contract/audit_contract_test.go"
Task: "T047 [US4] Add audit search and export integration test in backend/tests/integration/audit_query_test.go"
Task: "T048 [US4] Implement audit query repository and export job repository in backend/internal/repository/audit_repository.go and backend/internal/repository/audit_export_repository.go"
Task: "T051 [US4] Build audit query and export UI in frontend/src/features/audit/pages/AuditEventPage.tsx and frontend/src/features/audit/components/AuditExportModal.tsx"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. 完成 Governance Gates、Setup、Foundational
2. 完成 US1 的契约测试、集群接入、资源索引和前端总览
3. 独立验证多集群接入与统一资源浏览
4. 将分支推送到 GitHub 并更新中文 PR 说明

### Incremental Delivery

1. **MVP**: US1 提供多集群接入与资源总览
2. **Increment 2**: US2 增加平台级 RBAC 与工作空间/项目隔离
3. **Increment 3**: US3 增加受控运维操作和执行跟踪
4. **Increment 4**: US4 增加审计检索与导出
5. 每个增量完成后更新中文 PR、补充测试结果，并等待用户对合并动作的明确同意

### Parallel Team Strategy

1. 开发者 A 负责后端基础能力和 US1 / US3 服务层
2. 开发者 B 负责前端壳层、US1 / US2 页面与交互
3. 开发者 C 负责权限、审计与契约/集成测试
4. 如使用子代理，全部子代理必须固定为 `gpt-5.3-codex`

---

## Notes

- 所有任务均遵循 `- [ ] Txxx [P?] [US?] 描述 + 文件路径` 格式
- 用户故事阶段的每个任务都带有 `US1` 至 `US4` 标签
- 由于 `plan.md` 已明确测试策略，本任务清单包含契约测试和集成测试
- 最小 MVP 建议范围是 **US1**；若需要最小权限安全闭环，则扩展到 **US1 + US2**
- 合并前必须完成备份、中文 PR、用户批准记录三项治理证据
