# Tasks: 平台应用目录与扩展市场

**Input**: Design documents from `/specs/010-platform-marketplace/`  
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/openapi.yaml, quickstart.md

**Tests**: 本特性包含明确的独立验收标准、目录来源治理、模板分发、扩展兼容性和权限边界要求，任务清单包含后端契约/集成测试与前端页面测试任务。  
**Organization**: 任务按用户故事分组，保证每个故事可独立实现、独立验证。

**Constitutional Gates**: 必须满足功能分支、数据库备份证据、国内依赖源配置、中文 PR、远程推送、用户同意后合并。

## Format: `[ID] [P?] [Story] Description`

- `[P]`: 可并行执行（不同文件、无前置依赖）
- `[Story]`: 任务归属用户故事（US1/US2/US3）
- 每条任务必须包含明确文件路径

## Phase 0: Governance Gates

**Purpose**: 完成宪章门槛与实施前证据准备

- [X] T001 在 `artifacts/010-platform-marketplace/branch-check.txt` 记录当前分支、禁止在 `main/master` 开发、`009` 已合并完成以及“未获用户同意不得合并”门槛
- [X] T002 执行 010 开发前数据库备份并在 `artifacts/010-platform-marketplace/backup-manifest.txt` 记录命令、时间戳、产物路径与恢复抽样验证结果
- [X] T003 在 `artifacts/010-platform-marketplace/mirror-and-remote-check.txt` 记录 `GOPROXY=https://goproxy.cn,direct`、`https://registry.npmmirror.com`、目录/扩展来源联调镜像策略与 `git@github.com:baihua19941101/kbManage.git` PR 流程

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: 建立 010 模块骨架、配置入口和导航占位

- [X] T004 创建后端应用目录与扩展市场模块骨架 `backend/internal/service/marketplace/`、`backend/internal/api/handler/marketplace_handler.go`、`backend/internal/api/router/marketplace_routes.go`
- [X] T005 [P] 创建目录与扩展适配层目录 `backend/internal/integration/marketplace/`、`backend/internal/integration/marketplace/catalog_provider.go`、`backend/internal/integration/marketplace/extension_registry.go`
- [X] T006 [P] 创建前端模块骨架 `frontend/src/features/platform-marketplace/`、`frontend/src/services/platformMarketplace.ts` 与路由占位到 `frontend/src/app/router.tsx`
- [X] T007 在 `backend/config/config.example.yaml`、`backend/config/config.dev.yaml`、`frontend/.env.example`、`frontend/.env.development`、`README.md` 增加 `platformMarketplace.*` 配置说明
- [X] T008 在 `frontend/src/app/AuthorizedMenu.tsx`、`frontend/src/app/ProtectedRoute.tsx`、`frontend/src/features/auth/store.ts` 接入应用目录与扩展市场入口门禁占位

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: 完成所有用户故事共享且阻塞性的基础能力

**⚠️ CRITICAL**: US1/US2/US3 必须在本阶段完成后才可开始

- [X] T009 新增 010 数据库迁移 `backend/migrations/0012_platform_marketplace_core.sql`，落库目录来源、模板、模板版本、模板分发、安装记录、扩展包、兼容性结论、扩展生命周期与市场审计表
- [X] T010 [P] 在 `backend/internal/domain/platform_marketplace.go` 定义 `CatalogSource`、`ApplicationTemplate`、`TemplateVersion`、`TemplateReleaseScope`、`InstallationRecord`、`ExtensionPackage`、`CompatibilityStatement`、`ExtensionLifecycleRecord`
- [X] T011 [P] 创建仓储实现 `backend/internal/repository/catalog_source_repository.go`、`backend/internal/repository/application_template_repository.go`、`backend/internal/repository/template_version_repository.go`、`backend/internal/repository/template_release_scope_repository.go`、`backend/internal/repository/installation_record_repository.go`、`backend/internal/repository/extension_package_repository.go`、`backend/internal/repository/compatibility_statement_repository.go`、`backend/internal/repository/extension_lifecycle_repository.go`
- [X] T012 [P] 在 `backend/internal/integration/marketplace/catalog_provider.go`、`backend/internal/integration/marketplace/extension_registry.go` 定义目录来源访问抽象、扩展注册抽象、同步与兼容性结果模型
- [X] T013 在 `backend/internal/service/marketplace/service.go`、`backend/internal/service/marketplace/scope_service.go`、`backend/internal/service/auth/scope_authorizer.go` 建立市场范围过滤、目标范围映射和权限边界入口
- [X] T014 [P] 在 `backend/internal/repository/redis.go`、`backend/internal/service/marketplace/catalog_cache.go`、`backend/internal/service/marketplace/distribution_coordinator.go`、`backend/internal/service/marketplace/compatibility_cache.go` 建立目录缓存、分发协调、兼容性缓存与幂等键
- [X] T015 [P] 在 `backend/internal/service/auth/permission_service.go`、`backend/internal/api/middleware/authorization.go` 增加 010 权限语义 `marketplace:read`、`marketplace:manage-source`、`marketplace:publish-template`、`marketplace:manage-extension`
- [X] T016 在 `backend/internal/api/router/marketplace_routes.go`、`backend/internal/api/router/router.go` 注册 010 API 路由骨架
- [X] T017 [P] 在 `frontend/src/services/api/types.ts`、`frontend/src/services/api/client.ts`、`frontend/src/app/queryClient.ts` 增加 010 共享类型、查询 key 和错误归一化
- [X] T018 在 `backend/internal/service/audit/event_writer.go`、`backend/internal/service/audit/service.go` 预置 `platformmarketplace.*` 审计动作类型与查询维度映射

