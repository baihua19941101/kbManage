# Feature Specification: 身份与多租户治理中心

**Feature Branch**: `009-identity-tenancy`  
**Created**: 2026-04-19  
**Status**: Completed (Implementation Ready)  
**Input**: User description: "我要新增 009-identity-and-tenancy，严格对标 Rancher 和企业级平台的身份与多租户治理能力，覆盖 SSO、OIDC、LDAP、组织模型和细粒度 RBAC。面向平台管理员、安全管理员和组织治理负责人，在多集群 Kubernetes 平台中提供统一的身份接入、组织建模和授权治理中心。用户需要能够接入外部身份源，管理组织、团队、用户组和项目关系，定义平台级、组织级、工作空间级、项目级和资源级角色与权限，控制权限继承、授权范围、委派管理、临时授权和访问回收，并在统一视图中查看用户来源、角色分布、权限边界和访问风险。平台还需要支持本地账号与外部身份并存、登录方式切换、会话治理和权限变更追踪，确保扩展到企业级组织后仍能保持清晰的租户隔离和访问治理。首期范围聚焦身份接入、组织模型和细粒度 RBAC，不包含准入策略、合规扫描、应用发布和灾备恢复。"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - 接入外部身份源并统一登录 (Priority: P1)

作为平台管理员或安全管理员，我希望把企业已有身份源接入平台，并支持本地账号与外部身份并存，这样用户可以通过统一登录入口完成身份认证，而平台仍保留本地兜底管理能力。

**Why this priority**: 如果没有稳定的身份接入和登录控制，后续组织建模、角色授权和租户隔离都无法建立在可信身份基础上。

**Independent Test**: 接入一个外部身份源并保留本地管理员账号后，用户可以选择登录方式完成登录，管理员可以查看身份源状态、用户来源和会话状态。

**Acceptance Scenarios**:

1. **Given** 平台尚未接入外部身份源，**When** 管理员新增一个身份源配置并启用，**Then** 平台应记录其状态、支持的登录方式和可用范围。
2. **Given** 平台同时存在本地账号和外部身份源，**When** 用户访问登录入口，**Then** 平台应允许按规则选择登录方式，并明确展示用户来源。
3. **Given** 某个身份源不可用或被停用，**When** 用户尝试通过该方式登录，**Then** 平台应阻止登录并保留本地管理员兜底访问能力。

---

### User Story 2 - 建立组织与租户关系模型 (Priority: P1)

作为组织治理负责人，我希望在平台中维护组织、团队、用户组、工作空间和项目之间的归属与映射关系，这样企业级组织扩展后仍能保持清晰的租户边界和责任归属。

**Why this priority**: 身份接入只能解决“谁来登录”，组织模型才能解决“谁属于哪个租户、团队和业务边界”。

**Independent Test**: 创建一个组织、多个团队和用户组，并把它们映射到工作空间和项目后，管理员能够在统一视图中看到继承关系、归属关系和租户边界。

**Acceptance Scenarios**:

1. **Given** 平台已接入可用身份源，**When** 治理负责人创建组织、团队和用户组层级，**Then** 平台应记录成员归属、上下级关系和管理责任。
2. **Given** 组织模型已建立，**When** 管理员把团队或用户组映射到工作空间和项目，**Then** 平台应展示租户边界和可继承的授权范围。
3. **Given** 成员跨多个组织或项目存在访问关系，**When** 治理负责人查看成员详情，**Then** 平台应清晰展示其来源、归属、有效授权范围和边界冲突提示。

---

### User Story 3 - 管理细粒度 RBAC、委派和回收 (Priority: P2)

作为安全管理员，我希望定义多层级角色、控制权限继承与委派、配置临时授权并跟踪回收与变更记录，这样平台在复杂组织环境下仍能保持可控授权和可审计访问治理。

**Why this priority**: 企业级平台的关键挑战不是“有没有权限”，而是“权限是否继承合理、是否可回收、是否存在高风险扩散”。

**Independent Test**: 为用户或用户组分配平台级到资源级角色后，管理员能够查看权限边界、临时授权到期状态、委派链路和访问风险，并能执行回收。

**Acceptance Scenarios**:

1. **Given** 平台已存在多层级组织模型，**When** 安全管理员定义平台级、组织级、工作空间级、项目级和资源级角色，**Then** 平台应展示角色作用范围、可继承关系和限制条件。
2. **Given** 某个管理员被授予委派管理权限，**When** 其在授权范围内为团队发放临时授权，**Then** 平台应记录委派链路、生效范围、到期时间和回收条件。
3. **Given** 用户权限发生变更、到期或被手动回收，**When** 安全管理员查看统一视图，**Then** 平台应展示角色分布、权限边界变化、会话影响和访问风险提示。

### Edge Cases

