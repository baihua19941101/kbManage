# Implementation Plan: 身份与多租户治理中心

**Branch**: `009-identity-tenancy` | **Date**: 2026-04-19 | **Spec**: [/mnt/e/code/kbManage/specs/009-identity-tenancy/spec.md](/mnt/e/code/kbManage/specs/009-identity-tenancy/spec.md)
**Input**: Feature specification from `/specs/009-identity-tenancy/spec.md`

**Note**: 本文件由 `/speckit.plan` 生成，覆盖 009 的技术上下文、研究结论、设计边界和实施前置条件。

## Summary

在保持 001-008 既有主栈和模块化单体结构不变的前提下，为 kbManage 新增一个身份与多租户治理中心。009 聚焦“外部身份源接入 + 组织/团队/用户组/工作空间/项目关系建模 + 平台级到资源级细粒度 RBAC + 委派、临时授权、会话治理和访问风险视图”，通过统一身份目录、租户边界映射、授权继承和审计追踪，形成“接入身份 -> 建立组织边界 -> 分配角色 -> 委派/回收 -> 风险识别”的企业级访问治理闭环。平台继续作为控制面，不承担准入策略、合规扫描、应用发布和灾备恢复职责。

## Technical Context

**Language/Version**: Go 1.25；TypeScript 5.x；React 19.2  
**Primary Dependencies**: Gin；client-go；GORM；go-redis；身份源接入抽象（SSO、OIDC、LDAP、本地账号）；组织关系与授权评估抽象；React；Vite 8；Ant Design 6.3.x；React Router；TanStack Query；Zustand  
**Dependency Source**: Go 使用 `GOPROXY=https://goproxy.cn,direct`；前端使用 `https://registry.npmmirror.com`；如需接入身份目录 SDK、SSO 联调镜像或 LDAP 客户端组件，优先使用阿里云、DaoCloud 或已批准的国内代理，禁止默认境外源直连  
**Storage**: MySQL 8.4（身份源配置、组织模型、团队/用户组关系、角色定义、授权分配、委派记录、临时授权、会话治理索引、风险快照、审计索引）；Redis 8.x（登录会话、短时授权缓存、权限评估缓存、委派到期索引、会话回收协调）；外部身份目录和组织来源数据由平台外部身份源系统保存  
**Database Backup Plan**: 进入实现前执行 `mkdir -p artifacts/009-identity-tenancy && docker exec mysql8 sh -lc 'mysqldump -hlocalhost -P3306 -uadmin -p123456 --all-databases --single-transaction --quick --routines --events --triggers' > artifacts/009-identity-tenancy/mysql-backup-$(date +%Y%m%d-%H%M%S).sql`；若当前容器凭据与 `admin/123456` 不一致，必须按实际凭据执行并在 `artifacts/009-identity-tenancy/backup-manifest.txt` 记录差异、实际命令与恢复抽样验证结果；恢复校验使用临时 MySQL 容器导入备份并核对数据库与核心表可见性  
**Testing**: 后端使用 `go test ./...` + Testify + `httptest` + 身份源适配层 stub；前端使用 Vitest + React Testing Library；契约测试覆盖身份源、组织、角色、授权、委派、会话与风险接口；集成测试覆盖身份源切换、租户边界、授权继承、委派回收、临时授权到期和会话治理路径  
**Target Platform**: 桌面浏览器 Web 管理台；Linux 容器化后端服务；可访问多集群 Kubernetes 控制面和外部身份源的企业平台  
**Project Type**: 前后端分离的 Web application，模块化单体控制面  
**Git Workflow**: 当前规划分支为 `009-identity-tenancy`，远程为 `git@github.com:baihua19941101/kbManage.git`；后续交付必须推送该分支并提交中文 PR，合并前必须取得用户明确同意，所有提交说明必须专业且可审计  
**Performance Goals**: 90% 的用户权限查询、角色边界查看和风险摘要查询在 30 秒内返回结果；90% 的身份源登录与会话状态查询在 30 秒内返回明确结论；90% 的授权变更与访问回收在发起后 5 分钟内反映到有效权限视图  
**Constraints**: 保持 001-008 既有能力边界与模块化单体结构不变；所有身份源管理、组织变更、授权委派、临时授权和回收动作必须经后端统一授权与审计；首期只聚焦身份接入、组织模型和细粒度 RBAC，不承担准入策略、合规扫描、应用发布和灾备恢复；如使用子代理必须固定为 `gpt-5.3-codex`  
**Scale/Scope**: 首期覆盖至少 3 类身份来源、多个组织/团队/用户组层级、平台级到资源级角色体系，以及多租户工作空间和项目授权边界治理

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- `PASS`: 当前规划工作位于专用分支 `009-identity-tenancy`，明确禁止在 `main` 或 `master` 开发。
- `PASS`: 已在计划中定义进入实现前的数据库备份命令、产物路径、凭据差异处理和恢复抽样验证要求；实现前必须补齐 `artifacts/009-identity-tenancy/backup-manifest.txt`。
- `PASS`: Go、npm 以及后续身份源 SDK 或联调镜像相关依赖的国内镜像或代理来源已明确。
- `PASS`: 已定义 GitHub SSH 远程、中文 PR、专业提交说明和“用户明确同意后再合并”的审批门槛。
- `PASS`: 若后续使用子代理或并行代理，模型固定为 `gpt-5.3-codex`，并以宪章为准。
- `PASS`: 009 规格未遗留未解决的澄清标记，可进入研究与设计阶段。

### Post-Design Re-check

- `PASS`: `research.md` 已收敛身份源边界、组织模型、授权继承、委派治理、会话控制与风险可视化语义。
- `PASS`: `data-model.md` 已覆盖身份源、组织单元、租户映射、角色定义、授权分配、委派关系、会话记录与风险快照等核心实体。
- `PASS`: `contracts/openapi.yaml` 已覆盖身份源、组织、角色、授权、会话和风险接口。
- `PASS`: `quickstart.md` 已写明实施前备份、国内依赖源、最小联调路径与完整验收清单。

## Project Structure

### Documentation (this feature)

```text
specs/009-identity-tenancy/
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
│   │   ├── backuprestore/
│   │   └── identity/
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
│   │   ├── backuprestore/
│   │   └── identitytenancy/
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
│   │   ├── backup-restore-dr/
│   │   └── identity-tenancy/
│   ├── services/
│   │   ├── api/
│   │   └── identityTenancy.ts
│   └── components/
└── tests/
```

**Structure Decision**: 继续沿用现有前后端分离模块化单体结构。后端新增 `identitytenancy` 业务域与 `integration/identity` 身份源适配层，用于承载身份源接入、组织建模、角色授权、委派与会话治理流程；前端新增 `identity-tenancy` 功能域，承接统一身份与多租户治理入口，避免混入现有 `auth` 基础登录页或 005 的策略治理页面。

## Complexity Tracking

当前设计不存在需要特别豁免的宪章违规项。
