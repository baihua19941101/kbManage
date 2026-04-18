# Feature Specification: 多集群 Kubernetes 集群生命周期中心

**Feature Branch**: `007-cluster-lifecycle`  
**Created**: 2026-04-17  
**Status**: Implemented (Pending Remote Push / PR)  
**Input**: User description: "我要新增 007-cluster-lifecycle，严格对标 Rancher 的集群生命周期管理能力，覆盖注册、导入、创建、升级、停用和退役全链路，并预留能力矩阵与驱动扩展机制。面向平台管理员和平台工程团队，在多集群 Kubernetes 环境中提供统一的集群生命周期中心。用户需要能够导入现有集群、注册新集群、按基础设施类型或驱动能力创建集群，管理集群的基础信息、版本、节点池、接入状态、健康状态、升级计划和退役流程，并查看不同类型集群在网络、存储、身份、监控、安全、备份和发布等方面的能力矩阵与兼容状态。平台还需要支持驱动版本管理、驱动能力扩展、模板化创建和创建前校验，降低不同基础设施场景下的接入差异。该特性必须延续现有授权模型，并对创建、导入、升级、扩缩节点和退役等关键动作保留完整审计。首期范围聚焦集群全生命周期和能力矩阵管理，不包含平台级灾备、统一身份源整合、策略治理和应用市场。"

## 当前状态/执行说明（2026-04-17）

- 已完成 `/speckit.specify`，007 规格与质量清单已生成。
- 已完成 `/speckit.plan`，当前已生成 `plan.md`、`research.md`、`data-model.md`、`contracts/openapi.yaml` 和 `quickstart.md`。
- 已完成 `/speckit.tasks`，`tasks.md` 已生成并按 Governance、Setup、Foundational、US1、US2、US3、Polish 分阶段组织。
- 当前执行分支为 `007-cluster-lifecycle`。
- 已执行 `/speckit.implement`，后端路由、仓储、领域模型、核心服务、审计接入与前端生命周期页面/菜单/权限门禁已完成首轮落地。
- 已补齐 007 剩余 contract/integration 测试文件、缓存锁辅助文件、服务拆分包装文件，以及从集群总览进入生命周期中心的导航入口。
- 本地验证已完成：`cd backend && go test -run TestNonExistent -count=0 ./...`、`cd backend && go test ./tests/contract -count=1 -p 1`、`cd backend && go test ./tests/integration -count=1 -p 1`、`cd frontend && npm run lint`、`cd frontend && npm run build`。
- 前端 007 定向 `vitest` 在单 worker 模式下存在测试进程退出缓慢现象，已记录到验证材料，不作为本轮实现阻塞项。
- 当前状态为“实现已本地闭环验证完成，含命名/共享类型清理；待远程推送与 PR”；后续仍需继续遵守中文 PR、远程推送和用户明确同意后再合并的治理要求。

## User Scenarios & Testing *(mandatory)*

### User Story 1 - 导入与注册已有集群 (Priority: P1)

作为平台管理员，我希望在统一入口中导入或注册现有 Kubernetes 集群，并持续看到其接入状态、健康状态和基础元信息，这样我可以把存量集群纳入同一控制面管理。

**Why this priority**: 没有导入与注册能力，平台无法建立多集群生命周期管理的基础资产视图，后续创建、升级和退役流程都无从谈起。

**Independent Test**: 通过导入一个已存在集群和注册一个待接入集群，管理员能够在同一列表中看到它们的接入进度、健康摘要、版本和基础属性，并区分成功、失败和待处理状态。

**Acceptance Scenarios**:

1. **Given** 管理员持有可用的现有集群接入信息，**When** 提交导入请求，**Then** 平台应创建对应集群记录并显示接入状态、版本信息和最近校验结果。
2. **Given** 管理员需要注册一个新纳管集群，**When** 生成注册指引并完成接入动作，**Then** 平台应将该集群标记为已纳管并展示其基础健康与连通状态。
3. **Given** 某个集群导入或注册失败，**When** 管理员查看详情，**Then** 平台应明确展示失败阶段、原因说明和可重试入口。

---

### User Story 2 - 创建、升级与退役集群 (Priority: P1)

作为平台工程团队成员，我希望按基础设施类型或驱动能力模板化创建集群，并在后续统一管理版本升级、节点池调整、停用和退役流程，这样我可以降低不同基础设施场景下的运维差异并保持生命周期可控。

**Why this priority**: 集群生命周期管理的核心价值在于标准化创建和受控变更，如果只能接入而不能持续管理，平台无法对标 Rancher 的全链路能力。

**Independent Test**: 使用一个驱动模板创建集群，随后为该集群制定升级计划、调整节点池并执行停用或退役流程，管理员能够看到各阶段校验结果、执行状态和审计记录。

**Acceptance Scenarios**:

