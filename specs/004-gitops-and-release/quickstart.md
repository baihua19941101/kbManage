# Quickstart: 多集群 GitOps 与应用发布中心

## 当前执行状态（2026-04-13）

- 004 已完成 US1/US2/US3，进入 Final Phase。
- 治理与配置前置任务 T001/T002/T003/T004/T008 已完成。
- 当前重点为收尾文档、验证基线与 PR 交付材料，仍需遵守“中文 PR + 用户明确同意后再合并”。

## 1. 前提

- 当前工作分支必须是 `004-gitops-and-release` 或其后续专用实现分支，禁止在 `main` 上直接开发。
- `001`、`002`、`003` 已形成 004 的基础能力来源；004 继续遵守“中文 PR + 用户明确同意后再合并”的宪章要求。
- 任何数据库结构、交付来源、交付单元、环境推进记录或发布历史模型变更前，都必须执行新的数据库备份。

## 2. 实施前数据库备份

004 的备份与恢复抽样验证已在规划阶段执行完成。

执行命令：

```bash
mkdir -p artifacts/004-gitops-and-release

docker exec mysql8 sh -lc 'mysqldump -hlocalhost -P3306 -uroot -p123456 --all-databases --single-transaction --quick --routines --events --triggers' \
  > artifacts/004-gitops-and-release/mysql-backup-20260413-155243.sql
```

恢复抽样验证：

```bash
docker run -d --rm --name mysql8-restore-check-004 -e MYSQL_ROOT_PASSWORD=123456 mysql:8.0
cat artifacts/004-gitops-and-release/mysql-backup-20260413-155243.sql \
  | docker exec -i mysql8-restore-check-004 sh -lc 'MYSQL_PWD=123456 mysql -uroot'
docker exec mysql8-restore-check-004 sh -lc 'MYSQL_PWD=123456 mysql -N -B -uroot -e "SHOW DATABASES;"'
```

## 2026-04-13 备份执行记录

- 分支：`004-gitops-and-release`
- 备份文件名：`artifacts/004-gitops-and-release/mysql-backup-20260413-155243.sql`
- 凭据说明：`admin/123456` 在当前容器环境不可用；本次实际使用 `root/123456` 执行备份
- 恢复校验：已通过临时容器 `mysql8-restore-check-004` 完成完整导入，并成功列出数据库
- 详细记录：`artifacts/004-gitops-and-release/backup-manifest.txt`

## 3. 联调模式建议

### 模式 A：单应用交付单元最小闭环联调

适用于尽快验证来源接入、交付单元详情、状态/差异和一次基础同步或升级动作。

必需条件：
- 至少一个已接入且已授权的集群
- 至少一个可访问的 Git 仓库或应用发布来源
- 至少一个工作空间/项目和一个目标命名空间
- 平台可访问 Kubernetes API，并能够读取目标对象运行状态

适用范围：
- 来源管理
- 交付单元创建
- 状态聚合
- 差异查看
- 单环境安装/升级/同步

不覆盖：
- 多环境推进
- 大规模目标组发布
- 生产级权限回收验收

### 模式 B：004 完整验收联调

适用于按 004 规格做完整验收。

必需条件：
- 至少两个已接入集群
- 至少两个不同工作空间/项目范围
- 至少一个可复用的集群目标组
- 至少两个环境阶段，例如测试、预发、生产
- 至少一个可回滚的交付修订历史

适用范围：
- 来源接入与校验
- 目标组与环境推进
- 差异/漂移状态
- 同步、升级、暂停、恢复
- 交付历史与回滚
- 审计查询

## 4. 后端开发环境

依赖下载使用国内源：

```bash
export GOPROXY=https://goproxy.cn,direct
```

如需拉取 Git/Helm/OCI 联调镜像或本地控制器组件：
- 优先使用阿里云、DaoCloud 或已批准的国内代理镜像
- 禁止在未说明镜像源的情况下直接使用默认境外拉取

建议优先实现的后端模块：
- `gitops/sources`: 交付来源创建、校验、禁用与认证引用
- `gitops/targets`: 集群目标组与环境阶段建模
- `gitops/delivery`: 交付单元、覆盖合成、最终配置预览
- `gitops/status`: 状态聚合、差异视图、漂移判定、来源不可达与目标不可达降级
- `gitops/release`: 安装、升级、同步、暂停、恢复、回滚、卸载和环境推进
- `gitops/audit`: 发布动作与配置变更审计写入

建议新增配置项统一放到 `backend/config/config.dev.yaml` 与 `backend/config/config.example.yaml`：
- `gitops.sources.git_timeout`
- `gitops.sources.package_timeout`
- `gitops.sync.default_timeout`
- `gitops.sync.max_parallel_targets`
- `gitops.diff.cache_ttl`
- `gitops.release.history_retention_days`
- `gitops.audit.retention_days`

## 5. 前端开发环境

依赖下载使用国内源：

```bash
cd frontend
npm config set registry https://registry.npmmirror.com
```

建议优先实现的前端页面与模块：
- GitOps / 发布中心概览页
- 交付来源管理页
- 交付单元列表与详情页
- 差异/漂移视图面板
- 发布历史与回滚对话框
- 环境推进时间线与动作确认抽屉

首期交互原则：
- 所有发布动作都必须在交付单元上下文中发起，避免脱离目标范围执行
- 高影响动作必须展示目标环境、目标组、应用版本、配置版本和风险提示
- 暂停状态、漂移状态、部分成功和来源不可达必须有明确提示，不能只显示单一成功/失败标签

## 6. 推荐实施顺序

1. 落库交付来源、目标组、环境阶段、配置覆盖、发布修订和动作记录表。
2. 完成统一权限校验、交付单元建模和状态聚合读取链路。
3. 完成差异/漂移聚合和来源校验，形成最小可读闭环。
4. 完成同步、安装、升级、暂停、恢复和卸载等基础动作。
5. 完成多环境推进、部分成功归集和发布历史。
6. 完成回滚、审计、异常降级和契约/集成测试。

## 7. 最小验收清单

- 授权用户能够接入一个交付来源并完成来源校验。
- 能创建交付单元，并为其绑定目标组、环境阶段和配置覆盖。
- 能查看期望状态、实际状态、最近同步结果和漂移状态。
- 能查看交付差异，并定位差异环境与目标范围。
- 能执行同步、安装或升级动作，并看到进度与结果。
- 能暂停与恢复自动对齐，并在暂停期间继续看到待同步变更。
- 能查看交付修订历史，并对可恢复修订执行回滚。
- 未授权用户无法访问其他工作空间、项目或环境范围内的来源、交付单元和发布动作入口。

## 8. PR 交付要求

- 所有提交说明必须专业、中文、可审计。
- 推送远程后必须创建中文 PR，并附上数据库备份证据、验证结果和风险说明。
- 合并前必须取得用户明确同意。
