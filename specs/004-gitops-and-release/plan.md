# Implementation Plan: 多集群 GitOps 与应用发布中心

**Branch**: `004-gitops-and-release` | **Date**: 2026-04-13 | **Spec**: [/mnt/e/code/kbManage/specs/004-gitops-and-release/spec.md](/mnt/e/code/kbManage/specs/004-gitops-and-release/spec.md)
**Input**: Feature specification from `/specs/004-gitops-and-release/spec.md`

**Note**: 本文件由 `/speckit.plan` 生成，覆盖 004 的技术上下文、研究结论、设计边界和实施前置条件。

## 当前执行状态（2026-04-13）

- 已完成 `/speckit.specify`，004 规格已创建并通过质量清单校验。
- 当前工作分支为 `004-gitops-and-release`。
- 已完成 004 实施前数据库备份与恢复抽样验证，详见 `artifacts/004-gitops-and-release/backup-manifest.txt`。
- 已完成 `/speckit.plan`，本轮已生成 `research.md`、`data-model.md`、`contracts/openapi.yaml` 和 `quickstart.md`。
- 已完成 `/speckit.tasks`，`tasks.md` 已生成。
- Phase 0 治理任务已完成：T001、T002、T003、T004。
- Phase 1 配置任务已启动并完成 T008：已补齐 gitops.sources/sync/diff/release/audit 配置与文档说明。
- 当前状态为“tasks 阶段治理门槛已达成，implement 前置配置已启动并完成 T008”，可进入 `/speckit.implement` 阶段；后续实现仍需遵守“中文 PR + 用户明确同意后再合并”的宪章要求。

## Summary

在保持 001、002、003 既有主栈与模块化单体结构不变的前提下，为 kbManage 增加一个对标 Rancher GitOps / Fleet 的多集群 GitOps 与应用发布中心。该能力围绕交付来源接入、应用交付单元、集群组目标、多环境分层、配置覆盖、同步与漂移状态、发布历史、回滚以及暂停/恢复等生命周期动作展开。平台继续作为受控控制面，不承担通用 CI 流水线编排、制品仓库管理、终端运维、策略准入或合规扫描职责，而是通过现有工作空间/项目范围模型、统一发布动作链路和审计模型，形成“来源 -> 目标 -> 环境 -> 同步/发布 -> 回滚 -> 审计”的持续交付闭环。

## Technical Context

