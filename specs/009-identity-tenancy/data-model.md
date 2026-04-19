# Data Model: 身份与多租户治理中心

## 1. IdentitySource

### Purpose

统一表示平台接入的本地或外部身份源，用于控制登录方式、来源状态和可用范围。

### Fields

- `id`: 主键标识
- `name`: 身份源名称
- `sourceType`: `local | sso | oidc | ldap`
- `status`: `draft | active | disabled | unavailable`
- `loginMode`: `exclusive | optional | fallback`
- `scopeMode`: `platform | selected-organizations`
- `syncState`: `idle | syncing | succeeded | failed`
- `ownerUserId`: 责任人
- `lastCheckedAt`: 最近可用性检查时间
- `createdAt`: 创建时间
- `updatedAt`: 更新时间

### Relationships

- 一对多关联 `IdentityAccount`
- 一对多关联 `SessionRecord`
- 一对多关联 `IdentityAuditEvent`

### Validation Rules

- 同名身份源不可重复
- `sourceType=local` 时必须保留可用兜底状态
- `status=disabled` 时不可作为新的登录入口

## 2. IdentityAccount

### Purpose

表示用户在某个身份源中的身份映射，用于区分用户来源和登录方式。

### Fields

- `id`: 主键标识
- `userId`: 平台用户
- `identitySourceId`: 来源身份源
- `externalRef`: 外部目录引用
- `principalType`: `user | service-account`
- `status`: `active | disabled | locked | orphaned`
- `lastLoginAt`: 最近登录时间
- `createdAt`: 创建时间
- `updatedAt`: 更新时间

### Relationships

- 多对一关联 `IdentitySource`
- 多对一关联平台用户

### Validation Rules

- 同一身份源下 `externalRef` 必须唯一
- `status=orphaned` 时必须标记为风险对象

## 3. OrganizationUnit

### Purpose

表示组织、团队或用户组等治理单元，用于维护成员归属和组织层级。

### Fields

- `id`: 主键标识
- `unitType`: `organization | team | user-group`
- `name`: 单元名称
- `description`: 单元说明
- `parentUnitId`: 上级单元
- `identitySourceId`: 来源身份源，可为空
- `ownerUserId`: 管理责任人
- `status`: `active | disabled | archived`
- `createdAt`: 创建时间
- `updatedAt`: 更新时间

### Relationships

- 自关联形成树状层级
- 一对多关联 `OrganizationMembership`
- 一对多关联 `TenantScopeMapping`

### Validation Rules

- 同层级下名称必须唯一
- 禁止形成循环上下级关系

## 4. OrganizationMembership

### Purpose

表示用户或用户组与组织单元之间的成员关系。

### Fields

- `id`: 主键标识
- `unitId`: 组织单元
- `memberType`: `user | user-group`
- `memberRef`: 成员引用
- `membershipRole`: `owner | maintainer | member | observer`
- `status`: `active | suspended | removed`
- `joinedAt`: 加入时间
- `updatedAt`: 更新时间

### Relationships

- 多对一关联 `OrganizationUnit`

### Validation Rules

- 同一成员在同一组织单元中只能存在一条有效成员关系
- `status=removed` 时必须触发回收评估

## 5. TenantScopeMapping

### Purpose

表示组织对象与工作空间、项目或资源之间的租户边界映射关系。

### Fields

- `id`: 主键标识
- `unitId`: 组织单元
- `scopeType`: `workspace | project | resource`
- `scopeRef`: 目标范围引用
- `inheritanceMode`: `direct | inherited | restricted`
- `status`: `active | pending | revoked`
- `createdBy`: 创建人
- `createdAt`: 创建时间
- `updatedAt`: 更新时间

### Relationships

- 多对一关联 `OrganizationUnit`
- 一对多关联 `RoleAssignment`

### Validation Rules

- 同一组织单元与同一范围映射不可重复
- `scopeType=resource` 时必须明确上层租户边界

## 6. RoleDefinition

### Purpose

表示平台级到资源级的标准角色模板。

### Fields

