# Tasks: 多集群 Kubernetes 可观测中心

**Input**: Design documents from `/specs/002-observability-center/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/, quickstart.md

**Tests**: 本特性明确要求测试任务，必须覆盖后端 contract/integration 测试以及前端 Vitest 页面测试。

**Organization**: 任务按用户故事组织，确保 US1、US2、US3 可以独立实现与独立验收。

**Constitutional Gates**: 功能分支校验、数据库备份证据、国内源配置、中文 PR、远程推送和明确的用户合并授权必须在实现与交付时满足。

## Format: `[ID] [P?] [Story] Description`

- **[P]**: 可并行执行（不同文件、无未完成依赖）
- **[Story]**: 对应的用户故事标签（US1、US2、US3）
- 描述中包含明确文件路径，便于直接执行

## Phase 0: Governance Gates

**Purpose**: 在真正开始实现 002 之前完成宪章要求的治理门槛

- [x] T001 在 `artifacts/002-observability-center/branch-check.txt` 记录当前分支、`001` PR 完成状态检查结果和“未经用户批准不得合并”的门槛
- [x] T002 执行 002 本轮数据库备份并在 `artifacts/002-observability-center/backup-manifest.txt` 记录命令、时间戳、产物路径和恢复抽样验证结果
- [x] T003 在 `artifacts/002-observability-center/mirror-and-remote-check.txt` 记录 `GOPROXY=https://goproxy.cn,direct`、`https://registry.npmmirror.com`、观测联调镜像来源和 `git@github.com:baihua19941101/kbManage.git` PR 流程
- [x] T004 更新 `specs/002-observability-center/spec.md`、`specs/002-observability-center/plan.md` 和 `specs/002-observability-center/quickstart.md` 中的执行状态，注明已进入 tasks 阶段且实现前仍受 `001` PR gate 限制

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: 为 002 建立最小的模块骨架、配置入口和依赖准备

- [x] T005 创建后端可观测模块骨架 `backend/internal/service/observability/`、`backend/internal/integration/observability/alerts/`、`backend/internal/integration/observability/logs/`、`backend/internal/integration/observability/metrics/`
- [x] T006 [P] 创建前端可观测模块骨架 `frontend/src/features/observability/`、`frontend/src/services/observability/` 和对应路由占位文件 `frontend/src/app/router.tsx`
- [x] T007 [P] 在 `frontend/package.json`、`frontend/src/main.tsx` 和 `frontend/src/app/App.tsx` 引入 002 所需图表依赖与页面注册，保持 001 主栈不变
- [x] T008 在 `backend/config/config.example.yaml`、`backend/config/config.dev.yaml` 和 `README.md` 增加 `observability.metrics.*`、`observability.logs.*`、`observability.alerts.*`、`observability.cache.*` 的配置说明

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: 搭建 002 的共享数据模型、适配器接口、路由和授权基础，阻塞所有用户故事

**⚠️ CRITICAL**: US1、US2、US3 都必须在本阶段完成后才能开始

- [x] T009 新增 002 数据库迁移 `backend/migrations/0004_observability_core.sql`，落库数据源配置、告警规则、通知目标、静默窗口、告警快照和处理记录表
- [x] T010 [P] 在 `backend/internal/domain/observability.go` 定义 `ObservabilityDataSource`、`AlertRule`、`NotificationTarget`、`SilenceWindow`、`AlertIncidentSnapshot` 和 `AlertHandlingRecord` 领域模型
- [x] T011 [P] 创建仓储实现 `backend/internal/repository/observability_datasource_repository.go`、`backend/internal/repository/alert_rule_repository.go`、`backend/internal/repository/notification_target_repository.go`、`backend/internal/repository/silence_window_repository.go`、`backend/internal/repository/alert_incident_repository.go`
- [x] T012 [P] 在 `backend/internal/integration/observability/logs/provider.go`、`backend/internal/integration/observability/metrics/provider.go`、`backend/internal/integration/observability/alerts/provider.go` 定义 Loki/Prometheus/Alertmanager 兼容适配接口与错误归一化模型
- [x] T013 在 `backend/internal/service/observability/service.go`、`backend/internal/service/observability/scope_service.go` 和 `backend/internal/service/auth/scope_authorizer.go` 建立统一观测授权与范围过滤入口
- [x] T014 [P] 在 `backend/internal/api/handler/observability_handler.go`、`backend/internal/api/router/observability_routes.go` 和 `backend/internal/api/router/router.go` 注册 002 API 路由骨架
- [x] T015 [P] 在 `backend/internal/repository/redis.go`、`backend/internal/service/observability/query_cache.go` 和 `backend/internal/worker/observability_sync_worker.go` 建立查询缓存、短时上下文和告警同步基础设施
- [x] T016 [P] 在 `frontend/src/services/api/types.ts`、`frontend/src/services/api/client.ts` 和 `frontend/src/app/queryClient.ts` 增加 002 的共享类型、错误归一化和查询 key 约定
- [x] T017 在 `frontend/src/app/AuthorizedMenu.tsx`、`frontend/src/app/router.tsx` 和 `frontend/src/app/ProtectedRoute.tsx` 接入可观测功能入口和基础权限门控

