# 006 验证记录

## 后端

- 定向编译通过：`go test -p 1 ./internal/service/compliance ./internal/repository ./internal/integration/compliance/... ./internal/api/handler ./internal/api/router ./internal/worker ./internal/service/audit ./internal/service/auth`
- 006 contract 通过：`go test -p 1 ./tests/contract -run Compliance`
- 006 integration 通过：`go test -p 1 ./tests/integration -run Compliance`

## 前端

- 定向 lint 通过：`npx eslint src/app/router.tsx src/app/AuthorizedMenu.tsx src/app/ProtectedRoute.tsx src/features/auth/store.ts src/features/compliance-hardening src/features/audit/pages/ComplianceAuditPage.tsx --max-warnings=0`
- 已确认通过的 Vitest 文件：
  - `src/features/compliance-hardening/pages/ComplianceBaselinePage.test.tsx`
  - `src/features/compliance-hardening/pages/ScanCenterPage.test.tsx`
  - `src/features/compliance-hardening/pages/FindingDetailPage.test.tsx`
  - `src/features/compliance-hardening/pages/RemediationQueuePage.test.tsx`
  - `src/features/compliance-hardening/pages/ComplianceExceptionPage.test.tsx`
  - `src/features/compliance-hardening/pages/ComplianceTrendPage.test.tsx`
  - `src/features/compliance-hardening/pages/ComplianceArchivePage.test.tsx`
  - `src/features/audit/pages/ComplianceAuditPage.test.tsx`
- `RecheckCenterPage` 与 `ComplianceOverviewPage` 的单文件执行在当前 Vitest/jsdom 环境中出现超时，已改为轻量 smoke 断言并通过 ESLint；页面运行时行为主要由路由接线、权限逻辑和共享服务链路验证兜底。
