# Data Model: 多集群 Kubernetes 可观测中心

> 说明：002 复用 001 中已有的 `Cluster`、`Workspace`、`Project`、`PlatformUser` 和审计模型；本文件只描述 002 新增或被显著扩展的观测治理实体。原始日志全文、原始时序指标点和原始 Event 长期归档不属于平台持久化模型。

## 1. ObservabilityDataSource
- Purpose: 描述某个集群或全局范围内可用的观测后端连接信息和健康状态。
- Key Fields:
  - `id`
  - `cluster_id`: 可为空，表示全局共享数据源
  - `datasource_type`: `metrics | logs | alerts | events`
  - `provider_kind`: `prometheus-compatible | loki-compatible | alertmanager-compatible | k8s-events | custom-compatible`
  - `name`
  - `base_url`
  - `auth_secret_ref`: 凭据引用，不保存明文
  - `status`: `pending | healthy | degraded | unreachable | disabled`
  - `last_verified_at`
  - `last_error`
  - `created_by`
  - `created_at` / `updated_at`
- Relationships:
  - 可关联到一个 `Cluster`
- Validation Rules:
  - 同一 `cluster_id + datasource_type + name` 组合必须唯一
  - 明文密钥、令牌和密码不得直接落库
- State Transitions:
  - `pending -> healthy | degraded | unreachable | disabled`
  - `healthy -> degraded | unreachable | disabled`
  - `degraded -> healthy | unreachable | disabled`
  - `unreachable -> healthy | disabled`

## 2. ObservabilityScope
- Purpose: 统一描述一次观测操作或治理动作适用的业务范围。
- Key Fields:
  - `cluster_ids`
  - `workspace_ids`
  - `project_ids`
  - `namespaces`
  - `resource_kinds`
  - `resource_names`
  - `label_selector`
  - `time_range`
- Relationships:
  - 被 `AlertRule`、`SilenceWindow`、`LogQuerySession` 和 `AlertIncidentSnapshot` 复用
- Validation Rules:
  - 任何作用域都必须能映射到现有工作空间/项目授权边界
  - 时间范围不得为负值或空范围

## 3. LogQuerySession
- Purpose: 表示一次围绕资源上下文执行的日志筛选与检索请求。
- Key Fields:
  - `id`
  - `requested_by`
  - `scope_snapshot_json`
  - `keyword`
  - `start_at`
  - `end_at`
  - `direction`: `backward | forward`
  - `limit`
  - `result_status`: `ready | partial | unavailable | denied`
  - `executed_at`
  - `expires_at`
- Relationships:
  - 由 `PlatformUser` 发起
  - 依赖 `ObservabilityDataSource`
- Validation Rules:
  - 单次查询时间跨度必须受平台配置限制
  - 未授权范围不得生成成功结果
- Notes:
  - 该实体可落在 Redis 或短期缓存层，不要求长期保留

## 4. AlertRule
- Purpose: 平台定义的告警规则治理实体。
- Key Fields:
  - `id`
  - `name`
  - `description`
  - `severity`: `info | warning | critical`
  - `scope_snapshot_json`
  - `condition_expression`
  - `evaluation_window`
  - `notification_strategy_json`
  - `status`: `enabled | disabled`
  - `sync_status`: `pending | synced | failed`
  - `last_synced_at`
  - `last_sync_error`
  - `created_by`
  - `created_at` / `updated_at`
- Relationships:
  - 关联多个 `NotificationTarget`
  - 可生成多个 `AlertIncidentSnapshot`
- Validation Rules:
  - `name` 在同一业务范围内必须唯一
  - 禁用规则不得继续向外部后端投影为可执行状态

## 5. NotificationTarget
- Purpose: 平台统一管理的告警通知对象。
- Key Fields:
  - `id`
  - `name`
  - `target_type`: `webhook | email | chat | custom`
  - `config_secret_ref`
  - `status`: `active | disabled`
  - `scope_snapshot_json`
  - `created_by`
  - `created_at` / `updated_at`
- Relationships:
  - 可被多个 `AlertRule` 复用
- Validation Rules:
  - 同一作用域下通知目标名称必须唯一
  - 明文通知凭据不得直接落库

