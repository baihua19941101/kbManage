# Feature Specification: 多集群 Kubernetes 安全与策略治理中心

**Feature Branch**: `005-security-and-policy`  
**Created**: 2026-04-14  
**Status**: Draft  
**Input**: User description: "我要新增 005-security-and-policy，严格对标 Rancher 在安全与策略治理方面的能力，覆盖策略中心、准入控制和工作负载安全基线。面向平台管理员、安全管理员、合规负责人和项目管理员，在多集群 Kubernetes 环境中提供统一的策略定义、分发、执行和违规治理能力。用户需要能够管理平台级、工作空间级和项目级的安全策略与准入规则，按集群、命名空间、项目或资源类型分配策略，配置审计、告警、仅提示和强制执行等模式，查看策略命中结果、违规对象、风险级别、例外申请、例外时效和整改状态，并管理 Pod安全等级、镜像来源限制、资源约束、标签规范、网络和准入相关控制策略。平台必须支持分阶段启用策略、灰度验证和例外管理，避免一次性强制执行引发大面积影响，并且所有策略变更和违规处置都必须可审计、可追踪。首期范围聚焦策略治理与准入控制，不包含CIS 或 STIG 合规扫描、平台身份源整合、应用发布和灾备恢复。"

## 当前状态/执行说明（2026-04-14）

- 已完成 `/speckit.specify`，005 规格和质量清单已生成。
- 已完成 `/speckit.plan`，生成 `plan.md`、`research.md`、`data-model.md`、`contracts/openapi.yaml` 和 `quickstart.md`。
- 已完成 `/speckit.tasks`，`tasks.md` 已生成并按 Governance、Setup、Foundational、US1、US2、US3、Polish 分阶段组织。
- 已进入 `/speckit.implement` 阶段，并完成首批治理证据文件：`artifacts/005-security-and-policy/branch-check.txt`、`artifacts/005-security-and-policy/backup-manifest.txt`、`artifacts/005-security-and-policy/mirror-and-remote-check.txt`。
- 当前执行分支为 `005-security-and-policy`。
- 005 实现已启动；后续若发生新的高风险数据库变更，仍需再次执行备份并留存证据。
- implement 进展：已完成 Phase 0（T001-T004）、Phase 1（T005-T009）、Phase 2（T010-T019）、US1（T020-T032）、US2（T033-T044）、US3（T045-T055）和 Final Phase（T056-T059）全量任务。
- 已完成低并发验证：`cd backend && go test -p 1 ./...` 通过；`frontend` 安全策略新增页面测试与 ESLint 通过（测试使用 `--maxWorkers=1`）。
- 当前状态：005 implement 已完成，待你确认后进入 PR 提交流程（仍需“用户明确同意后再合并”）。
- 后续交付仍需遵守“中文 PR + 用户明确同意后再合并”的治理门槛。

## User Scenarios & Testing *(mandatory)*

### User Story 1 - 统一策略中心与分层策略管理 (Priority: P1)

作为平台管理员或安全管理员，我希望在统一入口中定义和维护平台级、工作空间级、项目级的安全策略，并按集群、命名空间、项目和资源类型进行分发，以便在多集群环境下形成一致且可控的策略治理基线。

**Why this priority**: 如果缺少统一策略建模与分发能力，后续准入控制和违规治理无法形成统一标准，平台会继续依赖分散、不可追踪的手工策略配置。

**Independent Test**: 创建一组平台级策略和两组范围化策略后，管理员可独立验证策略是否在目标范围内生效、继承关系是否清晰、以及不同范围的策略是否可被区分展示。

**Acceptance Scenarios**:

1. **Given** 平台管理员已进入策略中心，**When** 创建平台级策略并分配到多个目标集群，**Then** 平台应保存策略定义、作用范围、执行模式和当前状态，并显示生效范围摘要。
2. **Given** 安全管理员需要按组织层级管理策略，**When** 为工作空间或项目创建策略并绑定资源类型，**Then** 平台应明确显示策略层级、适用范围和与上层策略的关系。
3. **Given** 同一资源可能命中多条策略，**When** 用户查看目标对象的策略视图，**Then** 平台应展示最终适用策略集合和每条策略来源层级，避免范围混淆。

---

### User Story 2 - 准入控制模式与分阶段启用 (Priority: P1)

作为安全管理员或项目管理员，我希望在策略执行中配置审计、告警、仅提示和强制执行等模式，并支持灰度验证和分阶段启用，这样可以先评估影响再逐步收紧策略，避免一次性强制导致大面积拦截。