1. **Given** 平台已提供某类基础设施驱动和可用模板，**When** 管理员按模板填写参数并发起创建，**Then** 平台应在创建前完成校验，并在创建后展示阶段性进度和最终结果。
2. **Given** 某个集群处于可升级状态，**When** 管理员选择目标版本并发起升级计划，**Then** 平台应展示兼容性检查、影响提示、执行进度和升级结果。
3. **Given** 某个集群将停止提供服务，**When** 管理员执行停用或退役流程，**Then** 平台应要求完成必要确认并记录状态变更、残留风险和最终退役结论。

---

### User Story 3 - 能力矩阵与驱动扩展管理 (Priority: P2)

作为平台管理员，我希望查看不同类型集群在网络、存储、身份、监控、安全、备份和发布等方面的能力矩阵与兼容状态，并管理驱动版本、驱动能力和模板资产，这样我可以在创建或升级前做出可预期的选择。

**Why this priority**: 生命周期管理不仅是动作执行，还需要在动作前建立清晰的能力边界和驱动约束，避免在不同基础设施之间产生隐性兼容问题。

**Independent Test**: 为两个不同驱动类型准备能力定义和模板后，管理员能够比较其能力矩阵、驱动版本和兼容状态，并基于这些信息选择可用创建模板。

**Acceptance Scenarios**:

1. **Given** 平台中存在多个驱动与集群类型，**When** 管理员查看能力矩阵，**Then** 平台应按能力域展示支持状态、兼容结论和差异说明。
2. **Given** 平台引入新的驱动版本或能力扩展，**When** 管理员更新驱动信息，**Then** 平台应保留版本差异并影响后续模板可用性与创建前校验结果。
3. **Given** 某个模板依赖的驱动能力不满足当前目标环境，**When** 管理员尝试使用该模板创建集群，**Then** 平台应在创建前阻止提交并说明缺失能力。

### Edge Cases

- 当同一个集群被重复导入、重复注册或以不同标识重复接入时，平台必须识别冲突并阻止产生重复资产记录。
- 当驱动版本升级后导致既有模板部分字段失效时，平台必须明确标识模板风险，而不是继续按旧模板静默创建。
- 当集群正在升级或退役过程中又收到节点池扩缩或停用请求时，平台必须阻止冲突动作并说明当前锁定状态。
- 当目标基础设施校验只部分通过时，平台必须区分“可继续但有风险”和“必须阻断”两类结果。
- 当集群已失联但尚未正式退役时，平台必须保留历史状态和最近一次可用信息，不得直接视为已删除。
- 当用户在执行创建、导入、升级、扩缩节点或退役过程中权限被回收时，平台必须立即阻止后续关键动作。

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: 系统 MUST 提供统一的集群生命周期中心，用于查看和管理已纳管、待纳管、创建中、升级中、停用中和退役中的集群。
- **FR-002**: 系统 MUST 支持导入现有 Kubernetes 集群，并在导入过程中记录接入状态、最近校验结果、版本信息和失败原因。
- **FR-003**: 系统 MUST 支持注册新集群并提供明确的接入指引、注册状态和完成确认机制。
- **FR-004**: 系统 MUST 支持按基础设施类型或驱动能力创建新集群，并在创建前执行必需校验。
- **FR-005**: 系统 MUST 在集群创建前展示模板信息、依赖条件、校验结果和阻断项。
- **FR-006**: 系统 MUST 支持维护集群基础信息，至少包括名称、所属范围、基础设施类型、控制面版本、节点池概览、接入状态和健康状态。
- **FR-007**: 系统 MUST 支持查看和管理集群的节点池信息，包括节点池角色、容量状态和可执行的扩缩调整动作。
- **FR-008**: 系统 MUST 支持为集群制定和执行升级计划，并展示目标版本、兼容检查、影响提示、执行进度和结果结论。
- **FR-009**: 系统 MUST 在执行升级、扩缩节点、停用和退役等高影响动作前提供风险提示和确认步骤。
- **FR-010**: 系统 MUST 支持将集群标记为停用中、已停用、退役中和已退役，并保留各阶段状态说明。
- **FR-011**: 系统 MUST 在退役流程中要求记录退役原因、范围确认和最终结论，并保留可复盘信息。
- **FR-012**: 系统 MUST 支持展示不同类型集群在网络、存储、身份、监控、安全、备份和发布等能力域的支持状态与兼容结论。
- **FR-013**: 系统 MUST 支持维护驱动版本信息，并将驱动版本与可用模板、能力矩阵和创建前校验结果关联。
- **FR-014**: 系统 MUST 支持扩展驱动能力定义，使平台可以新增或调整不同驱动的能力描述与适用范围。
- **FR-015**: 系统 MUST 支持模板化创建集群，并允许模板与驱动版本、基础设施类型和能力要求建立关联。
- **FR-016**: 系统 MUST 在模板、驱动能力或目标环境不兼容时阻止创建或升级继续提交，并返回明确原因。
- **FR-017**: 系统 MUST 继承现有授权模型，确保用户仅能查看和操作其被授权范围内的集群生命周期对象。
- **FR-018**: 系统 MUST 在用户权限变化后立即限制其对集群创建、导入、升级、扩缩节点、停用和退役动作的访问。
- **FR-019**: 系统 MUST 对创建、导入、注册、升级、扩缩节点、停用、退役、驱动变更和模板变更等关键动作生成可检索审计记录。
- **FR-020**: 系统 MUST 为每条审计记录保留操作者、目标集群或模板对象、动作类型、阶段结果、发生时间和必要说明。
- **FR-021**: 系统 MUST 支持按时间、操作者、基础设施类型、驱动类型、动作类型、集群状态和结果筛选集群生命周期记录。
- **FR-022**: 系统 MUST 在导入失败、创建前校验失败、驱动不兼容、升级冲突、节点池调整失败和退役阻断等场景下返回可理解的状态说明。
- **FR-023**: 首期范围 MUST 聚焦集群注册、导入、创建、升级、停用、退役以及能力矩阵和驱动扩展管理。
- **FR-024**: 首期范围 MUST 不包含平台级灾备、统一身份源整合、策略治理和应用市场能力。
- **FR-025**: 首期范围 MUST 预留驱动扩展机制和能力矩阵建模空间，但不要求一次性覆盖所有可能基础设施类型。

