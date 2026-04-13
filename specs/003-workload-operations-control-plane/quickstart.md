# Quickstart: 多集群 Kubernetes 工作负载运维控制面

## 1. 前提

- 当前工作分支必须是 `003-workload-operations-control-plane` 或其后续专用实现分支，禁止在 `main` 上直接开发。
- `001` 和 `002` 已形成 003 的基础能力来源；003 继续遵守“中文 PR + 用户明确同意后再合并”的宪章要求。
- 任何数据库结构、终端会话审计表、批量任务表或动作模型变更前，都必须执行新的数据库备份。

## 2. 实施前数据库备份

在开始实现 003 的表结构、终端会话和动作执行链路前，执行以下命令：

```bash
mkdir -p artifacts/003-workload-operations-control-plane

docker exec mysql8 sh -lc 'mysqldump -hlocalhost -P3306 -uadmin -p123456 --all-databases --single-transaction --quick --routines --events --triggers' \
  > artifacts/003-workload-operations-control-plane/mysql-backup-$(date +%Y%m%d-%H%M%S).sql
```

如果当前容器环境中的 `admin/123456` 不可用，必须改用实际可用凭据执行，并记录差异与恢复验证结果。

恢复抽样验证建议：

```bash
docker run -d --rm --name mysql8-restore-check-003 -e MYSQL_ROOT_PASSWORD=123456 mysql:8.0
docker exec -i mysql8-restore-check-003 sh -lc 'MYSQL_PWD=123456 mysql -uroot' < artifacts/003-workload-operations-control-plane/<backup-file>.sql
docker exec mysql8-restore-check-003 sh -lc 'MYSQL_PWD=123456 mysql -N -B -uroot -e "SHOW DATABASES;"'
```

## 2026-04-12 备份执行记录

- 分支：`003-workload-operations-control-plane`
- 备份文件名：`artifacts/003-workload-operations-control-plane/mysql-backup-20260412-150123.sql`
- 凭据说明：`admin/123456` 在当前容器环境不可用；本次实际使用 `root/123456` 执行备份
- 恢复校验：已通过临时容器 `mysql8-restore-check-003` 完成完整导入，并核对数据库列表与 `kb_manage` 表数量一致
- 详细记录：`artifacts/003-workload-operations-control-plane/backup-manifest.txt`

## 2026-04-12 实施启动记录

- 已完成 Governance 证据文件：
  - `artifacts/003-workload-operations-control-plane/branch-check.txt`
  - `artifacts/003-workload-operations-control-plane/mirror-and-remote-check.txt`
- 已进入 `/speckit.implement` 执行，当前处于 Phase 0/1/2（治理+基础骨架）阶段。

## 3. 联调模式建议

### 模式 A：单资源运维诊断联调

适用于尽快验证单资源页、实例列表、日志跳转和终端入口。

必需条件：
- 至少一个已接入且已授权的集群
- 至少一个 `Deployment`、`StatefulSet` 或 `DaemonSet`
- 目标命名空间具备 Pod 与容器实例
- 平台可通过 Kubernetes API 访问 Pod 日志和 exec 通道

适用范围：
- 工作负载运维上下文页
- 实例列表
- 日志跳转联动
- 终端会话建立与关闭

不覆盖：
- 批量动作
- 回滚恢复
- 大规模权限回收验收

### 模式 B：003 完整验收联调

适用于按 003 规格做完整验收。

必需条件：
- 至少两个已接入集群
- 至少两个不同工作空间/项目范围
- 至少一个可产生 revision 历史的工作负载
- 至少一个可执行批量动作的资源集合

适用范围：
- 单资源诊断
- 单体动作
- 批量动作
- 发布历史与回滚
- 终端会话审计
- 权限回收即时生效

## 4. 后端开发环境

依赖下载使用国内源：

```bash
export GOPROXY=https://goproxy.cn,direct
```

如需拉取本地联调镜像：
- 优先使用阿里云、DaoCloud 或已批准的国内代理镜像
- 禁止在未说明镜像源的情况下直接使用默认境外拉取

建议优先实现的后端模块：
- `workloadops/context`: 单资源运维聚合视图、实例摘要、最近变更
- `workloadops/actions`: 动作提交、进度查询、失败归一化、回滚执行
- `workloadops/batch`: 批量任务编排、受控并发、子项结果汇总
- `workloadops/terminal`: 会话建立、关闭、超时回收、会话审计
- `workloadops/revisions`: 工作负载发布历史与可回滚目标识别

建议新增配置项统一放到 `backend/config/config.dev.yaml` 与 `backend/config/config.example.yaml`：
- `workloadops.actions.*`
- `workloadops.batch.*`
- `workloadops.terminal.*`
- `workloadops.audit.*`

## 5. 前端开发环境

依赖下载使用国内源：

```bash
cd frontend
npm config set registry https://registry.npmmirror.com
```

建议优先实现的前端页面与模块：
- 工作负载运维上下文页
- 实例列表与单实例诊断面板
- 发布历史与回滚对话框
- 批量动作提交与结果页
- 终端会话入口与会话状态提示

首期交互原则：
- 所有动作必须在资源上下文中发起，避免脱离目标对象执行
- 高风险动作必须展示影响范围和二次确认
- 终端关闭、超时和实例失效必须有明确提示

## 6. 推荐实施顺序

1. 落库终端会话、批量任务、动作扩展字段和必要审计索引。
2. 完成统一权限校验、工作负载上下文聚合和发布历史识别。
3. 完成单资源动作和进度查询，形成最小写路径闭环。
4. 完成终端会话管理和会话审计。
5. 完成批量动作、部分失败处理和结果归集。
6. 完成版本回滚、异常降级和契约/集成测试。

## 7. 最小验收清单

- 授权用户能够从目标工作负载进入统一运维上下文页。
- 能查看实例分布、最近变更、最近动作结果和发布进度。
- 能从资源上下文跳转查看目标实例日志。
- 能创建和关闭容器终端会话，并看到会话状态变化。
- 能执行扩缩容、重启或重新部署等单资源动作，并看到进度和结果。
- 能查看可用发布历史，并对存在历史版本的工作负载执行回滚。
- 能对多个工作负载提交同类批量动作，并看到总进度与单项结果。
- 未授权用户无法访问其他工作空间或项目范围内的工作负载、终端和高风险动作入口。

## 8. PR 交付要求

- 所有提交说明必须专业、中文、可审计。
- 推送远程后必须创建中文 PR，并附上数据库备份证据、验证结果和风险说明。
- 合并前必须取得用户明确同意。
