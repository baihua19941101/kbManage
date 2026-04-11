# Feature Specification: 多集群 Kubernetes 可视化管理平台

**Feature Branch**: `001-k8s-ops-platform-followup`  
**Created**: 2026-04-09  
**Status**: Follow-up 执行中（可执行）  
**Input**: User description: "一个基于 Web 的 Kubernetes 可视化管理平台，用于多集群下的资源管理、运维操作与审计追踪,整体功能设计可以参考Rancher"

## 当前状态/执行说明（2026-04-10）

- 本特性已从 Draft/规划态切换为 Follow-up 执行阶段，按现有规格继续落地，不变更主体范围。
- 当前执行分支：`001-k8s-ops-platform-followup`。
- Follow-up 期间仅补充状态与执行记录，需求、场景与验收标准维持原文定义。

## User Scenarios & Testing *(mandatory)*

### User Story 1 - 多集群统一接入与资源总览 (Priority: P1)

作为平台管理员，我希望在一个 Web 平台中接入并管理多个 Kubernetes 集群，统一查看集群健康、容量、命名空间和核心资源状态，从而减少在多个控制面之间来回切换的成本。

**Why this priority**: 多集群接入与统一视图是后续授权、运维操作和审计追踪的前提，没有这个能力，平台无法形成可用的管理入口。

**Independent Test**: 接入至少 2 个集群后，管理员能够在统一视图中按集群、命名空间和资源类型浏览资源，并识别异常集群与异常资源。

**Acceptance Scenarios**:

1. **Given** 管理员持有可用的集群接入信息，**When** 提交一个新集群接入请求，**Then** 平台应将该集群加入管理范围，并展示连接状态、健康摘要和基础资源统计。
2. **Given** 平台已接入多个集群，**When** 管理员按集群、命名空间、资源类型或关键字筛选资源，**Then** 平台应返回匹配结果并显示每个资源的当前状态与所属范围。
3. **Given** 某个已接入集群暂时不可达，**When** 管理员查看总览页面，**Then** 平台应明确标记该集群状态异常并保留最近一次成功同步的状态信息，不影响其他集群的浏览。

---

### User Story 2 - 团队授权与资源隔离 (Priority: P1)

作为平台管理员，我希望基于团队、项目或工作空间为不同角色授予差异化访问权限，使运维人员、开发团队和审计人员只看到并操作自己被授权的集群和资源。

**Why this priority**: 多集群平台如果没有权限隔离，会直接带来误操作和越权风险，因此必须与统一接入能力一起成为首批交付内容。

**Independent Test**: 创建至少两个独立的工作空间并分别授权给不同角色后，每个用户只能看到并操作自己授权范围内的资源。

**Acceptance Scenarios**:

1. **Given** 平台管理员已创建工作空间并绑定目标集群或命名空间，**When** 为某个团队成员授予运维权限，**Then** 该成员登录后只能看到授权范围内的资源与允许执行的操作。
2. **Given** 某个用户仅具备审计或只读角色，**When** 访问资源详情或操作入口，**Then** 平台应允许其查看被授权的数据，但不得显示或执行超出权限范围的变更操作。
3. **Given** 用户的授权范围被缩减或撤销，**When** 其再次访问原有资源范围，**Then** 平台应立即按新权限生效并拒绝越权访问。

---

### User Story 3 - 受控运维操作执行 (Priority: P2)

作为运维人员，我希望直接在平台内对被授权的资源执行常见运维操作，并获得清晰的风险提示、执行进度和结果反馈，从而提升处理故障和日常变更的效率。

**Why this priority**: 在具备统一视图和权限控制后，平台的核心价值在于把高频运维动作收敛到单一入口，减少切换工具与误操作成本。

**Independent Test**: 在授权范围内选择一个工作负载和一个节点，执行常见运维动作后，平台能展示确认、执行过程与最终结果。

**Acceptance Scenarios**:

1. **Given** 运维人员拥有目标工作负载的操作权限，**When** 发起扩缩容、重启或配置变更等常见运维动作并确认执行，**Then** 平台应记录请求、展示执行过程并返回明确结果。
2. **Given** 运维人员拥有节点维护权限，**When** 发起节点维护类操作，**Then** 平台应在执行前提示潜在影响，并在执行后反馈完成状态或失败原因。
3. **Given** 某项运维操作因资源冲突、权限不足或集群异常而失败，**When** 运维人员查看操作结果，**Then** 平台应展示可理解的失败原因并保留失败记录用于后续追踪。

