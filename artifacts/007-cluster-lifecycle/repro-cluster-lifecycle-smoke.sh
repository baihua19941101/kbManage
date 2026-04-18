#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"

echo "[007] backend compile smoke"
(cd "$ROOT_DIR/backend" && go test -run TestNonExistent -count=0 ./...)

echo "[007] backend targeted tests"
(cd "$ROOT_DIR/backend" && go test ./tests/contract -run TestClusterLifecycleContract_ImportRegisterListDetail -count=1 -p 1)
(cd "$ROOT_DIR/backend" && go test ./tests/integration -run TestClusterLifecycleIntegration_CreateUpgradeRetireFlow -count=1 -p 1)

echo "[007] frontend lint/build"
(cd "$ROOT_DIR/frontend" && npm run lint)
(cd "$ROOT_DIR/frontend" && npm run build)

echo "[007] smoke complete"
