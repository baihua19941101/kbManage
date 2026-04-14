# Feature Specification: 多集群 GitOps 与应用发布中心

**Feature Branch**: `004-gitops-and-release`  
**Created**: 2026-04-13  
**Status**: In Progress (Tasks Ready)  
**Input**: User description: "我要新增 004-gitops-and-release，严格对标 Rancher 在 GitOps 和应用发布管理上的能力，并参考 Fleet 的多集群持续交付模式。面向平台管理员、应用交付团队和运维人员，在多集群 Kubernetes 场景下提供统一的 GitOps 与发布中心。用户需要能够接入代码仓库和发布来源，定义应用交付单元、目标集群或集群组、环境分层、配置覆盖和发布策略，查看期望状态与实际状态差异、同步结果、漂移状态、发布历史、失败原因和回滚入口；同时需要能够管理应用版本、配置版本、发布节奏和多环境推进过程，完成安装、升级、回滚、暂停、恢复和卸载等发布生命周期动作。平台必须延续工作空间、项目和环境范围的权限隔离，并对每一次发布、同步、回滚和配置变更保留可检索审计记录。首期范围聚焦 GitOps 持续交付与发布生命周期管理，不包含通用 CI 流水线编排、制品仓库管理、终端运维、策略准入和合规扫描。"

## 当前状态/执行说明（2026-04-13）

- 已完成 `/speckit.specify`，004 规格已创建并通过质量清单校验。
- 已完成 `/speckit.plan`，生成 `plan.md`、`research.md`、`data-model.md`、`contracts/openapi.yaml` 和 `quickstart.md`。
- 当前执行分支：`004-gitops-and-release`。
- 004 实施前数据库备份已于 2026-04-13 执行完成，产物为 `artifacts/004-gitops-and-release/mysql-backup-20260413-155243.sql`，并已通过临时容器恢复抽样验证。
- 已完成 `/speckit.tasks`，`tasks.md` 已生成并按用户故事、依赖关系和并行机会完成拆解。
- Phase 0 治理任务已完成：T001（分支与合并门槛核对）、T002（备份清单复核补充）、T003（国内源与远程流程核对）、T004（文档状态补齐）。
- Phase 1 配置任务已启动：T008 已完成配置项定义与注释补齐（gitops.sources/sync/diff/release/audit）。
- 当前状态更新为“tasks 阶段已完成治理门槛并启动 implement 前置配置，具备进入 `/speckit.implement` 的执行条件”，后续实现仍需遵守“提交中文 PR + 用户明确批准后再合并”。
- 2026-04-13 已补齐 US1 前端缺失页面/组件文件与路由串联，修复 `SourceFormDrawer.test.tsx` 的 QueryClientProvider 依赖，`npm run test -- --run src/features/gitops` 与 `npm run lint` 均通过。
- 2026-04-13 已完成 backend gitops 服务冲突收敛：删除 `source_service.go`、`target_group_service.go`、`environment_service.go`、`overlay_service.go`、`delivery_unit_service.go`、`status_service.go`，统一保留 `service.go` 作为实现；执行 `cd backend && go test ./...` 全量通过。
- 2026-04-13 已完成 US1（建模与状态查看 MVP），并通过后端 `go test` 与前端 GitOps 相关测试/lint。
- 2026-04-13 已补齐 US2 后端测试文件（T031、T032）：新增 3 个 contract + 3 个 integration，用最小断言覆盖 actions(202+operationId)、diff(200)、releases(200 list)、operation 查询(200)；`cd backend && go test ./...` 全量通过。
- 2026-04-13 已完成 US2（T033-T042）并收敛回归问题：修复 `delivery_operation_worker.go` 的进度快照字段编译错误，修复 `ReleaseActionDrawer`/`RollbackDialog` 与 `useDeliveryOperation` 的重置依赖导致的前端渲染循环；`cd backend && go test ./...`、`cd frontend && npm run test -- --run src/features/gitops`、`cd frontend && npm run lint` 全量通过。
- 2026-04-13 已完成 US3（T043-T051）：新增 GitOps 权限与审计 contract/integration 测试，落地 GitOps 细粒度鉴权中间件与路由权限分层，新增 `gitops.*` 审计动作写入与分类，补齐前端 GitOps 访问门禁测试、权限回收锁定态与 GitOps 审计查询页；`cd backend && go test ./...`、`cd frontend && npm run test -- --run src/features/gitops/pages/GitOpsAccessGate.test.tsx`、`cd frontend && npm run lint` 通过。
- 2026-04-13 已完成 Final Phase（T052-T055）：完成 GitOps 服务层命名收敛、README/quickstart 状态刷新，补齐 `artifacts/004-gitops-and-release/verification.md`、`quickstart-validation.md`、`repro-gitops-smoke.sh`、`pr-summary.md`、`pr-readiness.md`。

