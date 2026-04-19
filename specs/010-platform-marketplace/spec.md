# Feature Specification: 平台应用目录与扩展市场

**Feature Branch**: `010-platform-marketplace`  
**Created**: 2026-04-19  
**Status**: Review Ready  
**Input**: User description: "我要新增 010-platform-marketplace，严格对标 Rancher 的应用目录和扩展生态能力，面向平台管理员、平台工程团队和应用交付团队，在多集群 Kubernetes 平台中提供统一的应用目录、模板中心和扩展机制。用户需要能够管理应用目录来源、模板分类、版本、依赖关系、参数表单、部署约束、适用范围和发布说明，能够将标准化应用模板发布到指定工作空间、项目或集群范围，并查看安装记录、升级入口、版本变更和下线状态。平台还需要支持扩展包、插件或集成模块的注册、启停、版本兼容性、权限声明和可见范围，让平台能够持续扩展新的产品能力而不破坏核心治理边界。该特性必须与现有权限模型、审计模型和发布模型衔接，保证模板发布和扩展安装都在受控范围内进行。首期范围聚焦应用目录、模板分发和扩展机制，不包含完整 GitOps 编排、集群生命周期、统一身份源和合规扫描。"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - 管理应用目录与模板中心 (Priority: P1)

作为平台管理员或平台工程团队成员，我希望集中管理应用目录来源、模板分类、版本、依赖关系、参数表单和发布说明，这样平台可以持续维护一套标准化、可复用、可审计的应用模板资产。

**Why this priority**: 如果没有统一的应用目录和模板中心，后续模板分发、安装记录和扩展能力都缺少可信来源与治理基础。

**Independent Test**: 新增一个目录来源并发布多个模板版本后，管理员可以在模板中心查看分类、版本、依赖关系、参数表单、部署约束和适用范围。

**Acceptance Scenarios**:

1. **Given** 平台尚未配置应用目录来源，**When** 管理员新增一个目录来源并完成启用，**Then** 平台应记录来源状态、同步结果和可见模板集合。
2. **Given** 某个模板存在多个版本和依赖关系，**When** 平台工程团队查看模板详情，**Then** 平台应展示版本列表、依赖说明、参数表单、部署约束和发布说明。
3. **Given** 某个模板已被标记下线，**When** 用户查看应用目录，**Then** 平台应明确显示下线状态并阻止新的安装入口继续使用该模板版本。

---

### User Story 2 - 按范围分发模板并跟踪安装升级 (Priority: P1)

作为应用交付团队成员，我希望把标准化应用模板发布到指定工作空间、项目或集群范围，并持续查看安装记录、升级入口和版本变化，这样交付动作可以在受控范围内复用平台标准模板。

**Why this priority**: 目录资产只有能够被安全分发和复用，才真正形成平台交付价值。

**Independent Test**: 将一个模板发布到目标工作空间或项目后，交付团队可以查看可用范围、安装记录、升级入口、版本变更和下线提示。

**Acceptance Scenarios**:

1. **Given** 平台已存在可用模板，**When** 管理员将模板发布到指定工作空间、项目或集群范围，**Then** 目标范围内的用户应只看到自己被授权使用的模板。
2. **Given** 某个模板已经在目标范围内被安装，**When** 模板发布了新版本，**Then** 平台应展示版本差异、升级入口和适用约束。
3. **Given** 某个模板已停止分发或下线，**When** 用户查看已安装记录，**Then** 平台应保留历史记录、标记当前状态并提示是否允许继续升级或仅保留运行态。

---

### User Story 3 - 注册扩展并控制平台扩展边界 (Priority: P2)

作为平台管理员，我希望注册扩展包、插件或集成模块，并管理其启停、版本兼容性、权限声明和可见范围，这样平台可以扩展新能力而不会破坏现有治理边界。

**Why this priority**: 扩展机制是平台持续演进的关键，但必须建立在权限、兼容性和范围控制都可治理的前提上。

**Independent Test**: 注册一个扩展模块并声明其兼容版本、权限范围和可见范围后，管理员可以启停该扩展，并查看安装记录、兼容状态和审计记录。

**Acceptance Scenarios**:

1. **Given** 平台支持扩展注册，**When** 管理员导入一个扩展包或插件元数据，**Then** 平台应记录其版本、兼容性、权限声明和可见范围。
2. **Given** 某个扩展与当前平台版本或目标范围不兼容，**When** 管理员尝试启用该扩展，**Then** 平台应阻止生效并明确给出兼容性原因。
3. **Given** 某个扩展已经安装并对部分范围可见，**When** 管理员停用或下线该扩展，**Then** 平台应保留安装记录、变化说明和审计轨迹。

### Edge Cases

