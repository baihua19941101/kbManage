#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"

echo "[010] backend compile"
(cd "$ROOT_DIR/backend" && go test -run TestNonExistent -count=0 ./...)

echo "[010] backend contract"
(cd "$ROOT_DIR/backend" && go test ./tests/contract -run TestPlatformMarketplaceContract -count=1 -p 1)

echo "[010] backend integration"
(cd "$ROOT_DIR/backend" && go test ./tests/integration -run TestPlatformMarketplaceIntegration -count=1 -p 1)

echo "[010] frontend lint"
(cd "$ROOT_DIR/frontend" && npm run lint)

echo "[010] frontend build"
(cd "$ROOT_DIR/frontend" && npm run build)

echo "[010] smoke complete"
