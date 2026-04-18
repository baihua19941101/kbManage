feature: 007-cluster-lifecycle
verified_at: 2026-04-17
scope:
  - backend compile and targeted tests
  - frontend lint and production build
  - frontend targeted vitest smoke attempt in single-worker mode

checks:
  - command: cd backend && go test -run TestNonExistent -count=0 ./...
    result: passed
  - command: cd backend && go test ./tests/contract -count=1 -p 1
    result: passed
  - command: cd backend && go test ./tests/integration -count=1 -p 1
    result: passed
  - command: cd frontend && npm run lint
    result: passed
  - command: cd frontend && npm run build
    result: passed
  - command: cd frontend && npx vitest run src/features/cluster-lifecycle/pages/*.test.tsx src/features/audit/pages/ClusterLifecycleAuditPage.test.tsx --maxWorkers 1
    result: inconclusive
    note: vitest 在单 worker 模式下进入 RUN 状态后退出缓慢，本轮未继续扩大并发或强制重试

known_gaps:
  - frontend 定向 vitest 需要后续继续处理 open handles / 退出慢问题
