# Data Model: 多集群 GitOps 与应用发布中心

> 说明：004 复用 001 中已有的 `Cluster`、`Workspace`、`Project`、基础授权与审计底座，也复用 003 的异步动作执行经验。本文件只描述 004 新增或被显著扩展的交付域实体。

## 1. DeliverySource
- Purpose: 表示一个可被平台接入和验证的交付来源，是 004 的来源入口对象。
- Key Fields:
  - `id`
  - `name`
  - `source_type`: `git | package`
  - `endpoint`
  - `default_ref`
  - `credential_ref`
  - `workspace_scope`
  - `status`: `pending | ready | failed | disabled`
  - `last_verified_at`
  - `last_error_message`
  - `created_at` / `updated_at`
- Relationships:
  - 可被多个 `ApplicationDeliveryUnit` 引用
- Validation Rules:
  - 同一可见范围内来源名称必须唯一
  - 禁用来源不得继续作为新交付单元的有效来源
- State Transitions:
  - `pending -> ready | failed | disabled`
  - `ready -> failed | disabled`

## 2. ClusterTargetGroup
- Purpose: 表示一组可复用的交付目标集合，用于复用多集群发布范围。
- Key Fields:
  - `id`
  - `name`
  - `workspace_id`
  - `project_id`
  - `cluster_refs`
  - `cluster_selector_snapshot`
  - `description`
  - `status`: `active | stale | disabled`
  - `created_at` / `updated_at`
- Relationships:
  - 可被多个 `ApplicationDeliveryUnit` 与 `EnvironmentStage` 引用
- Validation Rules:
  - 目标集合中每个集群都必须属于当前用户可授权范围
  - 禁用目标组不得再被新环境阶段引用

## 3. EnvironmentStage
- Purpose: 表示一个具备顺序语义的环境阶段，用于建模测试、预发、生产等推进顺序。
- Key Fields:
  - `id`
  - `delivery_unit_id`
  - `name`
  - `order_index`
  - `target_group_id`
  - `promotion_mode`: `manual | automatic`
  - `paused`
  - `status`: `idle | waiting | progressing | succeeded | failed | paused`
  - `last_entered_at`
  - `last_completed_at`
- Relationships:
  - 属于一个 `ApplicationDeliveryUnit`
  - 引用一个 `ClusterTargetGroup`
  - 可拥有多个 `ConfigurationOverlay`
- Validation Rules:
  - 同一交付单元内 `order_index` 必须唯一且连续
  - 被暂停阶段不得自动进入后续推进
- State Transitions:
  - `idle -> waiting | paused`
  - `waiting -> progressing | paused`
  - `progressing -> succeeded | failed | paused`

## 4. ConfigurationOverlay
- Purpose: 表示针对环境或目标范围追加的配置覆盖。
- Key Fields:
  - `id`
  - `delivery_unit_id`
  - `environment_stage_id`
  - `overlay_type`: `values | patch | manifest-snippet`
  - `overlay_ref`
  - `precedence`
  - `effective_scope`
  - `created_at` / `updated_at`
- Relationships:
  - 属于一个 `ApplicationDeliveryUnit`
  - 可挂载到一个 `EnvironmentStage`
- Validation Rules:
  - 同一作用域内覆盖优先级不得冲突
  - 平台必须能生成最终生效配置预览，否则该覆盖不得进入可发布状态

## 5. ApplicationDeliveryUnit
- Purpose: 表示一个被 GitOps 持续交付管理的应用单元，是 004 的核心对象。
- Key Fields:
  - `id`
  - `name`
  - `workspace_id`
  - `project_id`
  - `source_id`
  - `source_path`
  - `default_namespace`
  - `sync_mode`: `manual | auto`
  - `release_policy_json`
  - `desired_revision`
  - `desired_app_version`
  - `desired_config_version`
  - `paused`
  - `delivery_status`: `ready | progressing | degraded | out_of_sync | paused | unknown`
  - `last_synced_at`
  - `last_release_id`
  - `created_at` / `updated_at`
- Relationships:
  - 引用一个 `DeliverySource`
  - 拥有多个 `EnvironmentStage`、`ConfigurationOverlay`、`ReleaseRevision`、`DeliveryOperation`
- Validation Rules:
  - 交付单元必须至少绑定一个环境阶段
  - 未通过来源验证的交付单元不得进入自动同步

