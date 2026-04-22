# Tasks: 企业级治理报表与产品化交付收尾

**Input**: Design documents from `/specs/012-enterprise-polish/`  
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/openapi.yaml, quickstart.md

**Tests**: 本特性包含明确的独立验收标准、深度审计、治理报表和交付清单要求，任务清单包含后端契约/集成测试与前端页面测试任务。  
**Organization**: 任务按用户故事分组，保证每个故事可独立实现、独立验证。

**Constitutional Gates**: 必须满足功能分支、数据库备份证据、国内依赖源配置、中文 PR、远程推送、用户同意后合并。

## Format: `[ID] [P?] [Story] Description`

- `[P]`: 可并行执行（不同文件、无前置依赖）
- `[Story]`: 任务归属用户故事（US1/US2/US3）
- 每条任务必须包含明确文件路径

## Phase 0: Governance Gates

**Purpose**: 完成宪章门槛与实施前证据准备

- [X] T001 在 `artifacts/012-enterprise-polish/branch-check.txt` 记录当前分支、禁止在 `main/master` 开发、`011` 已合并完成以及“未获用户同意不得合并”门槛
- [X] T002 执行 012 开发前数据库备份并在 `artifacts/012-enterprise-polish/backup-manifest.txt` 记录命令、时间戳、产物路径与恢复抽样验证结果
- [X] T003 在 `artifacts/012-enterprise-polish/mirror-and-remote-check.txt` 记录 `GOPROXY=https://goproxy.cn,direct`、`https://registry.npmmirror.com`、报表/导出依赖镜像策略与 `git@github.com:baihua19941101/kbManage.git` PR 流程

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: 建立 012 模块骨架、配置入口和导航占位

- [X] T004 创建后端 enterprise 模块骨架 `backend/internal/service/enterprise/`、`backend/internal/api/handler/enterprise_handler.go`、`backend/internal/api/router/enterprise_routes.go`
- [X] T005 [P] 创建深度审计、报表生成与交付目录适配层目录 `backend/internal/integration/enterprise/`、`backend/internal/integration/enterprise/audit_provider.go`、`backend/internal/integration/enterprise/report_builder.go`、`backend/internal/integration/enterprise/delivery_catalog.go`
- [X] T006 [P] 创建前端模块骨架 `frontend/src/features/enterprise-polish/`、`frontend/src/services/enterprisePolish.ts` 与路由占位到 `frontend/src/app/router.tsx`
- [X] T007 在 `backend/config/config.example.yaml`、`backend/config/config.dev.yaml`、`frontend/.env.example`、`frontend/.env.development`、`README.md` 增加 `enterprisePolish.*` 配置说明
- [X] T008 在 `frontend/src/app/AuthorizedMenu.tsx`、`frontend/src/app/ProtectedRoute.tsx`、`frontend/src/features/auth/store.ts` 接入企业治理与交付入口门禁占位

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: 完成所有用户故事共享且阻塞性的基础能力

**⚠️ CRITICAL**: US1/US2/US3 必须在本阶段完成后才可开始

- [X] T009 新增 012 数据库迁移 `backend/migrations/0014_enterprise_polish_core.sql`，落库权限变更链路、关键操作追踪、跨团队授权快照、治理覆盖率快照、治理报表、导出记录、交付材料、交付清单与治理待办表
- [X] T010 [P] 在 `backend/internal/domain/enterprise_polish.go` 定义 `PermissionChangeTrail`、`KeyOperationTrace`、`CrossTeamAuthorizationSnapshot`、`GovernanceRiskEvent`、`GovernanceCoverageSnapshot`、`GovernanceReportPackage`、`ExportRecord`、`DeliveryArtifact`、`DeliveryReadinessBundle`、`DeliveryChecklistItem`、`GovernanceActionItem`
- [X] T011 [P] 创建仓储实现 `backend/internal/repository/permission_change_trail_repository.go`、`backend/internal/repository/key_operation_trace_repository.go`、`backend/internal/repository/cross_team_authorization_snapshot_repository.go`、`backend/internal/repository/governance_risk_event_repository.go`、`backend/internal/repository/governance_coverage_snapshot_repository.go`、`backend/internal/repository/governance_report_package_repository.go`、`backend/internal/repository/export_record_repository.go`、`backend/internal/repository/delivery_artifact_repository.go`、`backend/internal/repository/delivery_readiness_bundle_repository.go`、`backend/internal/repository/delivery_checklist_item_repository.go`、`backend/internal/repository/governance_action_item_repository.go`
- [X] T012 [P] 在 `backend/internal/integration/enterprise/audit_provider.go`、`backend/internal/integration/enterprise/report_builder.go`、`backend/internal/integration/enterprise/delivery_catalog.go` 定义深度审计聚合、报表模板编排和交付包目录接入抽象
- [X] T013 在 `backend/internal/service/enterprise/service.go`、`backend/internal/service/enterprise/scope_service.go`、`backend/internal/service/auth/scope_authorizer.go` 建立企业治理范围过滤、导出范围控制和交付就绪入口
- [X] T014 [P] 在 `backend/internal/repository/redis.go`、`backend/internal/service/enterprise/report_cache.go`、`backend/internal/service/enterprise/export_coordinator.go`、`backend/internal/service/enterprise/trend_cache.go` 建立报表生成缓存、导出协调、趋势快取与幂等键
- [X] T015 [P] 在 `backend/internal/service/auth/permission_service.go`、`backend/internal/api/middleware/authorization.go` 增加 012 权限语义 `enterprise:read`、`enterprise:manage-audit`、`enterprise:manage-reports`、`enterprise:manage-delivery`
- [X] T016 在 `backend/internal/api/router/enterprise_routes.go`、`backend/internal/api/router/router.go` 注册 012 API 路由骨架
- [X] T017 [P] 在 `frontend/src/services/api/types.ts`、`frontend/src/services/api/client.ts`、`frontend/src/app/queryClient.ts` 增加 012 共享类型、查询 key 和错误归一化
- [X] T018 在 `backend/internal/service/audit/event_writer.go`、`backend/internal/service/audit/service.go` 预置 `enterprise.*` 审计动作类型与查询维度映射