## User Scenarios & Testing *(mandatory)*

### User Story 1 - 统一交付源与多集群目标建模 (Priority: P1)

作为平台管理员或应用交付人员，我希望在一个统一入口中接入代码仓库与发布来源，定义应用交付单元、目标集群或集群组、环境分层和配置覆盖，这样我可以用声明式方式把同一应用稳定地分发到多个集群和环境。

**Why this priority**: GitOps 与发布中心的首要价值在于把“来源、目标、环境、覆盖关系”收敛到一个可持续维护的交付模型里；如果这个模型不存在，后续同步、差异检测和多环境推进都无法成立。

**Independent Test**: 接入至少一个代码仓库和一个发布来源后，授权用户能够创建一个应用交付单元，为其绑定多个目标集群或集群组、环境层级和配置覆盖，并看到每个目标的期望状态、最近同步结果和当前漂移状态。

**Acceptance Scenarios**:

1. **Given** 平台管理员已经配置可用的交付来源，**When** 应用交付人员创建新的应用交付单元并选择目标集群或集群组，**Then** 平台应保存其来源、目标范围、环境层级、配置覆盖和发布策略，并展示该单元的交付概览。
2. **Given** 同一应用需要面向不同环境使用不同配置，**When** 交付人员为测试、预发和生产环境分别配置覆盖项，**Then** 平台应明确展示继承关系、最终生效内容和适用目标范围。
3. **Given** 某个目标集群当前实际状态与期望状态不一致，**When** 用户查看该应用交付单元，**Then** 平台应标识差异、最近同步时间、同步结果和漂移状态，而不是只显示“已部署”。

---

### User Story 2 - 发布生命周期与多环境推进 (Priority: P1)

作为应用交付团队或运维人员，我希望在平台内完成安装、升级、回滚、暂停、恢复和卸载等发布生命周期动作，并按环境顺序推进版本与配置变更，这样我可以在多集群和多环境中有节奏地交付应用。

**Why this priority**: 统一建模之后，平台的核心价值在于把发布动作和环境推进纳入同一个受控流程。如果只能定义目标不能推进发布，就无法形成真正的交付中心。

**Independent Test**: 选择一个已建好的应用交付单元后，授权用户能够对其执行首次安装、版本升级、配置升级、暂停同步、恢复同步、按环境推进和版本回滚，并查看每一步的结果、失败原因和历史记录。

**Acceptance Scenarios**:

1. **Given** 某个应用交付单元已有可用版本与目标环境，**When** 交付人员发起安装或升级，**Then** 平台应展示影响范围、目标环境、目标版本、目标配置版本和执行结果。
2. **Given** 一个版本需要先在低风险环境验证再进入高风险环境，**When** 交付人员按既定顺序推进发布，**Then** 平台应记录每个环境阶段的开始时间、执行状态、完成结果和当前推进位置。
3. **Given** 最近一次升级导致异常，**When** 运维人员在发布历史中选择可恢复版本并发起回滚，**Then** 平台应提供明确的回滚入口、展示影响目标并返回成功或失败原因。
4. **Given** 某个应用需要暂时停止自动对齐，**When** 交付人员执行暂停或恢复动作，**Then** 平台应更新该单元的同步状态并保留暂停期间的变更可见性。

---

### User Story 3 - 权限隔离与发布审计闭环 (Priority: P2)

作为平台管理员或审计人员，我希望 GitOps 与发布能力继续遵守工作空间、项目和环境范围的授权边界，并对每一次发布、同步、回滚和配置变更形成可检索审计，这样既能支持持续交付，也能控制越权和追责风险。

**Why this priority**: 发布能力直接改变运行中环境，若没有范围隔离和审计闭环，会放大错误发布和越权操作风险，因此必须与发布能力同时成立。

**Independent Test**: 为两个不同工作空间或环境范围的用户分别授权后，他们只能看到并操作各自范围内的交付来源、应用交付单元、环境推进和发布历史；审计人员能够按时间、操作者、对象和结果检索完整发布记录。

**Acceptance Scenarios**:

1. **Given** 两个团队分别属于不同工作空间或项目范围，**When** 他们进入 GitOps 与发布中心，**Then** 平台应仅展示各自授权范围内的交付来源、应用交付单元、目标集群和发布记录。
2. **Given** 某个用户只有查看权限或仅能操作低风险环境，**When** 其尝试执行升级、回滚、暂停、恢复或卸载动作，**Then** 平台应拒绝未授权动作并返回明确原因。
3. **Given** 平台内已经发生多次发布与配置变更，**When** 审计人员按时间范围、操作者、工作空间、项目、环境、应用或动作类型检索，**Then** 平台应返回完整的记录上下文和执行结果。

