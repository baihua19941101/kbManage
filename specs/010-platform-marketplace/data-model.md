# Data Model: 平台应用目录与扩展市场

## 1. CatalogSource

- **Purpose**: 表示一个目录来源，负责提供模板集合和版本清单。
- **Fields**:
  - `id`
  - `name`
  - `sourceType`：如 Git 风格、Helm 风格、OCI 风格或平台内建目录
  - `endpointRef`
  - `status`：draft、active、degraded、disabled
  - `syncState`：idle、syncing、succeeded、failed
  - `lastSyncedAt`
  - `lastError`
  - `ownerUserId`
  - `visibilityScope`
  - `createdAt`
  - `updatedAt`
- **Validation Rules**:
  - 名称在平台内唯一。
  - 禁用或异常来源不得对外提供新的可分发模板版本。
- **Relationships**:
  - 一个 `CatalogSource` 包含多个 `ApplicationTemplate`。

## 2. ApplicationTemplate

- **Purpose**: 表示一个标准化应用模板。
- **Fields**:
  - `id`
  - `catalogSourceId`
  - `name`
  - `slug`
  - `category`
  - `summary`
  - `publishStatus`：active、disabled、retired、history-only
  - `defaultVersionId`
  - `supportedScopes`
  - `releaseNotesSummary`
  - `createdAt`
  - `updatedAt`
- **Validation Rules**:
  - 同一目录来源下 `slug` 唯一。
  - 模板必须至少有一个可解析版本后才能进入可用状态。
- **Relationships**:
  - 一个 `ApplicationTemplate` 包含多个 `TemplateVersion`。
  - 一个 `ApplicationTemplate` 可被多个 `TemplateReleaseScope` 分发。

## 3. TemplateVersion

- **Purpose**: 表示模板的具体版本。
- **Fields**:
  - `id`
  - `templateId`
  - `version`
  - `status`：draft、active、deprecated、retired
  - `dependencySnapshot`
  - `parameterSchemaSummary`
  - `deploymentConstraintSummary`
  - `releaseNotes`
  - `isUpgradeable`
  - `supersedesVersionId`
  - `createdAt`
  - `updatedAt`
- **Validation Rules**:
  - 同一模板下版本号唯一。
  - 依赖缺失或约束不满足的版本不得进入 active。
- **Relationships**:
  - 一个 `TemplateVersion` 对应多个 `CompatibilityStatement`。
  - 一个 `TemplateVersion` 可关联多个 `InstallationRecord`。

## 4. TemplateReleaseScope

- **Purpose**: 表示模板被发布到的工作空间、项目或集群范围。
- **Fields**:
  - `id`
  - `templateId`
  - `versionId`
  - `scopeType`：workspace、project、cluster
  - `scopeRef`
  - `status`：published、withdrawn、retired
  - `visibilityMode`
  - `publishedBy`
  - `publishedAt`
  - `withdrawnAt`
- **Validation Rules**:
  - 同一模板版本在同一范围内只能存在一个当前有效发布记录。
  - 超出授权范围的发布请求必须被阻止。
- **Relationships**:
  - 一个 `TemplateReleaseScope` 可对应多个 `InstallationRecord`。

## 5. InstallationRecord

- **Purpose**: 表示模板在某个范围内的安装、升级或下线历史。
- **Fields**:
  - `id`
  - `templateId`
  - `versionId`
  - `scopeType`
  - `scopeRef`
  - `releaseScopeId`
  - `lifecycleStatus`：installed、upgrade-available、upgraded、retired、orphaned
  - `currentInstalledVersion`
  - `upgradeTargetVersion`
  - `changeSummary`
  - `installedAt`
  - `lastChangedAt`
- **Validation Rules**:
  - 下线模板仍需保留历史记录。
  - 记录必须能够追溯到模板、版本和目标范围。
- **Relationships**:
  - 多条安装记录可关联同一模板版本。

## 6. ExtensionPackage

- **Purpose**: 表示一个扩展包、插件或集成模块。
- **Fields**:
  - `id`
  - `name`
  - `extensionType`
  - `version`
  - `status`：draft、registered、enabled、disabled、retired
  - `compatibilityPolicy`
  - `permissionDeclaration`
  - `visibilityScope`
  - `entrySummary`
  - `ownerUserId`
  - `createdAt`
  - `updatedAt`
- **Validation Rules**:
  - 扩展启用前必须具备兼容性结论和权限声明。
  - 权限声明超出平台治理边界时不得注册为可启用状态。
- **Relationships**:
  - 一个 `ExtensionPackage` 可对应多个 `CompatibilityStatement`。
  - 一个 `ExtensionPackage` 可对应多个 `ExtensionLifecycleRecord`。

## 7. CompatibilityStatement

- **Purpose**: 表示模板依赖或扩展兼容性的治理结论。
- **Fields**:
  - `id`
  - `ownerType`：template-version、extension-package
  - `ownerRef`
  - `targetType`：platform-version、dependency-template、scope
  - `targetRef`
  - `result`：compatible、warning、blocked
  - `summary`
  - `evaluatedAt`
  - `evaluator`
- **Validation Rules**:
  - blocked 结果必须提供明确原因。
  - 同一 owner 与 target 的最新结论需可检索。

## 8. ExtensionLifecycleRecord

- **Purpose**: 表示扩展注册、启用、停用和下线历史。
- **Fields**:
  - `id`
  - `extensionPackageId`
  - `action`：register、enable、disable、retire
  - `scopeType`
  - `scopeRef`
  - `outcome`
  - `reason`
  - `executedBy`
  - `executedAt`
- **Validation Rules**:
  - 所有扩展生命周期动作都必须留存。
  - 停用或下线时必须说明影响范围。

## State Transitions

- `CatalogSource.status`: draft -> active/degraded/disabled
- `CatalogSource.syncState`: idle -> syncing -> succeeded/failed
- `TemplateVersion.status`: draft -> active -> deprecated/retired
- `TemplateReleaseScope.status`: published -> withdrawn/retired
- `InstallationRecord.lifecycleStatus`: installed -> upgrade-available -> upgraded/retired/orphaned
- `ExtensionPackage.status`: draft -> registered -> enabled/disabled -> retired
