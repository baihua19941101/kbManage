# Tasks: 身份与多租户治理中心

**Input**: Design documents from `/specs/009-identity-tenancy/`  
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/openapi.yaml, quickstart.md

**Tests**: 本特性包含明确的独立验收标准、权限边界、委派回收和会话治理要求，任务清单包含后端契约/集成测试与前端页面测试任务。  
**Organization**: 任务按用户故事分组，保证每个故事可独立实现、独立验证。

**Constitutional Gates**: 必须满足功能分支、数据库备份证据、国内依赖源配置、中文 PR、远程推送、用户同意后合并。

## Format: `[ID] [P?] [Story] Description`

- `[P]`: 可并行执行（不同文件、无前置依赖）
- `[Story]`: 任务归属用户故事（US1/US2/US3）
- 每条任务必须包含明确文件路径

## Phase 0: Governance Gates

**Purpose**: 完成宪章门槛与实施前证据准备

- [X] T001 在 `artifacts/009-identity-tenancy/branch-check.txt` 记录当前分支、禁止在 `main/master` 开发、`008` 已合并完成以及“未获用户同意不得合并”门槛
- [X] T002 执行 009 开发前数据库备份并在 `artifacts/009-identity-tenancy/backup-manifest.txt` 记录命令、时间戳、产物路径与恢复抽样验证结果
- [X] T003 在 `artifacts/009-identity-tenancy/mirror-and-remote-check.txt` 记录 `GOPROXY=https://goproxy.cn,direct`、`https://registry.npmmirror.com`、身份源联调依赖镜像来源与 `git@github.com:baihua19941101/kbManage.git` PR 流程

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: 建立 009 模块骨架、配置入口和导航占位

- [X] T004 创建后端身份治理模块骨架 `backend/internal/service/identitytenancy/`、`backend/internal/api/handler/identity_tenancy_handler.go`、`backend/internal/api/router/identity_tenancy_routes.go`
- [X] T005 [P] 创建身份源适配层目录与占位 `backend/internal/integration/identity/`、`backend/internal/integration/identity/provider.go`、`backend/internal/integration/identity/sync_provider.go`
- [X] T006 [P] 创建前端模块骨架 `frontend/src/features/identity-tenancy/`、`frontend/src/services/identityTenancy.ts` 与路由占位到 `frontend/src/app/router.tsx`
- [X] T007 在 `backend/config/config.example.yaml`、`backend/config/config.dev.yaml`、`frontend/.env.example`、`frontend/.env.development`、`README.md` 增加 `identityTenancy.*` 配置说明
- [X] T008 在 `frontend/src/app/AuthorizedMenu.tsx`、`frontend/src/app/ProtectedRoute.tsx`、`frontend/src/features/auth/store.ts` 接入身份与租户治理中心入口门禁占位

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: 完成所有用户故事共享且阻塞性的基础能力

**⚠️ CRITICAL**: US1/US2/US3 必须在本阶段完成后才可开始