### Edge Cases

- 当同一应用在不同环境中使用不同版本与不同配置时，平台必须清晰展示各环境的生效版本和覆盖来源，避免把环境差异误判为异常漂移。
- 当代码仓库、发布来源或目标集群暂时不可达时，平台必须区分“来源不可用”“目标不可达”和“无待同步变更”三种状态。
- 当某次发布只在部分目标集群或集群组成员上成功时，平台必须明确显示整体状态与单目标结果，不能把部分成功显示为整体完成。
- 当平台外的人工变更导致实际状态偏离期望状态时，平台必须将其标记为漂移，并保留后续同步或忽略处理的可追踪记录。
- 当用户在环境推进过程中权限被收回时，平台必须立即阻止其继续推进后续环境或执行回滚动作。
- 当可回滚版本引用的源版本、配置版本或目标范围已经失效时，平台必须在执行前或执行后给出可理解的阻断或失败说明。
- 当应用处于暂停状态时，平台必须继续展示待同步差异和最近变更，不能把暂停误显示为“已完成对齐”。
- 当目标集群被移出集群组或环境绑定关系变更时，平台必须明确说明哪些目标将停止接收后续发布。

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: 系统 MUST 提供统一的多集群 GitOps 与应用发布中心，支持围绕交付来源、应用交付单元、目标集群或集群组、环境和发布记录进行集中管理。
- **FR-002**: 系统 MUST 支持平台管理员接入、启用、停用和移除交付来源，并维护每个来源的连接状态、身份信息和最近一次验证结果。
- **FR-003**: 系统 MUST 支持授权用户创建和维护应用交付单元，为其绑定来源、目标集群或集群组、环境层级、配置覆盖和发布策略。
- **FR-004**: 系统 MUST 支持将多个目标集群组织为可复用的交付目标集合，并允许应用交付单元引用这些集合进行批量发布。
- **FR-005**: 系统 MUST 支持为不同环境定义有序层级，并将目标集群、配置覆盖和发布节奏约束到具体环境范围。
- **FR-006**: 系统 MUST 支持为应用交付单元定义基础配置与环境级覆盖配置，并在发布前展示最终生效内容与适用范围。
- **FR-007**: 系统 MUST 展示每个应用交付单元在每个目标和环境下的期望状态、实际状态、最近同步时间、最近同步结果和当前漂移状态。
- **FR-008**: 系统 MUST 支持查看期望状态与实际状态的差异摘要，并帮助用户定位差异发生的环境、目标对象和变更来源。
- **FR-009**: 系统 MUST 支持授权用户触发同步、重新同步和状态刷新，并返回明确的执行状态与结果说明。
- **FR-010**: 系统 MUST 识别并标记平台外人工变更或目标状态偏离造成的漂移，并在后续视图中持续可见直到被处理。
- **FR-011**: 系统 MUST 支持应用交付单元的安装、升级、回滚、暂停、恢复和卸载等发布生命周期动作。
- **FR-012**: 系统 MUST 在执行高影响发布动作前展示影响范围、目标环境、目标集群或集群组、版本信息和风险提示。
- **FR-013**: 系统 MUST 支持同时管理应用版本与配置版本，并在每次发布记录中分别记录这两类版本标识。
- **FR-014**: 系统 MUST 支持按既定环境顺序推进发布，并展示每个环境阶段的当前状态、进入时间、完成时间和推进结果。
- **FR-015**: 首期多环境推进流程 MUST 采用按环境顺序的受控推进，不包含通用 CI 流水线编排或自定义工作流引擎。
- **FR-016**: 系统 MUST 为每个应用交付单元保留完整发布历史，至少包含发起人、环境、目标范围、应用版本、配置版本、动作类型、开始时间、结束时间、执行结果和失败原因。
- **FR-017**: 系统 MUST 为可恢复的发布历史提供明确的回滚入口，并在回滚前展示回滚目标和受影响范围。
- **FR-018**: 系统 MUST 支持暂停和恢复自动对齐或发布推进，并在暂停期间持续展示待处理差异与未完成变更。
- **FR-019**: 系统 MUST 继承现有工作空间、项目和环境范围隔离，对交付来源、应用交付单元、目标集合、环境推进、发布动作和历史记录统一执行权限校验。
- **FR-020**: 系统 MUST 在用户权限变化后立即按新权限限制其可见的来源、应用、目标、环境和可执行动作。
- **FR-021**: 系统 MUST 对来源变更、配置变更、同步、重新同步、安装、升级、回滚、暂停、恢复和卸载生成可检索审计记录。
- **FR-022**: 系统 MUST 支持审计人员按时间范围、操作者、工作空间、项目、环境、应用交付单元、目标范围、动作类型和结果筛选交付记录。
- **FR-023**: 系统 MUST 在来源不可达、目标不可达、版本缺失、配置冲突、部分成功、权限不足和对象已变化等场景下返回可理解的错误说明。
- **FR-024**: 首期交付来源类型 MUST 聚焦代码仓库与应用发布来源，不扩展到通用制品仓库管理。
- **FR-025**: 首期范围 MUST 聚焦 GitOps 持续交付与发布生命周期管理，不包含通用 CI 流水线编排、制品仓库管理、终端运维、策略准入和合规扫描。

