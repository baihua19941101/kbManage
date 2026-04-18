# Data Model: 多集群 Kubernetes 集群生命周期中心

## 1. ClusterLifecycleRecord

### Purpose

统一表示一个受平台管理的集群主记录，覆盖导入、注册、创建、升级、停用和退役全过程中的核心状态。

### Fields

- `id`: 主键标识
- `name`: 集群名称
- `displayName`: 展示名称
- `lifecycleMode`: `imported | registered | provisioned`
- `infrastructureType`: 基础设施类型，例如裸机、虚拟化、公有云、托管 Kubernetes
- `driverRef`: 当前驱动标识
- `driverVersion`: 当前驱动版本
- `workspaceId`: 所属工作空间
- `projectId`: 关联项目，可为空
- `status`: `pending | active | degraded | upgrading | disabled | retiring | retired | failed`
- `registrationStatus`: `not_required | pending | issued | connected | failed`
- `healthStatus`: `healthy | warning | critical | unknown`
- `kubernetesVersion`: 当前控制面版本
- `targetVersion`: 目标升级版本，可为空
- `nodePoolSummary`: 节点池摘要
- `lastValidationStatus`: 最近校验状态
- `lastValidationAt`: 最近校验时间
- `lastOperationId`: 最近关键动作引用
- `retirementReason`: 退役原因，可为空
- `createdBy`: 创建人
- `createdAt`: 创建时间
- `updatedAt`: 更新时间

### Relationships

- 一对多关联 `NodePoolProfile`
- 一对多关联 `UpgradePlan`
- 一对多关联 `LifecycleOperation`
- 多对一关联 `ClusterDriverVersion`
- 可选多对一关联 `ClusterTemplate`

### Validation Rules

- `name` 在同一授权范围内必须唯一
- `status=retired` 时不可再发起升级、扩缩或停用动作
- `targetVersion` 仅在升级计划存在时可设置

## 2. ClusterDriverVersion

### Purpose

表示一种集群驱动及其版本化定义，用于支撑导入、创建、升级和能力矩阵展示。

### Fields

- `id`: 主键标识
- `driverKey`: 驱动唯一键
- `version`: 驱动版本号
- `displayName`: 展示名称
- `providerType`: 适配的基础设施类型
- `status`: `draft | active | deprecated | disabled`
- `capabilityProfileVersion`: 关联能力定义版本
- `schemaVersion`: 参数模式版本
- `releaseNotes`: 版本说明
- `createdAt`: 创建时间
- `updatedAt`: 更新时间

### Relationships

- 一对多关联 `CapabilityMatrixEntry`
- 一对多关联 `ClusterTemplate`
- 一对多关联 `ClusterLifecycleRecord`

### Validation Rules

- 同一 `driverKey` 下 `version` 必须唯一
- `status=disabled` 时不可用于新建集群

## 3. CapabilityMatrixEntry

### Purpose

描述某个驱动版本或集群类型在某个能力域上的支持状态与兼容结论。

### Fields

- `id`: 主键标识
- `ownerType`: `driver | cluster-type`
- `ownerRef`: 驱动版本或集群类型引用
- `capabilityDomain`: `network | storage | identity | observability | security | backup | release`
- `supportLevel`: `native | extended | partial | unsupported`
- `compatibilityStatus`: `compatible | conditional | incompatible`
- `constraintsSummary`: 约束说明
- `recommendedFor`: 推荐使用场景
- `updatedAt`: 更新时间

### Validation Rules

- 同一 `ownerRef + capabilityDomain` 只能存在一条生效条目
- `supportLevel=unsupported` 时 `compatibilityStatus` 不能为 `compatible`

## 4. ClusterTemplate

### Purpose

表示模板化创建集群所使用的参数集合、能力要求和驱动兼容范围。

### Fields

- `id`: 主键标识
- `name`: 模板名称
- `description`: 模板说明
- `infrastructureType`: 目标基础设施类型
- `driverKey`: 依赖驱动键
- `driverVersionRange`: 适用驱动版本范围
- `requiredCapabilities`: 必需能力列表
- `parameterSchemaRef`: 参数模式引用
- `defaultValuesRef`: 默认参数引用
- `status`: `draft | active | deprecated | disabled`
- `createdBy`: 创建人
- `createdAt`: 创建时间
- `updatedAt`: 更新时间

