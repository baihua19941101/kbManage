# Tasks: 多集群 Kubernetes 安全与策略治理中心

**Input**: Design documents from `/specs/005-security-and-policy/`  
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/openapi.yaml, quickstart.md

**Tests**: 本特性包含明确的验收场景、权限隔离与审计闭环要求，任务清单包含后端契约/集成测试与前端页面测试任务。  
**Organization**: 任务按用户故事分组，保证每个故事可独立实现、独立验证。

**Constitutional Gates**: 必须满足功能分支、数据库备份证据、国内依赖源配置、中文 PR、远程推送、用户同意后合并。

## Format: `[ID] [P?] [Story] Description`

- `[P]`: 可并行执行（不同文件、无前置依赖）
- `[Story]`: 任务归属用户故事（US1/US2/US3）
- 每条任务必须包含明确文件路径

## Phase 0: Governance Gates

**Purpose**: 完成宪章门槛与实施前证据准备

- [X] T001 在 `artifacts/005-security-and-policy/branch-check.txt` 记录当前分支、禁止在 `main/master` 开发、以及“未获用户同意不得合并”门槛
- [X] T002 执行 005 开发前数据库备份并在 `artifacts/005-security-and-policy/backup-manifest.txt` 记录命令、时间戳、产物路径与恢复抽样验证结果
- [X] T003 在 `artifacts/005-security-and-policy/mirror-and-remote-check.txt` 记录 `GOPROXY=https://goproxy.cn,direct`、`https://registry.npmmirror.com`、策略相关镜像来源与 `git@github.com:baihua19941101/kbManage.git` PR 流程
- [X] T004 更新 `specs/005-security-and-policy/spec.md`、`specs/005-security-and-policy/plan.md`、`specs/005-security-and-policy/quickstart.md` 的执行状态，注明已进入 tasks 阶段

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: 建立 005 模块骨架、配置入口和导航占位

- [X] T005 创建后端策略治理模块骨架 `backend/internal/service/securitypolicy/`、`backend/internal/api/handler/security_policy_handler.go`、`backend/internal/api/router/security_policy_routes.go`
- [X] T006 [P] 创建策略集成层目录与占位 `backend/internal/integration/policy/engine/`、`backend/internal/integration/policy/distribution/`
- [X] T007 [P] 创建前端模块骨架 `frontend/src/features/security-policy/`、`frontend/src/services/securityPolicy.ts` 与路由占位 `frontend/src/app/router.tsx`
- [X] T008 在 `backend/config/config.example.yaml`、`backend/config/config.dev.yaml`、`frontend/.env.example`、`frontend/.env.development`、`README.md` 增加 `securityPolicy.*` 配置说明
- [X] T009 在 `frontend/src/app/AuthorizedMenu.tsx`、`frontend/src/app/ProtectedRoute.tsx` 接入安全与策略治理入口门禁占位

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: 完成所有用户故事共享且阻塞性的基础能力

**⚠️ CRITICAL**: US1/US2/US3 必须在本阶段完成后才可开始

- [X] T010 新增 005 数据库迁移 `backend/migrations/0007_security_policy_core.sql`，落库策略、策略版本、策略分配、命中记录、例外申请、整改动作表
- [X] T011 [P] 在 `backend/internal/domain/security_policy.go` 定义 `SecurityPolicy`、`PolicyVersion`、`PolicyAssignment`、`PolicyHitRecord`、`ExceptionRequest`、`RemediationAction`、`PolicyDistributionTask`
- [X] T012 [P] 创建仓储实现 `backend/internal/repository/security_policy_repository.go`、`backend/internal/repository/policy_assignment_repository.go`、`backend/internal/repository/policy_hit_repository.go`、`backend/internal/repository/policy_exception_repository.go`
- [X] T013 [P] 在 `backend/internal/integration/policy/engine/provider.go`、`backend/internal/integration/policy/distribution/provider.go` 定义策略评估与分发抽象以及错误归一化模型
- [X] T014 在 `backend/internal/service/securitypolicy/service.go`、`backend/internal/service/securitypolicy/scope_service.go` 建立统一策略授权与范围过滤入口
- [X] T015 [P] 在 `backend/internal/repository/redis.go`、`backend/internal/service/securitypolicy/distribution_cache.go`、`backend/internal/service/securitypolicy/exception_cache.go` 建立分发进度、例外时效与幂等缓存
- [X] T016 在 `backend/internal/api/router/security_policy_routes.go`、`backend/internal/api/router/router.go` 注册 005 API 路由骨架
- [X] T017 [P] 在 `backend/internal/service/auth/permission_service.go`、`backend/internal/api/middleware/authorization.go` 增加 005 权限语义 `policy:read`、`policy:manage`、`policy:assign`、`policy:approve-exception`、`policy:remediate`
- [X] T018 [P] 在 `frontend/src/services/api/types.ts`、`frontend/src/services/api/client.ts`、`frontend/src/app/queryClient.ts` 增加 005 共享类型、查询 key 和错误归一化
- [X] T019 在 `backend/internal/service/audit/event_writer.go`、`backend/internal/service/audit/service.go` 预置 `policy.*` 审计动作类型与查询维度映射