**Why this priority**: 准入控制直接影响业务部署链路，没有分阶段启用和灰度验证将显著放大变更风险，难以在生产环境安全落地策略。

**Independent Test**: 选择一条新策略先以仅提示模式在部分命名空间灰度验证，再切换到强制执行后，团队能够看到模式切换前后的命中变化与拦截结果。

**Acceptance Scenarios**:

1. **Given** 某条策略处于初次启用阶段，**When** 安全管理员将其配置为审计或仅提示模式，**Then** 平台应记录命中结果且不阻断目标对象创建或变更。
2. **Given** 团队完成灰度验证并准备收紧策略，**When** 将策略切换为强制执行并扩大到更多目标范围，**Then** 平台应按新模式阻断违规请求并标识执行模式变更时间点。
3. **Given** 某个项目短期内无法满足新策略，**When** 项目管理员申请并获得临时例外，**Then** 平台应在例外有效期内按例外规则处理命中结果，并在到期后自动恢复原策略约束。

---

### User Story 3 - 违规治理闭环与审计追踪 (Priority: P2)

作为合规负责人或审计人员，我希望持续查看策略命中结果、违规对象、风险级别、例外申请与整改状态，并追踪每次策略变更和违规处置记录，以便形成可审计、可追责、可复盘的治理闭环。

**Why this priority**: 安全治理不仅是“有策略”，更要“可执行、可跟踪、可复盘”；缺少违规闭环会导致策略效果无法评估，合规责任无法落实。

**Independent Test**: 在平台产生多条违规事件、例外申请和整改更新后，审计人员可按时间、策略、范围、风险等级和处置状态检索完整记录并导出复盘材料。

**Acceptance Scenarios**:

1. **Given** 平台已检测到多条策略违规，**When** 合规负责人按风险级别和处置状态筛选，**Then** 平台应返回对应违规对象、命中策略和当前整改进度。
2. **Given** 某条违规需要临时豁免，**When** 安全管理员审批或拒绝例外申请，**Then** 平台应记录申请人、审批人、原因、有效期和最终状态。
3. **Given** 审计人员需要复盘策略治理效果，**When** 按策略变更与违规处置链路查询，**Then** 平台应提供可追踪的时间线记录，覆盖策略变更、命中结果、例外和整改结论。

### Edge Cases

- 当同一资源同时命中平台级和项目级策略时，平台必须明确展示最终生效规则和冲突处理结果，避免执行歧义。
- 当策略执行模式从仅提示切换到强制执行时，平台必须提前标识影响范围和预计风险，防止误触发大面积拦截。
- 当例外申请处于审批中或已过期时，平台必须区分“待审批”“生效中”“已过期”“已撤销”状态，避免错误放行。
- 当目标集群或命名空间暂时不可达时，平台必须区分“策略未下发成功”和“策略无命中结果”两种状态。
- 当用户在治理过程中权限被回收时，平台必须立即阻止其继续执行策略变更、例外审批或整改状态更新。
- 当策略作用范围被修改后，平台必须保留变更前后范围差异，支持审计回溯。

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: 系统 MUST 提供统一的安全与策略治理中心，支持在多集群范围内集中管理策略定义、分发和执行状态。
- **FR-002**: 系统 MUST 支持平台级、工作空间级和项目级策略建模，并清晰展示各层级策略的适用范围。
- **FR-003**: 系统 MUST 支持按集群、命名空间、项目和资源类型分配策略。
- **FR-004**: 系统 MUST 支持在策略中配置执行模式，至少覆盖审计、告警、仅提示和强制执行。
- **FR-005**: 系统 MUST 支持策略分阶段启用，允许先在部分范围灰度验证再逐步扩大。
- **FR-006**: 系统 MUST 在策略模式或作用范围变更时记录变更前后差异，并可被追踪查询。
- **FR-007**: 系统 MUST 支持维护工作负载安全基线策略，至少覆盖 Pod 安全等级、镜像来源限制、资源约束、标签规范、网络和准入相关控制。
- **FR-008**: 系统 MUST 为每条策略提供启用、停用、更新和撤销能力，并保留每次状态变更记录。
- **FR-009**: 系统 MUST 在策略命中时记录目标对象、命中规则、风险级别、发生时间和当前处置状态。
- **FR-010**: 系统 MUST 支持按策略、集群、命名空间、项目、资源类型、风险级别、执行模式和时间范围筛选命中结果。
- **FR-011**: 系统 MUST 支持对违规对象记录整改状态，并跟踪从发现到关闭的治理流程。
- **FR-012**: 系统 MUST 支持例外申请流程，记录申请原因、申请范围、审批结果和有效期限。
- **FR-013**: 系统 MUST 在例外生效期间按例外规则处理命中结果，并在例外到期或撤销后自动恢复原策略执行。
- **FR-014**: 系统 MUST 支持查看例外申请全生命周期状态，至少包括待审批、生效中、已拒绝、已过期和已撤销。
- **FR-015**: 系统 MUST 支持在策略视图中展示策略命中趋势和主要违规分布，用于评估治理效果。
- **FR-016**: 系统 MUST 对策略管理动作执行权限隔离，确保用户仅能操作其被授权范围内的策略与对象。
- **FR-017**: 系统 MUST 在用户权限变化后立即限制其策略查看、策略变更、例外审批和整改更新能力。
- **FR-018**: 系统 MUST 对策略创建、更新、模式切换、范围调整、启停、例外审批和整改处置生成可检索审计记录。
- **FR-019**: 系统 MUST 支持审计人员按时间、操作者、策略层级、策略对象、动作类型和结果检索治理记录。
- **FR-020**: 系统 MUST 在策略未成功分发、范围无效、权限不足、对象已变化和执行冲突等场景下返回可理解的错误说明。
- **FR-021**: 首期范围 MUST 聚焦策略治理与准入控制，不包含 CIS 或 STIG 合规扫描。
- **FR-022**: 首期范围 MUST 不包含平台身份源整合、应用发布和灾备恢复能力。

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
- **GC-007**: If subagents are used for implementation, they MUST use `gpt-5.3-codex`.