**Checkpoint**: 基础能力完成，可开始用户故事实现

---

## Phase 3: User Story 1 - 管理应用目录与模板中心 (Priority: P1) 🎯 MVP

**Goal**: 提供目录来源管理、模板中心、模板版本与依赖治理能力。  
**Independent Test**: 新增一个目录来源并导入多个模板版本后，可在模板中心查看分类、版本、依赖关系、参数表单、部署约束和适用范围。

### Tests for User Story 1

- [X] T019 [P] [US1] 编写后端契约测试 `backend/tests/contract/marketplace_catalog_source_contract_test.go`、`backend/tests/contract/marketplace_template_catalog_contract_test.go`
- [X] T020 [P] [US1] 编写后端集成测试 `backend/tests/integration/marketplace_catalog_sync_flow_test.go`、`backend/tests/integration/marketplace_template_version_flow_test.go`
- [X] T021 [P] [US1] 编写前端 Vitest 页面测试 `frontend/src/features/platform-marketplace/pages/CatalogSourcePage.test.tsx`、`frontend/src/features/platform-marketplace/pages/TemplateCatalogPage.test.tsx`

### Implementation for User Story 1

- [X] T022 [P] [US1] 实现目录来源与模板目录服务 `backend/internal/service/marketplace/catalog_source_service.go`、`backend/internal/service/marketplace/template_catalog_service.go`
- [X] T023 [P] [US1] 实现模板版本、依赖与参数摘要服务 `backend/internal/service/marketplace/template_version_service.go`、`backend/internal/service/marketplace/template_dependency_service.go`
- [X] T024 [US1] 在 `backend/internal/api/handler/marketplace_handler.go`、`backend/internal/api/router/marketplace_routes.go` 落地 `/marketplace/catalog-sources`、`/marketplace/catalog-sources/{sourceId}/sync`、`/marketplace/templates`、`/marketplace/templates/{templateId}`
- [X] T025 [US1] 在 `backend/internal/service/marketplace/scope_service.go`、`backend/internal/api/middleware/authorization.go`、`backend/internal/service/auth/scope_authorizer.go` 落地目录来源与模板中心路径的范围过滤
- [X] T026 [P] [US1] 实现前端服务层 `frontend/src/services/platformMarketplace.ts`，覆盖目录来源、模板中心和模板详情接口
- [X] T027 [P] [US1] 实现目录来源页面与表单 `frontend/src/features/platform-marketplace/pages/CatalogSourcePage.tsx`、`frontend/src/features/platform-marketplace/components/CatalogSourceDrawer.tsx`
- [X] T028 [P] [US1] 实现模板中心与模板详情页面 `frontend/src/features/platform-marketplace/pages/TemplateCatalogPage.tsx`、`frontend/src/features/platform-marketplace/pages/TemplateDetailPage.tsx`
- [X] T029 [US1] 在 `frontend/src/app/router.tsx`、`frontend/src/features/platform-marketplace/pages/TemplateCatalogPage.tsx` 打通目录中心导航、模板状态标识和下线空态

**Checkpoint**: US1 完整可测，可作为 010 MVP 交付

---

## Phase 4: User Story 2 - 按范围分发模板并跟踪安装升级 (Priority: P1)

**Goal**: 提供模板发布到工作空间、项目或集群范围的能力，并跟踪安装记录、升级入口和版本变化。  
**Independent Test**: 将模板发布到目标范围后，可查看目标范围内的可用模板、安装记录、升级入口和下线状态。

### Tests for User Story 2