- 当同一用户同时来自本地账号和外部身份源，平台必须明确主身份来源与冲突处理结果。
- 当外部身份源中的组织结构与平台已有组织模型不一致时，平台必须提示映射冲突并阻止产生模糊租户边界。
- 当角色继承链导致权限超出原始授权边界时，平台必须阻止生效并给出冲突说明。
- 当临时授权到期但用户仍存在活跃会话时，平台必须定义会话影响并确保后续访问不再扩大权限。
- 当用户被移出组织、团队或用户组时，平台必须同步评估需回收的授权和残留访问风险。
- 当身份源停用、切换或连接异常时，平台必须保证至少保留安全管理员的应急访问路径。

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: 系统 MUST 提供统一的身份与多租户治理中心，用于管理身份源接入、组织模型和细粒度 RBAC。
- **FR-002**: 系统 MUST 支持接入外部身份源，并管理其启用状态、登录入口状态和可用性状态。
- **FR-003**: 系统 MUST 支持本地账号与外部身份并存，并允许管理员控制默认登录方式与可用登录方式。
- **FR-004**: 系统 MUST 在统一视图中展示用户的身份来源、当前登录方式和账号状态。
- **FR-005**: 系统 MUST 支持组织、团队、用户组、工作空间和项目之间的关系建模。
- **FR-006**: 系统 MUST 支持把组织、团队和用户组映射到工作空间和项目范围。
- **FR-007**: 系统 MUST 支持平台级、组织级、工作空间级、项目级和资源级角色定义。
- **FR-008**: 系统 MUST 为每种角色明确记录适用范围、可授权对象和权限边界。
- **FR-009**: 系统 MUST 支持权限继承规则，并明确展示继承来源、继承层级和最终生效权限。
- **FR-010**: 系统 MUST 支持授权范围控制，确保用户只能被授予授权人可管理范围内的权限。
- **FR-011**: 系统 MUST 支持委派管理，使指定管理员可以在受限范围内继续分配角色。
- **FR-012**: 系统 MUST 支持临时授权，并记录生效时间、到期时间、授权原因和回收状态。
- **FR-013**: 系统 MUST 支持访问回收，并在身份变更、组织变更、到期或管理员操作后收回对应权限。
- **FR-014**: 系统 MUST 支持会话治理，包括查看会话状态、识别高风险会话并在权限变更后处理受影响访问。
- **FR-015**: 系统 MUST 支持查看用户、团队、用户组和角色的分布情况与授权覆盖范围。
- **FR-016**: 系统 MUST 支持查看权限边界、越权风险、继承扩散风险和残留访问风险提示。
- **FR-017**: 系统 MUST 支持权限变更追踪，记录角色授予、委派、回收、到期和来源变更。
- **FR-018**: 系统 MUST 对身份源管理、组织变更、角色授予、临时授权、委派和回收动作生成可检索审计记录。
- **FR-019**: 系统 MUST 支持按用户来源、组织、团队、用户组、角色层级、授权范围和风险状态筛选治理对象。
- **FR-020**: 系统 MUST 保持多租户隔离，防止一个组织或项目的授权扩散到未授权租户范围。
- **FR-021**: 系统 MUST 在身份源异常、组织映射冲突、角色冲突或越权授权时阻止生效并返回明确原因。
- **FR-022**: 系统 MUST 支持统一查看单个用户的来源、组织归属、角色分布、有效权限、会话状态和风险摘要。
- **FR-023**: 首期范围 MUST 聚焦身份接入、组织模型和细粒度 RBAC。
- **FR-024**: 首期范围 MUST 不包含准入策略、合规扫描、应用发布和灾备恢复。
- **FR-025**: 首期范围 MUST 预留后续扩展空间，以支持更多身份源、更多租户层级和更细粒度的授权对象。

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
- **GC-007**: If subagents are used for implementation, they SHOULD use `gpt-5.4` with `medium` reasoning unless the user explicitly specifies another model.

### Key Entities *(include if feature involves data)*

- **Identity Source**: 表示一个外部或本地身份接入入口，包含来源类型、可用状态、登录方式和接入范围。
- **Organization Unit**: 表示组织、团队或用户组等层级治理对象，包含成员关系、上下级关系和映射范围。
- **Tenant Scope Mapping**: 表示组织对象与工作空间、项目或资源之间的租户边界映射关系。
- **Role Definition**: 表示某一级别的角色模板，包含适用层级、权限集合和边界限制。
- **Role Assignment**: 表示某个用户、团队或用户组在特定范围内获得的授权记录，包含来源、继承关系和有效期。
- **Delegation Grant**: 表示授权委派关系，包含委派人、被委派人、允许操作范围和限制条件。
- **Session Record**: 表示一次登录会话，包含身份来源、会话状态、风险状态和受权限变更影响情况。
- **Access Risk Snapshot**: 表示针对用户或租户范围计算出的访问风险摘要，包含越权、残留权限和高风险继承提示。

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 平台管理员能够在 30 分钟内完成至少一种外部身份源接入，并保留本地管理员兜底登录能力。
- **SC-002**: 组织治理负责人能够在 20 分钟内创建一个包含组织、团队、用户组、工作空间和项目映射关系的完整租户模型。
- **SC-003**: 安全管理员能够在 15 分钟内完成从平台级到项目级的一组角色定义与授权配置，并查看清晰的权限边界。
- **SC-004**: 100% 的角色授予、委派、临时授权、回收和身份源变更动作都能留存可检索审计记录。
- **SC-005**: 在试点组织中，90% 的用户权限查询、角色边界查看和风险摘要查看可在 30 秒内返回目标结果。
- **SC-006**: 权限验收中，100% 的超范围授权和跨租户越权访问请求都被拦截，并返回明确原因。
- **SC-007**: 至少 80% 的组织治理场景能够通过统一视图识别用户来源、角色分布、权限边界和高风险访问状态。

## Assumptions

- `001-k8s-ops-platform` 已提供基础用户、会话、工作空间、项目和审计能力，009 在其之上扩展身份接入与租户治理。
- 首期目标用户为平台管理员、安全管理员和组织治理负责人，不覆盖普通开发者自助申请身份接入或角色开通流程。
- 首期以桌面 Web 管理体验为主，不要求移动端专属交互。
- 企业已有外部身份目录和组织信息可作为接入来源，但平台首期只负责接入、映射和治理，不负责改造外部身份系统本身。
- 首期重点是统一身份接入、组织建模和细粒度 RBAC 闭环，而不是覆盖更广义的安全策略、准入控制或发布治理能力。
- 首期不包含准入策略、合规扫描、应用发布和灾备恢复，这些能力如需纳入应作为后续独立特性处理。
