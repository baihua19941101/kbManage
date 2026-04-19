# Quickstart: 平台级备份恢复与灾备中心

## 当前执行状态（2026-04-18）

- 008 已完成 `specify` 与 `plan` 阶段，当前已生成 `spec.md`、`plan.md`、`research.md`、`data-model.md`、`contracts/openapi.yaml` 和质量清单。
- 当前工作分支为 `008-backup-restore-dr`。
- 后续进入 `/speckit.tasks` 或实现前，仍需按仓库治理要求落实数据库备份、国内依赖源配置、中文 PR 和用户明确同意后再合并。

## 1. 前提

- 仓库位于功能分支 `008-backup-restore-dr`
- 后端配置文件统一使用 `backend/config/config.dev.yaml`
- 前端配置文件统一使用 `frontend/.env.development`
- Go 依赖通过 `GOPROXY=https://goproxy.cn,direct`
- npm 依赖通过 `https://registry.npmmirror.com`

## 2. 实施前数据库备份

在开始实现 008 之前，先执行：

```bash
mkdir -p artifacts/008-backup-restore-dr
docker exec mysql8 sh -lc 'mysqldump -hlocalhost -P3306 -uadmin -p123456 --all-databases --single-transaction --quick --routines --events --triggers' \
  > artifacts/008-backup-restore-dr/mysql-backup-$(date +%Y%m%d-%H%M%S).sql
```

如果当前环境 `admin/123456` 不可用，必须：

1. 改用容器实际可用凭据执行备份
2. 在 `artifacts/008-backup-restore-dr/backup-manifest.txt` 中记录差异原因、实际命令和产物路径
3. 使用临时 MySQL 容器做恢复抽样验证

恢复抽样示例：

```bash
docker run -d --rm --name mysql8-restore-check-008 -e MYSQL_ROOT_PASSWORD=123456 mysql:8.0
cat artifacts/008-backup-restore-dr/mysql-backup-<timestamp>.sql \
  | docker exec -i mysql8-restore-check-008 sh -lc 'MYSQL_PWD=123456 mysql -h127.0.0.1 -uroot'
docker exec mysql8-restore-check-008 sh -lc 'MYSQL_PWD=123456 mysql -h127.0.0.1 -N -B -uroot -e "SHOW DATABASES;"'
```

## 3. 国内依赖源配置

### Go

```bash
go env -w GOPROXY=https://goproxy.cn,direct
```

### npm

```bash
cd frontend
npm config set registry https://registry.npmmirror.com
```

### 其他联调依赖

- 若需要对象存储 SDK、备份执行器、快照工具或演练联调镜像，优先使用阿里云、DaoCloud 或已批准的国内代理。
- 禁止在计划和任务中省略镜像源说明。

## 4. 最小联调路径

完成实现后，最小联调至少覆盖：

1. 为平台元数据和关键业务命名空间各创建一条备份策略
2. 触发一次手动备份并确认生成恢复点、耗时、结果和一致性说明
3. 基于已有恢复点执行一次原地恢复预检查并查看冲突提示
4. 发起一次跨集群恢复或环境迁移并确认源目标映射关系可见
5. 创建一份灾备演练计划，写入 RPO/RTO、切换步骤和验证清单
6. 完成一次演练记录并生成演练报告
7. 检索至少一条备份恢复域审计记录

## 5. 完整验收清单

- 备份策略、恢复点、恢复任务、迁移任务和演练计划均有明确状态流
- 恢复和迁移前能区分冲突阻断、一致性说明和人工确认步骤
- 原地恢复、跨集群恢复、定向恢复和环境迁移在界面与审计上具备清晰区分
- 演练记录能够对比目标 RPO/RTO 与实际结果
- 报告包含目标达成情况、问题项和改进建议
- 权限回收后备份、恢复、迁移和演练动作立即锁定
- README、配置说明和端口配置方式在实现阶段同步更新

## 6. 规划产物

- 规格：`specs/008-backup-restore-dr/spec.md`
- 计划：`specs/008-backup-restore-dr/plan.md`
- 研究：`specs/008-backup-restore-dr/research.md`
- 数据模型：`specs/008-backup-restore-dr/data-model.md`
- 契约：`specs/008-backup-restore-dr/contracts/openapi.yaml`
