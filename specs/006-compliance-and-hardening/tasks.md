# Tasks: 多集群 Kubernetes 合规与加固中心

**Input**: Design documents from `/specs/006-compliance-and-hardening/`  
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/openapi.yaml, quickstart.md

**Tests**: 本特性包含明确的独立验收标准、权限隔离与审计闭环要求，任务清单包含后端契约/集成测试与前端页面测试任务。  
**Organization**: 任务按用户故事分组，保证每个故事可独立实现、独立验证。

**Constitutional Gates**: 必须满足功能分支、数据库备份证据、国内依赖源配置、中文 PR、远程推送、用户同意后合并。

## Format: `[ID] [P?] [Story] Description`

- `[P]`: 可并行执行（不同文件、无前置依赖）
- `[Story]`: 任务归属用户故事（US1/US2/US3）
- 每条任务必须包含明确文件路径

## Phase 0: Governance Gates

**Purpose**: 完成宪章门槛与实施前证据准备

- [X] T001 在 `artifacts/006-compliance-and-hardening/branch-check.txt` 记录当前分支、禁止在 `main/master` 开发、`005` 已合并完成以及“未获用户同意不得合并”门槛
- [X] T002 执行 006 开发前数据库备份并在 `artifacts/006-compliance-and-hardening/backup-manifest.txt` 记录命令、时间戳、产物路径与恢复抽样验证结果
- [X] T003 在 `artifacts/006-compliance-and-hardening/mirror-and-remote-check.txt` 记录 `GOPROXY=https://goproxy.cn,direct`、`https://registry.npmmirror.com`、扫描器/规则包镜像来源与 `git@github.com:baihua19941101/kbManage.git` PR 流程

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: 建立 006 模块骨架、配置入口和导航占位

- [X] T004 创建后端合规模块骨架 `backend/internal/service/compliance/`、`backend/internal/api/handler/compliance_handler.go`、`backend/internal/api/router/compliance_routes.go`
- [X] T005 [P] 创建合规集成层目录与占位 `backend/internal/integration/compliance/`、`backend/internal/integration/compliance/scanner/provider.go`、`backend/internal/integration/compliance/baseline/provider.go`
- [X] T006 [P] 创建前端模块骨架 `frontend/src/features/compliance-hardening/`、`frontend/src/services/compliance.ts` 与路由占位到 `frontend/src/app/router.tsx`
- [X] T007 在 `backend/config/config.example.yaml`、`backend/config/config.dev.yaml`、`frontend/.env.example`、`frontend/.env.development`、`README.md` 增加 `compliance.*` 配置说明
- [X] T008 在 `frontend/src/app/AuthorizedMenu.tsx`、`frontend/src/app/ProtectedRoute.tsx`、`frontend/src/features/auth/store.ts` 接入合规与加固入口门禁占位

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: 完成所有用户故事共享且阻塞性的基础能力

**⚠️ CRITICAL**: US1/US2/US3 必须在本阶段完成后才可开始

