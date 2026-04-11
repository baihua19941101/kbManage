#!/usr/bin/env bash
set -euo pipefail

# Reproducible smoke load script for local/dev verification.
# Usage:
#   API_BASE="http://localhost:8888/api/v1" TOKEN="xxx" ./artifacts/001-k8s-ops-platform/repro-perf-smoke.sh

API_BASE="${API_BASE:-http://localhost:8888/api/v1}"
TOKEN="${TOKEN:-}"
ROUNDS="${ROUNDS:-50}"

if [[ -z "${TOKEN}" ]]; then
  echo "TOKEN is required. Example:"
  echo '  API_BASE="http://localhost:8888/api/v1" TOKEN="<jwt>" ./artifacts/001-k8s-ops-platform/repro-perf-smoke.sh'
  exit 1
fi

echo "[info] API_BASE=${API_BASE}"
echo "[info] ROUNDS=${ROUNDS}"

req() {
  local url="$1"
  curl -sS -o /dev/null -w "%{http_code} %{time_total}\n" \
    -H "Authorization: Bearer ${TOKEN}" \
    "${url}"
}

run_case() {
  local name="$1"
  local url="$2"
  local tmp
  tmp="$(mktemp)"
  echo "[case] ${name}"
  for _ in $(seq 1 "${ROUNDS}"); do
    req "${url}" >> "${tmp}"
  done
  awk '
    BEGIN { n=0; sum=0; ok=0; }
    {
      n+=1;
      if ($1 ~ /^2/) ok+=1;
      sum+=$2;
    }
    END {
      avg=(n>0?sum/n:0);
      printf("  total=%d ok=%d avg_time=%.4fs\n", n, ok, avg);
    }
  ' "${tmp}"
  rm -f "${tmp}"
}

run_case "clusters list" "${API_BASE}/clusters"
run_case "resources list (cluster=1)" "${API_BASE}/clusters/1/resources?limit=50"
run_case "audit query (90d window)" "${API_BASE}/audits/events?startAt=2026-01-01T00:00:00Z&endAt=2026-03-31T23:59:59Z&limit=100"

echo "[done] smoke performance check complete."
