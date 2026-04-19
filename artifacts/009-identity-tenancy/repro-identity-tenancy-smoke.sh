#!/usr/bin/env bash
set -euo pipefail

echo "[009] backend compile check"
( cd backend && go test -run TestNonExistent -count=0 ./... )

echo "[009] backend contract check"
( cd backend && go test -run TestIdentityTenancyContract -count=1 -p 1 ./tests/contract )

echo "[009] backend integration check"
( cd backend && go test -run TestIdentityTenancyIntegration -count=1 -p 1 ./tests/integration )

echo "[009] frontend lint"
( cd frontend && npm run lint )

echo "[009] frontend build"
( cd frontend && npm run build )

echo "[009] note: vitest single-worker stability remains a known repo-wide issue"
