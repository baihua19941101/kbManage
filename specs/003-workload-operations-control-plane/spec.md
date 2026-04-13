# Feature Specification: 多集群 Kubernetes 工作负载运维控制面

**Feature Branch**: `003-workload-operations-control-plane`  
**Created**: 2026-04-12  
**Status**: In Progress  
**Input**: User description: "我要新增 003-workload-operations-control-plane，严格对标 Rancher 的工作负载运维能力。面向运维人员、SRE、项目负责人和平台管理员，在已接入并已授权的多集群 Kubernetes 环境中，提供围绕工作负载和实例的持续运维能力。用户需要能够在单个资源上下文中查看运行状态、发布进度、实例分布和最近变更，查看并跟踪 Pod 或容器日志，进入容器终端执行诊断命令，对工作负载执行扩缩容、重启、重新部署、实例替换、批量操作、发布历史查看和版本回滚等动作，并在执行前看到影响范围、在执行中看到进度、在执行后看到结果和失败原因。平台必须支持按工作空间、项目、命名空间和资源范围进行权限控制，并对每一次高风险操作、终端使用和回滚动作记录完整审计。首期范围聚焦工作负载资源的运维闭环，不包含全局日志中心、统一监控告警、GitOps 持续交付、Helm 发布生命周期、策略治理和集群生命周期管理。"

## 当前状态/执行说明（2026-04-12）

- 已完成 `/speckit.specify` 与 `/speckit.plan` 产物生成，`plan.md`、`research.md`、`data-model.md`、`contracts/openapi.yaml` 与 `quickstart.md` 已创建。
- 当前工作分支为 `003-workload-operations-control-plane`。
- 003 本轮实施前数据库备份已于 2026-04-12 执行，产物为 `artifacts/003-workload-operations-control-plane/mysql-backup-20260412-150123.sql`，并已通过临时容器恢复抽样验证。
- 已确认终端审计首期仅记录会话建立、关闭、目标容器、操作者、持续时长和结束原因，不记录完整命令与终端输出正文。
- 已完成 `/speckit.tasks`，任务清单位于 `tasks.md`，并已开始 `/speckit.implement`。
- 已完成 Phase 0/1/2：治理证据文件、后端 `workloadops` 基础骨架、前端 `workload-ops` 入口占位与共享底座已落库到代码仓。
- 已完成 US1（单资源诊断入口）、US2（动作执行与发布恢复）与 US3（权限隔离与高风险审计闭环），并通过当前自动化测试基线（后端 contract/integration + 前端 Vitest + 变更文件 ESLint）。
- 当前状态更新为“Implement 已完成（Governance + Setup + Foundational + US1 + US2 + US3 + Final Phase）”；后续交付仍需遵守“中文 PR + 用户明确批准后再合并”。

## User Scenarios & Testing *(mandatory)*

### User Story 1 - 单资源运维诊断入口 (Priority: P1)

作为运维人员或 SRE，我希望在单个工作负载资源上下文中持续查看运行状态、发布进度、实例分布、最近变更以及关联的 Pod/容器日志和终端入口，这样我可以在不切换多个工具的情况下快速完成定位和初步诊断。

**Why this priority**: 如果没有单资源运维诊断入口，后续所有扩缩容、重启、替换、回滚等动作都缺少可靠的上下文基础，平台无法形成可用的工作负载运维入口。

**Independent Test**: 在至少两个已接入并授权的集群中选择一个异常工作负载，用户能够在同一资源页面完成状态查看、发布进度跟踪、实例下钻、日志查看和容器终端进入，并据此识别异常实例。

**Acceptance Scenarios**:

1. **Given** 用户已被授权访问某个工作负载，**When** 打开该资源的运维视图，**Then** 平台应展示该资源的运行状态、发布进度、实例分布、最近变更和所属集群/工作空间/项目/命名空间上下文。
2. **Given** 某个工作负载下存在多个 Pod 或容器实例，**When** 用户切换实例查看日志或终端入口，**Then** 平台应保留当前资源上下文并明确显示所选实例与容器标识。
3. **Given** 某个实例正在重建、已退出或暂时不可连接，**When** 用户尝试查看日志或进入终端，**Then** 平台应返回可理解的状态说明，并明确区分对象已变化、无可用日志和终端不可建立等情况。

---

### User Story 2 - 工作负载动作执行与发布恢复 (Priority: P1)

作为运维人员、项目负责人或平台管理员，我希望围绕工作负载执行扩缩容、重启、重新部署、实例替换、批量操作、发布历史查看和版本回滚等动作，并在执行前看到影响范围、执行中看到进度、执行后看到结果和失败原因，这样我可以在平台内完成持续运维闭环。

**Why this priority**: 工作负载动作执行和发布恢复是 Rancher 对标范围内最核心的运维价值。如果只能看不能管，平台仍然只是只读视图，无法承担生产运维入口角色。

**Independent Test**: 选择一个被授权的工作负载并执行扩缩容、重启或回滚动作后，用户能够看到影响预览、执行状态流转、最终结果以及失败原因；对多个资源执行批量动作时能够区分整体结果和单项结果。

