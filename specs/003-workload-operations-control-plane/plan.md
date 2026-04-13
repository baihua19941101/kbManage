# Implementation Plan: 多集群 Kubernetes 工作负载运维控制面

**Branch**: `003-workload-operations-control-plane` | **Date**: 2026-04-12 | **Spec**: [/mnt/e/code/kbManage/specs/003-workload-operations-control-plane/spec.md](/mnt/e/code/kbManage/specs/003-workload-operations-control-plane/spec.md)
**Input**: Feature specification from `/specs/003-workload-operations-control-plane/spec.md`

**Note**: 本文件由 `/speckit.plan` 生成，覆盖 003 的技术上下文、研究结论、设计边界和实施前置条件。

## 当前执行状态（2026-04-12）

- 已完成 `/speckit.specify`，003 规格已创建并通过质量清单校验。
- 当前工作分支为 `003-workload-operations-control-plane`。
- 已完成 003 实施前数据库备份与恢复抽样验证，详见 `artifacts/003-workload-operations-control-plane/backup-manifest.txt`。
- 已完成 `/speckit.tasks`，任务清单见 `tasks.md`。
- 已启动 `/speckit.implement` 并完成全部阶段：Phase 0（Governance Gates）+ Phase 1（Setup）+ Phase 2（Foundational）+ Phase 3（US1）+ Phase 4（US2）+ Phase 5（US3）+ Final Phase（Polish & Delivery Readiness）。

## Summary

在保持 001 与 002 既有主栈和模块化单体结构不变的前提下，为 kbManage 增加一个对标 Rancher 的工作负载运维控制面。该能力围绕 `Deployment`、`StatefulSet`、`DaemonSet` 及其关联的 `Pod/Container` 展开，提供单资源诊断入口、实例日志与终端访问、扩缩容、重启、重新部署、实例替换、批量操作、发布历史查看和版本回滚能力。平台继续作为受控控制面，不承担 GitOps、Helm、统一监控告警或集群生命周期职责，而是通过现有工作空间/项目/命名空间授权模型、统一动作执行链路和审计模型，形成“诊断 -> 动作 -> 跟踪 -> 审计”的工作负载运维闭环。

## Technical Context

