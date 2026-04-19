#!/usr/bin/env bash

set -euo pipefail

BASE_URL="${BASE_URL:-http://127.0.0.1:8080/api/v1}"
TOKEN="${TOKEN:-}"
WORKSPACE_ID="${WORKSPACE_ID:-1}"
PROJECT_ID="${PROJECT_ID:-1}"

if [[ -z "$TOKEN" ]]; then
  echo "TOKEN is required"
  exit 1
fi

auth_header=(-H "Authorization: Bearer ${TOKEN}" -H "Content-Type: application/json")

echo "[1/6] create backup policy"
policy_json=$(curl -fsS "${auth_header[@]}" \
  -X POST "${BASE_URL}/backup-restore/policies" \
  -d "{
    \"name\":\"smoke-backup-policy\",
    \"scopeType\":\"namespace\",
    \"scopeRef\":\"orders-prod\",
    \"workspaceId\":${WORKSPACE_ID},
    \"projectId\":${PROJECT_ID},
    \"executionMode\":\"manual\",
    \"retentionRule\":\"14d\",
    \"consistencyLevel\":\"application-consistent\",
    \"status\":\"active\"
  }")
policy_id=$(printf '%s' "$policy_json" | sed -n 's/.*\"id\":\([0-9][0-9]*\).*/\1/p')
echo "policy_id=${policy_id}"

echo "[2/6] run backup policy"
restore_point_json=$(curl -fsS "${auth_header[@]}" \
  -X POST "${BASE_URL}/backup-restore/policies/${policy_id}/run")
restore_point_id=$(printf '%s' "$restore_point_json" | sed -n 's/.*\"id\":\([0-9][0-9]*\).*/\1/p')
echo "restore_point_id=${restore_point_id}"

echo "[3/6] create restore job"
restore_job_json=$(curl -fsS "${auth_header[@]}" \
  -X POST "${BASE_URL}/backup-restore/restore-jobs" \
  -d "{
    \"restorePointId\":${restore_point_id},
    \"jobType\":\"cross-cluster-restore\",
    \"sourceEnvironment\":\"prod\",
    \"targetEnvironment\":\"dr-site\",
    \"scopeSelection\":{\"namespaces\":[\"orders-prod\"]}
  }")
restore_job_id=$(printf '%s' "$restore_job_json" | sed -n 's/.*\"id\":\([0-9][0-9]*\).*/\1/p')
echo "restore_job_id=${restore_job_id}"

echo "[4/6] validate restore job"
curl -fsS "${auth_header[@]}" \
  -X POST "${BASE_URL}/backup-restore/restore-jobs/${restore_job_id}/validate" >/dev/null

echo "[5/6] create drill plan"
plan_json=$(curl -fsS "${auth_header[@]}" \
  -X POST "${BASE_URL}/backup-restore/drills/plans" \
  -d "{
    \"name\":\"smoke-drill-plan\",
    \"workspaceId\":${WORKSPACE_ID},
    \"projectId\":${PROJECT_ID},
    \"scopeSelection\":{\"namespaces\":[\"orders-prod\"]},
    \"rpoTargetMinutes\":15,
    \"rtoTargetMinutes\":30,
    \"roleAssignments\":[\"sre\",\"biz-owner\"],
    \"cutoverProcedure\":[\"freeze writes\",\"switch traffic\"],
    \"validationChecklist\":[\"verify api\",\"verify jobs\"]
  }")
plan_id=$(printf '%s' "$plan_json" | sed -n 's/.*\"id\":\([0-9][0-9]*\).*/\1/p')
echo "plan_id=${plan_id}"

echo "[6/6] query backup restore audit"
curl -fsS -H "Authorization: Bearer ${TOKEN}" \
  "${BASE_URL}/audit/backup-restore/events" >/dev/null

echo "backup restore smoke passed"
