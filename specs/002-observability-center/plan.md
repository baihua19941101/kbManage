# Implementation Plan: 多集群 Kubernetes 可观测中心

**Branch**: `002-observability-center` | **Date**: 2026-04-11 | **Spec**: [/mnt/e/code/kbManage/specs/002-observability-center/spec.md](/mnt/e/code/kbManage/specs/002-observability-center/spec.md)
**Input**: Feature specification from `/specs/002-observability-center/spec.md`

**Note**: 本文件由 `/speckit.plan` 生成，覆盖 002 的技术上下文、研究结论、设计边界和实施前置条件。

## 当前执行状态（2026-04-11）

- `001` 已确认合并 `main`，002 的实现 gate 已解除。
- 已完成实施前数据库备份与恢复抽样验证（`artifacts/002-observability-center/backup-manifest.txt`）。
- 当前进入实现阶段，优先执行 Governance、Setup 与 Foundational 骨架任务，保持 001 主栈不变。

## Summary

在保持 001 既有主栈不变的前提下，为 kbManage 增加一个对标 Rancher 的多集群可观测中心。平台继续作为控制面，围绕资源上下文统一串联日志、事件、指标和告警，不自建日志或指标存储；原始观测数据来自 Prometheus、Alertmanager 和 Loki 兼容后端，事件直接来自 Kubernetes API。平台自身只保存数据源配置、告警规则治理元数据、通知目标、静默窗口、告警快照、处理记录和审计信息，通过后端适配层和现有工作空间/项目授权模型，形成受控的故障发现、定位和告警处理闭环。

## Technical Context

**Language/Version**: Go 1.25；TypeScript 5.x；React 19.2  
**Primary Dependencies**: Gin；client-go；GORM；go-redis；React；Vite 8；Ant Design 6.3.x；React Router；TanStack Query；Zustand；Apache ECharts；Prometheus-compatible Query API；Alertmanager-compatible API；Loki-compatible Query API  
**Dependency Source**: Go 使用 `GOPROXY=https://goproxy.cn,direct`；前端使用 `https://registry.npmmirror.com`；如需本地联调可观测后端或图表相关依赖，容器镜像优先使用阿里云、DaoCloud 或已批准的国内代理，禁止直接裸连默认境外源  
**Storage**: MySQL 8.4（可观测数据源配置、告警规则治理元数据、通知目标、静默窗口、告警快照、处理记录、审计索引）；Redis 8.x（查询缓存、同步游标、短时上下文、告警状态协同）；原始日志/指标/运行态告警由外部 Prometheus/Alertmanager/Loki 兼容后端保存  
**Database Backup Plan**: 实施前执行 `mkdir -p artifacts/002-observability-center && docker exec mysql8 sh -lc 'mysqldump -hlocalhost -P3306 -uadmin -p123456 --all-databases --single-transaction --quick --routines --events --triggers' > artifacts/002-observability-center/mysql-backup-$(date +%Y%m%d-%H%M%S).sql`；恢复校验使用临时 MySQL 容器导入备份并核对数据库与表计数；`artifacts/001-k8s-ops-platform/` 中已有历史备份证据，但 002 进入实现前必须重新执行并留档。2026-04-11 已按容器实际凭据 `root/123456` 完成一次备份，产物为 `artifacts/002-observability-center/mysql-backup-20260411-214819.sql`，详见 `artifacts/002-observability-center/backup-manifest.txt`  
**Testing**: 后端使用 `go test ./...` + Testify + `httptest` + fake `client-go`；前端使用 Vitest + React Testing Library；契约测试覆盖日志、事件、指标、告警、规则、静默和通知目标 API；集成测试覆盖多集群查询、权限隔离、后端适配器降级和告警处理闭环  
**Target Platform**: 桌面浏览器中的 Web 管理台；Linux 容器化后端服务；已接入且可访问 Prometheus/Alertmanager/Loki 兼容观测后端的 Kubernetes 集群  
**Project Type**: 前后端分离的 Web application，模块化单体控制面  
**Git Workflow**: 当前规划分支为 `002-observability-center`，远程为 `git@github.com:baihua19941101/kbManage.git`；`001` 已确认合并到 `main`，002 可进入实现阶段；交付时必须推送该分支并提交中文 PR，合并前必须取得用户明确同意，所有提交说明必须专业且可审计  
**Performance Goals**: 常见 24 小时日志查询 90% 在 10 秒内返回首批可阅读结果；30 天告警检索 90% 在 15 秒内返回；授权用户在异常演练中可在 3 分钟内从目标资源进入关联日志、事件、指标和告警视图并完成初步定位  
**Constraints**: 保持 001 现有主栈与模块化单体结构不变；平台只做观测控制面，不承担日志/指标原始存储职责；所有观测查询必须经后端统一授权与审计；首期不包含终端、批量操作、回滚、GitOps、Helm、策略治理、合规扫描、集群创建导入和灾备能力；如使用子代理必须固定为 `gpt-5.3-codex`  
**Scale/Scope**: 首期覆盖多集群统一概览、日志检索、事件时间线、指标趋势、告警中心、规则治理、通知目标和静默窗口；目标环境至少 20 个集群、多个工作空间/项目、24 小时常见日志检索和 30 天告警分析；实现范围以资源上下文关联和告警闭环为中心，不扩展到 Day2 运维动作

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- `PASS`: 当前规划工作位于专用分支 `002-observability-center`，明确禁止在 `main` 或 `master` 上实施。
- `PASS`: 已为 002 记录新的数据库备份命令、产物路径和恢复校验步骤，并已于 2026-04-11 完成一次实际备份与恢复抽样验证；虽然 admin/123456 在当前容器中不可用，但已按容器实际凭据 root/123456 留档执行。
- `PASS`: Go、npm 以及本地可观测后端联调所需的镜像/代理来源已明确为国内镜像或经批准代理。
- `PASS`: 已记录 GitHub SSH 远程、中文 PR 约束、专业提交说明和“必须获得用户同意后才能合并”的审批门槛。
- `PASS`: 若后续使用子代理或并行代理，模型固定为 `gpt-5.3-codex`。
- `PASS`: 用户已确认 001 合并完成，且 002 的分支/备份/远程/镜像证据已补充到 artifacts，当前可进入实现。