**Language/Version**: Go 1.25；TypeScript 5.x；React 19.2  
**Primary Dependencies**: Gin；client-go；GORM；go-redis；React；Vite 8；Ant Design 6.3.x；React Router；TanStack Query；Zustand  
**Dependency Source**: Go 使用 `GOPROXY=https://goproxy.cn,direct`；前端使用 `https://registry.npmmirror.com`；如需终端流式调试依赖或本地联调镜像，优先使用阿里云、DaoCloud 或已批准的国内代理，禁止直接裸连默认境外源  
**Storage**: MySQL 8.4（工作负载动作请求、批量任务、发布历史快照索引、终端会话审计、扩展运维审计索引）；Redis 8.x（动作进度缓存、批量执行协调、终端短会话上下文、幂等键和短时结果缓存）；运行时工作负载状态、Pod/容器状态、日志流和 exec 通道来自 Kubernetes API  
**Database Backup Plan**: 实施前执行 `mkdir -p artifacts/003-workload-operations-control-plane && docker exec mysql8 sh -lc 'mysqldump -hlocalhost -P3306 -uadmin -p123456 --all-databases --single-transaction --quick --routines --events --triggers' > artifacts/003-workload-operations-control-plane/mysql-backup-$(date +%Y%m%d-%H%M%S).sql`；若当前容器凭据与 `admin/123456` 不一致，必须按容器实际凭据执行并在 `backup-manifest.txt` 记录差异与恢复抽样验证结果。2026-04-12 已按容器实际凭据 `root/123456` 完成一次备份，产物为 `artifacts/003-workload-operations-control-plane/mysql-backup-20260412-150123.sql`  
**Testing**: 后端使用 `go test ./...` + Testify + `httptest` + fake `client-go`；前端使用 Vitest + React Testing Library；契约测试覆盖工作负载运维总览、动作执行、批量任务、发布历史、回滚、终端会话和审计接口；集成测试覆盖权限隔离、动作状态流转、部分失败、回滚恢复和终端会话关闭  
**Target Platform**: 桌面浏览器中的 Web 管理台；Linux 容器化后端服务；已接入并可访问 Kubernetes API 的多集群环境  
**Project Type**: 前后端分离的 Web application，模块化单体控制面  
**Git Workflow**: 当前规划分支为 `003-workload-operations-control-plane`，远程为 `git@github.com:baihua19941101/kbManage.git`；交付时必须推送该分支并提交中文 PR，合并前必须获得用户明确同意，所有提交说明必须专业且可审计  
**Performance Goals**: 90% 的授权用户可在 3 分钟内从工作负载进入实例诊断视图；95% 的单资源常见运维动作在 2 分钟内返回明确结果；90% 的最多 50 个对象批量动作在 5 分钟内完成结果归集；90% 的可回滚工作负载回滚动作在 3 分钟内返回明确结果  
**Constraints**: 保持 001/002 主栈与模块化单体结构不变；所有日志、终端和运维动作必须经后端统一授权；终端首期只记录会话元数据审计，不记录完整命令与终端输出正文；首期不包含 GitOps、Helm、统一监控告警、策略治理和集群生命周期管理；如使用子代理必须固定为 `gpt-5.3-codex`  
**Scale/Scope**: 首期覆盖 `Deployment`、`StatefulSet`、`DaemonSet` 工作负载及关联 `Pod/Container` 的单资源诊断、单体动作、批量动作、发布历史、回滚和终端访问；目标环境至少 20 个集群、多个工作空间/项目、最多 50 个目标对象的同类批量动作

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- `PASS`: 当前规划工作位于专用分支 `003-workload-operations-control-plane`，明确禁止在 `main` 或 `master` 上实施。
- `PASS`: 已记录 003 的数据库备份命令、产物路径、凭据差异和恢复抽样验证步骤；2026-04-12 已按容器实际凭据 `root/123456` 完成一次备份与恢复抽样验证。
- `PASS`: Go、npm 和本地联调依赖的国内镜像或代理来源已明确。
- `PASS`: 已记录 GitHub SSH 远程、中文 PR 约束、专业提交说明和“必须获得用户同意后才能合并”的审批门槛。
- `PASS`: 若后续使用子代理或并行代理，模型固定为 `gpt-5.3-codex`。
- `PASS`: 003 规划前置输入已经齐备，不存在未解决的规格澄清项。

### Post-Design Re-check

- `PASS`: `research.md` 已收敛终端会话模型、发布历史来源、批量动作执行策略和回滚边界。
- `PASS`: `data-model.md` 已覆盖工作负载运维视图、动作请求、批量任务、终端会话和发布历史模型。
- `PASS`: `contracts/openapi.yaml` 已覆盖工作负载运维读取、动作执行、批量任务、回滚和终端会话接口。
- `PASS`: `quickstart.md` 已写明数据库备份、国内源设置、最小联调模式与完整验收模式。

## Project Structure

### Documentation (this feature)

```text
specs/003-workload-operations-control-plane/
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
│   │   ├── operation/
│   │   ├── observability/
│   │   └── workloadops/
│   ├── kube/
│   │   ├── adapter/
│   │   ├── client/
│   │   └── exec/
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
│   │   ├── resources/
│   │   ├── operations/
│   │   ├── observability/
│   │   └── workload-ops/
│   ├── services/
│   │   ├── observability/
│   │   └── workloadOps.ts
│   └── components/
└── tests/
```

**Structure Decision**: 继续沿用 001/002 的前后端分离模块化单体结构。后端新增 `workloadops` 业务域，集中承载工作负载诊断聚合、动作编排、批量执行、发布历史、回滚和终端会话管理；前端新增独立 `workload-ops` 功能域，承接单资源运维页、批量动作页、发布历史和终端会话入口，避免继续把 Day2 运维能力堆叠进 001 的通用 `operations` 页或 002 的 `observability` 页中。

## Complexity Tracking

当前设计不存在需要特别豁免的宪章违规项。
