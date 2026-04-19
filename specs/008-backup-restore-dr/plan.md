# Implementation Plan: 平台级备份恢复与灾备中心

**Branch**: `008-backup-restore-dr` | **Date**: 2026-04-18 | **Spec**: [/mnt/e/code/kbManage/specs/008-backup-restore-dr/spec.md](/mnt/e/code/kbManage/specs/008-backup-restore-dr/spec.md)
**Input**: Feature specification from `/specs/008-backup-restore-dr/spec.md`

**Note**: 本文件由 `/speckit.plan` 生成，覆盖 008 的技术上下文、研究结论、设计边界和实施前置条件。

## Summary

在保持 001-007 既有主栈和模块化单体结构不变的前提下，为 kbManage 新增一个平台级备份恢复与灾备中心。008 聚焦“备份策略与恢复点管理 + 原地恢复/跨集群恢复/环境迁移/定向恢复 + 灾备演练计划/执行/报告”，通过统一的受保护对象目录、恢复点建模、恢复校验、演练记录和审计链路，形成“定义保护范围 -> 生成恢复点 -> 恢复或迁移 -> 演练验证 -> 报告追溯”的企业级闭环。平台继续作为控制面，不承担集群创建导入、策略治理、统一身份源整合或应用发布编排职责。

## Technical Context

**Language/Version**: Go 1.25；TypeScript 5.x；React 19.2  
**Primary Dependencies**: Gin；client-go；GORM；go-redis；平台级备份编排抽象（备份、恢复、迁移、演练执行适配层）；对象范围与一致性评估抽象；React；Vite 8；Ant Design 6.3.x；React Router；TanStack Query；Zustand；Apache ECharts  
**Dependency Source**: Go 使用 `GOPROXY=https://goproxy.cn,direct`；前端使用 `https://registry.npmmirror.com`；如需接入对象存储 SDK、备份执行器、快照工具或演练联调镜像，优先使用阿里云、DaoCloud 或已批准的国内代理，禁止默认境外源直连  
**Storage**: MySQL 8.4（备份策略、恢复点、恢复任务、迁移任务、演练计划、演练记录、验证清单、报告索引、审计索引）；Redis 8.x（备份/恢复进度、短时一致性评估缓存、幂等键、互斥锁、演练步骤协调）；实际备份数据与恢复介质由外部对象存储、快照仓库或执行器适配层保存  
**Database Backup Plan**: 进入实现前执行 `mkdir -p artifacts/008-backup-restore-dr && docker exec mysql8 sh -lc 'mysqldump -hlocalhost -P3306 -uadmin -p123456 --all-databases --single-transaction --quick --routines --events --triggers' > artifacts/008-backup-restore-dr/mysql-backup-$(date +%Y%m%d-%H%M%S).sql`；若当前容器凭据与 `admin/123456` 不一致，必须按实际凭据执行并在 `artifacts/008-backup-restore-dr/backup-manifest.txt` 记录差异、实际命令与恢复抽样验证结果；恢复校验使用临时 MySQL 容器导入备份并核对数据库与核心表可见性  
**Testing**: 后端使用 `go test ./...` + Testify + `httptest` + 备份执行适配层 stub；前端使用 Vitest + React Testing Library；契约测试覆盖策略、恢复点、恢复、迁移、演练、报告和审计接口；集成测试覆盖恢复前校验、跨集群迁移路径、定向恢复、一致性说明、RPO/RTO 记录与权限隔离  
**Target Platform**: 桌面浏览器 Web 管理台；Linux 容器化后端服务；可访问 Kubernetes API、多集群环境及外部备份执行器/对象存储的企业平台  
**Project Type**: 前后端分离的 Web application，模块化单体控制面  
**Git Workflow**: 当前规划分支为 `008-backup-restore-dr`，远程为 `git@github.com:baihua19941101/kbManage.git`；后续交付必须推送该分支并提交中文 PR，合并前必须取得用户明确同意，所有提交说明必须专业且可审计  
**Performance Goals**: 90% 的恢复点查询和策略列表查询在 30 秒内返回结果；90% 的恢复或迁移请求在发起后 30 分钟内返回明确成功、部分成功或失败结论；90% 的演练记录与报告查询在 30 秒内返回目标结果  
**Constraints**: 保持 001-007 既有能力边界与模块化单体结构不变；所有备份、恢复、迁移和演练动作必须经后端统一授权与审计；首期只聚焦平台级备份恢复、跨集群迁移和灾备演练，不承担集群创建导入、策略治理、统一身份源整合和应用发布编排；如使用子代理必须固定为 `gpt-5.3-codex`  
**Scale/Scope**: 首期覆盖至少 3 类平台对象、10 个关键业务命名空间、多个目标集群和演练计划/演练报告闭环；支持平台元数据、权限配置、审计记录、集群配置和关键业务命名空间的保护范围

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- `PASS`: 当前规划工作位于专用分支 `008-backup-restore-dr`，明确禁止在 `main` 或 `master` 开发。
- `PASS`: 已在计划中定义进入实现前的数据库备份命令、产物路径、凭据差异处理和恢复抽样验证要求；实现前必须补齐 `artifacts/008-backup-restore-dr/backup-manifest.txt`。
- `PASS`: Go、npm 以及后续备份执行器/对象存储相关依赖的国内镜像或代理来源已明确。
- `PASS`: 已定义 GitHub SSH 远程、中文 PR、专业提交说明和“用户明确同意后再合并”的审批门槛。
- `PASS`: 若后续使用子代理或并行代理，模型固定为 `gpt-5.3-codex`，并以宪章为准。
- `PASS`: 008 规格未遗留未解决的澄清标记，可进入研究与设计阶段。