### Relationships

- 多对一关联 `ClusterDriverVersion`
- 一对多关联 `LifecycleOperation` 中的创建动作

### Validation Rules

- `status=disabled` 时不可用于新建集群
- 模板必须声明至少一个驱动依赖或基础设施类型约束

## 5. LifecycleOperation

### Purpose

统一表示导入、注册、创建、校验、升级、节点池调整、停用和退役等生命周期动作执行记录。

### Fields

- `id`: 主键标识
- `clusterId`: 目标集群，可为空（创建前）
- `operationType`: `import | register | create | validate | upgrade | scale-node-pool | disable | retire`
- `triggerSource`: `manual | scheduled | follow-up`
- `status`: `pending | running | partially_succeeded | succeeded | failed | canceled | blocked`
- `riskLevel`: `low | medium | high | critical`
- `requestedBy`: 发起人
- `requestSnapshot`: 请求快照
- `resultSummary`: 结果摘要
- `failureReason`: 失败原因
- `startedAt`: 开始时间
- `completedAt`: 完成时间

### Validation Rules

- 同一集群存在 `running` 的关键动作时，新关键动作必须被阻断
- `operationType=retire` 时必须记录退役原因或结论

## 6. UpgradePlan

### Purpose

表示一个集群升级计划及其执行状态。

### Fields

- `id`: 主键标识
- `clusterId`: 目标集群
- `fromVersion`: 当前版本
- `toVersion`: 目标版本
- `windowStart`: 升级窗口开始时间
- `windowEnd`: 升级窗口结束时间
- `precheckStatus`: `pending | passed | warning | failed`
- `impactSummary`: 影响摘要
- `status`: `draft | approved | running | succeeded | failed | canceled`
- `lastOperationId`: 关联执行动作
- `createdBy`: 创建人
- `createdAt`: 创建时间
- `updatedAt`: 更新时间

### Validation Rules

- `toVersion` 必须高于或区别于 `fromVersion`
- `precheckStatus=failed` 时不可进入执行状态

## 7. NodePoolProfile

### Purpose

表示集群节点池的期望配置和当前运行状态。

### Fields

- `id`: 主键标识
- `clusterId`: 所属集群
- `name`: 节点池名称
- `role`: `control-plane | worker | mixed`
- `desiredCount`: 期望节点数
- `currentCount`: 当前节点数
- `minCount`: 最小容量
- `maxCount`: 最大容量
- `version`: 节点池版本
- `zoneRefs`: 可用区或资源池列表
- `status`: `pending | active | scaling | upgrading | degraded | failed`
- `lastOperationId`: 最近节点池动作引用
- `updatedAt`: 更新时间

### Validation Rules

- `desiredCount` 必须在 `minCount` 与 `maxCount` 范围内
- 正在升级或退役中的集群，节点池不能进入新的扩缩动作

## 8. LifecycleAuditEvent

### Purpose

表示生命周期域的标准化审计记录。

### Fields

- `id`: 主键标识
- `action`: 动作类型
- `actorUserId`: 操作者
- `clusterId`: 目标集群，可为空
- `targetType`: `cluster | driver | template | node-pool | upgrade-plan`
- `targetRef`: 目标引用
- `outcome`: `succeeded | failed | blocked | canceled`
- `detailSnapshot`: 详情快照
- `occurredAt`: 发生时间

### Validation Rules

- 所有关键动作都必须生成对应审计事件
- `targetType` 与 `targetRef` 必须成对出现

## State Transitions

### ClusterLifecycleRecord.status

- `pending -> active`
- `pending -> failed`
- `active -> upgrading`
- `active -> disabled`
- `active -> retiring`
- `upgrading -> active`
- `upgrading -> failed`
- `disabled -> retiring`
- `retiring -> retired`
- `retiring -> failed`

### UpgradePlan.status

- `draft -> approved`
- `approved -> running`
- `running -> succeeded`
- `running -> failed`
- `draft -> canceled`
- `approved -> canceled`

### NodePoolProfile.status

- `pending -> active`
- `active -> scaling`
- `active -> upgrading`
- `scaling -> active`
- `upgrading -> active`
- `scaling -> failed`
- `upgrading -> failed`