- [X] T030 [P] [US2] 编写后端契约测试 `backend/tests/contract/marketplace_template_release_contract_test.go`、`backend/tests/contract/marketplace_installation_record_contract_test.go`
- [X] T031 [P] [US2] 编写后端集成测试 `backend/tests/integration/marketplace_template_release_flow_test.go`、`backend/tests/integration/marketplace_installation_upgrade_flow_test.go`
- [X] T032 [P] [US2] 编写前端 Vitest 页面测试 `frontend/src/features/platform-marketplace/pages/TemplateDistributionPage.test.tsx`、`frontend/src/features/platform-marketplace/pages/InstallationRecordPage.test.tsx`

### Implementation for User Story 2

- [X] T033 [P] [US2] 实现模板分发与安装记录服务 `backend/internal/service/marketplace/template_release_service.go`、`backend/internal/service/marketplace/installation_record_service.go`
- [X] T034 [P] [US2] 实现升级入口与版本变化摘要服务 `backend/internal/service/marketplace/upgrade_advisor_service.go`、`backend/internal/service/marketplace/change_log_service.go`
- [X] T035 [US2] 在 `backend/internal/api/handler/marketplace_handler.go`、`backend/internal/api/router/marketplace_routes.go` 落地 `/marketplace/templates/{templateId}/releases`、`/marketplace/installations`
- [X] T036 [US2] 在 `backend/internal/service/audit/event_writer.go`、`backend/internal/service/audit/service.go` 打通模板分发、撤回、安装与升级动作的审计写入与查询聚合
- [X] T037 [P] [US2] 扩展前端服务与 hooks `frontend/src/services/platformMarketplace.ts`、`frontend/src/features/platform-marketplace/hooks/useTemplateDistributionAction.ts`
- [X] T038 [P] [US2] 实现模板分发页面与表单 `frontend/src/features/platform-marketplace/pages/TemplateDistributionPage.tsx`、`frontend/src/features/platform-marketplace/components/TemplateReleaseDrawer.tsx`
- [X] T039 [P] [US2] 实现安装记录与升级入口页面 `frontend/src/features/platform-marketplace/pages/InstallationRecordPage.tsx`、`frontend/src/features/platform-marketplace/components/UpgradeEntryDrawer.tsx`
- [X] T040 [US2] 在 `frontend/src/app/router.tsx`、`frontend/src/features/platform-marketplace/pages/InstallationRecordPage.tsx` 落地版本变化说明、下线提示和升级受限空态

**Checkpoint**: US2 可独立验证模板范围分发和安装升级闭环

---

## Phase 5: User Story 3 - 注册扩展并控制平台扩展边界 (Priority: P2)

**Goal**: 提供扩展包注册、启停、兼容性、权限声明和可见范围治理能力。  
**Independent Test**: 注册一个扩展包并声明兼容性、权限范围和可见范围后，可执行启停并查看兼容状态和审计记录。

### Tests for User Story 3

- [X] T041 [P] [US3] 编写后端契约测试 `backend/tests/contract/marketplace_extension_registry_contract_test.go`、`backend/tests/contract/marketplace_extension_compatibility_contract_test.go`
- [X] T042 [P] [US3] 编写后端集成测试 `backend/tests/integration/marketplace_extension_enable_flow_test.go`、`backend/tests/integration/marketplace_extension_visibility_flow_test.go`
- [X] T043 [P] [US3] 编写前端 Vitest 页面测试 `frontend/src/features/platform-marketplace/pages/ExtensionCenterPage.test.tsx`、`frontend/src/features/platform-marketplace/pages/ExtensionCompatibilityPage.test.tsx`

### Implementation for User Story 3

- [X] T044 [P] [US3] 实现扩展注册与生命周期服务 `backend/internal/service/marketplace/extension_registry_service.go`、`backend/internal/service/marketplace/extension_lifecycle_service.go`
- [X] T045 [P] [US3] 实现兼容性评估与权限声明校验服务 `backend/internal/service/marketplace/compatibility_service.go`、`backend/internal/service/marketplace/permission_guard_service.go`
- [X] T046 [P] [US3] 实现扩展可见范围与影响分析服务 `backend/internal/service/marketplace/extension_visibility_service.go`、`backend/internal/service/marketplace/impact_summary_service.go`
- [X] T047 [US3] 在 `backend/internal/api/handler/marketplace_handler.go`、`backend/internal/api/router/marketplace_routes.go` 落地 `/marketplace/extensions`、`/marketplace/extensions/{extensionId}/enable`、`/marketplace/extensions/{extensionId}/disable`、`/marketplace/extensions/{extensionId}/compatibility`
- [X] T048 [US3] 在 `backend/internal/api/handler/audit_handler.go`、`backend/internal/api/router/audit_routes.go` 聚合并暴露 `/audit/platform-marketplace/events` 查询链路
- [X] T049 [P] [US3] 扩展前端服务与 hooks `frontend/src/services/platformMarketplace.ts`、`frontend/src/features/platform-marketplace/hooks/useExtensionAction.ts`
- [X] T050 [P] [US3] 实现扩展中心与注册表单 `frontend/src/features/platform-marketplace/pages/ExtensionCenterPage.tsx`、`frontend/src/features/platform-marketplace/components/ExtensionPackageDrawer.tsx`
- [X] T051 [P] [US3] 实现兼容性与影响分析页面 `frontend/src/features/platform-marketplace/pages/ExtensionCompatibilityPage.tsx`、`frontend/src/features/platform-marketplace/components/CompatibilitySummaryCard.tsx`
- [X] T052 [P] [US3] 实现市场审计页面 `frontend/src/features/audit/pages/PlatformMarketplaceAuditPage.tsx`、`frontend/src/features/audit/pages/PlatformMarketplaceAuditPage.test.tsx`
- [X] T053 [US3] 在 `frontend/src/app/AuthorizedMenu.tsx`、`frontend/src/app/router.tsx`、`frontend/src/features/platform-marketplace/pages/ExtensionCenterPage.tsx` 落地兼容性阻断提示、权限声明展示和扩展停用影响空态

