#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"

echo "[003-smoke] root: ${ROOT_DIR}"

echo "[003-smoke] backend contract"
(
  cd "${ROOT_DIR}/backend"
  go test ./tests/contract -run TestWorkloadOpsContract_AccessControlDeniedWithoutScopeBinding -count=1
)

echo "[003-smoke] backend integration"
(
  cd "${ROOT_DIR}/backend"
  go test ./tests/integration -run 'TestWorkloadOpsScopeAuthorizationIsolationAndRevocation|TestWorkloadOpsAuditIntegration_WritesActionAndTerminalEvents' -count=1
)

echo "[003-smoke] frontend vitest"
(
  cd "${ROOT_DIR}/frontend"
  npm test -- src/features/workload-ops/pages/WorkloadOperationsAccessGate.test.tsx src/features/workload-ops/components/RollbackDialog.test.tsx src/features/workload-ops/pages/WorkloadOperationsPage.test.tsx
)

echo "[003-smoke] frontend eslint"
(
  cd "${ROOT_DIR}/frontend"
  npx eslint src/app/AuthorizedMenu.tsx src/features/auth/store.ts src/features/workload-ops/pages/WorkloadOperationsPage.tsx src/features/workload-ops/components/TerminalSessionDrawer.tsx src/features/workload-ops/components/RollbackDialog.tsx src/features/workload-ops/pages/WorkloadOperationsAccessGate.test.tsx src/features/workload-ops/components/RollbackDialog.test.tsx
)

echo "[003-smoke] done"
