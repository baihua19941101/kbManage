# Implementation Plan: 多集群 Kubernetes 集群生命周期中心

**Branch**: `007-cluster-lifecycle` | **Date**: 2026-04-17 | **Spec**: [/mnt/e/code/kbManage/specs/007-cluster-lifecycle/spec.md](/mnt/e/code/kbManage/specs/007-cluster-lifecycle/spec.md)
**Input**: Feature specification from `/specs/007-cluster-lifecycle/spec.md`

**Note**: 本文件由 `/speckit.plan` 生成，覆盖 007 的技术上下文、研究结论、设计边界和实施前置条件。

## 当前执行状态（2026-04-17）

- 已完成 `/speckit.specify`，007 规格与质量清单已生成。
- 当前工作分支为 `007-cluster-lifecycle`。
- 已完成 `/speckit.plan`，本轮生成 `research.md`、`data-model.md`、`contracts/openapi.yaml` 和 `quickstart.md`。
- 007 当前处于 planning 完成、可进入 `/speckit.tasks` 的状态；后续进入实现前仍需落实数据库备份证据、中文 PR、远程推送和用户明确同意后再合并等治理门槛。

## Summary

在保持 001-006 既有主栈和模块化单体结构不变的前提下，为 kbManage 新增一个对标 Rancher Cluster Management 的集群生命周期中心。007 聚焦“导入/注册已有集群 + 驱动/模板化创建 + 节点池管理 + 升级计划 + 停用/退役 + 能力矩阵与驱动扩展管理”，通过统一的生命周期记录、创建前校验、驱动能力抽象和审计链路，形成“接入 -> 创建 -> 变更 -> 升级 -> 停用 -> 退役”的受控闭环。平台继续作为控制面，不承担平台级灾备、统一身份源整合、策略治理或应用市场职责。

## Technical Context