### Post-Design Re-check

- `PASS`: `research.md` 已收敛控制面/数据面分层、Prometheus/Alertmanager/Loki 兼容接入、告警治理元数据归属和查询编排策略。
- `PASS`: `data-model.md` 已定义 002 新增的主要治理实体，并明确哪些对象为平台持久化实体、哪些对象为请求级或派生级对象。
- `PASS`: `contracts/openapi.yaml` 已覆盖日志查询、事件时间线、指标视图、告警中心、规则治理、通知目标、静默窗口和集群观测配置接口。
- `PASS`: `quickstart.md` 已写明数据库备份、国内源设置、最小联调模式与完整验收模式；当前实施前数据库备份门槛已实际完成。

## Project Structure

### Documentation (this feature)

```text
specs/002-observability-center/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── openapi.yaml
└── tasks.md
```

### Source Code (repository root)

```text
backend/
├── cmd/
│   └── server/
├── config/
├── internal/
│   ├── api/
│   │   ├── handler/
│   │   ├── middleware/
│   │   └── router/
│   ├── domain/
│   ├── repository/
│   ├── service/
│   │   ├── auth/
│   │   ├── cluster/
│   │   ├── audit/
│   │   └── observability/
│   ├── integration/
│   │   └── observability/
│   │       ├── alerts/
│   │       ├── logs/
│   │       └── metrics/
│   ├── kube/
│   │   ├── adapter/
│   │   └── client/
│   └── worker/
├── migrations/
└── tests/
    ├── contract/
    ├── integration/
    └── testutil/

frontend/
├── src/
│   ├── app/
│   ├── features/
│   │   ├── auth/
│   │   ├── clusters/
│   │   ├── resources/
│   │   └── observability/
│   ├── services/
│   │   └── observability/
│   └── components/
└── tests/
```

**Structure Decision**: 继续沿用 001 的前后端分离模块化单体结构。后端新增 `observability` 业务域与 `integration/observability` 适配层，负责屏蔽 Prometheus/Alertmanager/Loki 兼容实现差异；前端新增独立 `observability` 功能域，用于统一总览、日志、事件、指标和告警中心页面，避免把 002 功能散落到 clusters/resources/operations 既有模块中。

## Complexity Tracking

当前设计不存在需要特别豁免的宪章违规项。
