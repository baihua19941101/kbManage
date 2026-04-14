#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
BACKEND_DIR="$ROOT_DIR/backend"
FRONTEND_DIR="$ROOT_DIR/frontend"
API_BASE_URL="${API_BASE_URL:-http://127.0.0.1:8888/api/v1}"

printf '\n[004-smoke] 1/4 backend tests\n'
(
  cd "$BACKEND_DIR"
  go test ./...
)

printf '\n[004-smoke] 2/4 frontend gitops tests\n'
(
  cd "$FRONTEND_DIR"
  npm run test -- --run src/features/gitops
)

printf '\n[004-smoke] 3/4 frontend lint\n'
(
  cd "$FRONTEND_DIR"
  npm run lint
)

printf '\n[004-smoke] 4/4 optional runtime probe (%s/healthz)\n' "$API_BASE_URL"
if command -v curl >/dev/null 2>&1; then
  if curl -fsS "${API_BASE_URL%/api/v1}/healthz" >/dev/null 2>&1; then
    echo "[004-smoke] runtime probe ok"
  else
    echo "[004-smoke] runtime probe skipped/failed (non-blocking)"
  fi
else
  echo "[004-smoke] curl not found, skip runtime probe"
fi

echo "[004-smoke] all done"
