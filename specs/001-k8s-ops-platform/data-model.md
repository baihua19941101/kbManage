# Data Model: 多集群 Kubernetes 可视化管理平台

## 1. PlatformUser
- Purpose: 平台登录主体，支持用户名密码登录和平台级身份治理。
- Key Fields:
  - `id`: 主键，雪花 ID 或自增主键
  - `username`: 登录名，全局唯一，3-64 个字符
  - `display_name`: 展示名称，1-128 个字符
  - `email`: 邮箱，可选但唯一
  - `phone`: 手机号，可选
  - `status`: `active | locked | disabled`
  - `password_hash`: 密码哈希，不保存明文
  - `password_changed_at`: 最近一次修改密码时间
  - `last_login_at`: 最近成功登录时间
  - `created_at` / `updated_at`
- Validation Rules:
  - `username` 必须唯一且不可重复使用大小写变体
  - 被 `disabled` 的用户不可登录
  - 密码更新后必须使历史刷新令牌失效

## 2. UserGroup
- Purpose: 用户分组，用于批量授予平台角色或工作空间/项目角色。
- Key Fields:
  - `id`
  - `name`: 全局唯一
  - `description`
  - `created_at` / `updated_at`
- Relationships:
  - 与 `PlatformUser` 多对多

## 3. PlatformRole
- Purpose: 平台级 RBAC 角色，描述全局治理能力。
- Key Fields:
  - `id`
  - `role_key`: 唯一键，首批固定为 `platform-admin`、`ops-operator`、`audit-reader`、`readonly`
  - `name`
  - `description`
  - `is_system`: 是否系统内置
- Relationships:
  - 与 `PlatformPermission` 多对多
  - 通过 `PlatformRoleBinding` 绑定到用户或用户组
- Validation Rules:
  - 首期仅允许上述 4 个系统内置角色键，新增角色需后续版本扩展

## 4. PlatformPermission
- Purpose: 平台级权限点，例如集群接入、平台用户管理、全局审计查看。
- Key Fields:
  - `id`
  - `permission_key`: 唯一键，例如 `cluster.manage`、`audit.read.global`
  - `resource_scope`: `platform | cluster | audit | identity`
  - `description`

## 5. PlatformRoleBinding
- Purpose: 将平台角色授予用户或用户组。
- Key Fields:
  - `id`
  - `subject_type`: `user | group`
  - `subject_id`
  - `platform_role_id`
  - `granted_by`
  - `created_at`
- Validation Rules:
  - 同一主体不可重复绑定同一平台角色

## 6. Workspace
- Purpose: 顶层业务隔离边界，对应团队、产品线或组织域。
- Key Fields:
  - `id`
  - `name`: 唯一键
  - `code`: 唯一短编码
  - `description`
  - `status`: `active | archived`
  - `owner_user_id`: 负责人
  - `created_at` / `updated_at`
- Relationships:
  - 一个工作空间下有多个 `Project`
  - 可通过 `WorkspaceClusterBinding` 关联多个集群

## 7. Project
- Purpose: 工作空间下的二级边界，承载环境、应用或业务模块。
- Key Fields:
  - `id`
  - `workspace_id`
  - `name`
  - `code`
  - `description`
  - `status`: `active | archived`
  - `owner_user_id`
  - `created_at` / `updated_at`
- Validation Rules:
  - `workspace_id + code` 唯一
  - 被归档项目默认只读

## 8. ScopeRole
- Purpose: 工作空间/项目层的角色定义，例如 `workspace-owner`、`project-operator`。
- Key Fields:
  - `id`
  - `scope_type`: `workspace | project`
  - `role_key`
  - `name`
  - `description`
  - `is_system`

## 9. ScopePermission
- Purpose: 工作空间/项目层的权限点，覆盖资源可见性和资源操作。
- Key Fields:
  - `id`
  - `scope_type`: `workspace | project`
  - `permission_key`: 例如 `resource.read`、`operation.execute`
  - `description`

## 10. ScopeRoleBinding
- Purpose: 将作用域角色授予用户或用户组，并限定到工作空间或项目。
- Key Fields:
  - `id`
  - `subject_type`: `user | group`
  - `subject_id`
  - `scope_type`: `workspace | project`
  - `scope_id`
  - `scope_role_id`
  - `granted_by`
  - `created_at`
- Validation Rules:
  - 同一主体在同一作用域内同一角色不能重复授予

## 11. Cluster
- Purpose: 被平台接入的 Kubernetes 集群。
- Key Fields:
  - `id`
  - `name`: 平台内唯一
  - `display_name`
  - `description`
  - `environment`: `prod | staging | test | other`
  - `api_server`: 集群 API 地址
  - `status`: `pending | healthy | degraded | unreachable | disabled`
  - `k8s_version`
  - `last_sync_at`
  - `created_by`
  - `created_at` / `updated_at`
- Relationships:
  - 一个 `Cluster` 可绑定多个 `Workspace`
  - 一个 `Cluster` 有一份 `ClusterCredential`
- State Transitions:
  - `pending -> healthy`
  - `healthy -> degraded | unreachable | disabled`
  - `degraded -> healthy | unreachable | disabled`
  - `unreachable -> healthy | disabled`