- [X] T009 新增 009 数据库迁移 `backend/migrations/0011_identity_tenancy_core.sql`，落库身份源、身份映射、组织单元、成员关系、租户映射、角色定义、授权分配、委派关系、会话治理、风险快照与身份治理审计表
- [X] T010 [P] 在 `backend/internal/domain/identity_tenancy.go` 定义 `IdentitySource`、`IdentityAccount`、`OrganizationUnit`、`OrganizationMembership`、`TenantScopeMapping`、`RoleDefinition`、`RoleAssignment`、`DelegationGrant`、`SessionRecord`、`AccessRiskSnapshot`
- [X] T011 [P] 创建仓储实现 `backend/internal/repository/identity_source_repository.go`、`backend/internal/repository/organization_unit_repository.go`、`backend/internal/repository/role_definition_repository.go`、`backend/internal/repository/role_assignment_repository.go`、`backend/internal/repository/delegation_grant_repository.go`、`backend/internal/repository/session_record_repository.go`、`backend/internal/repository/access_risk_repository.go`
- [X] T012 [P] 在 `backend/internal/integration/identity/provider.go`、`backend/internal/integration/identity/sync_provider.go` 定义身份源接入抽象、目录同步结果模型和错误归一化语义
- [X] T013 在 `backend/internal/service/identitytenancy/service.go`、`backend/internal/service/identitytenancy/scope_service.go`、`backend/internal/service/auth/scope_authorizer.go` 建立统一身份治理授权与租户边界过滤入口
- [X] T014 [P] 在 `backend/internal/repository/redis.go`、`backend/internal/service/identitytenancy/session_cache.go`、`backend/internal/service/identitytenancy/permission_cache.go`、`backend/internal/service/identitytenancy/revocation_coordinator.go` 建立会话缓存、权限评估缓存和回收协调
- [X] T015 [P] 在 `backend/internal/service/auth/permission_service.go`、`backend/internal/api/middleware/authorization.go` 增加 009 权限语义 `identity:read`、`identity:manage-source`、`identity:manage-org`、`identity:manage-role`、`identity:delegate`、`identity:session-govern`
- [X] T016 在 `backend/internal/api/router/identity_tenancy_routes.go`、`backend/internal/api/router/router.go` 注册 009 API 路由骨架
- [X] T017 [P] 在 `frontend/src/services/api/types.ts`、`frontend/src/services/api/client.ts`、`frontend/src/app/queryClient.ts` 增加 009 共享类型、查询 key 和错误归一化
- [X] T018 在 `backend/internal/service/audit/event_writer.go`、`backend/internal/service/audit/service.go` 预置 `identitytenancy.*` 审计动作类型与查询维度映射

**Checkpoint**: 基础能力完成，可开始用户故事实现

---

## Phase 3: User Story 1 - 接入外部身份源并统一登录 (Priority: P1) 🎯 MVP

**Goal**: 提供身份源接入、本地账号并存、统一登录方式和来源状态可视化能力。  
**Independent Test**: 接入一个外部身份源并保留本地管理员账号后，用户可切换登录方式，管理员可看到身份源状态、用户来源和会话状态。

### Tests for User Story 1

- [X] T019 [P] [US1] 编写后端契约测试 `backend/tests/contract/identity_source_contract_test.go`、`backend/tests/contract/identity_login_mode_contract_test.go`、`backend/tests/contract/session_governance_contract_test.go`
- [X] T020 [P] [US1] 编写后端集成测试 `backend/tests/integration/identity_source_lifecycle_test.go`、`backend/tests/integration/login_mode_switch_test.go`、`backend/tests/integration/local_fallback_access_test.go`
- [X] T021 [P] [US1] 编写前端 Vitest 页面测试 `frontend/src/features/identity-tenancy/pages/IdentitySourcePage.test.tsx`、`frontend/src/features/identity-tenancy/pages/SessionGovernancePage.test.tsx`

### Implementation for User Story 1

- [X] T022 [P] [US1] 实现身份源与身份映射服务 `backend/internal/service/identitytenancy/identity_source_service.go`、`backend/internal/service/identitytenancy/identity_account_service.go`
- [X] T023 [P] [US1] 实现登录方式切换与会话治理服务 `backend/internal/service/identitytenancy/login_mode_service.go`、`backend/internal/service/identitytenancy/session_service.go`
- [X] T024 [US1] 在 `backend/internal/api/handler/identity_tenancy_handler.go`、`backend/internal/api/router/identity_tenancy_routes.go` 落地 `/identity/sources`、`/identity/sources/{sourceId}`、`/identity/sessions`
- [X] T025 [US1] 在 `backend/internal/service/identitytenancy/scope_service.go`、`backend/internal/api/middleware/authorization.go`、`backend/internal/service/auth/scope_authorizer.go` 落地身份源、登录方式和会话治理路径的范围过滤
- [X] T026 [P] [US1] 实现前端服务层 `frontend/src/services/identityTenancy.ts`，覆盖身份源、登录方式和会话治理接口
- [X] T027 [P] [US1] 实现身份源页面与表单 `frontend/src/features/identity-tenancy/pages/IdentitySourcePage.tsx`、`frontend/src/features/identity-tenancy/components/IdentitySourceDrawer.tsx`
- [X] T028 [P] [US1] 实现会话治理页面 `frontend/src/features/identity-tenancy/pages/SessionGovernancePage.tsx`、`frontend/src/features/identity-tenancy/components/SessionRiskDrawer.tsx`
- [X] T029 [US1] 在 `frontend/src/app/router.tsx`、`frontend/src/features/auth/pages/LoginPage.tsx` 打通统一登录方式切换与身份治理中心导航入口