**Language/Version**: Go 1.25；TypeScript 5.x；React 19.2  
**Primary Dependencies**: Gin；client-go；GORM；go-redis；go-git 风格 Git 访问抽象；Helm SDK 风格发布源抽象；React；Vite 8；Ant Design 6.3.x；React Router；TanStack Query；Zustand  
**Dependency Source**: Go 使用 `GOPROXY=https://goproxy.cn,direct`；前端使用 `https://registry.npmmirror.com`；如需镜像、Helm/OCI 相关依赖或本地联调组件，优先使用阿里云、DaoCloud 或已批准的国内代理，禁止直接裸连默认境外源  
**Storage**: MySQL 8.4（交付来源、目标组、环境阶段、配置覆盖、发布历史、同步/发布动作、审计索引）；Redis 8.x（动作进度、分布式锁、幂等键、短时差异缓存、阶段推进上下文）；Git/Helm/OCI 来源内容与 Kubernetes 实际状态由外部来源和 Kubernetes API 提供  
**Database Backup Plan**: 已执行 `docker exec mysql8 sh -lc 'mysqldump -hlocalhost -P3306 -uroot -p123456 --all-databases --single-transaction --quick --routines --events --triggers' > artifacts/004-gitops-and-release/mysql-backup-20260413-155243.sql`；默认凭据 `admin/123456` 在当前容器环境不可用，本轮按实际可用凭据 `root/123456` 执行，并通过临时容器 `mysql8-restore-check-004` 完成完整导入和数据库列表校验，详见 `artifacts/004-gitops-and-release/backup-manifest.txt`  
**Testing**: 后端使用 `go test ./...` + Testify + `httptest` + fake `client-go`；前端使用 Vitest + React Testing Library；契约测试覆盖来源管理、交付单元、状态/差异、发布动作、阶段推进和审计接口；集成测试覆盖权限隔离、部分成功、暂停恢复、回滚和漂移处理  
**Target Platform**: 桌面浏览器中的 Web 管理台；Linux 容器化后端服务；已接入并可访问 Kubernetes API、Git 仓库和应用发布来源的多集群环境  
**Project Type**: 前后端分离的 Web application，模块化单体控制面  
**Git Workflow**: 当前规划分支为 `004-gitops-and-release`，远程为 `git@github.com:baihua19941101/kbManage.git`；交付时必须推送该分支并提交中文 PR，合并前必须获得用户明确同意，所有提交说明必须专业且可审计  
**Performance Goals**: 90% 的应用交付单元状态查询在 30 秒内返回目标环境的期望状态、实际状态、同步结果和漂移标记；95% 的最多覆盖 50 个目标集群的单次同步或发布动作在 5 分钟内返回明确结果；90% 的最近 90 天交付记录查询在 30 秒内返回结果集  
**Constraints**: 保持 001/002/003 主栈与模块化单体结构不变；所有来源接入、差异判断、同步和发布动作必须经后端统一授权；首期只聚焦 Git 仓库与应用发布来源，不做通用制品仓库治理；首期多环境推进采用按环境顺序的受控推进，不引入通用 CI 工作流引擎、策略准入或合规扫描；如使用子代理必须固定为 `gpt-5.3-codex`  
**Scale/Scope**: 首期覆盖至少 20 个已接入集群、多个工作空间/项目/环境阶段、最多 50 个目标集群的单次同步或发布动作；覆盖来源接入、目标分组、环境推进、差异/漂移、发布历史、回滚、暂停/恢复和审计闭环

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- `PASS`: 当前规划工作位于专用分支 `004-gitops-and-release`，明确禁止在 `main` 或 `master` 上实施。
- `PASS`: 已记录 004 的数据库备份命令、产物路径、凭据差异和恢复抽样验证步骤；2026-04-13 已按容器实际可用凭据 `root/123456` 完成一次备份与恢复抽样验证。
- `PASS`: Go、npm 以及后续 Git/Helm/OCI 联调依赖的国内镜像或代理来源已明确。
- `PASS`: 已记录 GitHub SSH 远程、中文 PR 约束、专业提交说明和“必须获得用户同意后才能合并”的审批门槛。
- `PASS`: 若后续使用子代理或并行代理，模型固定为 `gpt-5.3-codex`。
- `PASS`: 004 规划前置输入已经齐备，不存在未解决的规格澄清项。

### Post-Design Re-check

- `PASS`: `research.md` 已收敛 004 的来源模型、目标分组、环境推进、状态/漂移语义、回滚历史和权限边界。
- `PASS`: `data-model.md` 已覆盖交付来源、目标组、环境阶段、交付单元、配置覆盖、发布动作和审计模型。
- `PASS`: `contracts/openapi.yaml` 已覆盖来源管理、目标组、交付单元、状态/差异、动作执行、发布历史和审计接口。
- `PASS`: `quickstart.md` 已写明数据库备份、国内源设置、最小联调模式与完整验收模式。

## Project Structure

### Documentation (this feature)

```text
specs/004-gitops-and-release/
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
│   ├── integration/
│   │   ├── observability/
│   │   └── delivery/
│   │       ├── git/
│   │       ├── helm/
│   │       └── diff/
│   ├── kube/
│   │   ├── adapter/
│   │   └── client/
│   ├── repository/
│   ├── service/
│   │   ├── auth/
│   │   ├── audit/
│   │   ├── cluster/
│   │   ├── observability/
│   │   ├── workloadops/
│   │   └── gitops/
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
│   │   ├── observability/
│   │   ├── workload-ops/
│   │   └── gitops/
│   ├── services/
│   │   ├── api/
│   │   └── gitops.ts
│   └── components/
└── tests/
```

**Structure Decision**: 继续沿用 001/002/003 的前后端分离模块化单体结构。后端新增 `gitops` 业务域与 `integration/delivery` 适配层，集中承载来源管理、目标分组、环境推进、差异比较、同步/发布动作和回滚编排；前端新增独立 `gitops` 功能域，承接交付中心概览、来源管理、交付单元详情、差异视图、发布历史和环境推进入口，避免把 004 能力混入 003 的 `workload-ops` 或 001 的通用 `operations`。

## Complexity Tracking

当前设计不存在需要特别豁免的宪章违规项。