## 12. ClusterCredential
- Purpose: 保存访问集群所需的凭据引用。
- Key Fields:
  - `id`
  - `cluster_id`
  - `credential_type`: `kubeconfig | token | service-account`
  - `secret_ciphertext`: 加密后的凭据内容
  - `secret_version`
  - `last_verified_at`
  - `created_at` / `updated_at`
- Validation Rules:
  - 明文凭据禁止落库
  - 更新凭据后必须触发连通性校验

## 13. WorkspaceClusterBinding
- Purpose: 定义哪些工作空间可以管理哪些集群。
- Key Fields:
  - `id`
  - `workspace_id`
  - `cluster_id`
  - `default_namespaces`: 默认可见命名空间规则
  - `created_at`
- Validation Rules:
  - `workspace_id + cluster_id` 唯一

## 14. ResourceInventory
- Purpose: 平台用于跨集群浏览和筛选的资源索引快照。
- Key Fields:
  - `id`
  - `cluster_id`
  - `workspace_id`: 可为空，表示未映射
  - `project_id`: 可为空
  - `resource_uid`: 集群资源唯一标识
  - `api_group`
  - `api_version`
  - `resource_kind`
  - `namespace`
  - `name`
  - `labels_json`
  - `status_summary`
  - `raw_summary_json`: 摘要信息，不保存完整大对象
  - `last_observed_at`
- Validation Rules:
  - `cluster_id + resource_uid` 唯一
  - `resource_kind` 首期仅允许 `Deployment`、`StatefulSet`、`DaemonSet`、`Pod`、`Service`、`Ingress`、`Node`、`Namespace`
  - 仅存索引用于列表和筛选，详情页实时拉取最新资源

## 15. OperationRequest
- Purpose: 表示一次由平台用户发起的资源管理或运维操作。
- Key Fields:
  - `id`
  - `cluster_id`
  - `workspace_id`
  - `project_id`
  - `resource_uid`
  - `resource_kind`
  - `operation_type`: `create | update | delete | scale | restart | cordon | drain | custom`
  - `risk_level`: `low | medium | high | critical`
  - `requested_by`
  - `approved_by`: 预留字段，首期固定为空（不启用他人审批流）
  - `request_payload_json`
  - `status`: `pending_confirm | queued | running | succeeded | failed | canceled`
  - `failure_reason`
  - `created_at` / `updated_at` / `completed_at`
- State Transitions:
  - `pending_confirm -> queued | canceled`
  - `queued -> running | canceled`
  - `running -> succeeded | failed`
- Validation Rules:
  - 高风险动作必须经过二次确认后才能从 `pending_confirm` 进入 `queued`
  - 首期流程为“二次确认即执行”，不允许引入审批人或审批节点

## 16. AuditEvent
- Purpose: 审计记录，覆盖访问、授权变更、资源操作和结果。
- Key Fields:
  - `id`
  - `event_type`: `login | logout | permission_change | cluster_change | resource_read | operation_execute`
  - `actor_user_id`
  - `actor_username`
  - `actor_roles_json`
  - `scope_type`: `platform | workspace | project | cluster | resource`
  - `scope_id`
  - `cluster_id`
  - `resource_kind`
  - `resource_namespace`
  - `resource_name`
  - `result`: `success | failure | denied`
  - `request_id`
  - `detail_json`
  - `occurred_at`
- Validation Rules:
  - `occurred_at` 必须可排序并支持按时间范围检索
  - `detail_json` 仅保存审计必要信息，不保存敏感明文
  - 首期导出格式仅支持 CSV，且导出前必须对敏感字段脱敏（至少包含访问凭据、令牌、密钥、密码、手机号、邮箱）

## 17. RefreshSession
- Purpose: 管理刷新令牌与会话失效。
- Key Fields:
  - `id`
  - `user_id`
  - `session_key`
  - `client_ip`
  - `user_agent`
  - `status`: `active | revoked | expired`
  - `expires_at`
  - `created_at` / `updated_at`
- Validation Rules:
  - 密码重置、用户禁用或主动退出时，相关 `active` 会话必须被撤销

## 核心关系概览
- `PlatformUser` <-M:N-> `UserGroup`
- `PlatformRole` <-M:N-> `PlatformPermission`
- `PlatformRoleBinding` -> `PlatformUser | UserGroup`
- `Workspace` 1:N `Project`
- `Workspace` M:N `Cluster` via `WorkspaceClusterBinding`
- `ScopeRole` <-M:N-> `ScopePermission`
- `ScopeRoleBinding` -> `Workspace | Project`
- `Cluster` 1:1 `ClusterCredential`
- `Cluster` 1:N `ResourceInventory`
- `PlatformUser` 1:N `OperationRequest`
- `PlatformUser` 1:N `AuditEvent`
- `OperationRequest` 1:N `AuditEvent`

## 推荐索引
- `platform_users(username)` 唯一索引
- `clusters(name)` 唯一索引
- `projects(workspace_id, code)` 唯一索引
- `resource_inventory(cluster_id, resource_kind, namespace, name)` 组合索引
- `resource_inventory(cluster_id, workspace_id, project_id, last_observed_at)` 组合索引
- `audit_events(occurred_at, event_type, result)` 组合索引
- `audit_events(actor_user_id, occurred_at)` 组合索引
- `operation_requests(cluster_id, status, created_at)` 组合索引
