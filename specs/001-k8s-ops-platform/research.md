# Research: 多集群 Kubernetes 可视化管理平台

## 决策 1：前端采用 React 19 + TypeScript + Vite 8 + Ant Design 6
- Decision: 使用 React 19.2、TypeScript 5.x、Vite 8 和 Ant Design 6.3.x 作为控制台前端基础栈，并补充 React Router 与 TanStack Query 支撑路由和服务端状态管理。
- Rationale: React 19 代表当前官方文档主版本，Vite 8 是当前受支持版本，Ant Design 6 当前官方文档稳定版本为 6.3.5，明确面向企业级 Web 应用并提供 Vite 集成方式；这套组合最适合资源表格、树结构、表单、抽屉和权限后台类界面。
- Alternatives considered: Next.js 被排除，因为本项目不是 SEO/SSR 驱动型站点；Vue + Element Plus 被排除，因为你已经明确指定 React 路线。

## 决策 2：后端采用 Go 1.25 + Gin + client-go + GORM 的模块化单体
- Decision: 后端使用 Go 1.25，HTTP 层采用 Gin，Kubernetes 集成层使用 client-go，数据库访问层采用 GORM；整体保持前后端分离的模块化单体，不拆分微服务。
- Rationale: Go 是 Kubernetes 官方一等客户端生态语言，client-go 是官方维护的 Go 客户端；Gin 提供成熟的 HTTP 路由、中间件和 JSON 绑定；GORM 适合快速落地关系模型与事务逻辑。模块化单体能在保留清晰边界的同时降低分布式复杂度。
- Alternatives considered: 纯 net/http 被排除，因为在认证、中间件、统一错误处理上开发效率更低；微服务被排除，因为本阶段目标是尽快形成可用控制面，不需要过早引入服务治理成本。

## 决策 3：集群资源读取采用“平台元数据 + 实时访问 + 索引缓存”混合模式
- Decision: 平台基础数据（用户、权限、工作空间、项目、审计）存 MySQL；Kubernetes 资源浏览采用 client-go 的 list/watch 和 informer 模式维护资源索引，再在资源详情和高风险动作时向目标集群实时确认。
- Rationale: Kubernetes API 官方推荐客户端使用 list/watch 追踪变化。对 20 集群、10k 资源的场景，纯实时透传会导致列表筛选和审计关联性能不稳定；纯离线快照又会削弱操作时效性。混合模式更平衡。
- Alternatives considered: 全量实时代理被排除，因为跨集群筛选和分页成本高；全量数据库镜像被排除，因为实现和一致性成本过高。

## 决策 4：权限模型采用平台级 RBAC + 工作空间/项目双层权限
- Decision: 平台层负责全局角色与治理能力，工作空间/项目层负责业务边界内的资源可见性与可执行动作。
- Rationale: 平台存在两类权限问题：全局治理（集群接入、平台用户、全局审计）和业务隔离（谁能看哪个团队、哪个项目、哪个环境）。单一平台 RBAC 难以长期维持清晰边界，双层模型更符合 Rancher 类平台演进路径。
- Alternatives considered: 纯平台 RBAC 被排除，因为随着团队增多会出现角色爆炸；只做工作空间级权限被排除，因为平台级治理能力无法单独建模。

## 决策 5：认证先采用用户名密码登录，并使用访问令牌 + 刷新令牌机制
- Decision: 首期不接入 SSO 或 Keycloak，而是使用平台内建用户名密码认证；服务端保存密码哈希，登录后下发短期访问令牌和可撤销刷新令牌。
- Rationale: 你已明确要求首期不接入外部认证中心。前后端分离场景下，访问令牌 + 刷新令牌便于兼顾用户体验和会话可控性，也便于与 Redis 配合做会话失效和强制退出。
- Alternatives considered: 纯 Cookie Session 被排除，因为后续多实例与 API 客户端扩展时灵活性较弱；SSO/Keycloak 被排除，因为不在当前范围。

## 决策 6：MySQL 8.4 作为主库，Redis 8 用于会话、缓存和任务协同
- Decision: 采用 MySQL 8.4 存储平台主数据、角色绑定、审计事件和资源索引；采用 Redis 8.x 存储登录会话、短时缓存、操作去重锁和异步任务协调信息。
- Rationale: 你已指定 MySQL + Redis。MySQL 8.4 提供长期稳定的关系存储能力；Redis 的数据结构和流式能力适合会话与任务状态编排。
- Alternatives considered: PostgreSQL 被排除，因为与用户指定数据库不一致；只用 MySQL 不用 Redis 被排除，因为登录态、限流和异步协调会增加数据库压力。

## 决策 7：外部自定义后端服务采用适配器接口接入，不直接耦合到核心领域层
- Decision: 在后端保留 `integration/external` 和 `kube/adapter` 层，为 Kubernetes API 与自定义后端服务提供统一适配接口。
- Rationale: 规格里已经允许“对接 Kubernetes API / 自定义后端服务”。适配器层可以把集群访问、扩展运维动作和后续第三方能力隔离开，避免核心服务层被外部协议污染。
- Alternatives considered: 直接在业务服务中调用外部 HTTP 接口被排除，因为后期维护和测试成本过高。

## 决策 8：本轮设计明确排除 CI/CD、可观测与交付部署细节
- Decision: 在 plan、contract 和 quickstart 中仅保留最小化运行说明，不设计 Argo CD、Prometheus、Grafana、Loki 或生产部署方案。
- Rationale: 你已明确这几部分“暂不考虑”。过早设计会扩散当前需求范围，影响任务聚焦。
- Alternatives considered: 同步设计部署与监控方案被排除，因为与当前需求边界不一致。