**Acceptance Scenarios**:

1. **Given** 用户拥有目标工作负载的运维权限，**When** 发起扩缩容、重启、重新部署或实例替换动作，**Then** 平台应在执行前展示影响范围和风险提示，并在执行后返回明确结果。
2. **Given** 某个工作负载存在可用的发布历史，**When** 用户查看历史版本并发起回滚，**Then** 平台应展示可回滚目标、影响范围、执行进度和最终回滚结果。
3. **Given** 用户一次性选择多个工作负载执行同类动作，**When** 批量任务运行，**Then** 平台应展示批量任务总进度、每个目标对象的独立状态以及失败对象的具体原因。

---

### User Story 3 - 权限隔离与高风险审计闭环 (Priority: P2)

作为平台管理员或审计人员，我希望工作负载运维能力继续遵守工作空间、项目、命名空间和资源范围的授权边界，并对高风险动作、终端使用和回滚过程形成完整审计，这样可以在支持高效运维的同时控制越权和追责风险。

**Why this priority**: 工作负载运维能力直接触达生产运行对象，缺少权限隔离和审计会放大误操作与合规风险，因此必须和动作能力一起形成闭环。

**Independent Test**: 为两个不同授权范围的用户分别登录平台后，他们只能访问各自范围内的工作负载、实例、终端和动作入口；执行高风险动作、进入终端或发起回滚后，审计人员能够检索到完整记录。

**Acceptance Scenarios**:

1. **Given** 两个用户分别属于不同工作空间或项目范围，**When** 他们访问工作负载运维页面，**Then** 平台应只展示各自授权范围内的工作负载、实例、日志和可执行动作。
2. **Given** 某个用户仅具备查看权限或受限运维权限，**When** 访问终端、回滚或批量高风险动作入口，**Then** 平台应阻止未授权行为并返回明确原因，不暴露超出范围的对象细节。
3. **Given** 某个用户已经执行终端访问或高风险动作，**When** 审计人员按时间、操作者、资源或动作类型检索，**Then** 平台应返回包含目标对象、影响范围、执行结果和关键说明的完整记录。

### Edge Cases

- 当工作负载在用户浏览期间发生新一轮发布、被外部修改或已被删除时，平台必须提示当前视图基于过期状态，并阻止对失效版本继续执行高风险动作。
- 当某个批量动作只对部分目标成功时，平台必须明确展示成功、失败和未执行对象，避免把部分成功误解为整体完成。
- 当用户选中的 Pod、容器或发布版本在执行日志查看、终端进入或回滚之前已经不存在时，平台必须返回对象已变化的说明，而不是静默失败。
- 当终端会话空闲超时、连接中断或目标容器重启时，平台必须明确标记会话已失效，并要求重新建立连接。
- 当回滚目标版本与当前配置不兼容或缺少可恢复对象时，平台必须在执行前给出阻断性提示或在执行后返回可理解的失败原因。
- 当用户在操作过程中权限被收回时，平台必须立即阻止后续日志、终端和运维动作访问，并按新权限返回结果。

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: 系统 MUST 提供统一的工作负载运维入口，支持围绕单个工作负载查看运行状态、发布进度、实例分布、最近变更和资源上下文。
- **FR-002**: 系统 MUST 支持用户在工作负载上下文中查看关联实例列表，并从实例下钻到 Pod 或容器级别的状态、标识和最近变化。
- **FR-003**: 系统 MUST 支持按工作负载、实例和容器查看并跟踪运行日志，并明确展示所属集群、工作空间、项目、命名空间和资源对象。
- **FR-004**: 系统 MUST 支持授权用户从资源上下文进入容器终端执行诊断命令，并展示当前会话所关联的目标实例和容器。
- **FR-005**: 系统 MUST 在终端建立失败、会话超时、实例重建、容器退出或权限不足时返回可理解的状态说明。
- **FR-006**: 系统 MUST 支持授权用户对工作负载执行扩缩容、重启、重新部署和实例替换等持续运维动作。
- **FR-007**: 系统 MUST 支持批量选择多个工作负载执行同类运维动作，并为每个目标对象保留独立状态与结果。
- **FR-008**: 系统 MUST 在执行高影响动作前展示影响范围、目标对象、潜在风险和二次确认信息。
- **FR-009**: 系统 MUST 支持查看工作负载发布历史，并在可恢复条件满足时发起版本回滚。
- **FR-010**: 系统 MUST 在运维动作执行期间持续展示整体进度、当前状态、已完成项、失败项和失败原因。
- **FR-011**: 系统 MUST 在运维动作结束后保留结果摘要，至少包含发起人、目标对象、动作类型、开始时间、结束时间、执行结果和必要说明。
- **FR-012**: 系统 MUST 区分普通动作与高风险动作，并至少将批量高影响变更、实例替换、版本回滚和终端使用纳入高风险审计范围。
- **FR-013**: 系统 MUST 继承现有工作空间、项目、命名空间和资源范围授权模型，对工作负载视图、日志、终端和运维动作统一执行权限校验。
- **FR-014**: 系统 MUST 在用户权限变化后立即按新权限限制其对工作负载、实例、日志、终端和动作入口的访问。
- **FR-015**: 系统 MUST 对每一次高风险动作、终端使用和回滚动作生成可检索审计记录，记录操作者、目标对象、影响范围、发生时间、执行结果和关键说明。
- **FR-016**: 系统 MUST 支持审计人员按时间范围、操作者、工作空间、项目、命名空间、资源对象、动作类型和结果检索工作负载运维记录。
- **FR-017**: 系统 MUST 在工作负载运维视图中清晰标记对象当前状态，区分正常、发布中、部分可用、异常、已变化和不可操作等情况。
- **FR-018**: 系统 MUST 在首期支持 `Deployment`、`StatefulSet`、`DaemonSet` 作为主要工作负载对象，并允许围绕其关联的 Pod 与容器执行查看、诊断和运维动作。
- **FR-019**: 系统 MUST 将发布历史、版本回滚和持续运维闭环限定在工作负载资源范围内，不扩展到全局日志中心、统一监控告警、GitOps 持续交付、Helm 发布生命周期、策略治理和集群生命周期管理。