## 6. DeliveryStatusSnapshot
- Purpose: 表示交付单元在某个环境和目标上的状态聚合视图。
- Key Fields:
  - `delivery_unit_id`
  - `environment_stage_id`
  - `target_group_id`
  - `desired_state_summary`
  - `live_state_summary`
  - `sync_status`: `synced | pending | running | failed | paused | unknown`
  - `drift_status`: `in_sync | drifted | modified | unavailable | unknown`
  - `last_sync_result`
  - `last_error_message`
  - `last_observed_at`
- Relationships:
  - 由 `ApplicationDeliveryUnit` 与 `EnvironmentStage` 聚合生成
- Validation Rules:
  - 不可把来源不可达误表示为“无差异”
  - 部分成功必须保留目标级细分结果
- Notes:
  - 该实体属于查询结果模型，不要求长期完整持久化

## 7. DeliveryOperation
- Purpose: 表示一次同步、发布、推进、回滚、暂停、恢复或卸载动作请求。
- Key Fields:
  - `id`
  - `request_id`
  - `operator_id`
  - `delivery_unit_id`
  - `environment_stage_id`
  - `action_type`: `install | sync | resync | upgrade | promote | rollback | pause | resume | uninstall`
  - `target_release_id`
  - `payload_json`
  - `status`: `pending | running | partially_succeeded | succeeded | failed | canceled`
  - `progress_percent`
  - `result_summary`
  - `failure_reason`
  - `started_at`
  - `completed_at`
  - `created_at` / `updated_at`
- Relationships:
  - 可引用一个 `ReleaseRevision`
  - 可拥有多个 `DeliveryStageExecution`
- Validation Rules:
  - `rollback` 必须指向明确可恢复的 `ReleaseRevision`
  - 终态动作不得再次进入运行态
- State Transitions:
  - `pending -> running | canceled`
  - `running -> partially_succeeded | succeeded | failed | canceled`

## 8. DeliveryStageExecution
- Purpose: 表示一次交付动作在某个环境阶段上的执行结果。
- Key Fields:
  - `id`
  - `operation_id`
  - `environment_stage_id`
  - `target_group_id`
  - `status`: `pending | running | succeeded | failed | skipped | blocked`
  - `target_count`
  - `succeeded_count`
  - `failed_count`
  - `result_summary`
  - `failure_reason`
  - `started_at`
  - `completed_at`
- Relationships:
  - 属于一个 `DeliveryOperation`
  - 引用一个 `EnvironmentStage`
- Validation Rules:
  - 同一推进动作必须按环境顺序写入阶段执行记录
  - 被阻断阶段不得伪装成成功

## 9. ReleaseRevision
- Purpose: 表示平台内一次可识别、可回滚的交付修订。
- Key Fields:
  - `id`
  - `delivery_unit_id`
  - `source_revision`
  - `app_version`
  - `config_version`
  - `effective_scope_snapshot`
  - `release_notes_summary`
  - `created_by`
  - `created_at`
  - `rollback_available`
  - `status`: `active | historical | failed | rolled_back`
- Relationships:
  - 属于一个 `ApplicationDeliveryUnit`
  - 可被 `DeliveryOperation` 作为升级目标或回滚目标引用
- Validation Rules:
  - 同一交付单元内的修订身份必须由来源修订、应用版本、配置版本和目标快照组合唯一确定
  - 已失效来源或缺失目标快照的历史修订不得展示为可回滚
- State Transitions:
  - `active -> historical | failed | rolled_back`

## 10. GitOpsAuditEvent
- Purpose: 描述 004 中来源变更、同步、发布、推进、回滚和配置变更的标准化审计对象。
- Key Fields:
  - `action`: `gitops.source.create | gitops.source.update | gitops.source.verify | gitops.delivery.sync | gitops.delivery.promote | gitops.delivery.rollback | gitops.delivery.pause | gitops.delivery.resume | gitops.delivery.uninstall | gitops.config.update`
  - `operator_id`
  - `workspace_id`
  - `project_id`
  - `environment_stage_id`
  - `delivery_unit_id`
  - `target_group_id`
  - `release_revision_id`
  - `outcome`
  - `details_json`
  - `occurred_at`
- Relationships:
  - 复用 001 的审计查询与持久化能力
- Validation Rules:
  - 审计记录必须能映射回真实交付单元、环境或来源对象
  - 关键发布动作必须同时记录范围、目标版本和执行结果
