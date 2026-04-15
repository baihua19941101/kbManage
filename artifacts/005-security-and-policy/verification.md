# 005 验证基线

Date: 2026-04-14
Branch: 005-security-and-policy

## 后端

命令：

```bash
cd backend
go test -p 1 ./...
```

结果：PASS

## 前端（低并发）

命令：

```bash
cd frontend
npm run test -- --run \
  src/features/security-policy/pages/PolicyCenterPage.test.tsx \
  src/features/security-policy/components/PolicyScopeDrawer.test.tsx \
  src/features/security-policy/pages/PolicyRolloutPage.test.tsx \
  src/features/security-policy/components/ExceptionReviewDrawer.test.tsx \
  src/features/security-policy/pages/ViolationCenterPage.test.tsx \
  src/features/audit/pages/SecurityPolicyAuditPage.test.tsx \
  --maxWorkers=1
```

结果：PASS（6 files, 10 tests）

## 前端 Lint

命令（节选 005 相关范围）：

```bash
cd frontend
npx eslint src/features/security-policy/** src/features/audit/pages/SecurityPolicyAuditPage.tsx src/services/securityPolicy.ts src/services/audit.ts src/services/api/types.ts src/app/router.tsx src/app/AuthorizedMenu.tsx src/features/auth/store.ts
```

结果：PASS

## 说明

- 运行期间存在 Ant Design deprecation warning（Space/Drawer/Progress），不影响当前功能正确性。
- 测试已按低并发执行，避免内存峰值过高。
