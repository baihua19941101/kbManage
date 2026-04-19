# 008 PR 摘要

## 目标

新增“平台级备份恢复与灾备中心”，覆盖：
- 备份策略与恢复点管理
- 原地恢复、跨集群恢复、定向恢复与环境迁移
- 灾备演练计划、执行记录、报告与审计查询

## 主要改动

- 后端新增 `backuprestore` 领域模型、迁移、仓储、服务、执行器/校验器抽象、worker、API 和路由注册
- 前端新增 `backup-restore-dr` 功能域、页面、抽屉、hooks、服务封装与审计页面
- 接入 008 权限语义、菜单、路由、平台资源入口和全局审计查询
- 补齐 008 的 contract/integration 测试、前端页面测试、配置说明和验证材料

## 验证

- backend compile 通过
- backend contract/integration 定向测试通过
- frontend lint/build 通过
- frontend 定向 vitest 已单 worker 尝试，复现仓库现有长时间不退出问题，已记录

## 风险

- 前端定向 `vitest` 仍受仓库现有测试环境问题影响
- 当前尚未推送远程、未开 PR、未合并
