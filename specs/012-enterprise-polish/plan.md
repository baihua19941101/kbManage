# Implementation Plan: 企业级治理报表与产品化交付收尾

**Branch**: `012-enterprise-polish` | **Date**: 2026-04-19 | **Spec**: [spec.md](/mnt/e/code/kbManage/specs/012-enterprise-polish/spec.md)
**Input**: Feature specification from `/specs/012-enterprise-polish/spec.md`

**Note**: 本文件由 `/speckit.plan` 生成，覆盖 012 的技术上下文、研究结论、设计边界和实施前置条件。

## Summary

012 将在现有多集群 Kubernetes 平台中新增“企业级治理报表与产品化交付收尾中心”，聚焦深度权限审计、关键操作追踪、跨团队授权分布、高风险访问分析、治理覆盖率与长期趋势报表，以及面向管理汇报、审计复核和客户交付的标准化报告、导出材料和交付清单。实现上延续当前 Go + Gin + GORM + Redis 的后端模式与 React + Vite 前端模式，复用 009 的身份与授权关系、011 的审计与趋势治理语义、008 的报告与演练留痕方式、010 的标准化目录和分发经验，在后端新增治理报表与交付包领域模型，在前端新增企业交付收尾工作台。

## Technical Context

**Language/Version**: Go 1.25；TypeScript 5.x；React 19.2  
**Primary Dependencies**: Gin；GORM；go-redis；现有权限审计聚合抽象；治理覆盖率快照与报表生成抽象；导出记录与交付包编排抽象；Ant Design 6.3.x；React Router；TanStack Query；Zustand；Apache ECharts  
**Dependency Source**: Go 依赖使用 `GOPROXY=https://goproxy.cn,direct`；npm 使用 `https://registry.npmmirror.com`；若需补充报表导出、文档模板或图表依赖，优先使用国内镜像或已批准代理  
**Storage**: MySQL 8.4（权限变更链路、关键操作追踪索引、跨团队授权分布快照、治理覆盖率快照、治理报表、导出记录、交付材料目录、交付检查清单、审计索引）；Redis 8.x（报表生成上下文、导出任务短时状态、权限风险聚合缓存、幂等键、趋势查询缓存）  
**Database Backup Plan**: 实现前通过 `docker exec mysql8 sh -lc 'mysqldump -hlocalhost -P3306 -uadmin -p123456 --all-databases --single-transaction --quick --routines --events --triggers' > artifacts/012-enterprise-polish/mysql-backup-<timestamp>.sql` 执行备份；在 `artifacts/012-enterprise-polish/backup-manifest.txt` 记录命令、时间戳、产物路径；并使用临时 MySQL 容器执行恢复抽样验证  
**Testing**: `go test -run TestNonExistent -count=0 ./...`；012 定向 contract/integration 测试；前端 `npm run lint`、`npm run build`、定向 `vitest --maxWorkers=1`  
**Target Platform**: Linux 容器化后端 + 现代浏览器前端的多集群 Kubernetes 管理平台  
**Project Type**: Web application（`backend/` + `frontend/`）  
**Git Workflow**: 在 `012-enterprise-polish` 功能分支开发；推送到 `git@github.com:baihua19941101/kbManage.git`；使用中文提交与中文 PR 摘要；待用户明确同意后才允许合并 `main`  
**Performance Goals**: 深度审计查询和治理报表汇总在试点规模下满足规格中的 10-15 分钟人工复核窗口要求；标准化交付包在 30 分钟内可生成；100% 的报表生成与导出动作具备可检索审计轨迹  
**Constraints**: 必须延续现有权限与审计模型；必须先记录数据库备份证据再实施；不得引入新的核心业务域能力；首期不包含集群创建导入、GitOps 编排和灾备扩展；计划阶段不使用子代理实现工作  
**Scale/Scope**: 首期覆盖深度权限审计、关键操作追踪、跨团队授权分析、治理覆盖率快照、管理/审计/交付三类报表、导出留痕、交付材料目录和交付检查清单，支撑多团队、多环境和多客户交付场景

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
说明：012 已在功能分支 `012-enterprise-polish` 上规划；数据库备份方案、国内镜像源、GitHub PR 路径和中文提交要求已明确；当前计划阶段未使用子代理执行实现，因此不存在模型约束冲突。

### Post-Design Re-check

- `PASS`: `research.md` 已收敛深度审计边界、治理报表口径、导出分层策略和产品化交付包建模原则。
- `PASS`: `data-model.md` 已覆盖权限变更链路、治理风险事件、治理覆盖率快照、报表包、交付材料和检查清单等核心实体。
- `PASS`: `contracts/openapi.yaml` 已覆盖深度审计查询、治理报表生成、导出记录、交付包目录和就绪检查接口。
- `PASS`: `quickstart.md` 已写明实施前备份、国内依赖源、最小联调路径与交付验收清单。

## Project Structure

### Documentation (this feature)

```text
specs/012-enterprise-polish/
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
├── config/
├── internal/
│   ├── api/
│   │   ├── handler/
│   │   ├── middleware/
│   │   └── router/
│   ├── domain/
│   ├── integration/
│   │   └── enterprise/
│   ├── repository/
│   └── service/
│       └── enterprise/
├── migrations/
└── tests/
    ├── contract/
    └── integration/

frontend/
├── src/
│   ├── app/
│   ├── features/
│   │   ├── audit/
│   │   └── enterprise-polish/
│   └── services/
└── tests/
```

**Structure Decision**: 采用现有 Web 应用双端结构，在后端新增 `enterprise` 服务域与治理报表/交付包适配层，在前端新增 `enterprise-polish` 功能域，并复用全局路由、菜单、权限和审计页面接线。

## Complexity Tracking

当前设计不存在需要特别豁免的宪章违规项。