- 当模板依赖的上游模板版本不可用、被下线或与目标范围不兼容时，平台必须阻止发布并明确提示依赖原因。
- 当同一模板在不同范围内存在不同可见版本时，平台必须避免用户看到超出自身范围的版本入口。
- 当扩展声明的权限超出当前授权模型可承载的边界时，平台必须阻止注册或启用。
- 当模板已下线但仍有运行中的安装实例时，平台必须保留历史和状态说明，而不是直接删除记录。
- 当扩展停用会影响已启用的平台能力时，平台必须在执行前提示影响范围和受影响对象。
- 当目录来源同步失败、版本元数据不完整或发布说明缺失时，平台必须阻止将该版本作为可安装资产分发。

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: 系统 MUST 提供统一的平台应用目录与扩展市场中心，用于管理目录来源、模板中心、模板分发和扩展注册。
- **FR-002**: 系统 MUST 支持管理应用目录来源，包括新增、启停、同步状态查看和来源健康状态查看。
- **FR-003**: 系统 MUST 支持按分类管理应用模板，并记录模板名称、版本、依赖关系、参数表单、部署约束、适用范围和发布说明。
- **FR-004**: 系统 MUST 支持为同一模板维护多个版本，并明确展示版本之间的可见状态、升级关系和下线状态。
- **FR-005**: 系统 MUST 支持标记模板版本为可用、禁用、下线或仅保留历史状态。
- **FR-006**: 系统 MUST 支持将标准化模板发布到指定工作空间、项目或集群范围。
- **FR-007**: 系统 MUST 确保用户只能看到自己被授权范围内可使用的模板和版本。
- **FR-008**: 系统 MUST 记录每次模板发布、撤回、安装、升级和下线相关的历史记录。
- **FR-009**: 系统 MUST 为已安装模板展示安装记录、当前版本、可升级版本、版本变更说明和当前分发生命周期状态。
- **FR-010**: 系统 MUST 在模板存在依赖关系时校验依赖是否满足，并在不满足时阻止继续分发或安装。
- **FR-011**: 系统 MUST 支持为模板维护参数表单说明和部署约束，并在使用前展示这些约束。
- **FR-012**: 系统 MUST 支持注册扩展包、插件或集成模块，并维护其版本、兼容性、权限声明和可见范围。
- **FR-013**: 系统 MUST 支持启用、停用和下线扩展，并保留其生命周期记录。
- **FR-014**: 系统 MUST 在扩展启用前校验版本兼容性、权限声明和目标可见范围。
- **FR-015**: 系统 MUST 阻止不兼容扩展、超范围扩展或超权限声明扩展进入可用状态。
- **FR-016**: 系统 MUST 将模板发布、扩展安装和扩展启停动作纳入现有权限模型控制。
- **FR-017**: 系统 MUST 将目录来源管理、模板分发、扩展注册、扩展启停和版本变化动作纳入现有审计模型。
- **FR-018**: 系统 MUST 与现有发布模型衔接，使模板分发、安装记录和升级入口能够复用已有交付治理语义。
- **FR-019**: 系统 MUST 支持按目录来源、模板分类、模板状态、目标范围、扩展状态和兼容性状态筛选对象。
- **FR-020**: 系统 MUST 为模板和扩展分别展示适用范围、可见范围和受影响范围，防止跨租户或跨边界误用。
- **FR-021**: 系统 MUST 在目录同步失败、依赖缺失、版本冲突或权限不满足时返回明确原因。
- **FR-022**: 系统 MUST 支持查看模板和扩展的版本变更说明，帮助使用者理解升级影响和下线影响。
- **FR-023**: 首期范围 MUST 聚焦应用目录、模板分发和扩展机制。
- **FR-024**: 首期范围 MUST 不包含完整 GitOps 编排、集群生命周期、统一身份源和合规扫描。
- **FR-025**: 首期范围 MUST 为后续更多目录来源、更多扩展类型和更复杂分发生命周期预留扩展空间。

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

- **Catalog Source**: 表示一个应用目录来源，包含来源身份、同步状态、可用性和目录归属。
- **Application Template**: 表示一个可分发的标准化应用模板，包含分类、版本、依赖关系、参数表单、部署约束、适用范围和发布说明。
- **Template Release Scope**: 表示模板被发布到的工作空间、项目或集群范围，以及其可见性和生效状态。
- **Installation Record**: 表示模板在某个范围内的安装、升级、保留或下线历史。
- **Extension Package**: 表示一个扩展包、插件或集成模块，包含版本、兼容性、权限声明、启停状态和可见范围。
- **Compatibility Statement**: 表示模板或扩展与平台版本、目标范围或依赖对象之间的兼容结论。

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 平台管理员能够在 30 分钟内完成至少一个目录来源接入，并让目录中的模板在中心视图中可检索。
- **SC-002**: 平台工程团队能够在 20 分钟内完成一个标准模板的分类、版本整理和范围发布配置。
- **SC-003**: 应用交付团队能够在 10 分钟内定位目标模板、查看发布说明，并判断其是否可在自身范围内安装或升级。
- **SC-004**: 100% 的模板发布、撤回、升级入口变更、扩展注册、扩展启停和扩展下线动作都能被检索到审计记录。
- **SC-005**: 100% 的超范围模板分发、依赖不满足模板使用和不兼容扩展启用请求都会被阻止并返回明确原因。
- **SC-006**: 试点范围内至少 90% 的模板查询、安装记录查询和扩展兼容性查询能够在 30 秒内得到结果。
- **SC-007**: 平台团队能够在统一视图中识别模板状态、版本变化、目标范围、扩展兼容性和下线影响，而无需额外查阅外部台账。

## Assumptions

- `004-gitops-and-release` 已经提供基础发布模型和交付历史语义，010 在此之上扩展目录中心、模板分发与扩展机制。
- 首期目标用户为平台管理员、平台工程团队和应用交付团队，不覆盖普通业务用户自行上传扩展或自行维护目录来源。
- 首期以平台内标准化模板和平台扩展治理为主，不要求覆盖完整第三方市场结算或商业分发流程。
- 模板分发与扩展注册均需要延续既有权限边界、工作空间/项目边界和审计边界。
- 首期不包含完整 GitOps 编排、集群生命周期、统一身份源和合规扫描，这些能力如需联动将在后续特性中扩展。
