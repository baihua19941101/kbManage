# Tasks: 多集群 Kubernetes 工作负载运维控制面

**Input**: Design documents from `/specs/003-workload-operations-control-plane/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/, quickstart.md

**Tests**: 本特性明确要求测试任务，必须覆盖后端 contract/integration 测试以及前端 Vitest 页面测试。

**Organization**: 任务按用户故事组织，确保 US1、US2、US3 可以独立实现与独立验收。

**Constitutional Gates**: 功能分支校验、数据库备份证据、国内源配置、中文 PR、远程推送和明确的用户合并授权必须在实现与交付时满足。

## Format: `[ID] [P?] [Story] Description`

- **[P]**: 可并行执行（不同文件、无未完成依赖）
- **[Story]**: 对应的用户故事标签（US1、US2、US3）
- 描述中包含明确文件路径，便于直接执行

## Phase 0: Governance Gates

**Purpose**: 在真正开始实现 003 之前完成宪章要求的治理门槛

- [X] T001 在 `artifacts/003-workload-operations-control-plane/branch-check.txt` 记录当前分支、远程仓库和“未经用户批准不得合并”的门槛
- [X] T002 核对并补充 003 本轮数据库备份证据到 `artifacts/003-workload-operations-control-plane/backup-manifest.txt`
- [X] T003 在 `artifacts/003-workload-operations-control-plane/mirror-and-remote-check.txt` 记录 `GOPROXY=https://goproxy.cn,direct`、`https://registry.npmmirror.com`、联调镜像来源和 `git@github.com:baihua19941101/kbManage.git` PR 流程
- [X] T004 更新 `specs/003-workload-operations-control-plane/spec.md`、`specs/003-workload-operations-control-plane/plan.md` 和 `specs/003-workload-operations-control-plane/quickstart.md` 的执行状态，注明已进入 tasks 阶段

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: 为 003 建立最小的模块骨架、配置入口和依赖准备

- [X] T005 创建后端工作负载运维模块骨架 `backend/internal/service/workloadops/`、`backend/internal/api/handler/workload_ops_handler.go`、`backend/internal/api/router/workload_ops_routes.go`
- [X] T006 [P] 创建终端执行与会话基础目录 `backend/internal/kube/exec/`、`backend/internal/worker/`
- [X] T007 [P] 创建前端模块骨架 `frontend/src/features/workload-ops/`、`frontend/src/services/workloadOps.ts` 和路由占位到 `frontend/src/app/router.tsx`
- [X] T008 在 `backend/config/config.example.yaml`、`backend/config/config.dev.yaml` 和 `README.md` 增加 `workloadops.actions.*`、`workloadops.batch.*`、`workloadops.terminal.*`、`workloadops.audit.*` 配置说明

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: 搭建 003 的共享数据模型、动作语义、路由和权限基础，阻塞所有用户故事

**⚠️ CRITICAL**: US1、US2、US3 都必须在本阶段完成后才能开始

- [X] T009 新增 003 数据库迁移 `backend/migrations/0005_workload_operations_core.sql`，落库批量任务、批量子项、终端会话和动作扩展字段
- [X] T010 [P] 在 `backend/internal/domain/workloadops.go` 定义 `WorkloadActionRequest`、`BatchOperationTask`、`BatchOperationItem`、`TerminalSession`
- [X] T011 [P] 创建仓储实现 `backend/internal/repository/workload_action_repository.go`、`backend/internal/repository/batch_operation_repository.go`、`backend/internal/repository/terminal_session_repository.go`
- [X] T012 [P] 在 `backend/internal/service/auth/permission_service.go`、`backend/internal/api/middleware/authorization.go` 定义 003 权限语义 `workloadops:read`、`workloadops:execute`、`workloadops:terminal`、`workloadops:rollback`、`workloadops:batch`
- [X] T013 [P] 在 `backend/internal/api/router/workload_ops_routes.go`、`backend/internal/api/router/router.go` 注册 003 API 路由骨架
- [X] T014 [P] 在 `backend/internal/repository/redis.go`、`backend/internal/service/workloadops/progress_cache.go`、`backend/internal/service/workloadops/session_cache.go` 建立动作进度、批量协调和终端短会话缓存
- [X] T015 在 `backend/internal/service/workloadops/service.go`、`backend/internal/service/workloadops/scope_service.go` 建立统一工作负载运维授权和范围过滤入口
- [X] T016 [P] 在 `frontend/src/services/api/types.ts`、`frontend/src/services/api/client.ts`、`frontend/src/app/queryClient.ts` 增加 003 共享类型、查询 key 和错误归一化
- [X] T017 在 `frontend/src/app/AuthorizedMenu.tsx`、`frontend/src/app/ProtectedRoute.tsx`、`frontend/src/features/auth/store.ts` 接入工作负载运维入口和基础权限门控

