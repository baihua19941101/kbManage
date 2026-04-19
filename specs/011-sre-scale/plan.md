# Implementation Plan: 平台 SRE 与规模化治理

**Branch**: `011-sre-scale` | **Date**: 2026-04-19 | **Spec**: [spec.md](/mnt/e/code/kbManage/specs/011-sre-scale/spec.md)
**Input**: Feature specification from `/specs/011-sre-scale/spec.md`

**Note**: 本文件由 `/speckit.plan` 生成，覆盖 011 的技术上下文、研究结论、设计边界和实施前置条件。

## Summary

011 将在现有多集群 Kubernetes 平台中新增“平台自身 SRE 与规模化治理中心”，聚焦控制面高可用、故障切换、容量阈值、性能基线、限流保护、维护窗口、升级前检查、滚动升级、回退验证、运行手册、自诊断和容量预测。实现上延续当前 Go + Gin + GORM + Redis 的后端模式与 React + Vite 前端模式，复用 002 的健康与趋势视图能力、007 的生命周期检查与版本语义、008 的演练与恢复思路、009/010 的权限审计接线方式，在后端新增平台 SRE 领域模型与治理接口，在前端新增统一的 SRE 运维工作台。

## Technical Context

**Language/Version**: Go 1.25；TypeScript 5.x；React 19.2  
**Primary Dependencies**: Gin；GORM；go-redis；平台运行健康聚合抽象；升级前检查与版本兼容性评估抽象；容量预测与压测证据接入抽象；Ant Design 6.3.x；React Router；TanStack Query；Zustand；Apache ECharts  
**Dependency Source**: Go 依赖使用 `GOPROXY=https://goproxy.cn,direct`；npm 使用 `https://registry.npmmirror.com`；若需要引入压测工具客户端、时间序列 SDK 或诊断组件，优先选择国内镜像或已批准代理  
**Storage**: MySQL 8.4（高可用策略、维护窗口、平台组件健康快照索引、容量基线、升级计划、回退验证、运行手册、告警基线、自诊断摘要、审计索引）；Redis 8.x（健康聚合缓存、任务积压短时状态、升级协调、限流状态、幂等键、容量预测缓存）  
**Database Backup Plan**: 实现前通过 `docker exec mysql8 sh -lc 'mysqldump -hlocalhost -P3306 -uadmin -p123456 --all-databases --single-transaction --quick --routines --events --triggers' > artifacts/011-sre-scale/mysql-backup-<timestamp>.sql` 执行备份；在 `artifacts/011-sre-scale/backup-manifest.txt` 记录命令、时间戳、产物路径；并使用临时 MySQL 容器执行恢复抽样验证  
**Testing**: `go test -run TestNonExistent -count=0 ./...`；011 定向 contract/integration 测试；前端 `npm run lint`、`npm run build`、定向 `vitest --maxWorkers=1`  
**Target Platform**: Linux 容器化后端 + 现代浏览器前端的多集群 Kubernetes 管理平台  
**Project Type**: Web application（`backend/` + `frontend/`）  
**Git Workflow**: 在 `011-sre-scale` 功能分支开发；推送到 `git@github.com:baihua19941101/kbManage.git`；使用中文提交与中文 PR 摘要；待用户明确同意后才允许合并 `main`  
**Performance Goals**: 平台健康总览、升级检查结果、容量趋势和自诊断摘要在试点规模下满足规格中的 30 秒内可得结果目标；关键异常分类在 5 分钟内归类；主要容量风险可提前至少 7 天识别  
**Constraints**: 必须延续现有权限与审计模型；必须先记录数据库备份证据再实施；不得引入新的业务域能力；首期不包含身份治理、应用市场和合规扫描；如使用子代理，规划与实现记录需遵循仓库当前约定  
**Scale/Scope**: 首期覆盖平台控制面多实例高可用、组件与依赖健康治理、升级闭环、容量基线与预测、压测证据、自诊断与运行手册治理，支撑多集群和大规模资源场景下的平台自身运维

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
说明：011 已在功能分支 `011-sre-scale` 上规划；数据库备份方案、国内镜像源、GitHub PR 路径和中文提交要求已明确；当前计划阶段未使用子代理执行实现。

### Post-Design Re-check

- `PASS`: `research.md` 已收敛高可用状态语义、升级治理模式、容量预测可信度、压测证据边界和运行手册关联原则。
- `PASS`: `data-model.md` 已覆盖高可用策略、维护窗口、健康快照、容量基线、升级计划、回退验证、运行手册和自诊断摘要等核心实体。
- `PASS`: `contracts/openapi.yaml` 已覆盖高可用、健康总览、升级治理、容量趋势、运行手册和审计查询接口。
- `PASS`: `quickstart.md` 已写明实施前备份、国内依赖源、最小联调路径与完整验收清单。

## Project Structure

### Documentation (this feature)

```text
specs/011-sre-scale/
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
│   │   └── sre/
│   ├── repository/
│   └── service/
│       └── sre/
├── migrations/
└── tests/
    ├── contract/
    └── integration/

frontend/
├── src/
│   ├── app/
│   ├── features/
│   │   ├── audit/
│   │   └── sre-scale/
│   └── services/
└── tests/
```

**Structure Decision**: 采用现有 Web 应用双端结构，在后端新增 `sre` 服务域与平台健康/升级/容量证据适配层，在前端新增 `sre-scale` 功能域，并复用全局路由、菜单、权限和审计页面接线。

## Complexity Tracking

当前设计不存在需要特别豁免的宪章违规项。
