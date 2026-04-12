# Feature Specification: 多集群 Kubernetes 可观测中心

**Feature Branch**: `002-observability-center`  
**Created**: 2026-04-11  
**Status**: In Progress  
**Input**: User description: "我要新增 002-observability-center，严格对标 Rancher 的可观测中心能力。该特性面向平台管理员、SRE、运维和值班人员，在多集群 Kubernetes 环境中提供统一的日志、事件、监控和告警能力。用户可以从集群、工作空间、项目、命名空间、工作负载、Pod、容器和时间范围等维度统一查看、筛选和检索日志，查看 Kubernetes 事件时间线，查看集群、节点、命名空间和工作负载的健康状态、容量使用、趋势变化和异常指标，并在同一资源上下文中联动查看日志、事件、指标和告警，用于故障发现、定位和处理。平台需要提供告警规则、告警级别、通知目标、静默窗口、告警确认、恢复状态和处理记录等能力，形成完整的告警治理闭环。首期只做统一可观测入口、问题定位闭环和告警处理闭环，不包含终端、批量操作、回滚、GitOps、Helm、策略治理、合规扫描、集群创建导入和灾备能力。该特性必须继承现有工作空间和项目级权限隔离，保证用户只能访问自己授权范围内的可观测数据。"

## 当前状态/执行说明（2026-04-11）

- 已完成 `/speckit.plan` 规划产物，`plan.md`、`research.md`、`data-model.md`、`contracts/openapi.yaml` 和 `quickstart.md` 已生成。
- 已完成 `/speckit.tasks` 任务拆解，`tasks.md` 已生成，按 Setup、Foundational、US1、US2、US3、Polish 组织。
- 已确认 `001` 的 PR 流已结束并合并到 `main`；`002` 的实施前数据库备份已于 2026-04-11 执行完成，产物为 `artifacts/002-observability-center/mysql-backup-20260411-214819.sql`，并已通过临时容器恢复抽样验证。
- 002 已进入实现阶段，当前先行执行 Governance + Setup + Foundational 首批任务（骨架与路由占位）。
- 002 开始编码前的前置数据保护要求已满足；后续若涉及新的高风险数据库调整，仍需再次执行备份。
- 已补充治理门槛证据：`artifacts/002-observability-center/branch-check.txt`、`artifacts/002-observability-center/mirror-and-remote-check.txt`。
- 当前状态更新为“可进入 `/speckit.implement` 实施阶段”，仍需遵守“提交中文 PR + 用户明确批准后再合并”。

## User Scenarios & Testing *(mandatory)*

### User Story 1 - 统一可观测入口与问题定位 (Priority: P1)

作为 SRE、运维和值班人员，我希望在同一个平台里围绕集群、工作空间、项目和具体资源统一查看日志、事件、监控和告警，这样我在故障发生时可以更快判断影响范围、定位异常对象并找到处理切入点。

**Why this priority**: 统一可观测入口是所有后续诊断和处置动作的前提。如果仍需要在多个工具之间切换，平台无法形成对标 Rancher 的核心观测体验。

**Independent Test**: 在至少两个已接入集群中制造一个工作负载异常后，授权用户能够在同一平台内按资源上下文查看关联日志、事件、指标和当前告警，并据此定位异常对象。

**Acceptance Scenarios**:

1. **Given** 某个工作负载出现异常，**When** 值班人员进入该资源的可观测视图，**Then** 平台应同时展示该资源关联的日志、事件、指标摘要和相关告警，并明确其所属集群、工作空间和命名空间。
2. **Given** 用户需要检索某段时间内的运行日志，**When** 按集群、工作空间、项目、命名空间、工作负载、Pod、容器、时间范围或关键字筛选，**Then** 平台应返回匹配结果并保留筛选上下文，便于继续下钻。
3. **Given** 某个集群短时不可达或观测数据暂时中断，**When** 用户查看相关资源，**Then** 平台应明确标记数据陈旧、缺失或不可获取，而不是误报为资源不存在或状态正常。

---

### User Story 2 - 告警治理与值班闭环 (Priority: P1)

作为平台管理员或值班负责人，我希望在平台内统一管理告警规则、通知目标、静默窗口和处置记录，使团队能够从告警触发到确认、抑制、恢复和复盘形成闭环。

**Why this priority**: 没有告警治理，平台只能被动展示数据，无法承担企业级值班和问题处理入口的角色。

**Independent Test**: 创建一条告警规则并触发一次异常后，授权用户能够看到告警生成、通知、确认、静默、恢复和处理记录的完整链路。

**Acceptance Scenarios**:

