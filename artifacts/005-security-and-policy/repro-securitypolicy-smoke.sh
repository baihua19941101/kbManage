#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")/../.." && pwd)"

echo "[005 smoke] backend tests"
(
  cd "$ROOT_DIR/backend"
  go test -p 1 ./...
)

echo "[005 smoke] frontend lint"
(
  cd "$ROOT_DIR/frontend"
  npx eslint \
    src/features/security-policy/** \
    src/features/audit/pages/SecurityPolicyAuditPage.tsx \
    src/services/securityPolicy.ts \
    src/services/audit.ts \
    src/services/api/types.ts \
    src/app/router.tsx \
    src/app/AuthorizedMenu.tsx \
    src/features/auth/store.ts
)

echo "[005 smoke] frontend tests (low parallelism)"
(
  cd "$ROOT_DIR/frontend"
  npm run test -- --run \
    src/features/security-policy/pages/PolicyCenterPage.test.tsx \
    src/features/security-policy/components/PolicyScopeDrawer.test.tsx \
    src/features/security-policy/pages/PolicyRolloutPage.test.tsx \
    src/features/security-policy/components/ExceptionReviewDrawer.test.tsx \
    src/features/security-policy/pages/ViolationCenterPage.test.tsx \
    src/features/audit/pages/SecurityPolicyAuditPage.test.tsx \
    --maxWorkers=1
)

echo "[005 smoke] done"
