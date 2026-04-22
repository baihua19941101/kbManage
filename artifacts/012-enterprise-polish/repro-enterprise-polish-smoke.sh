#!/usr/bin/env bash
set -euo pipefail

echo "[012] backend compile"
(cd backend && go test -run TestNonExistent -count=0 ./...)

echo "[012] backend contract"
(cd backend && go test ./tests/contract -run TestEnterprisePolishContract -count=1 -p 1)

echo "[012] backend integration"
(cd backend && go test ./tests/integration -run TestEnterprisePolishIntegration -count=1 -p 1)

echo "[012] frontend lint/build"
(cd frontend && npm run lint && npm run build)
