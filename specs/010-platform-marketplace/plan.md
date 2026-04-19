# Implementation Plan: 平台应用目录与扩展市场

**Branch**: `010-platform-marketplace` | **Date**: 2026-04-19 | **Spec**: [spec.md](/mnt/e/code/kbManage/specs/010-platform-marketplace/spec.md)
**Input**: Feature specification from `/specs/010-platform-marketplace/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/plan-template.md` for the execution workflow.

## Summary

010 将在现有多集群平台中新增“应用目录、模板中心与扩展市场”能力，覆盖目录来源管理、模板版本与依赖治理、按工作空间/项目/集群范围分发模板，以及扩展包与插件的注册、启停、兼容性和权限声明治理。实现上沿用现有 Go + Gin + GORM + Redis 后端模式与 React + Vite 前端模式，复用 004 的发布治理语义、009 的权限边界与审计能力，新增目录来源、模板版本、模板分发、安装记录、扩展包和兼容性结论等模型及其接口。

## Technical Context

<!--
  ACTION REQUIRED: Replace the content in this section with the technical details
  for the project. The structure here is presented in advisory capacity to guide
  the iteration process.
-->

**Language/Version**: Go 1.25；TypeScript 5.x；React 19.2  
**Primary Dependencies**: Gin；GORM；go-redis；应用目录来源抽象（Git/Helm/OCI 风格目录访问）；扩展注册与兼容性评估抽象；Ant Design 6.3.x；React Router；TanStack Query；Zustand  
**Dependency Source**: Go 依赖使用 `GOPROXY=https://goproxy.cn,direct`；npm 使用 `https://registry.npmmirror.com`  
**Storage**: MySQL 8.4（目录来源、模板、版本、分发、安装记录、扩展包、兼容性结论、审计索引）；Redis 8.x（目录同步游标、短时模板缓存、分发协调、幂等键、兼容性结果缓存）  
**Database Backup Plan**: 实现前通过 `docker exec mysql8 sh -lc 'mysqldump -hlocalhost -P3306 -uadmin -p123456 --all-databases --single-transaction --quick --routines --events --triggers' > artifacts/010-platform-marketplace/mysql-backup-<timestamp>.sql` 执行备份；在 `artifacts/010-platform-marketplace/backup-manifest.txt` 记录命令、时间戳、产物路径；并用临时 MySQL 容器执行恢复抽样验证  
**Testing**: `go test -run TestNonExistent -count=0 ./...`；009/010 定向 contract/integration 测试；前端 `npm run lint`、`npm run build`、定向 `vitest --maxWorkers=1`  
**Target Platform**: Linux 容器化后端 + 现代浏览器前端的多集群 Kubernetes 管理平台  
**Project Type**: Web application（`backend/` + `frontend/`）  
**Git Workflow**: 在 `010-platform-marketplace` 功能分支开发；推送到 `git@github.com:baihua19941101/kbManage.git`；使用中文提交与中文 PR 摘要；待用户明确同意后才允许合并 `main`  
**Performance Goals**: 目录来源同步状态、模板查询、安装记录查询和扩展兼容性查询满足规格中的 30 秒内可得结果目标；目录中心首屏列表在常规试点规模下保持可交互  
**Constraints**: 必须复用现有权限、审计与发布模型；必须阻止跨租户超范围模板分发；必须阻止不兼容扩展启用；不得把完整 GitOps 编排、集群生命周期、统一身份源和合规扫描纳入首期  
**Scale/Scope**: 首期覆盖 1-3 个目录来源、每个来源数十个模板、每个模板多版本、多个工作空间/项目/集群范围分发，以及平台扩展包/插件的注册与启停治理

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- Dedicated feature branch is defined; implementation on `master` or `main` is forbidden.
- Database backup is planned before implementation begins, including the command run inside
  container `mysql8`, the artifact path, and any restore validation steps, or a justified N/A.
- Every dependency installation step specifies the China mirror or proxy configuration that will be used.
- Push and PR workflow to the GitHub remote is defined, including the Chinese PR summary and the
  explicit user-approval gate required before merge.
- Commit message expectations are documented and any delegated agent or subagent work is pinned
  to `gpt-5.3-codex`.

**Gate Result (Pre-Design)**: PASS  
说明：010 已在功能分支 `010-platform-marketplace` 上规划；数据库备份方案、国内镜像源、GitHub PR 路径和中文提交要求已明确；本阶段未使用子代理执行实现。

## Project Structure

### Documentation (this feature)

```text
specs/010-platform-marketplace/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── openapi.yaml
└── tasks.md
```

### Source Code (repository root)
<!--
  ACTION REQUIRED: Replace the placeholder tree below with the concrete layout
  for this feature. Delete unused options and expand the chosen structure with
  real paths (e.g., apps/admin, packages/something). The delivered plan must
  not include Option labels.
-->

```text
backend/
├── config/
├── internal/
│   ├── api/
│   │   ├── handler/
│   │   ├── middleware/
│   │   └── router/
│   ├── domain/
│   ├── integration/
│   │   └── marketplace/
│   ├── repository/
│   └── service/
│       └── marketplace/
├── migrations/
└── tests/
    ├── contract/
    └── integration/

frontend/
├── src/
│   ├── app/
│   ├── features/
│   │   ├── audit/
│   │   └── platform-marketplace/
│   └── services/
└── tests/
```

**Structure Decision**: 采用现有 Web 应用双端结构，在后端新增 `marketplace` 服务域和目录/扩展来源适配层，在前端新增 `platform-marketplace` 功能域，并复用全局路由、菜单、权限和审计页面接线。

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

无。010 的设计未违反宪章门槛。