**Checkpoint**: 003 的共享模型、权限、路由和缓存骨架就绪，用户故事可以开始实现

---

## Phase 3: User Story 1 - 单资源运维诊断入口 (Priority: P1) 🎯 MVP

**Goal**: 提供围绕单个工作负载的统一运维上下文、实例列表、日志联动和容器终端入口，形成最小可用诊断闭环

**Independent Test**: 在至少两个已接入并授权的集群中选择一个异常工作负载，用户能够在同一资源页面完成状态查看、发布进度跟踪、实例下钻、日志查看和容器终端进入，并据此识别异常实例。

### Tests for User Story 1

- [X] T018 [P] [US1] 编写后端契约测试 `backend/tests/contract/workload_ops_context_contract_test.go`、`backend/tests/contract/workload_ops_instances_contract_test.go`、`backend/tests/contract/workload_ops_terminal_contract_test.go`
- [X] T019 [P] [US1] 编写后端集成测试 `backend/tests/integration/workload_ops_context_test.go`、`backend/tests/integration/workload_terminal_session_test.go`
- [X] T020 [P] [US1] 编写前端 Vitest 页面测试 `frontend/src/features/workload-ops/pages/WorkloadOperationsPage.test.tsx`、`frontend/src/features/workload-ops/components/InstanceListPanel.test.tsx`

### Implementation for User Story 1

- [X] T021 [P] [US1] 实现工作负载运维上下文聚合服务 `backend/internal/service/workloadops/context_service.go`
- [X] T022 [P] [US1] 实现实例列表与状态聚合 `backend/internal/service/workloadops/instance_service.go`、`backend/internal/kube/adapter/`
- [X] T023 [P] [US1] 实现终端会话创建、关闭和超时处理 `backend/internal/service/workloadops/terminal_service.go`、`backend/internal/kube/exec/`
- [X] T024 [US1] 在 `backend/internal/api/handler/workload_ops_handler.go` 和 `backend/internal/api/router/workload_ops_routes.go` 落地 `/workload-ops/resources/context`、`/workload-ops/resources/instances`、`/workload-ops/terminal/sessions`
- [X] T025 [P] [US1] 创建前端服务层 `frontend/src/services/workloadOps.ts`
- [X] T026 [P] [US1] 实现工作负载运维页和实例面板 `frontend/src/features/workload-ops/pages/WorkloadOperationsPage.tsx`、`frontend/src/features/workload-ops/components/InstanceListPanel.tsx`
- [X] T027 [US1] 实现终端入口与会话状态提示 `frontend/src/features/workload-ops/components/TerminalSessionDrawer.tsx`
- [X] T028 [US1] 在 `frontend/src/features/resources/components/ResourceDetailDrawer.tsx`、`frontend/src/features/observability/pages/ResourceContextPage.tsx`、`frontend/src/app/router.tsx` 打通到 003 运维上下文的导航入口

**Checkpoint**: US1 完成后，平台应具备可独立演示的单资源运维诊断 MVP

---

## Phase 4: User Story 2 - 工作负载动作执行与发布恢复 (Priority: P1)

**Goal**: 提供单资源动作、批量动作、发布历史查看和版本回滚能力，形成工作负载运维写路径闭环

**Independent Test**: 选择一个被授权的工作负载并执行扩缩容、重启或回滚动作后，用户能够看到影响预览、执行状态流转、最终结果以及失败原因；对多个资源执行批量动作时能够区分整体结果和单项结果。

### Tests for User Story 2

- [X] T029 [P] [US2] 编写后端契约测试 `backend/tests/contract/workload_ops_actions_contract_test.go`、`backend/tests/contract/workload_ops_batches_contract_test.go`、`backend/tests/contract/workload_ops_revisions_contract_test.go`
- [X] T030 [P] [US2] 编写后端集成测试 `backend/tests/integration/workload_action_execution_test.go`、`backend/tests/integration/workload_batch_operation_test.go`、`backend/tests/integration/workload_rollback_test.go`
- [X] T031 [P] [US2] 编写前端 Vitest 页面测试 `frontend/src/features/workload-ops/pages/BatchOperationPage.test.tsx`、`frontend/src/features/workload-ops/components/RollbackDialog.test.tsx`