## Governance & Delivery Constraints *(mandatory)*

- **GC-001**: Feature work MUST occur on a dedicated feature branch; direct development on `master` or `main` is forbidden.
- **GC-002**: All user-facing communication, approval records, PR summaries, and delivery notes MUST be written in Chinese.
- **GC-003**: Any dependency or framework installation MUST document the China mirror or proxy configuration that will be used during implementation.
- **GC-004**: Before implementation begins, the feature specification or plan MUST record a database backup executed from container `mysql8` using `localhost:3306` and credentials `admin/123456`, or explicitly justify why the backup requirement is not applicable.
- **GC-005**: Delivery MUST include pushing the feature branch to the GitHub remote and opening or updating a PR; the next feature MUST NOT start until the current PR flow is complete.
- **GC-006**: Merge to the mainline branch MUST NOT occur without explicit user approval.
- **GC-007**: If subagents are used for implementation, they MUST use `gpt-5.3-codex`.

### Key Entities *(include if feature involves data)*

- **Workload Operation Context**: 一次围绕单个工作负载展开的运维视图上下文，包含资源标识、所属范围、运行状态、发布状态、实例分布和最近变更。
- **Workload Instance**: 某个工作负载下的具体运行实例，表示一个可被查看、诊断或替换的 Pod/容器执行单元。
- **Operation Batch**: 一次面向一个或多个工作负载发起的运维任务，记录目标集合、动作类型、整体进度、单项结果和失败原因。
- **Release Revision**: 某个工作负载的可识别发布历史版本，用于查看变更轨迹、判断恢复点和发起版本回滚。
- **Terminal Session Record**: 一次终端访问行为的审计对象，记录目标实例、会话起止、操作者、权限范围和结束状态。
- **Workload Audit Event**: 围绕高风险动作、终端访问、回滚和批量变更形成的标准化审计记录。

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 在至少 20 个已接入集群的环境中，90% 的授权用户能够在 3 分钟内从目标工作负载进入实例诊断视图并定位到目标 Pod 或容器。
- **SC-002**: 对单个工作负载发起的常见运维动作中，95% 能在 2 分钟内向用户返回明确的成功或失败结果。
- **SC-003**: 对最多 50 个目标对象的同类批量动作中，90% 的任务能够在 5 分钟内完成结果归集，并清晰区分成功项与失败项。
- **SC-004**: 在有可用发布历史的工作负载上，90% 的回滚操作可在 3 分钟内完成并返回明确的结果说明。
- **SC-005**: 在试点团队的日常排障场景中，至少 80% 的工作负载级诊断与恢复任务可仅通过该平台完成，无需切换到其他运维入口。
- **SC-006**: 在权限验收中，100% 的跨工作空间、跨项目、跨命名空间和超出资源范围的未授权工作负载运维访问都被拦截，且不会暴露目标对象的详细运行信息。

## Assumptions

- `001-k8s-ops-platform` 已提供多集群接入、资源索引、工作空间/项目范围模型、运维动作基础链路和审计基础能力。
- `002-observability-center` 已提供资源上下文跳转、实例相关日志查看和基础观测访问能力，本特性在此基础上向工作负载运维闭环延伸。
- 首期目标用户为运维人员、SRE、项目负责人和平台管理员，不单独覆盖移动端专属体验。
- 首期主要工作负载对象限定为 `Deployment`、`StatefulSet`、`DaemonSet`，Pod 与容器作为实例级诊断与执行上下文存在。
- 目标集群中的被管理工作负载已具备可读取运行状态、实例信息和发布历史的前提条件；若部分对象缺少历史信息，平台需要明确提示而不是静默降级。
- 首期不扩展到全局日志中心、统一监控告警、GitOps 持续交付、Helm 发布生命周期、策略治理和集群生命周期管理。
