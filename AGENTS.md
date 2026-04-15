# kbManage Development Guidelines

Auto-generated from all feature plans. Last updated: 2026-04-14

## Active Technologies
- Go 1.25；TypeScript 5.x；React 19.2 + Gin；client-go；GORM；React；Vite 8；Ant Design 6.3.x；React Router；TanStack Query (001-k8s-ops-platform)
- MySQL 8.4（平台元数据、权限、审计、资源索引）；Redis 8.x（会话、短时缓存、任务协同） (001-k8s-ops-platform)
- Go 1.25；TypeScript 5.x；React 19.2 + Gin；client-go；GORM；go-redis；React；Vite 8；Ant Design 6.3.x；React Router；TanStack Query；Zustand；Apache ECharts；Prometheus-compatible Query API；Alertmanager-compatible API；Loki-compatible Query API (002-observability-center)
- MySQL 8.4（可观测数据源配置、告警规则治理元数据、通知目标、静默窗口、告警快照、处理记录、审计索引）；Redis 8.x（查询缓存、同步游标、短时上下文、告警状态协同）；原始日志/指标/运行态告警由外部 Prometheus/Alertmanager/Loki 兼容后端保存 (002-observability-center)
- MySQL 8.4（工作负载动作请求、批量任务、发布历史快照索引、终端会话审计、扩展运维审计索引）；Redis 8.x（动作进度缓存、批量执行协调、终端短会话上下文、幂等键和短时结果缓存）；运行时工作负载状态、Pod/容器状态、日志流和 exec 通道来自 Kubernetes API (003-workload-operations-control-plane)
- Go 1.25；TypeScript 5.x；React 19.2 + Gin；client-go；GORM；go-redis；go-git 风格 Git 访问抽象；Helm SDK 风格发布源抽象；React；Vite 8；Ant Design 6.3.x；React Router；TanStack Query；Zustand (004-gitops-and-release)
- MySQL 8.4（交付来源、目标组、环境阶段、配置覆盖、发布历史、同步/发布动作、审计索引）；Redis 8.x（动作进度、分布式锁、幂等键、短时差异缓存、阶段推进上下文）；Git/Helm/OCI 来源内容与 Kubernetes 实际状态由外部来源和 Kubernetes API 提供 (004-gitops-and-release)
- Go 1.25；TypeScript 5.x；React 19.2 + Gin；client-go；GORM；go-redis；策略评估与准入执行抽象（平台内部策略模板 + 规则执行器适配层）；React；Vite 8；Ant Design 6.3.x；React Router；TanStack Query；Zustand (005-security-and-policy)
- MySQL 8.4（策略定义、策略分配、命中记录、例外申请、整改状态、审计索引）；Redis 8.x（策略分发进度、短时命中缓存、例外时效索引、幂等键）；运行时准入与对象状态来自 Kubernetes API (005-security-and-policy)

## Project Structure

```text
backend/
frontend/
specs/
```

## Commands

npm test && npm run lint

## Code Style

Go 1.25；TypeScript 5.x；React 19.2: Follow standard conventions

## Recent Changes
- 005-security-and-policy: Added Go 1.25；TypeScript 5.x；React 19.2 + Gin；client-go；GORM；go-redis；策略评估与准入执行抽象（平台内部策略模板 + 规则执行器适配层）；React；Vite 8；Ant Design 6.3.x；React Router；TanStack Query；Zustand
- 004-gitops-and-release: Added Go 1.25；TypeScript 5.x；React 19.2 + Gin；client-go；GORM；go-redis；go-git 风格 Git 访问抽象；Helm SDK 风格发布源抽象；React；Vite 8；Ant Design 6.3.x；React Router；TanStack Query；Zustand
- 003-workload-operations-control-plane: Added Go 1.25；TypeScript 5.x；React 19.2 + Gin；client-go；GORM；go-redis；React；Vite 8；Ant Design 6.3.x；React Router；TanStack Query；Zustand

<!-- MANUAL ADDITIONS START -->
## Configuration Conventions

- 后端所有可配置项必须集中在 `backend/config/*.yaml`，默认使用 `backend/config/config.dev.yaml`。
- 后端配置文件必须带注释，新增配置项时同步更新 `backend/config/config.example.yaml`。
- 后端仅允许通过 `CONFIG_FILE` 指定配置文件路径；如需临时调试可使用同名环境变量覆盖单项。
- 前端配置统一放在 `frontend/.env.development`（开发）与 `frontend/.env.example`（模板）。
- 前端启动端口必须通过 `VITE_PORT` 配置，不允许在脚本中硬编码固定端口。
- README 必须同步维护“配置说明 + 启动命令 + 端口配置方式”。
<!-- MANUAL ADDITIONS END -->


## 协作代理规则
- 必须使用中文与用户交流
- 每次新会话都必须读取spec了解项目功能和计划
- 每完成一个任务，必须更新spec
- 如果有新的需求或改动，必须先和用户进行讨论，经过用户允许之后必须立刻更新spec
- 完成新的需求或改动，必须立刻更新spec
- 所有语言框架安装依赖时必须使用国内源或代理下载

## 开发流程规则
- 新功能必须使用分支开发模式，从远程的main分支拉取代码，禁止在主分支进行开发，开发完成后，提交github
- 合并分支必须经过用户同意后才能进行合并
- 提交代码必须包含专业的说明描述

## 数据库备份流程
- 开发新功能前必须保证数据库备份完成，才能进行新功能开发
- 数据库备份操作可以登录到docker容器内进行