### Implementation for User Story 2

- [X] T032 [P] [US2] 实现 Kubernetes 工作负载动作执行器 `backend/internal/service/workloadops/executor.go`，覆盖 `scale`、`restart`、`redeploy`、`replace-instance`、`rollback`
- [X] T033 [P] [US2] 实现发布历史与回滚目标识别 `backend/internal/service/workloadops/revision_service.go`
- [X] T034 [P] [US2] 实现批量任务编排与子项结果归集 `backend/internal/service/workloadops/batch_service.go`、`backend/internal/worker/workload_batch_worker.go`
- [X] T035 [US2] 在 `backend/internal/service/workloadops/service.go`、`backend/internal/worker/operation_worker.go` 打通单资源动作状态流转、失败归一化和幂等处理
- [X] T036 [US2] 在 `backend/internal/api/handler/workload_ops_handler.go`、`backend/internal/api/router/workload_ops_routes.go` 落地 `/workload-ops/actions`、`/workload-ops/actions/{id}`、`/workload-ops/batches`、`/workload-ops/batches/{id}`、`/workload-ops/resources/revisions`
- [X] T037 [P] [US2] 实现前端动作提交与轮询 `frontend/src/features/workload-ops/components/ActionConfirmDrawer.tsx`、`frontend/src/features/workload-ops/components/BatchOperationDrawer.tsx`
- [X] T038 [P] [US2] 实现发布历史与回滚交互 `frontend/src/features/workload-ops/components/RevisionHistoryPanel.tsx`、`frontend/src/features/workload-ops/components/RollbackDialog.tsx`
- [X] T039 [US2] 实现批量任务结果页 `frontend/src/features/workload-ops/pages/BatchOperationPage.tsx`

**Checkpoint**: US2 完成后，平台应具备真实可追踪的工作负载动作执行与恢复链路

---

## Phase 5: User Story 3 - 权限隔离与高风险审计闭环 (Priority: P2)

**Goal**: 对工作负载视图、终端、单体动作、批量动作和回滚统一执行授权校验，并形成完整审计闭环

**Independent Test**: 为两个不同授权范围的用户分别登录平台后，他们只能访问各自范围内的工作负载、实例、终端和动作入口；执行高风险动作、进入终端或发起回滚后，审计人员能够检索到完整记录。

### Tests for User Story 3

- [X] T040 [P] [US3] 编写后端契约测试 `backend/tests/contract/workload_ops_access_control_contract_test.go`
- [X] T041 [P] [US3] 编写后端集成测试 `backend/tests/integration/workload_ops_scope_authorization_test.go`、`backend/tests/integration/workload_ops_audit_test.go`
- [X] T042 [P] [US3] 编写前端 Vitest 页面测试 `frontend/src/features/workload-ops/pages/WorkloadOperationsAccessGate.test.tsx`

### Implementation for User Story 3

- [X] T043 [P] [US3] 在 `backend/internal/service/auth/scope_authorizer.go`、`backend/internal/service/workloadops/scope_service.go` 落地工作负载级范围校验
- [X] T044 [P] [US3] 在 `backend/internal/service/audit/event_writer.go`、`backend/internal/service/audit/service.go` 增加 `workloadops.*` 审计动作类型，覆盖终端打开/关闭、回滚、批量任务和高风险动作
- [X] T045 [US3] 在 `backend/internal/service/workloadops/terminal_service.go` 中落实终端审计边界，仅记录会话建立、关闭、目标容器、操作者、持续时长、结束原因
- [X] T046 [US3] 在 `backend/internal/api/middleware/authorization.go`、`backend/internal/api/router/workload_ops_routes.go` 区分只读、终端、回滚和批量高风险动作权限
- [X] T047 [P] [US3] 在 `frontend/src/app/AuthorizedMenu.tsx`、`frontend/src/features/auth/store.ts` 实现 003 导航与动作级门控
- [X] T048 [US3] 在 `frontend/src/features/workload-ops/pages/WorkloadOperationsPage.tsx`、`frontend/src/features/workload-ops/components/TerminalSessionDrawer.tsx`、`frontend/src/features/workload-ops/components/RollbackDialog.tsx` 实现未授权空态、只读态和权限回收处理

**Checkpoint**: US3 完成后，003 的所有工作负载运维能力都应运行在现有租户隔离和审计模型之下

---

## Final Phase: Polish & Delivery Readiness

**Purpose**: 收尾文档、验证、性能与交付准备