### Post-Design Re-check

- `PASS`: `research.md` 已收敛受保护对象边界、恢复点一致性说明、跨集群迁移策略、演练模型和授权审计边界。
- `PASS`: `data-model.md` 已覆盖备份策略、恢复点、恢复/迁移任务、演练计划、演练记录、报告和审计事件等核心实体。
- `PASS`: `contracts/openapi.yaml` 已覆盖策略、恢复点、恢复、迁移、演练、报告和审计接口。
- `PASS`: `quickstart.md` 已写明实施前备份、国内依赖源、最小联调路径与完整验收清单。

## Project Structure

### Documentation (this feature)

```text
specs/008-backup-restore-dr/
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
│   │   ├── clusterlifecycle/
│   │   └── backuprestore/
│   ├── kube/
│   │   ├── adapter/
│   │   ├── cache/
│   │   ├── client/
│   │   └── exec/
│   ├── repository/
│   ├── service/
│   │   ├── auth/
│   │   ├── audit/
│   │   ├── cluster/
│   │   ├── observability/
│   │   ├── workloadops/
│   │   ├── gitops/
│   │   ├── securitypolicy/
│   │   ├── compliance/
│   │   ├── clusterlifecycle/
│   │   └── backuprestore/
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
│   │   ├── gitops/
│   │   ├── security-policy/
│   │   ├── compliance-hardening/
│   │   ├── cluster-lifecycle/
│   │   └── backup-restore-dr/
│   ├── services/
│   │   ├── api/
│   │   └── backupRestore.ts
│   └── components/
└── tests/
```

**Structure Decision**: 继续沿用现有前后端分离模块化单体结构。后端新增 `backuprestore` 业务域与 `integration/backuprestore` 备份执行适配层，用于承载备份策略、恢复点、恢复/迁移、演练与报告流程；前端新增 `backup-restore-dr` 功能域，承接平台级备份恢复中心、迁移和演练入口，避免混入 006 的合规扫描或 007 的集群生命周期页面。

## Complexity Tracking

当前设计不存在需要特别豁免的宪章违规项。