## Governance & Delivery Constraints *(mandatory)*

- **GC-001**: Feature work MUST occur on a dedicated feature branch; direct development on `master` or `main` is forbidden.
- **GC-002**: All user-facing communication, approval records, PR summaries, and delivery notes MUST be written in Chinese.
- **GC-003**: Any dependency or framework installation MUST document the China mirror or proxy configuration that will be used during implementation.
- **GC-004**: Before implementation begins, the feature specification or plan MUST record a database backup executed from container `mysql8` using `localhost:3306` and credentials `admin/123456`, or explicitly justify why the backup requirement is not applicable.
- **GC-005**: Delivery MUST include pushing the feature branch to the GitHub remote and opening or updating a PR; the next feature MUST NOT start until the current PR flow is complete.
- **GC-006**: Merge to the mainline branch MUST NOT occur without explicit user approval.
- **GC-007**: If subagents are used for implementation, they MUST use `gpt-5.3-codex`.

### Key Entities *(include if feature involves data)*

- **Delivery Source**: 一个可被平台接入和验证的交付来源，表示代码仓库或应用发布来源，并包含连接状态、身份信息和可见范围。
- **Application Delivery Unit**: 一个被持续交付管理的应用单元，包含来源引用、目标范围、环境层级、配置覆盖、发布策略、版本信息和当前交付状态。
- **Target Group**: 一组可复用的发布目标集合，用于表达多个集群或环境下的统一分发范围。
- **Environment Stage**: 一个可排序的交付环境阶段，定义发布顺序、目标范围和环境级配置边界。
- **Configuration Overlay**: 针对某个环境或目标范围追加的配置覆盖，用于表达基础配置与环境差异。
- **Sync Record**: 一次期望状态对齐或重新对齐的执行记录，包含触发原因、目标范围、执行状态、结果摘要和漂移处理信息。
- **Release Revision**: 一次具体的应用版本与配置版本发布记录，用于历史追踪、失败诊断和回滚定位。
- **Release Audit Event**: 围绕来源变更、配置变更、同步、发布和回滚形成的标准化审计对象。

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 平台管理员或交付人员能够在 15 分钟内完成一个新应用交付单元的来源接入、目标绑定、环境分层和首次交付配置定义。
- **SC-002**: 在至少 20 个已接入集群的环境中，90% 的应用交付单元状态查询能够在 30 秒内返回目标环境的期望状态、实际状态、同步结果和漂移标记。
- **SC-003**: 对于最多覆盖 50 个目标集群的单次同步或发布动作，95% 能在 5 分钟内向用户返回明确的成功、失败或部分成功结果。
- **SC-004**: 在试点团队的多环境发布演练中，90% 的版本推进能够按既定环境顺序完成并形成完整的阶段记录与回滚入口。
- **SC-005**: 在权限验收中，100% 的跨工作空间、跨项目和跨环境未授权发布访问都被拦截，且不会暴露目标对象的详细交付信息。
- **SC-006**: 审计人员针对最近 90 天的交付记录检索时，90% 的查询能够在 30 秒内返回满足筛选条件的记录集。
- **SC-007**: 在试点期内，至少 80% 的多集群应用发布、版本升级和回滚操作可仅通过该平台完成，无需切换到额外发布控制台。

## Assumptions

- `001-k8s-ops-platform` 已提供多集群接入、工作空间/项目级授权、基础资源范围模型和审计基础能力。
- `003-workload-operations-control-plane` 已提供工作负载级运维动作与回滚能力；本特性聚焦应用交付声明管理和发布生命周期，不替代运行时运维入口。
- 首期目标用户为平台管理员、应用交付团队和运维人员，不单独覆盖移动端专属体验。
- 首期默认按环境顺序进行受控推进，环境之间的自动审批流、通用工作流编排和 CI 构建过程不纳入本轮范围。
- 首期聚焦代码仓库与应用发布来源接入，不负责通用制品仓库本身的存储治理与制品生命周期管理。
- 平台默认管理的是组织已授权且可持续访问的交付来源与 Kubernetes 集群；若外部来源短时异常，平台以状态提示和失败记录呈现，而不是静默忽略。