---

### User Story 4 - 审计追踪与操作复盘 (Priority: P3)

作为审计人员或平台负责人，我希望按时间、操作者、集群、资源和结果检索平台内的访问与运维记录，以便进行安全追踪、责任界定和事后复盘。

**Why this priority**: 审计是平台可信运营的闭环能力，能够将多集群管理从“可操作”提升为“可追责、可复盘”。

**Independent Test**: 在平台产生多条访问和运维记录后，审计人员能筛选出指定时间段、指定操作者和指定资源的完整记录。

**Acceptance Scenarios**:

1. **Given** 平台已记录多集群下的访问与操作行为，**When** 审计人员按时间范围、操作者、集群或结果检索，**Then** 平台应返回匹配记录并支持查看每条记录的详细上下文。
2. **Given** 存在高风险或失败操作，**When** 审计人员打开记录详情，**Then** 平台应展示操作者、目标对象、操作类型、发生时间、执行结果和相关说明。
3. **Given** 审计人员需要形成阶段性审计报告，**When** 选择检索结果导出，**Then** 平台应输出结构化记录，便于审计留档与复盘。

### Edge Cases

- 当多个集群中存在同名资源时，平台必须始终展示其所属集群、命名空间或工作空间，避免误判目标对象。
- 当集群在浏览过程中短暂失联时，平台必须区分“当前不可达”和“资源不存在”两种状态。
- 当用户在会话期间权限被回收时，平台必须阻止其继续执行已不再被授权的操作。
- 当同一资源被多个操作者同时修改时，平台必须给出冲突提示并明确最终结果。
- 当资源在平台外被删除或变更后，平台必须在后续操作中提示状态已变化，避免对过期数据执行动作。
- 当审计查询跨越大量记录时，平台必须保持筛选条件清晰可见，避免返回不可理解或难以定位的结果集。

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: 系统 MUST 支持平台管理员接入、启用、停用和移除多个 Kubernetes 集群，并维护每个集群的连接状态与基本元数据。
- **FR-002**: 系统 MUST 为每个已接入集群展示统一的健康摘要、容量概览、节点状态和核心资源统计。
- **FR-003**: 系统 MUST 提供跨集群的统一资源浏览能力，支持按集群、工作空间、命名空间、资源类型、标签或关键字筛选。
- **FR-004**: 系统 MUST 支持平台管理员创建逻辑隔离单元，用于组织集群、命名空间和团队资源范围。
- **FR-005**: 系统 MUST 支持平台管理员为不同用户或用户组分配角色，并将角色约束到明确的资源范围。
- **FR-006**: 系统 MUST 在每一次资源查看、变更和运维动作中强制执行权限校验，禁止用户访问或操作未授权资源。
- **FR-007**: 系统 MUST 支持授权用户查看常见 Kubernetes 资源的详情、状态、关联关系和最近变更记录。
- **FR-007A**: 首期资源 Kind 范围 MUST 固定为 `Deployment`、`StatefulSet`、`DaemonSet`、`Pod`、`Service`、`Ingress`、`Node`、`Namespace`，超出范围的 Kind 不纳入首期资源索引与列表筛选。
- **FR-008**: 系统 MUST 支持授权用户在平台内发起常见资源管理动作，包括创建、编辑、删除、扩缩容、重启和维护类操作。
- **FR-009**: 系统 MUST 在执行高影响操作前提供明确的风险提示和二次确认，至少覆盖删除、节点维护和跨集群范围变更等动作。
- **FR-009A**: 首期高风险操作流程 MUST 采用“二次确认即执行”，不引入他人审批流或多级审批节点。
- **FR-010**: 系统 MUST 在运维操作执行期间向用户展示当前进度、最终结果及失败原因。
- **FR-011**: 系统 MUST 在集群不可达、权限不足、资源冲突或资源状态过期时向用户返回可理解的错误说明和后续建议。
- **FR-012**: 系统 MUST 对登录、登出、权限变更、集群接入、资源变更和运维动作生成可检索的审计记录。
- **FR-013**: 系统 MUST 为每条审计记录保存操作者、角色、目标资源、所属集群、动作类型、发生时间、执行结果和必要说明。
- **FR-014**: 系统 MUST 支持按时间范围、操作者、角色、集群、工作空间、资源类型、动作类型和结果筛选审计记录。
- **FR-015**: 系统 MUST 支持导出筛选后的审计记录，用于外部留档、审查和复盘。
- **FR-015A**: 首期审计导出格式 MUST 仅支持 `CSV`；导出内容 MUST 对敏感字段执行脱敏（至少包含访问凭据、令牌、密钥、密码、手机号、邮箱）。
- **FR-016**: 系统 MUST 保证不同工作空间、团队和角色之间的资源可见性与操作权限相互隔离。
- **FR-017**: 系统 MUST 在单个资源详情中展示其所属集群、命名空间、工作空间、当前状态以及最近一次相关操作结果。
- **FR-018**: 系统 MUST 将审计记录保留不少于 180 天，除非更长保留周期由组织治理要求覆盖。
- **FR-019**: 首批角色矩阵 MUST 固定为 `platform-admin`、`ops-operator`、`audit-reader`、`readonly`，角色语义和权限边界以该矩阵为准。

