# Quickstart: 多集群 Kubernetes 安全与策略治理中心

## 当前执行状态（2026-04-14）

- 005 已完成 `specify/plan/tasks`，当前处于 implement 阶段。
- 2026-04-14 已完成一次实施前数据库备份与恢复抽样验证，详见 `artifacts/005-security-and-policy/backup-manifest.txt`。
- 治理证据文件已补齐：`branch-check.txt`、`mirror-and-remote-check.txt`。

## 1. 前提

- 当前工作分支必须是 `005-security-and-policy`（或其后续实现分支），禁止在 `main/master` 直接开发。
- 001-004 已可用，尤其需要 001 的授权与审计底座。
- 在进入编码前，必须完成数据库备份并记录证据。

## 2. 实施前数据库备份（阻断门槛）

先创建证据目录：

```bash
mkdir -p artifacts/005-security-and-policy
```

标准备份命令：

```bash
docker exec mysql8 sh -lc 'mysqldump -hlocalhost -P3306 -uadmin -p123456 --all-databases --single-transaction --quick --routines --events --triggers' \
  > artifacts/005-security-and-policy/mysql-backup-$(date +%Y%m%d-%H%M%S).sql
```

恢复抽样验证（建议）：

```bash
docker run -d --rm --name mysql8-restore-check-005 -e MYSQL_ROOT_PASSWORD=123456 mysql:8.0
cat artifacts/005-security-and-policy/mysql-backup-<timestamp>.sql \
  | docker exec -i mysql8-restore-check-005 sh -lc 'MYSQL_PWD=123456 mysql -uroot'
docker exec mysql8-restore-check-005 sh -lc 'MYSQL_PWD=123456 mysql -N -B -uroot -e "SHOW DATABASES;"'
```

若环境中 `admin/123456` 不可用，必须在 `artifacts/005-security-and-policy/backup-manifest.txt` 记录：
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

## 4. 最小联调路径（MVP）

目标：先打通 P1 范围（策略中心 + 准入模式 + 灰度 + 例外）。

1. 创建一条平台级策略（`pod-security` 或 `image` 类）。
2. 为一个试点命名空间创建 `audit/warn` 分配，观察命中数据。
3. 将该分配切换到 `enforce`，验证违规对象被拦截。
4. 对阻断对象发起例外申请并审批通过，验证限时放行。
5. 例外到期后验证自动恢复原策略约束。
6. 在审计页检索策略创建、模式切换、例外审批、整改更新全链路事件。

## 5. 完整验收路径

适用于 005 全量验收。

- 至少 20 个已接入集群中的样本范围可见。
- 至少 2 个工作空间 + 2 个项目，用于隔离验证。
- 覆盖策略层级：平台级、工作空间级、项目级。
- 覆盖执行模式：`audit`、`alert`、`warn`、`enforce`。
- 覆盖例外状态：`pending`、`approved/active`、`rejected`、`expired`、`revoked`。
- 覆盖整改状态：`open` -> `in_progress` -> `closed`。

## 6. 推荐实施顺序

1. 落库 005 核心表：策略、版本、分配、命中、例外、整改、审计索引。
2. 后端策略服务与分配服务，先打通读写和权限校验。
3. 实现命中查询、模式切换任务和分阶段分配。
4. 实现例外申请与审批、有效期处理。
5. 实现整改状态更新和治理看板查询。
6. 补齐审计事件写入、契约测试、集成测试和前端门禁。

## 7. 最小验收清单

- 授权管理员可创建和更新策略。
- 可按集群/命名空间/项目/资源类型分配策略。
- 可在 `audit/warn/enforce` 之间切换并看到命中差异。
- 可查询违规对象、风险级别和整改状态。
- 可提交并审批例外，且到期后自动失效。
- 权限回收后，相关策略操作立即被拦截。
- 审计人员可检索策略变更与违规处置完整记录。

## 8. PR 交付要求

- 提交信息必须专业、中文、可审计。
- PR 说明必须包含：备份证据、测试结果、风险说明、未完成项。
- 未获用户明确同意前，禁止执行合并。
