# 006 Verification Strategy

- backend tests: use `go test -p 1` on affected packages only
- frontend tests: use `npm run test -- --run <target> --maxWorkers=1`
- lint: prefer file-scoped or feature-scoped verification when possible
- avoid running backend and frontend heavy tests concurrently
- avoid full-repo parallel test sweeps until implementation stabilizes
