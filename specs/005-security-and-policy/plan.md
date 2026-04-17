# Implementation Plan: 多集群 Kubernetes 安全与策略治理中心

**Branch**: `005-security-and-policy` | **Date**: 2026-04-14 | **Spec**: [/mnt/e/code/kbManage/specs/005-security-and-policy/spec.md](/mnt/e/code/kbManage/specs/005-security-and-policy/spec.md)
**Input**: Feature specification from `/specs/005-security-and-policy/spec.md`

**Note**: 本文件由 `/speckit.plan` 生成，覆盖 005 的技术上下文、研究结论、设计边界和实施前置条件。

## 当前执行状态（2026-04-15）

- 已完成 `/speckit.specify`、`/speckit.plan`、`/speckit.tasks` 和 `/speckit.implement`。
- 已补齐治理证据：`artifacts/005-security-and-policy/branch-check.txt`、`artifacts/005-security-and-policy/backup-manifest.txt`、`artifacts/005-security-and-policy/mirror-and-remote-check.txt`。
- 005 implement 全量任务已完成，当前进入 PR 准备阶段，仍需遵守“中文 PR + 用户明确同意后再合并”。

## Summary

在保持 001-004 既有主栈和模块化单体结构不变的前提下，新增一个对标 Rancher 安全与策略治理能力的 005 功能域。005 聚焦“策略中心 + 准入控制 + 工作负载安全基线 + 例外治理 + 审计闭环”，支持平台级/工作空间级/项目级策略分层管理、按范围分发、分阶段启用、灰度验证以及违规处置追踪。首期明确排除 CIS/STIG 合规扫描、平台身份源整合、应用发布与灾备恢复，避免与 004 发布域和未来潜在安全扫描域发生职责重叠。

## Technical Context

**Language/Version**: Go 1.25；TypeScript 5.x；React 19.2  
**Primary Dependencies**: Gin；client-go；GORM；go-redis；策略评估与准入执行抽象（平台内部策略模板 + 规则执行器适配层）；React；Vite 8；Ant Design 6.3.x；React Router；TanStack Query；Zustand  
**Dependency Source**: Go 使用 `GOPROXY=https://goproxy.cn,direct`；前端使用 `https://registry.npmmirror.com`；若引入策略样例包、准入测试镜像或安全规则数据，优先使用阿里云、DaoCloud 或已批准国内代理，禁止默认境外源直连  
**Storage**: MySQL 8.4（策略定义、策略分配、命中记录、例外申请、整改状态、审计索引）；Redis 8.x（策略分发进度、短时命中缓存、例外时效索引、幂等键）；运行时准入与对象状态来自 Kubernetes API  
**Database Backup Plan**: 在进入 `/speckit.implement` 前必须执行：`docker exec mysql8 sh -lc 'mysqldump -hlocalhost -P3306 -uadmin -p123456 --all-databases --single-transaction --quick --routines --events --triggers' > artifacts/005-security-and-policy/mysql-backup-<timestamp>.sql`；若容器环境凭据与规范不一致，需在 `artifacts/005-security-and-policy/backup-manifest.txt` 记录差异原因、实际命令与恢复抽样验证  
**Testing**: 后端 `go test ./...` + contract/integration；前端 Vitest + React Testing Library；策略命中、准入模式切换、例外时效与权限隔离覆盖契约与集成测试  
**Target Platform**: 桌面浏览器 Web 管理台；Linux 容器化后端；已接入并可访问 Kubernetes API 的多集群环境  
**Project Type**: 前后端分离的 Web application，模块化单体控制面  
**Git Workflow**: 当前分支 `005-security-and-policy`；远程 `git@github.com:baihua19941101/kbManage.git`；实现后必须推送分支并提交中文 PR；未获用户明确同意不得合并；提交信息需专业、可审计  
**Performance Goals**: 90% 的策略命中查询在 30 秒内返回；90% 的最近 90 天治理记录检索在 30 秒内返回；90% 的新策略可先灰度验证再切换目标执行模式  
**Constraints**: 保持 001-004 既有能力边界；策略模式切换必须支持分阶段启用；权限回收必须即时生效；首期不包含 CIS/STIG 扫描、身份源整合、应用发布和灾备恢复；若使用子代理，模型固定 `gpt-5.3-codex`  
**Scale/Scope**: 首期覆盖至少 20 个已接入集群、平台级/工作空间级/项目级策略分层、策略分配到集群/命名空间/项目/资源类型、违规处置与例外全生命周期

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- `PASS`: 当前规划工作位于专用分支 `005-security-and-policy`，明确禁止在 `main` 或 `master` 开发。
- `PASS`: 已给出实施前数据库备份命令、产物路径与恢复校验要求；进入 implement 前必须先落地并留存 `artifacts/005-security-and-policy/backup-manifest.txt`。
- `PASS`: Go/npm 与后续策略相关依赖的国内镜像与代理要求已明确。
- `PASS`: 已定义 GitHub 远程、中文 PR、专业提交说明和“用户明确同意后再合并”的审批门槛。
- `PASS`: 若后续使用子代理，模型固定为 `gpt-5.3-codex`。
- `PASS`: 005 规格未遗留 `NEEDS CLARIFICATION` 标记，可进入设计阶段。

### Post-Design Re-check

- `PASS`: `research.md` 已收敛策略层级模型、执行模式语义、例外治理边界与审计闭环方案。
- `PASS`: `data-model.md` 已覆盖策略定义、策略分配、命中记录、例外申请、整改动作与审计事件。
- `PASS`: `contracts/openapi.yaml` 已覆盖策略管理、分配、命中查询、例外审批、整改更新与审计检索接口。
- `PASS`: `quickstart.md` 已给出实施前备份、国内源配置、最小联调路径与完整验收清单。

## Project Structure

### Documentation (this feature)

```text
specs/005-security-and-policy/
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
│   │   ├── delivery/
│   │   └── policy/
│   ├── kube/
│   │   ├── adapter/
│   │   └── client/
│   ├── repository/
│   ├── service/
│   │   ├── auth/
│   │   ├── audit/
│   │   ├── observability/
│   │   ├── workloadops/
│   │   ├── gitops/
│   │   └── securitypolicy/
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
│   │   ├── observability/
│   │   ├── workload-ops/
│   │   ├── gitops/
│   │   └── security-policy/
│   ├── services/
│   │   ├── api/
│   │   └── securityPolicy.ts
│   └── components/
└── tests/
```

**Structure Decision**: 沿用现有前后端分离模块化单体结构。后端新增 `securitypolicy` 业务域与 `integration/policy` 适配层，承载策略定义、分发、命中与例外治理；前端新增 `security-policy` 功能域，承接策略中心、命中视图、例外审批、整改追踪与审计查询入口，避免与 004 的发布流和 003 的运行时运维流混杂。

## Complexity Tracking

当前设计不存在需要宪章豁免的复杂度例外。
