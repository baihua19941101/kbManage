# Quickstart: 多集群 Kubernetes 可观测中心

## 1. 前提
- 当前工作分支必须是 `002-observability-center` 或其后续专用实现分支，禁止在 `main` 上直接开发。
- 001 已确认合并到 `main`，002 可进入实现；仍需持续遵守“中文 PR + 用户明确同意后再合并”的宪章要求。
- 任何数据库结构或治理元数据落库改动前，都必须执行新的数据库备份。

## 2. 实施前数据库备份
在开始实现 002 的表结构、配置项和同步任务前，执行以下命令：

```bash
mkdir -p artifacts/002-observability-center

docker exec mysql8 sh -lc 'mysqldump -hlocalhost -P3306 -uadmin -p123456 --all-databases --single-transaction --quick --routines --events --triggers' \
  > artifacts/002-observability-center/mysql-backup-$(date +%Y%m%d-%H%M%S).sql
```

恢复抽样验证建议：

```bash
docker exec -i mysql8 sh -lc 'mysql -hlocalhost -P3306 -uadmin -p123456 -e "CREATE DATABASE IF NOT EXISTS restore_check_002;"'
docker exec -i mysql8 sh -lc 'mysql -hlocalhost -P3306 -uadmin -p123456 restore_check_002' < artifacts/002-observability-center/<backup-file>.sql
```

说明：
- `artifacts/001-k8s-ops-platform/` 中已有历史备份证据，可作为参考。
- 002 进入实现前必须重新执行本轮备份并记录文件名，不可直接复用 001 的历史备份结论。

## 2026-04-11 备份执行记录
- 分支：`002-observability-center`
- 备份文件名：`artifacts/002-observability-center/mysql-backup-20260411-214819.sql`
- 凭据说明：沿用 001 的历史结论，`admin/123456` 在当前容器环境不可用；本次实际使用 `root/123456` 执行备份
- 恢复校验：已通过临时容器 `mysql8-restore-check-002` 完成完整导入，并核对数据库列表与抽样表计数
- 详细记录：`artifacts/002-observability-center/backup-manifest.txt`

## 2026-04-11 实现启动记录
- 前置 gate：`001` 已合并 `main`，允许启动 002 实现
- 当前阶段：优先执行 `tasks.md` 的 Phase 0（剩余治理项）+ Phase 1（Setup）+ Phase 2（Foundational 骨架）
- 约束保持不变：首期仅实现观测中心范围，不引入终端、批量操作、回滚、GitOps、Helm、策略治理、合规扫描、集群生命周期和灾备能力

## 2026-04-11 治理门槛记录
- 分支与合并门槛：`artifacts/002-observability-center/branch-check.txt`
- 国内源与远程流程：`artifacts/002-observability-center/mirror-and-remote-check.txt`

## 3. 联调模式建议

### 模式 A：最小监控展示联调
适用于尽快验证资源上下文、指标概览和事件时间线。

必需条件：
- 至少一个已接入集群
- 一个 Prometheus-compatible 指标查询端点
- Kubernetes Event 可读权限

适用范围：
- 指标概览与趋势
- 事件时间线
- 资源上下文联动骨架

不覆盖：
- 日志检索
- 通知目标
- 静默窗口
- 完整告警治理闭环

### 模式 B：002 完整验收联调
适用于按 002 规格做完整验收。

必需条件：
- 至少一个 Prometheus-compatible 指标端点
- 至少一个 Alertmanager-compatible 告警端点
- 至少一个 Loki-compatible 日志端点
- Kubernetes Event 可读权限

适用范围：
- 日志、事件、指标、告警统一联动
- 告警规则治理
- 通知目标
- 静默窗口
- 告警确认与处理记录

## 4. 后端开发环境
建议继续在 `backend/` 目录内按 001 既有结构扩展 `observability` 域。

依赖下载使用国内源：

```bash
export GOPROXY=https://goproxy.cn,direct
```

如果需要为联调环境拉取可观测相关镜像：
- 优先使用阿里云、DaoCloud 或已批准的国内代理镜像
- 禁止在未说明代理或镜像源的情况下直接使用默认境外拉取

建议优先实现的后端模块：
- `observability/datasource`: 数据源配置、连通性校验、适配器路由
- `observability/metrics`: 统一概览、趋势查询、资源上下文指标摘要
- `observability/events`: Event 时间线查询、短时缓存、错误归一化
- `observability/logs`: Loki 兼容查询、时间范围约束、上下文透传
- `observability/alerts`: 告警列表、告警详情、确认、静默、处理记录
- `observability/admin`: 规则、通知目标和静默窗口治理

建议新增的后端配置项统一放到 `backend/config/config.dev.yaml` 与 `backend/config/config.example.yaml`：
- `observability.metrics.*`
- `observability.logs.*`
- `observability.alerts.*`
- `observability.cache.*`

## 5. 前端开发环境
依赖下载使用国内源：

```bash
cd frontend
npm config set registry https://registry.npmmirror.com
```

建议优先实现的前端页面与模块：
- 可观测总览页
- 日志检索页
- 事件时间线页
- 指标趋势与健康页
- 告警中心页
- 告警规则/通知目标/静默窗口管理页

首期视觉与交互原则：
- 所有视图围绕资源上下文统一跳转和筛选
- 图表、表格和时间线保持统一筛选状态
- 日志首期以查询与联动为主，不做终端式实时 tail

## 6. 推荐实施顺序
1. 落库数据源配置、规则、通知目标、静默窗口、告警快照和处理记录。
2. 完成数据源配置管理、统一授权校验和后端适配器骨架。
3. 完成指标概览与 Event 时间线，形成最小可见性入口。
4. 完成日志检索与资源上下文联动。
5. 完成告警中心、规则治理、通知目标和静默窗口。
6. 补齐审计、异常降级、缓存和契约/集成测试。

## 7. 最小验收清单
- 授权用户能够从目标资源进入统一的可观测上下文页。
- 能按集群、工作空间、项目、命名空间、工作负载、Pod、容器和时间范围查询日志。
- 能查看 Kubernetes Event 时间线并区分正常/告警事件。
- 能查看集群、节点、命名空间和工作负载的健康状态与趋势变化。
- 能看到当前告警列表、详情、时间线和处理记录。
- 能创建至少一条告警规则、一个通知目标和一个静默窗口，并看到治理状态变化。
- 未授权用户无法访问其他工作空间或项目范围内的观测数据。

## 8. PR 交付要求
- 所有提交说明必须专业、中文、可审计。
- 推送远程后必须创建中文 PR，并附上数据库备份证据、验证结果和风险说明。
- 合并前必须取得用户明确同意。
