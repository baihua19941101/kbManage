# Implementation Plan: 多集群 Kubernetes 可视化管理平台

**Branch**: `001-k8s-ops-platform-followup` | **Date**: 2026-04-09 | **Spec**: [/mnt/e/code/kbmanage/specs/001-k8s-ops-platform/spec.md](/mnt/e/code/kbmanage/specs/001-k8s-ops-platform/spec.md)
**Input**: Feature specification from `/specs/001-k8s-ops-platform/spec.md`

**Note**: 本文件由 `/speckit.plan` 生成，覆盖设计阶段的技术上下文、研究结论和实现边界。

## Summary

构建一个前后端分离的 Web 化 Kubernetes 多集群管理平台，首期提供统一集群接入与资源总览、平台级 RBAC 与工作空间/项目双层权限、受控运维操作执行、审计追踪与导出。首期资源 Kind 固定为 `Deployment`、`StatefulSet`、`DaemonSet`、`Pod`、`Service`、`Ingress`、`Node`、`Namespace`；首批角色矩阵固定为 `platform-admin`、`ops-operator`、`audit-reader`、`readonly`；高风险操作采用“二次确认即执行”；审计导出首期仅支持 CSV 且需敏感字段脱敏。技术上采用 React + TypeScript + Ant Design + Vite 的前端栈，Go + Gin + client-go + GORM 的后端栈，MySQL 负责平台元数据与审计持久化，Redis 负责会话、缓存与异步任务协同。系统按模块化单体实现，通过“平台控制面 + 集群访问/同步模块 + 操作执行模块 + 审计模块”组织代码，暂不纳入 CI/CD、可观测和交付部署设计。

## Technical Context

**Language/Version**: Go 1.25；TypeScript 5.x；React 19.2  
**Primary Dependencies**: Gin；client-go；GORM；React；Vite 8；Ant Design 6.3.x；React Router；TanStack Query  
**Dependency Source**: Go 使用 `GOPROXY=https://goproxy.cn,direct`；前端使用 `https://registry.npmmirror.com`；如需容器基础镜像，优先使用阿里云或 Docker Hub 国内代理镜像  
**Storage**: MySQL 8.4（平台元数据、权限、审计、资源索引）；Redis 8.x（会话、短时缓存、任务协同）  
**Database Backup Plan**: 实施前执行 `mkdir -p artifacts && docker exec mysql8 sh -lc 'mysqldump -hlocalhost -P3306 -uadmin -p123456 --all-databases --single-transaction --quick --routines --events --triggers' > artifacts/mysql-backup-before-001-k8s-ops-platform-$(date +%Y%m%d-%H%M%S).sql`；备份产物落盘到 `/mnt/e/code/kbmanage/artifacts/`；恢复校验通过在 `mysql8` 容器内创建临时库并导入抽样验证结构与核心表。2026-04-10 follow-up 已由 `root` 账号完成一次备份执行并落盘；同时已记录 `admin/123456` 与当前容器配置存在不一致，后续迁移按容器实际凭据执行。  
**Testing**: 前端使用 Vitest + React Testing Library；后端使用 `go test` + Testify；契约测试覆盖 REST API；集成测试覆盖权限、集群接入和操作审计主流程  
**Target Platform**: 桌面浏览器中的 Web 管理台；Linux 容器化后端服务  
**Project Type**: 前后端分离的 Web application，模块化单体  
**Git Workflow**: 当前 follow-up 执行分支为 `001-k8s-ops-platform-followup`（已从原设计分支切分），远程为 `git@github.com:baihua19941101/kbManage.git`；实施阶段继续在该专用分支开发，完成后推送远程并提交中文 PR，合并前必须获得用户明确同意。数据库备份已执行并留档，且已记录 `admin` 凭据与容器配置不一致问题。  
**Performance Goals**: 支持 20 个已接入集群与至少 10,000 个受管资源；资源筛选 1 分钟内完成；95% 常见运维操作 2 分钟内给出明确结果；90% 的 90 天审计查询 30 秒内返回。性能验收证据采用“测试环境压测报告 + 可复现实验脚本”组合  
**Constraints**: 暂不接入 SSO/Keycloak，仅支持用户名密码登录；权限模型固定为平台级 RBAC + 工作空间/项目双层权限；首批角色矩阵固定为 `platform-admin`、`ops-operator`、`audit-reader`、`readonly`；首期资源 Kind 固定为 `Deployment`、`StatefulSet`、`DaemonSet`、`Pod`、`Service`、`Ingress`、`Node`、`Namespace`；高风险动作不引入他人审批流，仅二次确认后执行；审计导出首期仅支持 CSV 且要求敏感字段脱敏；CI/CD、可观测和交付部署不纳入本轮设计；所有设计与交付文档必须中文；子代理如被使用必须固定为 `gpt-5.3-codex`  
**Scale/Scope**: 首期面向内部平台管理员、运维、项目管理员、审计与只读角色；目标覆盖多集群资源浏览、授权、常见运维动作与审计导出，不覆盖应用商店、计费、多云商业集成

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- `PASS`: 当前工作分支为 `001-k8s-ops-platform-followup`，且已完成从原设计分支切分，满足“禁止在 `master/main` 直接开发”的要求。
- `PASS`: 已记录实施前的 MySQL 备份命令、产物路径和恢复校验步骤；2026-04-10 已由 `root` 执行备份并留档，`admin/123456` 与容器实际配置不一致问题已登记。
- `PASS`: 已明确 Go、Node 依赖使用国内镜像/代理源，符合依赖受控安装要求。
- `PASS`: 已记录 GitHub SSH 远程地址、中文 PR 约束和合并需用户批准的流程。
- `PASS`: 设计文档使用中文；如后续使用子代理，计划限定为 `gpt-5.3-codex`。

