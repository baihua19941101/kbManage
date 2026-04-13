# 003 验证基线（2026-04-12）

## 范围

- US3：权限隔离与高风险审计闭环
- 覆盖后端 Contract/Integration、前端 Vitest 与变更文件 ESLint

## 验证命令与结果

1. `cd backend && go test ./tests/contract -run TestWorkloadOpsContract_AccessControlDeniedWithoutScopeBinding -count=1`
- 结果：`ok kbmanage/backend/tests/contract`

2. `cd backend && go test ./tests/integration -run 'TestWorkloadOpsScopeAuthorizationIsolationAndRevocation|TestWorkloadOpsAuditIntegration_WritesActionAndTerminalEvents' -count=1`
- 结果：`ok kbmanage/backend/tests/integration`

3. `cd frontend && npm test -- src/features/workload-ops/pages/WorkloadOperationsAccessGate.test.tsx src/features/workload-ops/components/RollbackDialog.test.tsx src/features/workload-ops/pages/WorkloadOperationsPage.test.tsx`
- 结果：`3 files passed, 6 tests passed`
- 备注：存在 Ant Design 运行时 deprecation warning（不影响结果）

4. `cd frontend && npx eslint src/app/AuthorizedMenu.tsx src/features/auth/store.ts src/features/workload-ops/pages/WorkloadOperationsPage.tsx src/features/workload-ops/components/TerminalSessionDrawer.tsx src/features/workload-ops/components/RollbackDialog.tsx src/features/workload-ops/pages/WorkloadOperationsAccessGate.test.tsx src/features/workload-ops/components/RollbackDialog.test.tsx`
- 结果：通过，无错误

## 结论

- US3 范围内的后端权限校验、审计行为与前端门控交互均已达到当前自动化测试基线。