**Checkpoint**: 基础能力完成，可开始用户故事实现

---

## Phase 3: User Story 1 - 统一策略中心与分层策略管理 (Priority: P1) 🎯 MVP

**Goal**: 提供平台级/工作空间级/项目级策略定义、分层展示与范围分配能力。  
**Independent Test**: 创建平台级与项目级策略并分配到不同范围后，可清晰查看策略层级关系、作用范围和最终适用集合。

### Tests for User Story 1

- [X] T020 [P] [US1] 编写后端契约测试 `backend/tests/contract/security_policy_management_contract_test.go`，覆盖策略 CRUD、层级字段与状态约束
- [X] T021 [P] [US1] 编写后端集成测试 `backend/tests/integration/security_policy_scope_modeling_test.go`，覆盖平台/工作空间/项目分层与范围分配隔离
- [X] T022 [P] [US1] 编写前端 Vitest 页面测试 `frontend/src/features/security-policy/pages/PolicyCenterPage.test.tsx`、`frontend/src/features/security-policy/components/PolicyScopeDrawer.test.tsx`

### Implementation for User Story 1

- [X] T023 [P] [US1] 实现策略定义服务 `backend/internal/service/securitypolicy/policy_service.go`，覆盖创建、更新、启停、归档和版本快照
- [X] T024 [P] [US1] 实现策略分配服务 `backend/internal/service/securitypolicy/assignment_service.go`，覆盖集群/命名空间/项目/资源类型绑定
- [X] T025 [US1] 实现策略层级聚合查询 `backend/internal/service/securitypolicy/hierarchy_service.go`，输出最终适用策略集合与来源层级
- [X] T026 [US1] 在 `backend/internal/api/handler/security_policy_handler.go`、`backend/internal/api/router/security_policy_routes.go` 落地 `/security-policies`、`/security-policies/{policyId}`、`/security-policies/{policyId}/assignments`
- [X] T027 [P] [US1] 实现前端服务层 `frontend/src/services/securityPolicy.ts`，覆盖策略列表/详情/创建/分配接口
- [X] T028 [P] [US1] 实现策略中心页 `frontend/src/features/security-policy/pages/PolicyCenterPage.tsx` 和策略表格 `frontend/src/features/security-policy/components/PolicyTable.tsx`
- [X] T029 [P] [US1] 实现策略编辑与分配抽屉 `frontend/src/features/security-policy/components/PolicyEditorDrawer.tsx`、`frontend/src/features/security-policy/components/PolicyScopeDrawer.tsx`
- [X] T030 [US1] 在 `frontend/src/features/security-policy/pages/PolicyCenterPage.tsx` 实现策略层级关系与最终适用策略可视化
- [X] T031 [US1] 在 `frontend/src/app/router.tsx`、`frontend/src/app/AuthorizedMenu.tsx` 打通 005 页面路由与导航
- [X] T032 [US1] 在 `frontend/src/features/security-policy/pages/PolicyCenterPage.tsx` 增加空态、错误态和权限不足态提示

**Checkpoint**: US1 完整可测，可作为 005 MVP 交付

---

## Phase 4: User Story 2 - 准入控制模式与分阶段启用 (Priority: P1)

**Goal**: 提供策略执行模式切换、灰度启用和例外时效管理能力。  
**Independent Test**: 将策略从 `warn` 灰度切换到 `enforce` 后可观察命中变化；例外在有效期内生效并在到期后自动恢复原约束。

### Tests for User Story 2