## 6. SilenceWindow
- Purpose: 对一组告警范围进行临时抑制或静默处理。
- Key Fields:
  - `id`
  - `name`
  - `scope_snapshot_json`
  - `reason`
  - `starts_at`
  - `ends_at`
  - `status`: `scheduled | active | expired | canceled`
  - `created_by`
  - `canceled_by`
  - `created_at` / `updated_at`
- Relationships:
  - 影响多个 `AlertIncidentSnapshot`
- Validation Rules:
  - `ends_at` 必须大于 `starts_at`
  - 取消后的静默窗口不得重新激活
- State Transitions:
  - `scheduled -> active | canceled`
  - `active -> expired | canceled`

## 7. AlertIncidentSnapshot
- Purpose: 表示一次告警触发、处理与恢复过程的快照记录。
- Key Fields:
  - `id`
  - `source_incident_key`
  - `rule_id`
  - `cluster_id`
  - `workspace_id`
  - `project_id`
  - `resource_kind`
  - `resource_name`
  - `namespace`
  - `severity`
  - `status`: `firing | acknowledged | silenced | resolved`
  - `summary`
  - `starts_at`
  - `acknowledged_at`
  - `resolved_at`
  - `last_synced_at`
  - `timeline_json`
- Relationships:
  - 可关联一个 `AlertRule`
  - 可拥有多个 `AlertHandlingRecord`
- Validation Rules:
  - 同一数据源中的 `source_incident_key` 必须唯一映射到单条快照
  - 已 `resolved` 的告警不得再写入新的确认动作
- State Transitions:
  - `firing -> acknowledged | silenced | resolved`
  - `acknowledged -> silenced | resolved`
  - `silenced -> firing | resolved`

## 8. AlertHandlingRecord
- Purpose: 记录告警确认、说明补充、交接和恢复复盘等处理动作。
- Key Fields:
  - `id`
  - `incident_id`
  - `action_type`: `acknowledge | note | handoff | recover-note`
  - `content`
  - `acted_by`
  - `acted_at`
- Relationships:
  - 属于一个 `AlertIncidentSnapshot`
  - 由一个 `PlatformUser` 产生
- Validation Rules:
  - 未授权用户不得写入处理记录
  - 已关闭告警仅允许补充复盘说明，不允许重新确认

## 9. ResourceContextView
- Purpose: 表示面向单个资源的统一观测联动视图，是 002 的核心读取对象。
- Key Fields:
  - `cluster_id`
  - `workspace_id`
  - `project_id`
  - `namespace`
  - `resource_kind`
  - `resource_name`
  - `log_summary`
  - `event_summary`
  - `metric_summary`
  - `active_alerts`
  - `data_freshness`
- Relationships:
  - 聚合 `LogQuerySession`、`EventTimelineItem`、`MetricInsight` 和 `AlertIncidentSnapshot`
- Validation Rules:
  - 任何聚合结果都必须通过统一授权过滤
- Notes:
  - 该实体属于查询结果模型，不要求长期持久化

## 10. MetricInsight
- Purpose: 面向集群、节点、命名空间或工作负载的指标摘要和趋势结果。
- Key Fields:
  - `subject_type`: `cluster | node | namespace | workload | pod`
  - `subject_ref`
  - `window`
  - `cpu_usage`
  - `memory_usage`
  - `restarts`
  - `availability`
  - `trend_direction`
  - `anomaly_flags`
  - `observed_at`
- Relationships:
  - 属于某个 `ResourceContextView` 或概览面板
- Validation Rules:
  - 指标窗口必须与用户请求的时间范围一致
  - 数据缺失时必须显式返回状态而不是填充伪造零值

## 11. EventTimelineItem
- Purpose: 资源或范围内的标准化事件时间线项。
- Key Fields:
  - `cluster_id`
  - `namespace`
  - `involved_kind`
  - `involved_name`
  - `event_type`: `normal | warning`
  - `reason`
  - `message`
  - `first_seen_at`
  - `last_seen_at`
  - `count`
- Relationships:
  - 可被 `ResourceContextView` 和范围查询复用
- Validation Rules:
  - 时间线项必须能映射回真实资源上下文
  - 警告事件不得在 UI 中与正常事件混淆