**Checkpoint**: 基础能力完成，可开始用户故事实现

---

## Phase 3: User Story 1 - 建立深度权限审计与风险追踪视图 (Priority: P1) 🎯 MVP

**Goal**: 提供权限变更链路、关键操作追踪、跨团队授权分布与高风险访问视图。  
**Independent Test**: 准备权限变更、关键操作和跨团队授权样本后，可独立查看链路还原、风险分类、责任归属和长期趋势。

### Tests for User Story 1

- [X] T019 [P] [US1] 编写后端契约测试 `backend/tests/contract/enterprise_permission_trails_contract_test.go`、`backend/tests/contract/enterprise_key_operations_contract_test.go`
- [X] T020 [P] [US1] 编写后端集成测试 `backend/tests/integration/enterprise_permission_audit_flow_test.go`、`backend/tests/integration/enterprise_cross_team_authorization_flow_test.go`
- [X] T021 [P] [US1] 编写前端 Vitest 页面测试 `frontend/src/features/enterprise-polish/pages/PermissionAuditPage.test.tsx`、`frontend/src/features/enterprise-polish/pages/RiskTrackingPage.test.tsx`

### Implementation for User Story 1

- [X] T022 [P] [US1] 实现权限变更链路与关键操作服务 `backend/internal/service/enterprise/permission_change_trail_service.go`、`backend/internal/service/enterprise/key_operation_trace_service.go`
- [X] T023 [P] [US1] 实现跨团队授权分布与高风险访问服务 `backend/internal/service/enterprise/cross_team_authorization_service.go`、`backend/internal/service/enterprise/governance_risk_service.go`
- [X] T024 [US1] 在 `backend/internal/api/handler/enterprise_handler.go`、`backend/internal/api/router/enterprise_routes.go` 落地 `/enterprise/audit/permission-trails`、`/enterprise/audit/key-operations`
- [X] T025 [US1] 在 `backend/internal/service/enterprise/scope_service.go`、`backend/internal/api/middleware/authorization.go`、`backend/internal/service/auth/scope_authorizer.go` 落地深度审计与风险视图的范围过滤和敏感信息控制
- [X] T026 [P] [US1] 实现前端服务层 `frontend/src/services/enterprisePolish.ts`，覆盖权限链路、关键操作、跨团队授权与风险追踪接口
- [X] T027 [P] [US1] 实现权限审计页面与筛选器 `frontend/src/features/enterprise-polish/pages/PermissionAuditPage.tsx`、`frontend/src/features/enterprise-polish/components/PermissionTrailTable.tsx`
- [X] T028 [P] [US1] 实现风险追踪页面与趋势面板 `frontend/src/features/enterprise-polish/pages/RiskTrackingPage.tsx`、`frontend/src/features/enterprise-polish/components/RiskTrendChart.tsx`
- [X] T029 [US1] 在 `frontend/src/app/router.tsx`、`frontend/src/app/AuthorizedMenu.tsx` 打通企业治理审计导航、跨团队授权分布视图和高风险访问告警入口

**Checkpoint**: US1 完整可测，可作为 012 MVP 交付

---

## Phase 4: User Story 2 - 生成治理报表与标准化导出材料 (Priority: P1)

**Goal**: 提供管理汇报、审计复核和客户交付三类治理报表及导出留痕闭环。  
**Independent Test**: 准备审计、覆盖率和趋势数据后，可独立生成标准化报表、导出材料并查看导出审计记录。