## Governance & Delivery Constraints *(mandatory)*

- **GC-001**: Feature work MUST occur on a dedicated feature branch; direct development on
  `master` or `main` is forbidden.
- **GC-002**: All user-facing communication, approval records, PR summaries, and delivery notes
  MUST be written in Chinese.
- **GC-003**: Any dependency or framework installation MUST document the China mirror or proxy
  configuration that will be used during implementation.
- **GC-004**: Before implementation begins, the feature specification or plan MUST record a
  database backup executed from container `mysql8` using `localhost:3306` and credentials
  `admin/123456`, or explicitly justify why the backup requirement is not applicable.
- **GC-005**: Delivery MUST include pushing the feature branch to the GitHub remote and opening or
  updating a PR; the next feature MUST NOT start until the current PR flow is complete.
- **GC-006**: Merge to the mainline branch MUST NOT occur without explicit user approval.
- **GC-007**: If subagents are used for implementation, they MUST default to `gpt-5.4` with `medium` reasoning unless the user explicitly overrides.

### Key Entities *(include if feature involves data)*

- **Cluster Lifecycle Record**: 表示一个受平台管理的集群对象，包含接入方式、基础设施类型、当前生命周期阶段、版本状态、节点池摘要和健康状态。
- **Cluster Driver**: 表示一种可用于创建或管理集群的驱动定义，包含驱动版本、支持能力、适用范围和兼容状态。
- **Cluster Template**: 表示可复用的集群创建模板，描述目标基础设施类型、依赖能力、校验规则和可用参数集合。
- **Capability Matrix Entry**: 表示某类集群或驱动在某个能力域下的支持状态、兼容结论和差异说明。
- **Upgrade Plan**: 表示一次集群版本升级安排，包含目标版本、前置检查、风险提示、执行状态和结果摘要。
- **Node Pool Profile**: 表示集群节点池的结构化信息，包含池角色、容量信息、伸缩目标和当前状态。
- **Lifecycle Audit Event**: 表示围绕导入、注册、创建、升级、扩缩节点、停用、退役和驱动变更产生的标准化审计记录。

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 平台管理员能够在 15 分钟内完成一个已有集群的导入或注册，并在统一生命周期中心看到明确的接入结果。
- **SC-002**: 对于受支持的基础设施类型，90% 的标准化集群创建请求能够在发起后 30 分钟内返回明确的成功或失败结论。
- **SC-003**: 在至少 20 个已纳管集群的环境中，90% 的集群升级计划查询和状态查看可在 30 秒内返回目标结果。
- **SC-004**: 对于需要停用或退役的集群，100% 的关键动作都能留存完整审计记录，并能按操作者、时间和结果检索。
- **SC-005**: 平台工程团队能够在 10 分钟内比较至少两类集群驱动在核心能力域上的兼容差异，并基于结果选择可用模板。
- **SC-006**: 权限验收中，100% 的未授权创建、导入、升级、扩缩节点和退役请求都被拦截，且不暴露超范围集群敏感细节。
- **SC-007**: 试点阶段至少 80% 的集群接入、创建、升级和退役流程可在平台内闭环完成，无需切换到分散入口。

## Assumptions

- `001-k8s-ops-platform` 已提供多集群资源范围、授权模型和审计基础能力，007 在其之上扩展集群生命周期管理。
- 首期目标用户为平台管理员和平台工程团队，不单独覆盖普通业务开发团队的自助集群申请流程。
- 首期以桌面 Web 管理体验为主，不要求移动端专属交互。
- 不同基础设施类型的驱动和模板存在能力差异，平台首期重点是统一展示和受控管理，而不是完全消除所有差异。
- 平台管理的目标环境具备最基本的集群接入或创建前提；若外部环境条件不满足，平台通过显式校验结果进行阻断或提示。
- 首期不覆盖平台级灾备、统一身份源整合、策略治理和应用市场，这些能力如需纳入应作为后续独立特性处理。