- [X] T049 [P] 清理共享类型与命名，在 `backend/internal/service/workloadops/`、`frontend/src/features/workload-ops/`、`frontend/src/services/workloadOps.ts` 做收敛
- [X] T050 [P] 刷新配置与启动文档，在 `README.md`、`backend/config/config.example.yaml`、`backend/config/config.dev.yaml`、`frontend/.env.example`、`frontend/.env.development` 补齐 003 说明
- [X] T051 [P] 记录验证基线到 `artifacts/003-workload-operations-control-plane/verification.md`、`artifacts/003-workload-operations-control-plane/quickstart-validation.md`、`artifacts/003-workload-operations-control-plane/repro-workloadops-smoke.sh`
- [X] T052 在 `artifacts/003-workload-operations-control-plane/pr-summary.md` 和 `artifacts/003-workload-operations-control-plane/pr-readiness.md` 准备中文 PR 摘要、备份证据、测试说明和风险清单

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 0: Governance Gates**: 无依赖，但 003 真正开始编码前必须完成
- **Phase 1: Setup**: 依赖 Phase 0，负责建立模块骨架和配置入口
- **Phase 2: Foundational**: 依赖 Phase 1，阻塞所有用户故事
- **Phase 3: US1**: 依赖 Phase 2，是 003 的 MVP
- **Phase 4: US2**: 依赖 Phase 2，最佳实践是在 US1 的读路径和终端入口稳定后推进
- **Phase 5: US3**: 依赖 US1 和 US2，因为需要对已存在的工作负载视图、终端和动作链路统一施加授权与审计
- **Final Phase**: 依赖目标用户故事完成

### User Story Dependencies

- **US1**: 仅依赖共享基础设施，是最小可演示范围
- **US2**: 依赖共享基础设施；可不等待 US1 全部 UI 完成，但会复用上下文、终端和动作基座
- **US3**: 依赖 US1、US2，因为权限隔离和审计需要覆盖已存在的终端、动作、回滚和批量任务

### Parallel Opportunities

- Phase 1 中，`T006`、`T007`、`T008` 可并行
- Phase 2 中，`T010`、`T011`、`T012`、`T013`、`T014`、`T016` 可并行
- US1 中，`T018`、`T019`、`T020` 可并行；`T021`、`T022`、`T023` 可并行；`T025`、`T026`、`T027` 可并行
- US2 中，`T029`、`T030`、`T031` 可并行；`T032`、`T033`、`T034` 可并行；`T037`、`T038`、`T039` 可并行
- US3 中，`T040`、`T041`、`T042` 可并行；`T043`、`T044`、`T047` 可并行
- Final Phase 中，`T049`、`T050`、`T051` 可并行

---

## Parallel Example: User Story 1

```bash
# 并行启动 US1 的测试任务
Task: "T018 [US1] 后端契约测试"
Task: "T019 [US1] 后端集成测试"
Task: "T020 [US1] 前端 Vitest 页面测试"

# 并行实现 US1 的三类后端能力
Task: "T021 [US1] 工作负载运维上下文聚合服务"
Task: "T022 [US1] 实例列表与状态聚合"
Task: "T023 [US1] 终端会话创建、关闭和超时处理"

# 并行实现 US1 的前端视图
Task: "T025 [US1] 前端服务层"
Task: "T026 [US1] 工作负载运维页和实例面板"
Task: "T027 [US1] 终端入口与会话状态提示"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. 完成 Phase 0: Governance Gates
2. 完成 Phase 1: Setup
3. 完成 Phase 2: Foundational
4. 完成 Phase 3: US1
5. **STOP and VALIDATE**：按 US1 独立验收标准验证单资源运维上下文、实例诊断、日志联动和终端入口

### Incremental Delivery

1. 完成 Governance + Setup + Foundational，形成 003 共享底座
2. 交付 US1，形成可演示的单资源运维诊断 MVP
3. 交付 US2，补齐动作执行、批量任务和回滚恢复
4. 交付 US3，补齐工作负载运维权限隔离与高风险审计闭环
5. 完成 Final Phase，准备 PR 证据与交付文档

### Notes

- 首期范围严格限定在工作负载级诊断、单体动作、批量动作、发布历史、回滚和终端访问
- 全局日志中心、统一监控告警、GitOps、Helm、策略治理和集群生命周期能力禁止混入本任务清单
- 如果后续需要扩展命令留痕、终端录屏或更高级发布治理能力，应进入新的 feature，而不是继续扩大 003 范围