**Checkpoint**: 002 的共享模型、配置、适配器接口、路由与授权骨架就绪，用户故事可以开始实现

---

## Phase 3: User Story 1 - 统一可观测入口与问题定位 (Priority: P1) 🎯 MVP

**Goal**: 提供围绕资源上下文的统一概览、日志查询、事件时间线和指标趋势入口，形成最小可观测问题定位闭环

**Independent Test**: 在至少两个已接入集群中制造一个工作负载异常后，授权用户能够在同一平台内查看关联日志、事件、指标和资源上下文，并据此完成初步定位

### Tests for User Story 1

- [x] T018 [P] [US1] 收紧并新增后端契约测试 `backend/tests/contract/observability_overview_contract_test.go`、`backend/tests/contract/observability_logs_contract_test.go`、`backend/tests/contract/observability_events_contract_test.go`、`backend/tests/contract/observability_metrics_contract_test.go`
- [x] T019 [P] [US1] 编写后端集成测试 `backend/tests/integration/observability_overview_test.go` 和 `backend/tests/integration/resource_context_test.go`，覆盖多集群概览、资源上下文联动、数据源降级和错误提示
- [x] T020 [P] [US1] 编写前端 Vitest 页面测试 `frontend/src/features/observability/pages/ObservabilityOverviewPage.test.tsx`、`frontend/src/features/observability/pages/LogExplorerPage.test.tsx`、`frontend/src/features/observability/pages/ResourceContextPage.test.tsx`

### Implementation for User Story 1

- [x] T021 [P] [US1] 实现集群观测配置读写与连通性校验 `backend/internal/service/observability/datasource_service.go`、`backend/internal/api/handler/observability_config_handler.go`、`backend/internal/api/router/cluster_routes.go`
- [x] T022 [P] [US1] 实现 Prometheus 指标查询适配器和指标摘要服务 `backend/internal/integration/observability/metrics/prometheus_provider.go`、`backend/internal/service/observability/metrics_service.go`
- [x] T023 [P] [US1] 实现 Loki 日志查询适配器和日志查询服务 `backend/internal/integration/observability/logs/loki_provider.go`、`backend/internal/service/observability/logs_service.go`
- [x] T024 [P] [US1] 实现 Kubernetes Event 查询与时间线聚合 `backend/internal/service/observability/events_service.go`、`backend/internal/kube/adapter/event_reader.go`
- [x] T025 [US1] 实现资源上下文聚合查询 `backend/internal/service/observability/resource_context_service.go`、`backend/internal/api/handler/observability_handler.go`
- [x] T026 [US1] 在 `backend/internal/api/handler/observability_handler.go` 和 `backend/internal/api/router/observability_routes.go` 落地 `/observability/overview`、`/observability/logs/query`、`/observability/events`、`/observability/metrics/series`、`/observability/resources/context`
- [x] T027 [P] [US1] 创建前端服务层 `frontend/src/services/observability/overview.ts`、`frontend/src/services/observability/logs.ts`、`frontend/src/services/observability/events.ts`、`frontend/src/services/observability/metrics.ts`
- [x] T028 [P] [US1] 实现总览与图表组件 `frontend/src/features/observability/components/OverviewCards.tsx`、`frontend/src/features/observability/components/MetricsTrendChart.tsx`、`frontend/src/features/observability/pages/ObservabilityOverviewPage.tsx`
- [x] T029 [P] [US1] 实现日志检索页和筛选组件 `frontend/src/features/observability/components/LogFilters.tsx`、`frontend/src/features/observability/components/LogTable.tsx`、`frontend/src/features/observability/pages/LogExplorerPage.tsx`
- [x] T030 [P] [US1] 实现事件时间线和资源上下文页 `frontend/src/features/observability/components/EventTimeline.tsx`、`frontend/src/features/observability/components/ResourceContextPanel.tsx`、`frontend/src/features/observability/pages/ResourceContextPage.tsx`
- [x] T031 [US1] 在 `frontend/src/app/router.tsx`、`frontend/src/features/resources/components/ResourceDetailDrawer.tsx` 和 `frontend/src/features/clusters/pages/ClusterOverviewPage.tsx` 打通从资源与集群进入可观测上下文的导航入口