- [X] T009 新增 006 数据库迁移 `backend/migrations/0008_compliance_hardening_core.sql`，落库基线标准、扫描配置、扫描执行、失败项、证据、整改任务、例外、复检、趋势快照和归档导出任务表
- [X] T010 [P] 在 `backend/internal/domain/compliance.go` 定义 `ComplianceBaseline`、`ScanProfile`、`ScanExecution`、`ComplianceFinding`、`EvidenceRecord`、`RemediationTask`、`ComplianceExceptionRequest`、`RecheckTask`、`ComplianceTrendSnapshot`、`ArchiveExportTask`
- [X] T011 [P] 创建仓储实现 `backend/internal/repository/compliance_baseline_repository.go`、`backend/internal/repository/compliance_scan_profile_repository.go`、`backend/internal/repository/compliance_scan_execution_repository.go`、`backend/internal/repository/compliance_finding_repository.go`、`backend/internal/repository/compliance_remediation_repository.go`、`backend/internal/repository/compliance_exception_repository.go`、`backend/internal/repository/compliance_recheck_repository.go`、`backend/internal/repository/compliance_trend_repository.go`、`backend/internal/repository/compliance_export_repository.go`
- [X] T012 [P] 在 `backend/internal/integration/compliance/scanner/provider.go`、`backend/internal/integration/compliance/baseline/provider.go` 定义扫描执行抽象、基线包抽象和错误归一化模型
- [X] T013 在 `backend/internal/service/compliance/service.go`、`backend/internal/service/compliance/scope_service.go`、`backend/internal/service/auth/scope_authorizer.go` 建立统一合规授权与范围过滤入口
- [X] T014 [P] 在 `backend/internal/repository/redis.go`、`backend/internal/service/compliance/progress_cache.go`、`backend/internal/service/compliance/export_cache.go`、`backend/internal/service/compliance/schedule_cache.go` 建立扫描进度、导出状态和调度协调缓存
- [X] T015 [P] 在 `backend/internal/service/auth/permission_service.go`、`backend/internal/api/middleware/authorization.go` 增加 006 权限语义 `compliance:read`、`compliance:manage-baseline`、`compliance:execute-scan`、`compliance:manage-remediation`、`compliance:review-exception`、`compliance:export-archive`
- [X] T016 在 `backend/internal/api/router/compliance_routes.go`、`backend/internal/api/router/router.go` 注册 006 API 路由骨架
- [X] T017 [P] 在 `frontend/src/services/api/types.ts`、`frontend/src/services/api/client.ts`、`frontend/src/app/queryClient.ts` 增加 006 共享类型、查询 key 和错误归一化
- [X] T018 在 `backend/internal/service/audit/event_writer.go`、`backend/internal/service/audit/service.go` 预置 `compliance.*` 审计动作类型与查询维度映射

**Checkpoint**: 基础能力完成，可开始用户故事实现

---

## Phase 3: User Story 1 - 统一基线选择与多范围合规扫描 (Priority: P1) 🎯 MVP

**Goal**: 提供基线标准管理、扫描配置、按需/计划性扫描执行以及失败项与证据查看能力。  
**Independent Test**: 创建一条基线标准和一个节点/命名空间范围的扫描配置，触发一次按需扫描后，可查看得分、失败项、证据详情与基线快照。

### Tests for User Story 1

- [X] T019 [P] [US1] 编写后端契约测试 `backend/tests/contract/compliance_baselines_contract_test.go`、`backend/tests/contract/compliance_scan_profiles_contract_test.go`、`backend/tests/contract/compliance_scans_contract_test.go`、`backend/tests/contract/compliance_findings_contract_test.go`
- [X] T020 [P] [US1] 编写后端集成测试 `backend/tests/integration/compliance_scan_execution_test.go`、`backend/tests/integration/compliance_scheduled_scan_test.go`、`backend/tests/integration/compliance_finding_snapshot_test.go`
- [X] T021 [P] [US1] 编写前端 Vitest 页面测试 `frontend/src/features/compliance-hardening/pages/ComplianceBaselinePage.test.tsx`、`frontend/src/features/compliance-hardening/pages/ScanCenterPage.test.tsx`、`frontend/src/features/compliance-hardening/pages/FindingDetailPage.test.tsx`

### Implementation for User Story 1

