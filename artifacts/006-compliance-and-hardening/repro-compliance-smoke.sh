#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${BASE_URL:-http://127.0.0.1:8888/api/v1}"
TOKEN="${TOKEN:-}"
if [[ -z "$TOKEN" ]]; then
  echo "TOKEN is required" >&2
  exit 1
fi

curl -sS -H "Authorization: Bearer $TOKEN" "$BASE_URL/compliance/baselines" >/dev/null
curl -sS -H "Authorization: Bearer $TOKEN" "$BASE_URL/compliance/scan-profiles" >/dev/null
curl -sS -H "Authorization: Bearer $TOKEN" "$BASE_URL/compliance/scans" >/dev/null
curl -sS -H "Authorization: Bearer $TOKEN" "$BASE_URL/compliance/overview" >/dev/null
curl -sS -H "Authorization: Bearer $TOKEN" "$BASE_URL/audit/compliance/events" >/dev/null
echo "compliance smoke passed"