**Checkpoint**: US1 完整可测，可作为 009 MVP 交付

---

## Phase 4: User Story 2 - 建立组织与租户关系模型 (Priority: P1)

**Goal**: 提供组织、团队、用户组、工作空间和项目之间的关系建模与租户边界映射能力。  
**Independent Test**: 创建一个组织、多个团队和用户组并映射到工作空间和项目后，可看到归属关系、租户边界和成员视图。

### Tests for User Story 2

- [X] T030 [P] [US2] 编写后端契约测试 `backend/tests/contract/organization_unit_contract_test.go`、`backend/tests/contract/organization_mapping_contract_test.go`、`backend/tests/contract/membership_query_contract_test.go`
- [X] T031 [P] [US2] 编写后端集成测试 `backend/tests/integration/organization_tree_flow_test.go`、`backend/tests/integration/tenant_scope_mapping_flow_test.go`、`backend/tests/integration/membership_boundary_query_test.go`
- [X] T032 [P] [US2] 编写前端 Vitest 页面测试 `frontend/src/features/identity-tenancy/pages/OrganizationModelPage.test.tsx`、`frontend/src/features/identity-tenancy/pages/TenantMappingPage.test.tsx`

### Implementation for User Story 2

- [X] T033 [P] [US2] 实现组织与成员关系服务 `backend/internal/service/identitytenancy/organization_service.go`、`backend/internal/service/identitytenancy/membership_service.go`
- [X] T034 [P] [US2] 实现租户边界映射服务 `backend/internal/service/identitytenancy/tenant_mapping_service.go`、`backend/internal/service/identitytenancy/boundary_view_service.go`
- [X] T035 [US2] 在 `backend/internal/api/handler/identity_tenancy_handler.go`、`backend/internal/api/router/identity_tenancy_routes.go` 落地 `/identity/organizations`、`/identity/organizations/{unitId}/mappings`
- [X] T036 [US2] 在 `backend/internal/service/audit/event_writer.go`、`backend/internal/service/audit/service.go` 打通组织变更、成员调整和租户映射动作的审计写入与查询聚合
- [X] T037 [P] [US2] 扩展前端服务与 hooks `frontend/src/services/identityTenancy.ts`、`frontend/src/features/identity-tenancy/hooks/useOrganizationAction.ts`
- [X] T038 [P] [US2] 实现组织模型页面与表单 `frontend/src/features/identity-tenancy/pages/OrganizationModelPage.tsx`、`frontend/src/features/identity-tenancy/components/OrganizationUnitDrawer.tsx`
- [X] T039 [P] [US2] 实现租户边界映射页面与详情 `frontend/src/features/identity-tenancy/pages/TenantMappingPage.tsx`、`frontend/src/features/identity-tenancy/components/TenantScopeDrawer.tsx`
- [X] T040 [US2] 在 `frontend/src/app/router.tsx`、`frontend/src/features/identity-tenancy/pages/OrganizationModelPage.tsx` 落地租户边界可视化、边界冲突空态和成员来源视图

**Checkpoint**: US2 可独立验证组织模型和租户边界闭环

---

