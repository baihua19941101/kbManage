# 003 Quickstart 验证记录（2026-04-12）

## 前置条件

- 当前分支：`003-workload-operations-control-plane`
- 已完成数据库备份：`artifacts/003-workload-operations-control-plane/mysql-backup-20260412-150123.sql`
- 依赖镜像：
  - Go: `GOPROXY=https://goproxy.cn,direct`
  - npm: `https://registry.npmmirror.com`

## 验证步骤

1. 后端 Contract 校验
- 执行：`cd backend && go test ./tests/contract -run TestWorkloadOpsContract_AccessControlDeniedWithoutScopeBinding -count=1`
- 结果：通过

2. 后端 Integration 校验
- 执行：`cd backend && go test ./tests/integration -run 'TestWorkloadOpsScopeAuthorizationIsolationAndRevocation|TestWorkloadOpsAuditIntegration_WritesActionAndTerminalEvents' -count=1`
- 结果：通过

3. 前端页面权限门控校验
- 执行：`cd frontend && npm test -- src/features/workload-ops/pages/WorkloadOperationsAccessGate.test.tsx src/features/workload-ops/components/RollbackDialog.test.tsx src/features/workload-ops/pages/WorkloadOperationsPage.test.tsx`
- 结果：通过

4. 前端静态检查
- 执行：`cd frontend && npx eslint src/app/AuthorizedMenu.tsx src/features/auth/store.ts src/features/workload-ops/pages/WorkloadOperationsPage.tsx src/features/workload-ops/components/TerminalSessionDrawer.tsx src/features/workload-ops/components/RollbackDialog.tsx src/features/workload-ops/pages/WorkloadOperationsAccessGate.test.tsx src/features/workload-ops/components/RollbackDialog.test.tsx`
- 结果：通过

## 结果摘要

- US3 的“权限隔离 + 高风险审计闭环”在当前 quickstart 路径下可复现通过。