- [X] T022 [P] [US1] 实现基线标准服务 `backend/internal/service/compliance/baseline_service.go`，覆盖创建、更新、启停和版本快照读取
- [X] T023 [P] [US1] 实现扫描配置服务 `backend/internal/service/compliance/scan_profile_service.go`，覆盖集群/节点/命名空间/关键资源范围配置与计划性调度字段校验
- [X] T024 [P] [US1] 实现扫描执行编排 `backend/internal/service/compliance/scan_execution_service.go`、`backend/internal/integration/compliance/scanner/provider.go`、`backend/internal/worker/compliance_scan_worker.go`，覆盖按需执行、基线快照和部分成功归一化
- [X] T025 [US1] 实现失败项与证据查询服务 `backend/internal/service/compliance/finding_service.go`、`backend/internal/service/compliance/evidence_service.go`
- [X] T026 [US1] 在 `backend/internal/api/handler/compliance_handler.go`、`backend/internal/api/router/compliance_routes.go` 落地 `/compliance/baselines`、`/compliance/scan-profiles`、`/compliance/scan-profiles/{profileId}/execute`、`/compliance/scans`、`/compliance/findings`
- [X] T027 [US1] 在 `backend/internal/service/compliance/scope_service.go`、`backend/internal/api/middleware/authorization.go`、`backend/internal/service/auth/scope_authorizer.go` 落地读取路径的范围过滤与敏感证据访问控制
- [X] T028 [P] [US1] 实现前端服务层 `frontend/src/services/compliance.ts`，覆盖基线、扫描配置、扫描执行、失败项和证据查询接口
- [X] T029 [P] [US1] 实现基线管理与扫描中心页面 `frontend/src/features/compliance-hardening/pages/ComplianceBaselinePage.tsx`、`frontend/src/features/compliance-hardening/pages/ScanCenterPage.tsx`、`frontend/src/features/compliance-hardening/components/BaselineFormDrawer.tsx`、`frontend/src/features/compliance-hardening/components/ScanProfileDrawer.tsx`
- [X] T030 [P] [US1] 实现失败项详情与证据视图 `frontend/src/features/compliance-hardening/pages/FindingDetailPage.tsx`、`frontend/src/features/compliance-hardening/components/FindingTable.tsx`、`frontend/src/features/compliance-hardening/components/EvidenceDrawer.tsx`
- [X] T031 [US1] 在 `frontend/src/app/router.tsx`、`frontend/src/features/clusters/pages/ClusterOverviewPage.tsx`、`frontend/src/features/resources/components/ResourceDetailDrawer.tsx` 打通从集群/资源进入合规上下文的导航入口

**Checkpoint**: US1 完整可测，可作为 006 MVP 交付

---

## Phase 4: User Story 2 - 不合规项整改、例外与复检闭环 (Priority: P1)

**Goal**: 提供整改任务、例外审批和复检闭环能力。  
**Independent Test**: 对一次扫描中的失败项创建整改任务、提交并审批例外、发起复检后，可追踪每个对象从发现到关闭的完整状态链路。

### Tests for User Story 2

- [X] T032 [P] [US2] 编写后端契约测试 `backend/tests/contract/compliance_remediation_contract_test.go`、`backend/tests/contract/compliance_exceptions_contract_test.go`、`backend/tests/contract/compliance_rechecks_contract_test.go`
- [X] T033 [P] [US2] 编写后端集成测试 `backend/tests/integration/compliance_remediation_lifecycle_test.go`、`backend/tests/integration/compliance_exception_expiry_test.go`、`backend/tests/integration/compliance_recheck_lifecycle_test.go`
- [X] T034 [P] [US2] 编写前端 Vitest 页面测试 `frontend/src/features/compliance-hardening/pages/RemediationQueuePage.test.tsx`、`frontend/src/features/compliance-hardening/pages/ComplianceExceptionPage.test.tsx`、`frontend/src/features/compliance-hardening/pages/RecheckCenterPage.test.tsx`

### Implementation for User Story 2

- [X] T035 [P] [US2] 实现整改任务服务 `backend/internal/service/compliance/remediation_service.go`，覆盖任务创建、状态流转、逾期与关闭结论
- [X] T036 [P] [US2] 实现例外申请与审批服务 `backend/internal/service/compliance/exception_service.go`、`backend/internal/worker/compliance_exception_expiry_worker.go`，覆盖申请、审批、拒绝、撤销、到期失效
- [X] T037 [P] [US2] 实现复检服务 `backend/internal/service/compliance/recheck_service.go`、`backend/internal/worker/compliance_recheck_worker.go`，覆盖复检触发、结果回写和失败项状态恢复
- [X] T038 [US2] 在 `backend/internal/api/handler/compliance_handler.go`、`backend/internal/api/router/compliance_routes.go` 落地 `/compliance/remediation-tasks`、`/compliance/exceptions`、`/compliance/exceptions/{exceptionId}/review`、`/compliance/rechecks`
- [X] T039 [US2] 在 `backend/internal/service/audit/event_writer.go`、`backend/internal/service/audit/service.go` 打通整改、例外和复检的审计写入与查询聚合
- [X] T040 [P] [US2] 扩展前端服务与 hooks `frontend/src/services/compliance.ts`、`frontend/src/features/compliance-hardening/hooks/useComplianceAction.ts`
- [X] T041 [P] [US2] 实现整改任务工作台 `frontend/src/features/compliance-hardening/pages/RemediationQueuePage.tsx`、`frontend/src/features/compliance-hardening/components/RemediationTaskTable.tsx`、`frontend/src/features/compliance-hardening/components/RemediationTaskDrawer.tsx`
- [X] T042 [P] [US2] 实现例外审批与复检页面 `frontend/src/features/compliance-hardening/pages/ComplianceExceptionPage.tsx`、`frontend/src/features/compliance-hardening/components/ComplianceExceptionReviewDrawer.tsx`、`frontend/src/features/compliance-hardening/pages/RecheckCenterPage.tsx`、`frontend/src/features/compliance-hardening/components/RecheckTaskTable.tsx`
- [X] T043 [US2] 在 `frontend/src/app/ProtectedRoute.tsx`、`frontend/src/features/compliance-hardening/pages/RemediationQueuePage.tsx`、`frontend/src/features/compliance-hardening/pages/ComplianceExceptionPage.tsx` 落地动作级权限门控、空态和权限回收后的状态处理