- `id`: 主键标识
- `name`: 角色名称
- `roleLevel`: `platform | organization | workspace | project | resource`
- `description`: 角色说明
- `permissionSummary`: 权限摘要
- `inheritancePolicy`: `none | upward-blocked | downward-allowed | bounded`
- `delegable`: 是否允许委派
- `status`: `draft | active | deprecated | disabled`
- `createdBy`: 创建人
- `createdAt`: 创建时间
- `updatedAt`: 更新时间

### Relationships

- 一对多关联 `RoleAssignment`

### Validation Rules

- 同一层级下角色名称必须唯一
- `delegable=false` 时不得出现在委派授权链路中

## 7. RoleAssignment

### Purpose

表示用户、团队或用户组在某个租户范围内获得的具体授权记录。

### Fields

- `id`: 主键标识
- `subjectType`: `user | team | user-group`
- `subjectRef`: 授权对象引用
- `roleDefinitionId`: 角色模板
- `scopeType`: `platform | organization | workspace | project | resource`
- `scopeRef`: 作用范围引用
- `sourceType`: `direct | inherited | delegated | temporary`
- `delegationGrantId`: 委派来源，可为空
- `validFrom`: 生效时间
- `validUntil`: 到期时间，可为空
- `status`: `active | pending | expired | revoked`
- `grantedBy`: 授权人
- `createdAt`: 创建时间
- `updatedAt`: 更新时间

### Relationships

- 多对一关联 `RoleDefinition`
- 可关联 `DelegationGrant`

### Validation Rules

- 授权范围不得超出授权人管理边界
- `sourceType=temporary` 时必须存在 `validUntil`
- `status=expired` 或 `status=revoked` 时不得继续生效

## 8. DelegationGrant

### Purpose

表示一条可审计的委派管理关系。

### Fields

- `id`: 主键标识
- `grantorRef`: 委派人
- `delegateRef`: 被委派人
- `allowedRoleLevels`: 可委派层级范围
- `allowedScopeSnapshot`: 可操作范围快照
- `status`: `active | suspended | expired | revoked`
- `validFrom`: 生效时间
- `validUntil`: 到期时间
- `reason`: 委派原因
- `createdAt`: 创建时间
- `updatedAt`: 更新时间

### Relationships

- 一对多关联 `RoleAssignment`

### Validation Rules

- 被委派范围不得超出委派人当前有效权限边界
- 到期后必须自动停止新的授权分发

## 9. SessionRecord

### Purpose

表示一次登录会话及其治理状态。

### Fields

- `id`: 主键标识
- `userId`: 用户
- `identitySourceId`: 身份来源
- `loginMethod`: 登录方式
- `status`: `active | idle | revoked | expired | risk-blocked`
- `riskLevel`: `low | medium | high | critical`
- `permissionVersion`: 登录时权限版本
- `lastSeenAt`: 最近活动时间
- `revokedAt`: 回收时间
- `createdAt`: 创建时间

### Relationships

- 多对一关联 `IdentitySource`
- 一对多关联 `IdentityAuditEvent`

### Validation Rules

- `status=revoked` 时必须阻止后续受保护访问
- `riskLevel=critical` 时必须在统一视图中显著提示

## 10. AccessRiskSnapshot

### Purpose

表示针对用户或租户范围聚合后的访问风险摘要。

### Fields

- `id`: 主键标识
- `subjectType`: `user | organization | workspace | project`
- `subjectRef`: 目标引用
- `riskType`: `over-privileged | residual-access | inheritance-sprawl | stale-session | source-conflict`
- `severity`: `low | medium | high | critical`
- `summary`: 风险摘要
- `recommendedAction`: 建议动作
- `status`: `open | acknowledged | mitigated | accepted`
- `generatedAt`: 生成时间
- `updatedAt`: 更新时间

### Validation Rules

- 同一对象的同类风险快照应可被更新而不是无限重复堆积
- `status=mitigated` 时必须保留处理说明

## State Transitions

- `IdentitySource`: `draft -> active -> disabled -> active`，不可直接从 `disabled` 删除
- `OrganizationUnit`: `active -> disabled -> archived`
- `RoleAssignment`: `pending -> active -> expired/revoked`
- `DelegationGrant`: `active -> suspended/expired/revoked`
- `SessionRecord`: `active -> idle -> expired`，或 `active -> revoked/risk-blocked`
- `AccessRiskSnapshot`: `open -> acknowledged -> mitigated/accepted`
