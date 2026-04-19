# Research: 身份与多租户治理中心

> 研究输入基于 009 规格、001-008 既有能力边界以及企业级平台身份与租户治理通用语义，目标是在不扩展到准入策略、合规扫描、应用发布和灾备恢复的前提下，形成可实施的首期方案。

## Decision 1: 009 采用独立 `identitytenancy` 业务域，而不是继续扩展基础 `auth`

- Decision: 后端新增 `service/identitytenancy`、`integration/identity` 与对应 handler/router/repository，前端新增 `features/identity-tenancy`；复用 001 的用户、会话和审计底座，但身份源、组织模型和多层级授权语义独立承载。
- Rationale: 基础 `auth` 更偏向登录会话与用户自身份认证，009 的核心对象是身份源、组织层级、角色定义、授权边界和委派关系，职责边界已明显超出基础认证。
- Alternatives considered:
  - 继续扩展 `auth`：会把基础登录和企业级租户治理耦合，后续维护成本过高。
  - 混入 005：策略治理关注准入与策略执行，不适合作为身份和授权中心。

## Decision 2: 身份接入采用“本地账号 + 外部身份源并存”的双目录模型

- Decision: 平台保留本地账号作为管理兜底，同时接入外部身份源；统一视图展示用户来源、登录方式和有效状态。
- Rationale: 企业级平台必须避免单一外部身份源异常时造成平台完全失管。
- Alternatives considered:
  - 全量切换到外部身份源：无法保证应急兜底访问。
  - 只保留本地账号：无法满足企业级统一身份接入诉求。

## Decision 3: 外部身份源首期统一抽象为 IdentitySource

- Decision: 无论是 SSO、OIDC 还是 LDAP，都统一建模为 `IdentitySource`，由来源类型区分配置和能力差异。
- Rationale: 这样可以在不暴露实现细节的情况下，为后续新增更多身份源保留一致的治理入口。
- Alternatives considered:
  - 为每类身份源单独建模：管理入口和审计动作会重复膨胀。
  - 只做单一身份源：不符合企业级扩展预期。

## Decision 4: 组织、团队、用户组采用统一组织单元模型

- Decision: 组织、团队和用户组统一归入 `OrganizationUnit`，通过单元类型和上下级关系表达不同治理层级。
- Rationale: 企业组织树通常需要同时表达管理归属和授权归属，统一模型更利于展示和映射。
- Alternatives considered:
  - 各自独立表述：关系跨表过多，展示和审计不清晰。
  - 仅保留组织和项目两级：无法表达团队和用户组治理。

## Decision 5: 租户边界采用“组织对象到工作空间/项目/资源”的映射模型

- Decision: 平台通过 `TenantScopeMapping` 记录组织对象与工作空间、项目或资源之间的边界映射，而不是把租户边界直接硬编码到角色中。
- Rationale: 租户边界和角色职责是两个不同维度，拆分后更易于复用和风险评估。
- Alternatives considered:
  - 把边界写死在角色定义中：角色复用性差。
  - 只依赖工作空间和项目归属：无法表达资源级或跨层级边界。

## Decision 6: 角色定义与授权分配分离建模

- Decision: `RoleDefinition` 表示角色模板，`RoleAssignment` 表示具体用户、团队或用户组在某一范围内获得的授权。
- Rationale: 企业级 RBAC 需要既能定义标准角色，又能追踪具体授权来源、有效期和继承链路。
- Alternatives considered:
  - 只保留角色模板：无法落地谁在何处被授权。
  - 把角色内容直接嵌入授权记录：难以统一治理角色体系。

## Decision 7: 权限继承、委派和临时授权必须显式建模

- Decision: 继承关系通过授权链路可见，委派单独建模为 `DelegationGrant`，临时授权在授权分配中保留有效期和回收状态。
- Rationale: 企业级访问治理的关键不是静态角色，而是权限如何扩散、何时到期、由谁继续委派。
- Alternatives considered:
  - 只做静态 RBAC：无法满足委派和回收治理要求。
  - 把委派和临时授权当作普通授权备注：难以审计和自动回收。

## Decision 8: 会话治理独立建模，不把会话状态混在用户状态中

- Decision: 登录会话采用 `SessionRecord` 独立承载，记录身份来源、风险状态、会话有效性和权限变更影响。
- Rationale: 用户账号状态和单次会话状态不是同一维度；权限回收后的会话处置也需要单独追踪。
- Alternatives considered:
  - 只记录用户状态：无法表示多会话和会话风险。
  - 完全依赖现有 token 机制：缺少企业级治理视图。

## Decision 9: 风险视图采用快照式聚合，而不是实时逐条解释

- Decision: 平台对越权、残留权限、高风险继承和异常会话采用 `AccessRiskSnapshot` 聚合展示，在统一视图中呈现风险摘要。
- Rationale: 治理负责人需要快速识别问题用户和问题租户，而不是手动拼接多张列表。
- Alternatives considered:
  - 只保留原始授权记录：难以快速定位高风险对象。
  - 完全实时计算不落库：审计和复盘成本高。

## Decision 10: 009 的审计事件单独归类为身份与租户治理域

- Decision: 身份源管理、组织变更、角色授予、委派、临时授权、回收和会话治理动作单独归入身份治理域审计。
- Rationale: 这类动作对象以身份、组织和授权关系为中心，与 005 的策略动作和 008 的备份恢复动作明显不同。
- Alternatives considered:
  - 混入通用审计流不做域区分：查询维度不清晰，不利于安全审计。
  - 只保留粗粒度登录日志：无法覆盖组织与授权治理动作。
