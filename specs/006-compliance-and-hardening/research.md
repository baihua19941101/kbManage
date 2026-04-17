# Research: 多集群 Kubernetes 合规与加固中心

> 研究输入基于 006 规格、001-005 既有能力边界以及企业级 Kubernetes 合规治理通用语义，目标是在不扩展到准入阻断、自动原位修复和统一身份治理的前提下，形成可实施的首期方案。

## Decision 1: 006 采用独立 `compliance` 业务域，而不是复用 005 `securitypolicy`

- Decision: 后端新增 `service/compliance`、`integration/compliance` 与对应 handler/router/repository/worker，前端新增 `features/compliance-hardening`；复用 001 的授权与审计底座，但扫描与整改复核模型独立。
- Rationale: 006 的核心对象是“扫描执行、失败项证据、整改/例外/复检与趋势复盘”，与 005 的“准入策略定义、执行模式和命中治理”不同。独立业务域可避免权限语义、页面入口和异步任务模型相互污染。
- Alternatives considered:
  - 并入 `securitypolicy`: 会把“评估报告”与“实时策略”混为一体，导致对象生命周期过长且职责不清。
  - 并入 `observability`: 会把合规治理误降为只读报表，缺少整改与审批闭环。

## Decision 2: 扫描执行采用“平台统一编排 + 外部扫描器/基线包适配层”模式

- Decision: 平台负责定义基线、扫描配置、调度执行、结果归集、证据快照与治理流转；具体检查动作通过 `integration/compliance` 适配层对接外部扫描器或基线执行包，不在首期将某一个扫描器实现写死到业务模型中。
- Rationale: 这样既能保持控制面一致性，又能避免首期被单一工具绑定；后续扩展不同基线执行来源时，只需新增适配实现而不改业务语义。
- Alternatives considered:
  - 平台内嵌固定扫描器逻辑：工具耦合过高，不利于后续基线扩展。
  - 完全依赖外部系统回传报告：平台将失去调度、权限、审计和闭环控制能力。

## Decision 3: 基线标准与扫描配置分离建模，且基线版本必须显式可见

- Decision: `ComplianceBaseline` 负责描述 `CIS`、`STIG` 或平台基线模板的标准与版本口径；`ScanProfile` 负责目标范围、执行频率和过滤配置；每次 `ScanExecution` 必须固化所使用的基线版本快照。
- Rationale: 分离模型便于同一基线被多个范围复用，也能避免历史结果在基线版本变化后失去可比性。
- Alternatives considered:
  - 将范围和基线定义写成一个对象：复用性差，变更成本高。
  - 不记录基线版本快照：历史对比会失真，难以支撑审计复盘。

## Decision 4: 首期关键资源评估对象固定为“工作负载 + 身份/网络控制 + 基础范围”

- Decision: 除 `Cluster`、`Node`、`Namespace` 外，首期关键资源评估对象固定聚焦 `Deployment`、`StatefulSet`、`DaemonSet`、`Pod`、`ServiceAccount`、`Role`、`RoleBinding`、`ClusterRole`、`ClusterRoleBinding`、`NetworkPolicy`。
- Rationale: 这些对象覆盖工作负载安全、权限暴露、网络边界和运行范围四类常见加固场景，能够在首期提供足够治理价值，同时避免对象面过宽导致实现失控。
- Alternatives considered:
  - 全量覆盖所有 Kubernetes 资源：首期范围不可控，证据模型和检查口径会显著膨胀。
  - 只做节点与命名空间：无法体现工作负载与授权边界相关的加固价值。

## Decision 5: 失败项与证据采用不可变快照模型，而不是引用实时对象状态

- Decision: 每次扫描产生的 `ComplianceFinding` 与 `EvidenceRecord` 以快照形式保存，记录采集时间、证据摘要、可信度和对象上下文；当目标对象后续删除、迁移或变更时，历史结果仍保持可追溯。
- Rationale: 合规复盘和审计依赖“当时看到什么”，不能因对象后续变化而丢失证据语义。
- Alternatives considered:
  - 仅保存对象引用：对象删除或变更后，历史记录无法解释。
  - 保存完整原始报告正文：首期存储和敏感信息管理成本过高。

## Decision 6: 整改、例外和复检保持独立一等对象，并以失败项为中心串联

- Decision: `RemediationTask`、`ComplianceExceptionRequest` 和 `RecheckTask` 作为独立对象建模，并都显式关联 `ComplianceFinding`；例外语义参考 005，但不直接复用 005 的实体表。
- Rationale: 006 的例外是“针对扫描发现的治理例外”，而 005 的例外偏向“针对策略命中的运行态例外”。独立建模有利于后续统计、复检和归档查询。
- Alternatives considered:
  - 直接复用 005 例外表：会混淆来源域、审批语义和生命周期。
  - 只用任务评论串联处置：无法支持结构化查询和到期恢复。

## Decision 7: 趋势与汇总视图采用定期快照聚合，而不是每次查询实时全量重算

- Decision: 通过 `ComplianceTrendSnapshot` 定期汇总覆盖率、风险分布、整改进度和遗留风险，趋势页面优先读取快照；细粒度追溯再回到原始扫描和失败项。
- Rationale: 面向 20+ 集群和 90 天历史范围时，实时全量重算会显著增加查询成本，不利于管理汇报场景的稳定响应。
- Alternatives considered:
  - 每次查询实时扫描全量执行记录：实现简单但查询代价高、结果不稳定。
  - 只展示单次扫描快照：无法满足长期趋势和汇报需求。

## Decision 8: 首期明确不做自动修复与实时阻断，闭环止于“受控处置 + 复检”

- Decision: 006 首期的治理闭环止于整改任务、例外审批、复检和归档；不自动执行原位修复、不批量强制下发加固配置，也不把扫描结果直接转为实时阻断策略。
- Rationale: 这样可以在不放大生产变更风险的前提下，先建立评估和治理闭环，后续若需要自动修复，可作为独立 feature 讨论。
- Alternatives considered:
  - 扫描后直接自动修复：风险过高，且需要更多审批与回滚机制。
  - 扫描后直接联动准入阻断：会与 005 准入策略域重叠，并改变现有业务边界。

## Decision 9: 归档采用“结构化导出任务 + 原始视图复用”模式

- Decision: 006 不额外建立独立的只读归档存储域，而是通过 `ArchiveExportTask` 生成结构化归档产物；历史扫描、趋势与治理动作的在线查看仍复用原始查询接口，导出任务负责形成可复盘归档包。
- Rationale: 这样可以在首期满足“复盘归档”需求，同时避免引入第二套数据读取模型；在线查看与归档导出职责清晰分离。
- Alternatives considered:
  - 建立独立归档数据库或专用只读库：首期复杂度过高，且与现有查询能力重复。
  - 只保留在线列表不提供导出：难以满足审计留档与管理汇报的外部归档需求。