**Checkpoint**: US2 可独立验证整改、例外和复检闭环

---

## Phase 5: User Story 3 - 合规覆盖率、趋势复盘与审计汇报 (Priority: P2)

**Goal**: 提供覆盖率总览、趋势比较、归档导出和合规审计能力。  
**Independent Test**: 在累计多次扫描与治理动作后，可按集群/团队查看覆盖率与趋势，生成归档导出任务，并检索完整审计记录。

### Tests for User Story 3

- [X] T044 [P] [US3] 编写后端契约测试 `backend/tests/contract/compliance_overview_contract_test.go`、`backend/tests/contract/compliance_trends_contract_test.go`、`backend/tests/contract/compliance_archive_exports_contract_test.go`、`backend/tests/contract/compliance_audit_contract_test.go`
- [X] T045 [P] [US3] 编写后端集成测试 `backend/tests/integration/compliance_reporting_test.go`、`backend/tests/integration/compliance_archive_export_test.go`、`backend/tests/integration/compliance_scope_authorization_test.go`
- [X] T046 [P] [US3] 编写前端 Vitest 页面测试 `frontend/src/features/compliance-hardening/pages/ComplianceOverviewPage.test.tsx`、`frontend/src/features/compliance-hardening/pages/ComplianceTrendPage.test.tsx`、`frontend/src/features/compliance-hardening/pages/ComplianceArchivePage.test.tsx`、`frontend/src/features/audit/pages/ComplianceAuditPage.test.tsx`

### Implementation for User Story 3

- [X] T047 [P] [US3] 实现覆盖率与趋势服务 `backend/internal/service/compliance/overview_service.go`、`backend/internal/service/compliance/trend_service.go`、`backend/internal/worker/compliance_trend_snapshot_worker.go`
- [X] T048 [P] [US3] 实现归档导出服务 `backend/internal/service/compliance/archive_export_service.go`、`backend/internal/worker/compliance_export_worker.go`、`backend/internal/repository/compliance_export_repository.go`
- [X] T049 [US3] 在 `backend/internal/service/audit/service.go`、`backend/internal/service/audit/event_writer.go`、`backend/internal/api/handler/audit_handler.go` 聚合并暴露 `/audit/compliance/events` 查询链路
- [X] T050 [US3] 在 `backend/internal/api/handler/compliance_handler.go`、`backend/internal/api/router/compliance_routes.go` 落地 `/compliance/overview`、`/compliance/trends`、`/compliance/archive-exports` 和 `/audit/compliance/events`
- [X] T051 [P] [US3] 扩展前端服务与 hooks `frontend/src/services/compliance.ts`、`frontend/src/features/compliance-hardening/hooks/useArchiveExport.ts`
- [X] T052 [P] [US3] 实现覆盖率与趋势页面 `frontend/src/features/compliance-hardening/pages/ComplianceOverviewPage.tsx`、`frontend/src/features/compliance-hardening/pages/ComplianceTrendPage.tsx`、`frontend/src/features/compliance-hardening/components/ComplianceOverviewCards.tsx`、`frontend/src/features/compliance-hardening/components/ComplianceTrendChart.tsx`
- [X] T053 [P] [US3] 实现归档导出与审计页面 `frontend/src/features/compliance-hardening/pages/ComplianceArchivePage.tsx`、`frontend/src/features/compliance-hardening/components/ArchiveExportDrawer.tsx`、`frontend/src/features/audit/pages/ComplianceAuditPage.tsx`
- [X] T054 [US3] 在 `frontend/src/app/AuthorizedMenu.tsx`、`frontend/src/app/router.tsx`、`frontend/src/features/compliance-hardening/pages/ComplianceArchivePage.tsx` 落地报表/导出权限门控、筛选持久化和未授权空态