**Language/Version**: Go 1.25；TypeScript 5.x；React 19.2  
**Primary Dependencies**: Gin；client-go；GORM；go-redis；集群驱动访问抽象（导入/注册/创建/升级/节点池操作适配层）；模板与能力矩阵建模抽象；React；Vite 8；Ant Design 6.3.x；React Router；TanStack Query；Zustand  
**Dependency Source**: Go 使用 `GOPROXY=https://goproxy.cn,direct`；前端使用 `https://registry.npmmirror.com`；如需接入基础设施驱动 SDK、集群创建依赖或联调镜像，优先使用阿里云、DaoCloud 或已批准的国内代理，禁止默认境外源直连  
**Storage**: MySQL 8.4（集群生命周期记录、驱动版本、模板、能力矩阵、升级计划、节点池快照、审计索引）；Redis 8.x（创建/升级进度、幂等键、短时校验缓存、异步任务协调）；运行时集群状态与节点信息来自 Kubernetes API 与基础设施驱动适配层  
**Database Backup Plan**: 进入实现前执行 `mkdir -p artifacts/007-cluster-lifecycle && docker exec mysql8 sh -lc 'mysqldump -hlocalhost -P3306 -uadmin -p123456 --all-databases --single-transaction --quick --routines --events --triggers' > artifacts/007-cluster-lifecycle/mysql-backup-$(date +%Y%m%d-%H%M%S).sql`；若当前容器凭据与 `admin/123456` 不一致，必须按实际凭据执行并在 `artifacts/007-cluster-lifecycle/backup-manifest.txt` 记录差异、实际命令与恢复抽样验证结果；恢复校验使用临时 MySQL 容器导入备份并核对数据库与核心表可见性  
**Testing**: 后端使用 `go test ./...` + Testify + `httptest` + fake `client-go` + 驱动适配层 stub；前端使用 Vitest + React Testing Library；契约测试覆盖导入、注册、创建、升级、节点池调整、退役、驱动管理、模板管理和能力矩阵接口；集成测试覆盖创建前校验、权限隔离、冲突动作阻断、升级/退役状态流与审计链路  
**Target Platform**: 桌面浏览器 Web 管理台；Linux 容器化后端服务；可访问 Kubernetes API 和多类基础设施驱动适配层的多集群环境  
**Project Type**: 前后端分离的 Web application，模块化单体控制面  
**Git Workflow**: 当前规划分支为 `007-cluster-lifecycle`，远程为 `git@github.com:baihua19941101/kbManage.git`；后续交付必须推送该分支并提交中文 PR，合并前必须取得用户明确同意，所有提交说明必须专业且可审计  
**Performance Goals**: 90% 的导入/注册状态查询在 30 秒内返回明确结果；90% 的标准化创建或升级请求在 30 分钟内返回明确成功或失败结论；90% 的最近 90 天生命周期审计查询在 30 秒内返回结果集  
**Constraints**: 保持 001-006 既有能力边界与模块化单体结构不变；所有创建、导入、升级、节点池调整和退役动作必须经后端统一授权与审计；首期只聚焦生命周期与能力矩阵，不承担统一身份源、灾备、策略治理和应用市场；驱动扩展必须通过抽象层接入，避免写死单一基础设施实现；若使用子代理必须固定为 `gpt-5.3-codex`  
**Scale/Scope**: 首期覆盖至少 20 个已纳管集群、多种基础设施/驱动类型、模板化创建、升级计划、节点池调整、停用/退役流程，以及网络/存储/身份/监控/安全/备份/发布等能力域的矩阵展示

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- `PASS`: 当前规划工作位于专用分支 `007-cluster-lifecycle`，明确禁止在 `main` 或 `master` 开发。
- `PASS`: 已在计划中定义进入实现前的数据库备份命令、产物路径、凭据差异处理和恢复抽样验证要求；实现前必须补齐 `artifacts/007-cluster-lifecycle/backup-manifest.txt`。
- `PASS`: Go、npm 以及后续驱动适配层相关依赖的国内镜像或代理来源已明确。
- `PASS`: 已定义 GitHub SSH 远程、中文 PR、专业提交说明和“用户明确同意后再合并”的审批门槛。
- `PASS`: 若后续使用子代理或并行代理，模型固定为 `gpt-5.3-codex`。
- `PASS`: 007 规格未遗留未解决的澄清标记，可进入研究与设计阶段。

### Post-Design Re-check

- `PASS`: `research.md` 已收敛导入/注册语义、驱动扩展模型、模板与能力矩阵关系、升级与退役状态机以及权限审计边界。
- `PASS`: `data-model.md` 已覆盖集群生命周期记录、驱动、模板、能力矩阵、升级计划、节点池与审计事件等核心实体。
- `PASS`: `contracts/openapi.yaml` 已覆盖导入、注册、创建、升级、节点池变更、退役、驱动管理、模板管理、能力矩阵与审计接口。
- `PASS`: `quickstart.md` 已写明实施前备份、国内依赖源、最小联调路径与完整验收清单。

## Project Structure

### Documentation (this feature)

```text
specs/007-cluster-lifecycle/
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
│   │   └── clusterlifecycle/
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
│   │   └── clusterlifecycle/
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
│   │   └── cluster-lifecycle/
│   ├── services/
│   │   ├── api/
│   │   └── clusterLifecycle.ts
│   └── components/
└── tests/
```

**Structure Decision**: 继续沿用现有前后端分离模块化单体结构。后端新增 `clusterlifecycle` 业务域与 `integration/clusterlifecycle` 驱动适配层，用于承载导入/注册、模板化创建、升级编排、节点池操作、能力矩阵和退役流程；前端新增 `cluster-lifecycle` 功能域，承接集群生命周期中心、驱动与模板管理、能力矩阵与升级/退役入口，避免混入 001 的通用集群总览页面。

## Complexity Tracking

当前设计不存在需要特别豁免的宪章违规项。