### Post-Design Re-check

- `PASS`: `research.md` 已解析全部关键技术决策，不存在未解决的澄清项。
- `PASS`: `data-model.md` 已覆盖平台用户、双层权限、集群、资源索引、操作与审计等核心实体。
- `PASS`: `contracts/openapi.yaml` 已定义认证、集群、工作空间、项目、角色绑定、运维操作和审计查询接口。
- `PASS`: `quickstart.md` 已包含国内源设置、MySQL 备份命令、开发顺序和 PR 交付要求。

## Project Structure

### Documentation (this feature)

```text
specs/001-k8s-ops-platform/
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
├── internal/
│   ├── api/
│   │   ├── handler/
│   │   ├── middleware/
│   │   └── router/
│   ├── service/
│   │   ├── auth/
│   │   ├── cluster/
│   │   ├── workspace/
│   │   ├── project/
│   │   ├── operation/
│   │   └── audit/
│   ├── domain/
│   ├── repository/
│   ├── kube/
│   │   ├── client/
│   │   ├── cache/
│   │   └── adapter/
│   ├── integration/
│   │   └── external/
│   └── worker/
├── migrations/
└── tests/
    ├── integration/
    └── contract/

frontend/
├── src/
│   ├── app/
│   ├── pages/
│   ├── components/
│   ├── features/
│   │   ├── auth/
│   │   ├── clusters/
│   │   ├── workspaces/
│   │   ├── projects/
│   │   ├── resources/
│   │   ├── operations/
│   │   └── audit/
│   ├── services/
│   ├── hooks/
│   └── utils/
└── tests/
    ├── unit/
    └── integration/
```

**Structure Decision**: 采用前后端分离的模块化单体。后端按业务域和基础设施分层，避免早期微服务拆分成本；前端按功能域拆分模块，以支持集群管理、权限、运维与审计界面的并行开发。

## Complexity Tracking

当前设计不存在需要特别豁免的宪章违规项。
