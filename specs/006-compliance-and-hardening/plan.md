# Implementation Plan: 多集群 Kubernetes 合规与加固中心

**Branch**: `006-compliance-and-hardening` | **Date**: 2026-04-15 | **Spec**: [/mnt/e/code/kbManage/specs/006-compliance-and-hardening/spec.md](/mnt/e/code/kbManage/specs/006-compliance-and-hardening/spec.md)
**Input**: Feature specification from `/specs/006-compliance-and-hardening/spec.md`

**Note**: 本文件由 `/speckit.plan` 生成，覆盖 006 的技术上下文、研究结论、设计边界和实施前置条件。

## 当前执行状态（2026-04-15）

- 已完成 `/speckit.specify`、`/speckit.plan` 与 `/speckit.tasks`，本轮已生成 `research.md`、`data-model.md`、`contracts/openapi.yaml`、`quickstart.md` 和 `tasks.md`。
- 当前工作分支为 `006-compliance-and-hardening`。
- 已确认 `005-security-and-policy` 的主 PR 流已完成并合并到 `main`，006 的规划 gate 已解除。
- 006 当前进入“tasks ready”阶段，后续进入实现前仍需落实数据库备份证据、中文 PR 约束和用户明确同意后再合并。

## Summary

在保持 001-005 既有主栈和模块化单体结构不变的前提下，为 kbManage 新增一个对标企业级 Kubernetes 合规与加固能力的 006 功能域。006 聚焦“基线标准管理 + 计划性/按需扫描 + 失败项证据 + 整改/例外/复检闭环 + 趋势复盘与审计归档”，支持围绕 `CIS`、`STIG` 和平台基线模板对多集群环境进行统一评估。平台继续作为受控控制面，不承担准入策略执行、自动原位修复、统一身份治理、应用发布或集群创建导入职责，而是通过统一扫描编排、证据快照、治理对象流转和审计模型，形成“发现 -> 分析 -> 整改/例外 -> 复检 -> 归档”的合规治理闭环。

## Technical Context

