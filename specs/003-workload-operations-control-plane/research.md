# Research: 多集群 Kubernetes 工作负载运维控制面

## Decision 1: 003 采用独立 `workloadops` 业务域，而不是继续扩展 001 的 `operation` 或 002 的 `observability`

- Decision: 在后端新增 `service/workloadops`、对应 handler/router/worker 以及前端 `features/workload-ops` 功能域；复用 001 的操作队列与审计底座、复用 002 的日志/资源上下文能力，但不把 003 混入已有域。
- Rationale: 001 的 `operation` 更偏通用动作执行基座，002 的 `observability` 更偏读路径和治理视图。003 同时包含资源诊断聚合、批量动作、回滚和终端会话，职责已经超出原有域的自然边界。独立域更利于隔离权限、页面、契约和测试。
- Alternatives considered:
  - 继续在 `operation` 域中横向扩展：会把发布历史、终端会话和批量动作都塞入通用操作模型，语义过宽。
  - 继续在 `observability` 域中扩展：会把“读路径”和“写路径/高风险动作”混在一起，权限和审计边界不清晰。

## Decision 2: 终端访问采用“浏览器 -> 平台后端 -> Kubernetes API”的短生命周期代理模型

- Decision: 平台通过短生命周期终端会话代理到 Kubernetes `exec` 通道；浏览器不直接访问集群。首期终端审计只记录会话建立、关闭、目标容器、操作者、持续时长和结束原因，不记录完整命令与终端输出正文。
- Rationale: 终端能力必须经过平台统一鉴权和审计，才能满足多租户与追责要求。只记录会话元数据能满足首期审计闭环，同时避免把 003 扩展成命令录屏或长期留痕系统。
- Alternatives considered:
  - 浏览器直连集群：绕过平台鉴权与审计，违反控制面定位。
  - 记录完整命令与终端输出：合规与存储成本更高，也会显著扩大首期范围。
  - 长连接常驻会话服务：复杂度高于首期目标，且对断线回收和资源占用要求更高。

## Decision 3: 发布历史与回滚能力基于 Kubernetes 原生 revision 语义，不引入独立发布系统

- Decision: `Deployment` 的历史版本基于其控制的 `ReplicaSet` revision 语义；`StatefulSet` 与 `DaemonSet` 的历史版本基于 `ControllerRevision`。回滚仅针对工作负载的 Pod 模板相关变更进行恢复，不把扩缩容等非模板修改当成新 revision。
- Rationale: 官方 Kubernetes 行为已经提供稳定的一致语义，平台应尽量复用，而不是再发明一套版本体系。这样可以保证平台展示的历史、回滚结果与集群真实控制器状态一致。
- Alternatives considered:
  - 平台独立维护“发布版本”表：与 Kubernetes 原生 revision 语义容易漂移，增加同步复杂度。
  - 仅支持 Deployment 回滚：虽然实现简单，但会削弱 Rancher 风格工作负载运维一致性。
  - 把扩缩容纳入 revision 历史：不符合 Kubernetes 原生 revision 规则，也会误导用户。

## Decision 4: 单资源动作与批量动作统一走异步任务模型，批量动作保留独立子项结果

- Decision: 003 的单资源动作与批量动作都使用异步提交、进度轮询和终态查询模型。批量动作在逻辑上拆成一个 `BatchOperationTask` 与多个子项，每个目标对象保留独立状态、结果和失败原因；平台层采用受控并发窗口执行，而不是一次性全并发。
- Rationale: 工作负载运维动作具有延迟、失败和部分成功的天然特征。异步任务模型能复用 001 已有的队列/状态流转基础，也更适合展示影响范围、执行进度和部分失败。
- Alternatives considered:
  - 全同步接口：不适合长耗时动作，也不适合前端持续刷新状态。
  - 批量动作只返回整体结果：会掩盖部分失败，不符合审计与复盘需求。
  - 无并发限制的全量并行：对集群和平台控制面冲击过大。

## Decision 5: 003 的动作契约采用统一提交接口，UI 层再做动作类型细分

- Decision: 后端接口层提供统一的动作提交契约，通过 `actionType` 区分 `scale`、`restart`、`redeploy`、`replace-instance`、`rollback` 等类型；前端页面按具体动作呈现不同确认表单和风险提示。
- Rationale: 当前 001 已有通用 `operation` 提交模式，003 继续采用统一动作模型更利于共享幂等、队列、审计和状态流转能力，同时避免为每一种动作重复定义大量接口骨架。
- Alternatives considered:
  - 为每一种动作单独设计一套提交流程：可读性更强，但重复接口和状态模型过多。
  - 把所有动作继续压进 001 旧接口：无法清晰表达 003 新增的批量任务、回滚和终端语义。

## Decision 6: 权限模型在现有 RBAC 基础上细分为读取、终端、高风险动作和批量变更四类能力

- Decision: 003 在现有工作空间/项目范围模型上新增 `workloadops:read`、`workloadops:execute`、`workloadops:terminal`、`workloadops:rollback`、`workloadops:batch` 等动作级权限语义，并要求终端、回滚和批量高风险动作显式鉴权。
- Rationale: 仅用 `operation:execute` 无法细分终端和回滚等更高敏感动作。动作级权限更符合“能看不一定能进终端，能重启不一定能批量回滚”的实际治理需求。
- Alternatives considered:
  - 沿用 `operation:execute` 单一权限：粒度不足。
  - 把终端权限并入 `observability:write`：会混淆读域和写域边界。
