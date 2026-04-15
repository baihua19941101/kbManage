# Data Model: 多集群 Kubernetes 安全与策略治理中心

> 说明：005 复用 001 的 `Cluster`、`Workspace`、`Project`、权限与审计底座。本文件仅描述 005 新增或显著扩展的策略治理域实体。

## 1. SecurityPolicy
- Purpose: 表示一条可复用的安全策略定义。
- Key Fields:
  - `id`
  - `name`
  - `scope_level`: `platform | workspace | project`
  - `category`: `pod-security | image | resource | label | network | admission`
  - `rule_template`
  - `risk_level`: `low | medium | high | critical`
  - `default_enforcement_mode`: `audit | alert | warn | enforce`
  - `status`: `draft | active | disabled | archived`
  - `created_by`
  - `created_at` / `updated_at`
- Relationships:
  - 可关联多个 `PolicyAssignment`
  - 可关联多个 `PolicyVersion`
- Validation Rules:
  - 同一作用域下策略名称唯一
  - `archived` 状态策略不可再被新分配引用
- State Transitions:
  - `draft -> active | disabled`
  - `active -> disabled | archived`

## 2. PolicyVersion
- Purpose: 记录策略规则变更历史，支持追溯和回放。
- Key Fields:
  - `id`
  - `policy_id`
  - `version`
  - `rule_snapshot`
  - `change_summary`
  - `changed_by`
  - `changed_at`
- Relationships:
  - 属于一个 `SecurityPolicy`
- Validation Rules:
  - 同一策略下版本号递增且唯一

## 3. PolicyAssignment
- Purpose: 定义策略在目标范围内的实际分配与执行模式。
- Key Fields:
  - `id`
  - `policy_id`
  - `workspace_id`
  - `project_id`
  - `cluster_refs`
  - `namespace_refs`
  - `resource_kinds`
  - `enforcement_mode`: `audit | alert | warn | enforce`
  - `rollout_stage`: `pilot | canary | broad | full`
  - `status`: `pending | active | failed | paused`
  - `effective_from`
  - `effective_to`
  - `created_at` / `updated_at`
- Relationships:
  - 属于一个 `SecurityPolicy`
  - 可关联多个 `PolicyHitRecord`
- Validation Rules:
  - 目标范围必须在操作者授权边界内
  - `enforce` 模式必须有可回滚切换记录
- State Transitions:
  - `pending -> active | failed | paused`
  - `active -> paused | failed`

## 4. PolicyHitRecord
- Purpose: 表示一次策略命中或违规事实。
- Key Fields:
  - `id`
  - `policy_id`
  - `assignment_id`
  - `cluster_id`
  - `namespace`
  - `resource_kind`
  - `resource_name`
  - `resource_uid`
  - `hit_result`: `pass | warn | block`
  - `risk_level`
  - `message`
  - `detected_at`
  - `remediation_status`: `open | in_progress | mitigated | closed`
- Relationships:
  - 可关联零或一个 `ExceptionRequest`
  - 可关联多个 `RemediationAction`
- Validation Rules:
  - 被阻断事件必须保留命中规则快照
  - 状态变更必须可追踪到操作者或系统动作

## 5. ExceptionRequest
- Purpose: 表示针对策略命中的临时例外申请。
- Key Fields:
  - `id`
  - `policy_id`
  - `hit_record_id`
  - `scope_snapshot`
  - `reason`
  - `requested_by`
  - `reviewed_by`
  - `status`: `pending | approved | rejected | active | expired | revoked`
  - `starts_at`
  - `expires_at`
  - `review_comment`
  - `created_at` / `updated_at`
- Relationships:
  - 关联一个 `PolicyHitRecord`
- Validation Rules:
  - `expires_at` 必须晚于 `starts_at`
  - `active` 状态例外必须在有效期内
- State Transitions:
  - `pending -> approved | rejected`
  - `approved -> active`
  - `active -> expired | revoked`

## 6. RemediationAction
- Purpose: 表示针对违规对象的整改动作。
- Key Fields:
  - `id`
  - `hit_record_id`
  - `action_type`: `config-update | image-replace | resource-adjust | label-fix | network-fix | exception-linked | manual-note`
  - `owner`
  - `status`: `todo | in_progress | done | canceled`
  - `summary`
  - `evidence_ref`
  - `created_at`
  - `completed_at`
- Relationships:
  - 属于一个 `PolicyHitRecord`
- Validation Rules:
  - `done` 状态必须记录完成时间和结论摘要

## 7. PolicyDistributionTask
- Purpose: 表示策略分发或模式切换的一次异步任务。
- Key Fields:
  - `id`
  - `policy_id`
  - `assignment_batch_id`
  - `operator_id`
  - `operation`: `assign | mode-switch | pause | resume | revoke`
  - `status`: `pending | running | partially_succeeded | succeeded | failed`
  - `target_count`
  - `succeeded_count`
  - `failed_count`
  - `result_summary`
  - `failure_reason`
  - `started_at`
  - `completed_at`
- Relationships:
  - 可关联多个 `PolicyAssignment`
- Validation Rules:
  - 部分成功必须输出失败目标明细

## 8. PolicyAuditEvent
- Purpose: 描述策略治理关键动作的标准化审计对象。
- Key Fields:
  - `action`: `policy.create | policy.update | policy.activate | policy.disable | policy.assign | policy.mode.switch | policy.exception.request | policy.exception.approve | policy.exception.reject | policy.exception.revoke | policy.remediation.update`
  - `operator_id`
  - `workspace_id`
  - `project_id`
  - `policy_id`
  - `assignment_id`
  - `hit_record_id`
  - `exception_request_id`
  - `outcome`
  - `details_json`
  - `occurred_at`
- Relationships:
  - 复用 001 审计查询框架
- Validation Rules:
  - 关键变更必须具备变更前后摘要和影响范围