### Key Entities *(include if feature involves data)*

- **Security Policy**: 一条可被分发和执行的安全规则，包含策略层级、作用范围、执行模式、风险定义和当前状态。
- **Policy Assignment**: 策略与目标范围之间的绑定关系，描述策略作用到哪些集群、命名空间、项目或资源类型。
- **Admission Rule Profile**: 准入控制策略集合，定义审计、告警、仅提示和强制执行等模式的启用与切换规则。
- **Violation Record**: 一次策略命中或违规结果，包含目标对象、命中策略、风险级别、发生时间、处置状态和整改进度。
- **Exception Request**: 针对特定策略命中提出的例外申请，包含申请原因、适用范围、审批结论、生效时段和失效时间。
- **Remediation Action**: 针对违规对象的整改动作记录，包含责任人、处理步骤、当前状态和关闭结论。
- **Policy Audit Event**: 策略治理过程中产生的审计事件，覆盖策略变更、模式调整、范围调整、例外审批和违规处置。

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 平台管理员能够在 20 分钟内完成一条策略从定义、范围分配、灰度启用到状态确认的完整流程。
- **SC-002**: 在至少 20 个已接入集群的环境中，90% 的策略命中查询可在 30 秒内返回筛选结果。
- **SC-003**: 在试点期内，90% 的新策略能够先以非强制模式完成灰度验证，再按计划切换到目标执行模式。
- **SC-004**: 对于已发现的高风险违规项，80% 能在 24 小时内完成例外审批或整改动作登记。
- **SC-005**: 权限验收中，100% 的跨工作空间、跨项目未授权策略变更和例外审批请求都被拦截，且不暴露敏感治理细节。
- **SC-006**: 审计人员针对最近 90 天的策略治理记录检索时，90% 的查询可在 30 秒内返回可复盘结果集。
- **SC-007**: 试点团队中至少 80% 的准入策略治理与违规处置流程可在平台内闭环完成，无需依赖分散工具。

## Assumptions

- `001-k8s-ops-platform` 已提供多集群接入、资源范围模型、权限隔离与审计基础能力，本特性在其之上扩展策略治理能力。
- `002-observability-center` 与 `003-workload-operations-control-plane` 提供的资源上下文与处置入口可作为违规定位和整改跟踪参考，但不改变本特性的首期范围。
- 首期目标用户为平台管理员、安全管理员、合规负责人和项目管理员，不单独覆盖移动端专属体验。
- 首期优先满足策略中心、准入控制、灰度启用和例外管理闭环，不引入独立的外部合规评测体系。
- 组织已具备基础的角色与授权体系，本特性重点解决策略治理流程和违规处置可追踪性。
- 首期不包含 CIS 或 STIG 合规扫描、平台身份源整合、应用发布和灾备恢复能力。
