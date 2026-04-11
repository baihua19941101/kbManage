# Test Environment Load Report - 001-k8s-ops-platform

## 执行日期
- 2026-04-10（Asia/Shanghai）

## 测试环境
- Backend: Go 1.25，Gin + GORM，MySQL/Redis 本地开发实例
- Frontend: Node 20 + Vite 8
- 分支：`001-k8s-ops-platform-followup`

## 验证命令与结果
- `cd backend && go test ./...`：通过
- `cd frontend && npm run lint`：通过
- `cd frontend && npm run test -- --run`：通过（4 文件 / 7 用例）
- `cd frontend && npm run build`：通过

## 前端构建体积观测（路由懒加载后）
- `LoginPage`：约 1.23kB
- `ClusterOverviewPage`：约 13.92kB
- `ResourceListPage`：约 43.84kB
- `AuditEventPage`：约 215.95kB
- 共享 `client` chunk：约 606.55kB（高于默认 500k 告警阈值）

## 说明
- 本报告记录当前测试环境下可复现的构建与回归结果。
- API 维度压测通过脚本复现（见 `repro-perf-smoke.sh`），在目标环境启动后可直接执行。
