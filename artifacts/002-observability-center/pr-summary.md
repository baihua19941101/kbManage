# PR Summary - 002 可观测中心

## 变更范围

本 PR 完成 `002-observability-center` 全部用户故事与收尾项：

1. US1：统一可观测入口
- 总览、日志、事件、指标、资源上下文页面与后端接口。
- 资源/集群入口联动至可观测上下文。

2. US2：告警治理闭环
- 告警中心、规则治理、通知目标、静默窗口。
- 告警确认、处理记录、运行态同步与审计接入。

3. US3：权限隔离与动作级门控
- 后端对 observability read/write 接口统一授权。
- 范围映射到 workspace/project 并执行隔离。
- 前端支持未授权空态、只读态、权限回收后阻断。

4. Final Phase 收尾
- observability 前端服务 query 构造逻辑收敛复用。
- README、配置注释、验证基线、quickstart 校验与 smoke 脚本补齐。

## 关键文件

- 后端路由与中间件：
  - `backend/internal/api/router/observability_routes.go`
  - `backend/internal/api/middleware/authorization.go`
- 后端服务：
  - `backend/internal/service/observability/*`
  - `backend/internal/service/auth/*`
  - `backend/internal/service/audit/*`
- 前端页面与权限门控：
  - `frontend/src/features/observability/pages/*`
  - `frontend/src/app/ProtectedRoute.tsx`
  - `frontend/src/app/AuthorizedMenu.tsx`
  - `frontend/src/features/auth/store.ts`

## 测试与验证

- Backend: `cd backend && go test ./...` ✅
- Frontend: `cd frontend && npm test` ✅
- Frontend lint: `cd frontend && npm run lint` ✅
- 一键复现：`bash artifacts/002-observability-center/repro-observability-smoke.sh` ✅

## 治理与证据

- 分支/远程/镜像证据：
  - `artifacts/002-observability-center/branch-check.txt`
  - `artifacts/002-observability-center/mirror-and-remote-check.txt`
- 数据库备份证据：
  - `artifacts/002-observability-center/backup-manifest.txt`
  - `artifacts/002-observability-center/mysql-backup-20260411-214819.sql`
- 验证记录：
  - `artifacts/002-observability-center/verification.md`
  - `artifacts/002-observability-center/quickstart-validation.md`

## 合并说明

- 当前状态：**待用户明确批准后合并**。
- 未经用户确认，不执行 merge 到 `main`。
