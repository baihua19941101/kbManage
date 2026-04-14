# PR 摘要（004-gitops-and-release）

## 主题
完成 004「多集群 GitOps 与应用发布中心」US1/US2/US3，实现来源建模、发布生命周期、权限隔离与审计闭环。

## 核心变更
- 后端：
  - 新增 GitOps 领域模型、仓储、服务、路由与 handler。
  - 完成发布动作执行链路（install/sync/resync/upgrade/promote/rollback/pause/resume/uninstall）。
  - 新增 GitOps 细粒度鉴权中间件（source/target-group/unit/action/operation）。
  - 新增 `gitops.*` 审计动作分类与写入。
- 前端：
  - 完成 GitOps 概览、交付单元详情、差异面板、发布历史、动作抽屉、回滚交互。
  - 新增权限回收锁定态与 GitOps access gate 测试。
  - 新增 GitOps 审计查询页 `/audit-events/gitops`。
- 测试：
  - 新增 GitOps contract/integration 与前端 Vitest 场景覆盖。

## 验证结果
- `cd backend && go test ./...` ✅
- `cd frontend && npm run test -- --run src/features/gitops` ✅
- `cd frontend && npm run lint` ✅

## 宪章与治理
- 开发分支：`004-gitops-and-release`
- 备份证据：`artifacts/004-gitops-and-release/mysql-backup-20260413-155243.sql`
- 国内源配置：`GOPROXY=https://goproxy.cn,direct`、`https://registry.npmmirror.com`
- 合并规则：需用户明确批准后合并
