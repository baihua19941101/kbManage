#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"

echo "[011] backend compile"
(cd "$ROOT_DIR/backend" && go test -run TestNonExistent -count=0 ./...)

echo "[011] frontend lint"
(cd "$ROOT_DIR/frontend" && npm run lint)

echo "[011] frontend build"
(cd "$ROOT_DIR/frontend" && npm run build)

echo "[011] smoke complete"
