# kbManage

多集群 Kubernetes 可视化管理平台（开发中）。

## 目录结构

```text
backend/   # Go + Gin 后端
frontend/  # React + Vite 前端
specs/     # 规格、计划与任务
artifacts/ # 交付与治理证据
```

## 环境要求

- Go 1.25+
- Node.js 20+
- MySQL 8+
- Redis 8+

## 依赖镜像（中国网络）

后端（Go）：

```bash
go env -w GOPROXY=https://goproxy.cn,direct
```

前端（npm）：

```bash
cd frontend
npm config set registry https://registry.npmmirror.com
```

## 后端配置（单配置文件）

后端所有可配置项统一在配置文件管理：

- 默认配置文件：`backend/config/config.dev.yaml`
- 模板文件：`backend/config/config.example.yaml`
- 可通过环境变量指定文件：`CONFIG_FILE=/path/to/config.yaml`

配置项包括：

- `server.http_addr`
- `mysql.host / port / user / password / database / parse_time`
- `redis.addr / password / db`
- `security.jwt_secret / access_token_ttl / refresh_token_ttl`
- `observability.metrics.base_url / auth_secret_ref / timeout`
- `observability.logs.base_url / auth_secret_ref / timeout`
- `observability.alerts.base_url / auth_secret_ref / timeout`
- `observability.cache.query_ttl / sync_interval`
- `workloadops.actions.default_timeout / idempotency_ttl`
- `workloadops.batch.max_targets / concurrency`
- `workloadops.terminal.idle_timeout / token_ttl`
- `workloadops.audit.retention_days`
- `gitops.sources.git_timeout / package_timeout / healthcheck_timeout`
- `gitops.sync.default_timeout / max_parallel_targets / drift_scan_interval`
- `gitops.diff.cache_ttl / max_entries / ignore_last_applied_annotation`
- `gitops.release.history_retention_days / max_concurrent_operations / stage_progress_deadline`
- `gitops.audit.retention_days / redact_secrets`
- `securityPolicy.policy.distribution_timeout / max_parallel_targets`
- `securityPolicy.exception.max_duration_hours / expiry_scan_interval`
- `securityPolicy.cache.distribution_ttl / exception_ttl`
- `securityPolicy.audit.retention_days`
- `compliance.baselines.source_ref / snapshot_cache_ttl`
- `compliance.scan.default_timeout / max_parallel_executions / scheduler_interval`
- `compliance.export.retention_ttl`
- `compliance.audit.retention_days`

说明：

- 配置文件已包含中文注释。
- 支持少量环境变量临时覆盖（如 `MYSQL_HOST`、`REDIS_ADDR`、`HTTP_ADDR`），用于 CI 或调试。

## 前端配置（env 文件）

前端统一使用 `env` 文件配置：

- 开发默认：`frontend/.env.development`
- 模板：`frontend/.env.example`

关键项：

- `VITE_API_BASE_URL`：后端 API 前缀
- `VITE_HOST`：前端 dev server host
- `VITE_PORT`：前端 dev server 端口（必须通过该变量配置，不在脚本中硬编码）
- `VITE_GITOPS_OPERATION_POLL_INTERVAL`：GitOps 动作轮询间隔（毫秒）
- `VITE_GITOPS_DIFF_REFRESH_INTERVAL`：GitOps 差异面板刷新间隔（毫秒）
- `VITE_COMPLIANCE_REFRESH_INTERVAL`：合规与加固页面刷新间隔（毫秒）

## 启动方式

### 1) 启动后端

```bash
cd backend
# 使用默认配置文件 config/config.dev.yaml
go run ./cmd/server
```

自定义配置文件：

```bash
cd backend
CONFIG_FILE=./config/config.dev.yaml go run ./cmd/server
```

### 2) 启动前端

```bash
cd frontend
npm install
npm run dev
```

如果需要改端口，修改 `frontend/.env.development` 中的 `VITE_PORT` 即可，例如：

```env
VITE_PORT=5180
```

## 常用命令

```bash
# 仓库根目录
make backend-test
make frontend-test
make lint
make test
```

前端完整验证命令（与当前任务基线一致）：

```bash
cd frontend
npm run lint
npm run test -- --run
npm run build
```

## 当前状态说明

- 002-observability-center 已完成 US1/US2/US3（统一观测入口、告警治理闭环、工作空间/项目级权限隔离）。
- 可观测后端接口已按读写动作区分权限：`observability:read`（读取）与 `observability:write`（治理动作）。
- 前端已补齐可观测权限空态/只读态/权限回收处理，并通过 Vitest 与 ESLint 验证。
- 003-workload-operations-control-plane 已完成 US1/US2/US3（单资源诊断、动作执行与回滚恢复、权限隔离与高风险审计闭环），当前进入 Final Phase 收尾与交付准备。
- 004-gitops-and-release 已完成 US1/US2/US3，当前进入 Final Phase（文档收尾、验证基线与 PR 就绪材料）。
- GitOps 审计查询入口：`/audit-events/gitops`。
- 005-security-and-policy 已进入 implement 阶段，已落地策略中心 US1 前后端主干与治理证据。

## 002 可观测联调要点

- 后端默认读取 `backend/config/config.dev.yaml`，按需配置：
  - `observability.metrics.base_url`
  - `observability.logs.base_url`
  - `observability.alerts.base_url`
  - `observability.cache.query_ttl`
  - `observability.cache.sync_interval`
- 前端入口：
  - `/observability`（总览）
  - `/observability/logs`、`/observability/events`、`/observability/metrics`
  - `/observability/alerts`、`/observability/alert-rules`、`/observability/silences`
- 如需复现实验基线，执行：

```bash
bash artifacts/002-observability-center/repro-observability-smoke.sh
```

## 006 合规与加固联调要点

- 后端默认读取以下配置：
  - `compliance.baselines.source_ref`
  - `compliance.scan.default_timeout`
  - `compliance.scan.max_parallel_executions`
  - `compliance.scan.scheduler_interval`
  - `compliance.export.retention_ttl`
  - `compliance.audit.retention_days`
- 前端入口：
  - `/compliance-hardening/baselines`
  - `/compliance-hardening/scans`
  - `/compliance-hardening/remediation`
  - `/compliance-hardening/exceptions`
  - `/compliance-hardening/rechecks`
  - `/compliance-hardening/overview`
  - `/compliance-hardening/trends`
  - `/compliance-hardening/archive`
  - `/audit-events/compliance`