- [X] T033 [P] [US2] 编写后端契约测试 `backend/tests/contract/security_policy_enforcement_contract_test.go`、`backend/tests/contract/security_policy_exception_contract_test.go`
- [X] T034 [P] [US2] 编写后端集成测试 `backend/tests/integration/security_policy_mode_rollout_test.go`、`backend/tests/integration/security_policy_exception_lifecycle_test.go`
- [X] T035 [P] [US2] 编写前端 Vitest 页面测试 `frontend/src/features/security-policy/pages/PolicyRolloutPage.test.tsx`、`frontend/src/features/security-policy/components/ExceptionReviewDrawer.test.tsx`

### Implementation for User Story 2

- [X] T036 [P] [US2] 实现策略模式切换服务 `backend/internal/service/securitypolicy/enforcement_service.go`，支持 `audit/alert/warn/enforce`
- [X] T037 [P] [US2] 实现灰度分阶段启用服务 `backend/internal/service/securitypolicy/rollout_service.go`，支持 `pilot/canary/broad/full`
- [X] T038 [P] [US2] 实现例外申请与审批服务 `backend/internal/service/securitypolicy/exception_service.go`，覆盖申请、审批、拒绝、撤销、到期失效
- [X] T039 [US2] 实现例外到期回收 worker `backend/internal/worker/policy_exception_expiry_worker.go`
- [X] T040 [US2] 在 `backend/internal/api/handler/security_policy_handler.go`、`backend/internal/api/router/security_policy_routes.go` 落地 `/security-policies/{policyId}/mode-switch`、`/security-policies/hits/{hitId}/exceptions`、`/security-policies/exceptions`、`/security-policies/exceptions/{exceptionId}/review`
- [X] T041 [P] [US2] 扩展前端服务与 hooks `frontend/src/services/securityPolicy.ts`、`frontend/src/features/security-policy/hooks/usePolicyRollout.ts`
- [X] T042 [P] [US2] 实现模式切换与灰度页面 `frontend/src/features/security-policy/pages/PolicyRolloutPage.tsx`、`frontend/src/features/security-policy/components/ModeSwitchDrawer.tsx`
- [X] T043 [P] [US2] 实现例外申请/审批交互 `frontend/src/features/security-policy/components/ExceptionRequestDrawer.tsx`、`frontend/src/features/security-policy/components/ExceptionReviewDrawer.tsx`
- [X] T044 [US2] 在 `frontend/src/features/security-policy/pages/PolicyRolloutPage.tsx` 落地例外状态（待审批/生效/过期/撤销）可视化与到期提示

**Checkpoint**: US2 可独立验证模式切换、灰度与例外生命周期

---

## Phase 5: User Story 3 - 违规治理闭环与审计追踪 (Priority: P2)

**Goal**: 提供违规查询、整改跟踪与策略治理审计闭环。  
**Independent Test**: 在产生违规后可查询风险级别与整改状态，并按策略变更与处置链路检索完整审计记录。

### Tests for User Story 3

- [X] T045 [P] [US3] 编写后端契约测试 `backend/tests/contract/security_policy_hits_contract_test.go`、`backend/tests/contract/security_policy_audit_contract_test.go`
- [X] T046 [P] [US3] 编写后端集成测试 `backend/tests/integration/security_policy_violation_lifecycle_test.go`、`backend/tests/integration/security_policy_audit_query_test.go`
- [X] T047 [P] [US3] 编写前端 Vitest 页面测试 `frontend/src/features/security-policy/pages/ViolationCenterPage.test.tsx`、`frontend/src/features/audit/pages/SecurityPolicyAuditPage.test.tsx`

### Implementation for User Story 3

- [X] T048 [P] [US3] 实现违规查询服务 `backend/internal/service/securitypolicy/violation_service.go`，覆盖风险级别、范围、模式、时间筛选
- [X] T049 [P] [US3] 实现整改跟踪服务 `backend/internal/service/securitypolicy/remediation_service.go`，覆盖整改状态流转与处理记录
- [X] T050 [US3] 在 `backend/internal/api/handler/security_policy_handler.go`、`backend/internal/api/router/security_policy_routes.go` 落地 `/security-policies/hits`、`/security-policies/hits/{hitId}/remediation`
- [X] T051 [US3] 在 `backend/internal/service/audit/event_writer.go`、`backend/internal/service/audit/service.go` 完成 `policy.*` 审计写入和查询聚合
- [X] T052 [US3] 在 `backend/internal/api/handler/audit_handler.go`、`backend/internal/api/router/audit_routes.go` 落地 `/audit/security-policies/events`
- [X] T053 [P] [US3] 实现前端违规中心页 `frontend/src/features/security-policy/pages/ViolationCenterPage.tsx` 和组件 `frontend/src/features/security-policy/components/ViolationTable.tsx`
- [X] T054 [P] [US3] 实现整改更新交互 `frontend/src/features/security-policy/components/RemediationUpdateDrawer.tsx` 和风险分布视图 `frontend/src/features/security-policy/components/ViolationRiskChart.tsx`
- [X] T055 [US3] 实现策略审计页面 `frontend/src/features/audit/pages/SecurityPolicyAuditPage.tsx` 并扩展 `frontend/src/services/audit.ts` 查询接口