## Phase 5: User Story 3 - 管理细粒度 RBAC、委派和回收 (Priority: P2)

**Goal**: 提供多层级角色定义、授权分配、委派、临时授权、访问回收和风险视图能力。  
**Independent Test**: 为用户或用户组分配多层级角色后，可看到权限边界、委派链路、临时授权到期状态并执行回收。

### Tests for User Story 3

- [X] T041 [P] [US3] 编写后端契约测试 `backend/tests/contract/role_definition_contract_test.go`、`backend/tests/contract/role_assignment_contract_test.go`、`backend/tests/contract/delegation_grant_contract_test.go`、`backend/tests/contract/access_risk_contract_test.go`
- [X] T042 [P] [US3] 编写后端集成测试 `backend/tests/integration/role_inheritance_flow_test.go`、`backend/tests/integration/delegation_lifecycle_test.go`、`backend/tests/integration/temporary_access_revocation_test.go`、`backend/tests/integration/access_risk_query_test.go`
- [X] T043 [P] [US3] 编写前端 Vitest 页面测试 `frontend/src/features/identity-tenancy/pages/RoleCatalogPage.test.tsx`、`frontend/src/features/identity-tenancy/pages/RoleAssignmentPage.test.tsx`、`frontend/src/features/identity-tenancy/pages/AccessRiskPage.test.tsx`

### Implementation for User Story 3

- [X] T044 [P] [US3] 实现角色定义与授权服务 `backend/internal/service/identitytenancy/role_definition_service.go`、`backend/internal/service/identitytenancy/role_assignment_service.go`
- [X] T045 [P] [US3] 实现委派、临时授权与访问回收服务 `backend/internal/service/identitytenancy/delegation_service.go`、`backend/internal/service/identitytenancy/revocation_service.go`
- [X] T046 [P] [US3] 实现权限边界与风险聚合服务 `backend/internal/service/identitytenancy/risk_service.go`、`backend/internal/service/identitytenancy/effective_permission_service.go`
- [X] T047 [US3] 在 `backend/internal/api/handler/identity_tenancy_handler.go`、`backend/internal/api/router/identity_tenancy_routes.go` 落地 `/identity/roles`、`/identity/assignments`、`/identity/delegations`、`/identity/access-risks`
- [X] T048 [US3] 在 `backend/internal/api/handler/audit_handler.go`、`backend/internal/api/router/audit_routes.go` 聚合并暴露 `/audit/identity/events` 查询链路
- [X] T049 [P] [US3] 扩展前端服务与 hooks `frontend/src/services/identityTenancy.ts`、`frontend/src/features/identity-tenancy/hooks/useRoleGovernanceAction.ts`
- [X] T050 [P] [US3] 实现角色目录与授权分配页面 `frontend/src/features/identity-tenancy/pages/RoleCatalogPage.tsx`、`frontend/src/features/identity-tenancy/pages/RoleAssignmentPage.tsx`、`frontend/src/features/identity-tenancy/components/RoleAssignmentDrawer.tsx`
- [X] T051 [P] [US3] 实现委派和临时授权页面 `frontend/src/features/identity-tenancy/pages/DelegationPage.tsx`、`frontend/src/features/identity-tenancy/components/DelegationGrantDrawer.tsx`
- [X] T052 [P] [US3] 实现风险视图与审计页面 `frontend/src/features/identity-tenancy/pages/AccessRiskPage.tsx`、`frontend/src/features/audit/pages/IdentityGovernanceAuditPage.tsx`
- [X] T053 [US3] 在 `frontend/src/app/AuthorizedMenu.tsx`、`frontend/src/app/router.tsx`、`frontend/src/features/identity-tenancy/pages/RoleAssignmentPage.tsx` 落地角色边界、委派限制、临时授权到期提示和回收空态

**Checkpoint**: US3 完成后形成身份源、组织、RBAC 与访问风险审计闭环

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: 收敛质量、文档、验证证据与 PR 交付材料