**Checkpoint**: US3 完成后形成覆盖率、趋势、归档与审计闭环

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: 收敛质量、文档、验证证据与 PR 交付材料

- [X] T055 [P] 收敛命名与共享类型，在 `backend/internal/service/compliance/`、`backend/internal/integration/compliance/`、`frontend/src/features/compliance-hardening/`、`frontend/src/services/compliance.ts` 清理重复字段与错误文案
- [X] T056 [P] 刷新配置与启动文档，在 `README.md`、`backend/config/config.example.yaml`、`backend/config/config.dev.yaml`、`frontend/.env.example`、`frontend/.env.development` 补齐 006 说明
- [X] T057 [P] 记录验证基线到 `artifacts/006-compliance-and-hardening/verification.md`、`artifacts/006-compliance-and-hardening/quickstart-validation.md`、`artifacts/006-compliance-and-hardening/repro-compliance-smoke.sh`
- [X] T058 在 `artifacts/006-compliance-and-hardening/pr-summary.md`、`artifacts/006-compliance-and-hardening/pr-readiness.md` 准备中文 PR 摘要、备份证据、测试说明、风险清单与用户合并确认项

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 0 (Governance Gates)**: 无依赖，必须先完成
- **Phase 1 (Setup)**: 依赖 Phase 0
- **Phase 2 (Foundational)**: 依赖 Phase 1，阻塞所有用户故事
- **Phase 3/4/5 (User Stories)**: 均依赖 Phase 2；US1 作为 MVP 优先，US2 和 US3 在 US1 主干稳定后推进
- **Phase 6 (Polish)**: 依赖已完成的用户故事范围
- **Release / Merge**: 依赖远程推送、PR 更新、评审完成与用户明确同意

### User Story Dependencies

- **US1 (P1)**: 无用户故事前置依赖，Foundational 完成后可立即开始
- **US2 (P1)**: 依赖 US1 产出的扫描结果与失败项对象，才能形成整改、例外和复检闭环
- **US3 (P2)**: 依赖 US1 的扫描与趋势数据；若需要完整整改/例外统计，则同时集成 US2 的治理状态

### Parallel Opportunities

- **Phase 1**: T005/T006 可并行
- **Phase 2**: T010/T011/T012/T014/T015/T017 可并行
- **US1**: T019/T020/T021 并行，T022/T023/T024 并行，T028/T029/T030 并行
- **US2**: T032/T033/T034 并行，T035/T036/T037 并行，T040/T041/T042 并行
- **US3**: T044/T045/T046 并行，T047/T048 并行，T051/T052/T053 并行

---

## Parallel Example: User Story 1

```bash
# 并行测试任务
Task: "T019 [US1] backend/tests/contract/compliance_baselines_contract_test.go"
Task: "T020 [US1] backend/tests/integration/compliance_scan_execution_test.go"
Task: "T021 [US1] frontend/src/features/compliance-hardening/pages/ComplianceBaselinePage.test.tsx"

# 并行实现任务
Task: "T022 [US1] backend/internal/service/compliance/baseline_service.go"
Task: "T023 [US1] backend/internal/service/compliance/scan_profile_service.go"
Task: "T028 [US1] frontend/src/services/compliance.ts"
```

---

## Implementation Strategy

### MVP First (US1 Only)

1. 完成 Phase 0-2
2. 完成 US1（Phase 3）
3. 按 Independent Test 验证 US1
4. 产出阶段性演示或 PR 更新

### Incremental Delivery

1. 先交付 US1：统一基线选择、多范围扫描与失败项详情
2. 再交付 US2：整改任务、例外审批与复检闭环
3. 最后交付 US3：覆盖率、趋势、归档导出与审计汇报
4. 最终执行 Phase 6 文档与验证收尾

### Notes

- `[P]` 任务代表可并行，但仍需满足前置依赖
- 每个用户故事都可独立验收
- 所有提交说明与 PR 摘要必须为中文
- 未获用户明确同意前禁止合并