### Tests for User Story 2

- [X] T030 [P] [US2] 编写后端契约测试 `backend/tests/contract/enterprise_governance_reports_contract_test.go`、`backend/tests/contract/enterprise_export_records_contract_test.go`
- [X] T031 [P] [US2] 编写后端集成测试 `backend/tests/integration/enterprise_report_generation_flow_test.go`、`backend/tests/integration/enterprise_export_audit_flow_test.go`
- [X] T032 [P] [US2] 编写前端 Vitest 页面测试 `frontend/src/features/enterprise-polish/pages/GovernanceReportsPage.test.tsx`、`frontend/src/features/audit/pages/EnterpriseAuditPage.test.tsx`

### Implementation for User Story 2

- [X] T033 [P] [US2] 实现治理覆盖率与统一待办服务 `backend/internal/service/enterprise/governance_coverage_service.go`、`backend/internal/service/enterprise/governance_action_item_service.go`
- [X] T034 [P] [US2] 实现治理报表与导出协调服务 `backend/internal/service/enterprise/governance_report_service.go`、`backend/internal/service/enterprise/export_record_service.go`
- [X] T035 [US2] 在 `backend/internal/api/handler/enterprise_handler.go`、`backend/internal/api/router/enterprise_routes.go` 落地 `/enterprise/governance/coverage`、`/enterprise/governance/action-items`、`/enterprise/reports`、`/enterprise/reports/{reportId}/exports`
- [X] T036 [US2] 在 `backend/internal/service/audit/event_writer.go`、`backend/internal/service/audit/service.go` 打通报表生成、导出执行、可见范围裁剪和交付留痕动作的审计写入与查询聚合
- [X] T037 [P] [US2] 扩展前端服务与 hooks `frontend/src/services/enterprisePolish.ts`、`frontend/src/features/enterprise-polish/hooks/useReportActions.ts`
- [X] T038 [P] [US2] 实现治理报表页面与生成表单 `frontend/src/features/enterprise-polish/pages/GovernanceReportsPage.tsx`、`frontend/src/features/enterprise-polish/components/ReportBuilderDrawer.tsx`
- [X] T039 [P] [US2] 实现导出记录与统一待办页面 `frontend/src/features/enterprise-polish/pages/ExportCenterPage.tsx`、`frontend/src/features/enterprise-polish/components/ActionItemList.tsx`
- [X] T040 [US2] 在 `frontend/src/features/audit/pages/EnterpriseAuditPage.tsx`、`frontend/src/app/router.tsx` 落地企业治理审计查询、导出结果追踪和管理汇报快捷入口

**Checkpoint**: US2 可独立验证治理报表与导出闭环

---

## Phase 5: User Story 3 - 形成可复制的产品化交付包 (Priority: P2)

**Goal**: 提供交付材料目录、交付就绪包和交付检查清单，形成可复制的产品化交付收尾能力。  
**Independent Test**: 基于同一套交付材料，可独立查看材料目录、生成交付包并完成交付检查清单核验。

### Tests for User Story 3

- [X] T041 [P] [US3] 编写后端契约测试 `backend/tests/contract/enterprise_delivery_artifacts_contract_test.go`、`backend/tests/contract/enterprise_delivery_bundles_contract_test.go`
- [X] T042 [P] [US3] 编写后端集成测试 `backend/tests/integration/enterprise_delivery_bundle_flow_test.go`、`backend/tests/integration/enterprise_delivery_checklist_flow_test.go`
- [X] T043 [P] [US3] 编写前端 Vitest 页面测试 `frontend/src/features/enterprise-polish/pages/DeliveryArtifactsPage.test.tsx`、`frontend/src/features/enterprise-polish/pages/DeliveryReadinessPage.test.tsx`

### Implementation for User Story 3

