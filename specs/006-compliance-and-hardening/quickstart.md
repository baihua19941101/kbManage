# Quickstart: 多集群 Kubernetes 合规与加固中心

## 当前执行状态（2026-04-15）

- 006 已完成 `specify/plan/tasks`，当前处于任务清单就绪、待实现阶段。
- 已确认 `005-security-and-policy` 的主 PR 流已完成并合并到 `main`，006 不再受前序 feature gate 阻塞。
- 进入实现前仍需补齐数据库备份与恢复抽样验证证据。

## 1. 前提

- 当前工作分支必须是 `006-compliance-and-hardening`（或其后续实现分支），禁止在 `main/master` 直接开发。
- 001-005 已可用，尤其需要 001 的授权与审计底座、005 的策略治理语义参考。
- 进入编码前，必须先完成数据库备份并记录证据。
- 需要具备至少一个可访问的扫描执行来源或模拟适配层，用于验证计划性/按需扫描链路。

## 2. 实施前数据库备份（阻断门槛）

先创建证据目录：

```bash
mkdir -p artifacts/006-compliance-and-hardening
```

标准备份命令：

```bash
docker exec mysql8 sh -lc 'mysqldump -hlocalhost -P3306 -uadmin -p123456 --all-databases --single-transaction --quick --routines --events --triggers' \
  > artifacts/006-compliance-and-hardening/mysql-backup-$(date +%Y%m%d-%H%M%S).sql
```

恢复抽样验证（建议）：

```bash
docker run -d --rm --name mysql8-restore-check-006 -e MYSQL_ROOT_PASSWORD=123456 mysql:8.0
cat artifacts/006-compliance-and-hardening/mysql-backup-<timestamp>.sql \
  | docker exec -i mysql8-restore-check-006 sh -lc 'MYSQL_PWD=123456 mysql -uroot'
docker exec mysql8-restore-check-006 sh -lc 'MYSQL_PWD=123456 mysql -N -B -uroot -e "SHOW DATABASES;"'
```

若环境中 `admin/123456` 不可用，必须在 `artifacts/006-compliance-and-hardening/backup-manifest.txt` 记录：
- 实际使用凭据与命令
- 差异原因
- 恢复验证结果

## 3. 国内源配置

后端依赖：

```bash
export GOPROXY=https://goproxy.cn,direct
```

前端依赖：

```bash
cd frontend
npm config set registry https://registry.npmmirror.com
```

若引入扫描器镜像、基线规则包或联调组件，优先使用阿里云、DaoCloud 或经批准的国内代理。

## 4. 最小联调路径（MVP）

目标：先打通 P1 范围（基线选择 + 按需扫描 + 失败项整改/例外/复检闭环）。

1. 创建一条 `CIS` 或平台基线模板。
2. 为一个试点集群或命名空间创建扫描配置。
3. 手动触发一次按需扫描并查看合规得分、失败项数量和风险等级分布。
4. 打开一个高风险失败项，查看证据摘要和受影响对象。
5. 将该失败项流转为整改任务，并记录责任人和计划完成时间。
6. 对另一个暂时无法处理的失败项提交例外申请并完成审批。
7. 在一个已整改对象上发起复检，验证结果从失败转为通过或继续失败。
8. 在审计页检索扫描发起、整改更新、例外审批和复检完成全链路事件。
9. 生成一次合规归档导出任务，验证可按筛选条件打包扫描、趋势和审计结果。

## 5. 完整验收路径

适用于 006 全量验收。

- 至少 20 个已接入集群中的样本范围可见。
- 至少 2 个工作空间 + 2 个项目，用于隔离验证。
- 覆盖基线类型：`CIS`、`STIG`、平台基线模板。
- 覆盖扫描模式：按需扫描、计划性扫描、复检触发。
- 覆盖结果状态：`succeeded`、`partially_succeeded`、`failed`、`partial coverage`。
- 覆盖例外状态：`pending`、`approved/active`、`rejected`、`expired`、`revoked`。
- 覆盖整改状态：`todo` -> `in_progress` -> `done`。
- 覆盖复检状态：`pending` -> `running` -> `passed/failed`。
- 覆盖归档导出状态：`pending` -> `running` -> `succeeded/failed`。

## 6. 推荐实施顺序

1. 落库 006 核心表：基线、扫描配置、扫描执行、失败项、证据索引、整改任务、例外、复检、趋势快照、审计索引。
2. 后端基线与扫描配置服务，先打通读写、授权与计划/手动触发入口。
3. 实现扫描执行 worker、结果归集与证据快照模型，覆盖部分成功与失败原因归一化。
4. 实现失败项检索、整改任务、例外审批与复检服务。
5. 实现趋势快照聚合、覆盖率汇总与管理视图查询。
6. 实现归档导出任务、产物状态查询和审计留痕。
7. 补齐审计事件写入、契约测试、集成测试和前端门禁。

## 7. 最小验收清单

- 授权管理员可创建和更新基线标准与扫描配置。
- 可对集群、节点、命名空间和关键资源范围发起按需或计划性扫描。
- 可查看扫描得分、失败项、风险分布和证据详情。
- 可对失败项创建整改任务、提交并审批例外、发起复检。
- 例外到期后可自动恢复该项为待处理风险。
- 可创建并查询合规归档导出任务。
- 权限回收后，相关扫描结果、证据和治理动作立即被拦截。
- 审计人员可检索扫描、整改、例外、复检完整记录。
- 管理视图可按集群或团队展示覆盖率、遗留风险和整改进度。

## 8. PR 交付要求

- 提交信息必须专业、中文、可审计。
- PR 说明必须包含：备份证据、测试结果、风险说明、未完成项。
- 未获用户明确同意前，禁止执行合并。