**Checkpoint**: US1 完成后，平台应具备可独立演示的统一观测入口与问题定位 MVP

---

## Phase 4: User Story 2 - 告警治理与值班闭环 (Priority: P1)

**Goal**: 提供告警中心、规则治理、通知目标、静默窗口和处理记录，形成从触发到恢复的值班闭环

**Independent Test**: 创建一条告警规则并触发一次异常后，授权用户能够看到告警生成、确认、静默、恢复和处理记录的完整链路

### Tests for User Story 2

- [x] T032 [P] [US2] 编写后端契约测试 `backend/tests/contract/observability_alerts_contract_test.go`、`backend/tests/contract/observability_alert_rules_contract_test.go`、`backend/tests/contract/observability_notification_targets_contract_test.go`、`backend/tests/contract/observability_silences_contract_test.go`
- [x] T033 [P] [US2] 编写后端集成测试 `backend/tests/integration/alert_center_test.go` 和 `backend/tests/integration/alert_governance_test.go`，覆盖规则生命周期、通知目标、静默窗口、处理记录和恢复状态
- [x] T034 [P] [US2] 编写前端 Vitest 页面测试 `frontend/src/features/observability/pages/AlertCenterPage.test.tsx`、`frontend/src/features/observability/pages/AlertRulePage.test.tsx`、`frontend/src/features/observability/pages/SilenceWindowPage.test.tsx`

### Implementation for User Story 2

- [x] T035 [P] [US2] 实现 Alertmanager 告警查询与静默适配器 `backend/internal/integration/observability/alerts/alertmanager_provider.go`、`backend/internal/service/observability/alert_sync_service.go`
- [x] T036 [P] [US2] 实现规则治理与通知目标服务 `backend/internal/service/observability/alert_rule_service.go`、`backend/internal/service/observability/notification_target_service.go`
- [x] T037 [P] [US2] 实现静默窗口和处理记录服务 `backend/internal/service/observability/silence_service.go`、`backend/internal/service/observability/handling_record_service.go`
- [x] T038 [US2] 在 `backend/internal/api/handler/observability_alert_handler.go`、`backend/internal/api/handler/observability_admin_handler.go` 和 `backend/internal/api/router/observability_routes.go` 落地 `/observability/alerts`、`/observability/alert-rules`、`/observability/notification-targets`、`/observability/silences`
- [x] T039 [US2] 在 `backend/internal/worker/observability_sync_worker.go`、`backend/internal/service/audit/event_writer.go` 和 `backend/internal/service/audit/service.go` 打通告警同步、治理动作审计和运行态快照回写
- [x] T040 [P] [US2] 创建前端服务层 `frontend/src/services/observability/alerts.ts`、`frontend/src/services/observability/alertRules.ts`、`frontend/src/services/observability/notificationTargets.ts`、`frontend/src/services/observability/silences.ts`
- [x] T041 [P] [US2] 实现告警中心页面与处理交互 `frontend/src/features/observability/components/AlertTable.tsx`、`frontend/src/features/observability/components/AlertDetailDrawer.tsx`、`frontend/src/features/observability/pages/AlertCenterPage.tsx`
- [x] T042 [P] [US2] 实现规则治理页面 `frontend/src/features/observability/components/AlertRuleForm.tsx`、`frontend/src/features/observability/pages/AlertRulePage.tsx`
- [x] T043 [P] [US2] 实现通知目标与静默窗口页面 `frontend/src/features/observability/components/NotificationTargetForm.tsx`、`frontend/src/features/observability/components/SilenceWindowForm.tsx`、`frontend/src/features/observability/pages/SilenceWindowPage.tsx`
- [x] T044 [US2] 在 `frontend/src/app/router.tsx` 和 `frontend/src/app/AuthorizedMenu.tsx` 接入告警中心、规则、通知目标和静默窗口导航与状态联动

