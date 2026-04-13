# PR 摘要（003-workload-operations-control-plane）

## 背景

003 目标是补齐对标 Rancher 的工作负载运维闭环，覆盖：
- 单资源诊断入口（US1）
- 动作执行与回滚恢复（US2）
- 权限隔离与高风险审计闭环（US3）

## 本次新增/完善

1. 后端权限与范围
- 新增/完善 `workloadops:read|execute|terminal|rollback|batch` 权限语义。
- 动作接口按 `actionType` 动态区分回滚与普通动作权限。
- 服务层统一执行 cluster + workspace/project 范围校验。

2. 后端审计
- 补齐 `workloadops.*` 审计动作链路，覆盖动作、批量、回滚、终端会话。
- 终端审计边界明确：仅会话元数据（建立/关闭、目标容器、操作者、持续时长、结束原因），不记录命令正文与终端输出正文。

3. 前端门控与权限回收
- 导航与动作入口按能力分级门控（读取、执行、终端、回滚、批量）。
- 页面、终端抽屉、回滚弹窗支持只读态与权限回收锁定。

4. 规格与交付文档同步
- 003 `spec.md / plan.md / tasks.md` 已更新到 US3 完成状态。
- 新增验证基线与复现脚本：
  - `verification.md`
  - `quickstart-validation.md`
  - `repro-workloadops-smoke.sh`

## 验证

- 后端 Contract：通过
- 后端 Integration：通过
- 前端 Vitest（workload-ops 关键用例）：通过
- 变更文件 ESLint：通过

详细命令见：`artifacts/003-workload-operations-control-plane/verification.md`