## Governance & Delivery Constraints *(mandatory)*

- **GC-001**: Feature work MUST occur on a dedicated feature branch; direct development on `master` or `main` is forbidden.
- **GC-002**: All user-facing communication, approval records, PR summaries, and delivery notes MUST be written in Chinese.
- **GC-003**: Any dependency or framework installation MUST document the China mirror or proxy configuration that will be used during implementation.
- **GC-004**: Before implementation begins, the feature specification or plan MUST record a database backup executed from container `mysql8` using `localhost:3306` and credentials `admin/123456`, or explicitly justify why the backup requirement is not applicable.
- **GC-005**: Delivery MUST include pushing the feature branch to the GitHub remote and opening or updating a PR; the next feature MUST NOT start until the current PR flow is complete.
- **GC-006**: Merge to the mainline branch MUST NOT occur without explicit user approval.
- **GC-007**: If subagents are used for implementation, they MUST use `gpt-5.3-codex`.

### Key Entities *(include if feature involves data)*

- **Cluster**: 一个被平台接入的 Kubernetes 集群，包含标识信息、连接状态、健康摘要、容量数据和所属范围。
- **Workspace**: 用于隔离团队、项目或环境的逻辑单元，定义可见资源范围和授权边界。
- **Resource Item**: 平台展示或操作的具体资源对象，包含资源类型、名称、所属集群、命名空间、状态和关联关系。
- **Role Binding**: 用户或用户组与角色、资源范围之间的授权关系，用于限制可见性和可执行动作。
- **Operation Record**: 一次资源管理或运维动作的请求与执行结果，包含发起人、目标对象、风险级别、进度和最终状态。
- **Audit Event**: 用于审计追踪的标准化记录，覆盖访问、授权变更、集群接入、资源操作和结果信息。

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 平台管理员能够在单次管理周期内接入并验证至少 20 个集群，且每个集群在接入后 10 分钟内出现在统一总览中。
- **SC-002**: 在至少 10,000 个受管资源的环境下，授权用户能够在 1 分钟内定位目标集群资源并完成筛选。
- **SC-003**: 95% 的常见运维操作在用户提交后 2 分钟内返回明确的成功或失败结果。
- **SC-004**: 审计人员针对 90 天时间范围执行检索时，90% 的查询能够在 30 秒内返回目标结果集。
- **SC-005**: 试点团队至少 90% 的日常多集群资源管理与运维任务可直接在该平台完成，无需切换到额外管理界面。
- **SC-006**: 性能验收证据 MUST 包含“测试环境压测报告 + 可复现实验脚本”，并可在相同测试环境复核关键指标。

## Assumptions

- 首个版本面向平台管理员、运维人员、项目管理员、审计人员和只读查看者等内部角色。
- 首批角色矩阵固定为 `platform-admin`、`ops-operator`、`audit-reader`、`readonly`，后续扩展角色需作为新版本需求处理。
- 首个版本聚焦通用 Kubernetes 集群管理能力，不包含应用市场、计费结算或云厂商专属商业能力。
- 首批资源 Kind 固定为 `Deployment`、`StatefulSet`、`DaemonSet`、`Pod`、`Service`、`Ingress`、`Node`、`Namespace`。
- 用户身份来源与组织关系已由现有企业身份体系提供，本功能重点解决资源管理、操作控制与审计追踪。
- 平台默认管理的是组织可合法接入并持续连通的 Kubernetes 集群。
- 移动端专属体验不在首个版本范围内，首个版本以桌面浏览器中的 Web 管理体验为主。
- 审计记录的保留周期以 180 天为默认下限，若组织治理有更高要求，后续版本可扩展。
- 首期高风险动作不引入他人审批流，统一采用二次确认后立即执行。