- [X] T044 [P] [US3] 实现交付材料目录与交付包服务 `backend/internal/service/enterprise/delivery_artifact_service.go`、`backend/internal/service/enterprise/delivery_readiness_bundle_service.go`
- [X] T045 [P] [US3] 实现交付检查清单与适用边界服务 `backend/internal/service/enterprise/delivery_checklist_service.go`、`backend/internal/service/enterprise/delivery_scope_service.go`
- [X] T046 [US3] 在 `backend/internal/api/handler/enterprise_handler.go`、`backend/internal/api/router/enterprise_routes.go` 落地 `/enterprise/delivery/artifacts`、`/enterprise/delivery/bundles`、`/enterprise/delivery/bundles/{bundleId}/checklists`
- [X] T047 [US3] 在 `backend/internal/api/handler/audit_handler.go`、`backend/internal/api/router/audit_routes.go` 聚合并暴露 `/audit/enterprise/events` 查询链路
- [X] T048 [P] [US3] 扩展前端服务与 hooks `frontend/src/services/enterprisePolish.ts`、`frontend/src/features/enterprise-polish/hooks/useDeliveryBundleActions.ts`
- [X] T049 [P] [US3] 实现交付材料目录页面 `frontend/src/features/enterprise-polish/pages/DeliveryArtifactsPage.tsx`、`frontend/src/features/enterprise-polish/components/DeliveryArtifactCatalog.tsx`
- [X] T050 [P] [US3] 实现交付就绪与检查清单页面 `frontend/src/features/enterprise-polish/pages/DeliveryReadinessPage.tsx`、`frontend/src/features/enterprise-polish/components/DeliveryChecklistBoard.tsx`
- [X] T051 [P] [US3] 实现交付包适用边界和版本差异提示组件 `frontend/src/features/enterprise-polish/components/DeliveryScopeNotice.tsx`、`frontend/src/features/enterprise-polish/components/ReadinessSummaryCard.tsx`
- [X] T052 [US3] 在 `frontend/src/app/AuthorizedMenu.tsx`、`frontend/src/app/router.tsx`、`frontend/src/features/enterprise-polish/pages/DeliveryReadinessPage.tsx` 落地交付包导航、缺失项提醒和交付完成状态空态

**Checkpoint**: US3 完成后形成深度审计、治理报表与产品化交付闭环

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: 收敛质量、文档、验证证据与 PR 交付材料

- [X] T053 [P] 收敛命名与共享类型，在 `backend/internal/service/enterprise/`、`backend/internal/integration/enterprise/`、`frontend/src/features/enterprise-polish/`、`frontend/src/services/enterprisePolish.ts` 清理重复字段与错误文案
- [X] T054 [P] 刷新配置与启动文档，在 `README.md`、`backend/config/config.example.yaml`、`backend/config/config.dev.yaml`、`frontend/.env.example`、`frontend/.env.development` 补齐 012 说明
- [X] T055 [P] 记录验证基线到 `artifacts/012-enterprise-polish/verification.md`、`artifacts/012-enterprise-polish/quickstart-validation.md`、`artifacts/012-enterprise-polish/repro-enterprise-polish-smoke.sh`
- [X] T056 在 `artifacts/012-enterprise-polish/pr-summary.md`、`artifacts/012-enterprise-polish/pr-readiness.md` 准备中文 PR 摘要、备份证据、测试说明、风险清单与用户合并确认项

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 0 (Governance Gates)**: 无依赖，必须先完成
- **Phase 1 (Setup)**: 依赖 Phase 0
- **Phase 2 (Foundational)**: 依赖 Phase 1，阻塞所有用户故事
- **Phase 3/4/5 (User Stories)**: 均依赖 Phase 2；US1 作为 MVP 优先，US2 在 US1 主干稳定后推进，US3 在 US1/US2 的治理语义稳定后推进
- **Phase 6 (Polish)**: 依赖已完成的用户故事范围
- **Release / Merge**: 依赖远程推送、PR 更新、评审完成与用户明确同意

### User Story Dependencies

- **US1 (P1)**: 无用户故事前置依赖，Foundational 完成后可立即开始
- **US2 (P1)**: 依赖 US1 产出的深度审计、风险分类与授权分布治理语义，才能形成标准化报表与导出闭环
- **US3 (P2)**: 依赖 Foundational 的材料目录、报表和待办模型，以及 US2 的导出与可见范围语义

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
Task: "T019 [US1] backend/tests/contract/enterprise_permission_trails_contract_test.go"
Task: "T020 [US1] backend/tests/integration/enterprise_permission_audit_flow_test.go"
Task: "T021 [US1] frontend/src/features/enterprise-polish/pages/PermissionAuditPage.test.tsx"

# 并行实现任务
Task: "T022 [US1] backend/internal/service/enterprise/permission_change_trail_service.go"
Task: "T023 [US1] backend/internal/service/enterprise/governance_risk_service.go"
Task: "T026 [US1] frontend/src/services/enterprisePolish.ts"
```

---

## Implementation Strategy

### MVP First (US1 Only)

1. 完成 Phase 0-2
2. 完成 US1（Phase 3）
3. 按 Independent Test 验证 US1
4. 产出阶段性演示或 PR 更新

### Incremental Delivery

1. 先交付 US1：深度权限审计、关键操作追踪、跨团队授权分布与高风险访问视图
2. 再交付 US2：治理覆盖率、统一待办、标准化治理报表与导出留痕
3. 最后交付 US3：交付材料目录、交付就绪包与交付检查清单
4. 最终执行 Phase 6 文档与验证收尾

### Notes

- `[P]` 任务代表可并行，但仍需满足前置依赖
- 每个用户故事都可独立验收
- 所有提交说明与 PR 摘要必须为中文
- 未获用户明确同意前禁止合并