1. **Given** 管理员已经配置面向某个资源范围的告警规则和通知目标，**When** 目标对象满足触发条件，**Then** 平台应生成对应告警，展示级别、影响范围、触发时间和当前状态。
2. **Given** 某条告警正在处理中，**When** 授权用户对其执行确认、补充处理说明或设置静默窗口，**Then** 平台应记录处理人、处理时间、原因说明和生效状态，并在后续视图中可追踪。
3. **Given** 异常已经消除，**When** 告警恢复，**Then** 平台应更新告警状态、保留从触发到恢复的完整时间线，并支持后续检索和复盘。

---

### User Story 3 - 权限隔离下的可观测访问 (Priority: P2)

作为平台管理员，我希望可观测数据继续遵守工作空间和项目级权限隔离，让不同团队只能访问自己被授权范围内的日志、事件、指标和告警，从而避免跨团队数据泄露和错误处置。

**Why this priority**: 可观测数据往往包含业务运行信息和敏感上下文，如果没有范围隔离，会直接破坏平台现有的多租户治理边界。

**Independent Test**: 为两个不同工作空间分别授权后，两个用户登录平台时只能看到各自范围内的日志、事件、指标和告警；权限被回收后访问立即失效。

**Acceptance Scenarios**:

1. **Given** 两个工作空间分别绑定了不同项目和资源范围，**When** 不同用户进入可观测中心，**Then** 平台应只显示各自授权范围内的数据，不得暴露其他范围的资源线索或告警摘要。
2. **Given** 某个用户仅具备查看权限，**When** 访问告警治理入口，**Then** 平台应允许其查看授权范围内的告警和上下文，但不得修改规则、静默窗口或处理记录。
3. **Given** 用户的授权范围被缩减或撤销，**When** 其再次查看原有资源或告警，**Then** 平台应立即按新权限拒绝访问，并阻止继续查看相关观测数据。

### Edge Cases

- 当不同集群或命名空间中存在同名工作负载或 Pod 时，平台必须始终展示其完整资源上下文，避免误判目标对象。
- 当 Pod 已重建、日志源已轮转或历史数据部分缺失时，平台必须区分“无结果”“历史缺口”和“数据源不可用”三种状态。
- 当指标暂时中断、采样延迟或事件尚未到达时，平台必须提示数据时效性，避免把延迟数据解释为资源健康。
- 当同一告警在短时间内频繁触发和恢复时，平台必须清晰展示告警状态变化，避免值班人员误判为多个独立问题。
- 当静默窗口生效期间出现新的同类异常时，平台必须标明其被静默覆盖，而不是直接隐藏相关风险。
- 当用户在处置告警过程中权限被收回时，平台必须阻止其继续修改告警状态、静默窗口或处理记录。

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: 系统 MUST 提供统一的多集群可观测入口，支持围绕集群、工作空间、项目、命名空间和具体资源查看相关观测数据。
- **FR-002**: 系统 MUST 支持按集群、工作空间、项目、命名空间、工作负载、Pod、容器、时间范围和关键字筛选与检索日志。
- **FR-003**: 系统 MUST 在日志结果中展示日志来源上下文，至少包括所属集群、命名空间、资源对象、实例对象和发生时间。
- **FR-004**: 系统 MUST 支持用户从资源上下文直接进入关联日志视图，并继承当前资源范围和时间范围，减少重复筛选。
- **FR-005**: 系统 MUST 提供 Kubernetes 事件时间线视图，支持按集群、命名空间、资源对象、事件级别和时间范围查看事件。
- **FR-006**: 系统 MUST 展示集群、节点、命名空间和工作负载的健康状态、容量使用、趋势变化和异常指标。
- **FR-007**: 系统 MUST 在同一资源上下文中关联展示日志、事件、指标和当前相关告警，帮助用户进行问题定位。
- **FR-008**: 系统 MUST 明确标识观测数据的时间范围、最后更新时间和当前数据状态，区分实时、延迟、缺失和不可获取等情况。
- **FR-009**: 系统 MUST 支持用户在统一可观测入口中切换资源层级和时间范围，并保持筛选条件在同一分析流程内连续可见。
- **FR-010**: 系统 MUST 支持授权管理员创建、编辑、启用、停用和删除告警规则，并为规则定义适用范围和触发条件。
- **FR-011**: 系统 MUST 支持为告警规则配置告警级别、通知目标和静默窗口等治理属性。
- **FR-012**: 系统 MUST 在满足触发条件时生成告警事件，并为每条告警保留触发时间、影响范围、当前状态和恢复状态。
- **FR-013**: 系统 MUST 允许授权用户对告警执行确认、补充处理说明和记录处理过程。
- **FR-014**: 系统 MUST 支持创建、更新和取消静默窗口，并记录静默范围、生效时段、设置原因和操作人。
- **FR-015**: 系统 MUST 为每条告警展示从触发、通知、确认、静默到恢复的完整状态时间线和处理记录。
- **FR-016**: 系统 MUST 支持按告警状态、级别、集群、工作空间、项目、资源类型、时间范围和处理人筛选告警。
- **FR-017**: 系统 MUST 继承现有工作空间和项目级权限隔离，对日志、事件、指标、告警、规则和静默窗口统一执行授权校验。
- **FR-018**: 系统 MUST 在用户权限变化后立即按新权限限制其可见的观测数据和可执行的告警治理动作。
- **FR-019**: 系统 MUST 对告警规则变更、静默窗口变更、告警确认和关键观测访问行为生成可检索的审计记录。
- **FR-020**: 系统 MUST 为告警和观测视图中的异常状态提供可理解的说明，至少覆盖数据源中断、数据延迟、权限不足和对象已变化等情况。
- **FR-021**: 首期范围 MUST 聚焦统一可观测入口、问题定位闭环和告警处理闭环，不包含终端、批量操作、回滚、GitOps、Helm、策略治理、合规扫描、集群创建导入和灾备能力。