- [X] T054 [P] 收敛命名与共享类型，在 `backend/internal/service/identitytenancy/`、`backend/internal/integration/identity/`、`frontend/src/features/identity-tenancy/`、`frontend/src/services/identityTenancy.ts` 清理重复字段与错误文案
- [X] T055 [P] 刷新配置与启动文档，在 `README.md`、`backend/config/config.example.yaml`、`backend/config/config.dev.yaml`、`frontend/.env.example`、`frontend/.env.development` 补齐 009 说明
- [X] T056 [P] 记录验证基线到 `artifacts/009-identity-tenancy/verification.md`、`artifacts/009-identity-tenancy/quickstart-validation.md`、`artifacts/009-identity-tenancy/repro-identity-tenancy-smoke.sh`
- [X] T057 在 `artifacts/009-identity-tenancy/pr-summary.md`、`artifacts/009-identity-tenancy/pr-readiness.md` 准备中文 PR 摘要、备份证据、测试说明、风险清单与用户合并确认项

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 0 (Governance Gates)**: 无依赖，必须先完成
- **Phase 1 (Setup)**: 依赖 Phase 0
- **Phase 2 (Foundational)**: 依赖 Phase 1，阻塞所有用户故事
- **Phase 3/4/5 (User Stories)**: 均依赖 Phase 2；US1 作为 MVP 优先，US2 在 US1 主干稳定后推进，US3 在 US1/US2 的身份目录、组织模型和授权边界语义完成后推进
- **Phase 6 (Polish)**: 依赖已完成的用户故事范围
- **Release / Merge**: 依赖远程推送、PR 更新、评审完成与用户明确同意

### User Story Dependencies

- **US1 (P1)**: 无用户故事前置依赖，Foundational 完成后可立即开始
- **US2 (P1)**: 依赖 US1 产出的身份源与用户来源语义，才能形成组织模型和租户映射闭环
- **US3 (P2)**: 依赖 Foundational 的角色与会话基础模型，以及 US1/US2 的身份来源和租户边界语义

### Parallel Opportunities

- **Phase 1**: T005/T006 可并行
- **Phase 2**: T010/T011/T012/T014/T015/T017 可并行
- **US1**: T019/T020/T021 并行，T022/T023 并行，T026/T027/T028 并行
- **US2**: T030/T031/T032 并行，T033/T034 并行，T037/T038/T039 并行
- **US3**: T041/T042/T043 并行，T044/T045/T046 并行，T049/T050/T051/T052 并行

---

## Parallel Example: User Story 1

```bash
# 并行测试任务
Task: "T019 [US1] backend/tests/contract/identity_source_contract_test.go"
Task: "T020 [US1] backend/tests/integration/identity_source_lifecycle_test.go"
Task: "T021 [US1] frontend/src/features/identity-tenancy/pages/IdentitySourcePage.test.tsx"

# 并行实现任务
Task: "T022 [US1] backend/internal/service/identitytenancy/identity_source_service.go"
Task: "T023 [US1] backend/internal/service/identitytenancy/login_mode_service.go"
Task: "T026 [US1] frontend/src/services/identityTenancy.ts"
```

---

## Implementation Strategy

### MVP First (US1 Only)

1. 完成 Phase 0-2
2. 完成 US1（Phase 3）
3. 按 Independent Test 验证 US1
4. 产出阶段性演示或 PR 更新

### Incremental Delivery

1. 先交付 US1：身份源接入、本地账号并存、统一登录方式和会话治理
2. 再交付 US2：组织、团队、用户组、工作空间和项目之间的租户关系建模
3. 最后交付 US3：细粒度 RBAC、委派、临时授权、回收和访问风险视图
4. 最终执行 Phase 6 文档与验证收尾

### Notes

- `[P]` 任务代表可并行，但仍需满足前置依赖
- 每个用户故事都可独立验收
- 所有提交说明与 PR 摘要必须为中文
- 未获用户明确同意前禁止合并
