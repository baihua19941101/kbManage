# kbManage Development Guidelines

Auto-generated from all feature plans. Last updated: 2026-04-19

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
- Go 1.25；TypeScript 5.x；React 19.2 + Gin；client-go；GORM；go-redis；合规扫描执行抽象（平台编排 + 外部扫描器/基线包适配层）；React；Vite 8；Ant Design 6.3.x；React Router；TanStack Query；Zustand；Apache ECharts (006-compliance-and-hardening)
- MySQL 8.4（基线标准、扫描配置、扫描执行、失败项、证据索引、整改任务、例外审批、复检任务、趋势快照、审计索引）；Redis 8.x（扫描调度队列、进度缓存、短时证据缓存、幂等键、复检协调）；运行时证据与原始检查结果来自 Kubernetes API 与外部扫描器/基线执行适配层 (006-compliance-and-hardening)
- Go 1.25；TypeScript 5.x；React 19.2 + Gin；client-go；GORM；go-redis；集群驱动访问抽象（导入/注册/创建/升级/节点池操作适配层）；模板与能力矩阵建模抽象；React；Vite 8；Ant Design 6.3.x；React Router；TanStack Query；Zustand (007-cluster-lifecycle)
- MySQL 8.4（集群生命周期记录、驱动版本、模板、能力矩阵、升级计划、节点池快照、审计索引）；Redis 8.x（创建/升级进度、幂等键、短时校验缓存、异步任务协调）；运行时集群状态与节点信息来自 Kubernetes API 与基础设施驱动适配层 (007-cluster-lifecycle)
- Go 1.25；TypeScript 5.x；React 19.2 + Gin；client-go；GORM；go-redis；平台级备份编排抽象（备份、恢复、迁移、演练执行适配层）；对象范围与一致性评估抽象；React；Vite 8；Ant Design 6.3.x；React Router；TanStack Query；Zustand；Apache ECharts (008-backup-restore-dr)
- MySQL 8.4（备份策略、恢复点、恢复任务、迁移任务、演练计划、演练记录、验证清单、报告索引、审计索引）；Redis 8.x（备份/恢复进度、短时一致性评估缓存、幂等键、互斥锁、演练步骤协调）；实际备份数据与恢复介质由外部对象存储、快照仓库或执行器适配层保存 (008-backup-restore-dr)
- Go 1.25；TypeScript 5.x；React 19.2 + Gin；client-go；GORM；go-redis；身份源接入抽象（SSO、OIDC、LDAP、本地账号）；组织关系与授权评估抽象；React；Vite 8；Ant Design 6.3.x；React Router；TanStack Query；Zustand (009-identity-tenancy)
- MySQL 8.4（身份源配置、组织模型、团队/用户组关系、角色定义、授权分配、委派记录、临时授权、会话治理索引、风险快照、审计索引）；Redis 8.x（登录会话、短时授权缓存、权限评估缓存、委派到期索引、会话回收协调）；外部身份目录和组织来源数据由平台外部身份源系统保存 (009-identity-tenancy)
- Go 1.25；TypeScript 5.x；React 19.2 + Gin；GORM；go-redis；应用目录来源抽象（Git/Helm/OCI 风格目录访问）；扩展注册与兼容性评估抽象；Ant Design 6.3.x；React Router；TanStack Query；Zustand (010-platform-marketplace)
- MySQL 8.4（目录来源、模板、版本、分发、安装记录、扩展包、兼容性结论、审计索引）；Redis 8.x（目录同步游标、短时模板缓存、分发协调、幂等键、兼容性结果缓存） (010-platform-marketplace)
- Go 1.25；TypeScript 5.x；React 19.2 + Gin；GORM；go-redis；平台运行健康聚合抽象；升级前检查与版本兼容性评估抽象；容量预测与压测证据接入抽象；Ant Design 6.3.x；React Router；TanStack Query；Zustand；Apache ECharts (011-sre-scale)
- MySQL 8.4（高可用策略、维护窗口、平台组件健康快照索引、容量基线、升级计划、回退验证、运行手册、告警基线、自诊断摘要、审计索引）；Redis 8.x（健康聚合缓存、任务积压短时状态、升级协调、限流状态、幂等键、容量预测缓存） (011-sre-scale)
- Go 1.25；TypeScript 5.x；React 19.2 + Gin；GORM；go-redis；现有权限审计聚合抽象；治理覆盖率快照与报表生成抽象；导出记录与交付包编排抽象；Ant Design 6.3.x；React Router；TanStack Query；Zustand；Apache ECharts (012-enterprise-polish)
- MySQL 8.4（权限变更链路、关键操作追踪索引、跨团队授权分布快照、治理覆盖率快照、治理报表、导出记录、交付材料目录、交付检查清单、审计索引）；Redis 8.x（报表生成上下文、导出任务短时状态、权限风险聚合缓存、幂等键、趋势查询缓存） (012-enterprise-polish)

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
- 012-enterprise-polish: Added Go 1.25；TypeScript 5.x；React 19.2 + Gin；GORM；go-redis；现有权限审计聚合抽象；治理覆盖率快照与报表生成抽象；导出记录与交付包编排抽象；Ant Design 6.3.x；React Router；TanStack Query；Zustand；Apache ECharts
- 011-sre-scale: Added Go 1.25；TypeScript 5.x；React 19.2 + Gin；GORM；go-redis；平台运行健康聚合抽象；升级前检查与版本兼容性评估抽象；容量预测与压测证据接入抽象；Ant Design 6.3.x；React Router；TanStack Query；Zustand；Apache ECharts
- 010-platform-marketplace: Added Go 1.25；TypeScript 5.x；React 19.2 + Gin；GORM；go-redis；应用目录来源抽象（Git/Helm/OCI 风格目录访问）；扩展注册与兼容性评估抽象；Ant Design 6.3.x；React Router；TanStack Query；Zustand

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