**Checkpoint**: US2 完成后，平台应具备完整的告警治理与值班闭环

---

## Phase 5: User Story 3 - 权限隔离下的可观测访问 (Priority: P2)

**Goal**: 对日志、事件、指标、告警和治理动作统一执行工作空间/项目级授权校验，确保不同团队只能看到自己范围内的可观测数据

**Independent Test**: 为两个不同工作空间分别授权后，不同用户只能看到各自范围内的日志、事件、指标和告警；权限被回收后访问立即失效

### Tests for User Story 3

- [x] T045 [P] [US3] 编写后端契约测试 `backend/tests/contract/observability_access_control_contract_test.go`，覆盖未授权日志、事件、指标、告警和治理接口访问
- [x] T046 [P] [US3] 编写后端集成测试 `backend/tests/integration/observability_scope_authorization_test.go`，覆盖跨工作空间/项目隔离、权限回收即时生效和只读角色限制
- [x] T047 [P] [US3] 编写前端 Vitest 页面测试 `frontend/src/features/observability/pages/ObservabilityAccessGate.test.tsx` 和 `frontend/src/features/observability/pages/AlertCenterPermissions.test.tsx`

### Implementation for User Story 3

- [x] T048 [P] [US3] 在 `backend/internal/service/auth/scope_authorizer.go`、`backend/internal/service/observability/scope_service.go` 和 `backend/internal/service/cluster/service.go` 实现观测范围到工作空间/项目授权边界的映射校验
- [x] T049 [P] [US3] 在 `backend/internal/api/middleware/authorization.go` 和 `backend/internal/api/router/observability_routes.go` 增加观测接口的统一授权中间件和只读/治理动作区分
- [x] T050 [US3] 在 `backend/internal/service/observability/logs_service.go`、`backend/internal/service/observability/events_service.go`、`backend/internal/service/observability/metrics_service.go`、`backend/internal/service/observability/alert_rule_service.go` 落地后端范围过滤与敏感错误归一化
- [x] T051 [US3] 在 `backend/internal/service/audit/service.go`、`backend/internal/service/audit/event_writer.go` 和 `backend/internal/domain/audit.go` 增加关键观测访问、规则变更、静默和确认动作审计字段
- [x] T052 [P] [US3] 在 `frontend/src/app/AuthorizedMenu.tsx`、`frontend/src/app/ProtectedRoute.tsx` 和 `frontend/src/features/auth/store.ts` 实现可观测模块的导航可见性与动作级权限门控
- [x] T053 [US3] 在 `frontend/src/features/observability/pages/ObservabilityOverviewPage.tsx`、`frontend/src/features/observability/pages/AlertCenterPage.tsx` 和 `frontend/src/features/observability/pages/AlertRulePage.tsx` 实现未授权空态、只读态和权限回收后的错误处理

**Checkpoint**: US3 完成后，002 的所有可观测读取与治理动作都应运行在现有租户隔离模型之下

---

## Final Phase: Polish & Cross-Cutting Concerns

**Purpose**: 收尾文档、验证、性能与交付准备

