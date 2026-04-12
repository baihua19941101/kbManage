# 002 Observability Verification Baseline

- Feature: `002-observability-center`
- Date: 2026-04-12
- Branch: `002-observability-center`
- Scope: US1 + US2 + US3 + Final Phase baseline validation

## Environment

- Backend: Go 1.25+
- Frontend: Node.js 20+ (current runtime supports Vitest 3.x)
- MySQL: 8.x (dev)
- Redis: 8.x (dev)
- Config file: `backend/config/config.dev.yaml`

## Executed Validation Commands

### Backend

```bash
cd backend
go test ./...
```

Result: PASS

### Frontend

```bash
cd frontend
npm test
npm run lint
```

Result: PASS

## Key Functional Baselines

1. Unified observability read APIs (`overview/logs/events/metrics/context`) are routed and protected by read scope checks.
2. Alert governance APIs (`alerts/rules/targets/silences`) are available and split by read/write authorization.
3. Scope isolation (workspace/project) is enforced in backend middleware and scope service mapping.
4. Frontend observability pages handle:
   - unauthorized empty state
   - readonly action gating
   - permission-revoked warning after action/API denial
5. Audit events capture observability critical actions and sync activity.

## Known Non-Blocking Notes

- Vitest output includes Ant Design deprecation warnings (`Space direction`, `Drawer width`, etc.); these do not block test pass.
- SQLite test environment does not start observability sync worker to avoid lock contention (production DB unaffected).
