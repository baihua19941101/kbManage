#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"

echo "[1/3] backend test"
(
  cd "$ROOT_DIR/backend"
  go test ./...
)

echo "[2/3] frontend test"
(
  cd "$ROOT_DIR/frontend"
  npm test
)

echo "[3/3] frontend lint"
(
  cd "$ROOT_DIR/frontend"
  npm run lint
)

echo "observability smoke validation finished"