## Governance & Delivery Constraints *(mandatory)*

- **GC-001**: Feature work MUST occur on a dedicated feature branch; direct development on `master` or `main` is forbidden.
- **GC-002**: All user-facing communication, approval records, PR summaries, and delivery notes MUST be written in Chinese.
- **GC-003**: Any dependency or framework installation MUST document the China mirror or proxy configuration that will be used during implementation.
- **GC-004**: Before implementation begins, the feature specification or plan MUST record a database backup executed from container `mysql8` using `localhost:3306` and credentials `admin/123456`, or explicitly justify why the backup requirement is not applicable.
- **GC-005**: Delivery MUST include pushing the feature branch to the GitHub remote and opening or updating a PR; the next feature MUST NOT start until the current PR flow is complete.
- **GC-006**: Merge to the mainline branch MUST NOT occur without explicit user approval.
- **GC-007**: If subagents are used for implementation, they MUST use `gpt-5.3-codex`.

### Key Entities *(include if feature involves data)*

- **Observability Scope**: 一次观测视图的目标范围，包含集群、工作空间、项目、命名空间和资源对象等上下文，用于统一约束日志、事件、指标和告警的展示边界。
- **Log Query Session**: 一次日志筛选与检索行为，包含筛选条件、时间范围、关键字和结果上下文，用于支持连续分析和问题定位。
- **Event Timeline Item**: 某个资源或范围内发生的事件记录，描述事件级别、发生时间、关联对象和关键说明。
- **Metric Insight**: 面向集群、节点、命名空间或工作负载的健康、容量、趋势和异常信号摘要，用于帮助用户识别异常变化。
- **Alert Rule**: 用于定义何时触发告警及其治理属性的规则实体，包含适用范围、触发条件、告警级别和通知目标。
- **Alert Incident**: 一次具体告警实例，记录触发时间、影响范围、当前状态、恢复状态和处理进展。
- **Silence Window**: 对特定范围和时段生效的告警抑制规则，记录原因、生效时间、操作者和目标范围。
- **Handling Record**: 围绕告警确认、说明补充、静默设置和恢复复盘形成的处理记录，用于交接和值班复盘。

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 在至少 20 个已接入集群的环境中，授权用户能够在 3 分钟内从目标资源进入关联日志、事件、指标和告警视图并完成初步问题定位。
- **SC-002**: 在试点团队的值班演练中，90% 的告警可在 1 分钟内完成确认、静默或补充处理说明等首个处置动作。
- **SC-003**: 针对最近 24 小时范围内的常见日志查询，90% 的查询能够在 10 秒内向用户返回首批可阅读结果。
- **SC-004**: 针对最近 30 天范围内的告警检索，90% 的查询能够在 15 秒内返回符合筛选条件的结果集。
- **SC-005**: 在试点期内，至少 80% 的常见故障定位场景可仅通过该平台完成，不再依赖多个分散工具来回切换。
- **SC-006**: 在权限验收中，100% 的跨工作空间和跨项目未授权观测访问请求都被拦截，且不会暴露目标对象的详细运行数据。

## Assumptions

- 现有 `001-k8s-ops-platform` 已提供基础的多集群接入、资源范围模型、工作空间和项目级权限隔离能力，本特性在其之上扩展观测与告警能力。
- 首期目标用户为平台管理员、SRE、运维和值班人员，不单独覆盖移动端专属体验。
- 各被接入集群已经具备可供平台读取的日志、事件和指标来源；若部分来源缺失，平台仍需明确提示数据缺口而不是静默失败。
- 通知目标的具体外部渠道可以因组织而异，但首期重点是统一管理目标对象、告警级别、静默和处理闭环，而不是替代外部协同系统的全部能力。
- 首期不扩展到终端登录、批量运维、回滚、GitOps、Helm、策略治理、合规扫描、集群生命周期和灾备恢复等其他业务域能力。
- 该特性的实现和交付仍需遵守“完成当前 PR 流后再启动下一个 feature 实施”的仓库治理要求。