**Language/Version**: Go 1.25；TypeScript 5.x；React 19.2  
**Primary Dependencies**: Gin；client-go；GORM；go-redis；合规扫描执行抽象（平台编排 + 外部扫描器/基线包适配层）；React；Vite 8；Ant Design 6.3.x；React Router；TanStack Query；Zustand；Apache ECharts  
**Dependency Source**: Go 使用 `GOPROXY=https://goproxy.cn,direct`；前端使用 `https://registry.npmmirror.com`；若引入扫描器镜像、基线规则包或趋势图表相关依赖，优先使用阿里云、DaoCloud 或经批准的国内代理，禁止默认境外源直连  
**Storage**: MySQL 8.4（基线标准、扫描配置、扫描执行、失败项、证据索引、整改任务、例外审批、复检任务、趋势快照、审计索引）；Redis 8.x（扫描调度队列、进度缓存、短时证据缓存、幂等键、复检协调）；运行时证据与原始检查结果来自 Kubernetes API 与外部扫描器/基线执行适配层  
**Database Backup Plan**: 进入实现前执行 `mkdir -p artifacts/006-compliance-and-hardening && docker exec mysql8 sh -lc 'mysqldump -hlocalhost -P3306 -uadmin -p123456 --all-databases --single-transaction --quick --routines --events --triggers' > artifacts/006-compliance-and-hardening/mysql-backup-$(date +%Y%m%d-%H%M%S).sql`；若当前容器凭据与 `admin/123456` 不一致，必须按容器实际凭据执行并在 `artifacts/006-compliance-and-hardening/backup-manifest.txt` 记录差异、实际命令与恢复抽样验证结果；恢复校验使用临时 MySQL 容器导入备份并核对数据库与核心表可见性  
**Testing**: 后端使用 `go test ./...` + Testify + `httptest` + fake `client-go` + 扫描适配层 stub；前端使用 Vitest + React Testing Library；契约测试覆盖基线、扫描配置、扫描执行、失败项、整改、例外、复检、趋势与审计接口；集成测试覆盖计划扫描、按需扫描、部分成功、权限隔离、整改闭环和趋势聚合  
**Target Platform**: 桌面浏览器 Web 管理台；Linux 容器化后端；已接入并可访问 Kubernetes API、节点/资源元数据与外部扫描器适配层的多集群环境  
**Project Type**: 前后端分离的 Web application，模块化单体控制面  
**Git Workflow**: 当前规划分支为 `006-compliance-and-hardening`，远程为 `git@github.com:baihua19941101/kbManage.git`；`005` 已确认合并到 `main`；006 后续交付必须推送该分支并提交中文 PR，合并前必须取得用户明确同意，所有提交说明必须专业且可审计  
**Performance Goals**: 90% 的常见合规扫描结果在 30 分钟内向用户返回明确摘要；90% 的最近 90 天趋势查询在 30 秒内返回；80% 的高风险失败项在 24 小时内进入整改、例外或复检安排中的一种受控处置状态  
**Constraints**: 保持 001-005 既有能力边界；006 仅负责评估、报告、整改复核与归档，不承担实时准入阻断；所有扫描与证据访问必须经后端统一授权与审计；首期不包含自动原位修复、批量强制加固、统一身份治理、应用发布和集群创建导入；若使用子代理必须固定为 `gpt-5.3-codex`  
**Scale/Scope**: 首期覆盖至少 20 个已接入集群、多个工作空间/项目团队、`CIS`/`STIG`/平台基线模板、计划性与按需扫描、最近 90 天趋势分析，以及围绕集群、节点、命名空间和关键资源的治理闭环

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- `PASS`: 当前规划工作位于专用分支 `006-compliance-and-hardening`，明确禁止在 `main` 或 `master` 开发。
- `PASS`: 已在计划中定义进入实现前的数据库备份命令、产物路径、凭据差异处理和恢复抽样验证要求；实现前必须补齐 `artifacts/006-compliance-and-hardening/backup-manifest.txt`。
- `PASS`: Go/npm 与扫描器镜像或规则包的国内镜像/代理要求已明确。
- `PASS`: 已定义 GitHub 远程、中文 PR、专业提交说明和“用户明确同意后再合并”的审批门槛。
- `PASS`: 若后续使用子代理，模型固定为 `gpt-5.3-codex`。
- `PASS`: 006 规格未遗留未决澄清标记，可进入设计阶段。

### Post-Design Re-check

- `PASS`: `research.md` 已收敛扫描编排模式、基线版本语义、失败项与证据模型、整改/例外/复检边界和趋势聚合方案。
- `PASS`: `data-model.md` 已覆盖基线标准、扫描配置、扫描执行、失败项、证据、整改任务、例外、复检与趋势快照等核心实体。
- `PASS`: `contracts/openapi.yaml` 已覆盖基线管理、扫描配置、扫描执行、失败项检索、整改流转、例外审批、复检、趋势与审计接口。
- `PASS`: `quickstart.md` 已给出实施前备份、国内源配置、最小联调路径与完整验收清单。

## Project Structure

### Documentation (this feature)

```text
specs/006-compliance-and-hardening/
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
│   │   ├── policy/
│   │   └── compliance/
│   ├── kube/
│   │   ├── adapter/
│   │   ├── cache/
│   │   ├── client/
│   │   └── exec/
│   ├── repository/
│   ├── service/
│   │   ├── auth/
│   │   ├── audit/
│   │   ├── observability/
│   │   ├── workloadops/
│   │   ├── gitops/
│   │   ├── securitypolicy/
│   │   └── compliance/
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
│   │   ├── security-policy/
│   │   └── compliance-hardening/
│   ├── services/
│   │   ├── api/
│   │   └── compliance.ts
│   └── components/
└── tests/
```

**Structure Decision**: 沿用现有前后端分离模块化单体结构。后端新增 `compliance` 业务域与 `integration/compliance` 适配层，承载基线标准、扫描编排、证据汇聚、整改流转与趋势聚合；前端新增 `compliance-hardening` 功能域和 `services/compliance.ts`，承接基线管理、扫描中心、失败项详情、整改/例外交互、复检和趋势汇报入口，避免与 005 的准入策略页和 002 的可观测页职责混杂。

## Complexity Tracking

当前设计不存在需要宪章豁免的复杂度例外。