**Checkpoint**: US3 完成后形成策略治理与违规处置可追踪闭环

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: 收敛质量、文档、验证证据与 PR 交付材料

- [X] T056 [P] 收敛命名与共享类型，在 `backend/internal/service/securitypolicy/`、`frontend/src/features/security-policy/`、`frontend/src/services/securityPolicy.ts` 清理重复字段与错误文案
- [X] T057 [P] 刷新配置与启动文档，在 `README.md`、`backend/config/config.example.yaml`、`backend/config/config.dev.yaml`、`frontend/.env.example`、`frontend/.env.development` 补齐 005 说明
- [X] T058 [P] 记录验证基线到 `artifacts/005-security-and-policy/verification.md`、`artifacts/005-security-and-policy/quickstart-validation.md`、`artifacts/005-security-and-policy/repro-securitypolicy-smoke.sh`
- [X] T059 在 `artifacts/005-security-and-policy/pr-summary.md`、`artifacts/005-security-and-policy/pr-readiness.md` 准备中文 PR 摘要、备份证据、测试说明、风险清单与用户合并确认项

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 0 (Governance Gates)**: 无依赖，必须先完成
- **Phase 1 (Setup)**: 依赖 Phase 0
- **Phase 2 (Foundational)**: 依赖 Phase 1，阻塞所有用户故事
- **Phase 3/4/5 (User Stories)**: 依赖 Phase 2；US1 优先作为 MVP，US2/US3 可在 US1 主干稳定后并行推进
- **Phase 6 (Polish)**: 依赖已完成的用户故事范围

### User Story Dependencies

- **US1 (P1)**: 无用户故事前置依赖，Foundational 完成后可立即开始
- **US2 (P1)**: 依赖 US1 的策略定义与分配基础对象
- **US3 (P2)**: 依赖 US1/US2 产生命中、例外与状态数据后形成完整闭环

### Parallel Opportunities

- Phase 1: T006/T007 可并行
- Phase 2: T011/T012/T013/T015/T017/T018 可并行
- US1: T020/T021/T022 并行，T023/T024 并行，T027/T028/T029 并行
- US2: T033/T034/T035 并行，T036/T037/T038 并行，T041/T042/T043 并行
- US3: T045/T046/T047 并行，T048/T049 并行，T053/T054 并行

---

## Parallel Example: User Story 1

```bash
# 并行测试任务
Task: "T020 [US1] backend/tests/contract/security_policy_management_contract_test.go"
Task: "T021 [US1] backend/tests/integration/security_policy_scope_modeling_test.go"
Task: "T022 [US1] frontend/src/features/security-policy/pages/PolicyCenterPage.test.tsx"

# 并行实现任务
Task: "T023 [US1] backend/internal/service/securitypolicy/policy_service.go"
Task: "T024 [US1] backend/internal/service/securitypolicy/assignment_service.go"
Task: "T027 [US1] frontend/src/services/securityPolicy.ts"
```

---

## Implementation Strategy

### MVP First (US1 Only)

1. 完成 Phase 0-2
2. 完成 US1（Phase 3）
3. 按 Independent Test 验证 US1
4. 产出阶段性演示或 PR 更新

### Incremental Delivery

1. 先交付 US1：统一策略中心与分层分配
2. 再交付 US2：执行模式/灰度/例外
3. 最后交付 US3：违规与审计闭环
4. 最终执行 Phase 6 文档与验证收尾

### Notes

- `[P]` 任务代表可并行，但仍需满足前置依赖
- 每个用户故事都可独立验收
- 所有提交说明与 PR 摘要必须为中文
- 未获用户明确同意前禁止合并