**Checkpoint**: US3 完成后形成应用目录、模板分发与扩展治理审计闭环

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: 收敛质量、文档、验证证据与 PR 交付材料

- [X] T054 [P] 收敛命名与共享类型，在 `backend/internal/service/marketplace/`、`backend/internal/integration/marketplace/`、`frontend/src/features/platform-marketplace/`、`frontend/src/services/platformMarketplace.ts` 清理重复字段与错误文案
- [X] T055 [P] 刷新配置与启动文档，在 `README.md`、`backend/config/config.example.yaml`、`backend/config/config.dev.yaml`、`frontend/.env.example`、`frontend/.env.development` 补齐 010 说明
- [X] T056 [P] 记录验证基线到 `artifacts/010-platform-marketplace/verification.md`、`artifacts/010-platform-marketplace/quickstart-validation.md`、`artifacts/010-platform-marketplace/repro-platform-marketplace-smoke.sh`
- [X] T057 在 `artifacts/010-platform-marketplace/pr-summary.md`、`artifacts/010-platform-marketplace/pr-readiness.md` 准备中文 PR 摘要、备份证据、测试说明、风险清单与用户合并确认项

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 0 (Governance Gates)**: 无依赖，必须先完成
- **Phase 1 (Setup)**: 依赖 Phase 0
- **Phase 2 (Foundational)**: 依赖 Phase 1，阻塞所有用户故事
- **Phase 3/4/5 (User Stories)**: 均依赖 Phase 2；US1 作为 MVP 优先，US2 在 US1 主干稳定后推进，US3 在 US1/US2 的目录资产、范围分发和权限边界语义完成后推进
- **Phase 6 (Polish)**: 依赖已完成的用户故事范围
- **Release / Merge**: 依赖远程推送、PR 更新、评审完成与用户明确同意

### User Story Dependencies

- **US1 (P1)**: 无用户故事前置依赖，Foundational 完成后可立即开始
- **US2 (P1)**: 依赖 US1 产出的目录来源、模板版本和范围语义，才能形成模板分发闭环
- **US3 (P2)**: 依赖 Foundational 的扩展注册与兼容性模型，以及 US1/US2 的模板/范围治理语义

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
Task: "T019 [US1] backend/tests/contract/marketplace_catalog_source_contract_test.go"
Task: "T020 [US1] backend/tests/integration/marketplace_catalog_sync_flow_test.go"
Task: "T021 [US1] frontend/src/features/platform-marketplace/pages/CatalogSourcePage.test.tsx"

# 并行实现任务
Task: "T022 [US1] backend/internal/service/marketplace/catalog_source_service.go"
Task: "T023 [US1] backend/internal/service/marketplace/template_version_service.go"
Task: "T026 [US1] frontend/src/services/platformMarketplace.ts"
```

---

## Implementation Strategy

### MVP First (US1 Only)

1. 完成 Phase 0-2
2. 完成 US1（Phase 3）
3. 按 Independent Test 验证 US1
4. 产出阶段性演示或 PR 更新

### Incremental Delivery

1. 先交付 US1：目录来源接入、模板中心、版本与依赖治理
2. 再交付 US2：模板按范围分发、安装记录、升级入口和下线状态
3. 最后交付 US3：扩展注册、兼容性、启停、权限声明和扩展审计
4. 最终执行 Phase 6 文档与验证收尾

### Notes

- `[P]` 任务代表可并行，但仍需满足前置依赖
- 每个用户故事都可独立验收
- 所有提交说明与 PR 摘要必须为中文
- 未获用户明确同意前禁止合并