- [x] T054 [P] 优化跨故事共享类型与组件复用，在 `backend/internal/service/observability/`、`frontend/src/features/observability/components/` 和 `frontend/src/services/observability/` 做代码清理与命名收敛
- [x] T055 [P] 补齐配置与启动文档，在 `README.md`、`backend/config/config.example.yaml`、`backend/config/config.dev.yaml`、`frontend/.env.example` 中刷新 002 相关说明
- [x] T056 [P] 记录验证基线与联调结果到 `artifacts/002-observability-center/verification.md`、`artifacts/002-observability-center/quickstart-validation.md` 和 `artifacts/002-observability-center/repro-observability-smoke.sh`
- [x] T057 在 `artifacts/002-observability-center/pr-summary.md` 和 `artifacts/002-observability-center/pr-readiness.md` 准备中文 PR 摘要、备份证据、测试说明、风险清单和“待用户批准合并”的交付记录

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 0: Governance Gates**: 无依赖，但 002 真正开始编码前必须完成，并确认 `001` PR 流已结束
- **Phase 1: Setup**: 依赖 Phase 0，负责建立模块骨架和配置入口
- **Phase 2: Foundational**: 依赖 Phase 1，阻塞所有用户故事
- **Phase 3: US1**: 依赖 Phase 2，是 002 的 MVP
- **Phase 4: US2**: 依赖 Phase 2，可与 US1 部分并行，但最佳实践是在 US1 的共享查询与配置能力落稳后推进
- **Phase 5: US3**: 依赖 US1 和 US2，因为需要对已存在的观测读取与治理动作统一施加授权
- **Final Phase**: 依赖目标用户故事完成

### User Story Dependencies

- **US1**: 仅依赖共享基础设施，是最小可演示范围
- **US2**: 依赖共享基础设施；可不等待 US1 全部 UI 完成，但会复用数据源配置、适配器和路由基础
- **US3**: 依赖 US1、US2，因为权限隔离需要覆盖日志、事件、指标、告警和治理接口

### Parallel Opportunities

- Phase 1 中，`T006`、`T007`、`T008` 可并行
- Phase 2 中，`T010`、`T011`、`T012`、`T014`、`T016` 可并行
- US1 中，`T018`、`T019`、`T020` 可并行；`T022`、`T023`、`T024` 可并行；`T027`、`T028`、`T029`、`T030` 可并行
- US2 中，`T032`、`T033`、`T034` 可并行；`T035`、`T036`、`T037` 可并行；`T040`、`T041`、`T042`、`T043` 可并行
- US3 中，`T045`、`T046`、`T047` 可并行；`T048`、`T049`、`T052` 可并行
- Final Phase 中，`T054`、`T055`、`T056` 可并行

---

## Parallel Example: User Story 1

```bash
# 并行启动 US1 的测试任务
Task: "T018 [US1] 后端契约测试"
Task: "T019 [US1] 后端集成测试"
Task: "T020 [US1] 前端 Vitest 页面测试"

# 并行实现 US1 的三类查询后端
Task: "T022 [US1] Prometheus 指标适配器与指标服务"
Task: "T023 [US1] Loki 日志适配器与日志服务"
Task: "T024 [US1] Kubernetes Event 查询与时间线聚合"

# 并行实现 US1 的前端视图
Task: "T028 [US1] 可观测总览页"
Task: "T029 [US1] 日志检索页"
Task: "T030 [US1] 事件时间线与资源上下文页"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. 完成 Phase 0: Governance Gates
2. 完成 Phase 1: Setup
3. 完成 Phase 2: Foundational
4. 完成 Phase 3: US1
5. **STOP and VALIDATE**：按 US1 独立验收标准验证多集群概览、日志、事件、指标和资源上下文联动

### Incremental Delivery

1. 完成 Governance + Setup + Foundational，形成 002 共享底座
2. 交付 US1，形成可演示 MVP
3. 交付 US2，补齐告警治理与值班闭环
4. 交付 US3，补齐租户隔离与动作级授权
5. 完成 Final Phase，准备 PR 证据与交付文档

### Notes

- 首期范围严格限定在 Prometheus + Alertmanager + Loki 兼容接入、Kubernetes Event 查询、统一授权校验、审计记录和资源上下文联动
- 终端、批量操作、回滚、GitOps、Helm、策略治理、合规扫描、集群生命周期和灾备任务禁止混入本任务清单
- 如果后续需要扩展日志实时 tail、终端、批量动作等 Day2 运维能力，应进入 `003-workload-day2-ops`
