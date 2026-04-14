# Research: 多集群 GitOps 与应用发布中心

> 研究输入同时参考了 kbManage 现有 001/002/003 设计边界，以及 Rancher / Fleet 官方文档中关于 GitRepo、Bundle 状态、Cluster Group、Target Customization 和多集群持续交付的公开语义。

## Decision 1: 004 采用独立 `gitops` 业务域，而不是继续扩展 001 的 `operation` 或 003 的 `workloadops`

- Decision: 在后端新增 `service/gitops`、`integration/delivery/*`、对应 handler/router/worker 以及前端 `features/gitops` 功能域；复用 001 的权限/审计/资源范围底座和 003 的异步动作执行经验，但不把 004 混入已有业务域。
- Rationale: 001 的 `operation` 更偏通用动作基座，003 的 `workloadops` 更偏运行中对象的 Day2 运维。004 的职责包括来源接入、目标分组、环境推进、差异/漂移、发布历史和回滚，已经形成独立的交付域边界。
- Alternatives considered:
  - 继续在 `operation` 域中横向扩展：会把来源管理、交付历史、环境推进和 GitOps 状态都塞进通用动作模型，语义过宽。
  - 继续在 `workloadops` 域中扩展：会把“声明式交付”和“运行时运维”混在一起，权限、页面和审计边界不清晰。

## Decision 2: 交付来源采用“平台抽象 + Fleet 风格运行时语义”的双层模型

- Decision: 平台的一等对象定义为 `DeliverySource`，首期支持两类来源：代码仓库来源与应用发布来源；运行时语义参考 Fleet 的 `GitRepo` 与发布来源管理方式，但平台对外不直接暴露底层 CRD 名称作为唯一用户模型。
- Rationale: 用户需要的是“交付来源”和“应用交付单元”，而不是直接学习底层控制器资源名。保留平台抽象可以让 UI 和权限模型保持稳定，同时又能借鉴 Fleet 在 Git 仓库路径、分支/修订、包来源等方面的成熟语义。
- Alternatives considered:
  - 直接把 Fleet `GitRepo`、`Bundle` 暴露为 UI 主模型：更贴近控制器，但对非 Fleet 背景用户不够友好，也会把底层实现细节耦合到产品界面。
  - 只支持 Git 仓库来源：无法覆盖用户明确提出的“代码仓库和发布来源”双来源要求。

## Decision 3: 目标分发与环境差异采用“Cluster Group + Environment Stage + Overlay”的组合建模

- Decision: 004 将 `ClusterTargetGroup` 作为可复用发布目标集合，将 `EnvironmentStage` 作为一等环境推进阶段，将 `ConfigurationOverlay` 作为环境或目标级覆盖模型；整体语义参考 Fleet 的 Cluster Group 与 Target Customization，但在产品层显式保留“环境顺序”概念。
- Rationale: Fleet 的 Cluster Group / selector / target customization 很适合表达“同一应用分发到不同集群或不同集群集合”的需求；而用户明确提出了“环境分层、配置覆盖和多环境推进过程”，因此需要在平台模型中加入比单纯 selector 更清晰的环境阶段概念。
- Alternatives considered:
  - 只用集群组，不建环境阶段：无法清晰表达测试、预发、生产的推进顺序。
  - 只用环境，不建目标组：不同环境下跨多个集群复用目标集合会重复配置，难以维护。
  - 每个环境单独维护一整套交付单元：会放大配置重复和版本漂移风险。

## Decision 4: 状态与漂移语义采用“期望状态 / 实际状态 / 同步结果 / 漂移状态”四层视图

- Decision: 004 的状态聚合视图同时展示交付单元级与目标级的 `desired state`、`live state`、最近同步结果、当前漂移状态和错误摘要；漂移语义对齐 Fleet `Bundle` / `BundleDeployment` / `GitRepo` 的状态思路，明确区分未同步、已修改、部分就绪、等待应用、不可达等状态。
- Rationale: GitOps 交付场景里，用户最关心的不是“是否发布过”，而是“当前是否对齐”“为什么没对齐”“差异发生在哪个环境/目标”。如果状态模型只返回一个整体布尔值，会掩盖部分成功、暂停、人工漂移和来源不可达等关键诊断信息。
- Alternatives considered:
  - 只展示最近一次动作成功/失败：无法体现持续对齐状态。
  - 只展示 diff：无法解释差异是否已经被暂停、待同步或目标不可达所阻断。

## Decision 5: 发布生命周期动作统一采用异步操作模型，环境推进作为同一动作域中的专门类型

- Decision: 安装、升级、同步、重新同步、暂停、恢复、回滚、卸载和环境推进统一落为 `DeliveryOperation` 异步动作模型；前端按动作类型呈现不同确认信息，后端统一处理状态流转、幂等和审计。
- Rationale: 004 的动作天然具有长耗时、部分成功和跨多目标聚合的特征，和 003 的工作负载动作一样，更适合异步提交与轮询查询。把环境推进视为动作域的一部分，可以复用统一的进度、审计和失败归一化能力。
- Alternatives considered:
  - 每种动作一个完全独立的执行模型：接口和状态模型重复度高。
  - 采用同步 API：不适合跨多集群、多环境和长耗时的交付动作。

## Decision 6: 发布历史与回滚基于“平台发布修订”而不是直接复用 Kubernetes 工作负载 revision

- Decision: 004 定义平台级 `ReleaseRevision`，以“来源修订 + 应用版本 + 配置版本 + 目标范围快照 + 环境推进结果”为一条可回滚的交付历史；回滚恢复的是交付声明与目标版本组合，而不是 003 中工作负载原生 revision。
- Rationale: GitOps 场景下真正需要回滚的是“交付声明”的组合，而不是单个工作负载控制器历史。直接复用 Kubernetes revision 会丢失来源修订、配置版本和环境推进上下文，也无法表达一次多集群交付的整体版本身份。
- Alternatives considered:
  - 直接依赖 003 的工作负载 revision：只覆盖运行时对象，不覆盖 Git/配置版本。
  - 不建立平台发布历史，只靠来源仓库历史：无法提供平台内的发布节奏、环境推进、失败原因和可回滚入口。

## Decision 7: 权限模型在现有范围隔离基础上细分到来源管理、同步发布、环境推进和回滚

- Decision: 004 在现有工作空间/项目范围模型上新增动作级权限语义，例如 `gitops:read`、`gitops:manage-source`、`gitops:sync`、`gitops:promote`、`gitops:rollback`、`gitops:override`；环境推进和回滚必须显式鉴权，并以环境范围为附加约束。
- Rationale: GitOps 与发布动作直接改变多环境运行状态，不能只沿用 001 的粗粒度读写权限。交付团队常见的真实治理需求是“能看全部测试环境、只能推进到预发、不能回滚生产”，因此需要动作级和环境级的组合授权。
- Alternatives considered:
  - 沿用单一发布权限：无法区分来源管理、同步、推进和回滚等敏感度不同的动作。
  - 只按工作空间授权，不区分环境：会导致生产环境与低风险环境无法做差异化控制。
